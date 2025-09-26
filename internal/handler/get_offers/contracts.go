package get_offers

import (
	"context"

	"github.com/4udiwe/subscription-service/internal/entity"
)

type OfferService interface {
	GetAllOffers(ctx context.Context, page int, pageSize int) (offers []entity.Offer, total int, err error)
}
