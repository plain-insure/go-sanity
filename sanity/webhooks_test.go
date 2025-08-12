package sanity

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestWebhooksService_List(t *testing.T) {
	// Create a test server that returns a list of webhooks
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		if r.URL.Path != "/hooks/projects/test-project" {
			t.Errorf("Expected /hooks/projects/test-project path, got %s", r.URL.Path)
		}

		webhooks := []Webhook{
			{
				Id:            "webhook1",
				ProjectId:     "test-project",
				Type:          "document",
				Name:          "Test Webhook",
				Dataset:       "production",
				URL:           "https://example.com/webhook",
				HttpMethod:    "POST",
				ApiVersion:    "v2025-02-19",
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

	// Create a client and set test base URL
	client := NewClient(http.DefaultClient)
	client.Webhooks.testBaseURL = ts.URL

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
		if r.URL.Path != "/hooks/projects/test-project" {
			t.Errorf("Expected /hooks/projects/test-project path, got %s", r.URL.Path)
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
			Type:          req.Type,
			Name:          req.Name,
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

	// Create a client and set test base URL
	client := NewClient(http.DefaultClient)
	client.Webhooks.testBaseURL = ts.URL

	// Test the Create method
	ctx := context.Background()
	req := &CreateWebhookRequest{
		Type:          "document",
		Name:          "Test Webhook",
		Dataset:       "production",
		URL:           "https://example.com/webhook",
		HttpMethod:    "POST",
		ApiVersion:    "v2025-02-19",
		IncludeDrafts: NewBool(true),
	}

	webhook, err := client.Webhooks.Create(ctx, "test-project", req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if webhook.Id != "new-webhook" {
		t.Errorf("Expected webhook ID 'new-webhook', got '%s'", webhook.Id)
	}
	if webhook.Type != "document" {
		t.Errorf("Expected webhook type 'document', got '%s'", webhook.Type)
	}
	if webhook.Name != "Test Webhook" {
		t.Errorf("Expected webhook name 'Test Webhook', got '%s'", webhook.Name)
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

func TestCreateWebhookRequest_RequiredFields(t *testing.T) {
	// Test that CreateWebhookRequest includes all required fields
	req := &CreateWebhookRequest{
		Type:    "document",
		Name:    "Required Webhook Name",
		Dataset: "production",
		URL:     "https://example.com/webhook",
	}

	// Verify that required fields are present and accessible
	if req.Type != "document" {
		t.Errorf("Expected type field to be 'document', got '%s'", req.Type)
	}
	if req.Name != "Required Webhook Name" {
		t.Errorf("Expected name field to be 'Required Webhook Name', got '%s'", req.Name)
	}
	if req.Dataset != "production" {
		t.Errorf("Expected dataset field to be 'production', got '%s'", req.Dataset)
	}
	if req.URL != "https://example.com/webhook" {
		t.Errorf("Expected url field to be 'https://example.com/webhook', got '%s'", req.URL)
	}

	// Test JSON marshalling includes all required fields
	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal CreateWebhookRequest: %v", err)
	}

	jsonStr := string(jsonData)
	if !strings.Contains(jsonStr, `"type":"document"`) {
		t.Errorf("Expected JSON to contain type field, got: %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"name":"Required Webhook Name"`) {
		t.Errorf("Expected JSON to contain name field, got: %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"dataset":"production"`) {
		t.Errorf("Expected JSON to contain dataset field, got: %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"url":"https://example.com/webhook"`) {
		t.Errorf("Expected JSON to contain url field, got: %s", jsonStr)
	}
}

func TestUpdateWebhookRequest_NameField(t *testing.T) {
	// Test that UpdateWebhookRequest includes the name field
	req := &UpdateWebhookRequest{
		Name: "Updated Webhook Name",
		URL:  "https://example.com/updated-webhook",
	}

	// Verify that name field is present and accessible
	if req.Name != "Updated Webhook Name" {
		t.Errorf("Expected name field to be 'Updated Webhook Name', got '%s'", req.Name)
	}

	// Test JSON marshalling includes the name field
	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal UpdateWebhookRequest: %v", err)
	}

	jsonStr := string(jsonData)
	if !strings.Contains(jsonStr, `"name":"Updated Webhook Name"`) {
		t.Errorf("Expected JSON to contain name field, got: %s", jsonStr)
	}
}

func TestWebhookRule_Structure(t *testing.T) {
	// Test that WebhookRule has the correct structure
	rule := &WebhookRule{
		On:         []string{"create", "update"},
		Filter:     "_type == 'post'",
		Projection: "{title, slug}",
	}

	// Verify fields are accessible
	if len(rule.On) != 2 || rule.On[0] != "create" || rule.On[1] != "update" {
		t.Errorf("Expected On field to be ['create', 'update'], got %v", rule.On)
	}
	if rule.Filter != "_type == 'post'" {
		t.Errorf("Expected Filter field to be '_type == 'post'', got '%s'", rule.Filter)
	}
	if rule.Projection != "{title, slug}" {
		t.Errorf("Expected Projection field to be '{title, slug}', got '%s'", rule.Projection)
	}

	// Test JSON marshalling
	jsonData, err := json.Marshal(rule)
	if err != nil {
		t.Fatalf("Failed to marshal WebhookRule: %v", err)
	}

	jsonStr := string(jsonData)
	if !strings.Contains(jsonStr, `"on":["create","update"]`) {
		t.Errorf("Expected JSON to contain on field, got: %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"filter":"_type == 'post'"`) {
		t.Errorf("Expected JSON to contain filter field, got: %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"projection":"{title, slug}"`) {
		t.Errorf("Expected JSON to contain projection field, got: %s", jsonStr)
	}
}

func TestCreateWebhookRequest_WithRule(t *testing.T) {
	// Test CreateWebhookRequest with Rule
	rule := &WebhookRule{
		On:         []string{"create"},
		Filter:     "_type == 'article'",
		Projection: "{title, _id}",
	}

	req := &CreateWebhookRequest{
		Type:             "document",
		Name:             "Test Webhook with Rule",
		Dataset:          "production",
		URL:              "https://example.com/webhook",
		Rule:             rule,
		IsDisabledByUser: NewBool(false),
	}

	// Test JSON marshalling includes rule and isDisabledByUser
	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal CreateWebhookRequest: %v", err)
	}

	jsonStr := string(jsonData)
	if !strings.Contains(jsonStr, `"rule":{`) {
		t.Errorf("Expected JSON to contain rule field, got: %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"isDisabledByUser":false`) {
		t.Errorf("Expected JSON to contain isDisabledByUser field, got: %s", jsonStr)
	}
}

func TestUpdateWebhookRequest_WithoutIsDisabled(t *testing.T) {
	// Test that UpdateWebhookRequest does not include IsDisabled field
	req := &UpdateWebhookRequest{
		Name:             "Updated Webhook",
		IsDisabledByUser: NewBool(true),
	}

	// Test JSON marshalling
	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal UpdateWebhookRequest: %v", err)
	}

	jsonStr := string(jsonData)
	// Should contain isDisabledByUser but not isDisabled
	if !strings.Contains(jsonStr, `"isDisabledByUser":true`) {
		t.Errorf("Expected JSON to contain isDisabledByUser field, got: %s", jsonStr)
	}
	if strings.Contains(jsonStr, `"isDisabled"`) {
		t.Errorf("Expected JSON to NOT contain isDisabled field, got: %s", jsonStr)
	}
}
