// Package handler provides HTTP handlers for API endpoints.
// @description This is a sample API server for handling NASA API requests.
// @BasePath /api
package handler

import (
	"context"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os"

	_ "github.com/gh0st3e/NASA_API/docs"
	"github.com/gh0st3e/NASA_API/internal/entity"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ApodActions interface {
	RetrieveAllApods(ctx context.Context) ([]entity.Apod, error)
	RetrieveApodByDate(ctx context.Context, date string) (*entity.Apod, error)
}

type ImageActions interface {
	GetImage(ctx context.Context, fileName string) (string, error)
}

type Handler struct {
	log          *logrus.Logger
	apodService  ApodActions
	imageService ImageActions
}

func NewHandler(log *logrus.Logger, apodService ApodActions, imageService ImageActions) *Handler {
	return &Handler{
		log:          log,
		apodService:  apodService,
		imageService: imageService,
	}
}

// Mount registers API routes and handlers.
// @title NASA API Server
// @version 1.0
// @host localhost:8080
// @BasePath /api
func (h *Handler) Mount(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")

	api.GET("/apods", h.RetrieveAllApods)
	api.GET("/apod", h.RetrieveApodByDate)

	img := api.Group("img")

	img.GET("/apod", h.RetrieveApodImageByDate)

}

// RetrieveAllApods retrieves all APODs (Astronomy Picture of the Day).
// @Summary Get all APODs
// @Description Retrieves all Astronomy Picture of the Day.
// @Produce json
// @Success 200 {object} entity.Apod
// @Failure 500 {object} error
// @Router /apods [get]
func (h *Handler) RetrieveAllApods(ctx *gin.Context) {
	apods, err := h.apodService.RetrieveAllApods(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"apods": apods})
}

// RetrieveApodImageByDate retrieves an image for a specific APOD by date.
// @Summary Get APOD image by date
// @Description Retrieves the image for the Astronomy Picture of the Day for a specific date.
// @Produce octet-stream
// @Param date query string true "Date for the APOD image in the format YYYY-MM-DD"
// @Success 200 {string} file
// @Failure 500 {object} error
// @Router /img/apod [get]
func (h *Handler) RetrieveApodImageByDate(ctx *gin.Context) {
	date := ctx.Query("date")

	filePath, err := h.imageService.GetImage(ctx, date)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.File(filePath)

	err = os.Remove(filePath)
	if err != nil {
		h.log.Errorf("Error while deleting image: %s", err.Error())
		return
	}
}

// RetrieveApodByDate retrieves an APOD (Astronomy Picture of the Day) for a specific date.
// @Summary Get APOD by date
// @Description Retrieves the Astronomy Picture of the Day for a specific date.
// @Produce json
// @Param date query string true "Date for the APOD in the format YYYY-MM-DD"
// @Success 200 {object} entity.Apod
// @Failure 500 {object} error
// @Router /apod [get]
func (h *Handler) RetrieveApodByDate(ctx *gin.Context) {
	date := ctx.Query("date")

	apod, err := h.apodService.RetrieveApodByDate(ctx, date)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"apod": apod})
}
