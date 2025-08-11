package sanity

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// WebhooksService is a client for the Sanity Webhooks API.
//
// Refer to https://www.sanity.io/docs/webhooks for more information.
type WebhooksService struct {
	service
	// testBaseURL is used for testing to override the default URL construction
	testBaseURL string
}

// getWebhookBaseURL returns the base URL for webhook operations for a given project.
func (s *WebhooksService) getWebhookBaseURL(projectId string) string {
	if s.testBaseURL != "" {
		return s.testBaseURL
	}
	return fmt.Sprintf("https://%s.api.sanity.io/v2025-02-19", projectId)
}

// A Webhook represents a webhook configuration for a Sanity project.
type Webhook struct {
	// Id is the unique identifier for the webhook.
	Id string `json:"id"`

	// ProjectId is the identifier of the project this webhook belongs to.
	ProjectId string `json:"projectId"`

	// Dataset is the dataset this webhook is configured for.
	Dataset string `json:"dataset"`

	// URL is the endpoint that will receive webhook notifications.
	URL string `json:"url"`

	// HttpMethod is the HTTP method used for webhook requests (typically POST).
	HttpMethod string `json:"httpMethod"`

	// ApiVersion is the API version used for webhook payloads.
	ApiVersion string `json:"apiVersion"`

	// IncludeDrafts indicates whether draft documents trigger webhook notifications.
	IncludeDrafts bool `json:"includeDrafts"`

	// Headers are custom HTTP headers sent with webhook requests.
	Headers map[string]string `json:"headers,omitempty"`

	// Filter is a GROQ filter expression to determine which documents trigger the webhook.
	Filter string `json:"filter,omitempty"`

	// CreatedAt is the time the webhook was created.
	CreatedAt time.Time `json:"createdAt"`

	// UpdatedAt is the time the webhook was last updated.
	UpdatedAt time.Time `json:"updatedAt"`

	// Secret is used for webhook signature verification.
	Secret string `json:"secret,omitempty"`

	// IsDisabled indicates whether the webhook is currently disabled.
	IsDisabled bool `json:"isDisabled"`
}

// CreateWebhookRequest represents the payload for creating a new webhook.
type CreateWebhookRequest struct {
	// Dataset is the dataset this webhook is configured for.
	Dataset string `json:"dataset"`

	// URL is the endpoint that will receive webhook notifications.
	URL string `json:"url"`

	// HttpMethod is the HTTP method used for webhook requests (typically POST).
	HttpMethod string `json:"httpMethod,omitempty"`

	// ApiVersion is the API version used for webhook payloads.
	ApiVersion string `json:"apiVersion,omitempty"`

	// IncludeDrafts indicates whether draft documents trigger webhook notifications.
	IncludeDrafts *bool `json:"includeDrafts,omitempty"`

	// Headers are custom HTTP headers sent with webhook requests.
	Headers map[string]string `json:"headers,omitempty"`

	// Filter is a GROQ filter expression to determine which documents trigger the webhook.
	Filter string `json:"filter,omitempty"`

	// Secret is used for webhook signature verification.
	Secret string `json:"secret,omitempty"`
}

// UpdateWebhookRequest represents the payload for updating an existing webhook.
type UpdateWebhookRequest struct {
	// URL is the endpoint that will receive webhook notifications.
	URL string `json:"url,omitempty"`

	// HttpMethod is the HTTP method used for webhook requests.
	HttpMethod string `json:"httpMethod,omitempty"`

	// ApiVersion is the API version used for webhook payloads.
	ApiVersion string `json:"apiVersion,omitempty"`

	// IncludeDrafts indicates whether draft documents trigger webhook notifications.
	IncludeDrafts *bool `json:"includeDrafts,omitempty"`

	// Headers are custom HTTP headers sent with webhook requests.
	Headers map[string]string `json:"headers,omitempty"`

	// Filter is a GROQ filter expression to determine which documents trigger the webhook.
	Filter string `json:"filter,omitempty"`

	// Secret is used for webhook signature verification.
	Secret string `json:"secret,omitempty"`

	// IsDisabled indicates whether the webhook is currently disabled.
	IsDisabled *bool `json:"isDisabled,omitempty"`
}

// List fetches and returns all webhooks for the specified project.
func (s *WebhooksService) List(ctx context.Context, projectId string) ([]Webhook, error) {
	url := fmt.Sprintf("%s/hooks/projects/%s", s.getWebhookBaseURL(projectId), projectId)

	var webhooks []Webhook
	err := do(ctx, s.client.client, url, http.MethodGet, nil, &webhooks)

	return webhooks, err
}

// Create generates a new webhook for the specified project.
func (s *WebhooksService) Create(ctx context.Context, projectId string, r *CreateWebhookRequest) (*Webhook, error) {
	url := fmt.Sprintf("%s/hooks/projects/%s", s.getWebhookBaseURL(projectId), projectId)

	var webhook Webhook
	err := do(ctx, s.client.client, url, http.MethodPost, r, &webhook)

	return &webhook, err
}

// Get fetches a webhook by its unique identifier.
func (s *WebhooksService) Get(ctx context.Context, projectId, webhookId string) (*Webhook, error) {
	url := fmt.Sprintf("%s/hooks/projects/%s/%s", s.getWebhookBaseURL(projectId), projectId, webhookId)

	var webhook Webhook
	err := do(ctx, s.client.client, url, http.MethodGet, nil, &webhook)

	return &webhook, err
}

// Update applies the requested changes to the specified webhook.
func (s *WebhooksService) Update(ctx context.Context, projectId, webhookId string, r *UpdateWebhookRequest) (*Webhook, error) {
	url := fmt.Sprintf("%s/hooks/projects/%s/%s", s.getWebhookBaseURL(projectId), projectId, webhookId)

	var webhook Webhook
	err := do(ctx, s.client.client, url, http.MethodPatch, r, &webhook)

	return &webhook, err
}

// Delete removes the specified webhook without prompt.
func (s *WebhooksService) Delete(ctx context.Context, projectId, webhookId string) (bool, error) {
	url := fmt.Sprintf("%s/hooks/projects/%s/%s", s.getWebhookBaseURL(projectId), projectId, webhookId)

	type response struct {
		Deleted bool `json:"deleted"`
	}

	var resp response
	err := do(ctx, s.client.client, url, http.MethodDelete, nil, &resp)
	return resp.Deleted, err
}

