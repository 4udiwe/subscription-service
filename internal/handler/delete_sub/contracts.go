package delete_sub

import (
	"context"

	"github.com/google/uuid"
)

type SubscriptionService interface {
	DeleteSubscription(ctx context.Context, subID uuid.UUID) error
}
