package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gh0st3e/NASA_API/internal/entity"

	"github.com/go-co-op/gocron"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	dirPermission  = 0o777
	imagesPath     = "assets/images"
	filePermission = 0o644
	fileExtension  = ".jpg"
)

type ClientActions interface {
	GetApod() (*entity.Apod, error)
}

type ApodSaver interface {
	SaveApod(ctx context.Context, apod *entity.Apod) error
}

type ImageSaver interface {
	PutImage(ctx context.Context, objectName, filePath, contentType string) error
}

type Worker struct {
	log        *logrus.Logger
	nasaClient ClientActions
	apodSaver  ApodSaver
	imageSaver ImageSaver
}

func NewWorker(log *logrus.Logger, nasaClient ClientActions, apodSaver ApodSaver, imageSaver ImageSaver) (*Worker, error) {
	err := os.MkdirAll(imagesPath, dirPermission)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			return nil, err
		}
	}

	return &Worker{
		log:        log,
		nasaClient: nasaClient,
		apodSaver:  apodSaver,
		imageSaver: imageSaver,
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

	imgResponse, err := http.Get(apod.URL)
	if err != nil {
		w.log.Errorf("error while request img url: %s", err.Error())
		return
	}
	defer imgResponse.Body.Close()

	body, err := io.ReadAll(imgResponse.Body)
	if err != nil {
		w.log.Errorf("error while reading response: %s", err.Error())
		return
	}

	err = os.WriteFile(fmt.Sprintf("%s/%s%s", imagesPath, apod.Date, fileExtension), body, filePermission)
	if err != nil {
		w.log.Errorf("error while writing file: %s", err.Error())
		return
	}
	defer func() {
		err := os.Remove(fmt.Sprintf("%s/%s%s", imagesPath, apod.Date, fileExtension))
		if err != nil {
			w.log.Errorf("error while deleting image")
		}
	}()

	objectName := apod.Date + fileExtension
	filePath := filepath.Join(imagesPath, objectName)
	contentType := "image/jpg"
	ctx := context.Background()

	err = w.imageSaver.PutImage(ctx, objectName, filePath, contentType)
	if err != nil {
		w.log.Error(err)
		return
	}

	err = w.apodSaver.SaveApod(context.Background(), apod)
	if err != nil {
		w.log.Error(err)
		return
	}
}
