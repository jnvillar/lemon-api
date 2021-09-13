package handler

import (
	"net/http"

	"lemonapp/logger"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/heartbeat", func(c *gin.Context) { h.heartbeat(c) })
	router.GET("/ping", func(c *gin.Context) { h.ping(c) })
}

func (h *HealthHandler) ping(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(http.StatusOK)
	_, err := c.Writer.Write([]byte("pong"))
	if err != nil {
		logger.Sugar.Errorf("ping error %v", err)
	}
}

func (h *HealthHandler) heartbeat(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(http.StatusOK)
	_, err := c.Writer.Write([]byte("{}"))
	if err != nil {
		logger.Sugar.Errorf("heartbeat error %v", err)
	}
}
