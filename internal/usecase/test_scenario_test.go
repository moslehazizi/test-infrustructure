package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository/mocks"
	repoMocks "control-panel-service/internal/repository/mocks"
	"control-panel-service/internal/server/dto/request"
	"control-panel-service/pkg"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTestScenarioUsecase_Init(t *testing.T) {
	mockRepo := new(mocks.MockTestScenario)
	mockTestCatRepo := new(mocks.MockTestCategory)
	mockTestServiceConfig := new(mocks.MockTestServiceConfig)
	mockMotherService := new(mocks.MockMotherService)
	mockTestServiceRepo := new(mocks.MockTestServiceRepository)

	service := NewTestScenarioUsecase(
		getMockDB(t),
		mockRepo,
		mockTestCatRepo,
		mockTestServiceConfig,
		mockMotherService,
		mockTestServiceRepo,
	)
	assert.NotNil(t, service)

	assert.NotNil(t, service.testScenarioRepository)
}

func TestTestScenarioUsecase_Create(t *testing.T) {
	t.Run("failed_case_test_service_config_could_not_be_null", func(t *testing.T) {
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)

		sampleInt64 := int64(1)

		testSci := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			MaxTestServiceCount: &sampleInt64,
			TestServiceConfig:   nil,
			NumSteps:            2,
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
		}
		testCat := &entity.TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasNumSteps:            true,
		}
		motherService := &entity.MotherService{
			ID:                testSci.MotherServiceID,
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockMotherService.On("GetByID", mock.Anything, testSci.MotherServiceID).Return(motherService, nil)
		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(testCat, nil)

		err := service.Create(context.Background(), testSci)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrTestServiceConfigIsRequired)
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockRepo.AssertExpectations(t)
		mockTestCatRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})
	t.Run("success_case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)
		sampleInt := 1
		sampleInt64 := int64(1)
		expectedID := uint64(1)
		databaseName := "db1"
		databaseTableName := "factorial"

		testServiceConfig := &entity.TestServiceConfig{
			TestScenarioID:         expectedID,
			MaxRequests:            1,
			MaxDuration:            1,
			RequestDelayDuration:   &sampleInt,
			FixedTestNumber:        &sampleInt,
			DatabaseName:           databaseName,
			DatabaseTableName:      databaseTableName,
			IncreaseFixedInput:     1,
			ExecNumMultiFixedInput: 1,
		}
		testCat := &entity.TestCategory{
			ID:                     1,
			Name:                   entity.STRESS,
			HasMaxTestServiceCount: true,
			HasNumSteps:            true,
		}
		testSci := &entity.TestScenario{
			ID:                  uint64(4),
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			MaxTestServiceCount: &sampleInt64,
			TestServiceConfig:   testServiceConfig,
			TestCategory:        testCat,
			NumSteps:            2,
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
		}
		motherService := &entity.MotherService{
			ID:                testSci.MotherServiceID,
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockMotherService.On("GetByID", mock.Anything, testSci.MotherServiceID).Return(motherService, nil)
		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(testCat, nil)
		mockRepo.On("CreateScenarioAndConfig", mock.Anything, testSci).Return(nil)

		err := service.Create(ctx, testSci)

		assert.Nil(t, err)
		mockRepo.AssertExpectations(t)
		mockTestCatRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})
	t.Run("failed_case_when_create_test_scenario", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)
		databaseName := "db1"
		databaseTableName := "factorial"

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)
		sampleInt := 1
		sampleInt64 := int64(1)
		testServiceCfg := &entity.TestServiceConfig{
			MaxRequests:            1,
			MaxDuration:            1,
			RequestDelayDuration:   &sampleInt,
			FixedTestNumber:        &sampleInt,
			DatabaseName:           databaseName,
			DatabaseTableName:      databaseTableName,
			IncreaseFixedInput:     1,
			ExecNumMultiFixedInput: 1,
		}
		testSci := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			MaxTestServiceCount: &sampleInt64,
			TestServiceConfig:   testServiceCfg,
			NumSteps:            2,
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
		}
		testCat := &entity.TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasNumSteps:            true,
		}
		motherService := &entity.MotherService{
			ID:                testSci.MotherServiceID,
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockMotherService.On("GetByID", mock.Anything, testSci.MotherServiceID).Return(motherService, nil)
		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(testCat, nil)
		mockRepo.On("CreateScenarioAndConfig", mock.Anything, testSci).Return(errors.New("something went wrong"))

		err := service.Create(ctx, testSci)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToCreateTestScenario)
		mockRepo.AssertCalled(t, "CreateScenarioAndConfig", mock.Anything, testSci)
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockRepo.AssertExpectations(t)
		mockTestCatRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})
	t.Run("failed_case_validation_error_max_test_service_count_less_than_one", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)

		sampleInt64 := int64(-1)
		testSci := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			MaxTestServiceCount: &sampleInt64,
			NumSteps:            2,
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
		}
		testCat := &entity.TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasNumSteps:            true,
		}
		motherService := &entity.MotherService{
			ID:                testSci.MotherServiceID,
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockMotherService.On("GetByID", mock.Anything, testSci.MotherServiceID).Return(motherService, nil)
		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(testCat, nil)

		err := service.Create(ctx, testSci)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMaxTestServiceCountLessThanOne)
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockTestCatRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})
	t.Run("failed_case_error_get_test_category_not_found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)

		testSci := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			NumSteps:            2,
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
		}
		motherService := &entity.MotherService{
			ID:                testSci.MotherServiceID,
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockMotherService.On("GetByID", mock.Anything, testSci.MotherServiceID).Return(motherService, nil)

		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(nil, pkg.ErrTestCategoryNotFound)

		err := service.Create(ctx, testSci)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrTestCategoryNotFound)
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockTestCatRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})
	t.Run("failed_case_error_get_test_category_unknown", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)
		testSci := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			NumSteps:            2,
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
		}
		motherService := &entity.MotherService{
			ID:                testSci.MotherServiceID,
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockMotherService.On("GetByID", mock.Anything, testSci.MotherServiceID).Return(motherService, nil)

		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(nil, errors.New("error happened"))

		err := service.Create(ctx, testSci)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestCategory)
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockTestCatRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})
	t.Run("failed_case_validation_error_max_service_count_not_set", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)
		testSci := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			NumSteps:            2,
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
		}
		testCat := &entity.TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: true,
			HasNumSteps:            true,
		}
		motherService := &entity.MotherService{
			ID:                testSci.MotherServiceID,
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockMotherService.On("GetByID", mock.Anything, testSci.MotherServiceID).Return(motherService, nil)
		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(testCat, nil)

		err := service.Create(ctx, testSci)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMaxTestServiceCountNotSet)
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockTestCatRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})
	t.Run("failed_case_validation_error_test_service_config", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)

		databaseName := "db1"
		databaseTableName := "factorial"
		testSci := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
			TestServiceConfig: &entity.TestServiceConfig{
				MaxRequests:            1,
				MaxDuration:            0,
				BadValueRate:           -1,
				DatabaseName:           databaseName,
				DatabaseTableName:      databaseTableName,
				IncreaseFixedInput:     1,
				ExecNumMultiFixedInput: 1,
			},
			NumSteps: 2,
		}
		testCat := &entity.TestCategory{
			ID:                     testSci.TestCategoryID,
			Name:                   "load",
			Label:                  "my load",
			HasMaxTestServiceCount: false,
			HasNumSteps:            true,
		}
		motherService := &entity.MotherService{
			ID:                testSci.MotherServiceID,
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockMotherService.On("GetByID", mock.Anything, testSci.MotherServiceID).Return(motherService, nil)
		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(testCat, nil)

		err := service.Create(ctx, testSci)

		assert.Error(t, err)
		assert.Equal(t, err, fmt.Errorf("%w: %w", pkg.ErrFailedToValidateTestSvcCfg, pkg.ErrInvalidBadValueRate))
		mockTestCatRepo.AssertCalled(t, "GetByID", mock.Anything, testSci.TestCategoryID)
		mockTestCatRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})
	t.Run("failed_case_mother_service_not_found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockMotherService := new(mocks.MockMotherService)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)

		sampleInt := 1
		sampleInt64 := int64(1)
		expectedID := uint64(1)
		databaseName := "db1"
		databaseTableName := "factorial"

		testServiceConfig := &entity.TestServiceConfig{
			TestScenarioID:         expectedID,
			MaxRequests:            1,
			MaxDuration:            1,
			RequestDelayDuration:   &sampleInt,
			FixedTestNumber:        &sampleInt,
			DatabaseName:           databaseName,
			DatabaseTableName:      databaseTableName,
			IncreaseFixedInput:     1,
			ExecNumMultiFixedInput: 1,
		}
		testSci := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     0,
			MaxTestServiceCount: &sampleInt64,
			TestServiceConfig:   testServiceConfig,
			NumSteps:            2,
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
		}

		mockMotherService.On("GetByID", mock.Anything, testSci.MotherServiceID).Return(nil, pkg.ErrMotherServiceNotFound)

		err := service.Create(ctx, testSci)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMotherServiceNotFound)
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
		mockTestCatRepo.AssertExpectations(t)
	})
	t.Run("success_case_when_test_category_has_num_steps_is_false", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)
		sampleInt := 1
		sampleInt64 := int64(1)
		expectedID := uint64(1)
		databaseName := "db1"
		databaseTableName := "factorial"

		testServiceConfig := &entity.TestServiceConfig{
			TestScenarioID:         expectedID,
			MaxRequests:            1,
			MaxDuration:            1,
			RequestDelayDuration:   &sampleInt,
			FixedTestNumber:        &sampleInt,
			DatabaseName:           databaseName,
			DatabaseTableName:      databaseTableName,
			IncreaseFixedInput:     1,
			ExecNumMultiFixedInput: 1,
		}
		testCat := &entity.TestCategory{
			ID:                     1,
			Name:                   entity.STRESS,
			HasMaxTestServiceCount: true,
			HasNumSteps:            false,
		}
		testSci := &entity.TestScenario{
			ID:                  uint64(4),
			Name:                "load1",
			TestCategoryID:      uint64(2),
			MotherServiceID:     uint64(1),
			MaxTestServiceCount: &sampleInt64,
			TestServiceConfig:   testServiceConfig,
			TestCategory:        testCat,
			NumSteps:            2,
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
		}
		motherService := &entity.MotherService{
			ID:                testSci.MotherServiceID,
			Name:              "mother1",
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		mockMotherService.On("GetByID", mock.Anything, testSci.MotherServiceID).Return(motherService, nil)
		mockTestCatRepo.On("GetByID", mock.Anything, testSci.TestCategoryID).Return(testCat, nil)
		mockRepo.On("CreateScenarioAndConfig", mock.Anything, testSci).Return(nil)

		err := service.Create(ctx, testSci)

		assert.Nil(t, err)
		assert.Equal(t, testSci.NumSteps, int64(1))
		mockRepo.AssertExpectations(t)
		mockTestCatRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})
}

