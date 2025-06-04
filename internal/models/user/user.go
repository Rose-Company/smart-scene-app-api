package user

import (
	"smart-scene-app-api/common"
	models "smart-scene-app-api/internal/models"

	"github.com/google/uuid"
)

// User represents the user model
type User struct {
    models.Base
    Email    string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
    Password string    `gorm:"type:varchar(255);not null" json:"password"`
    FullName string    `gorm:"type:varchar(255);not null" json:"full_name"`
    RoleID   uuid.UUID `gorm:"type:uuid;not null" json:"role"`
    Status   string    `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
}

func (User) TableName() string {
	return common.POSTGRES_TABLE_NAME_USERS
}
