package postgres

import (
	"context"
	"database/sql"
	"errors"
	"test_nanimai/backend/domain"
	"time"
)

var (
	ErrNotEnoughFunds = errors.New("not enough funds")
	ErrNotFound       = errors.New("not found")
	ErrForbidden      = errors.New("forbidden")
	ErrExpired        = errors.New("reservation expired")
)

type BalanceStorage struct {
	db *sql.DB
}

func (bs *BalanceStorage) GetDb() *sql.DB {
	return bs.db
}

func NewBalanceStorage(connStr string) (*BalanceStorage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			db.Close()
		}
	}()
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &BalanceStorage{db: db}, nil
}

func (s *BalanceStorage) UpdateLimit(ctx context.Context, accountID int64, delta int64) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE accounts
		SET max_amount = max_amount + $1
		WHERE id = $2
	`, delta, accountID)
	return err
}

func (s *BalanceStorage) UpdateBalance(ctx context.Context, accountID int64, delta int64) error {
	// Запрет уйти ниже 0 и выше max_amount
	cmd, err := s.db.ExecContext(ctx, `
		UPDATE accounts
		SET current_amount = current_amount + $1
		WHERE id = $2
		  AND (current_amount + $1) >= 0
		  AND (current_amount + $1) <= max_amount
	`, delta, accountID)
	if err != nil {
		return err
	}
	rows, _ := cmd.RowsAffected()
	if rows == 0 {
		return ErrNotEnoughFunds
	}
	return nil
}

func (r *BalanceStorage) OpenReservation(ctx context.Context, ownerServiceID, accountID int64, amount int64, idempotencyKey string, timeout time.Duration) (*domain.Reservation, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var existing domain.Reservation
	err = tx.QueryRowContext(ctx, `
		SELECT id, account_id, owner_service_id, amount, status, idempotency_key, expires_at, created_at
		FROM reservations
		WHERE owner_service_id = $1 AND idempotency_key = $2
	`, ownerServiceID, idempotencyKey).Scan(
		&existing.ID, &existing.AccountID, &existing.OwnerServiceID, &existing.Amount,
		&existing.Status, &existing.IdempotencyKey, &existing.ExpiresAt, &existing.CreatedAt,
	)
	if err == nil {
		// Уже есть такая транзакция
		return &existing, nil
	}

	// Блокируем аккаунт
	var acc domain.Account
	err = tx.QueryRowContext(ctx, `
		SELECT id, user_id, current_amount, max_amount, reserved_amount
		FROM accounts
		WHERE id = $1
		FOR UPDATE
	`, accountID).Scan(
		&acc.ID, &acc.UserID, &acc.CurrentAmount, &acc.MaxAmount, &acc.ReservedAmount,
	)
	if err != nil {
		return nil, ErrNotFound
	}

	if (acc.CurrentAmount - acc.ReservedAmount) < amount {
		return nil, ErrNotEnoughFunds
	}

	// Создаём резерв
	var res domain.Reservation
	err = tx.QueryRowContext(ctx, `
		INSERT INTO reservations (account_id, owner_service_id, amount, status, idempotency_key, expires_at)
		VALUES ($1, $2, $3, 'ACTIVE', $4, now() + $5::interval)
		RETURNING id, account_id, owner_service_id, amount, status, idempotency_key, expires_at, created_at
	`, accountID, ownerServiceID, amount, idempotencyKey, timeout.String()).Scan(
		&res.ID, &res.AccountID, &res.OwnerServiceID, &res.Amount,
		&res.Status, &res.IdempotencyKey, &res.ExpiresAt, &res.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Увеличиваем reserved_amount
	_, err = tx.ExecContext(ctx, `
		UPDATE accounts
		SET reserved_amount = reserved_amount + $1
		WHERE id = $2
	`, amount, accountID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &res, nil
}

// 4. Подтверждение транзакции
func (r *BalanceStorage) ConfirmReservation(ctx context.Context, reservationID int64, ownerServiceID int64) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var accID int64
	var amount int64
	var status string
	var expiresAt time.Time
	err = tx.QueryRowContext(ctx, `
		SELECT account_id, amount, status, expires_at
		FROM reservations
		WHERE id = $1 AND owner_service_id = $2
		FOR UPDATE
	`, reservationID, ownerServiceID).Scan(&accID, &amount, &status, &expiresAt)
	if err != nil {
		return ErrNotFound
	}

	if status != "ACTIVE" {
		return errors.New("reservation not active")
	}
	if time.Now().After(expiresAt) {
		return ErrExpired
	}

	// Списываем средства и уменьшаем reserved_amount
	_, err = tx.ExecContext(ctx, `
		UPDATE accounts
		SET current_amount = current_amount - $1,
		    reserved_amount = reserved_amount - $1
		WHERE id = $2
	`, amount, accID)
	if err != nil {
		return err
	}

	// Меняем статус резерва
	_, err = tx.ExecContext(ctx, `
		UPDATE reservations
		SET status = 'CONFIRMED'
		WHERE id = $1
	`, reservationID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// 5. Отмена транзакции
func (r *BalanceStorage) CancelReservation(ctx context.Context, reservationID int64, ownerServiceID int64) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var accID int64
	var amount int64
	var status string
	err = tx.QueryRowContext(ctx, `
		SELECT account_id, amount, status
		FROM reservations
		WHERE id = $1 AND owner_service_id = $2
		FOR UPDATE
	`, reservationID, ownerServiceID).Scan(&accID, &amount, &status)
	if err != nil {
		return ErrNotFound
	}

	if status != "ACTIVE" {
		return errors.New("reservation not active")
	}

	// Возвращаем средства (уменьшаем reserved_amount)
	_, err = tx.ExecContext(ctx, `
		UPDATE accounts
		SET reserved_amount = reserved_amount - $1
		WHERE id = $2
	`, amount, accID)
	if err != nil {
		return err
	}

	// Меняем статус резерва
	_, err = tx.ExecContext(ctx, `
		UPDATE reservations
		SET status = 'CANCELLED'
		WHERE id = $1
	`, reservationID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
