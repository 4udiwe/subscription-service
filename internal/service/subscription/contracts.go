package subscription

import (
	"context"
	"time"

	"github.com/4udiwe/subscription-serivce/internal/entity"
	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, userID, offerID uuid.UUID, startDate, endDate time.Time) (entity.Subscription, error)
	GetAll(ctx context.Context) ([]entity.Subscription, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Subscription, error)
	//GetById(ctx context.Context, id uuid.UUID) (entity.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetAllByUserIDAndSubscriptionName(
		ctx context.Context,
		userID uuid.UUID,
		subscriptionName string,
		startPeriod *time.Time,
		endPeriod *time.Time,
	) ([]entity.Subscription, error)
}

type OfferRepository interface {
	Create(ctx context.Context, name string, price int, durationMonth int) (entity.Offer, error)
	GetByID(ctx context.Context, id uuid.UUID) (entity.Offer, error)
	GetByNameAndPrice(ctx context.Context, name string, price int) (entity.Offer, error)
}
