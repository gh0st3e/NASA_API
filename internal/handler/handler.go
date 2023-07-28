package handler

import (
	"context"
	"fmt"
	"github.com/gh0st3e/NASA_API/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ServiceActions interface {
	RetrieveAllApods(ctx context.Context) ([]entity.Apod, error)
	RetrieveApodByDate(ctx context.Context, date string) (*entity.Apod, error)
}

type Handler struct {
	log     *logrus.Logger
	service ServiceActions
}

func NewHandler(log *logrus.Logger, service ServiceActions) *Handler {
	return &Handler{
		log:     log,
		service: service,
	}
}

func (h *Handler) Mount(r *gin.Engine) {
	api := r.Group("/api")

	api.GET("/apods", h.RetrieveAllApods)
	api.GET("/apod", h.RetrieveApodByDate)
}

func (h *Handler) RetrieveAllApods(ctx *gin.Context) {
	apods, err := h.service.RetrieveAllApods(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"apods": apods})
}

func (h *Handler) RetrieveApodByDate(ctx *gin.Context) {
	date := ctx.Query("date")

	fmt.Println(date)

	apod, err := h.service.RetrieveApodByDate(ctx, date)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"apod": apod})
}
