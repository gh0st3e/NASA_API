package main

import (
	"github.com/gh0st3e/NASA_API/internal/clients"
	"github.com/gh0st3e/NASA_API/internal/config"
	"github.com/gh0st3e/NASA_API/internal/db"
	"github.com/gh0st3e/NASA_API/internal/handler"
	"github.com/gh0st3e/NASA_API/internal/service"
	"github.com/gh0st3e/NASA_API/internal/store"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @title NASA_API
// @version 1.0
// @description app provide API to get APOD from NASA

// @host localhost:8090
// @BasePath /api

func main() {
	log := logrus.New()

	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	psql, err := db.InitPSQL(cfg.PSQLDatabase)
	if err != nil {
		log.Fatal(err)
	}

	minio, err := db.InitMinio(cfg.MinioConfig)
	if err != nil {
		log.Fatal(err)
	}

	apodStore := store.NewApodStore(psql)
	imageStore := store.NewImageStore(minio, cfg.MinioConfig.BucketName)

	apodService := service.NewApodService(log, apodStore)
	imageService := service.NewImageService(log, imageStore)

	hdlr := handler.NewHandler(log, apodService, imageService)

	fileLog, err := config.InitFileLogger()
	if err != nil {
		log.Fatal(err)
	}

	nasaClient := clients.NewNasaClient(fileLog, cfg.NasaClientConfig)
	worker, err := service.NewWorker(fileLog, nasaClient, apodService, imageService)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := worker.RunWorker()
		if err != nil {
			log.Error(err)
		}
	}()

	server := gin.New()

	hdlr.Mount(server)

	err = server.Run(":" + cfg.Server.Port)
	if err != nil {
		panic(err)
	}
}
