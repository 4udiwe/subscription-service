package get_subs_by_user_subname

import (
	"context"
	"time"

	"github.com/4udiwe/subscription-service/internal/entity"
	"github.com/google/uuid"
)

type SubscriptionService interface {
	GetAllWithPriceByUserIDAndSubscriptionName(
		ctx context.Context,
		userID uuid.UUID,
		subscriptionName string,
		startPeriod *time.Time,
		endPeriod *time.Time,
	) (subs []entity.SubscriptionFullInfo, price int, err error)
}
