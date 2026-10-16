package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	keeperapp "cpa-usage-keeper/internal/app"
	"cpa-usage-keeper/internal/entities"
)

func TestPricingUpgradeCompatibilityBlocksBeforeM1AndBusinessDDL(t *testing.T) {
	cfg := testAppConfig(t)
	seed := seedPublishedPricingRuntime(t, cfg)
	closePricingStartupPools(t, seed, seed)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v0/management/debug" {
			t.Errorf("unexpected compatibility request: %s", r.URL)
		}
		w.Header().Set("X-CPA-VERSION", "v7.9.9")
	}))
	defer server.Close()
	cfg.CPABaseURL = server.URL
	application, err := keeperapp.NewWithConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = application.Close() })
	if err := application.Initialize(context.Background()); err == nil || !strings.Contains(err.Error(), "minimum CPA 8.0.12") {
		t.Fatalf("incompatible CPA did not stop before migration: %v", err)
	}
	var state entities.PricingMigrationState
	if err := application.DB.Take(&state).Error; err != nil || state.BackupPath != nil || state.BaselineJSON != nil || state.SchemaComplete || state.DataComplete {
		t.Fatalf("incompatible CPA crossed M1: %+v error=%v", state, err)
	}
	if application.DB.Migrator().HasColumn("usage_events", "cost_usd") || application.Router != nil {
		t.Fatal("incompatible CPA reached destructive migration or business initialization")
	}
	if entries, err := os.ReadDir(cfg.BackupDir); err == nil && len(entries) != 0 {
		t.Fatal("incompatible CPA created upgrade backup")
	}
}

func TestPricingUpgradeCompatibilityDoesNotProbeFreshOrCompletedDatabase(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { requests++; w.Header().Set("X-CPA-VERSION", "v7.9.9") }))
	defer server.Close()
	cfg := testAppConfig(t)
	cfg.CPABaseURL = server.URL
	for range 2 {
		application, err := keeperapp.NewWithConfig(cfg)
		if err != nil {
			t.Fatal(err)
		}
		if err := application.Initialize(context.Background()); err != nil {
			t.Fatalf("fresh/completed startup gained compatibility network dependency: %v", err)
		}
		if err := application.Close(); err != nil {
			t.Fatal(err)
		}
	}
	if requests != 0 {
		t.Fatalf("fresh/completed startup made %d CPA probes", requests)
	}
}
