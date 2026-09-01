package routes

import (
	"github.com/anoop-dryad/canopy/app/infra/http/handlers"
	"github.com/gin-gonic/gin"
)

func Device(router *gin.RouterGroup, h *handlers.DeviceHandler) {
	router.GET("", h.List)
	router.GET("/:id", h.Get)
	router.POST("", h.Create)
	router.PUT("/:id", h.Update)
	router.DELETE("/:id", h.Delete)
}
