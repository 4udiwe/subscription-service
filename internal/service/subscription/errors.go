package subscription

import "errors"

var (
	ErrOfferNotFound        = errors.New("offer not found")
	ErrSubscriptionNotFound = errors.New("subscription not found")
)
