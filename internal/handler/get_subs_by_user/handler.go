package get_subs_by_user

import (
	"math"
	"net/http"

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

type GetSubscriptionsByUserRequest struct {
	UserID   uuid.UUID `json:"user_id" validate:"required,uuid"`
	Page     int       `query:"page"`
	PageSize int       `query:"page_size"`
}

type GetSubscriptionsByUserResponse struct {
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

// Get all subscriptions by user ID
// @Summary Получение подписок по ID пользователя
// @Description Получение списка подписок для указанного пользователя
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user body GetSubscriptionsByUserRequest true "user ID"
// @Param page query int false "Номер страницы" default(1)
// @Param page_size query int false "Размер страницы" default(10) maximum(100)
// @Success 200 {object} GetSubscriptionsByUserResponse
// @Failure 500 {string} ErrorResponse
// @Router /subscriptions/by-user [get]
func (h *handler) Handle(c echo.Context, in GetSubscriptionsByUserRequest) error {
	if in.Page == 0 {
		in.Page = PAGE_NUMBER
	}

	if in.PageSize <= 0 {
		in.PageSize = PAGE_SIZE
	} else if in.PageSize > 100 {
		in.PageSize = 100
	}

	sub, totalCount, err := h.s.GetAllSubscriptionsByUserID(c.Request().Context(), in.UserID, in.Page, in.PageSize)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(in.PageSize)))

	return c.JSON(http.StatusOK, GetSubscriptionsByUserResponse{
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
