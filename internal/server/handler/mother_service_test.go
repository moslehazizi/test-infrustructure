package handler

import (
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/server/dto/request"
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
		app.Post("/mother-services", handler.Create())

		reqBody := `{
        "name": "my-service",
        "exception_rate": 10,
        "response_delay_rate": 20,
        "database_name": "service_db",
        "database_table_name": "data"
    }`

		mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(svc *entity.MotherService) bool {
			return svc.Name == "my-service" &&
				svc.ExceptionRate == 10 &&
				svc.ResponseDelayRate == 20 &&
				svc.DatabaseName == "service_db" &&
				svc.DatabaseTableName == "data"
		})).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/mother-services", strings.NewReader(reqBody))
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
		app.Post("/mother-services", handler.Create())

		reqBody := `sample`

		req := httptest.NewRequest(http.MethodPost, "/mother-services", strings.NewReader(reqBody))
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
		app.Post("/mother-services", handler.Create())

		reqBody := `{
        "response_delay_rate": 10,
        "database_name": "service_db",
        "database_table_name": "data",
    }`
		req := httptest.NewRequest(http.MethodPost, "/mother-services", strings.NewReader(reqBody))
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
		app.Post("/mother-services", handler.Create())

		reqBody := `{
        "name": "my-service",
        "exception_rate": 10,
        "response_delay_rate": 20,
		"response_delay_duration": 1,
		"random_response_delay_min": 2,
		"random_response_delay_max": 3,
        "database_name": "service_db",
        "database_table_name": "data"
    }`

		sampleOne := 1
		sampleTwo := 2
		sampleThree := 3

		mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(svc *entity.MotherService) bool {
			return svc.Name == "my-service" &&
				svc.ExceptionRate == 10 &&
				svc.ResponseDelayRate == 20 &&
				svc.ResponseDelayDuration != nil && *svc.ResponseDelayDuration == sampleOne &&
				svc.RandomResponseDelayMin != nil && *svc.RandomResponseDelayMin == sampleTwo &&
				svc.RandomResponseDelayMax != nil && *svc.RandomResponseDelayMax == sampleThree &&
				svc.DatabaseName == "service_db" &&
				svc.DatabaseTableName == "data"
		})).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/mother-services", strings.NewReader(reqBody))
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
		app.Post("/mother-services", handler.Create())

		reqBody := `{
        "name": "my-service",
        "exception_rate": 10,
        "response_delay_rate": 20,
        "database_name": "service_db",
        "database_table_name": "data"
    }`

		mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(svc *entity.MotherService) bool {
			return svc.Name == "my-service" &&
				svc.ExceptionRate == 10 &&
				svc.ResponseDelayRate == 20 &&
				svc.DatabaseName == "service_db" &&
				svc.DatabaseTableName == "data"
		})).Return(pkg.ErrMotherServiceAlreadyExist)

		req := httptest.NewRequest(http.MethodPost, "/mother-services", strings.NewReader(reqBody))
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
		app.Post("/mother-services", handler.Create())

		reqBody := `{
        "name": "my-service",
        "exception_rate": 10,
        "response_delay_rate": 20,
        "database_name": "service_db",
        "database_table_name": "data"
    }`

		mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(svc *entity.MotherService) bool {
			return svc.Name == "my-service" &&
				svc.ExceptionRate == 10 &&
				svc.ResponseDelayRate == 20 &&
				svc.DatabaseName == "service_db" &&
				svc.DatabaseTableName == "data"
		})).Return(errors.New("internal error"))

		req := httptest.NewRequest(http.MethodPost, "/mother-services", strings.NewReader(reqBody))
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

	t.Run("failed case - request validation error delay rate is negative", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)
		handler := NewMotherServiceHandler(mockSvc)

		app := fiber.New(fiber.Config{})
		app.Post("/mother-services", handler.Create())

		reqBody := `{
        "name": "my-service",
        "exception_rate": 10,
        "response_delay_rate": 20,
        "database_name": "service_db",
        "database_table_name": "data"
    }`

		mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(svc *entity.MotherService) bool {
			return svc.Name == "my-service" &&
				svc.ExceptionRate == 10 &&
				svc.ResponseDelayRate == 20 &&
				svc.DatabaseName == "service_db" &&
				svc.DatabaseTableName == "data"
		})).Return(fmt.Errorf("failed to validate request: %w", pkg.ErrInvalidResponseDelayRate))

		req := httptest.NewRequest(http.MethodPost, "/mother-services", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, pkg.InvalidResponseDelayRate, result["error"])
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - request validation error exception rate is negative", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)
		handler := NewMotherServiceHandler(mockSvc)

		app := fiber.New(fiber.Config{})
		app.Post("/mother-services", handler.Create())

		reqBody := `{
        "name": "my-service",
        "exception_rate": 10,
        "response_delay_rate": 20,
        "database_name": "service_db",
        "database_table_name": "data"
    }`

		mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(svc *entity.MotherService) bool {
			return svc.Name == "my-service" &&
				svc.ExceptionRate == 10 &&
				svc.ResponseDelayRate == 20 &&
				svc.DatabaseName == "service_db" &&
				svc.DatabaseTableName == "data"
		})).Return(fmt.Errorf("failed to validate request: %w", pkg.ErrInvalidExceptionRate))

		req := httptest.NewRequest(http.MethodPost, "/mother-services", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, pkg.InvalidExceptionRate, result["error"])
		mockSvc.AssertExpectations(t)
	})

	t.Run("failed case - request validation error fixed delay is set but rate is 0", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)
		handler := NewMotherServiceHandler(mockSvc)

		app := fiber.New(fiber.Config{})
		app.Post("/mother-services", handler.Create())

		reqBody := `{
        "name": "my-service",
        "exception_rate": 10,
        "response_delay_rate": 20,
        "database_name": "service_db",
        "database_table_name": "data"
    }`

		mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(svc *entity.MotherService) bool {
			return svc.Name == "my-service" &&
				svc.ExceptionRate == 10 &&
				svc.ResponseDelayRate == 20 &&
				svc.DatabaseName == "service_db" &&
				svc.DatabaseTableName == "data"
		})).Return(fmt.Errorf("failed to validate request: %w", pkg.ErrInvalidDelayConfiguration))

		req := httptest.NewRequest(http.MethodPost, "/mother-services", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, pkg.InvalidDelayConfiguration, result["error"])
		mockSvc.AssertExpectations(t)
	})
	t.Run("failed case - request validation error min is greater than max", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)
		handler := NewMotherServiceHandler(mockSvc)

		app := fiber.New(fiber.Config{})
		app.Post("/mother-services", handler.Create())

		reqBody := `{
        "name": "my-service",
        "exception_rate": 10,
        "response_delay_rate": 20,
        "database_name": "service_db",
        "database_table_name": "data"
    }`

		mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(svc *entity.MotherService) bool {
			return svc.Name == "my-service" &&
				svc.ExceptionRate == 10 &&
				svc.ResponseDelayRate == 20 &&
				svc.DatabaseName == "service_db" &&
				svc.DatabaseTableName == "data"
		})).Return(fmt.Errorf("failed to validate request: %w", pkg.ErrInvalidRandomDelayRange))

		req := httptest.NewRequest(http.MethodPost, "/mother-services", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result map[string]interface{}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Equal(t, pkg.InvalidRandomDelayRange, result["error"])
		mockSvc.AssertExpectations(t)
	})
}

