package entity

import (
	"time"

	"gorm.io/gorm"
)

type TestServiceStatus string

const (
	TestServiceStatusPending TestServiceStatus = "pending"
	TestServiceStatusRunning TestServiceStatus = "running"
	TestServiceStatusPaused  TestServiceStatus = "paused"
	TestServiceStatusStopped TestServiceStatus = "stopped"
	TestServiceStatusDeleted TestServiceStatus = "deleted"
	TestServiceStatusSucceed TestServiceStatus = "succeed"
)

type TestService struct {
	ID        uint64            `gorm:"primaryKey;autoIncrement;column:id"`
	CreatedAt time.Time         `gorm:"column:created_at"`
	UpdatedAt time.Time         `gorm:"column:updated_at"`
	DeletedAt *gorm.DeletedAt   `gorm:"column:deleted_at"`
	Status    TestServiceStatus `gorm:"column:status"`
}

func (TestService) TableName() string {
	return "test_services"
}
