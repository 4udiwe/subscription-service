package delete_sub

import (
	"errors"
	"net/http"

	h "github.com/4udiwe/subscription-service/internal/handler"
	"github.com/4udiwe/subscription-service/internal/handler/decorator"
	"github.com/4udiwe/subscription-service/internal/service/subscription"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type handler struct {
	s SubscriptionService
}

func New(s SubscriptionService) h.Handler {
	return decorator.NewBindAndValidateDecorator(&handler{s: s})
}

type DeleteSubscriptionRequest struct {
	SubscriptionID uuid.UUID `json:"subscription_id" validate:"required,uuid"`
}

// Delete subscription
// @Summary Удаление подписки
// @Description Удаление подписки по ID. Не удаляет предложение, на которое была оформлена подписка.
// @Tags subscriptions
// @Accept json
// @Param subscription body DeleteSubscriptionRequest true "subscription to delete"
// @Success 202 {string} string "No Content"
// @Failure 404 {string} ErrorResponse
// @Failure 500 {string} ErrorResponse
// @Router /subscriptions [delete]
func (h *handler) Handle(c echo.Context, in DeleteSubscriptionRequest) error {
	err := h.s.DeleteSubscription(c.Request().Context(), in.SubscriptionID)

	if err != nil {
		if errors.Is(err, subscription.ErrSubscriptionNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusAccepted)
}
