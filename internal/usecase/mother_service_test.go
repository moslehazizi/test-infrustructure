package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMotherService(t *testing.T) {
	mockRepo := new(mocks.MockMotherService)
	service := NewMotherService(mockRepo)

	assert.NotNil(t, service)
}

func TestMotherServiceUsecase_Create(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockMotherService)
		service := NewMotherService(mockRepo)

		sampleMS := &entity.MotherService{
			Name:               "mother1",
			ProvisioningStatus: entity.ProvisioningStatusFailed,
			DatabaseName:       "db1",
			DatabaseTableName:  "factorial",
		}

		err := service.Create(ctx, sampleMS)

		assert.NoError(t, err)
	})
}
