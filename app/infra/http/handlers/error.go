package handlers

import (
	"errors"
	"net/http"

	apperror "github.com/anoop-dryad/canopy/app/internal/errors"
	"github.com/gin-gonic/gin"
)

// ErrorResponse — every non-2xx response returns this shape.
// The agent's gate branches on `code`, so keep the codes stable.
type ErrorResponse struct {
	Error string `json:"error"`        // human-readable, SAFE message (no internal detail)
	Code  string `json:"code"`         // machine-stable code for branching
	ID    string `json:"id,omitempty"` // the resource ID involved, if any
}

// writeError maps a domain error to the wire contract.
// THE single mapping point. Add a new error type here and nowhere else.
func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, apperror.ErrDeviceNotFound):
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "device not found",
			Code:  "DEVICE_NOT_FOUND",
			ID:    c.Param("id"),
		})
	case errors.Is(err, apperror.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid input",
			Code:  "INVALID_INPUT",
		})
	case errors.Is(err, apperror.ErrDuplicate):
		c.JSON(http.StatusConflict, ErrorResponse{
			Error: "device already exists",
			Code:  "DUPLICATE_DEVICE",
			ID:    c.Param("id"),
		})
	default:
		// Never leak internal error detail to the client.
		// Log the real err upstream (in the service); return a generic message here.
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "internal error",
			Code:  "INTERNAL",
		})
	}
}

// writeBadInput handles validation failures that are NOT domain errors —
// just a message string from DTO.Validate(). Kept separate because these
// carry a specific, safe-to-show reason ("battery_pct must be 0-100").
func writeBadInput(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error: msg,
		Code:  "INVALID_INPUT",
	})
}
