package service

import (
	"context"
	"fmt"

	"github.com/gh0st3e/NASA_API/internal/entity"

	"github.com/sirupsen/logrus"
)

type StoreActions interface {
	SaveApod(ctx context.Context, apod entity.Apod) error
	RetrieveAPODByDate(ctx context.Context, date string) (*entity.Apod, error)
	RetrieveAllApods(ctx context.Context) ([]entity.Apod, error)
}

type Service struct {
	log   *logrus.Logger
	store StoreActions
}

func NewService(log *logrus.Logger, store StoreActions) *Service {
	return &Service{
		log:   log,
		store: store,
	}
}

func (s *Service) SaveApod(ctx context.Context, apod entity.Apod) error {
	s.log.Info("[SaveApod] started")

	err := s.store.SaveApod(ctx, apod)
	if err != nil {
		s.log.Errorf("[SaveApod] error while save apod to db: \n%s", err.Error())
		return fmt.Errorf("error while save apod to db: \n%w", err)
	}

	s.log.Info("[SaveApod] ended")

	return nil
}

func (s *Service) RetrieveAllApods(ctx context.Context) ([]entity.Apod, error) {
	s.log.Info("[RetrieveAllApods] started")

	apods, err := s.store.RetrieveAllApods(ctx)
	if err != nil {
		s.log.Errorf("[RetrieveAllApods] error while retrieving apods from db: \n%s", err.Error())
		return nil, fmt.Errorf("error while retrieving apods from db: %w", err)
	}

	s.log.Info(apods)
	s.log.Info("[RetrieveAllApods] ended")

	return apods, nil
}

func (s *Service) RetrieveApodByDate(ctx context.Context, date string) (*entity.Apod, error) {
	s.log.Info("[RetireveApodByDate] started")

	apod, err := s.store.RetrieveAPODByDate(ctx, date)
	if err != nil {
		s.log.Errorf("[RetireveApodByDate] error while retrieving apod from db: \n%s", err.Error())
		return nil, fmt.Errorf("error while retrieving apod from db: %w", err)
	}

	s.log.Info(apod)
	s.log.Info("[RetireveApodByDate] ended")

	return apod, nil
}