func TestTestScenarioUsecase_GetByID(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)
		someTime := time.Date(2026, 01, 13, 14, 10, 0, 0, time.Now().Location())
		num := 10
		sampleInt64 := int64(1)
		sampleID := uint64(4)

		expectedTestScenario := &entity.TestScenario{
			ID:              sampleID,
			Name:            "load1",
			TestCategoryID:  uint64(1),
			MotherServiceID: uint64(2),
			NumSteps:        2,
			TestCategory: &entity.TestCategory{
				ID:                     1,
				Name:                   "load",
				Label:                  "Load Test",
				CreatedAt:              someTime,
				UpdatedAt:              someTime,
				HasMaxTestServiceCount: true,
				HasNumSteps:            true,
			},
			MotherService: &entity.MotherService{
				ID:   2,
				Name: "m2",
			},
			Status:              entity.ScenarioStatusSucceed,
			MaxTestServiceCount: &sampleInt64,
			TestServiceConfig: &entity.TestServiceConfig{
				ID:                    100,
				TestScenarioID:        1,
				CreatedAt:             someTime,
				UpdatedAt:             someTime,
				MaxRequests:           100,
				MaxDuration:           0,
				RequestDelayDuration:  nil,
				RandomRequestDelayMin: nil,
				RandomRequestDelayMax: nil,
				FixedTestNumber:       &num,
				RandomTestNumberMin:   nil,
				RandomTestNumberMax:   nil,
				BadValueRate:          0,
				NegativeValueRate:     0,
				RealValueRate:         0,
				ZeroValueRate:         0,
				StringValueRate:       0,
				LongStringValueRate:   0,
				NullValueRate:         0,
			},
			Editable: true,
		}

		mockRepo.On("GetByID", mock.Anything, sampleID).Return(expectedTestScenario, nil)

		result, err := service.GetByID(ctx, sampleID)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedTestScenario, result)
		assert.Equal(t, expectedTestScenario.ID, result.ID)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_not_found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)
		sampleID := uint64(4)

		mockRepo.On("GetByID", mock.Anything, sampleID).Return(nil, pkg.ErrTestScenarioNotFound)

		result, err := service.GetByID(ctx, sampleID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, pkg.ErrTestScenarioNotFound)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_repository_unknown_error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)
		sampleID := uint64(4)

		mockRepo.On("GetByID", mock.Anything, sampleID).Return(nil, errors.New("error happened"))

		result, err := service.GetByID(ctx, sampleID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestScenario)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case_repository_unknown_error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)
		sampleID := uint64(4)

		mockRepo.On("GetByID", mock.Anything, sampleID).Return(nil, errors.New("error happened"))

		result, err := service.GetByID(ctx, sampleID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestScenario)
		mockRepo.AssertCalled(t, "GetByID", mock.Anything, sampleID)
		mockRepo.AssertExpectations(t)
	})
}

