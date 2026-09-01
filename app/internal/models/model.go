package models

import "time"

// DeviceStatus is a constrained set — not a free string.
// This matters: the agent's gate logic branches on status,
// so the set of possible values must be known and finite.
type DeviceStatus string

const (
	StatusOnline   DeviceStatus = "online"
	StatusOffline  DeviceStatus = "offline"
	StatusDegraded DeviceStatus = "degraded"
)

// Device is the core resource.
type Device struct {
	ID         string       `db:"id"          json:"id"`
	Name       string       `db:"name"        json:"name"`
	Status     DeviceStatus `db:"status"      json:"status"`
	BatteryPct int          `db:"battery_pct" json:"battery_pct"` // 0–100
	LastSeen   time.Time    `db:"last_seen"   json:"last_seen"`   // RFC3339 in JSON
	CreatedAt  time.Time    `db:"created_at"  json:"created_at"`
	UpdatedAt  time.Time    `db:"updated_at"  json:"updated_at"`
}

type UpdateFields struct {
	Name       *string       // nil = don't change; non-nil = set to this
	Status     *DeviceStatus // nil = don't change
	BatteryPct *int          // nil = don't change
}
