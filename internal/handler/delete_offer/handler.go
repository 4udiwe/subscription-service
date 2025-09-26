package delete_offer

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
	s OfferService
}

func New(s OfferService) h.Handler {
	return decorator.NewBindAndValidateDecorator(&handler{s: s})
}

type DeleteOfferRequest struct {
	OfferID uuid.UUID `json:"offer_id" validate:"required,uuid"`
}

// Delete offer
// @Summary Удаление предложения
// @Description Удаление предложения по ID. Если есть активные подписки на это предложение, оно не будет удалено.
// @Tags offers
// @Accept json
// @Param offer body DeleteOfferRequest true "offer to delete"
// @Success 202 {string} string "No Content"
// @Failure 404 {string} ErrorResponse
// @Failure 500 {string} ErrorResponse
// @Router /offers [delete]
func (h *handler) Handle(c echo.Context, in DeleteOfferRequest) error {
	err := h.s.DeleteOffer(c.Request().Context(), in.OfferID)

	if err != nil {
		if errors.Is(err, subscription.ErrOfferNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusAccepted)
}
