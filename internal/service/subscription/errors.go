package subscription

import "errors"

var (
	ErrOfferNotFound     = errors.New("offer not found")
	ErrCannotFindOffer   = errors.New("cannot find offer")
	ErrCannotCreateOffer = errors.New("cannot create offer")

	ErrSubscriptionNotFound     = errors.New("subscription not found")
	ErrCannotFindSubscription   = errors.New("cannot find subscription")
	ErrCannotCreateSubscription = errors.New("cannot create subscription")
	ErrCannotFetchSubscriptions = errors.New("cannot fetch subscriptions")
	ErrCannotDeleteSubscription = errors.New("cannot delete subscription")

	ErrUserAlreadyHasActiveSubscription = errors.New("user already has an active subscription for the given offer and date")
	ErrCannotCheckActiveSubscription    = errors.New("cannot check active subscription")
)
