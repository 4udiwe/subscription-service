package entity

import (
	"time"

	"github.com/google/uuid"
)

type Offer struct {
	ID            uuid.UUID `db:"id"`
	Name          string    `db:"name"`
	Price         int       `db:"price"`
	DurationMonth int       `db:"duration_month"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
