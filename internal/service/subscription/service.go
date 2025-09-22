package subscription

import (
	"context"
	"errors"
	"time"

	"github.com/4udiwe/subscription-serivce/internal/entity"
	offer_repo "github.com/4udiwe/subscription-serivce/internal/repository/offer"
	subscription_repo "github.com/4udiwe/subscription-serivce/internal/repository/subscription"
	"github.com/4udiwe/subscription-serivce/pkg/transactor"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type SubscriptionService struct {
	subRepository   SubscriptionRepository
	offerRepository OfferRepository
	txManager       transactor.Transactor
}

func New(subRepo SubscriptionRepository, offerRepo OfferRepository, txManager transactor.Transactor) *SubscriptionService {
	return &SubscriptionService{
		subRepository:   subRepo,
		offerRepository: offerRepo,
		txManager:       txManager,
	}
}

func (s *SubscriptionService) CreateSubscription(
	ctx context.Context,
	userID uuid.UUID,
	serviceName string,
	price int,
	startDate time.Time,
	endDate *time.Time,
) (entity.Subscription, error) {

	logrus.Infof("SubscriptionService.CreateSubscription called: userID=%s, serviceName=%s, price=%d, startDate=%v, endDate=%v", userID, serviceName, price, startDate, endDate)
	var sub entity.Subscription

	err := s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		// check if offer with given name and prise exists
		offer, err := s.offerRepository.GetByNameAndPrice(ctx, serviceName, price)
		if err != nil {
			if errors.Is(err, offer_repo.ErrOfferNotFound) {
				// if not -> create it
				durationMonth := 1
				if endDate != nil {
					durationMonth = int(endDate.Sub(startDate).Hours() / (24 * 30))
				}

				offer, err = s.offerRepository.Create(ctx, serviceName, price, durationMonth)
				if err != nil {
					logrus.Errorf("SubscriptionService.CreateSubscription error creating offer: %v", err)
					return err
				}
			}
			return err
		}
		// create subscription
		sub, err = s.subRepository.Create(ctx, userID, offer.ID, startDate, startDate.AddDate(0, offer.DurationMonth, 0))
		if err != nil {
			logrus.Errorf("SubscriptionService.CreateSubscription error creating subscription: %v", err)
		}
		return err
	})

	if err != nil {
		logrus.Errorf("SubscriptionService.CreateSubscription error: %v", err)
		return entity.Subscription{}, err
	}
	logrus.Infof("SubscriptionService.CreateSubscription success: id=%s", sub.ID)
	return sub, nil
}

func (s *SubscriptionService) CreateSubscriptionByOfferID(ctx context.Context, userID, offerID uuid.UUID, startDate time.Time) (entity.Subscription, error) {
	logrus.Infof("SubscriptionService.CreateSubscriptionByOfferID called: userID=%s, offerID=%s, startDate=%v", userID, offerID, startDate)
	offer, err := s.offerRepository.GetByID(ctx, offerID)
	if err != nil {
		if errors.Is(err, offer_repo.ErrOfferNotFound) {
			logrus.Errorf("SubscriptionService.CreateSubscriptionByOfferID error: offer not found")
			return entity.Subscription{}, ErrOfferNotFound
		}
		logrus.Errorf("SubscriptionService.CreateSubscriptionByOfferID error: %v", err)
		return entity.Subscription{}, err
	}

	sub, err := s.subRepository.Create(ctx, userID, offer.ID, startDate, startDate.AddDate(0, offer.DurationMonth, 0))
	if err != nil {
		logrus.Errorf("SubscriptionService.CreateSubscriptionByOfferID error creating subscription: %v", err)
		return entity.Subscription{}, err
	}
	logrus.Infof("SubscriptionService.CreateSubscriptionByOfferID success: id=%s", sub.ID)
	return sub, nil
}

func (s *SubscriptionService) GetAllSubscriptions(ctx context.Context) ([]entity.Subscription, error) {
	logrus.Info("SubscriptionService.GetAllSubscriptions called")
	subs, err := s.subRepository.GetAll(ctx)
	if err != nil {
		logrus.Errorf("SubscriptionService.GetAllSubscriptions error: %v", err)
		return nil, err
	}
	logrus.Infof("SubscriptionService.GetAllSubscriptions success: count=%d", len(subs))
	return subs, nil
}

func (s *SubscriptionService) GetAllByUserIDAndSubscriptionName(
	ctx context.Context,
	userID uuid.UUID,
	subscriptionName string,
	startPeriod *time.Time,
	endPeriod *time.Time,
) ([]entity.Subscription, error) {
	logrus.Infof("SubscriptionService.GetAllByUserIDAndSubscriptionName called: userID=%s, subscriptionName=%s, startPeriod=%v, endPeriod=%v", userID, subscriptionName, startPeriod, endPeriod)
	subs, err := s.subRepository.GetAllByUserIDAndSubscriptionName(ctx, userID, subscriptionName, startPeriod, endPeriod)
	if err != nil {
		logrus.Errorf("SubscriptionService.GetAllByUserIDAndSubscriptionName error: %v", err)
		return nil, err
	}
	logrus.Infof("SubscriptionService.GetAllByUserIDAndSubscriptionName success: count=%d", len(subs))
	return subs, nil
}

func (s *SubscriptionService) DeleteSubscription(ctx context.Context, subID uuid.UUID) error {
	logrus.Infof("SubscriptionService.DeleteSubscription called: subID=%s", subID)
	err := s.subRepository.Delete(ctx, subID)
	if err != nil {
		if errors.Is(err, subscription_repo.ErrSubscriptionNotFound) {
			logrus.Errorf("SubscriptionService.DeleteSubscription error: subscription not found")
			return ErrSubscriptionNotFound
		}
		logrus.Errorf("SubscriptionService.DeleteSubscription error: %v", err)
		return err
	}
	logrus.Infof("SubscriptionService.DeleteSubscription success: subID=%s deleted", subID)
	return nil
}

func (s *SubscriptionService) GetAllSubscriptionsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Subscription, error) {
	logrus.Infof("SubscriptionService.GetAllSubscriptionsByUserID called: userID=%s", userID)

	subs, err := s.subRepository.GetAllByUserID(ctx, userID)
	if err != nil {
		logrus.Errorf("SubscriptionService.GetAllSubscriptionsByUserID error: %v", err)
		return nil, err
	}

	logrus.Infof("SubscriptionService.GetAllSubscriptionsByUserID success: count=%d", len(subs))
	return subs, nil
}
