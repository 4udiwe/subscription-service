package post_offer

import (
	"errors"
	"net/http"

	h "github.com/4udiwe/subscription-service/internal/handler"
	"github.com/4udiwe/subscription-service/internal/handler/decorator"
	service "github.com/4udiwe/subscription-service/internal/service/offer"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type handler struct {
	s OfferService
}

func New(s OfferService) h.Handler {
	return decorator.NewBindAndValidateDecorator(&handler{s: s})
}

type Request struct {
	ServiceName    string `json:"service_name" validate:"required"`
	Price          int    `json:"price" validate:"required,min=0"`
	DurationMonths int    `json:"duration_months" validate:"required,min=1"`
}

type Response struct {
	OfferID        uuid.UUID `json:"offer_id"`
	ServiceName    string    `json:"service_name"`
	Price          int       `json:"price"`
	DurationMonths int       `json:"duration_months"`
	CreatedAt      string    `json:"created_at"`
}

// Create a new offer
// @Summary Создание нового предложения
// @Description Создание нового предложения с указанными параметрами
// @Tags offers
// @Accept json
// @Produce json
// @Param offer body Request true "Offer details"
// @Success 201 {object} Response
// @Failure 409 {string} ErrorResponse
// @Failure 500 {string} ErrorResponse
// @Router /offers [post]
func (h *handler) Handle(c echo.Context, in Request) error {
	offer, err := h.s.CreateOffer(c.Request().Context(), in.ServiceName, in.Price, in.DurationMonths)
	if err != nil {
		if errors.Is(err, service.ErrOfferWithNameAndPriceAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, Response{
		OfferID:        offer.ID,
		ServiceName:    offer.Name,
		Price:          offer.Price,
		DurationMonths: offer.DurationMonths,
		CreatedAt:      offer.CreatedAt.Format("2006-01-02"),
	})
}
