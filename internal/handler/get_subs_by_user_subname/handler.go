package get_subs_by_user_subname

import (
	"math"
	"net/http"
	"time"

	"github.com/4udiwe/subscription-service/internal/entity"
	h "github.com/4udiwe/subscription-service/internal/handler"
	"github.com/4udiwe/subscription-service/internal/handler/decorator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

const PAGE_NUMBER = 1
const PAGE_SIZE = 10

type handler struct {
	s SubscriptionService
}

func New(s SubscriptionService) h.Handler {
	return decorator.NewBindAndValidateDecorator(&handler{s: s})
}

type GetSubsByUserAndServiceNameRequest struct {
	UserID    uuid.UUID `json:"user_id" validate:"required,uuid"`
	OfferName string    `json:"offer_name" validate:"required"`
	StartDate string    `json:"start_date" validate:"omitempty"`
	EndDate   string    `json:"end_date" validate:"omitempty"`
	Page      int       `query:"page"`
	PageSize  int       `query:"page_size"`
}

type GetSubsByUserAndServiceNameResponse struct {
	TotalPrice    int            `json:"total_price"`
	Subscriptions []Subscription `json:"subscriptions"`
	Page          int            `json:"page"`
	PageSize      int            `json:"page_size"`
	TotalItems    int            `json:"total_items"`
	TotalPages    int            `json:"total_pages"`
}

type Subscription struct {
	SubscriptionID uuid.UUID `json:"subscription_id"`
	UserID         uuid.UUID `json:"user_id"`
	OfferName      string    `json:"offer_name"`
	Price          int       `json:"price"`
	StartDate      string    `json:"start_date"`
	EndDate        string    `json:"end_date"`
}

// Get all subscriptions by user ID and subscription name
// @Summary Получение подписок по ID пользователя и названию подписки
// @Description Получение списка подписок для указанного пользователя и названия подписки с возможностью фильтрации по дате начала и окончания
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user body GetSubsByUserAndServiceNameRequest true "user ID and subscription name"
// @Param page query int false "Номер страницы" default(1)
// @Param page_size query int false "Размер страницы" default(10) maximum(100)
// @Success 200 {object} GetSubsByUserAndServiceNameResponse
// @Failure 400 {string} ErrorResponse
// @Failure 500 {string} ErrorResponse
// @Router /subscriptions/by-user-and-subname [get]
func (h *handler) Handle(c echo.Context, in GetSubsByUserAndServiceNameRequest) error {
	if in.Page == 0 {
		in.Page = PAGE_NUMBER
	}

	if in.PageSize <= 0 {
		in.PageSize = PAGE_SIZE
	} else if in.PageSize > 100 {
		in.PageSize = 100
	}

	var startDate, endDate *time.Time

	if in.StartDate != "" {
		parsedStartDate, err := time.Parse("2006-01-02", in.StartDate)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid start_date format")
		}
		startDate = &parsedStartDate
	}
	if in.EndDate != "" {
		parsedEndDate, err := time.Parse("2006-01-02", in.EndDate)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid end_date format")
		}
		endDate = &parsedEndDate
	}

	sub, totalPrice, totalCount, err := h.s.GetAllWithPriceByUserIDAndSubscriptionName(c.Request().Context(), in.UserID, in.OfferName, startDate, endDate, in.Page, in.PageSize)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(in.PageSize)))

	return c.JSON(http.StatusOK, GetSubsByUserAndServiceNameResponse{
		TotalPrice: totalPrice,
		Subscriptions: lo.Map(sub, func(s entity.SubscriptionFullInfo, _ int) Subscription {
			return Subscription{
				SubscriptionID: s.ID,
				UserID:         s.UserID,
				OfferName:      s.OfferName,
				Price:          s.Price,
				StartDate:      s.StartDate.Format("2006-01-02"),
				EndDate:        s.EndDate.Format("2006-01-02"),
			}
		}),
		Page:       in.Page,
		PageSize:   in.PageSize,
		TotalItems: totalCount,
		TotalPages: totalPages,
	})
}
