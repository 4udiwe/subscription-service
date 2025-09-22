package post_offer

import (
	"context"

	"github.com/4udiwe/subscription-service/internal/entity"
)

type OfferService interface {
	CreateOffer(ctx context.Context, name string, price int, durationMonths int) (entity.Offer, error)
}
