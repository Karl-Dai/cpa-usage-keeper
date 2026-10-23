package test

import (
	. "cpa-usage-keeper/internal/config"
	"strings"
	"testing"
)

func TestLoadParsesMultipleCPAInstancesAndFillsPrimaryFields(t *testing.T) {
	withIsolatedEnvFiles(t)
	t.Setenv("CPA_INSTANCES", `[
		{"name":"daichangyu","base_url":"http://127.0.0.1:8317","management_key":"key-a"},
		{"name":"xuxiaozhi","base_url":"http://127.0.0.1:8318","management_key":"key-b","redis_queue_addr":"127.0.0.1:8319"}
	]`)

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv returned error: %v", err)
	}

	if len(cfg.Instances) != 2 {
		t.Fatalf("expected two instances, got %+v", cfg.Instances)
	}
	if !cfg.MultiInstance() {
		t.Fatalf("expected multi instance mode")
	}
	primary := cfg.PrimaryInstance()
	if primary.Name != "daichangyu" || primary.BaseURL != "http://127.0.0.1:8317" || primary.ManagementKey != "key-a" {
		t.Fatalf("unexpected primary instance: %+v", primary)
	}
	// 未显式配置单实例字段时，必须用第一个实例回填，保证配额/凭证等仍在主实例工作。
	if cfg.CPABaseURL != "http://127.0.0.1:8317" || cfg.CPAManagementKey != "key-a" {
		t.Fatalf("expected primary instance to fill legacy fields, got base=%q key=%q", cfg.CPABaseURL, cfg.CPAManagementKey)
	}
	if cfg.Instances[1].RedisQueueAddr != "127.0.0.1:8319" {
		t.Fatalf("expected per-instance redis queue addr, got %+v", cfg.Instances[1])
	}
}

func TestLoadCPAInstancesKeepsExplicitSingleInstanceFields(t *testing.T) {
	withIsolatedEnvFiles(t)
	t.Setenv("CPA_BASE_URL", "http://legacy.example.com:8317")
	t.Setenv("CPA_MANAGEMENT_KEY", "legacy-key")
	t.Setenv("CPA_INSTANCES", `[{"name":"first","base_url":"http://127.0.0.1:8317","management_key":"key-a"}]`)

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv returned error: %v", err)
	}
	// 显式单实例字段优先于 CPA_INSTANCES[0]，避免升级时静默改写已有对外地址。
	if cfg.CPABaseURL != "http://legacy.example.com:8317" || cfg.CPAManagementKey != "legacy-key" {
		t.Fatalf("expected explicit legacy fields to win, got base=%q key=%q", cfg.CPABaseURL, cfg.CPAManagementKey)
	}
	if cfg.MultiInstance() {
		t.Fatalf("single configured instance must not report multi instance")
	}
}

func TestLoadWithoutCPAInstancesKeepsSingleInstanceCompatibility(t *testing.T) {
	withIsolatedEnvFiles(t)
	t.Setenv("CPA_BASE_URL", "http://127.0.0.1:8317")
	t.Setenv("CPA_MANAGEMENT_KEY", "secret")
	t.Setenv("REDIS_QUEUE_ADDR", "127.0.0.1:9999")

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv returned error: %v", err)
	}
	if len(cfg.Instances) != 1 {
		t.Fatalf("expected synthesized single instance, got %+v", cfg.Instances)
	}
	// 单实例兼容模式实例名必须为空，前端标签与事件数据保持旧行为。
	if cfg.Instances[0].Name != "" {
		t.Fatalf("expected empty instance name for compatibility mode, got %q", cfg.Instances[0].Name)
	}
	if cfg.MultiInstance() {
		t.Fatalf("compatibility mode must not report multi instance")
	}
	if cfg.Instances[0].RedisQueueAddr != "127.0.0.1:9999" {
		t.Fatalf("expected synthesized instance to inherit redis queue addr, got %+v", cfg.Instances[0])
	}
}

func TestPrimaryInstanceFallsBackWhenInstancesUnset(t *testing.T) {
	cfg := Config{CPABaseURL: "http://127.0.0.1:8317", CPAManagementKey: "secret", RedisQueueAddr: "127.0.0.1:9999"}
	primary := cfg.PrimaryInstance()
	if primary.Name != "" || primary.BaseURL != cfg.CPABaseURL || primary.ManagementKey != cfg.CPAManagementKey || primary.RedisQueueAddr != cfg.RedisQueueAddr {
		t.Fatalf("unexpected fallback primary instance: %+v", primary)
	}
	if len(cfg.InstanceList()) != 1 {
		t.Fatalf("expected single synthesized instance, got %+v", cfg.InstanceList())
	}
}

func TestLoadRejectsInvalidCPAInstances(t *testing.T) {
	cases := []struct {
		name    string
		raw     string
		wantErr string
	}{
		{name: "invalid json", raw: `{"name":"a"}`, wantErr: "valid JSON array"},
		{name: "unknown field", raw: `[{"name":"a","base_url":"http://127.0.0.1:8317","management_key":"k","extra":1}]`, wantErr: "valid JSON array"},
		{name: "empty array", raw: `[]`, wantErr: "at least one instance"},
		{name: "missing name", raw: `[{"base_url":"http://127.0.0.1:8317","management_key":"k"}]`, wantErr: "name is required"},
		{name: "duplicate name", raw: `[{"name":"a","base_url":"http://127.0.0.1:8317","management_key":"k"},{"name":"a","base_url":"http://127.0.0.1:8318","management_key":"k2"}]`, wantErr: "duplicate name"},
		{name: "missing base url", raw: `[{"name":"a","management_key":"k"}]`, wantErr: "base_url is required"},
		{name: "missing management key", raw: `[{"name":"a","base_url":"http://127.0.0.1:8317"}]`, wantErr: "management_key is required"},
		{name: "name too long", raw: `[{"name":"` + strings.Repeat("x", CPAInstanceNameMaxLength+1) + `","base_url":"http://127.0.0.1:8317","management_key":"k"}]`, wantErr: "at most"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			withIsolatedEnvFiles(t)
			t.Setenv("CPA_INSTANCES", tc.raw)

			_, err := LoadFromEnv()
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("expected error containing %q, got %v", tc.wantErr, err)
			}
		})
	}
}
