//go:build linux

package service

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/nayiswftw/helm/helm-core/internal/domain"
	"github.com/nayiswftw/helm/helm-core/internal/integration/system"
)

var (
	ErrActionNotFound = errors.New("action not found")
)

// ActionService manages predefined action registry and execution.
type ActionService struct {
	system        *system.System
	logger        *slog.Logger
	notifications *NotificationService
	actions       map[string]domain.Action
	deviceID      string
}

// NewActionService creates an ActionService initialized with safe system actions.
func NewActionService(sys *system.System, localDeviceID string, logger *slog.Logger, notifications *NotificationService) *ActionService {
	actions := map[string]domain.Action{
		"reboot": {
			ID:          "reboot",
			Name:        "Reboot Server",
			Description: "Gracefully reboots the physical or virtual server",
			DeviceID:    localDeviceID,
			Dangerous:   true,
		},
		"shutdown": {
			ID:          "shutdown",
			Name:        "Shutdown Server",
			Description: "Powers off the server immediately",
			DeviceID:    localDeviceID,
			Dangerous:   true,
		},
	}

	return &ActionService{
		system:        sys,
		logger:        logger,
		notifications: notifications,
		actions:       actions,
		deviceID:      localDeviceID,
	}
}

// ListActions returns all predefined administrative actions.
func (s *ActionService) ListActions() []domain.Action {
	result := make([]domain.Action, 0, len(s.actions))
	for _, act := range s.actions {
		result = append(result, act)
	}
	return result
}

// Execute triggers a predefined action securely.
func (s *ActionService) Execute(actionID string) (domain.ActionResult, error) {
	act, ok := s.actions[actionID]
	if !ok {
		return domain.ActionResult{}, ErrActionNotFound
	}

	s.logger.Warn("executing action", "action_id", actionID, "device_id", act.DeviceID)

	switch actionID {
	case "reboot":
		// Asynchronous execution so HTTP response can return 202 Accepted
		go func() {
			if err := s.system.Reboot(); err != nil {
				s.logger.Error("reboot action failed", "error", err)
			}
		}()

		s.notifications.Notify(domain.EventActionExecuted, "Reboot initiated on "+act.DeviceID, map[string]string{
			"action_id": actionID,
			"device_id": act.DeviceID,
		})

		return domain.ActionResult{
			ActionID: actionID,
			Status:   "accepted",
			Message:  "Reboot sequence initiated",
		}, nil

	case "shutdown":
		go func() {
			if err := s.system.Shutdown(); err != nil {
				s.logger.Error("shutdown action failed", "error", err)
			}
		}()

		s.notifications.Notify(domain.EventActionExecuted, "Shutdown initiated on "+act.DeviceID, map[string]string{
			"action_id": actionID,
			"device_id": act.DeviceID,
		})

		return domain.ActionResult{
			ActionID: actionID,
			Status:   "accepted",
			Message:  "Shutdown sequence initiated",
		}, nil

	default:
		return domain.ActionResult{}, fmt.Errorf("unhandled action execution: %s", actionID)
	}
}

