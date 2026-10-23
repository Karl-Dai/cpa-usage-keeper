package repository

import (
	"strings"
	"time"

	"cpa-usage-keeper/internal/entities"
	"cpa-usage-keeper/internal/helper"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// SyncCPAAPIKeys 保留单实例签名，按空实例同步；多实例调用方使用 SyncCPAAPIKeysForInstance。
func SyncCPAAPIKeys(db *gorm.DB, keys []string, syncedAt time.Time) error {
	return SyncCPAAPIKeysForInstance(db, "", false, keys, syncedAt)
}

// SyncCPAAPIKeysForInstance 只同步指定实例的管理 API key。
// adoptUnattributed 仅对命名主实例开启：把历史单实例行（instance 为空）认领到主实例名下，
// 避免升级后同一个 key 因实例标签不同而重复插入。
func SyncCPAAPIKeysForInstance(db *gorm.DB, instance string, adoptUnattributed bool, keys []string, syncedAt time.Time) error {
	instance = strings.TrimSpace(instance)
	seen := make(map[string]struct{}, len(keys))
	uniqueKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		uniqueKeys = append(uniqueKeys, key)
	}

	var existingRows []struct {
		ID         int64
		Instance   string
		APIKey     string
		DisplayKey string
		IsDeleted  bool
	}
	// 只读取当前实例的行；命名主实例额外认领历史空实例行，保持单实例升级连续性。
	existingQuery := db.Clauses(dbresolver.Read).Model(&entities.CPAAPIKey{}).
		Select("id, instance, api_key, display_key, is_deleted")
	switch {
	case instance == "":
		existingQuery = existingQuery.Where("instance = ?", "")
	case adoptUnattributed:
		existingQuery = existingQuery.Where("instance IN ?", []string{instance, ""})
	default:
		existingQuery = existingQuery.Where("instance = ?", instance)
	}
	// 先在独立 reader 上比较完整 key 集合，只有变化才进入 writer 事务。
	if err := existingQuery.Find(&existingRows).Error; err != nil {
		return err
	}
	type keyUpdate struct {
		id     int64
		fields map[string]any
	}
	existingByKey := make(map[string]int, len(existingRows))
	for index, row := range existingRows {
		existingByKey[row.APIKey] = index
	}
	incoming := make(map[string]struct{}, len(uniqueKeys))
	toCreate := make([]entities.CPAAPIKey, 0)
	toUpdate := make([]keyUpdate, 0)
	for _, key := range uniqueKeys {
		incoming[key] = struct{}{}
		if index, ok := existingByKey[key]; ok {
			row := existingRows[index]
			fields := make(map[string]any)
			if display := helper.RedactSensitiveValue(key); row.DisplayKey != display {
				fields["display_key"] = display
			}
			if row.IsDeleted {
				fields["is_deleted"] = false
			}
			// 认领历史空实例行时同步写入实例名，后续轮次无需再次匹配空实例。
			if row.Instance != instance {
				fields["instance"] = instance
			}
			if len(fields) > 0 {
				fields["last_synced_at"] = &syncedAt
				fields["updated_at"] = syncedAt
				toUpdate = append(toUpdate, keyUpdate{id: row.ID, fields: fields})
			}
			continue
		}
		toCreate = append(toCreate, entities.CPAAPIKey{
			Instance:     instance,
			APIKey:       key,
			DisplayKey:   helper.RedactSensitiveValue(key),
			IsDeleted:    false,
			LastSyncedAt: &syncedAt,
		})
	}
	staleIDs := make([]int64, 0)
	for _, row := range existingRows {
		if row.IsDeleted {
			continue
		}
		if _, ok := incoming[row.APIKey]; !ok {
			staleIDs = append(staleIDs, row.ID)
		}
	}
	if len(toCreate) == 0 && len(toUpdate) == 0 && len(staleIDs) == 0 {
		return nil
	}
	return db.Transaction(func(tx *gorm.DB) error {
		for _, change := range toUpdate {
			if err := tx.Model(&entities.CPAAPIKey{}).Where("id = ?", change.id).Updates(change.fields).Error; err != nil {
				return err
			}
		}
		if len(toCreate) > 0 {
			if err := tx.Create(&toCreate).Error; err != nil {
				return err
			}
		}
		if len(staleIDs) == 0 {
			return nil
		}
		return tx.Model(&entities.CPAAPIKey{}).Where("id IN ?", staleIDs).Updates(map[string]any{"is_deleted": true, "updated_at": syncedAt}).Error
	})
}

func ListActiveCPAAPIKeys(db *gorm.DB) ([]entities.CPAAPIKey, error) {
	var rows []entities.CPAAPIKey
	err := db.Where("is_deleted = ?", false).Order("id asc").Find(&rows).Error
	return rows, err
}

func FindActiveCPAAPIKeyByID(db *gorm.DB, id int64) (entities.CPAAPIKey, error) {
	var row entities.CPAAPIKey
	err := db.Where("id = ? AND is_deleted = ?", id, false).First(&row).Error
	return row, err
}

func FindActiveCPAAPIKeyByValue(db *gorm.DB, apiKey string) (entities.CPAAPIKey, error) {
	var row entities.CPAAPIKey
	err := db.Where("api_key = ? AND is_deleted = ?", apiKey, false).First(&row).Error
	return row, err
}

func UpdateCPAAPIKeyAlias(db *gorm.DB, id int64, keyAlias string) error {
	result := db.Model(&entities.CPAAPIKey{}).Where("id = ? AND is_deleted = ?", id, false).Update("key_alias", strings.TrimSpace(keyAlias))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
