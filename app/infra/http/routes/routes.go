package routes

import (
	"log"
	"os"

	"github.com/anoop-dryad/canopy/app/infra/http/handlers"
	"github.com/anoop-dryad/canopy/app/infra/http/middleware"
	"github.com/anoop-dryad/canopy/app/infra/http/swagger"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine, deps handlers.Dependencies, isProduction bool) {
	apiKey := os.Getenv("CANOPY_API_KEY")
	if apiKey == "" {
		log.Fatal("CANOPY_API_KEY not set")
	}

	v1 := engine.Group("/v1")
	device := v1.Group("/devices")
	device.Use(middleware.APIKeyAuth(apiKey))

	Health(v1)
	Device(device, deps.DeviceHandler)

	if !isProduction {
		swagger.Register(engine, "v1")
	}
}
