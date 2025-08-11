package sanity

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWebhooksService_List(t *testing.T) {
	// Create a test server that returns a list of webhooks
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		if r.URL.Path != "/v2021-06-07/projects/test-project/webhooks" {
			t.Errorf("Expected /v2021-06-07/projects/test-project/webhooks path, got %s", r.URL.Path)
		}

		webhooks := []Webhook{
			{
				Id:            "webhook1",
				ProjectId:     "test-project",
				Dataset:       "production",
				URL:           "https://example.com/webhook",
				HttpMethod:    "POST",
				ApiVersion:    "v2021-06-07",
				IncludeDrafts: false,
				IsDisabled:    false,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(webhooks)
	}))
	defer ts.Close()

	// Create a client with the test server URL
	client := NewClient(http.DefaultClient)
	client.baseURL = ts.URL

	// Test the List method
	ctx := context.Background()
	webhooks, err := client.Webhooks.List(ctx, "test-project")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(webhooks) != 1 {
		t.Fatalf("Expected 1 webhook, got %d", len(webhooks))
	}

	webhook := webhooks[0]
	if webhook.Id != "webhook1" {
		t.Errorf("Expected webhook ID 'webhook1', got '%s'", webhook.Id)
	}
	if webhook.URL != "https://example.com/webhook" {
		t.Errorf("Expected URL 'https://example.com/webhook', got '%s'", webhook.URL)
	}
}

func TestWebhooksService_Create(t *testing.T) {
	// Create a test server that handles webhook creation
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		if r.URL.Path != "/v2021-06-07/projects/test-project/webhooks" {
			t.Errorf("Expected /v2021-06-07/projects/test-project/webhooks path, got %s", r.URL.Path)
		}

		// Parse the request body
		var req CreateWebhookRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		// Create a response webhook
		webhook := Webhook{
			Id:            "new-webhook",
			ProjectId:     "test-project",
			Dataset:       req.Dataset,
			URL:           req.URL,
			HttpMethod:    req.HttpMethod,
			ApiVersion:    req.ApiVersion,
			IncludeDrafts: false,
			IsDisabled:    false,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		if req.IncludeDrafts != nil {
			webhook.IncludeDrafts = *req.IncludeDrafts
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(webhook)
	}))
	defer ts.Close()

	// Create a client with the test server URL
	client := NewClient(http.DefaultClient)
	client.baseURL = ts.URL

	// Test the Create method
	ctx := context.Background()
	req := &CreateWebhookRequest{
		Dataset:       "production",
		URL:           "https://example.com/webhook",
		HttpMethod:    "POST",
		ApiVersion:    "v2021-06-07",
		IncludeDrafts: NewBool(true),
	}

	webhook, err := client.Webhooks.Create(ctx, "test-project", req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if webhook.Id != "new-webhook" {
		t.Errorf("Expected webhook ID 'new-webhook', got '%s'", webhook.Id)
	}
	if webhook.Dataset != "production" {
		t.Errorf("Expected dataset 'production', got '%s'", webhook.Dataset)
	}
	if webhook.URL != "https://example.com/webhook" {
		t.Errorf("Expected URL 'https://example.com/webhook', got '%s'", webhook.URL)
	}
	if !webhook.IncludeDrafts {
		t.Errorf("Expected IncludeDrafts to be true, got false")
	}
}

func TestClient_WebhooksService(t *testing.T) {
	// Test that the client properly initializes the Webhooks service
	client := NewClient(nil)

	if client.Webhooks == nil {
		t.Fatal("Expected Webhooks service to be initialized")
	}

	if client.Webhooks.client != client {
		t.Error("Expected Webhooks service to have reference to client")
	}
}