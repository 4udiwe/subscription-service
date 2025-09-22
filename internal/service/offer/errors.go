package offer

import "errors"

var (
	ErrOfferWithNameAndPriceAlreadyExists = errors.New("offer with given name and price already exists")
	ErrActiveSubscriptionsExist           = errors.New("active subscriptions exist for given offer, could not delete")
	ErrOfferNotFound                      = errors.New("offer not found")
)
