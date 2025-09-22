package post_sub_by_name

import (
	"context"
	"time"

	"github.com/4udiwe/subscription-service/internal/entity"
	"github.com/google/uuid"
)

type SubscriptionService interface {
	CreateSubscription(
		ctx context.Context,
		userID uuid.UUID,
		serviceName string,
		price int,
		startDate time.Time,
		endDate *time.Time,
	) (entity.SubscriptionFullInfo, error)
}
