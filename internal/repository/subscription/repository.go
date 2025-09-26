package subscription_repo

import (
	"context"
	"fmt"
	"time"

	"github.com/4udiwe/subscription-service/internal/entity"
	"github.com/4udiwe/subscription-service/pkg/postgres"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	*postgres.Postgres
}

func New(postgres *postgres.Postgres) *Repository {
	return &Repository{postgres}
}

func (r *Repository) Create(ctx context.Context, userID, offerID uuid.UUID, startDate, endDate time.Time) (entity.Subscription, error) {
	logrus.Infof("SubscriptionRepository.Create called: userID=%s, offerID=%s", userID, offerID)
	query, args, _ := r.Builder.
		Insert("subscription").
		Columns("user_id", "offer_id", "start_date", "end_date").
		Values(userID, offerID, startDate, endDate).
		Suffix("RETURNING id, created_at, updated_at").
		ToSql()

	sub := entity.Subscription{
		UserID:    userID,
		OfferID:   offerID,
		StartDate: startDate,
		EndDate:   endDate,
	}
	err := r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(
		&sub.ID, &sub.CreatedAt, &sub.UpdatedAt,
	)
	logrus.Debugf("Scanned values: ID=%s, CreatedAt=%s, UpdatedAt=%s", sub.ID.String(), sub.CreatedAt.String(), sub.UpdatedAt.String())
	if err != nil {
		logrus.Error("SubscriptionRepository.Create error: ", err)
		return entity.Subscription{}, fmt.Errorf("SubscriptionRepository.Create - failed to create subscription: %w", err)
	}
	logrus.Infof("SubscriptionRepository.Create success: id=%s", sub.ID.String())
	return sub, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]entity.SubscriptionFullInfo, error) {
	logrus.Info("SubscriptionRepository.GetAll called")
	query, args, _ := r.Builder.
		Select("s.id", "s.user_id", "s.offer_id", "s.start_date", "s.end_date", "s.created_at", "s.updated_at", "o.name", "o.price").
		From("subscription s").
		Join("offer o ON s.offer_id = o.id").
		ToSql()

	rows, err := r.GetTxManager(ctx).Query(ctx, query, args...)
	if err != nil {
		logrus.Error("SubscriptionRepository.GetAll error: ", err)
		return nil, fmt.Errorf("SubscriptionRepository.GetAll - failed to get subscriptions: %w", err)
	}
	defer rows.Close()

	var subs []entity.SubscriptionFullInfo
	for rows.Next() {
		var sub entity.SubscriptionFullInfo
		if err := rows.Scan(&sub.ID, &sub.UserID, &sub.OfferID, &sub.StartDate, &sub.EndDate, &sub.CreatedAt, &sub.UpdatedAt, &sub.OfferName, &sub.Price); err != nil {
			logrus.Error("SubscriptionRepository.GetAll scan error: ", err)
			return nil, fmt.Errorf("SubscriptionRepository.GetAll - scan error: %w", err)
		}
		subs = append(subs, sub)
	}
	logrus.Infof("SubscriptionRepository.GetAll success: count=%d", len(subs))
	return subs, nil
}

func (r *Repository) GetById(ctx context.Context, id uuid.UUID) (entity.Subscription, error) {
	logrus.Infof("SubscriptionRepository.GetById called: id=%s", id)
	query, args, _ := r.Builder.
		Select("id", "user_id", "offer_id", "start_date", "end_date", "created_at", "updated_at").
		From("subscription").
		Where("id = ?", id).
		ToSql()

	var sub entity.Subscription
	err := r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(
		&sub.ID, &sub.UserID, &sub.OfferID, &sub.StartDate, &sub.EndDate, &sub.CreatedAt, &sub.UpdatedAt,
	)
	if err != nil {
		logrus.Error("SubscriptionRepository.GetById error: ", err)
		return entity.Subscription{}, ErrSubscriptionNotFound
	}
	logrus.Infof("SubscriptionRepository.GetById success: id=%s", sub.ID)
	return sub, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	logrus.Infof("SubscriptionRepository.Delete called: id=%s", id)
	query, args, _ := r.Builder.
		Delete("subscription").
		Where("id = ?", id).
		ToSql()

	result, err := r.GetTxManager(ctx).Exec(ctx, query, args...)
	if err != nil {
		logrus.Error("SubscriptionRepository.Delete error: ", err)
		return fmt.Errorf("SubscriptionRepository.Delete - failed to delete subscription: %w", err)
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		logrus.Error("SubscriptionRepository.Delete error: no subscription found with id=", id)
		return ErrSubscriptionNotFound
	}
	logrus.Infof("SubscriptionRepository.Delete success: id=%s", id)
	return nil
}

