package config

import (
	"errors"
	"os"

	"github.com/sirupsen/logrus"
)

func InitFileLogger() (*logrus.Logger, error) {
	log := logrus.New()

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			_, err := os.Create("app.log")
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	log.SetOutput(file)

	return log, nil
}
