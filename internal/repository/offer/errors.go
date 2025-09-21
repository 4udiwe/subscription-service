package offer_repo

import "errors"

var (
	ErrOfferNotFound                      = errors.New("offer not found")
	ErrOfferWithNameAndPriceAlreadyExists = errors.New("offer with the same name and price already exists")
)
