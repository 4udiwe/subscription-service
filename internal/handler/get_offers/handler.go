package get_offers

import (
	"net/http"

	"github.com/4udiwe/subscription-service/internal/entity"
	h "github.com/4udiwe/subscription-service/internal/handler"
	decorator "github.com/4udiwe/subscription-service/internal/handler/decorator"
	"github.com/labstack/echo/v4"
)

type handler struct {
	s OfferService
}

func New(s OfferService) h.Handler {
	return decorator.NewBindAndValidateDecorator(&handler{s: s})
}

type Request struct{}

type Response struct {
	Offers []entity.Offer `json:"offers"`
}

// Get all offers
// @Summary Получение всех офферов
// @Description Получение списка всех офферов
// @Tags offers
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {string} ErrorResponse
// @Failure 500 {string} ErrorResponse
// @Router /offers [get]
func (h *handler) Handle(c echo.Context, in Request) error {
	offers, err := h.s.GetAllOffers(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{Offers: offers})
}
