package service

import (
	"context"
	"test_nanimai/backend/domain"
	"time"
)

type Balance interface {
	UpdateLimit(ctx context.Context, accountID int64, delta int64) error
	UpdateBalance(ctx context.Context, accountID int64, delta int64) error
	OpenReservation(ctx context.Context, ownerServiceID, accountID int64, amount int64, idempotencyKey string, timeout time.Duration) (*domain.Reservation, error)
	ConfirmReservation(ctx context.Context, reservationID int64, ownerServiceID int64) error
	CancelReservation(ctx context.Context, reservationID int64, ownerServiceID int64) error
}
