package get_subs_by_user_subname

import (
	"net/http"
	"time"

	"github.com/4udiwe/subscription-service/internal/entity"
	h "github.com/4udiwe/subscription-service/internal/handler"
	"github.com/4udiwe/subscription-service/internal/handler/decorator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type handler struct {
	s SubscriptionService
}

func New(s SubscriptionService) h.Handler {
	return decorator.NewBindAndValidateDecorator(&handler{s: s})
}

type Request struct {
	UserID    uuid.UUID `json:"user_id" validate:"required,uuid"`
	OfferName string    `json:"offer_name" validate:"required"`
	StartDate string    `json:"start_date" validate:"omitempty"`
	EndDate   string    `json:"end_date" validate:"omitempty"`
}

type Response struct {
	TotalPrice    int            `json:"total_price"`
	Subscriptions []Subscription `json:"subscriptions"`
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
// @Param user body Request true "user ID and subscription name"
// @Success 200 {object} Response
// @Failure 400 {string} ErrorResponse
// @Failure 500 {string} ErrorResponse
// @Router /subscriptions/by-user-and-subname [post]
func (h *handler) Handle(c echo.Context, in Request) error {
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

	sub, totalPrice, err := h.s.GetAllWithPriceByUserIDAndSubscriptionName(c.Request().Context(), in.UserID, in.OfferName, startDate, endDate)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
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
	})
}
