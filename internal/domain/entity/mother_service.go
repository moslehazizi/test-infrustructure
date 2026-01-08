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
	ID                       uint64             `gorm:"primaryKey;autoIncrement;column:id"`
	CreatedAt                time.Time          `gorm:"column:created_at"`
	UpdatedAt                time.Time          `gorm:"column:updated_at"`
	DeletedAt                *gorm.DeletedAt    `gorm:"index;column:deleted_at"`
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
	KafkaLiveFeedTopic       string             `gorm:"column:kafka_livefeed_topic"`
	KafkaFactorialTopic      string             `gorm:"column:kafka_factorial_topic"`
	StoppedAt                *time.Time         `gorm:"column:stopped_at"`
	RestartedAt              *time.Time         `gorm:"column:restarted_at"`
	StartedAt                *time.Time         `gorm:"column:started_at"`
}

// if both values of page and per page be zero then all items will be returned
type PaginationRequest struct {
	Page    int
	PerPage int
}

func (MotherService) TableName() string {
	return "mother_services"
}
