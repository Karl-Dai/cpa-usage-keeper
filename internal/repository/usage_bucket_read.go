package repository

import (
	"time"

	"cpa-usage-keeper/internal/entities"
	"cpa-usage-keeper/internal/overview"
	"cpa-usage-keeper/internal/timeutil"

	"gorm.io/gorm"
)

// usageBucketRead 只解释既有 Overview 存储键。真实日界不能直接替代旧日桶键，
// 因此查询用 keys，返回的日序列用 displayStarts；两者都不写回数据库。
type usageBucketRead struct {
	keys          []string
	displayStarts map[int64]time.Time
}

func overviewBucketRead(start, end time.Time, grain string) usageBucketRead {
	start, end = timeutil.NormalizeStorageTime(start), timeutil.NormalizeStorageTime(end)
	result := usageBucketRead{displayStarts: make(map[int64]time.Time)}
	if grain == "daily" {
		for date := timeutil.DateOf(start); ; date = date.AddDays(1) {
			bounds, exists := timeutil.DayBounds(date, start.Location())
			if !exists {
				continue
			}
			if !bounds.Start.Before(end) {
				break
			}
			if bounds.Start.Before(start) {
				continue
			}
			// 调用现有纯归桶函数，保持历史 midnight gap/fold 的存储规则。
			_, key := overview.BucketKeysForEvent(entities.UsageEvent{Timestamp: bounds.Start})
			result.keys = append(result.keys, timeutil.FormatStorageTime(key.BucketStart))
			result.displayStarts[key.BucketStart.Unix()] = bounds.Start
		}
		return result
	}
	// Overview 按绝对整小时分桶，不能替换成 Latency 的当地钟面整点。
	for at := start.Truncate(time.Hour); at.Before(end); at = at.Add(time.Hour) {
		if !at.Before(start) {
			result.keys = append(result.keys, timeutil.FormatStorageTime(at))
		}
	}
	return result
}

func (b usageBucketRead) query(query *gorm.DB) *gorm.DB {
	// 页面范围有界（小时至多跨 DST 日、日 <=365）；IN 使用既有 bucket_start 索引，
	// 不对时间列套函数，也不依赖含 offset 字符串的大小顺序。
	return query.Where("bucket_start IN ?", b.keys)
}

func (b usageBucketRead) displayStart(stored time.Time) time.Time {
	if start, ok := b.displayStarts[stored.Unix()]; ok {
		return start
	}
	return timeutil.NormalizeStorageTime(stored)
}
