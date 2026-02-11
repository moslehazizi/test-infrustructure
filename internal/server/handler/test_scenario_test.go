package handler

import (
	"control-panel-service/config"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/server/dto/response"
	"control-panel-service/internal/usecase/mocks"
	"control-panel-service/pkg"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTestServiceConfigHandler_New(t *testing.T) {
	mockSrv := new(mocks.MockTestScenario)

	handler := NewTestScenarioHandler(mockSrv)

	assert.NotNil(t, handler)

	assert.NotNil(t, handler.testScenario)
}

func TestTestScenarioHandler_Create(t *testing.T) {
	t.Run("success case - bad values set", func(t *testing.T) {
		mockSvc := new(mocks.MockTestScenario)
		handler := NewTestScenarioHandler(mockSvc)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios", handler.Create())

		reqBody := `{
			"name": "load1",
			"test_category_id": 1,
			"mother_service_id": 1,
			"max_test_service_count": 1,
			"execution_duration": 1,
			"auto_step_change_rate": 1,
			"test_service_config": {
				"max_requests": 1,
				"max_duration": 1,
				"request_delay_duration": null,
				"random_request_delay_min": 10,
				"random_request_delay_max": 20,
				"fixed_test_number": null,
				"random_test_number_min": 10,
				"random_test_number_max": 20,
				"bad_value_rate": 50,
				"negative_value_rate": 10,
				"real_value_rate": 20,
				"zero_value_rate": 40,
				"string_value_rate": 10,
				"long_string_value_rate": 10,
				"null_value_rate": 10 , 
		        "database_name": "test_service_db",
        		"database_table_name": "test_service_table"
			}
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleInt64 := int64(1)
		min := 10
		max := 20
		dur := time.Millisecond * 1
		databaseName := "test_service_db"
		databaseTableName := "test_service_table"
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt64,
			ExecutionDuration:   &dur,
			AutoStepChangeRate:  &sampleInt64,
			TestServiceConfig: &entity.TestServiceConfig{
				MaxRequests:           sampleInt,
				MaxDuration:           dur,
				RandomRequestDelayMin: &min,
				RandomRequestDelayMax: &max,
				RandomTestNumberMin:   &min,
				RandomTestNumberMax:   &max,
				BadValueRate:          50,
				NegativeValueRate:     10,
				RealValueRate:         20,
				ZeroValueRate:         40,
				StringValueRate:       10,
				LongStringValueRate:   10,
				NullValueRate:         10,
				DatabaseName:          databaseName,
				DatabaseTableName:     databaseTableName,
			},
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result response.SuccessResponse
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, pkg.CreateTestScenarioSuccessfully, result.Message)
		mockSvc.AssertExpectations(t)
	})
	t.Run("success case - random test number ", func(t *testing.T) {
		mockSvc := new(mocks.MockTestScenario)
		handler := NewTestScenarioHandler(mockSvc)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios", handler.Create())

		reqBody := `{
			"name": "load1",
			"test_category_id": 1,
			"mother_service_id": 1,
			"max_test_service_count": 1,
			"execution_duration": 1,
			"auto_step_change_rate": 1,
			"test_service_config": {
				"max_requests": 1,
				"max_duration": 1,
				"request_delay_duration": null,
				"random_request_delay_min": 10,
				"random_request_delay_max": 20,
				"fixed_test_number": null,
				"random_test_number_min": 10,
				"random_test_number_max": 20,
				"bad_value_rate": 0,
				"negative_value_rate": 0,
				"real_value_rate": 0,
				"zero_value_rate": 0,
				"string_value_rate": 0,
				"long_string_value_rate": 0,
				"null_value_rate": 0, 
				"database_name": "test_service_db",
        		"database_table_name": "test_service_table"
			}
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleInt64 := int64(1)
		min := 10
		max := 20
		dur := time.Millisecond * 1
		databaseName := "test_service_db"
		databaseTableName := "test_service_table"
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt64,
			ExecutionDuration:   &dur,
			AutoStepChangeRate:  &sampleInt64,
			TestServiceConfig: &entity.TestServiceConfig{
				MaxRequests:           sampleInt,
				MaxDuration:           dur,
				RandomRequestDelayMin: &min,
				RandomRequestDelayMax: &max,
				RandomTestNumberMin:   &min,
				RandomTestNumberMax:   &max,
				DatabaseName:          databaseName,
				DatabaseTableName:     databaseTableName,
			},
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result response.SuccessResponse
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, pkg.CreateTestScenarioSuccessfully, result.Message)
		mockSvc.AssertExpectations(t)
	})
	t.Run("success case - random request delay ", func(t *testing.T) {
		mockSvc := new(mocks.MockTestScenario)
		handler := NewTestScenarioHandler(mockSvc)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios", handler.Create())

		reqBody := `{
			"name": "load1",
			"test_category_id": 1,
			"mother_service_id": 1,
			"max_test_service_count": 1,
			"execution_duration": 1,
			"auto_step_change_rate": 1,
			"test_service_config": {
				"max_requests": 1,
				"max_duration": 1,
				"request_delay_duration": null,
				"random_request_delay_min": 10,
				"random_request_delay_max": 20,
				"fixed_test_number": 1,
				"random_test_number_min": null,
				"random_test_number_max": null,
				"bad_value_rate": 0,
				"negative_value_rate": 0,
				"real_value_rate": 0,
				"zero_value_rate": 0,
				"string_value_rate": 0,
				"long_string_value_rate": 0,
				"null_value_rate": 0,
				"database_name": "test_service_db",
        		"database_table_name": "test_service_table"
			}
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleInt64 := int64(1)
		min := 10
		max := 20
		dur := time.Millisecond * 1
		databaseName := "test_service_db"
		databaseTableName := "test_service_table"
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt64,
			ExecutionDuration:   &dur,
			AutoStepChangeRate:  &sampleInt64,
			TestServiceConfig: &entity.TestServiceConfig{
				MaxRequests:           sampleInt,
				MaxDuration:           dur,
				RandomRequestDelayMin: &min,
				RandomRequestDelayMax: &max,
				FixedTestNumber:       &sampleInt,
				DatabaseName:          databaseName,
				DatabaseTableName:     databaseTableName,
			},
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result response.SuccessResponse
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, pkg.CreateTestScenarioSuccessfully, result.Message)
		mockSvc.AssertExpectations(t)
	})
	t.Run("success case - request delay and fixed number set", func(t *testing.T) {
		mockSvc := new(mocks.MockTestScenario)
		handler := NewTestScenarioHandler(mockSvc)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios", handler.Create())

		reqBody := `{
			"name": "load1",
			"test_category_id": 1,
			"mother_service_id": 1,
			"max_test_service_count": 1,
			"execution_duration": 1,
			"auto_step_change_rate": 1,
			"test_service_config": {
				"max_requests": 1,
				"max_duration": 1,
				"request_delay_duration": 1,
				"random_request_delay_min": null,
				"random_request_delay_max": null,
				"fixed_test_number": 1,
				"random_test_number_min": null,
				"random_test_number_max": null,
				"bad_value_rate": 0,
				"negative_value_rate": 0,
				"real_value_rate": 0,
				"zero_value_rate": 0,
				"string_value_rate": 0,
				"long_string_value_rate": 0,
				"null_value_rate": 0,
				"database_name": "test_service_db",
        		"database_table_name": "test_service_table"
			}
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleInt64 := int64(1)
		dur := time.Millisecond * 1
		databaseName := "test_service_db"
		databaseTableName := "test_service_table"
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt64,
			ExecutionDuration:   &dur,
			AutoStepChangeRate:  &sampleInt64,
			TestServiceConfig: &entity.TestServiceConfig{
				MaxRequests:          sampleInt,
				MaxDuration:          dur,
				RequestDelayDuration: &sampleInt,
				FixedTestNumber:      &sampleInt,
				DatabaseName:         databaseName,
				DatabaseTableName:    databaseTableName,
			},
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result response.SuccessResponse
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, pkg.CreateTestScenarioSuccessfully, result.Message)
		mockSvc.AssertExpectations(t)
	})
	t.Run("success case - don't fill null fields in test scenario", func(t *testing.T) {
		mockSvc := new(mocks.MockTestScenario)
		handler := NewTestScenarioHandler(mockSvc)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios", handler.Create())

		reqBody := `{
			"name": "load1",
			"test_category_id": 1,
			"mother_service_id": 1,
			"max_test_service_count": null,
			"execution_duration": null,
			"auto_step_change_rate": null,
			"test_service_config": {
				"max_requests": 1,
				"max_duration": 1,
				"request_delay_duration": 1,
				"random_request_delay_min": null,
				"random_request_delay_max": null,
				"fixed_test_number": 1,
				"random_test_number_min": null,
				"random_test_number_max": null,
				"bad_value_rate": 0,
				"negative_value_rate": 0,
				"real_value_rate": 0,
				"zero_value_rate": 0,
				"string_value_rate": 0,
				"long_string_value_rate": 0,
				"null_value_rate": 0,
				"database_name": "test_service_db",
        		"database_table_name": "test_service_table"
			}
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		databaseName := "test_service_db"
		databaseTableName := "test_service_table"
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:            "load1",
			TestCategoryID:  sampleUin64,
			MotherServiceID: sampleUin64,
			TestServiceConfig: &entity.TestServiceConfig{
				MaxRequests:          sampleInt,
				MaxDuration:          time.Millisecond * 1,
				RequestDelayDuration: &sampleInt,
				FixedTestNumber:      &sampleInt,
				DatabaseName:         databaseName,
				DatabaseTableName:    databaseTableName,
			},
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result response.SuccessResponse
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, pkg.CreateTestScenarioSuccessfully, result.Message)
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - body parser bad request", func(t *testing.T) {
		mockSvc := new(mocks.MockTestScenario)
		handler := NewTestScenarioHandler(mockSvc)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios", handler.Create())

		reqBody := ``

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result response.ErrorResponse
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, pkg.InvalidReqBody, result.Error)
	})
	t.Run("failed case - test category not found", func(t *testing.T) {
		mockSvc := new(mocks.MockTestScenario)
		handler := NewTestScenarioHandler(mockSvc)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios", handler.Create())

		reqBody := `{
			"name": "load1",
			"test_category_id": 1,
			"mother_service_id": 1,
			"max_test_service_count": 1,
			"execution_duration": 1,
			"auto_step_change_rate": 1,
			"test_service_config": {
				"max_requests": 1,
				"max_duration": 1,
				"request_delay_duration": 1,
				"random_request_delay_min": null,
				"random_request_delay_max": null,
				"fixed_test_number": 1,
				"random_test_number_min": null,
				"random_test_number_max": null,
				"bad_value_rate": 0,
				"negative_value_rate": 0,
				"real_value_rate": 0,
				"zero_value_rate": 0,
				"string_value_rate": 0,
				"long_string_value_rate": 0,
				"null_value_rate": 0, 
				"database_name": "test_service_db",
        		"database_table_name": "test_service_table"
			}
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleInt64 := int64(1)
		dur := time.Millisecond * 1
		databaseName := "test_service_db"
		databaseTableName := "test_service_table"
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt64,
			ExecutionDuration:   &dur,
			AutoStepChangeRate:  &sampleInt64,
			TestServiceConfig: &entity.TestServiceConfig{
				MaxRequests:          sampleInt,
				MaxDuration:          dur,
				RequestDelayDuration: &sampleInt,
				FixedTestNumber:      &sampleInt,
				DatabaseName:         databaseName,
				DatabaseTableName:    databaseTableName,
			},
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest).Return(pkg.ErrTestCategoryNotFound)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result response.ErrorResponse
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.Equal(t, pkg.TestCategoryNotFound, result.Error)
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - failed to get test category", func(t *testing.T) {
		mockSvc := new(mocks.MockTestScenario)
		handler := NewTestScenarioHandler(mockSvc)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios", handler.Create())

		reqBody := `{
			"name": "load1",
			"test_category_id": 1,
			"mother_service_id": 1,
			"max_test_service_count": 1,
			"execution_duration": 1,
			"auto_step_change_rate": 1,
			"test_service_config": {
				"max_requests": 1,
				"max_duration": 1,
				"request_delay_duration": 1,
				"random_request_delay_min": null,
				"random_request_delay_max": null,
				"fixed_test_number": 1,
				"random_test_number_min": null,
				"random_test_number_max": null,
				"bad_value_rate": 0,
				"negative_value_rate": 0,
				"real_value_rate": 0,
				"zero_value_rate": 0,
				"string_value_rate": 0,
				"long_string_value_rate": 0,
				"null_value_rate": 0, 
				"database_name": "test_service_db",
        		"database_table_name": "test_service_table"
			}
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleInt64 := int64(1)
		dur := time.Millisecond * 1
		databaseName := "test_service_db"
		databaseTableName := "test_service_table"
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt64,
			ExecutionDuration:   &dur,
			AutoStepChangeRate:  &sampleInt64,
			TestServiceConfig: &entity.TestServiceConfig{
				MaxRequests:          sampleInt,
				MaxDuration:          dur,
				RequestDelayDuration: &sampleInt,
				FixedTestNumber:      &sampleInt,
				DatabaseName:         databaseName,
				DatabaseTableName:    databaseTableName,
			},
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest).Return(pkg.ErrFailedToGetTestCategory)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result response.ErrorResponse
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, result.Error, pkg.InternalServerErrorMessage)
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - validation error handler", func(t *testing.T) {
		mockSvc := new(mocks.MockTestScenario)
		handler := NewTestScenarioHandler(mockSvc)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios", handler.Create())

		reqBody := `{
			"name": "load1",
			"test_category_id": 1,
			"mother_service_id": 1,
			"max_test_service_count": -1,
			"execution_duration": 1,
			"auto_step_change_rate": 1,
			"test_service_config": {
				"max_requests": 1,
				"max_duration": 1,
				"request_delay_duration": 1,
				"random_request_delay_min": null,
				"random_request_delay_max": null,
				"fixed_test_number": 1,
				"random_test_number_min": null,
				"random_test_number_max": null,
				"bad_value_rate": 0,
				"negative_value_rate": 0,
				"real_value_rate": 0,
				"zero_value_rate": 0,
				"string_value_rate": 0,
				"long_string_value_rate": 0,
				"null_value_rate": 0, 
				"database_name": "test_service_db",
        		"database_table_name": "test_service_table"
			}
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleInt64 := int64(1)
		ncnt := int64(-1)
		dur := time.Millisecond * 1
		databaseName := "test_service_db"
		databaseTableName := "test_service_table"
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &ncnt,
			ExecutionDuration:   &dur,
			AutoStepChangeRate:  &sampleInt64,
			TestServiceConfig: &entity.TestServiceConfig{
				MaxRequests:          sampleInt,
				MaxDuration:          dur,
				RequestDelayDuration: &sampleInt,
				FixedTestNumber:      &sampleInt,
				DatabaseName:         databaseName,
				DatabaseTableName:    databaseTableName,
			},
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest).Return(pkg.ErrMaxTestServiceCountLessThanOne)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result response.ErrorResponse
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, result.Error, pkg.MaxTestServiceCountLessThanOne)
		mockSvc.AssertExpectations(t)
	})
}

