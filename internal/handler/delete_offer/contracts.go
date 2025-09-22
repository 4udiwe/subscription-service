package delete_offer

import (
	"context"

	"github.com/google/uuid"
)

type OfferService interface {
	DeleteOffer(ctx context.Context, offerID uuid.UUID) error
}
