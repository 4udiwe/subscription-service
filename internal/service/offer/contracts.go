package offer

import (
	"context"

	"github.com/4udiwe/subscription-service/internal/entity"
	"github.com/google/uuid"
)

type OfferRepository interface {
	Create(ctx context.Context, name string, price int, durationMonths int) (entity.Offer, error)
	GetAll(ctx context.Context, limit int, offset int) (offers []entity.Offer, total int, err error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type SubscriptionRepository interface {
	GetAllByOfferID(ctx context.Context, offerID uuid.UUID) ([]entity.Subscription, error)
}
