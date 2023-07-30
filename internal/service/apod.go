package service

import (
	"context"
	"fmt"

	"github.com/gh0st3e/NASA_API/internal/entity"

	"github.com/sirupsen/logrus"
)

type ApodActions interface {
	SaveApod(ctx context.Context, apod *entity.Apod) error
	RetrieveApodByDate(ctx context.Context, date string) (*entity.Apod, error)
	RetrieveAllApods(ctx context.Context) ([]entity.Apod, error)
}

type ApodService struct {
	log       *logrus.Logger
	apodStore ApodActions
}

func NewApodService(log *logrus.Logger, apodStore ApodActions) *ApodService {
	return &ApodService{
		log:       log,
		apodStore: apodStore,
	}
}

func (a *ApodService) SaveApod(ctx context.Context, apod *entity.Apod) error {
	a.log.Info("[SaveApod] started")

	err := a.apodStore.SaveApod(ctx, apod)
	if err != nil {
		a.log.Errorf("[SaveApod] error while save apod to db: \n%s", err.Error())
		return fmt.Errorf("error while save apod to db: \n%w", err)
	}

	a.log.Info("[SaveApod] ended")

	return nil
}

func (a *ApodService) RetrieveAllApods(ctx context.Context) ([]entity.Apod, error) {
	a.log.Info("[RetrieveAllApods] started")

	apods, err := a.apodStore.RetrieveAllApods(ctx)
	if err != nil {
		a.log.Errorf("[RetrieveAllApods] error while retrieving apods from db: \n%s", err.Error())
		return nil, fmt.Errorf("error while retrieving apods from db: %w", err)
	}

	a.log.Info(apods)
	a.log.Info("[RetrieveAllApods] ended")

	return apods, nil
}

func (a *ApodService) RetrieveApodByDate(ctx context.Context, date string) (*entity.Apod, error) {
	a.log.Info("[RetireveApodByDate] started")

	apod, err := a.apodStore.RetrieveApodByDate(ctx, date)
	if err != nil {
		a.log.Errorf("[RetireveApodByDate] error while retrieving apod from db: \n%s", err.Error())
		return nil, fmt.Errorf("error while retrieving apod from db: %w", err)
	}

	a.log.Info(apod)
	a.log.Info("[RetireveApodByDate] ended")

	return apod, nil
}
