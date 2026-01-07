package entity

import (
	"time"

	"gorm.io/gorm"
)

type ProvisioningStatus string

const (
	ProvisioningStatusPending       ProvisioningStatus = "pending"
	ProvisioningStatusProvisioning  ProvisioningStatus = "provisioning"
	ProvisioningStatusProvisioned   ProvisioningStatus = "provisioned"
	ProvisioningStatusFailed        ProvisioningStatus = "failed"
	ProvisioningStatusDeProvisioned ProvisioningStatus = "de-provisioned"
)

type MotherService struct {
	gorm.Model
	ID                       uint64             `gorm:"primaryKey;autoIncrement"`
	Name                     string             `gorm:"column:name"`
	ExceptionRate            float64            `gorm:"column:exception_rate"`
	ResponseDelayRate        float64            `gorm:"column:response_delay_rate"`
	ResponseDelayDuration    *int               `gorm:"column:response_delay_duration"`
	RandomResponseDelayMin   *int               `gorm:"column:random_response_delay_min"`
	RandomResponseDelayMax   *int               `gorm:"column:random_response_delay_max"`
	ProvisioningStatus       ProvisioningStatus `gorm:"column:provisioning_status"`
	ServiceDeploymentAddress *string            `gorm:"column:service_deployment_address"`
	DatabaseName             string             `gorm:"column:database_name"`
	DatabaseTableName        string             `gorm:"column:database_table_name"`
	StoppedAt                *time.Time         `gorm:"column:stopped_at"`
	RestartedAt              *time.Time         `gorm:"column:restarted_at"`
	StartedAt                *time.Time         `gorm:"column:started_at"`
}

func (MotherService) TableName() string {
	return "mother_services"
}