func TestMotherServiceHandler_GetByID(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)
		handler := NewMotherServiceHandler(mockSvc)

		expectedSvcResp := &entity.MotherService{
			Name:              "mother1",
			Status:            entity.MotherServiceStatusRunning,
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		app := fiber.New()
		app.Get("/mother-services/:id", handler.GetByID())

		mockSvc.On("GetByID", mock.Anything, uint64(1)).Return(expectedSvcResp, nil)

		req := httptest.NewRequest(http.MethodGet, "/mother-services/1", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.Nil(t, err)

		var response struct {
			Data response.MotherService `json:"data"`
		}

		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &response)
		assert.Nil(t, err)

		assert.Equal(t, resp.StatusCode, http.StatusOK)
		assert.NotNil(t, response)
		assert.Equal(t, expectedSvcResp.Name, response.Data.Name)
		assert.Equal(t, string(expectedSvcResp.Status), response.Data.Status)
		assert.Equal(t, expectedSvcResp.DatabaseName, response.Data.DatabaseName)
		assert.Equal(t, expectedSvcResp.DatabaseTableName, response.Data.DatabaseTableName)
		assert.Nil(t, response.Data.ServiceDeploymentAddress)

		mockSvc.AssertExpectations(t)
	})

	t.Run("success case - with pointer values", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)
		handler := NewMotherServiceHandler(mockSvc)
		sampleString := "service-address"
		sampleNum := 1

		expectedSvcResp := &entity.MotherService{
			Name:                     "mother1",
			Status:                   entity.MotherServiceStatusRunning,
			DatabaseName:             "db1",
			DatabaseTableName:        "factorial",
			ServiceDeploymentAddress: &sampleString,
			ResponseDelayDuration:    &sampleNum,
			RandomResponseDelayMin:   &sampleNum,
			RandomResponseDelayMax:   &sampleNum,
		}

		app := fiber.New()
		app.Get("/mother-services/:id", handler.GetByID())

		mockSvc.On("GetByID", mock.Anything, uint64(1)).Return(expectedSvcResp, nil)

		req := httptest.NewRequest(http.MethodGet, "/mother-services/1", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.Nil(t, err)

		var response struct {
			Data response.MotherService `json:"data"`
		}

		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &response)
		assert.Nil(t, err)

		assert.Equal(t, resp.StatusCode, http.StatusOK)
		assert.NotNil(t, response)
		assert.Equal(t, expectedSvcResp.ServiceDeploymentAddress, response.Data.ServiceDeploymentAddress)

		mockSvc.AssertExpectations(t)
	})

	t.Run("failed case - invalid id", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)
		handler := NewMotherServiceHandler(mockSvc)

		app := fiber.New()
		app.Get("/mother-services/:id", handler.GetByID())

		req := httptest.NewRequest(http.MethodGet, "/mother-services/sd12", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.Nil(t, err)

		var response struct {
			Error string `json:"error"`
		}

		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &response)
		assert.Nil(t, err)

		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
		assert.Equal(t, response.Error, pkg.InvalidIDInParams)
	})

	t.Run("failed case - not found", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)
		handler := NewMotherServiceHandler(mockSvc)

		app := fiber.New()
		app.Get("/mother-services/:id", handler.GetByID())

		mockSvc.On("GetByID", mock.Anything, uint64(1)).Return(nil, pkg.ErrMotherServiceNotFound)

		req := httptest.NewRequest(http.MethodGet, "/mother-services/1", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.Nil(t, err)

		var response struct {
			Error string `json:"error"`
		}

		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &response)
		assert.Nil(t, err)

		assert.Equal(t, resp.StatusCode, http.StatusNotFound)
		assert.Equal(t, response.Error, pkg.MotherServiceNotFound)

		mockSvc.AssertExpectations(t)
	})

	t.Run("failed case - internal server error", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)
		handler := NewMotherServiceHandler(mockSvc)

		app := fiber.New()
		app.Get("/mother-services/:id", handler.GetByID())

		mockSvc.On("GetByID", mock.Anything, uint64(1)).Return(nil, errors.New("error happened"))

		req := httptest.NewRequest(http.MethodGet, "/mother-services/1", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.Nil(t, err)

		var response struct {
			Error string `json:"error"`
		}

		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &response)
		assert.Nil(t, err)

		assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
		assert.Equal(t, response.Error, pkg.InternalServerErrorMessage)

		mockSvc.AssertExpectations(t)
	})
}

