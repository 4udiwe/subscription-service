package app

import (
	"github.com/4udiwe/subscription-service/internal/service/offer"
	"github.com/4udiwe/subscription-service/internal/service/subscription"
)

func (app *App) OfferService() *offer.OfferService {
	if app.offerService != nil {
		return app.offerService
	}
	app.offerService = offer.New(app.OfferRepo(), app.SubscriptionRepo(), app.Postgres())
	return app.offerService
}

func (app *App) SubscriptionService() *subscription.SubscriptionService {
	if app.subService != nil {
		return app.subService
	}
	app.subService = subscription.New(app.SubscriptionRepo(), app.OfferRepo(), app.Postgres())
	return app.subService
}
