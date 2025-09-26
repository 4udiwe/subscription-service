package subscription

import (
	"context"
	"errors"
	"time"

	"github.com/4udiwe/subscription-service/internal/entity"
	offer_repo "github.com/4udiwe/subscription-service/internal/repository/offer"
	subscription_repo "github.com/4udiwe/subscription-service/internal/repository/subscription"
	"github.com/4udiwe/subscription-service/pkg/transactor"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const defaultDurationMonths = 1

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
) (entity.SubscriptionFullInfo, error) {
	logrus.Infof("SubscriptionService.CreateSubscription called: userID=%s, serviceName=%s, price=%d, startDate=%v, endDate=%v", userID, serviceName, price, startDate, endDate)
	var sub entity.SubscriptionFullInfo

	err := s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		// check if offer with given name and price exists
		offer, err := s.offerRepository.GetByNameAndPrice(ctx, serviceName, price)
		if err != nil && !errors.Is(err, offer_repo.ErrOfferNotFound) {
			logrus.Errorf("SubscriptionService.GetByNameAndPrice error getting offer: %v", err)
			return ErrCannotFindOffer
		}

		if errors.Is(err, offer_repo.ErrOfferNotFound) {
			// if not -> create it
			durationMonths := defaultDurationMonths
			if endDate != nil {
				durationMonths = int(endDate.Sub(startDate).Hours() / (24 * 30))
			}

			offer, err = s.offerRepository.Create(ctx, serviceName, price, durationMonths)
			if err != nil {
				logrus.Errorf("SubscriptionService.CreateSubscription error creating offer: %v", err)
				return ErrCannotCreateOffer
			}
		}

		// check if user has active subscription for the offer on the start date
		hasActive, err := s.subRepository.HasActiveSubscriptionOnServiceForDate(ctx, userID, serviceName, startDate)
		if err != nil {
			logrus.Errorf("SubscriptionService.CreateSubscription error checking active subscription: %v", err)
			return ErrCannotCheckActiveSubscription
		}
		if hasActive {
			logrus.Errorf("SubscriptionService.CreateSubscription error: user already has an active subscription for this offer on the start date")
			return ErrUserAlreadyHasActiveSubscription
		}

		// create subscription
		sub.Subscription, err = s.subRepository.Create(ctx, userID, offer.ID, startDate, startDate.AddDate(0, offer.DurationMonths, 0))
		sub.OfferName = offer.Name
		sub.Price = offer.Price

		if err != nil {
			logrus.Errorf("SubscriptionService.CreateSubscription error creating subscription: %v", err)
			return ErrCannotCreateSubscription
		}
		return nil
	})

	if err != nil {
		return entity.SubscriptionFullInfo{}, err
	}

	logrus.Infof("SubscriptionService.CreateSubscription success: id=%s", sub.ID)
	return sub, nil
}

func (s *SubscriptionService) CreateSubscriptionByOfferID(ctx context.Context, userID, offerID uuid.UUID, startDate time.Time) (entity.SubscriptionFullInfo, error) {
	logrus.Infof("SubscriptionService.CreateSubscriptionByOfferID called: userID=%s, offerID=%s, startDate=%v", userID, offerID, startDate)
	var subFullInfo entity.SubscriptionFullInfo

	err := s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		offer, err := s.offerRepository.GetByID(txCtx, offerID)
		if err != nil && !errors.Is(err, offer_repo.ErrOfferNotFound) {
			logrus.Errorf("SubscriptionService.CreateSubscriptionByOfferID error: %v", err)
			return ErrCannotFindOffer
		}

		if errors.Is(err, offer_repo.ErrOfferNotFound) {
			logrus.Errorf("SubscriptionService.CreateSubscriptionByOfferID error: offer not found")
			return ErrOfferNotFound
		}

		// check if user has active subscription for the offer on the start date
		hasActive, err := s.subRepository.HasActiveSubscriptionOnServiceForDate(ctx, userID, offer.Name, startDate)
		if err != nil {
			logrus.Errorf("SubscriptionService.CreateSubscription error checking active subscription: %v", err)
			return ErrCannotCheckActiveSubscription
		}

		if hasActive {
			logrus.Errorf("SubscriptionService.CreateSubscription error: user already has an active subscription for this offer on the start date")
			return ErrUserAlreadyHasActiveSubscription
		}

		sub, err := s.subRepository.Create(txCtx, userID, offer.ID, startDate, startDate.AddDate(0, offer.DurationMonths, 0))
		if err != nil {
			logrus.Errorf("SubscriptionService.CreateSubscriptionByOfferID error creating subscription: %v", err)
			return ErrCannotCreateSubscription
		}

		subFullInfo = entity.SubscriptionFullInfo{
			Subscription: sub,
			OfferName:    offer.Name,
			Price:        offer.Price,
		}

		return nil
	})

	if err != nil {
		return entity.SubscriptionFullInfo{}, err
	}

	logrus.Infof("SubscriptionService.CreateSubscriptionByOfferID success: id=%s", subFullInfo.ID)
	return subFullInfo, nil
}

