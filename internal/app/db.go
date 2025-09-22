package app

import (
	offer_repo "github.com/4udiwe/subscription-service/internal/repository/offer"
	subscription_repo "github.com/4udiwe/subscription-service/internal/repository/subscription"
	"github.com/4udiwe/subscription-service/pkg/postgres"
)

func (app *App) Postgres() *postgres.Postgres {
	return app.postgres
}

func (app *App) OfferRepo() *offer_repo.Repository {
	if app.offerRepo != nil {
		return app.offerRepo
	}
	app.offerRepo = offer_repo.New(app.Postgres())
	return app.offerRepo
}

func (app *App) SubscriptionRepo() *subscription_repo.Repository {
	if app.subRepo != nil {
		return app.subRepo
	}
	app.subRepo = subscription_repo.New(app.Postgres())
	return app.subRepo
}
