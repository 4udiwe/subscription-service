package get_subs

import (
	"net/http"

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

type Request struct{}

type Response struct {
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

func (h *handler) Handle(c echo.Context, in Request) error {
	sub, err := h.s.GetAllSubscriptions(c.Request().Context())

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
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
