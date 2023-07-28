package service

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gh0st3e/NASA_API/internal/entity"

	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
)

const (
	dirPermission  = 0777
	imagesPath     = "assets/images"
	filePermission = 0644
	fileExtension  = ".jpg"
)

type ClientActions interface {
	GetApod() (*entity.Apod, error)
}

type Saver interface {
	SaveApod(ctx context.Context, apod entity.Apod) error
}

type Worker struct {
	log        *logrus.Logger
	nasaClient ClientActions
	saver      Saver
}

func NewWorker(log *logrus.Logger, nasaClient ClientActions, saver Saver) (*Worker, error) {
	err := os.MkdirAll(imagesPath, dirPermission)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			return nil, err
		}
	}

	return &Worker{
		log:        log,
		nasaClient: nasaClient,
		saver:      saver,
	}, nil
}

func (w *Worker) RunWorker() error {
	timer := gocron.NewScheduler(time.UTC)

	_, err := timer.Every(24).Hours().Do(func() {
		w.GetApod()
	})
	if err != nil {
		return err
	}

	timer.StartBlocking()

	return nil
}

func (w *Worker) GetApod() {
	apod, err := w.nasaClient.GetApod()
	if err != nil {
		w.log.Error("error from client: %s", err.Error())
		return
	}

	img, err := http.Get(apod.Url)
	if err != nil {
		w.log.Errorf("error while request img url: %s", err.Error())
		return
	}

	body, err := io.ReadAll(img.Body)
	if err != nil {
		w.log.Errorf("error while reading response: %s", err.Error())
		return
	}

	err = os.WriteFile(fmt.Sprintf("%s/%s%s", imagesPath, apod.Date, fileExtension), body, filePermission)
	if err != nil {
		w.log.Errorf("error while writing file: %s", err.Error())
		return
	}

	err = w.saver.SaveApod(context.Background(), *apod)
	if err != nil {
		w.log.Error(err)
		return
	}
}
