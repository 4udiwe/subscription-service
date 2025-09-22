package offer

import (
	"context"

	"github.com/4udiwe/subscription-service/internal/entity"
	"github.com/google/uuid"
)

type OfferRepository interface {
	GetAll(ctx context.Context) ([]entity.Offer, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type SubscriptionRepository interface {
	GetAllByOfferID(ctx context.Context, offerID uuid.UUID) ([]entity.Subscription, error)
}
