package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gh0st3e/NASA_API/internal/config"
	"github.com/gh0st3e/NASA_API/internal/entity"

	"github.com/sirupsen/logrus"
)

type NasaClient struct {
	log    *logrus.Logger
	apiKey string
}

func NewNasaClient(log *logrus.Logger, cfg config.NasaClientConfig) *NasaClient {
	return &NasaClient{
		log:    log,
		apiKey: cfg.ApiKey,
	}
}

func (n *NasaClient) GetApod() (*entity.Apod, error) {
	n.log.Info("[GetApod] started")

	url := fmt.Sprintf("https://api.nasa.gov/planetary/apod?api_key=%s", n.apiKey)

	resp, err := http.Get(url)
	if err != nil {
		n.log.Errorf("[GetApod] error while request api.nasa.gov: %s", err.Error())
		return nil, fmt.Errorf("error while request api.nasa.gov: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		n.log.Errorf("[GetApod] error while parse request: %s", err.Error())
		return nil, fmt.Errorf("error while arse request: %w", err)
	}

	var apod = &entity.Apod{}

	err = json.Unmarshal(body, apod)
	if err != nil {
		n.log.Errorf("[GetApod] error while unmarshal data: %s", err.Error())
		return nil, fmt.Errorf("error while unmarshal data: %w", err)
	}

	n.log.Info(apod)
	n.log.Info("[GetApod] ended")

	return apod, nil
}
