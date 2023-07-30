package handler

import (
	"context"
	"net/http"
	"os"

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

func (h *Handler) Mount(r *gin.Engine) {
	api := r.Group("/api")

	api.GET("/apods", h.RetrieveAllApods)
	api.GET("/apod", h.RetrieveApodByDate)

	img := api.Group("img")

	img.GET("/apod", h.RetrieveApodImageByDate)

}

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
