package test

import (
	"context"
	"testing"
	"time"

	"cpa-usage-keeper/internal/cpa/dto/cpaapikeys"
	"cpa-usage-keeper/internal/cpa/dto/response"
	"cpa-usage-keeper/internal/entities"
	"cpa-usage-keeper/internal/repository"
	"cpa-usage-keeper/internal/repository/dto"
	. "cpa-usage-keeper/internal/service"
)

func TestSyncManagementAPIKeysOnlyTagsInstanceAndSkipsOtherMetadata(t *testing.T) {
	db := openMetadataTestDatabase(t, "instance-management-api-keys.db")
	now := time.Date(2026, 10, 23, 10, 0, 0, 0, time.UTC)
	fetcher := newMetadataTestFetcher()
	fetcher.managementAPIKeysResult = &response.ManagementAPIKeysResult{StatusCode: 200, Payload: cpaapikeys.ManagementAPIKeysResponse{APIKeys: []string{"sk-xuxiaozhi"}}}
	syncer := NewSyncServiceWithOptions(db, SyncServiceOptions{
		BaseURL:         "https://cpa.example.com",
		MetadataFetcher: fetcher,
		Now:             func() time.Time { return now },
		InstanceName:    "xuxiaozhi",
	})

	if err := syncer.SyncManagementAPIKeysOnly(context.Background()); err != nil {
		t.Fatalf("SyncManagementAPIKeysOnly returned error: %v", err)
	}
	// 附属实例只能读取管理 API key，不得触碰主实例的 credentials／provider 身份数据。
	if fetcher.callCount("credentials") != 0 {
		t.Fatalf("expected no credential fetch, got %d", fetcher.callCount("credentials"))
	}
	if fetcher.callCount("management-api-keys") != 1 {
		t.Fatalf("expected one management api key fetch, got %d", fetcher.callCount("management-api-keys"))
	}

	var row entities.CPAAPIKey
	if err := db.Where("api_key = ?", "sk-xuxiaozhi").First(&row).Error; err != nil {
		t.Fatalf("expected synced row: %v", err)
	}
	if row.Instance != "xuxiaozhi" || row.IsDeleted {
		t.Fatalf("unexpected instance API key row: %+v", row)
	}
}

func TestProcessRedisUsageInboxPropagatesInstanceToUsageEvents(t *testing.T) {
	db := openSyncTestDatabase(t)
	poppedAt := time.Date(2026, 10, 23, 10, 0, 0, 0, time.UTC)
	if _, err := repository.InsertRedisUsageInboxMessages(db, []dto.RedisInboxInsert{{
		Source:     redisUsageInboxTestSource,
		Instance:   "zhaimenghui",
		PoppedAt:   poppedAt,
		RawMessage: `{"timestamp":"2026-10-23T10:00:00Z","provider":"claude","model":"sonnet","execution_id":"fixture-execution","request_id":"instance-event","tokens":{"input_tokens":1,"output_tokens":2}}`,
	}}); err != nil {
		t.Fatalf("seed instance inbox message: %v", err)
	}
	syncer := NewSyncServiceWithOptions(db, SyncServiceOptions{PricingCatalog: emptyPricingCatalogForTest(), BaseURL: "https://cpa.example.com"})

	result, err := syncer.ProcessRedisUsageInbox(context.Background())
	if err != nil {
		t.Fatalf("ProcessRedisUsageInbox returned error: %v", err)
	}
	if result == nil || result.InsertedEvents != 1 {
		t.Fatalf("unexpected process result: %+v", result)
	}

	var event entities.UsageEvent
	if err := db.Where("event_key = ?", "instance-event").First(&event).Error; err != nil {
		t.Fatalf("load usage event: %v", err)
	}
	if event.Instance != "zhaimenghui" {
		t.Fatalf("expected usage event to carry inbox instance label, got %+v", event)
	}
}