func TestTestScenarioHandler_GetPaginated(t *testing.T) {
	t.Run("success case - with nil values", func(t *testing.T) {
		mockSvc := new(mocks.MockTestScenario)
		handler := NewTestScenarioHandler(mockSvc)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/search", handler.GetPaginated())

		payload := entity.TestScenarioPaginationRequest{
			Page:    1,
			PerPage: 2,
		}
		count := int64(2)
		someTime := time.Date(2026, 01, 12, 16, 36, 22, 0, time.UTC)

		reqBody := fmt.Sprintf(`{"page": %d,"per_page": %d}`, payload.Page, payload.PerPage)

		cnt := int64(100)
		rate := int64(50)
		dur := time.Millisecond * 500000
		serviceResult := []*entity.TestScenario{
			{
				ID:                  1,
				Name:                "load test 123",
				CreatedAt:           someTime,
				UpdatedAt:           someTime,
				StartedAt:           &someTime,
				Status:              entity.ScenarioStatusPending,
				MaxTestServiceCount: &cnt,
				AutoStepChangeRate:  &rate,
				ExecutionDuration:   &dur,
				TestCategoryID:      100,
				TestCategory: &entity.TestCategory{
					ID:                     100,
					Name:                   "load",
					Label:                  "Load Test",
					CreatedAt:              someTime,
					UpdatedAt:              someTime,
					HasMaxTestServiceCount: true,
					HasExecutionDuration:   true,
					HasAutoStepChangeRate:  false,
				},
				MotherServiceID: 200,
				MotherService: &entity.MotherService{
					ID:                     200,
					CreatedAt:              someTime,
					UpdatedAt:              someTime,
					Name:                   "mother200",
					ExceptionRate:          0,
					ResponseDelayRate:      0,
					ResponseDelayDuration:  nil,
					RandomResponseDelayMin: nil,
					RandomResponseDelayMax: nil,
					Status:                 entity.MotherServiceStatusRunning,
					ServiceDeploymentAddress: func() *string {
						addr := "m200.svc"
						return &addr
					}(),
					DatabaseName:      "m200",
					DatabaseTableName: "t200",
				},
			},
			{
				ID:           2,
				TestCategory: nil,
			},
		}

		mockSvc.On("GetPaginated", mock.Anything, payload).Return(serviceResult, count, nil)

		expected := response.PaginatedTestScenario{
			Page:    1,
			PerPage: 2,
			Data: []response.TestScenario{
				{
					ID:                  1,
					Name:                "load test 123",
					CreatedAt:           someTime,
					UpdatedAt:           someTime,
					Status:              entity.ScenarioStatusPending,
					MaxTestServiceCount: &cnt,
					AutoStepChangeRate:  &rate,
					ExecutionDuration:   &dur,
					StartedAt:           &someTime,
					TestCategory: &response.TestCategory{
						ID:                     100,
						Name:                   "load",
						Label:                  "Load Test",
						CreatedAt:              someTime,
						UpdatedAt:              someTime,
						HasMaxTestServiceCount: true,
						HasExecutionDuration:   true,
						HasAutoStepChangeRate:  false,
					},
					MotherService: &response.MotherService{
						ID:                     200,
						CreatedAt:              someTime,
						UpdatedAt:              someTime,
						Name:                   "mother200",
						ExceptionRate:          0,
						ResponseDelayRate:      0,
						ResponseDelayDuration:  nil,
						RandomResponseDelayMin: nil,
						RandomResponseDelayMax: nil,
						Status:                 entity.MotherServiceStatusRunning,
						ServiceDeploymentAddress: func() *string {
							addr := "m200.svc"
							return &addr
						}(),
						DatabaseName:      "m200",
						DatabaseTableName: "t200",
					},
				},
				{
					// this object is for testing nil checking on TestCategory and MotherService
					ID:            2,
					TestCategory:  nil,
					MotherService: nil,
				},
			},
			Total: count,
		}

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/search", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var got response.PaginatedTestScenario
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)
		json.Unmarshal(bts, &got)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expected.Page, got.Page)
		assert.Equal(t, got.Total, count)
		assert.Len(t, got.Data, 2)
		assert.Equal(t, expected, got)

		// for i, want := range expected.Data {
		// 	assert.Equal(t, want.ID, got.Data[i].ID)
		// 	assert.Equal(t, want.CreatedAt.Format("2006-01-02"), got.Data[i].CreatedAt.Format("2006-01-02"))
		// 	assert.Equal(t, want.UpdatedAt.Format("2006-01-02"), got.Data[i].UpdatedAt.Format("2006-01-02"))
		// 	assert.Equal(t, want.Status, got.Data[i].Status)

		// 	assert.Equal(t, want.MaxTestServiceCount, got.Data[i].MaxTestServiceCount)
		// 	assert.Equal(t, want.AutoStepChangeRate, got.Data[i].AutoStepChangeRate)
		// 	assert.Equal(t, want.ExecutionDuration, got.Data[i].ExecutionDuration)

		// 	if want.TestCategory != nil {
		// 		assert.Equal(t, want.TestCategory.ID, got.Data[i].TestCategory.ID)
		// 		assert.Equal(t, want.TestCategory.Name, got.Data[i].TestCategory.Name)
		// 		assert.Equal(t, want.TestCategory.Label, got.Data[i].TestCategory.Label)
		// 		assert.Equal(t, want.TestCategory.HasMaxTestServiceCount, got.Data[i].TestCategory.HasMaxTestServiceCount)
		// 		assert.Equal(t, want.TestCategory.HasExecutionDuration, got.Data[i].TestCategory.HasExecutionDuration)
		// 		assert.Equal(t, want.TestCategory.HasAutoStepChangeRate, got.Data[i].TestCategory.HasAutoStepChangeRate)
		// 		assert.Equal(t, want.TestCategory.CreatedAt.Format("2006-01-02"), got.Data[i].TestCategory.CreatedAt.Format("2006-01-02"))
		// 		assert.Equal(t, want.TestCategory.UpdatedAt.Format("2006-01-02"), got.Data[i].TestCategory.UpdatedAt.Format("2006-01-02"))
		// 	} else {
		// 		assert.Nil(t, got.Data[i].TestCategory)
		// 	}
		// }
	})
	t.Run("error: invalid request body", func(t *testing.T) {
		mockSvc := new(mocks.MockTestScenario)
		handler := NewTestScenarioHandler(mockSvc)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/search", handler.GetPaginated())

		reqBody := "invalid body"

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/search", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result response.ErrorResponse
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, result.Error, pkg.InvalidReqBody)
	})
	t.Run("error from service layer", func(t *testing.T) {
		mockSvc := new(mocks.MockTestScenario)
		handler := NewTestScenarioHandler(mockSvc)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/search", handler.GetPaginated())

		payload := entity.TestScenarioPaginationRequest{
			Page:    1,
			PerPage: 2,
		}
		reqBody := fmt.Sprintf(`{"page": %d,"per_page": %d}`, payload.Page, payload.PerPage)
		var serviceResult []*entity.TestScenario

		mockSvc.On("GetPaginated", mock.Anything, payload).Return(serviceResult, int64(0), errors.New("something went wrong"))

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/search", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var got response.PaginatedTestScenario
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)
		json.Unmarshal(bts, &got.Data)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestTestScenario_GetByID(t *testing.T) {
	_, err := config.LoadConfig()
	assert.Nil(t, err)

	t.Run("error: missing id in param", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		h := NewTestScenarioHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Get("/test-scenarios", h.GetByID())

		req := httptest.NewRequest(http.MethodGet, "/test-scenarios", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
	t.Run("error: invalid id data in param", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)

		h := NewTestScenarioHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Get("/test-scenarios/:id", h.GetByID())

		req := httptest.NewRequest(http.MethodGet, "/test-scenarios/invalid", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		var response response.ErrorResponse
		err = json.Unmarshal(bts, &response)
		assert.Nil(t, err)

		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
		assert.Equal(t, response.Error, pkg.InvalidIDInParams)
	})

	t.Run("error on getting data from service layer", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		var want *entity.TestScenario
		srv.On("GetByID", mock.Anything, uint64(1)).Return(want, errors.New("something went wrong"))
		h := NewTestScenarioHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Get("/test-scenarios/:id", h.GetByID())

		req := httptest.NewRequest(http.MethodGet, "/test-scenarios/1", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
	t.Run("error item not found", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		var want *entity.TestScenario
		srv.On("GetByID", mock.Anything, uint64(1)).Return(want, pkg.ErrTestScenarioNotFound)
		h := NewTestScenarioHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Get("/test-scenarios/:id", h.GetByID())

		req := httptest.NewRequest(http.MethodGet, "/test-scenarios/1", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		var response response.ErrorResponse
		err = json.Unmarshal(bts, &response)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.Equal(t, response.Error, pkg.TestScenarioNotFound)
	})
	t.Run("success case", func(t *testing.T) {
		someTime := time.Date(2026, 01, 13, 11, 00, 00, 0, time.UTC)
		cnt := int64(100)
		rate := int64(50)
		exe := time.Millisecond * 500000
		num := 10

		srv := new(mocks.MockTestScenario)
		item := &entity.TestScenario{
			ID:                  1,
			Name:                "load test 123",
			CreatedAt:           someTime,
			UpdatedAt:           someTime,
			StartedAt:           &someTime,
			Status:              entity.ScenarioStatusPending,
			MaxTestServiceCount: &cnt,
			AutoStepChangeRate:  &rate,
			ExecutionDuration:   &exe,
			TestCategoryID:      100,
			MotherServiceID:     200,
			TestCategory: &entity.TestCategory{
				ID:                     100,
				Name:                   "load",
				Label:                  "Load Test",
				CreatedAt:              someTime,
				UpdatedAt:              someTime,
				HasMaxTestServiceCount: true,
				HasExecutionDuration:   true,
				HasAutoStepChangeRate:  false,
			},
			MotherService: &entity.MotherService{
				ID:                     200,
				CreatedAt:              someTime,
				UpdatedAt:              someTime,
				Name:                   "mother200",
				ExceptionRate:          0,
				ResponseDelayRate:      0,
				ResponseDelayDuration:  nil,
				RandomResponseDelayMin: nil,
				RandomResponseDelayMax: nil,
				Status:                 entity.MotherServiceStatusRunning,
				ServiceDeploymentAddress: func() *string {
					addr := "m200.svc"
					return &addr
				}(),
				DatabaseName:      "m200",
				DatabaseTableName: "t200",
			},
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
		}
		srv.On("GetByID", mock.Anything, uint64(1)).Return(item, nil)
		h := NewTestScenarioHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Get("/test-scenarios/:id", h.GetByID())

		req := httptest.NewRequest(http.MethodGet, "/test-scenarios/1", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		bts, _ := io.ReadAll(resp.Body)

		expected := response.TestScenario{
			ID:                  1,
			Name:                "load test 123",
			CreatedAt:           someTime,
			UpdatedAt:           someTime,
			Status:              entity.ScenarioStatusPending,
			StartedAt:           &someTime,
			MaxTestServiceCount: &cnt,
			AutoStepChangeRate:  &rate,
			ExecutionDuration:   &exe,
			TestCategory: &response.TestCategory{
				ID:                     100,
				Name:                   "load",
				Label:                  "Load Test",
				CreatedAt:              someTime,
				UpdatedAt:              someTime,
				HasMaxTestServiceCount: true,
				HasExecutionDuration:   true,
				HasAutoStepChangeRate:  false,
			},
			MotherService: &response.MotherService{
				ID:                     200,
				CreatedAt:              someTime,
				UpdatedAt:              someTime,
				Name:                   "mother200",
				ExceptionRate:          0,
				ResponseDelayRate:      0,
				ResponseDelayDuration:  nil,
				RandomResponseDelayMin: nil,
				RandomResponseDelayMax: nil,
				Status:                 entity.MotherServiceStatusRunning,
				ServiceDeploymentAddress: func() *string {
					addr := "m200.svc"
					return &addr
				}(),
				DatabaseName:      "m200",
				DatabaseTableName: "t200",
			},
			TestServiceConfig: &response.TestServiceConfig{
				ID:                    100,
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
		}

		var got response.TestScenarioResponseByID
		err = json.Unmarshal(bts, &got)
		assert.NoError(t, err)
		assert.Equal(t, expected, got.Data)
	})
	t.Run("success case => category and mother service are null", func(t *testing.T) {
		someTime := time.Date(2026, 01, 13, 11, 00, 00, 0, time.UTC)
		cnt := int64(100)
		rate := int64(50)
		exe := time.Millisecond * 500000

		srv := new(mocks.MockTestScenario)
		item := &entity.TestScenario{
			ID:                  1,
			Name:                "load test 123",
			CreatedAt:           someTime,
			UpdatedAt:           someTime,
			Status:              entity.ScenarioStatusPending,
			MaxTestServiceCount: &cnt,
			AutoStepChangeRate:  &rate,
			ExecutionDuration:   &exe,
			TestCategoryID:      100,
			MotherServiceID:     200,
			TestCategory:        nil,
			MotherService:       nil,
			TestServiceConfig:   nil,
		}
		srv.On("GetByID", mock.Anything, uint64(1)).Return(item, nil)
		h := NewTestScenarioHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Get("/test-scenarios/:id", h.GetByID())

		req := httptest.NewRequest(http.MethodGet, "/test-scenarios/1", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		bts, _ := io.ReadAll(resp.Body)

		expected := response.TestScenario{
			ID:                  1,
			Name:                "load test 123",
			CreatedAt:           someTime,
			UpdatedAt:           someTime,
			Status:              entity.ScenarioStatusPending,
			MaxTestServiceCount: &cnt,
			AutoStepChangeRate:  &rate,
			ExecutionDuration:   &exe,
			TestCategory:        nil,
			MotherService:       nil,
			TestServiceConfig:   nil,
		}

		var got response.TestScenarioResponseByID
		err = json.Unmarshal(bts, &got)
		assert.NoError(t, err)
		assert.Equal(t, expected, got.Data)
	})
}

func TestTestScenario_Start(t *testing.T) {
	_, err := config.LoadConfig()
	assert.Nil(t, err)

	t.Run("error: missing id in param", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		h := NewTestScenarioHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/:id/start", h.Start())

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
	t.Run("error: invalid id data in param", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)

		h := NewTestScenarioHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/:id/start", h.Start())

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/invalid/start", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		var response response.ErrorResponse
		err = json.Unmarshal(bts, &response)
		assert.Nil(t, err)

		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
		assert.Equal(t, response.Error, pkg.InvalidIDInParams)
	})

	t.Run("error on getting data from service layer", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Start", mock.Anything, uint64(1)).Return(errors.New("something went wrong"))
		h := NewTestScenarioHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/:id/start", h.Start())

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/start", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
	t.Run("error item not found", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Start", mock.Anything, uint64(1)).Return(pkg.ErrTestScenarioNotFound)
		h := NewTestScenarioHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/:id/start", h.Start())

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/start", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		var response response.ErrorResponse
		err = json.Unmarshal(bts, &response)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.Equal(t, response.Error, pkg.TestScenarioNotFound)
	})

	t.Run("success case", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Start", mock.Anything, uint64(1)).Return(nil)
		h := NewTestScenarioHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/:id/start", h.Start())

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/start", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		var response response.SuccessResponse
		err = json.Unmarshal(bts, &response)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, response.Message, pkg.TestScenarioStarted)
	})

}
