package delete_sub

import (
	"net/http"

	h "github.com/4udiwe/subscription-service/internal/handler"
	"github.com/4udiwe/subscription-service/internal/handler/decorator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type handler struct {
	s SubscriptionService
}

func New(s SubscriptionService) h.Handler {
	return decorator.NewBindAndValidateDecorator(&handler{s: s})
}

type Request struct {
	SubscriptionID uuid.UUID `json:"subscription_id" validate:"required,uuid"`
}

func (h *handler) Handle(c echo.Context, in Request) error {
	err := h.s.DeleteSubscription(c.Request().Context(), in.SubscriptionID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusAccepted)
}
