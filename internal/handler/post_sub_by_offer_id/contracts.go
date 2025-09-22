package post_sub_by_offer_id

import (
	"context"
	"time"

	"github.com/4udiwe/subscription-service/internal/entity"
	"github.com/google/uuid"
)

type SubscriptionService interface {
	CreateSubscriptionByOfferID(ctx context.Context, userID, offerID uuid.UUID, startDate time.Time) (entity.SubscriptionFullInfo, error)
}
