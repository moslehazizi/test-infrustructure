package handler

import (
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/usecase/mocks"
	"control-panel-service/pkg"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
			"null_value_rate": 10
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		min := 10
		max := 20
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:           sampleInt,
			MaxDuration:           sampleInt,
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
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, pkg.CreateTestScenarioSuccessfully, result["message"])
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		min := 10
		max := 20
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:           sampleInt,
			MaxDuration:           sampleInt,
			RandomRequestDelayMin: &min,
			RandomRequestDelayMax: &max,
			RandomTestNumberMin:   &min,
			RandomTestNumberMax:   &max,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, pkg.CreateTestScenarioSuccessfully, result["message"])
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		min := 10
		max := 20
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:           sampleInt,
			MaxDuration:           sampleInt,
			RandomRequestDelayMin: &min,
			RandomRequestDelayMax: &max,
			FixedTestNumber:       &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, pkg.CreateTestScenarioSuccessfully, result["message"])
		mockSvc.AssertExpectations(t)
	})
	t.Run("success case", func(t *testing.T) {
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, pkg.CreateTestScenarioSuccessfully, result["message"])
		mockSvc.AssertExpectations(t)
	})
	t.Run("success case - fill null fields in test scenario", func(t *testing.T) {
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:            "load1",
			TestCategoryID:  sampleUin64,
			MotherServiceID: sampleUin64,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, pkg.CreateTestScenarioSuccessfully, result["message"])
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

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, pkg.InvalidReqBody, result["error"])
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(pkg.ErrTestCategoryNotFound)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.Equal(t, pkg.TestCategoryNotFound, result["error"])
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(pkg.ErrFailedToGetTestCategory)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, result["error"], pkg.InternalServerErrorMessage)
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - test scenario validation - max test service count less than one", func(t *testing.T) {
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(pkg.ErrMaxTestServiceCountLessThanOne)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, result["error"], pkg.MaxTestServiceCountLessThanOne)
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - test scenario validation - execution duration couldn't be less than 1", func(t *testing.T) {
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(pkg.ErrExecutionDurationLessThanOne)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, result["error"], pkg.ExecutionDurationLessThanOne)
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - test scenario validation - auto step change rate couldn't be less than one", func(t *testing.T) {
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(pkg.ErrAutoStepChangeRateLessThanOne)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, result["error"], pkg.AutoStepChangeRateLessThanOne)
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - test scenario validation - no need to set test service count", func(t *testing.T) {
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(pkg.ErrNoNeedMaxTestServiceCount)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, result["error"], pkg.NoNeedMaxTestServiceCount)
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - test scenario validation - execution duration not set", func(t *testing.T) {
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(pkg.ErrExecutionDurationNotSet)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, result["error"], pkg.ExecutionDurationNotSet)
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - test scenario validation - no need to set execution duration", func(t *testing.T) {
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(pkg.ErrNoNeedExecutionDuration)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, result["error"], pkg.NoNeedExecutionDuration)
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - test scenario validation - auto step change rate not set", func(t *testing.T) {
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(pkg.ErrAutoStepChangeNotSet)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, result["error"], pkg.AutoStepChangeNotSet)
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - test scenario validation - auto step change rate not set", func(t *testing.T) {
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(pkg.ErrNoNeedAutoStepChange)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, result["error"], pkg.NoNeedAutoStepChange)
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - test service config validation - invalid max request", func(t *testing.T) {
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(pkg.ErrInvalidMaxRequest)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, result["error"], pkg.InvalidMaxRequest)
		mockSvc.AssertExpectations(t)
	})

	t.Run("failed case - test service config validation - invalid max duration", func(t *testing.T) {
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(pkg.ErrInvalidMaxDuration)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, result["error"], pkg.InvalidMaxDuration)
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - test service config validation - invalid bad value rate", func(t *testing.T) {
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(pkg.ErrInvalidBadValueRate)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, result["error"], pkg.InvalidBadValueRate)
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - test service config validation - invalid negative value rate", func(t *testing.T) {
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(pkg.ErrInvalidNegativeValueRate)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, result["error"], pkg.InvalidNegativeValueRate)
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - test service config validation - invalid real value rate", func(t *testing.T) {
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(pkg.ErrInvalidRealValueRate)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, result["error"], pkg.InvalidRealValueRate)
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - test service config validation - invalid zero value rate", func(t *testing.T) {
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(pkg.ErrInvalidZeroValueRate)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, result["error"], pkg.InvalidZeroValueRate)
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - test service config validation - invalid string value rate", func(t *testing.T) {
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(pkg.ErrInvalidStringValueRate)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, result["error"], pkg.InvalidStringValueRate)
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - test service config validation - invalid long string value rate", func(t *testing.T) {
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(pkg.ErrInvalidLongStringValueRate)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, result["error"], pkg.InvalidLongStringValueRate)
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - test service config validation - invalid null value rate", func(t *testing.T) {
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(pkg.ErrInvalidNullValueRate)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, result["error"], pkg.InvalidNullValueRate)
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - test service config validation - invalid request delay duration config", func(t *testing.T) {
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(pkg.ErrInvalidRequestDelayDurationConfig)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, result["error"], pkg.InvalidRequestDelayDurationConfig)
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - test service config validation - invalid request delay duration value", func(t *testing.T) {
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(pkg.ErrInvalidRequestDelayDuration)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, result["error"], pkg.InvalidRequestDelayDuration)
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - test service config validation - min random request delay duration couldn't be more than max", func(t *testing.T) {
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(pkg.ErrMinDelayDurationMoreThanMax)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, result["error"], pkg.MinDelayDurationMoreThanMax)
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - test service config validation - invalid select test number config", func(t *testing.T) {
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(pkg.ErrInvalidTestNumberConfig)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, result["error"], pkg.InvalidTestNumberConfig)
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - test service config validation - invalid fixed test number", func(t *testing.T) {
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(pkg.ErrInvalidFixedTestNumber)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, result["error"], pkg.InvalidFixedTestNumber)
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - test service config validation - invalid select fixed or random test number config", func(t *testing.T) {
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(pkg.ErrInvalidFixedTestNumberConfig)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, result["error"], pkg.InvalidFixedTestNumberConfig)
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - test service config validation - invalid min or max in random test number", func(t *testing.T) {
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(pkg.ErrInvalidMinOrMaxRandomTestNumber)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, result["error"], pkg.InvalidMinOrMaxRandomTestNumber)
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - test service config validation - min test number couldn't be more than max", func(t *testing.T) {
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(pkg.ErrMinRandomTestNumberMoreThanMax)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, result["error"], pkg.MinRandomTestNumberMoreThanMax)
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - test service config validation - when bad value rate is zero then sum of bad values be zero", func(t *testing.T) {
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(pkg.ErrInvalidZeroSumOfBadValues)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, result["error"], pkg.InvalidZeroSumOfBadValues)
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - test service config validation - when bad value rate more than zero then sum of bad values be 100", func(t *testing.T) {
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
			"null_value_rate": 0
    	}`

		sampleUin64 := uint64(1)
		sampleInt := 1
		sampleTestScenarioRequest := &entity.TestScenario{
			Name:                "load1",
			TestCategoryID:      sampleUin64,
			MotherServiceID:     sampleUin64,
			MaxTestServiceCount: &sampleInt,
			ExecutionDuration:   &sampleInt,
			AutoStepChangeRate:  &sampleInt,
		}
		sampleTestServiceConfig := &entity.TestServiceConfig{
			MaxRequests:          sampleInt,
			MaxDuration:          sampleInt,
			RequestDelayDuration: &sampleInt,
			FixedTestNumber:      &sampleInt,
		}

		mockSvc.On("Create", mock.Anything, sampleTestScenarioRequest, sampleTestServiceConfig).Return(pkg.ErrInvalid100SumOfBadValues)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, result["error"], pkg.Invalid100SumOfBadValues)
		mockSvc.AssertExpectations(t)
	})

}
