package entity

import "time"

type OutboxStatus string

const (
	OutboxStatusPending    OutboxStatus = "pending"    // waiting to be processed
	OutboxStatusProcessing OutboxStatus = "processing" // claimed by a worker, in progress
	OutboxStatusCompleted  OutboxStatus = "completed"  // processed successfully
	OutboxStatusFailed     OutboxStatus = "failed"     // permanently failed after exhausting max_attempts
)

type OutboxAggregateType string

const (
	OutboxAggregateTypeMotherService OutboxAggregateType = "mother_service"
)

type OutboxOperationType string

const (
	OutboxOperationProvisionMotherService OutboxOperationType = "provision_mother_service"
)

// Outbox is a durable, generic unit of work recorded in the same local
// transaction as the business row it describes. A background worker later
// claims pending rows and dispatches them to the handler for OperationType.
type Outbox struct {
	ID            uint64              `gorm:"primaryKey;autoIncrement;column:id"`
	AggregateType OutboxAggregateType `gorm:"column:aggregate_type"`
	AggregateID   uint64              `gorm:"column:aggregate_id"`
	OperationType OutboxOperationType `gorm:"column:operation_type"`
	Payload       *string             `gorm:"column:payload"`
	Status        OutboxStatus        `gorm:"column:status"`
	Attempts      int                 `gorm:"column:attempts"`
	MaxAttempts   int                 `gorm:"column:max_attempts"`
	LastError     *string             `gorm:"column:last_error"`
	AvailableAt   time.Time           `gorm:"column:available_at"`
	CreatedAt     time.Time           `gorm:"column:created_at"`
	UpdatedAt     time.Time           `gorm:"column:updated_at"`
}

func (Outbox) TableName() string {
	return "outbox"
}
