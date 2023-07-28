package main

import (
	"github.com/gh0st3e/NASA_API/internal/config"
	"github.com/gh0st3e/NASA_API/internal/db"
	"github.com/gh0st3e/NASA_API/internal/handler"
	"github.com/gh0st3e/NASA_API/internal/service"
	"github.com/gh0st3e/NASA_API/internal/store"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

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

	apodStore := store.MewStore(psql)
	apodService := service.NewService(log, apodStore)
	apodHandler := handler.NewHandler(log, apodService)

	server := gin.New()

	apodHandler.Mount(server)

	err = server.Run(cfg.Server.Address)
	if err != nil {
		panic(err)
	}

}
