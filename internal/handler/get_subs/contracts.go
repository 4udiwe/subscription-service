package get_subs

import (
	"context"

	"github.com/4udiwe/subscription-service/internal/entity"
)

type SubscriptionService interface {
	GetAllSubscriptions(ctx context.Context, page int, pageSize int) ([]entity.SubscriptionFullInfo, int, error)
}