func TestTestScenarioUsecase_GetPaginated(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)

		sampleInt64 := int64(2)
		pagReq := entity.TestScenarioPaginationRequest{
			Page:    2,
			PerPage: 2,
		}
		count := int64(2)
		expectedResult := []*entity.TestScenario{
			{
				ID:              uint64(2),
				Name:            "load1",
				TestCategoryID:  uint64(1),
				MotherServiceID: uint64(2),
				NumSteps:        2,
				TestCategory: &entity.TestCategory{
					ID:                     1,
					Name:                   "load",
					Label:                  "Load Test",
					CreatedAt:              time.Now(),
					UpdatedAt:              time.Now(),
					HasMaxTestServiceCount: true,
					HasNumSteps:            true,
				},
				MotherService: &entity.MotherService{
					ID:   2,
					Name: "m2",
				},
				Status:              entity.ScenarioStatusSucceed,
				MaxTestServiceCount: &sampleInt64,
				Editable:            true,
			},
			{
				ID:              uint64(1),
				Name:            "smoke2",
				TestCategoryID:  uint64(1),
				MotherServiceID: uint64(2),
				TestCategory: &entity.TestCategory{
					ID:                     1,
					Name:                   "load",
					Label:                  "Load Test",
					CreatedAt:              time.Now(),
					UpdatedAt:              time.Now(),
					HasMaxTestServiceCount: true,
					HasNumSteps:            true,
				},
				MotherService: &entity.MotherService{
					ID:   2,
					Name: "m2",
				},
				Status:              entity.ScenarioStatusSucceed,
				MaxTestServiceCount: &sampleInt64,
				Editable:            true,
			},
		}

		mockRepo.On("GetPaginated", mock.Anything, pagReq).Return(expectedResult, count, nil)

		result, total, err := service.GetPaginated(ctx, pagReq)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, total, count)
		assert.Equal(t, len(result), 2)

		assert.Equal(t, expectedResult, result)

		mockRepo.AssertCalled(t, "GetPaginated", mock.Anything, pagReq)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed_case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(mocks.MockTestScenario)
		mockTestCatRepo := new(mocks.MockTestCategory)
		mockTestServiceConfig := new(mocks.MockTestServiceConfig)
		mockMotherService := new(mocks.MockMotherService)
		mockTestServiceRepo := new(mocks.MockTestServiceRepository)

		service := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)
		pagReq := entity.TestScenarioPaginationRequest{
			Page:    2,
			PerPage: 2,
		}

		mockRepo.On("GetPaginated", mock.Anything, pagReq).Return(nil, int64(0), errors.New("error happened"))

		result, count, err := service.GetPaginated(ctx, pagReq)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, count, int64(0))
		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestScenarios)
		mockRepo.AssertCalled(t, "GetPaginated", mock.Anything, pagReq)
		mockRepo.AssertExpectations(t)
	})
}

