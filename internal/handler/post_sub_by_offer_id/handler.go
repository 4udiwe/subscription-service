package post_sub_by_offer_id

import (
	"errors"
	"net/http"
	"time"

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

type PostSubscriptionByOfferIDRequest struct {
	UserID    uuid.UUID `json:"user_id" validate:"required,uuid"`
	OfferID   uuid.UUID `json:"offer_id" validate:"required,uuid"`
	StartDate string    `json:"start_date" validate:"required,datetime=2006-01-02"`
}

type PostSubscriptionByOfferIDResponse struct {
	SubscriptionID uuid.UUID `json:"subscription_id"`
	UserID         uuid.UUID `json:"user_id"`
	OfferName      string    `json:"offer_name"`
	Price          int       `json:"price"`
	StartDate      string    `json:"start_date"`
	EndDate        string    `json:"end_date"`
}

// Create a new subscription by offer ID
// @Summary Создание новой подписки по ID предложения
// @Description Создание новой подписки для пользователя по ID предложения, полученного из ендпоинта всех предложений
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body PostSubscriptionByOfferIDRequest true "subscription info"
// @Success 201 {object} PostSubscriptionByOfferIDResponse
// @Failure 400 {string} ErrorResponse
// @Failure 500 {string} ErrorResponse
// @Router /subscriptions/by-offer [post]
func (h *handler) Handle(c echo.Context, in PostSubscriptionByOfferIDRequest) error {
	startDate, err := time.Parse("2006-01-02", in.StartDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid start_date format")
	}

	sub, err := h.s.CreateSubscriptionByOfferID(c.Request().Context(), in.UserID, in.OfferID, startDate)

	if err != nil {
		if errors.Is(err, subscription.ErrOfferNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if errors.Is(err, subscription.ErrUserAlreadyHasActiveSubscription) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, PostSubscriptionByOfferIDResponse{
		SubscriptionID: sub.ID,
		UserID:         sub.UserID,
		OfferName:      sub.OfferName,
		Price:          sub.Price,
		StartDate:      sub.StartDate.Format("2006-01-02"),
		EndDate:        sub.EndDate.Format("2006-01-02"),
	})
}
