package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository/mocks"
	"control-panel-service/pkg"
	"errors"
	"testing"
	"time"

	provisionProvider "control-panel-service/internal/provider/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOutboxProcessor_ProcessPending(t *testing.T) {
	t.Run("failed_case_claim_pending_error", func(t *testing.T) {
		ctx := context.Background()
		mockOutbox := new(mocks.MockOutbox)
		mockRepo := new(mocks.MockMotherService)
		mockProvision := new(provisionProvider.MockProvisioningService)

		processor := NewOutboxProcessor(mockOutbox, mockRepo, mockProvision, 10, time.Second)

		mockOutbox.On("ClaimPending", mock.Anything, 10).Return(nil, errors.New("connection lost"))

		err := processor.ProcessPending(ctx)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToClaimOutboxItems)
		mockOutbox.AssertExpectations(t)
	})

	t.Run("success_case_no_pending_items", func(t *testing.T) {
		ctx := context.Background()
		mockOutbox := new(mocks.MockOutbox)
		mockRepo := new(mocks.MockMotherService)
		mockProvision := new(provisionProvider.MockProvisioningService)

		processor := NewOutboxProcessor(mockOutbox, mockRepo, mockProvision, 10, time.Second)

		mockOutbox.On("ClaimPending", mock.Anything, 10).Return([]*entity.Outbox{}, nil)

		err := processor.ProcessPending(ctx)

		assert.NoError(t, err)
		mockOutbox.AssertExpectations(t)
	})

	t.Run("success_case_provision_mother_service", func(t *testing.T) {
		ctx := context.Background()
		mockOutbox := new(mocks.MockOutbox)
		mockRepo := new(mocks.MockMotherService)
		mockProvision := new(provisionProvider.MockProvisioningService)

		processor := NewOutboxProcessor(mockOutbox, mockRepo, mockProvision, 10, time.Second)

		motherService := &entity.MotherService{ID: uint64(1), Name: "mother1"}
		item := &entity.Outbox{
			ID:            uint64(100),
			AggregateType: entity.OutboxAggregateTypeMotherService,
			AggregateID:   motherService.ID,
			OperationType: entity.OutboxOperationProvisionMotherService,
			Attempts:      1,
			MaxAttempts:   5,
		}

		mockOutbox.On("ClaimPending", mock.Anything, 10).Return([]*entity.Outbox{item}, nil)
		mockRepo.On("GetByID", mock.Anything, motherService.ID).Return(motherService, nil)
		mockProvision.On("ProvisionMotherService", mock.Anything, motherService).Return(nil)
		mockRepo.On("SetStatus", mock.Anything, motherService.ID, entity.MotherServiceStatusRunning).Return(nil)
		mockOutbox.On("MarkCompleted", mock.Anything, item.ID).Return(nil)

		err := processor.ProcessPending(ctx)

		assert.NoError(t, err)
		mockOutbox.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockProvision.AssertExpectations(t)
	})

	t.Run("success_case_provision_succeeds_but_set_status_and_mark_completed_best_effort", func(t *testing.T) {
		ctx := context.Background()
		mockOutbox := new(mocks.MockOutbox)
		mockRepo := new(mocks.MockMotherService)
		mockProvision := new(provisionProvider.MockProvisioningService)

		processor := NewOutboxProcessor(mockOutbox, mockRepo, mockProvision, 10, time.Second)

		motherService := &entity.MotherService{ID: uint64(1), Name: "mother1"}
		item := &entity.Outbox{
			ID:            uint64(100),
			AggregateType: entity.OutboxAggregateTypeMotherService,
			AggregateID:   motherService.ID,
			OperationType: entity.OutboxOperationProvisionMotherService,
			Attempts:      1,
			MaxAttempts:   5,
		}

		mockOutbox.On("ClaimPending", mock.Anything, 10).Return([]*entity.Outbox{item}, nil)
		mockRepo.On("GetByID", mock.Anything, motherService.ID).Return(motherService, nil)
		mockProvision.On("ProvisionMotherService", mock.Anything, motherService).Return(nil)
		mockRepo.On("SetStatus", mock.Anything, motherService.ID, entity.MotherServiceStatusRunning).Return(errors.New("db unavailable"))
		mockOutbox.On("MarkCompleted", mock.Anything, item.ID).Return(nil)

		err := processor.ProcessPending(ctx)

		assert.NoError(t, err)
		mockOutbox.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_unknown_operation_type_marks_permanently_failed", func(t *testing.T) {
		ctx := context.Background()
		mockOutbox := new(mocks.MockOutbox)
		mockRepo := new(mocks.MockMotherService)
		mockProvision := new(provisionProvider.MockProvisioningService)

		processor := NewOutboxProcessor(mockOutbox, mockRepo, mockProvision, 10, time.Second)

		item := &entity.Outbox{
			ID:            uint64(101),
			AggregateType: entity.OutboxAggregateTypeMotherService,
			AggregateID:   uint64(1),
			OperationType: entity.OutboxOperationType("unknown_operation"),
			Attempts:      1,
			MaxAttempts:   5,
		}

		mockOutbox.On("ClaimPending", mock.Anything, 10).Return([]*entity.Outbox{item}, nil)
		mockOutbox.On("MarkFailedPermanently", mock.Anything, item.ID, "unknown operation type").Return(nil)

		err := processor.ProcessPending(ctx)

		assert.NoError(t, err)
		mockOutbox.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockProvision.AssertExpectations(t)
	})

	t.Run("failed_case_get_mother_service_reschedules_when_attempts_remain", func(t *testing.T) {
		ctx := context.Background()
		mockOutbox := new(mocks.MockOutbox)
		mockRepo := new(mocks.MockMotherService)
		mockProvision := new(provisionProvider.MockProvisioningService)

		processor := NewOutboxProcessor(mockOutbox, mockRepo, mockProvision, 10, time.Second)

		item := &entity.Outbox{
			ID:            uint64(102),
			AggregateType: entity.OutboxAggregateTypeMotherService,
			AggregateID:   uint64(1),
			OperationType: entity.OutboxOperationProvisionMotherService,
			Attempts:      1,
			MaxAttempts:   5,
		}

		mockOutbox.On("ClaimPending", mock.Anything, 10).Return([]*entity.Outbox{item}, nil)
		mockRepo.On("GetByID", mock.Anything, item.AggregateID).Return(nil, errors.New("not found"))
		mockOutbox.On("Reschedule", mock.Anything, item.ID, mock.AnythingOfType("time.Time"), "not found").Return(nil)

		err := processor.ProcessPending(ctx)

		assert.NoError(t, err)
		mockOutbox.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockProvision.AssertNotCalled(t, "DeprovisionMotherService", mock.Anything, mock.Anything)
	})

	t.Run("failed_case_get_mother_service_permanently_fails_when_attempts_exhausted", func(t *testing.T) {
		ctx := context.Background()
		mockOutbox := new(mocks.MockOutbox)
		mockRepo := new(mocks.MockMotherService)
		mockProvision := new(provisionProvider.MockProvisioningService)

		processor := NewOutboxProcessor(mockOutbox, mockRepo, mockProvision, 10, time.Second)

		item := &entity.Outbox{
			ID:            uint64(103),
			AggregateType: entity.OutboxAggregateTypeMotherService,
			AggregateID:   uint64(1),
			OperationType: entity.OutboxOperationProvisionMotherService,
			Attempts:      5,
			MaxAttempts:   5,
		}

		mockOutbox.On("ClaimPending", mock.Anything, 10).Return([]*entity.Outbox{item}, nil)
		mockRepo.On("GetByID", mock.Anything, item.AggregateID).Return(nil, errors.New("not found"))
		mockOutbox.On("MarkFailedPermanently", mock.Anything, item.ID, "not found").Return(nil)

		err := processor.ProcessPending(ctx)

		assert.NoError(t, err)
		mockOutbox.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockProvision.AssertNotCalled(t, "DeprovisionMotherService", mock.Anything, mock.Anything)
	})

	t.Run("failed_case_provision_reschedules_when_attempts_remain", func(t *testing.T) {
		ctx := context.Background()
		mockOutbox := new(mocks.MockOutbox)
		mockRepo := new(mocks.MockMotherService)
		mockProvision := new(provisionProvider.MockProvisioningService)

		processor := NewOutboxProcessor(mockOutbox, mockRepo, mockProvision, 10, time.Second)

		motherService := &entity.MotherService{ID: uint64(1), Name: "mother1"}
		item := &entity.Outbox{
			ID:            uint64(104),
			AggregateType: entity.OutboxAggregateTypeMotherService,
			AggregateID:   motherService.ID,
			OperationType: entity.OutboxOperationProvisionMotherService,
			Attempts:      2,
			MaxAttempts:   5,
		}

		mockOutbox.On("ClaimPending", mock.Anything, 10).Return([]*entity.Outbox{item}, nil)
		mockRepo.On("GetByID", mock.Anything, motherService.ID).Return(motherService, nil)
		mockProvision.On("ProvisionMotherService", mock.Anything, motherService).Return(errors.New("apply deployment failed"))
		mockOutbox.On("Reschedule", mock.Anything, item.ID, mock.AnythingOfType("time.Time"), "apply deployment failed").Return(nil)

		err := processor.ProcessPending(ctx)

		assert.NoError(t, err)
		mockOutbox.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockProvision.AssertExpectations(t)
		mockProvision.AssertNotCalled(t, "DeprovisionMotherService", mock.Anything, mock.Anything)
	})

	t.Run("failed_case_provision_permanently_fails_deprovisions_and_marks_mother_service_failed", func(t *testing.T) {
		ctx := context.Background()
		mockOutbox := new(mocks.MockOutbox)
		mockRepo := new(mocks.MockMotherService)
		mockProvision := new(provisionProvider.MockProvisioningService)

		processor := NewOutboxProcessor(mockOutbox, mockRepo, mockProvision, 10, time.Second)

		motherService := &entity.MotherService{ID: uint64(1), Name: "mother1"}
		item := &entity.Outbox{
			ID:            uint64(105),
			AggregateType: entity.OutboxAggregateTypeMotherService,
			AggregateID:   motherService.ID,
			OperationType: entity.OutboxOperationProvisionMotherService,
			Attempts:      5,
			MaxAttempts:   5,
		}

		mockOutbox.On("ClaimPending", mock.Anything, 10).Return([]*entity.Outbox{item}, nil)
		mockRepo.On("GetByID", mock.Anything, motherService.ID).Return(motherService, nil)
		mockProvision.On("ProvisionMotherService", mock.Anything, motherService).Return(errors.New("apply deployment failed"))
		mockProvision.On("DeprovisionMotherService", mock.Anything, motherService).Return(nil)
		mockRepo.On("SetStatus", mock.Anything, motherService.ID, entity.MotherServiceStatusFailed).Return(nil)
		mockOutbox.On("MarkFailedPermanently", mock.Anything, item.ID, "apply deployment failed").Return(nil)

		err := processor.ProcessPending(ctx)

		assert.NoError(t, err)
		mockOutbox.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockProvision.AssertExpectations(t)
	})

	t.Run("failed_case_provision_permanently_fails_deprovision_error_is_best_effort", func(t *testing.T) {
		ctx := context.Background()
		mockOutbox := new(mocks.MockOutbox)
		mockRepo := new(mocks.MockMotherService)
		mockProvision := new(provisionProvider.MockProvisioningService)

		processor := NewOutboxProcessor(mockOutbox, mockRepo, mockProvision, 10, time.Second)

		motherService := &entity.MotherService{ID: uint64(1), Name: "mother1"}
		item := &entity.Outbox{
			ID:            uint64(106),
			AggregateType: entity.OutboxAggregateTypeMotherService,
			AggregateID:   motherService.ID,
			OperationType: entity.OutboxOperationProvisionMotherService,
			Attempts:      5,
			MaxAttempts:   5,
		}

		mockOutbox.On("ClaimPending", mock.Anything, 10).Return([]*entity.Outbox{item}, nil)
		mockRepo.On("GetByID", mock.Anything, motherService.ID).Return(motherService, nil)
		mockProvision.On("ProvisionMotherService", mock.Anything, motherService).Return(errors.New("apply deployment failed"))
		mockProvision.On("DeprovisionMotherService", mock.Anything, motherService).Return(errors.New("cleanup failed"))
		mockRepo.On("SetStatus", mock.Anything, motherService.ID, entity.MotherServiceStatusFailed).Return(nil)
		mockOutbox.On("MarkFailedPermanently", mock.Anything, item.ID, "apply deployment failed").Return(nil)

		err := processor.ProcessPending(ctx)

		assert.NoError(t, err)
		mockOutbox.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockProvision.AssertExpectations(t)
	})
}
