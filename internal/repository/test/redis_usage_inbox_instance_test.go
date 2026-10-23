package test

import (
	"testing"
	"time"

	"cpa-usage-keeper/internal/repository"
)

func TestInsertRedisUsageInboxRawMessagesForInstanceTagsRows(t *testing.T) {
	db := openTestDatabase(t)
	poppedAt := time.Date(2026, 10, 23, 9, 0, 0, 0, time.UTC)

	if _, err := repository.InsertRedisUsageInboxRawMessagesForInstance(db, "xuxiaozhi", "redis_pull:usage", []string{`{"request_id":"instance-a"}`}, poppedAt); err != nil {
		t.Fatalf("insert instance inbox rows returned error: %v", err)
	}
	if _, err := repository.InsertRedisUsageInboxRawMessages(db, "redis_pull:usage", []string{`{"request_id":"legacy"}`}, poppedAt); err != nil {
		t.Fatalf("insert legacy inbox rows returned error: %v", err)
	}

	rows, err := repository.ListProcessableRedisUsageInbox(db, 10)
	if err != nil {
		t.Fatalf("ListProcessableRedisUsageInbox returned error: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected two processable rows, got %+v", rows)
	}
	// 处理列必须包含 instance，否则解码阶段拿不到实例标签。
	if rows[0].Instance != "xuxiaozhi" {
		t.Fatalf("expected instance tag on first row, got %+v", rows[0])
	}
	if rows[1].Instance != "" {
		t.Fatalf("expected legacy row to stay unattributed, got %+v", rows[1])
	}
}
