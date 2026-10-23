package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	keeperapi "cpa-usage-keeper/internal/api"
	"cpa-usage-keeper/internal/entities"
	"cpa-usage-keeper/internal/helper"
	"cpa-usage-keeper/internal/repository"
	"cpa-usage-keeper/internal/service"
)

func TestCPAAPIKeyOptionsExposeEveryInstanceWithPrefixedLabels(t *testing.T) {
	db := openAPITestDatabase(t)
	syncedAt := time.Date(2026, 10, 23, 10, 0, 0, 0, time.UTC)
	if err := repository.SyncCPAAPIKeysForInstance(db, "daichangyu", false, []string{"sk-daichangyu123456"}, syncedAt); err != nil {
		t.Fatalf("seed primary key: %v", err)
	}
	if err := repository.SyncCPAAPIKeysForInstance(db, "xuxiaozhi", false, []string{"sk-xuxiaozhi654321"}, syncedAt); err != nil {
		t.Fatalf("seed secondary key: %v", err)
	}
	router := keeperapi.NewRouter(nil, nil, nil, nil, keeperapi.AuthConfig{}, nil, "", keeperapi.OptionalProviders{CPAAPIKeys: service.NewCPAAPIKeyService(db)})

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/usage/api-keys/options", nil)
	router.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", resp.Code, resp.Body.String())
	}

	var parsed struct {
		Options []struct {
			ID       string `json:"id"`
			Instance string `json:"instance"`
			Label    string `json:"label"`
		} `json:"options"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &parsed); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(parsed.Options) != 2 {
		t.Fatalf("expected options from every instance, got %+v", parsed.Options)
	}
	if parsed.Options[0].Instance != "daichangyu" || parsed.Options[1].Instance != "xuxiaozhi" {
		t.Fatalf("expected instance tags in options, got %+v", parsed.Options)
	}
	if !strings.HasPrefix(parsed.Options[0].Label, "daichangyu · ") || !strings.HasPrefix(parsed.Options[1].Label, "xuxiaozhi · ") {
		t.Fatalf("expected instance prefix in labels, got %+v", parsed.Options)
	}
}

func TestCPAAPIKeyLabelKeepsCompatibilityWithoutInstance(t *testing.T) {
	row := entities.CPAAPIKey{APIKey: "sk-alpha123456", KeyAlias: "Primary Key"}
	if got := helper.CPAAPIKeyDisplayName(row); got != "Primary Key" {
		t.Fatalf("expected alias without instance prefix, got %q", got)
	}
	row.KeyAlias = ""
	if got := helper.CPAAPIKeyDisplayName(row); got != "sk-*********123456" {
		t.Fatalf("expected masked key without instance prefix, got %q", got)
	}
	row.KeyAlias = "Primary Key"
	row.Instance = "yangbing"
	if got := helper.CPAAPIKeyDisplayName(row); got != "yangbing · Primary Key" {
		t.Fatalf("expected instance prefix, got %q", got)
	}
}