func (r *Repository) GetAllByUserIDAndSubscriptionName(
	ctx context.Context,
	userID uuid.UUID,
	subscriptionName string,
	startPeriod *time.Time,
	endPeriod *time.Time,
) (subs []entity.SubscriptionFullInfo, totalPrice int, err error) {
	logrus.Infof("SubscriptionRepository.GetByUserIDAndSubscriptionName called: userID=%s, subscriptionName=%s, startDate=%v, endDate=%v", userID, subscriptionName, startPeriod, endPeriod)

	builder := r.Builder.
		Select("s.id", "s.user_id", "s.offer_id", "s.start_date", "s.end_date", "s.created_at", "s.updated_at", "o.name", "o.price", "SUM(o.price) OVER() AS total_price").
		From("subscription s").
		Join("offer o ON s.offer_id = o.id").
		Where("s.user_id = ?", userID).
		Where("o.name = ?", subscriptionName)

	if startPeriod != nil {
		builder = builder.Where("s.start_date >= ?", *startPeriod)
	}
	if endPeriod != nil {
		builder = builder.Where("s.start_date <= ?", *endPeriod)
	}

	query, args, _ := builder.ToSql()

	rows, err := r.GetTxManager(ctx).Query(ctx, query, args...)
	if err != nil {
		logrus.Error("SubscriptionRepository.GetByUserIDAndSubscriptionName error: ", err)
		return nil, 0, fmt.Errorf("SubscriptionRepository.GetByUserIDAndSubscriptionName - failed to get subscriptions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var sub entity.SubscriptionFullInfo
		if err := rows.Scan(&sub.ID, &sub.UserID, &sub.OfferID, &sub.StartDate, &sub.EndDate, &sub.CreatedAt, &sub.UpdatedAt, &sub.OfferName, &sub.Price, &totalPrice); err != nil {
			logrus.Error("SubscriptionRepository.GetByUserIDAndSubscriptionName scan error: ", err)
			return nil, 0, fmt.Errorf("SubscriptionRepository.GetByUserIDAndSubscriptionName - scan error: %w", err)
		}
		subs = append(subs, sub)
	}
	logrus.Infof("SubscriptionRepository.GetByUserIDAndSubscriptionName success: count=%d", len(subs))
	return subs, totalPrice, nil
}

func (r *Repository) GetAllByOfferID(ctx context.Context, offerID uuid.UUID) ([]entity.Subscription, error) {
	logrus.Infof("SubscriptionRepository.GetAllByOfferID called: offerID=%s", offerID)
	query, args, _ := r.Builder.
		Select("s.id", "s.user_id", "s.offer_id", "s.start_date", "s.end_date", "s.created_at", "s.updated_at").
		From("subscription s").
		Join("offer o ON s.offer_id = o.id").
		Where("s.offer_id = ?", offerID).
		ToSql()

	rows, err := r.GetTxManager(ctx).Query(ctx, query, args...)
	if err != nil {
		logrus.Error("SubscriptionRepository.GetAllByOfferID error: ", err)
		return nil, fmt.Errorf("SubscriptionRepository.GetAllByOfferID - failed to get subscriptions: %w", err)
	}
	defer rows.Close()

	var subs []entity.Subscription
	for rows.Next() {
		var sub entity.Subscription
		if err := rows.Scan(&sub.ID, &sub.UserID, &sub.OfferID, &sub.StartDate, &sub.EndDate, &sub.CreatedAt, &sub.UpdatedAt); err != nil {
			logrus.Error("SubscriptionRepository.GetAllByOfferID scan error: ", err)
			return nil, fmt.Errorf("SubscriptionRepository.GetAllByOfferID - scan error: %w", err)
		}
		subs = append(subs, sub)
	}
	logrus.Infof("SubscriptionRepository.GetAllByOfferID success: count=%d", len(subs))
	return subs, nil
}

func (r *Repository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]entity.SubscriptionFullInfo, error) {
	logrus.Infof("SubscriptionRepository.GetAllByUserID called: userID=%s", userID)

	query, args, _ := r.Builder.
		Select("s.id", "s.user_id", "s.offer_id", "s.start_date", "s.end_date", "s.created_at", "s.updated_at", "o.name", "o.price").
		From("subscription s").
		Join("offer o ON s.offer_id = o.id").
		Where("s.user_id = ?", userID).
		ToSql()

	rows, err := r.GetTxManager(ctx).Query(ctx, query, args...)
	if err != nil {
		logrus.Error("SubscriptionRepository.GetAllByUserID error: ", err)
		return nil, fmt.Errorf("SubscriptionRepository.GetAllByUserID - failed to get subscriptions: %w", err)
	}
	defer rows.Close()

	var subs []entity.SubscriptionFullInfo

	for rows.Next() {
		var sub entity.SubscriptionFullInfo
		if err := rows.Scan(&sub.ID, &sub.UserID, &sub.OfferID, &sub.StartDate, &sub.EndDate, &sub.CreatedAt, &sub.UpdatedAt, &sub.OfferName, &sub.Price); err != nil {
			logrus.Error("SubscriptionRepository.GetAllByUserID scan error: ", err)
			return nil, fmt.Errorf("SubscriptionRepository.GetAllByUserID - scan error: %w", err)
		}
		subs = append(subs, sub)
	}

	logrus.Infof("SubscriptionRepository.GetAllByUserID success: count=%d", len(subs))
	return subs, nil
}

func (r *Repository) HasActiveSubscriptionOnServiceForDate(ctx context.Context, userID uuid.UUID, serviceName string, date time.Time) (bool, error) {
	logrus.Infof("SubscriptionRepository.HasActiveSubscriptionOnServiceForDate called: userID=%s, serviceName=%s, onDate=%s", userID, serviceName, date)

	var count int
	query, args, _ := r.Builder.
		Select("COUNT(*)").
		From("subscription s").
		Join("offer o ON s.offer_id = o.id").
		Where("s.user_id = ?", userID).
		Where("o.name = ?", serviceName).
		Where("s.start_date <= ? AND s.end_date > ?", date, date).
		ToSql()

	err := r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		logrus.Error("SubscriptionRepository.HasActiveSubscriptionOnServiceForDate error: ", err)
		return false, fmt.Errorf("SubscriptionRepository.HasActiveSubscriptionOnServiceForDate - failed to check active subscription: %w", err)
	}

	return count > 0, nil
}
