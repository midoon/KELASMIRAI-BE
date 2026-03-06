package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type WebhookLog struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Provider       string         `gorm:"type:text;not null;uniqueIndex:unique_provider_external"`
	ExternalID     *string        `gorm:"type:text;uniqueIndex:unique_provider_external"`
	PayloadJSON    datatypes.JSON `gorm:"type:jsonb;not null"`
	SignatureValid bool           `gorm:"not null;default:false"`
	Processed      bool           `gorm:"not null;default:false"`
	CreatedAt      time.Time      `gorm:"type:timestamptz;not null;default:now()"`
}

func (WebhookLog) TableName() string {
	return "webhook_logs"
}

func (m *WebhookLog) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}
