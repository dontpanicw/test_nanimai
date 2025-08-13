package service

import "time"

type AccountDTO struct {
	UserID        int64
	CurrentAmount int64
	MaxAmount     int64
}

type ReservationDTO struct {
	ID             int64
	AccountID      int64
	OwnerServiceID int64
	Amount         int64
	Status         string
	ExpiresAt      time.Time
}

type UpdateBalanceInput struct {
	AccountID int64
	Delta     int64
}

type UpdateLimitInput struct {
	AccountID int64
	Delta     int64
}

type OpenReservationInput struct {
	AccountID      int64
	OwnerServiceID int64
	Amount         int64
	IdempotencyKey string
	Timeout        time.Duration
}
