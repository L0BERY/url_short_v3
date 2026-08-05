package handler

import (
	"errors"
	"net/http"
	"url_shortener_v3/internal/dto"
	"url_shortener_v3/internal/service"

	"github.com/gin-gonic/gin"
)

type HandlerManager struct {
	srv service.Service
}

func NewHandlerManager(srv service.Service) *HandlerManager {
	return &HandlerManager{srv: srv}
}

func (h *HandlerManager) Shorten(ctx *gin.Context) {
	var req dto.AddNewURLRequests

	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	code, err := h.srv.AddNewURL(ctx.Request.Context(), req.URL)
	if errors.Is(err, service.ErrEmptyUrl) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	} else if errors.Is(err, service.ErrTooManyAttempts) || errors.Is(err, service.ErrServer) {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": code})
}

func (h *HandlerManager) Redirect(ctx *gin.Context) {
	code := ctx.Param("code")

	url, err := h.srv.GetURL(ctx.Request.Context(), code)
	if errors.Is(err, service.ErrUrlNotFound) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	} else if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	ctx.Redirect(http.StatusMovedPermanently, url)
}
