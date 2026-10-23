package dto

import "time"

// RedisInboxInsert 是 Redis usage inbox 的入库参数，source 原样落库。
type RedisInboxInsert struct {
	Source     string
	RawMessage string
	PoppedAt   time.Time
	// Instance 标识消息来源实例；空字符串表示单实例兼容模式。
	Instance string
}

// RedisUsageInboxCleanupResult 是 Redis usage inbox 的清理结果。
type RedisUsageInboxCleanupResult struct {
	ProcessedDeleted int64
	FailedDeleted    int64
}
