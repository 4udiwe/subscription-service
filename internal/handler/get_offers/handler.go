package get_offers

import (
	"math"
	"net/http"

	"github.com/4udiwe/subscription-service/internal/entity"
	h "github.com/4udiwe/subscription-service/internal/handler"
	decorator "github.com/4udiwe/subscription-service/internal/handler/decorator"
	"github.com/labstack/echo/v4"
)

const PAGE_NUMBER = 1
const PAGE_SIZE = 10

type handler struct {
	s OfferService
}

func New(s OfferService) h.Handler {
	return decorator.NewBindAndValidateDecorator(&handler{s: s})
}

type GetAllOffersRequest struct {
	Page     int `query:"page"`
	PageSize int `query:"page_size"`
}

type GetAllOffersResponse struct {
	Offers     []entity.Offer `json:"offers"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalItems int            `json:"total_items"`
	TotalPages int            `json:"total_pages"`
}

// Get all offers
// @Summary Получение всех офферов
// @Description Получение списка всех офферов
// @Tags offers
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы" default(1)
// @Param page_size query int false "Размер страницы" default(10) maximum(100)
// @Success 200 {object} GetAllOffersResponse
// @Failure 400 {string} ErrorResponse
// @Failure 500 {string} ErrorResponse
// @Router /offers [get]
func (h *handler) Handle(c echo.Context, in GetAllOffersRequest) error {
	if in.Page == 0 {
		in.Page = PAGE_NUMBER
	}

	if in.PageSize <= 0 {
		in.PageSize = PAGE_SIZE
	} else if in.PageSize > 100 {
		in.PageSize = 100
	}

	offers, totalCount, err := h.s.GetAllOffers(c.Request().Context(), in.Page, in.PageSize)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(in.PageSize)))

	return c.JSON(http.StatusOK, GetAllOffersResponse{
		Offers:     offers,
		Page:       in.Page,
		PageSize:   in.PageSize,
		TotalItems: totalCount,
		TotalPages: totalPages,
	})
}
