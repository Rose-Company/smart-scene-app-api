package models

import (
	"smart-scene-app-api/common"
	"strings"
	"time"

	"github.com/google/uuid"
)

type QuerySort struct {
	Origin string
}

// Parse the query string to order string (Ex: http://example.com/messages?sort=created_at.asc,updated_at.acs
// => order string: created_at asc,updated_at acs)
func (s QuerySort) Parse() string {
	return strings.ReplaceAll(s.Origin, ".", " ")
}

type QueryParams struct {
	Limit  int
	Offset int
	QuerySort
	Preload  []common.Preload
	Selected []string
}

type BaseRequestParamsUri struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Sort     string `form:"sort"`
}

type BaseListResponse struct {
	Total    int         `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	Items    interface{} `json:"items"`
	Extra    interface{} `json:"extra"`
}

type Base struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}


func (b *BaseRequestParamsUri) VerifyPaging() {
	if b.Page <= 0 {
		b.Page = 1
	}
	if b.PageSize <= 0 {
		b.PageSize = 10
	}
}

