package migration

import (
	"fmt"

	"cpa-usage-keeper/internal/entities"
	"gorm.io/gorm"
)

// addCPAInstanceLabelsMigration 为多实例汇总补充实例标签列，并把单 key 唯一索引改为实例内唯一。
// 已存在的单实例数据保留空实例值，由命名主实例首次同步时认领，不在这里猜测实例名。
func addCPAInstanceLabelsMigration(tx *gorm.DB) error {
	// usage_events 是汇总视图的数据源，必须能区分事件来自哪个实例。
	if err := addInstanceColumnIfMissing(tx, &entities.UsageEvent{}, "usage_events"); err != nil {
		return err
	}
	// archive 与热表共享复制列清单，必须同步补齐，否则归档 INSERT SELECT 会缺列。
	if err := addInstanceColumnIfMissing(tx, &entities.UsageEventArchive{}, "usage_events_archive"); err != nil {
		return err
	}
	// inbox 行在解码前就携带实例标签，保证事件落库时标签已经确定。
	if err := addInstanceColumnIfMissing(tx, &entities.RedisUsageInbox{}, "redis_usage_inboxes"); err != nil {
		return err
	}
	if !tx.Migrator().HasTable("cpa_api_keys") {
		return nil
	}
	// cpa_api_keys 按实例分别保存管理 key，允许不同实例出现相同 key 值。
	if err := addInstanceColumnIfMissing(tx, &entities.CPAAPIKey{}, "cpa_api_keys"); err != nil {
		return err
	}
	// 旧唯一索引只约束 api_key，会阻止多实例保存同一个 key，必须先删除。
	if err := tx.Exec("DROP INDEX IF EXISTS uniq_cpa_api_keys_api_key").Error; err != nil {
		return fmt.Errorf("drop cpa_api_keys api_key unique index: %w", err)
	}
	// 新索引把实例名纳入唯一键，单实例部署等价于原来的 api_key 唯一。
	if err := tx.Exec("CREATE UNIQUE INDEX IF NOT EXISTS uniq_cpa_api_keys_instance_api_key ON cpa_api_keys(instance, api_key)").Error; err != nil {
		return fmt.Errorf("create cpa_api_keys instance/api_key unique index: %w", err)
	}
	return nil
}

// addInstanceColumnIfMissing 幂等补充 instance 列，历史行统一落到空实例。
func addInstanceColumnIfMissing(tx *gorm.DB, model any, table string) error {
	if !tx.Migrator().HasTable(table) {
		return nil
	}
	if tx.Migrator().HasColumn(table, "instance") {
		return nil
	}
	// 非空列必须带默认值，否则已有行会在 SQLite 增列时违反约束。
	if err := tx.Migrator().AddColumn(model, "Instance"); err != nil {
		return fmt.Errorf("add %s.instance column: %w", table, err)
	}
	return nil
}
