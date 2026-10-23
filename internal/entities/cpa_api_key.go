package entities

import "time"

// CPAAPIKey 保存 CPA 管理接口同步到本地的 API-Key，完整 key 仅供后端内部查询使用。
type CPAAPIKey struct {
	ID int64 `gorm:"primaryKey"`
	// Instance 标识该 key 所属的 CPA 实例；单实例部署为空字符串。
	Instance     string `gorm:"not null;default:'';uniqueIndex:uniq_cpa_api_keys_instance_api_key,priority:1"`
	APIKey       string `gorm:"uniqueIndex:uniq_cpa_api_keys_instance_api_key,priority:2"`
	DisplayKey   string
	KeyAlias     string
	IsDeleted    bool       `gorm:"index:idx_cpa_api_keys_is_deleted"`
	LastSyncedAt *time.Time `gorm:"serializer:storageTime"`
	CreatedAt    time.Time  `gorm:"serializer:storageTime"`
	UpdatedAt    time.Time  `gorm:"serializer:storageTime"`
}
