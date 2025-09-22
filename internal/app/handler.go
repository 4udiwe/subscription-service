package app

import (
	"github.com/4udiwe/subscription-service/internal/handler"
	"github.com/4udiwe/subscription-service/internal/handler/delete_sub"
	"github.com/4udiwe/subscription-service/internal/handler/get_offers"
	"github.com/4udiwe/subscription-service/internal/handler/get_subs"
	"github.com/4udiwe/subscription-service/internal/handler/get_subs_by_user"
	"github.com/4udiwe/subscription-service/internal/handler/get_subs_by_user_subname"
	"github.com/4udiwe/subscription-service/internal/handler/post_sub_by_name"
	"github.com/4udiwe/subscription-service/internal/handler/post_sub_by_offer_id"
)

func (app *App) DeleteProductHandler() handler.Handler {
	if app.deleteSubscriptionHandler != nil {
		return app.deleteSubscriptionHandler
	}
	app.deleteSubscriptionHandler = delete_sub.New(app.SubscriptionService())
	return app.deleteSubscriptionHandler
}

func (app *App) GetOffersHandler() handler.Handler {
	if app.getOffersHandler != nil {
		return app.getOffersHandler
	}
	app.getOffersHandler = get_offers.New(app.OfferService())
	return app.getOffersHandler
}

func (app *App) GetSubscriptionsHandler() handler.Handler {
	if app.getSubscriptionsHandler != nil {
		return app.getSubscriptionsHandler
	}
	app.getSubscriptionsHandler = get_subs.New(app.SubscriptionService())
	return app.getSubscriptionsHandler
}

func (app *App) GetSubscriptionsByUserHandler() handler.Handler {
	if app.getSubscriptionsByUserHandler != nil {
		return app.getSubscriptionsByUserHandler
	}
	app.getSubscriptionsByUserHandler = get_subs_by_user.New(app.SubscriptionService())
	return app.getSubscriptionsByUserHandler
}

func (app *App) GetSubscriptionsByUserAndSubNameHandler() handler.Handler {
	if app.getSubscriptionsByUserAndSubNameHandler != nil {
		return app.getSubscriptionsByUserAndSubNameHandler
	}
	app.getSubscriptionsByUserAndSubNameHandler = get_subs_by_user_subname.New(app.SubscriptionService())
	return app.getSubscriptionsByUserAndSubNameHandler
}

func (app *App) PostSubciptionByNameHandler() handler.Handler {
	if app.postSubciptionByNameHandler != nil {
		return app.postSubciptionByNameHandler
	}
	app.postSubciptionByNameHandler = post_sub_by_name.New(app.SubscriptionService())
	return app.postSubciptionByNameHandler
}

func (app *App) PostSubciptionByOfferIDHandler() handler.Handler {
	if app.postSubciptionByOfferIDHandler != nil {
		return app.postSubciptionByOfferIDHandler
	}
	app.postSubciptionByOfferIDHandler = post_sub_by_offer_id.New(app.SubscriptionService())
	return app.postSubciptionByOfferIDHandler
}
