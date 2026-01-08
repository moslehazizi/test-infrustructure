package entity

import "time"

type TestCategory struct {
	ID                      uint `gorm:"primarykey"`
	CreatedAt               time.Time
	UpdatedAt               time.Time
	Name                    string
	Label                   string
	HasMaxTestServiceCount  bool
	HasExecutionDuration    bool
	HasAutoStepIncreaseRate bool
}
