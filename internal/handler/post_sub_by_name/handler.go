package post_sub_by_name

import (
	"net/http"
	"time"

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
	UserID      uuid.UUID `json:"user_id" validate:"required,uuid"`
	ServiceName string    `json:"service_name" validate:"required"`
	Price       int       `json:"price" validate:"required,min=0"`
	StartDate   string    `json:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate     *string   `json:"end_date" validate:"omitempty,datetime=2006-01-02"`
}

type Response struct {
	SubscriptionID uuid.UUID `json:"subscription_id"`
	UserID         uuid.UUID `json:"user_id"`
	OfferName      string    `json:"offer_name"`
	Price          int       `json:"price"`
	StartDate      string    `json:"start_date"`
	EndDate        string    `json:"end_date"`
}

func (h *handler) Handle(c echo.Context, in Request) error {
	startDate, err := time.Parse("2006-01-02", in.StartDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid start_date format")
	}
	var endDate *time.Time
	if in.EndDate != nil {
		parsedEndDate, err := time.Parse("2006-01-02", *in.EndDate)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid end_date format")
		}
		endDate = &parsedEndDate
	}
	sub, err := h.s.CreateSubscription(c.Request().Context(), in.UserID, in.ServiceName, in.Price, startDate, endDate)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		SubscriptionID: sub.ID,
		UserID:         sub.UserID,
		OfferName:      sub.OfferName,
		Price:          sub.Price,
		StartDate:      sub.StartDate.Format("2006-01-02"),
		EndDate:        sub.EndDate.Format("2006-01-02"),
	})
}
