package domain

import "time"

type Account struct {
	ID             int64
	UserID         int64
	CurrentAmount  int64
	MaxAmount      int64
	ReservedAmount int64
}

type Reservation struct {
	ID             int64
	AccountID      int64
	OwnerServiceID int64
	Amount         int64
	Status         string
	IdempotencyKey string
	ExpiresAt      time.Time
	CreatedAt      time.Time
}
