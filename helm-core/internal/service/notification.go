//go:build linux

package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/nayiswftw/helm/helm-core/internal/domain"
	"github.com/nayiswftw/helm/helm-core/internal/integration/notify"
)

// NotificationService dispatches event notifications to configured channels.
// Notifications are sent asynchronously to avoid blocking callers.
type NotificationService struct {
	webhook *notify.WebhookNotifier
	logger  *slog.Logger
}

// NewNotificationService creates a new NotificationService.
func NewNotificationService(webhook *notify.WebhookNotifier, logger *slog.Logger) *NotificationService {
	return &NotificationService{
		webhook: webhook,
		logger:  logger,
	}
}

// Notify dispatches a notification asynchronously.
// Failures are logged but never propagated to the caller.
func (s *NotificationService) Notify(event string, message string, metadata map[string]string) {
	if !s.webhook.IsAvailable() {
		return
	}

	notification := domain.Notification{
		Event:     event,
		Message:   message,
		Timestamp: time.Now().UTC(),
		Metadata:  metadata,
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.webhook.Send(ctx, notification); err != nil {
			s.logger.Warn("failed to send notification",
				"event", event,
				"error", err,
			)
		} else {
			s.logger.Debug("notification sent",
				"event", event,
				"message", message,
			)
		}
	}()
}

// SendTest dispatches a test notification synchronously for debugging.
func (s *NotificationService) SendTest() error {
	if !s.webhook.IsAvailable() {
		return ErrNotificationsUnavailable
	}

	notification := domain.Notification{
		Event:     "test",
		Message:   "Helm notification test — if you see this, it works!",
		Timestamp: time.Now().UTC(),
		Metadata: map[string]string{
			"source": "helm",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.webhook.Send(ctx, notification)
}

// ErrNotificationsUnavailable is returned when no notification channel is configured.
var ErrNotificationsUnavailable = errNotificationsUnavailable("notifications are not configured")

type errNotificationsUnavailable string

func (e errNotificationsUnavailable) Error() string { return string(e) }