func TestMotherServiceHandler_GetPaginated(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)
		handler := NewMotherServiceHandler(mockSvc)

		app := fiber.New(fiber.Config{})
		app.Post("/mother-services/search", handler.GetPaginated())

		sampleString := "sample"
		sampleNum := 1
		sampleReq := request.PaginationRequest{
			Page:    1,
			PerPage: 2,
		}
		sampleSvcReq := entity.PaginationRequest{
			Page:    sampleReq.Page,
			PerPage: sampleReq.PerPage,
		}
		reqBody := fmt.Sprintf(`{"page": %d,"per_page": %d}`, sampleReq.Page, sampleReq.PerPage)
		expectedMotherServices := []*entity.MotherService{
			{
				ID:                       uint64(5),
				Name:                     "mother1",
				Status:                   entity.MotherServiceStatusRunning,
				DatabaseName:             "db1",
				DatabaseTableName:        "factorial",
				ServiceDeploymentAddress: &sampleString,
				ResponseDelayDuration:    &sampleNum,
				RandomResponseDelayMin:   &sampleNum,
				RandomResponseDelayMax:   &sampleNum,
			},
			{
				ID:                       uint64(4),
				Name:                     "mother2",
				Status:                   entity.MotherServiceStatusRunning,
				DatabaseName:             "db1",
				DatabaseTableName:        "factorial",
				ServiceDeploymentAddress: &sampleString,
				ResponseDelayDuration:    &sampleNum,
				RandomResponseDelayMin:   &sampleNum,
				RandomResponseDelayMax:   &sampleNum,
			},
		}

		mockSvc.On("GetPaginated", mock.Anything, sampleSvcReq).Return(expectedMotherServices, nil)

		req := httptest.NewRequest(http.MethodPost, "/mother-services/search", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result struct {
			Data    []response.MotherService `json:"data"`
			Page    int                      `json:"page"`
			PerPage int                      `json:"per_page"`
		}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, result.Data[0].Name, expectedMotherServices[0].Name)
		assert.Equal(t, result.Data[1].Name, expectedMotherServices[1].Name)
		mockSvc.AssertExpectations(t)
	})

	t.Run("success case - with nil values", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)
		handler := NewMotherServiceHandler(mockSvc)

		app := fiber.New(fiber.Config{})
		app.Post("/mother-services/search", handler.GetPaginated())

		sampleReq := request.PaginationRequest{
			Page:    1,
			PerPage: 2,
		}
		sampleSvcReq := entity.PaginationRequest{
			Page:    sampleReq.Page,
			PerPage: sampleReq.PerPage,
		}
		reqBody := fmt.Sprintf(`{"page": %d,"per_page": %d}`, sampleReq.Page, sampleReq.PerPage)
		expectedMotherServices := []*entity.MotherService{
			{
				ID:                uint64(5),
				Name:              "mother1",
				Status:            entity.MotherServiceStatusRunning,
				DatabaseName:      "db1",
				DatabaseTableName: "factorial",
			},
			{
				ID:                uint64(4),
				Name:              "mother2",
				Status:            entity.MotherServiceStatusRunning,
				DatabaseName:      "db1",
				DatabaseTableName: "factorial",
			},
		}

		mockSvc.On("GetPaginated", mock.Anything, sampleSvcReq).Return(expectedMotherServices, nil)

		req := httptest.NewRequest(http.MethodPost, "/mother-services/search", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result struct {
			Data    []response.MotherService `json:"data"`
			Page    int                      `json:"page"`
			PerPage int                      `json:"per_page"`
		}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, result.Data[0].Name, expectedMotherServices[0].Name)
		assert.Equal(t, result.Data[1].Name, expectedMotherServices[1].Name)
		mockSvc.AssertExpectations(t)
	})

	t.Run("failed case - invalid request", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)
		handler := NewMotherServiceHandler(mockSvc)

		app := fiber.New(fiber.Config{})
		app.Post("/mother-services/search", handler.GetPaginated())

		reqBody := `{"page": ewy,"per_page": erw}`

		req := httptest.NewRequest(http.MethodPost, "/mother-services/search", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result struct {
			Error string `json:"error"`
		}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("failed case - invalid request negative page or per page", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)
		handler := NewMotherServiceHandler(mockSvc)

		app := fiber.New(fiber.Config{})
		app.Post("/mother-services/search", handler.GetPaginated())

		reqBody := `{"page": -1,"per_page": 3}`

		req := httptest.NewRequest(http.MethodPost, "/mother-services/search", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result struct {
			Error string `json:"error"`
		}
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("failed case - internal server error", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)
		handler := NewMotherServiceHandler(mockSvc)

		app := fiber.New(fiber.Config{})
		app.Post("/mother-services/search", handler.GetPaginated())

		sampleReq := request.PaginationRequest{
			Page:    1,
			PerPage: 2,
		}
		sampleSvcReq := entity.PaginationRequest{
			Page:    sampleReq.Page,
			PerPage: sampleReq.PerPage,
		}
		reqBody := fmt.Sprintf(`{"page": %d,"per_page": %d}`, sampleReq.Page, sampleReq.PerPage)

		mockSvc.On("GetPaginated", mock.Anything, sampleSvcReq).Return(nil, errors.New("error happened"))

		req := httptest.NewRequest(http.MethodPost, "/mother-services/search", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result struct {
			Error string `json:"error"`
		}

		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockSvc.AssertExpectations(t)
	})

	t.Run("success case - empty result", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)
		handler := NewMotherServiceHandler(mockSvc)

		app := fiber.New(fiber.Config{})
		app.Post("/mother-services/search", handler.GetPaginated())

		sampleReq := request.PaginationRequest{
			Page:    0,
			PerPage: 0,
		}
		sampleSvcReq := entity.PaginationRequest{
			Page:    sampleReq.Page,
			PerPage: sampleReq.PerPage,
		}
		reqBody := fmt.Sprintf(`{"page": %d,"per_page": %d}`, sampleReq.Page, sampleReq.PerPage)

		mockSvc.On("GetPaginated", mock.Anything, sampleSvcReq).Return([]*entity.MotherService{}, nil)

		req := httptest.NewRequest(http.MethodPost, "/mother-services/search", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result struct {
			Data    []response.MotherService `json:"data"`
			Page    int                      `json:"page"`
			PerPage int                      `json:"per_page"`
		}

		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, result.Data, []response.MotherService(nil))
		mockSvc.AssertExpectations(t)
	})
}
