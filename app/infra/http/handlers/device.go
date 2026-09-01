package handlers

import (
	"net/http"

	"github.com/anoop-dryad/canopy/app/infra/http/dto"
	"github.com/anoop-dryad/canopy/app/internal/service"
	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	svc *service.Service
}

func NewDeviceHandler(svc *service.Service) *DeviceHandler {
	return &DeviceHandler{svc: svc}
}

// List godoc
// @Summary      List all devices
// @Description  Returns every device with a total count
// @Tags         devices
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "devices and count"
// @Failure      500  {object}  ErrorResponse
// @Router       /devices [get]
func (h *DeviceHandler) List(c *gin.Context) {
	devices, err := h.svc.List(c.Request.Context())
	if err != nil {
		writeError(c, err)
		return
	}
	dtos := make([]dto.DeviceDTO, len(devices))
	for i, d := range devices {
		dtos[i] = dto.ToDTO(d)
	}
	c.JSON(http.StatusOK, gin.H{"devices": dtos, "count": len(dtos)})
}

// Get godoc
// @Summary      Get a device by ID
// @Description  Returns a single device, or 404 if it does not exist
// @Tags         devices
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Device ID"
// @Success      200  {object}  dto.DeviceDTO
// @Failure      404  {object}  ErrorResponse
// @Router       /devices/{id} [get]
func (h *DeviceHandler) Get(c *gin.Context) {
	id := c.Param("id")
	d, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		writeError(c, err) // maps ErrNotFound → 404
		return
	}
	c.JSON(http.StatusOK, dto.ToDTO(d))
}

// Create godoc
// @Summary      Create a device
// @Description  Registers a new device. Server assigns id and timestamps.
// @Tags         devices
// @Accept       json
// @Produce      json
// @Param        device  body      dto.CreateDeviceRequest  true  "Device to create"
// @Success      201     {object}  dto.DeviceDTO
// @Failure      400     {object}  ErrorResponse
// @Failure      409     {object}  ErrorResponse
// @Router       /devices [post]
func (h *DeviceHandler) Create(c *gin.Context) {
	var req dto.CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeBadInput(c, "malformed JSON body")
		return
	}
	if msg, ok := req.Validate(); !ok {
		writeBadInput(c, msg)
		return
	}

	created, err := h.svc.Create(c.Request.Context(), req.ToDomain())
	if err != nil {
		writeError(c, err) // maps ErrDuplicate → 409
		return
	}
	c.JSON(http.StatusCreated, dto.ToDTO(created))
}

// Update godoc
// @Summary      Update a device
// @Description  Partially updates a device. Omitted fields are left unchanged.
// @Tags         devices
// @Accept       json
// @Produce      json
// @Param        id      path      string               true  "Device ID"
// @Param        device  body      dto.UpdateDeviceRequest  true  "Fields to update"
// @Success      200     {object}  dto.DeviceDTO
// @Failure      400     {object}  ErrorResponse
// @Failure      404     {object}  ErrorResponse
// @Router       /devices/{id} [put]
func (h *DeviceHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeBadInput(c, "malformed JSON body")
		return
	}
	if msg, ok := req.Validate(); !ok {
		writeBadInput(c, msg)
		return
	}

	updated, err := h.svc.Update(c.Request.Context(), id, req.ToFields())
	if err != nil {
		writeError(c, err) // maps ErrNotFound → 404
		return
	}
	c.JSON(http.StatusOK, dto.ToDTO(updated))
}

// Delete godoc
// @Summary      Delete a device
// @Description  Permanently removes a device. Returns 204 on success.
// @Tags         devices
// @Accept       json
// @Produce      json
// @Param        id   path  string  true  "Device ID"
// @Success      204  "No Content"
// @Failure      404  {object}  ErrorResponse
// @Router       /devices/{id} [delete]
func (h *DeviceHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		writeError(c, err) // maps ErrNotFound → 404
		return
	}
	c.Status(http.StatusNoContent) // 204, no body
}
