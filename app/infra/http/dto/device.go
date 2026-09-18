package dto

import (
	"time"

	"github.com/anoop-dryad/canopy/app/internal/models"
)

type DeviceDTO struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	BatteryPct int    `json:"battery_pct"`
	LastSeen   string `json:"last_seen"` // formatted RFC3339
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

func ToDTO(d models.Device) DeviceDTO {
	return DeviceDTO{
		ID:         d.ID,
		Name:       d.Name,
		Status:     string(d.Status),
		BatteryPct: d.BatteryPct,
		LastSeen:   d.LastSeen.Format(time.RFC3339),
		CreatedAt:  d.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  d.UpdatedAt.Format(time.RFC3339),
	}
}

// ---- inbound: create ----
type CreateDeviceRequest struct {
	Name       string `json:"name"`
	Status     string `json:"status"`
	BatteryPct int    `json:"battery_pct"`
}

func (r CreateDeviceRequest) Validate() (string, bool) {
	if r.Name == "" {
		return "name is required", false
	}
	if !validStatus(r.Status) {
		return "status must be one of: online, offline, degraded", false
	}
	if r.BatteryPct < 0 || r.BatteryPct > 100 {
		return "battery_pct must be between 0 and 100", false
	}
	return "", true
}

func (r CreateDeviceRequest) ToDomain() models.Device {
	return models.Device{
		Name:       r.Name,
		Status:     models.DeviceStatus(r.Status),
		BatteryPct: r.BatteryPct,
	}
}

// ---- inbound: update (pointers = partial) ----
type UpdateDeviceRequest struct {
	Name       *string `json:"name,omitempty"`
	Status     *string `json:"status,omitempty"`
	BatteryPct *int    `json:"battery_pct,omitempty"`
}

func (r UpdateDeviceRequest) Validate() (string, bool) {
	if r.Status != nil && !validStatus(*r.Status) {
		return "status must be one of: online, offline, degraded", false
	}
	if r.BatteryPct != nil && (*r.BatteryPct < 0 || *r.BatteryPct > 100) {
		return "battery_pct must be between 0 and 100", false
	}
	return "", true
}

func (r UpdateDeviceRequest) ToFields() models.UpdateFields {
	f := models.UpdateFields{Name: r.Name, BatteryPct: r.BatteryPct}
	if r.Status != nil {
		s := models.DeviceStatus(*r.Status)
		f.Status = &s
	}
	return f
}

func validStatus(s string) bool {
	switch models.DeviceStatus(s) {
	case models.StatusOnline, models.StatusOffline, models.StatusDegraded:
		return true
	}
	return false
}

type DownlinkDTO struct {
	ID        string `json:"id"`
	Command   string `json:"name"`
	Status    string `json:"status"`
	DeviceID  string `json:"battery_pct"`
	CreatedAt string `json:"created_at"`
}

func ToDownlinkDTO(d models.Downlink) DownlinkDTO {
	return DownlinkDTO{
		ID:        d.ID,
		DeviceID:  d.DeviceID,
		Command:   d.Command,
		Status:    string(d.Status),
		CreatedAt: d.CreatedAt.Format(time.RFC3339),
	}
}

// ---- inbound: create downlink ----
type CreateDownlinkRequest struct {
	Command string `json:"command"` // "reset", "calibrate", ...
	// DeviceID comes from the URL path (/devices/{id}/downlink), not the body
}

func (r CreateDownlinkRequest) Validate() (string, bool) {
	if r.Command == "" {
		return "command is required", false
	}
	return "", true
}
