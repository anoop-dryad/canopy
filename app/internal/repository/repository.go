package repository

import (
	"context"
	"database/sql"
	"errors"

	apperror "github.com/anoop-dryad/canopy/app/internal/errors"
	"github.com/anoop-dryad/canopy/app/internal/models"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

type RepositoryInterface interface {
	List(ctx context.Context) ([]models.Device, error)
	GetByID(ctx context.Context, id string) (models.Device, error)
	Create(ctx context.Context, d models.Device) (models.Device, error)
	Update(ctx context.Context, id string, fields models.UpdateFields) (models.Device, error)
	Delete(ctx context.Context, id string) error
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context) ([]models.Device, error) {
	var devices []models.Device
	err := r.db.SelectContext(ctx, &devices,
		`SELECT id, name, status, battery_pct, last_seen, created_at, updated_at
		 FROM devices ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (models.Device, error) {
	var d models.Device
	err := r.db.GetContext(ctx, &d,
		`SELECT id, name, status, battery_pct, last_seen, created_at, updated_at
		 FROM devices WHERE id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Device{}, apperror.ErrDeviceNotFound
	}
	if err != nil {
		return models.Device{}, err
	}
	return d, nil
}

func (r *Repository) Create(ctx context.Context, d models.Device) (models.Device, error) {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO devices (id, name, status, battery_pct, last_seen, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		d.ID, d.Name, d.Status, d.BatteryPct, d.LastSeen, d.CreatedAt, d.UpdatedAt)
	if err != nil {
		return models.Device{}, err // TODO: map 23505 unique violation → ErrDuplicate
	}
	return d, nil
}

func (r *Repository) Update(ctx context.Context, id string, f models.UpdateFields) (models.Device, error) {
	var d models.Device
	err := r.db.GetContext(ctx, &d,
		`UPDATE devices SET
			name        = COALESCE($2, name),
			status      = COALESCE($3, status),
			battery_pct = COALESCE($4, battery_pct),
			updated_at  = now()
		 WHERE id = $1
		 RETURNING id, name, status, battery_pct, last_seen, created_at, updated_at`,
		id, f.Name, f.Status, f.BatteryPct)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Device{}, apperror.ErrDeviceNotFound
	}
	if err != nil {
		return models.Device{}, err
	}
	return d, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM devices WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperror.ErrDeviceNotFound
	}
	return nil
}

var _ RepositoryInterface = (*Repository)(nil)
