package offer_repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/4udiwe/subscription-serivce/internal/database"
	"github.com/4udiwe/subscription-serivce/internal/entity"
	"github.com/4udiwe/subscription-serivce/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	*postgres.Postgres
}

func New(postgres *postgres.Postgres) *Repository {
	return &Repository{postgres}
}

func (r *Repository) Create(ctx context.Context, name string, price int, durationMonth int) (entity.Offer, error) {
	logrus.Infof("OfferRepository.Create called: name=%s, price=%d, durationMonth=%d", name, price, durationMonth)

	query, args, _ := r.Builder.
		Insert("offer").
		Columns("name", "price", "duration_month").
		Values(name, price, durationMonth).
		Suffix("RETURNING id, created_at").
		ToSql()

	offer := entity.Offer{
		Name:          name,
		Price:         price,
		DurationMonth: durationMonth,
	}

	err := r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(
		&offer.ID, &offer.CreatedAt,
	)
	if err != nil {
		logrus.Error("OfferRepository.Create error: ", err)
		if database.IsUniqueViolation(err) {
			return entity.Offer{}, ErrOfferWithNameAndPriceAlreadyExists
		}
		return entity.Offer{}, fmt.Errorf("OfferRepository.Create - failed to create offer: %w", err)
	}

	logrus.Infof("OfferRepository.Create success: offer created with ID=%d", offer.ID)
	return offer, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]entity.Offer, error) {
	logrus.Info("OfferRepository.GetAll called")
	query, args, _ := r.Builder.
		Select("id", "name", "price", "duration_month", "created_at, updated_at").
		From("offer").
		ToSql()

	rows, err := r.GetTxManager(ctx).Query(ctx, query, args...)
	if err != nil {
		logrus.Error("OfferRepository.GetAll error: ", err)
		return nil, fmt.Errorf("OfferRepository.GetAll - failed to get offers: %w", err)
	}
	defer rows.Close()

	var offers []entity.Offer
	for rows.Next() {
		var offer entity.Offer
		if err := rows.Scan(&offer.ID, &offer.Name, &offer.Price, &offer.DurationMonth, &offer.CreatedAt, &offer.UpdatedAt); err != nil {
			logrus.Error("OfferRepository.GetAll scan error: ", err)
			return nil, fmt.Errorf("OfferRepository.GetAll - scan error: %w", err)
		}
		offers = append(offers, offer)
	}

	logrus.Infof("OfferRepository.GetAll success: offers count=%d", len(offers))
	return offers, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (entity.Offer, error) {
	logrus.Infof("OfferRepository.GetById called: id=%d", id)
	query, args, _ := r.Builder.
		Select("id", "name", "price", "duration_month", "created_at").
		From("offer").
		Where("id = ?", id).
		ToSql()

	var offer entity.Offer

	err := r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(
		&offer.ID, &offer.Name, &offer.Price, &offer.DurationMonth, &offer.CreatedAt,
	)
	if err != nil {
		logrus.Error("OfferRepository.GetById error: ", err)
		return entity.Offer{}, ErrOfferNotFound
	}

	logrus.Infof("OfferRepository.GetById success: id=%d", offer.ID)
	return offer, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	logrus.Infof("OfferRepository.Delete called: id=%d", id)
	query, args, _ := r.Builder.
		Delete("offer").
		Where("id = ?", id).
		ToSql()

	result, err := r.GetTxManager(ctx).Exec(ctx, query, args...)
	if err != nil {
		logrus.Error("OfferRepository.Delete error: ", err)
		return fmt.Errorf("OfferRepository.Delete - failed to delete offer: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		logrus.Error("OfferRepository.Delete error: no offer found with id=", id)
		return ErrOfferNotFound
	}

	logrus.Infof("OfferRepository.Delete success: id=%d", id)
	return nil
}

func (r *Repository) GetByNameAndPrice(ctx context.Context, name string, price int) (entity.Offer, error) {
	logrus.Infof("OfferRepository.GetByNameAndPrice called: name=%s, price=%d", name, price)
	query, args, _ := r.Builder.
		Select("id", "name", "price", "duration_month", "created_at", "updated_at").
		From("offer").
		Where("name = ? AND price = ?", name, price).
		ToSql()

	var offer entity.Offer

	err := r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(&offer.ID, &offer.Name, &offer.Price, &offer.DurationMonth, &offer.CreatedAt, &offer.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Offer{}, ErrOfferNotFound
		}
		logrus.Error("OfferRepository.GetByNameAndPrice error: ", err)
		return entity.Offer{}, fmt.Errorf("OfferRepository.GetByNameAndPrice - failed to get offer: %w", err)
	}

	logrus.Infof("OfferRepository.GetByNameAndPrice success: id=%s", offer.ID)
	return offer, nil
}
