package common

import (
	"gorm.io/gorm"
)

type Preload struct {
	Model    string
	Selected []string
	Conds    map[string]interface{}
	Order    string
	Join     string
	Limit    int
	Offset   int
}

func ApplyPreload(db *gorm.DB, preloadConfig Preload) *gorm.DB {
	return db.Preload(preloadConfig.Model, func(db *gorm.DB) *gorm.DB {
		if preloadConfig.Selected != nil {
			db = db.Select(preloadConfig.Selected)
		}

		if preloadConfig.Limit != 0 {
			db.Limit(preloadConfig.Limit)
			db.Offset(preloadConfig.Offset * preloadConfig.Limit)
		}

		for k, v := range preloadConfig.Conds {
			db = db.Where(k, v)
		}

		if preloadConfig.Order != "" {
			db = db.Order(preloadConfig.Order)
		}

		if preloadConfig.Join != "" {
			db.Joins(preloadConfig.Join)
		}

		return db
	})
}

type BaseModel struct {
	CreateAt int64 `json:"create_at,omitempty"`
	UpdateAt int64 `json:"update_at,omitempty"`
}
