package balance

import (
	"context"
	"time"

	"test_nanimai/backend/domain"
	"test_nanimai/backend/internal/repository"
)

type BalanceService struct {
	balanceRepo repository.Balance
}

func NewBalanceService(balanceRepo repository.Balance) *BalanceService {
	return &BalanceService{balanceRepo: balanceRepo}
}

func (s *BalanceService) UpdateLimit(ctx context.Context, accountID int64, delta int64) error {
	return s.balanceRepo.UpdateLimit(ctx, accountID, delta)
}

func (s *BalanceService) UpdateBalance(ctx context.Context, accountID int64, delta int64) error {
	return s.balanceRepo.UpdateBalance(ctx, accountID, delta)
}

func (s *BalanceService) OpenReservation(ctx context.Context, ownerServiceID, accountID int64, amount int64, idempotencyKey string, timeout time.Duration) (*domain.Reservation, error) {
	return s.balanceRepo.OpenReservation(ctx, ownerServiceID, accountID, amount, idempotencyKey, timeout)
}

func (s *BalanceService) ConfirmReservation(ctx context.Context, reservationID, ownerServiceID int64) error {
	return s.balanceRepo.ConfirmReservation(ctx, reservationID, ownerServiceID)
}

func (s *BalanceService) CancelReservation(ctx context.Context, reservationID, ownerServiceID int64) error {
	return s.balanceRepo.CancelReservation(ctx, reservationID, ownerServiceID)
}
