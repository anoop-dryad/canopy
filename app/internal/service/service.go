package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	apperror "github.com/anoop-dryad/canopy/app/internal/errors"
	"github.com/anoop-dryad/canopy/app/internal/models"
	"github.com/anoop-dryad/canopy/app/internal/repository"
)

type Service struct {
	repo repository.RepositoryInterface
	log  *zap.Logger
}

func NewService(repo repository.RepositoryInterface, log *zap.Logger) *Service {
	return &Service{
		repo: repo,
		log:  log.With(zap.String("domain", "canopy")), // scoped logger
	}
}

func (s *Service) List(ctx context.Context) ([]models.Device, error) {
	devices, err := s.repo.List(ctx)
	if err != nil {
		s.log.Error("list devices failed", zap.Error(err))
		return nil, err
	}
	return devices, nil
}

func (s *Service) Get(ctx context.Context, id string) (models.Device, error) {
	d, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrDeviceNotFound) {
			s.log.Debug("device not found", zap.String("id", id))
		} else {
			s.log.Error("get device failed", zap.String("id", id), zap.Error(err))
		}
		return models.Device{}, err
	}
	return d, nil
}

func (s *Service) Create(ctx context.Context, d models.Device) (models.Device, error) {
	now := time.Now().UTC()
	d.ID = "dev-" + uuid.NewString()
	d.CreatedAt = now
	d.UpdatedAt = now
	d.LastSeen = now // no telemetry yet; created time is a sensible default

	created, err := s.repo.Create(ctx, d)
	if err != nil {
		s.log.Error("create device failed", zap.String("id", d.ID), zap.Error(err))
		return models.Device{}, err
	}
	s.log.Info("device created", zap.String("id", created.ID))
	return created, nil
}

func (s *Service) Update(ctx context.Context, id string, fields models.UpdateFields) (models.Device, error) {
	updated, err := s.repo.Update(ctx, id, fields)
	if err != nil {
		s.log.Error("update device failed", zap.String("id", id), zap.Error(err))
		return models.Device{}, err
	}
	s.log.Info("device updated", zap.String("id", id))
	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.log.Error("delete device failed", zap.String("id", id), zap.Error(err))
		return err
	}
	s.log.Info("device deleted", zap.String("id", id))
	return nil
}
