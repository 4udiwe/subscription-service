package offer

import "errors"

var (
	ErrOfferNotFound = errors.New("offer not found")

	ErrCannotCreateOffer = errors.New("cannot create offer")
	ErrCannotFindOffer   = errors.New("cannot find offer")
	ErrCannotDeleteOffer = errors.New("cannot delete offer")
	ErrCannotFetchOffers = errors.New("cannot fetch offers")

	ErrCannotCheckActiveSubscriptions     = errors.New("cannot check active subscriptions for offer")
	ErrOfferWithNameAndPriceAlreadyExists = errors.New("offer with given name and price already exists")
	ErrActiveSubscriptionsExist           = errors.New("active subscriptions exist for given offer, could not delete")
)
