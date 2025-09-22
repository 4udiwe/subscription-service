package app

import (
	"fmt"
	"net/http"

	"github.com/4udiwe/subscription-service/pkg/validator"
	"github.com/labstack/echo/v4"
)

func (app *App) EchoHandler() *echo.Echo {
	if app.echoHandler != nil {
		return app.echoHandler
	}

	handler := echo.New()
	handler.Validator = validator.NewCustomValidator()

	app.configureRouter(handler)

	for _, r := range handler.Routes() {
		fmt.Printf("%s %s\n", r.Method, r.Path)
	}

	app.echoHandler = handler
	return app.echoHandler
}

func (app *App) configureRouter(handler *echo.Echo) {
	offersGroup := handler.Group("offers")
	{
		offersGroup.GET("", app.GetOffersHandler().Handle)
	}

	subsGroup := handler.Group("subscriptions")
	{
		subsGroup.GET("", app.GetSubscriptionsHandler().Handle)
		subsGroup.GET("/by_user", app.GetSubscriptionsByUserHandler().Handle)
		subsGroup.GET("/by_user_service_name", app.GetSubscriptionsByUserAndSubNameHandler().Handle)
		subsGroup.POST("/by_name", app.PostSubciptionByNameHandler().Handle)
		subsGroup.POST("/by_offer_id", app.PostSubciptionByOfferIDHandler().Handle)
		subsGroup.DELETE("", app.DeleteProductHandler().Handle)
	}

	handler.GET("/health", func(c echo.Context) error { return c.NoContent(http.StatusOK) })
}
