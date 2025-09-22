package offer

import (
	"context"
	"errors"

	"github.com/4udiwe/subscription-serivce/internal/entity"
	offer_repo "github.com/4udiwe/subscription-serivce/internal/repository/offer"
	"github.com/4udiwe/subscription-serivce/pkg/transactor"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type OfferService struct {
	offerRepository OfferRepository
	subRepository   SubscriptionRepository
	txManager       transactor.Transactor
}

func New(offerRepository OfferRepository, subRepository SubscriptionRepository, txManager transactor.Transactor) *OfferService {
	return &OfferService{
		offerRepository: offerRepository,
		subRepository:   subRepository,
		txManager:       txManager,
	}
}

func (s *OfferService) GetAllOffers(ctx context.Context) ([]entity.Offer, error) {
	logrus.Info("OfferService.GetAllOffers called")

	offers, err := s.offerRepository.GetAll(ctx)
	if err != nil {
		logrus.Errorf("OfferService.GetAllOffers error: %v", err)
		return nil, err
	}

	logrus.Info("OfferService.GetAllOffers success")
	return offers, nil
}

func (s *OfferService) DeleteOffer(ctx context.Context, offerID uuid.UUID) error {
	logrus.Infof("OfferService.DeleteOffer called: id=%s", offerID)

	err := s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		// check if offer have referring subscriptions
		subs, err := s.subRepository.GetAllByOfferID(txCtx, offerID)
		if err != nil {
			logrus.Errorf("OfferService.DeleteOffer error fetching subscriptions: %v", err)
			return err
		}

		// checking amount of subs
		if len(subs) > 0 {
			logrus.Errorf("Offer.DeleteOffer error: active subscriptions exist for this offer (ID=%v), could not delete", offerID)
			return ErrActiveSubscriptionsExist
		}

		// if its zero -> delete
		return s.offerRepository.Delete(txCtx, offerID)
	})

	if err != nil {
		if errors.Is(err, offer_repo.ErrOfferNotFound) {
			return ErrOfferNotFound
		}
		logrus.Errorf("OfferService.DeleteOffer error deleting offer: %v", err)
		return err
	}

	logrus.Infof("OfferService.DeleteOffer success: offer with ID=%s deleted", offerID)
	return nil
}
