//go:build linux

package service

import (
	"context"
	"errors"

	"github.com/nayiswftw/helm/helm-core/internal/domain"
	"github.com/nayiswftw/helm/helm-core/internal/integration/docker"
)

var ErrDockerUnavailable = errors.New("docker is not available or socket is missing")

// ContainerService manages container business logic and delegates to Docker integration.
type ContainerService struct {
	dockerClient *docker.Client
}

// NewContainerService creates a new ContainerService instance.
func NewContainerService(dockerClient *docker.Client) *ContainerService {
	return &ContainerService{
		dockerClient: dockerClient,
	}
}

// List returns all Docker containers.
func (s *ContainerService) List(ctx context.Context) ([]domain.Container, error) {
	if !s.dockerClient.IsAvailable() {
		return nil, ErrDockerUnavailable
	}
	return s.dockerClient.ListContainers(ctx)
}

// Start starts a container by ID or name.
func (s *ContainerService) Start(ctx context.Context, id string) error {
	if !s.dockerClient.IsAvailable() {
		return ErrDockerUnavailable
	}
	return s.dockerClient.StartContainer(ctx, id)
}

// Stop stops a container by ID or name.
func (s *ContainerService) Stop(ctx context.Context, id string) error {
	if !s.dockerClient.IsAvailable() {
		return ErrDockerUnavailable
	}
	return s.dockerClient.StopContainer(ctx, id)
}

// Restart restarts a container by ID or name.
func (s *ContainerService) Restart(ctx context.Context, id string) error {
	if !s.dockerClient.IsAvailable() {
		return ErrDockerUnavailable
	}
	return s.dockerClient.RestartContainer(ctx, id)
}
