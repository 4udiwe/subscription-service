package get_subs_by_user

import (
	"context"

	"github.com/4udiwe/subscription-service/internal/entity"
	"github.com/google/uuid"
)

type SubscriptionService interface {
	GetAllSubscriptionsByUserID(ctx context.Context, userID uuid.UUID, page int, pageSize int) ([]entity.SubscriptionFullInfo, int, error)
}
