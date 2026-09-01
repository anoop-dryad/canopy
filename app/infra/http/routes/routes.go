package routes

import (
	"github.com/anoop-dryad/canopy/app/infra/http/handlers"
	"github.com/anoop-dryad/canopy/app/infra/http/swagger"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine, deps handlers.Dependencies, isProduction bool) {
	v1 := engine.Group("/v1")
	device := v1.Group("/devices")

	Health(v1)
	Device(device, deps.DeviceHandler)

	if !isProduction {
		swagger.Register(engine, "v1")
	}
}
