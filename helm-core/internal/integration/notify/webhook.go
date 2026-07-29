//go:build linux

package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/nayiswftw/helm/helm-core/internal/domain"
)

// WebhookNotifier sends notifications to an external webhook URL via HTTP POST.
// Compatible with Discord, Slack, ntfy.sh, Gotify, and any generic JSON webhook.
type WebhookNotifier struct {
	url        string
	httpClient *http.Client
}

// NewWebhookNotifier creates a new WebhookNotifier.
// If url is empty, the notifier reports as unavailable.
func NewWebhookNotifier(url string) *WebhookNotifier {
	return &WebhookNotifier{
		url: url,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// IsAvailable reports whether a webhook URL has been configured.
func (w *WebhookNotifier) IsAvailable() bool {
	return w.url != ""
}

// Send dispatches a notification to the configured webhook URL.
// The notification is sent as a JSON POST body.
func (w *WebhookNotifier) Send(ctx context.Context, notification domain.Notification) error {
	body, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("marshaling notification: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("creating webhook request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending webhook notification: %w", err)
	}
	defer resp.Body.Close()

	// Drain response body so connection can be reused.
	io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}
	return nil
}