func TestTestScenarioUsecase_Update(t *testing.T) {
	sampleCategory := &entity.TestCategory{
		ID:                     1,
		Name:                   "stress",
		HasMaxTestServiceCount: false,
		HasNumSteps:            true,
	}

	baseExistingScenario := func() *entity.TestScenario {
		return &entity.TestScenario{
			ID:              1,
			Name:            "Old Name",
			MotherServiceID: 10,
			MotherService:   &entity.MotherService{ID: 10, Name: "Old Mother"},
			TestCategory:    sampleCategory,
			Status:          entity.ScenarioStatusReady,
			TestServiceConfig: &entity.TestServiceConfig{
				ID: 100,
			},
			NumSteps: 2,
		}
	}

	baseMotherService := &entity.MotherService{ID: 12}

	// baseConfig := func() *entity.TestServiceConfig {
	// 	return &entity.TestServiceConfig{ID: 100}
	// }

	fixedTestNumber := 5
	requestDelay := 180

	baseRequest := func() *request.TestScenarioUpdateRequest {
		return &request.TestScenarioUpdateRequest{
			ID:                  1,
			Name:                "Updated Load Test Scenario",
			MotherServiceID:     12,
			NumSteps:            2,
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   1,
			Config: &request.TestServiceConfigRequest{
				MaxRequests:            1000,
				MaxDuration:            300,
				RequestDelayDuration:   &requestDelay,
				FixedTestNumber:        &fixedTestNumber,
				BadValueRate:           0,
				NegativeValueRate:      0,
				RealValueRate:          0,
				ZeroValueRate:          0,
				StringValueRate:        0,
				LongStringValueRate:    0,
				NullValueRate:          0,
				DatabaseName:           "test_db",
				DatabaseTableName:      "transactions",
				IncreaseFixedInput:     1,
				ExecNumMultiFixedInput: 1,
			},
		}
	}
	t.Run("success", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)
		mockRepo.On("UpdateScenarioAndConfig", mock.Anything, mock.AnythingOfType("*entity.TestScenario")).Return(nil)

		err := svc.Update(context.Background(), baseRequest())

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})

	t.Run("failed_case_pending_or_ready_scenarios_can_be_updated", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)

		sampleScenario := &entity.TestScenario{
			ID:              1,
			Name:            "Old Name",
			Status:          entity.ScenarioStatusPaused,
			MotherServiceID: 10,
			MotherService:   &entity.MotherService{ID: 10, Name: "Old Mother"},
			TestCategory:    sampleCategory,
			TestServiceConfig: &entity.TestServiceConfig{
				ID: 100,
			},
			NumSteps: 2,
		}

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(sampleScenario, nil)

		err := svc.Update(context.Background(), baseRequest())

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrScenariosCanNotBeUpdated)

		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})

	t.Run("success_category_requires_all_optional_fields", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)

		catWithAll := &entity.TestCategory{
			ID:                     2,
			Name:                   "advanced",
			HasMaxTestServiceCount: true,
			HasNumSteps:            true,
		}
		maxCount := int64(10)

		existing := baseExistingScenario()
		existing.TestCategory = catWithAll

		req := baseRequest()
		req.MaxTestServiceCount = &maxCount

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(existing, nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)
		mockRepo.On("UpdateScenarioAndConfig", mock.Anything, mock.AnythingOfType("*entity.TestScenario")).Return(nil)

		err := svc.Update(context.Background(), req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})

	t.Run("success_random_delay_and_random_test_number_config", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)

		minDelay := 100
		maxDelay := 500
		minTest := 1
		maxTest := 10

		req := baseRequest()
		req.Config.RequestDelayDuration = nil
		req.Config.RandomRequestDelayMin = &minDelay
		req.Config.RandomRequestDelayMax = &maxDelay
		req.Config.FixedTestNumber = nil
		req.Config.RandomTestNumberMin = &minTest
		req.Config.RandomTestNumberMax = &maxTest

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)
		mockRepo.On("UpdateScenarioAndConfig", mock.Anything, mock.AnythingOfType("*entity.TestScenario")).Return(nil)

		err := svc.Update(context.Background(), req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})

	t.Run("success_bad_value_rate_enabled_with_valid_sub_rates_summing_to_100", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)

		req := baseRequest()
		req.Config.BadValueRate = 30
		req.Config.NegativeValueRate = 20
		req.Config.RealValueRate = 20
		req.Config.ZeroValueRate = 20
		req.Config.StringValueRate = 10
		req.Config.LongStringValueRate = 10
		req.Config.NullValueRate = 20

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)
		mockRepo.On("UpdateScenarioAndConfig", mock.Anything, mock.AnythingOfType("*entity.TestScenario")).Return(nil)

		err := svc.Update(context.Background(), req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})

	t.Run("error_nil_config_returns_ErrTestServiceConfigIsRequired_immediately", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)

		req := baseRequest()
		req.Config = nil

		err := svc.Update(context.Background(), req)

		assert.ErrorIs(t, err, pkg.ErrTestServiceConfigIsRequired)
		mockRepo.AssertNotCalled(t, "GetByID")
		mockMotherService.AssertNotCalled(t, "GetByID")
	})

	t.Run("error_scenario_not_found", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(nil, errors.New("not found"))

		err := svc.Update(context.Background(), baseRequest())

		assert.ErrorIs(t, err, pkg.ErrTestScenarioNotFound)
		mockMotherService.AssertNotCalled(t, "GetByID")
		mockRepo.AssertExpectations(t)
	})

	t.Run("error_mother_service_not_found", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(nil, errors.New("not found"))

		err := svc.Update(context.Background(), baseRequest())

		assert.ErrorIs(t, err, pkg.ErrFailedToGetMotherService)
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})

	t.Run("error_scenario_validation_MaxTestServiceCount_<_1", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)

		cat := &entity.TestCategory{HasMaxTestServiceCount: true, HasNumSteps: true}
		existing := baseExistingScenario()
		existing.TestCategory = cat

		invalidCount := int64(0)
		req := baseRequest()
		req.MaxTestServiceCount = &invalidCount

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(existing, nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)

		err := svc.Update(context.Background(), req)

		assert.ErrorIs(t, err, pkg.ErrFailedToUpdateTestScenario)
		assert.ErrorIs(t, err, pkg.ErrMaxTestServiceCountLessThanOne)
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})

	t.Run("error_test_service_config_not_preloaded_(ID=0)", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)

		existing := baseExistingScenario()
		existing.TestServiceConfig = &entity.TestServiceConfig{ID: 0}

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(existing, nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)

		err := svc.Update(context.Background(), baseRequest())

		assert.ErrorIs(t, err, pkg.ErrFailedToGetTestServiceConfig)
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})

	t.Run("error_config_validation_fixed_delay_and_random_delay_both_set", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)

		minDelay := 100
		maxDelay := 500

		req := baseRequest()
		// RequestDelayDuration already set — mixing with random is invalid
		req.Config.RandomRequestDelayMin = &minDelay
		req.Config.RandomRequestDelayMax = &maxDelay

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)

		err := svc.Update(context.Background(), req)

		assert.ErrorIs(t, err, pkg.ErrFailedToValidateTestSvcCfg)
		assert.ErrorIs(t, err, pkg.ErrInvalidRequestDelayDurationConfig)
		mockRepo.AssertNotCalled(t, "Update")
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})

	t.Run("error_config_validation_random_delay_min_>=_max", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)

		minDelay := 500
		maxDelay := 100 // min > max

		req := baseRequest()
		req.Config.RequestDelayDuration = nil
		req.Config.RandomRequestDelayMin = &minDelay
		req.Config.RandomRequestDelayMax = &maxDelay

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)

		err := svc.Update(context.Background(), req)

		assert.ErrorIs(t, err, pkg.ErrFailedToValidateTestSvcCfg)
		assert.ErrorIs(t, err, pkg.ErrMinDelayDurationMoreThanMax)
		mockRepo.AssertNotCalled(t, "Update")
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})

	t.Run("error_config_validation_no_test_number_config_set", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)

		req := baseRequest()
		req.Config.FixedTestNumber = nil
		req.Config.RandomTestNumberMin = nil
		req.Config.RandomTestNumberMax = nil

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)

		err := svc.Update(context.Background(), req)

		assert.ErrorIs(t, err, pkg.ErrFailedToValidateTestSvcCfg)
		assert.ErrorIs(t, err, pkg.ErrInvalidTestNumberConfig)
		mockRepo.AssertNotCalled(t, "Update")
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})

	t.Run("error_config_validation_fixed_test_number_<=_0", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)

		zero := 0
		req := baseRequest()
		req.Config.FixedTestNumber = &zero

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)

		err := svc.Update(context.Background(), req)

		assert.ErrorIs(t, err, pkg.ErrFailedToValidateTestSvcCfg)
		assert.ErrorIs(t, err, pkg.ErrInvalidFixedTestNumber)
		mockRepo.AssertNotCalled(t, "Update")
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})

	t.Run("error_config_validation_fixed_and_random_test_number_both_set", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)

		minTest := 1
		maxTest := 10

		req := baseRequest()
		// FixedTestNumber already set in baseRequest — mixing with random is invalid
		req.Config.RandomTestNumberMin = &minTest
		req.Config.RandomTestNumberMax = &maxTest

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)

		err := svc.Update(context.Background(), req)

		assert.ErrorIs(t, err, pkg.ErrFailedToValidateTestSvcCfg)
		assert.ErrorIs(t, err, pkg.ErrInvalidFixedTestNumberConfig)
		mockRepo.AssertNotCalled(t, "Update")
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})

	t.Run("error_config_validation_random_test_number_min_>=_max", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)

		minTest := 10
		maxTest := 1 // min > max

		req := baseRequest()
		req.Config.FixedTestNumber = nil
		req.Config.RandomTestNumberMin = &minTest
		req.Config.RandomTestNumberMax = &maxTest

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)

		err := svc.Update(context.Background(), req)

		assert.ErrorIs(t, err, pkg.ErrFailedToValidateTestSvcCfg)
		assert.ErrorIs(t, err, pkg.ErrMinRandomTestNumberMoreThanMax)
		mockRepo.AssertNotCalled(t, "Update")
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})

	t.Run("error_config_validation_bad_value_rate_>_0_but_sub_rates_do_not_sum_to_100", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)

		req := baseRequest()
		req.Config.BadValueRate = 30
		req.Config.NegativeValueRate = 10 // sum = 10, not 100

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)

		err := svc.Update(context.Background(), req)

		assert.ErrorIs(t, err, pkg.ErrFailedToValidateTestSvcCfg)
		assert.ErrorIs(t, err, pkg.ErrInvalid100SumOfBadValues)
		mockRepo.AssertNotCalled(t, "Update")
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})

	t.Run("error_config_validation_bad_value_rate_==_0_but_sub_rates_are_non_zero", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)

		req := baseRequest()
		req.Config.BadValueRate = 0
		req.Config.NegativeValueRate = 10 // must be 0 when BadValueRate == 0

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)

		err := svc.Update(context.Background(), req)

		assert.ErrorIs(t, err, pkg.ErrFailedToValidateTestSvcCfg)
		assert.ErrorIs(t, err, pkg.ErrInvalidZeroSumOfBadValues)
		mockRepo.AssertNotCalled(t, "Update")
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})

	t.Run("error_config_validation_empty_database_name", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)

		req := baseRequest()
		req.Config.DatabaseName = ""

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)

		err := svc.Update(context.Background(), req)

		assert.ErrorIs(t, err, pkg.ErrFailedToValidateTestSvcCfg)
		assert.ErrorIs(t, err, pkg.ErrInvalidDatabaseName)
		mockRepo.AssertNotCalled(t, "Update")
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})

	t.Run("error_config_validation_empty_database_table_name", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)

		req := baseRequest()
		req.Config.DatabaseTableName = ""

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)

		err := svc.Update(context.Background(), req)

		assert.ErrorIs(t, err, pkg.ErrFailedToValidateTestSvcCfg)
		assert.ErrorIs(t, err, pkg.ErrInvalidDatabaseTableName)
		mockRepo.AssertNotCalled(t, "Update")
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})

	t.Run("error_scenario_repository_update_fails", func(t *testing.T) {
		mockRepo := new(repoMocks.MockTestScenario)
		mockTestCatRepo := new(repoMocks.MockTestCategory)
		mockTestServiceConfig := new(repoMocks.MockTestServiceConfig)
		mockMotherService := new(repoMocks.MockMotherService)
		mockTestServiceRepo := new(repoMocks.MockTestServiceRepository)

		svc := NewTestScenarioUsecase(
			getMockDB(t),
			mockRepo,
			mockTestCatRepo,
			mockTestServiceConfig,
			mockMotherService,
			mockTestServiceRepo,
		)

		mockRepo.On("GetByID", mock.Anything, uint64(1)).Return(baseExistingScenario(), nil)
		mockMotherService.On("GetByID", mock.Anything, uint64(12)).Return(baseMotherService, nil)
		mockRepo.On("UpdateScenarioAndConfig", mock.Anything, mock.AnythingOfType("*entity.TestScenario")).Return(errors.New("update failed"))

		err := svc.Update(context.Background(), baseRequest())

		assert.ErrorIs(t, err, pkg.ErrFailedToUpdateTestScenario)
		mockRepo.AssertExpectations(t)
		mockMotherService.AssertExpectations(t)
	})
}
