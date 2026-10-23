package test

import (
	"testing"
	"time"

	"cpa-usage-keeper/internal/entities"
	"cpa-usage-keeper/internal/repository"
)

func TestSyncCPAAPIKeysForInstanceKeepsSameKeyPerInstance(t *testing.T) {
	db := openTestDatabase(t)
	firstSync := time.Date(2026, 5, 13, 10, 0, 0, 0, time.UTC)
	secondSync := firstSync.Add(time.Hour)

	if err := repository.SyncCPAAPIKeysForInstance(db, "daichangyu", false, []string{"sk-shared"}, firstSync); err != nil {
		t.Fatalf("sync primary instance returned error: %v", err)
	}
	// 不同实例允许出现相同 key 值，唯一索引必须按实例生效。
	if err := repository.SyncCPAAPIKeysForInstance(db, "xuxiaozhi", false, []string{"sk-shared", "sk-xuxiaozhi"}, firstSync); err != nil {
		t.Fatalf("sync secondary instance returned error: %v", err)
	}

	var rows []entities.CPAAPIKey
	if err := db.Order("instance asc, api_key asc").Find(&rows).Error; err != nil {
		t.Fatalf("load synced rows: %v", err)
	}
	if len(rows) != 3 {
		t.Fatalf("expected three rows across two instances, got %+v", rows)
	}
	if rows[0].Instance != "daichangyu" || rows[1].Instance != "xuxiaozhi" || rows[2].Instance != "xuxiaozhi" {
		t.Fatalf("unexpected instance tags: %+v", rows)
	}

	// 附属实例不再返回 sk-shared 时只标记自己那一行删除，主实例行必须保持 active。
	if err := repository.SyncCPAAPIKeysForInstance(db, "xuxiaozhi", false, []string{"sk-xuxiaozhi"}, secondSync); err != nil {
		t.Fatalf("second secondary sync returned error: %v", err)
	}
	var primaryShared entities.CPAAPIKey
	if err := db.Where("instance = ? AND api_key = ?", "daichangyu", "sk-shared").First(&primaryShared).Error; err != nil {
		t.Fatalf("expected primary shared row: %v", err)
	}
	if primaryShared.IsDeleted {
		t.Fatalf("stale marking must stay scoped to the synced instance, got %+v", primaryShared)
	}
	var secondaryShared entities.CPAAPIKey
	if err := db.Where("instance = ? AND api_key = ?", "xuxiaozhi", "sk-shared").First(&secondaryShared).Error; err != nil {
		t.Fatalf("expected secondary shared row: %v", err)
	}
	if !secondaryShared.IsDeleted {
		t.Fatalf("expected missing secondary key to be marked deleted, got %+v", secondaryShared)
	}
}

func TestSyncCPAAPIKeysForInstanceAdoptsUnattributedLegacyRows(t *testing.T) {
	db := openTestDatabase(t)
	syncedAt := time.Date(2026, 5, 13, 10, 0, 0, 0, time.UTC)

	// 历史单实例数据落在空实例上，升级后命名主实例必须认领而不是重复插入。
	if err := repository.SyncCPAAPIKeys(db, []string{"sk-legacy"}, syncedAt); err != nil {
		t.Fatalf("legacy sync returned error: %v", err)
	}
	if err := repository.UpdateCPAAPIKeyAlias(db, 1, "Legacy Alias"); err != nil {
		t.Fatalf("seed legacy alias: %v", err)
	}
	if err := repository.SyncCPAAPIKeysForInstance(db, "daichangyu", true, []string{"sk-legacy"}, syncedAt.Add(time.Hour)); err != nil {
		t.Fatalf("adopting sync returned error: %v", err)
	}

	var rows []entities.CPAAPIKey
	if err := db.Find(&rows).Error; err != nil {
		t.Fatalf("load rows: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected legacy row to be adopted instead of duplicated, got %+v", rows)
	}
	if rows[0].Instance != "daichangyu" || rows[0].ID != 1 || rows[0].KeyAlias != "Legacy Alias" || rows[0].IsDeleted {
		t.Fatalf("unexpected adopted row: %+v", rows[0])
	}
}

func TestSyncCPAAPIKeysForInstanceWithoutAdoptionIgnoresLegacyRows(t *testing.T) {
	db := openTestDatabase(t)
	syncedAt := time.Date(2026, 5, 13, 10, 0, 0, 0, time.UTC)

	if err := repository.SyncCPAAPIKeys(db, []string{"sk-legacy"}, syncedAt); err != nil {
		t.Fatalf("legacy sync returned error: %v", err)
	}
	// 附属实例不认领历史行，必须新建自己的行且不动空实例行，避免误改主实例归属。
	if err := repository.SyncCPAAPIKeysForInstance(db, "xuxiaozhi", false, []string{"sk-legacy"}, syncedAt); err != nil {
		t.Fatalf("secondary sync returned error: %v", err)
	}

	var rows []entities.CPAAPIKey
	if err := db.Order("id asc").Find(&rows).Error; err != nil {
		t.Fatalf("load rows: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected new secondary row plus untouched legacy row, got %+v", rows)
	}
	if rows[0].Instance != "" || rows[0].IsDeleted {
		t.Fatalf("expected legacy row to stay unattributed and active, got %+v", rows[0])
	}
	if rows[1].Instance != "xuxiaozhi" {
		t.Fatalf("expected secondary row to carry instance tag, got %+v", rows[1])
	}
}

func TestListActiveCPAAPIKeysReturnsEveryInstance(t *testing.T) {
	db := openTestDatabase(t)
	syncedAt := time.Date(2026, 5, 13, 10, 0, 0, 0, time.UTC)

	if err := repository.SyncCPAAPIKeysForInstance(db, "daichangyu", false, []string{"sk-a"}, syncedAt); err != nil {
		t.Fatalf("sync instance a returned error: %v", err)
	}
	if err := repository.SyncCPAAPIKeysForInstance(db, "yangbing", false, []string{"sk-b"}, syncedAt); err != nil {
		t.Fatalf("sync instance b returned error: %v", err)
	}

	rows, err := repository.ListActiveCPAAPIKeys(db)
	if err != nil {
		t.Fatalf("ListActiveCPAAPIKeys returned error: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected keys from every instance, got %+v", rows)
	}
	if rows[0].Instance != "daichangyu" || rows[1].Instance != "yangbing" {
		t.Fatalf("expected instance tags to be preserved, got %+v", rows)
	}
}
