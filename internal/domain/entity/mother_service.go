package entity

import (
	"control-panel-service/pkg"
	"time"

	"gorm.io/gorm"
)

type MotherServiceStatus string

const (
	MotherServiceStatusPending MotherServiceStatus = "pending" // mother service just created
	MotherServiceStatusRunning MotherServiceStatus = "running" // test is running on application level (sending level)
	MotherServiceStatusPaused  MotherServiceStatus = "paused"  // application level pause on sending request
	MotherServiceStatusStopped MotherServiceStatus = "stopped" // stop container but can start scenario again.
	MotherServiceStatusAborted MotherServiceStatus = "aborted" //stop and delete containers. can not start again.
)

type MotherService struct {
	ID                       uint64              `gorm:"primaryKey;autoIncrement;column:id"`
	CreatedAt                time.Time           `gorm:"column:created_at"`
	UpdatedAt                time.Time           `gorm:"column:updated_at"`
	DeletedAt                *gorm.DeletedAt     `gorm:"column:deleted_at"`
	Name                     string              `gorm:"column:name"`
	ExceptionRate            int                 `gorm:"column:exception_rate"`
	ResponseDelayRate        int                 `gorm:"column:response_delay_rate"`
	ResponseDelayDuration    *int                `gorm:"column:response_delay_duration"`
	RandomResponseDelayMin   *int                `gorm:"column:random_response_delay_min"`
	RandomResponseDelayMax   *int                `gorm:"column:random_response_delay_max"`
	Status                   MotherServiceStatus `gorm:"column:status"`
	ServiceDeploymentAddress *string             `gorm:"column:service_deployment_address"`
	DatabaseName             string              `gorm:"column:database_name"`
	DatabaseTableName        string              `gorm:"column:database_table_name"`
}

// If both values of page and per page be zero then all items will be returned.
type PaginationRequest struct {
	Page    int
	PerPage int
}

func (MotherService) TableName() string {
	return "mother_services"
}

// nolint
func (m *MotherService) Validate() error {
	if m.Name == "" {
		return pkg.ErrInvalidName
	}

	if m.ExceptionRate < 0 || m.ExceptionRate > 100 {
		return pkg.ErrInvalidExceptionRate
	}

	if m.ResponseDelayRate < 0 || m.ResponseDelayRate > 100 {
		return pkg.ErrInvalidResponseDelayRate
	}

	isNoDelay := m.ResponseDelayRate == 0 &&
		m.ResponseDelayDuration == nil &&
		m.RandomResponseDelayMin == nil &&
		m.RandomResponseDelayMax == nil

	isFixedDelay := m.ResponseDelayRate > 0 &&
		m.ResponseDelayDuration != nil && *m.ResponseDelayDuration > 0 &&
		m.RandomResponseDelayMin == nil &&
		m.RandomResponseDelayMax == nil

	isRandomDelay := m.ResponseDelayRate > 0 &&
		m.ResponseDelayDuration == nil &&
		m.RandomResponseDelayMin != nil && *m.RandomResponseDelayMin >= 0 &&
		m.RandomResponseDelayMax != nil && *m.RandomResponseDelayMax > 0

	if !(isNoDelay || isFixedDelay || isRandomDelay) {
		return pkg.ErrInvalidDelayConfiguration
	}

	if m.RandomResponseDelayMin != nil || m.RandomResponseDelayMax != nil {
		if m.RandomResponseDelayMin == nil || m.RandomResponseDelayMax == nil {
			return pkg.ErrInvalidDelayConfiguration
		}
		if *m.RandomResponseDelayMin >= *m.RandomResponseDelayMax {
			return pkg.ErrInvalidRandomDelayRange
		}
	}

	return nil
}
