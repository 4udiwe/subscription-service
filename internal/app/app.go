package app

import (
	"context"
	"os"

	"github.com/4udiwe/subscription-service/config"
	"github.com/4udiwe/subscription-service/internal/database"
	"github.com/4udiwe/subscription-service/internal/handler"
	offer_repo "github.com/4udiwe/subscription-service/internal/repository/offer"
	subscription_repo "github.com/4udiwe/subscription-service/internal/repository/subscription"
	"github.com/4udiwe/subscription-service/internal/service/offer"
	"github.com/4udiwe/subscription-service/internal/service/subscription"
	"github.com/4udiwe/subscription-service/pkg/httpserver"
	"github.com/4udiwe/subscription-service/pkg/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type App struct {
	cfg       *config.Config
	interrupt <-chan os.Signal

	// DB
	postgres *postgres.Postgres

	// Echo
	echoHandler *echo.Echo

	// Repositories
	offerRepo *offer_repo.Repository
	subRepo   *subscription_repo.Repository

	// Services
	offerService *offer.OfferService
	subService   *subscription.SubscriptionService

	// Handlers
	deleteSubscriptionHandler handler.Handler

	getOffersHandler                        handler.Handler
	getSubscriptionsHandler                 handler.Handler
	getSubscriptionsByUserHandler           handler.Handler
	getSubscriptionsByUserAndSubNameHandler handler.Handler

	postSubciptionByNameHandler    handler.Handler
	postSubciptionByOfferIDHandler handler.Handler
}

func New(configPath string) *App {
	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatalf("app - New - config.New: %v", err)
	}

	initLogger(cfg.Log.Level)

	return &App{
		cfg: cfg,
	}
}

func (app *App) Start() {
	// Postgres
	log.Info("Connecting to PostgreSQL...")

	postgres, err := postgres.New(app.cfg.Postgres.URL, postgres.ConnAttempts(5))

	if err != nil {
		log.Fatalf("app - Start - Postgres failed:%v", err)
	}
	app.postgres = postgres

	defer postgres.Close()

	// Migrations
	if err := database.RunMigrations(context.Background(), app.postgres.Pool); err != nil {
		log.Errorf("app - Start - Migrations failed: %v", err)
	}

	// App server
	log.Info("Starting app server...")
	httpServer := httpserver.New(app.EchoHandler(), httpserver.Port(app.cfg.HTTP.Port))
	httpServer.Start()
	log.Debugf("Server port: %s", app.cfg.HTTP.Port)

	defer func() {
		if err := httpServer.Shutdown(); err != nil {
			log.Errorf("HTTP server shutdown error: %v", err)
		}
	}()

	select {
	case s := <-app.interrupt:
		log.Infof("app - Start - signal: %v", s)
	case err := <-httpServer.Notify():
		log.Errorf("app - Start - server error: %v", err)
	}

	log.Info("Shutting down...")
}