func (s *SubscriptionService) GetAllSubscriptions(ctx context.Context, page int, pageSize int) ([]entity.SubscriptionFullInfo, int, error) {
	logrus.Info("SubscriptionService.GetAllSubscriptions called")

	limit := pageSize
	offset := (page - 1) * pageSize

	subs, total, err := s.subRepository.GetAll(ctx, limit, offset)
	if err != nil {
		logrus.Errorf("SubscriptionService.GetAllSubscriptions error: %v", err)
		return nil, 0, ErrCannotFetchSubscriptions
	}

	logrus.Infof("SubscriptionService.GetAllSubscriptions success: count=%d", len(subs))
	return subs, total, nil
}

func (s *SubscriptionService) GetAllWithPriceByUserIDAndSubscriptionName(
	ctx context.Context,
	userID uuid.UUID,
	subscriptionName string,
	startPeriod *time.Time,
	endPeriod *time.Time,
	page int,
	pageSize int,
) (subs []entity.SubscriptionFullInfo, price int, totalCount int, err error) {
	logrus.Infof("SubscriptionService.GetAllWithPriceByUserIDAndSubscriptionName called: userID=%s, subscriptionName=%s, startPeriod=%v, endPeriod=%v", userID, subscriptionName, startPeriod, endPeriod)

	limit := pageSize
	offset := (page - 1) * pageSize

	subs, price, totalCount, err = s.subRepository.GetAllByUserIDAndSubscriptionName(ctx, userID, subscriptionName, startPeriod, endPeriod, limit, offset)
	if err != nil {
		logrus.Errorf("SubscriptionService.GetAllWithPriceByUserIDAndSubscriptionName error: %v", err)
		return nil, 0, 0, ErrCannotFetchSubscriptions
	}

	logrus.Infof("SubscriptionService.GetAllWithPriceByUserIDAndSubscriptionName success: count=%d", len(subs))
	return subs, price, totalCount, nil
}

func (s *SubscriptionService) DeleteSubscription(ctx context.Context, subID uuid.UUID) error {
	logrus.Infof("SubscriptionService.DeleteSubscription called: subID=%s", subID)
	err := s.subRepository.Delete(ctx, subID)
	if err != nil && !errors.Is(err, subscription_repo.ErrSubscriptionNotFound) {
		logrus.Errorf("SubscriptionService.DeleteSubscription error: %v", err)
		return ErrCannotDeleteSubscription
	}

	if errors.Is(err, subscription_repo.ErrSubscriptionNotFound) {
		logrus.Errorf("SubscriptionService.DeleteSubscription error: subscription not found")
		return ErrSubscriptionNotFound
	}

	logrus.Infof("SubscriptionService.DeleteSubscription success: subID=%s deleted", subID)
	return nil
}

func (s *SubscriptionService) GetAllSubscriptionsByUserID(ctx context.Context, userID uuid.UUID, page int, pageSize int) ([]entity.SubscriptionFullInfo, int, error) {
	logrus.Infof("SubscriptionService.GetAllSubscriptionsByUserID called: userID=%s", userID)

	limit := pageSize
	offset := (page - 1) * pageSize

	subs, totalCount, err := s.subRepository.GetAllByUserID(ctx, userID, limit, offset)
	if err != nil {
		logrus.Errorf("SubscriptionService.GetAllSubscriptionsByUserID error: %v", err)
		return nil, 0, ErrCannotFetchSubscriptions
	}

	logrus.Infof("SubscriptionService.GetAllSubscriptionsByUserID success: count=%d", len(subs))
	return subs, totalCount, nil
}
