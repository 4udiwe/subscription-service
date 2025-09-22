package offer

import "errors"

var (
	ErrActiveSubscriptionsExist = errors.New("active subscriptions exist for given offer, could not delete")
	ErrOfferNotFound            = errors.New("offer not found")
)
