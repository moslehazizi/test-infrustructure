package handler

import (
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/usecase/mocks"
	"control-panel-service/pkg"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMotherServiceHandler_New(t *testing.T) {
	mockSrv := new(mocks.MockMotherService)

	handler := NewMotherServiceHandler(mockSrv)

	assert.NotNil(t, handler)
}

func TestMotherServiceHandler_Create(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)
		handler := NewMotherServiceHandler(mockSvc)

		app := fiber.New(fiber.Config{})
		app.Post("/mother-service", handler.Create())

		reqBody := `{
        "name": "my-service",
        "exception_rate": 0.0,
        "response_delay_rate": 0.0,
        "provisioning_status": "pending",
        "database_name": "service_db",
        "database_table_name": "data",
        "kafka_livefeed_topic": "livefeed",
        "kafka_factorial_topic": "factorial"
    }`

		mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(svc *entity.MotherService) bool {
			return svc.Name == "my-service" &&
				svc.ExceptionRate == 0.0 &&
				svc.ResponseDelayRate == 0.0 &&
				svc.ProvisioningStatus == entity.ProvisioningStatusPending &&
				svc.DatabaseName == "service_db" &&
				svc.DatabaseTableName == "data" &&
				svc.KafkaLiveFeedTopic == "livefeed" &&
				svc.KafkaFactorialTopic == "factorial"
		})).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/mother-service", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, pkg.CreateMotherServiceSuccessfully, result["message"])
		mockSvc.AssertExpectations(t)
	})

	t.Run("failed case - invalid request", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)
		handler := NewMotherServiceHandler(mockSvc)

		app := fiber.New(fiber.Config{})
		app.Post("/mother-service", handler.Create())

		reqBody := `sample`

		req := httptest.NewRequest(http.MethodPost, "/mother-service", strings.NewReader(reqBody))
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

	t.Run("failed case - required fields in request body", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)
		handler := NewMotherServiceHandler(mockSvc)

		app := fiber.New(fiber.Config{})
		app.Post("/mother-service", handler.Create())

		reqBody := `{
        "response_delay_rate": 0.0,
        "provisioning_status": "pending",
        "database_name": "service_db",
        "database_table_name": "data",
    }`
		req := httptest.NewRequest(http.MethodPost, "/mother-service", strings.NewReader(reqBody))
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

	t.Run("success case - with nullable values", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)
		handler := NewMotherServiceHandler(mockSvc)

		app := fiber.New(fiber.Config{})
		app.Post("/mother-service", handler.Create())

		reqBody := `{
        "name": "my-service",
        "exception_rate": 0.0,
        "response_delay_rate": 0.0,
		"response_delay_duration": 1,
		"random_response_delay_min": 2,
		"random_response_delay_max": 3,
        "provisioning_status": "pending",
        "database_name": "service_db",
        "database_table_name": "data",
        "kafka_livefeed_topic": "livefeed",
        "kafka_factorial_topic": "factorial"
    }`

		sampleOne := 1
		sampleTwo := 2
		sampleThree := 3

		mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(svc *entity.MotherService) bool {
			return svc.Name == "my-service" &&
				svc.ExceptionRate == 0.0 &&
				svc.ResponseDelayRate == 0.0 &&
				svc.ResponseDelayDuration != nil && *svc.ResponseDelayDuration == sampleOne &&
				svc.RandomResponseDelayMin != nil && *svc.RandomResponseDelayMin == sampleTwo &&
				svc.RandomResponseDelayMax != nil && *svc.RandomResponseDelayMax == sampleThree &&
				svc.ProvisioningStatus == entity.ProvisioningStatusPending &&
				svc.DatabaseName == "service_db" &&
				svc.DatabaseTableName == "data" &&
				svc.KafkaLiveFeedTopic == "livefeed" &&
				svc.KafkaFactorialTopic == "factorial"
		})).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/mother-service", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, pkg.CreateMotherServiceSuccessfully, result["message"])
		mockSvc.AssertExpectations(t)
	})

	t.Run("failed case - duplication", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)
		handler := NewMotherServiceHandler(mockSvc)

		app := fiber.New(fiber.Config{})
		app.Post("/mother-service", handler.Create())

		reqBody := `{
        "name": "my-service",
        "exception_rate": 0.0,
        "response_delay_rate": 0.0,
        "provisioning_status": "pending",
        "database_name": "service_db",
        "database_table_name": "data",
        "kafka_livefeed_topic": "livefeed",
        "kafka_factorial_topic": "factorial"
    }`

		mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(svc *entity.MotherService) bool {
			return svc.Name == "my-service" &&
				svc.ExceptionRate == 0.0 &&
				svc.ResponseDelayRate == 0.0 &&
				svc.ProvisioningStatus == entity.ProvisioningStatusPending &&
				svc.DatabaseName == "service_db" &&
				svc.DatabaseTableName == "data" &&
				svc.KafkaLiveFeedTopic == "livefeed" &&
				svc.KafkaFactorialTopic == "factorial"
		})).Return(pkg.ErrMotherServiceAlreadyExist)

		req := httptest.NewRequest(http.MethodPost, "/mother-service", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
		assert.Equal(t, pkg.MotherServiceAlreadyExist, result["error"])
		mockSvc.AssertExpectations(t)
	})

	t.Run("failed case - internal error", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)
		handler := NewMotherServiceHandler(mockSvc)

		app := fiber.New(fiber.Config{})
		app.Post("/mother-service", handler.Create())

		reqBody := `{
        "name": "my-service",
        "exception_rate": 0.0,
        "response_delay_rate": 0.0,
        "provisioning_status": "pending",
        "database_name": "service_db",
        "database_table_name": "data",
        "kafka_livefeed_topic": "livefeed",
        "kafka_factorial_topic": "factorial"
    }`

		mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(svc *entity.MotherService) bool {
			return svc.Name == "my-service" &&
				svc.ExceptionRate == 0.0 &&
				svc.ResponseDelayRate == 0.0 &&
				svc.ProvisioningStatus == entity.ProvisioningStatusPending &&
				svc.DatabaseName == "service_db" &&
				svc.DatabaseTableName == "data" &&
				svc.KafkaLiveFeedTopic == "livefeed" &&
				svc.KafkaFactorialTopic == "factorial"
		})).Return(errors.New("internal error"))

		req := httptest.NewRequest(http.MethodPost, "/mother-service", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, pkg.InternalServerErrorMessage, result["error"])
		mockSvc.AssertExpectations(t)
	})
}
