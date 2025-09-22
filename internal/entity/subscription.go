package entity

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	OfferID   uuid.UUID `db:"offer_id"`
	StartDate time.Time `db:"start_date"`
	EndDate   time.Time `db:"end_date"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type SubscriptionFullInfo struct {
	Subscription
	OfferName string `db:"offer_name"`
	Price     int    `db:"price"`
}
