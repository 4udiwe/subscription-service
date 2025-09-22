package subscription

import (
	"context"
	"time"

	"github.com/4udiwe/subscription-service/internal/entity"
	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, userID, offerID uuid.UUID, startDate, endDate time.Time) (entity.Subscription, error)
	GetAll(ctx context.Context) ([]entity.SubscriptionFullInfo, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]entity.SubscriptionFullInfo, error)
	// GetById(ctx context.Context, id uuid.UUID) (entity.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetAllByUserIDAndSubscriptionName(
		ctx context.Context,
		userID uuid.UUID,
		subscriptionName string,
		startPeriod *time.Time,
		endPeriod *time.Time,
	) ([]entity.SubscriptionFullInfo, error)
}

type OfferRepository interface {
	Create(ctx context.Context, name string, price int, durationMonths int) (entity.Offer, error)
	GetByID(ctx context.Context, id uuid.UUID) (entity.Offer, error)
	GetByNameAndPrice(ctx context.Context, name string, price int) (entity.Offer, error)
}
