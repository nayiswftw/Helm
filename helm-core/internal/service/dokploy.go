//go:build linux

package service

import (
	"context"
	"errors"

	"github.com/nayiswftw/helm/helm-core/internal/domain"
	"github.com/nayiswftw/helm/helm-core/internal/integration/dokploy"
)

var ErrDokployUnavailable = errors.New("dokploy is not configured or unavailable")

// DokployService manages Dokploy project and application operations.
type DokployService struct {
	client        *dokploy.Client
	notifications *NotificationService
}

// NewDokployService creates a new DokployService.
func NewDokployService(client *dokploy.Client, notifications *NotificationService) *DokployService {
	return &DokployService{
		client:        client,
		notifications: notifications,
	}
}

// ListProjects returns all projects from the Dokploy instance.
func (s *DokployService) ListProjects(ctx context.Context) ([]domain.DokployProject, error) {
	if !s.client.IsAvailable() {
		return nil, ErrDokployUnavailable
	}
	return s.client.ListProjects(ctx)
}

// GetApplication returns details for a specific application.
func (s *DokployService) GetApplication(ctx context.Context, appID string) (domain.DokployApplication, error) {
	if !s.client.IsAvailable() {
		return domain.DokployApplication{}, ErrDokployUnavailable
	}
	return s.client.GetApplication(ctx, appID)
}

// Deploy triggers a deployment for the specified application.
func (s *DokployService) Deploy(ctx context.Context, appID string) error {
	if !s.client.IsAvailable() {
		return ErrDokployUnavailable
	}

	if err := s.client.DeployApplication(ctx, appID); err != nil {
		s.notifications.Notify(domain.EventDokployFailed, "Deployment failed for application "+appID, map[string]string{
			"application_id": appID,
			"operation":      "deploy",
		})
		return err
	}

	s.notifications.Notify(domain.EventDokployDeployed, "Deployment triggered for application "+appID, map[string]string{
		"application_id": appID,
		"operation":      "deploy",
	})
	return nil
}

// Redeploy triggers a redeployment for the specified application.
func (s *DokployService) Redeploy(ctx context.Context, appID string) error {
	if !s.client.IsAvailable() {
		return ErrDokployUnavailable
	}

	if err := s.client.RedeployApplication(ctx, appID); err != nil {
		s.notifications.Notify(domain.EventDokployFailed, "Redeployment failed for application "+appID, map[string]string{
			"application_id": appID,
			"operation":      "redeploy",
		})
		return err
	}

	s.notifications.Notify(domain.EventDokployDeployed, "Redeployment triggered for application "+appID, map[string]string{
		"application_id": appID,
		"operation":      "redeploy",
	})
	return nil
}

// ListDeployments returns deployment history for an application.
func (s *DokployService) ListDeployments(ctx context.Context, appID string) ([]domain.DokployDeployment, error) {
	if !s.client.IsAvailable() {
		return nil, ErrDokployUnavailable
	}
	return s.client.ListDeployments(ctx, appID)
}
