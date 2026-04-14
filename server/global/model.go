package global

import (
	"time"

	"gorm.io/gorm"
)

type GVA_MODEL struct {
	ID        uint           `gorm:"primarykey" json:"ID"` // primary key ID
	CreatedAt time.Time      // created at
	UpdatedAt time.Time      // updated at
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // deleted at
}
