package handler

import (
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/server/dto/response"
	"control-panel-service/internal/usecase/mocks"
	"control-panel-service/pkg"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMotherServiceHandler_New(t *testing.T) {
	mockSrv := new(mocks.MockMotherService)

	handler := NewMotherServiceHandler(mockSrv)

	assert.NotNil(t, handler)
}

func TestMotherServiceHandler_Create(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)

		handler := NewMotherServiceHandler(mockSvc)

		r := chi.NewRouter()
		r.Post("/mother-services", handler.Create)

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

		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		var result response.SuccessResponse
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, pkg.CreateMotherServiceSuccessfully, result.Message)

		mockSvc.AssertExpectations(t)
	})

	t.Run("failed_case_invalid_request", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)

		handler := NewMotherServiceHandler(mockSvc)

		r := chi.NewRouter()
		r.Post("/mother-services", handler.Create)

		reqBody := `sample`

		req := httptest.NewRequest(http.MethodPost, "/mother-services", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		var result response.ErrorResponse

		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, pkg.InvalidReqBody, result.Error)
	})

	t.Run("failed_case_required_fields_in_request_body", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)

		handler := NewMotherServiceHandler(mockSvc)

		r := chi.NewRouter()
		r.Post("/mother-services", handler.Create)

		reqBody := `{
	    "response_delay_rate": 10,
	    "database_name": "service_db",
	    "database_table_name": "data",
	}`
		req := httptest.NewRequest(http.MethodPost, "/mother-services", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		var result response.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, pkg.InvalidReqBody, result.Error)
	})

	t.Run("success_case_with_nullable_values", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)

		handler := NewMotherServiceHandler(mockSvc)

		r := chi.NewRouter()
		r.Post("/mother-services", handler.Create)

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

		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		var result response.SuccessResponse
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, pkg.CreateMotherServiceSuccessfully, result.Message)
		mockSvc.AssertExpectations(t)
	})

	t.Run("failed_case_duplication", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)

		handler := NewMotherServiceHandler(mockSvc)

		r := chi.NewRouter()
		r.Post("/mother-services", handler.Create)

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

		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		var result response.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusConflict, rec.Code)
		assert.Equal(t, pkg.MotherServiceAlreadyExist, result.Error)
		mockSvc.AssertExpectations(t)
	})

	t.Run("failed_case_internal_error", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)

		handler := NewMotherServiceHandler(mockSvc)

		r := chi.NewRouter()
		r.Post("/mother-services", handler.Create)

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

		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		var result response.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, pkg.InternalServerErrorMessage, result.Error)
		mockSvc.AssertExpectations(t)
	})

	t.Run("failed_case_request_validation_error_delay_rate_is_negative", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)

		handler := NewMotherServiceHandler(mockSvc)

		r := chi.NewRouter()
		r.Post("/mother-services", handler.Create)

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

		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		var result response.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		assert.Equal(t, pkg.InvalidResponseDelayRate, result.Error)
		mockSvc.AssertExpectations(t)
	})

	t.Run("failed_case_request_validation_error_exception_rate_is_negative", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)

		handler := NewMotherServiceHandler(mockSvc)

		r := chi.NewRouter()
		r.Post("/mother-services", handler.Create)

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

		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		var result response.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		assert.Equal(t, pkg.InvalidExceptionRate, result.Error)
		mockSvc.AssertExpectations(t)
	})

	t.Run("failed_case_request_validation_error_fixed_delay_is_set_but_rate_is_0", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)

		handler := NewMotherServiceHandler(mockSvc)

		r := chi.NewRouter()
		r.Post("/mother-services", handler.Create)

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

		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		var result response.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		assert.Equal(t, pkg.InvalidDelayConfiguration, result.Error)
		mockSvc.AssertExpectations(t)
	})

	t.Run("failed_case_request_validation_error_min_is_greater_than_max", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)

		handler := NewMotherServiceHandler(mockSvc)

		r := chi.NewRouter()
		r.Post("/mother-services", handler.Create)

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

		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		var result response.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		assert.Equal(t, pkg.InvalidRandomDelayRange, result.Error)
		mockSvc.AssertExpectations(t)
	})
}

func TestMotherServiceHandler_GetByID(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)

		handler := NewMotherServiceHandler(mockSvc)

		expectedSvcResp := &entity.MotherService{
			Name:              "mother1",
			Status:            entity.MotherServiceStatusRunning,
			DatabaseName:      "db1",
			DatabaseTableName: "factorial",
		}

		r := chi.NewRouter()
		r.Get("/mother-services/{id}", handler.GetByID)

		mockSvc.On("GetByID", mock.Anything, uint64(1)).Return(expectedSvcResp, nil)

		req := httptest.NewRequest(http.MethodGet, "/mother-services/1", nil)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var response response.MotherServiceResponseByID

		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NotNil(t, response)
		assert.Equal(t, expectedSvcResp.Name, response.Data.Name)
		assert.Equal(t, string(expectedSvcResp.Status), string(response.Data.Status))
		assert.Equal(t, expectedSvcResp.DatabaseName, response.Data.DatabaseName)
		assert.Equal(t, expectedSvcResp.DatabaseTableName, response.Data.DatabaseTableName)
		assert.Nil(t, response.Data.ServiceDeploymentAddress)

		mockSvc.AssertExpectations(t)
	})

	t.Run("success_case_with_pointer_values", func(t *testing.T) {
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

		r := chi.NewRouter()
		r.Get("/mother-services/{id}", handler.GetByID)

		mockSvc.On("GetByID", mock.Anything, uint64(1)).Return(expectedSvcResp, nil)

		req := httptest.NewRequest(http.MethodGet, "/mother-services/1", nil)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var response response.MotherServiceResponseByID

		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, rec.Code, http.StatusOK)
		assert.NotNil(t, response)
		assert.Equal(t, expectedSvcResp.ServiceDeploymentAddress, response.Data.ServiceDeploymentAddress)

		mockSvc.AssertExpectations(t)
	})

	t.Run("failed_case_invalid_id", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)

		handler := NewMotherServiceHandler(mockSvc)

		r := chi.NewRouter()
		r.Get("/mother-services/{id}", handler.GetByID)

		req := httptest.NewRequest(http.MethodGet, "/mother-services/sd12", nil)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var response response.ErrorResponse

		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, rec.Code, http.StatusBadRequest)
		assert.Equal(t, response.Error, pkg.InvalidIDInParams)
	})

	t.Run("failed_case_not_found", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)

		handler := NewMotherServiceHandler(mockSvc)

		r := chi.NewRouter()
		r.Get("/mother-services/{id}", handler.GetByID)

		mockSvc.On("GetByID", mock.Anything, uint64(1)).Return(nil, pkg.ErrMotherServiceNotFound)

		req := httptest.NewRequest(http.MethodGet, "/mother-services/1", nil)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var response response.ErrorResponse

		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, rec.Code, http.StatusNotFound)
		assert.Equal(t, response.Error, pkg.MotherServiceNotFound)

		mockSvc.AssertExpectations(t)
	})

	t.Run("failed_case_internal_server_error", func(t *testing.T) {
		mockSvc := new(mocks.MockMotherService)

		handler := NewMotherServiceHandler(mockSvc)

		r := chi.NewRouter()
		r.Get("/mother-services/{id}", handler.GetByID)

		mockSvc.On("GetByID", mock.Anything, uint64(1)).Return(nil, errors.New("error happened"))

		req := httptest.NewRequest(http.MethodGet, "/mother-services/1", nil)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var response response.ErrorResponse

		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, rec.Code, http.StatusInternalServerError)
		assert.Equal(t, response.Error, pkg.InternalServerErrorMessage)

		mockSvc.AssertExpectations(t)
	})
}

// func TestMotherServiceHandler_GetPaginated(t *testing.T) {
// 	t.Run("success_case", func(t *testing.T) {
// 		mockSvc := new(mocks.MockMotherService)

// 		handler := NewMotherServiceHandler(mockSvc)

// 		app := fiber.New(fiber.Config{})
// 		app.Post("/mother-services/search", handler.GetPaginated())

// 		sampleString := "sample"
// 		sampleNum := 1
// 		sampleReq := request.PaginationRequest{
// 			Page:    1,
// 			PerPage: 2,
// 		}
// 		sampleSvcReq := entity.PaginationRequest{
// 			Page:    sampleReq.Page,
// 			PerPage: sampleReq.PerPage,
// 		}
// 		count := int64(2)
// 		reqBody := fmt.Sprintf(`{"page": %d,"per_page": %d}`, sampleReq.Page, sampleReq.PerPage)
// 		expectedMotherServices := []*entity.MotherService{
// 			{
// 				ID:                       uint64(5),
// 				Name:                     "mother1",
// 				Status:                   entity.MotherServiceStatusRunning,
// 				DatabaseName:             "db1",
// 				DatabaseTableName:        "factorial",
// 				ServiceDeploymentAddress: &sampleString,
// 				ResponseDelayDuration:    &sampleNum,
// 				RandomResponseDelayMin:   &sampleNum,
// 				RandomResponseDelayMax:   &sampleNum,
// 			},
// 			{
// 				ID:                       uint64(4),
// 				Name:                     "mother2",
// 				Status:                   entity.MotherServiceStatusRunning,
// 				DatabaseName:             "db1",
// 				DatabaseTableName:        "factorial",
// 				ServiceDeploymentAddress: &sampleString,
// 				ResponseDelayDuration:    &sampleNum,
// 				RandomResponseDelayMin:   &sampleNum,
// 				RandomResponseDelayMax:   &sampleNum,
// 			},
// 		}

// 		mockSvc.On("GetPaginated", mock.Anything, sampleSvcReq).Return(expectedMotherServices, count, nil)

// 		req := httptest.NewRequest(http.MethodPost, "/mother-services/search", strings.NewReader(reqBody))
// 		req.Header.Set("Content-Type", "application/json")

// 		resp, _ := app.Test(req)
// 		defer resp.Body.Close()

// 		var result response.PaginatedMotherServices
// 		bts, err := io.ReadAll(resp.Body)
// 		assert.Nil(t, err)

// 		err = json.Unmarshal(bts, &result)
// 		assert.Nil(t, err)

// 		assert.Equal(t, http.StatusOK, resp.StatusCode)
// 		assert.Equal(t, result.Total, count)
// 		assert.Equal(t, result.Data[0].Name, expectedMotherServices[0].Name)
// 		assert.Equal(t, result.Data[1].Name, expectedMotherServices[1].Name)
// 		mockSvc.AssertExpectations(t)
// 	})

// 	t.Run("success_case_with_nil_values", func(t *testing.T) {
// 		mockSvc := new(mocks.MockMotherService)

// 		handler := NewMotherServiceHandler(mockSvc)

// 		app := fiber.New(fiber.Config{})
// 		app.Post("/mother-services/search", handler.GetPaginated())

// 		sampleReq := request.PaginationRequest{
// 			Page:    1,
// 			PerPage: 2,
// 		}
// 		sampleSvcReq := entity.PaginationRequest{
// 			Page:    sampleReq.Page,
// 			PerPage: sampleReq.PerPage,
// 		}
// 		count := int64(2)
// 		reqBody := fmt.Sprintf(`{"page": %d,"per_page": %d}`, sampleReq.Page, sampleReq.PerPage)
// 		expectedMotherServices := []*entity.MotherService{
// 			{
// 				ID:                uint64(5),
// 				Name:              "mother1",
// 				Status:            entity.MotherServiceStatusRunning,
// 				DatabaseName:      "db1",
// 				DatabaseTableName: "factorial",
// 			},
// 			{
// 				ID:                uint64(4),
// 				Name:              "mother2",
// 				Status:            entity.MotherServiceStatusRunning,
// 				DatabaseName:      "db1",
// 				DatabaseTableName: "factorial",
// 			},
// 		}

// 		mockSvc.On("GetPaginated", mock.Anything, sampleSvcReq).Return(expectedMotherServices, count, nil)

// 		req := httptest.NewRequest(http.MethodPost, "/mother-services/search", strings.NewReader(reqBody))
// 		req.Header.Set("Content-Type", "application/json")

// 		resp, _ := app.Test(req)
// 		defer resp.Body.Close()

// 		var result response.PaginatedMotherServices
// 		bts, err := io.ReadAll(resp.Body)
// 		assert.Nil(t, err)

// 		err = json.Unmarshal(bts, &result)
// 		assert.Nil(t, err)

// 		assert.Equal(t, http.StatusOK, resp.StatusCode)
// 		assert.Equal(t, result.Data[0].Name, expectedMotherServices[0].Name)
// 		assert.Equal(t, result.Total, count)
// 		assert.Equal(t, result.Data[1].Name, expectedMotherServices[1].Name)
// 		mockSvc.AssertExpectations(t)
// 	})

// 	t.Run("failed_case_invalid_request", func(t *testing.T) {
// 		mockSvc := new(mocks.MockMotherService)

// 		handler := NewMotherServiceHandler(mockSvc)

// 		app := fiber.New(fiber.Config{})
// 		app.Post("/mother-services/search", handler.GetPaginated())

// 		reqBody := `{"page": ewy,"per_page": erw}`

// 		req := httptest.NewRequest(http.MethodPost, "/mother-services/search", strings.NewReader(reqBody))
// 		req.Header.Set("Content-Type", "application/json")

// 		resp, _ := app.Test(req)
// 		defer resp.Body.Close()

// 		var result response.ErrorResponse
// 		bts, err := io.ReadAll(resp.Body)
// 		assert.Nil(t, err)

// 		err = json.Unmarshal(bts, &result)
// 		assert.Nil(t, err)

// 		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
// 	})

// 	t.Run("failed_case_invalid_request_negative_page_or_per_page", func(t *testing.T) {
// 		mockSvc := new(mocks.MockMotherService)

// 		handler := NewMotherServiceHandler(mockSvc)

// 		app := fiber.New(fiber.Config{})
// 		app.Post("/mother-services/search", handler.GetPaginated())

// 		reqBody := `{"page": -1,"per_page": 3}`

// 		req := httptest.NewRequest(http.MethodPost, "/mother-services/search", strings.NewReader(reqBody))
// 		req.Header.Set("Content-Type", "application/json")

// 		resp, _ := app.Test(req)
// 		defer resp.Body.Close()

// 		var result response.ErrorResponse
// 		bts, err := io.ReadAll(resp.Body)
// 		assert.Nil(t, err)

// 		err = json.Unmarshal(bts, &result)
// 		assert.Nil(t, err)

// 		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
// 	})

// 	t.Run("failed_case_internal_server_error", func(t *testing.T) {
// 		mockSvc := new(mocks.MockMotherService)

// 		handler := NewMotherServiceHandler(mockSvc)

// 		app := fiber.New(fiber.Config{})
// 		app.Post("/mother-services/search", handler.GetPaginated())

// 		sampleReq := request.PaginationRequest{
// 			Page:    1,
// 			PerPage: 2,
// 		}
// 		sampleSvcReq := entity.PaginationRequest{
// 			Page:    sampleReq.Page,
// 			PerPage: sampleReq.PerPage,
// 		}
// 		reqBody := fmt.Sprintf(`{"page": %d,"per_page": %d}`, sampleReq.Page, sampleReq.PerPage)

// 		mockSvc.On("GetPaginated", mock.Anything, sampleSvcReq).Return(nil, int64(0), errors.New("error happened"))

// 		req := httptest.NewRequest(http.MethodPost, "/mother-services/search", strings.NewReader(reqBody))
// 		req.Header.Set("Content-Type", "application/json")

// 		resp, _ := app.Test(req)
// 		defer resp.Body.Close()

// 		var result response.ErrorResponse

// 		bts, err := io.ReadAll(resp.Body)
// 		assert.Nil(t, err)

// 		err = json.Unmarshal(bts, &result)
// 		assert.Nil(t, err)

// 		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
// 		mockSvc.AssertExpectations(t)
// 	})

// 	t.Run("success_case_empty_result", func(t *testing.T) {
// 		mockSvc := new(mocks.MockMotherService)

// 		handler := NewMotherServiceHandler(mockSvc)

// 		app := fiber.New(fiber.Config{})
// 		app.Post("/mother-services/search", handler.GetPaginated())

// 		sampleReq := request.PaginationRequest{
// 			Page:    0,
// 			PerPage: 0,
// 		}
// 		sampleSvcReq := entity.PaginationRequest{
// 			Page:    sampleReq.Page,
// 			PerPage: sampleReq.PerPage,
// 		}
// 		reqBody := fmt.Sprintf(`{"page": %d,"per_page": %d}`, sampleReq.Page, sampleReq.PerPage)

// 		mockSvc.On("GetPaginated", mock.Anything, sampleSvcReq).Return([]*entity.MotherService{}, int64(0), nil)

// 		req := httptest.NewRequest(http.MethodPost, "/mother-services/search", strings.NewReader(reqBody))
// 		req.Header.Set("Content-Type", "application/json")

// 		resp, _ := app.Test(req)
// 		defer resp.Body.Close()

// 		var result response.PaginatedMotherServices

// 		bts, err := io.ReadAll(resp.Body)
// 		assert.Nil(t, err)

// 		err = json.Unmarshal(bts, &result)
// 		assert.Nil(t, err)

// 		assert.Equal(t, http.StatusOK, resp.StatusCode)
// 		assert.Equal(t, result.Total, int64(0))
// 		assert.Equal(t, result.Data, []response.MotherService(nil))
// 		mockSvc.AssertExpectations(t)
// 	})
// }

// func TestMotherServiceHandler_Delete(t *testing.T) {
// 	_, err := config.LoadConfig()
// 	assert.Nil(t, err)

// 	t.Run("error_missing_id_in_param", func(t *testing.T) {
// 		mockSvc := new(mocks.MockMotherService)

// 		h := NewMotherServiceHandler(mockSvc)

// 		app := fiber.New(fiber.Config{})
// 		app.Post("/mother-services/:id/delete", h.Delete())

// 		req := httptest.NewRequest(http.MethodPost, "/mother-services/", nil)

// 		resp, _ := app.Test(req)
// 		defer resp.Body.Close()

// 		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
// 	})

// 	t.Run("error_invalid_id_data_in_param", func(t *testing.T) {
// 		mockSvc := new(mocks.MockMotherService)

// 		h := NewMotherServiceHandler(mockSvc)

// 		app := fiber.New(fiber.Config{})
// 		app.Post("/mother-services/:id/delete", h.Delete())

// 		req := httptest.NewRequest(http.MethodPost, "/mother-services/1sdf/delete", nil)

// 		resp, _ := app.Test(req)
// 		defer resp.Body.Close()

// 		bts, err := io.ReadAll(resp.Body)
// 		assert.Nil(t, err)

// 		var response response.ErrorResponse
// 		err = json.Unmarshal(bts, &response)
// 		assert.Nil(t, err)

// 		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
// 		assert.Equal(t, response.Error, pkg.InvalidIDInParams)
// 	})

// 	t.Run("error_on_getting_data_from_service_layer", func(t *testing.T) {
// 		mockSvc := new(mocks.MockMotherService)
// 		mockSvc.On("Delete", mock.Anything, uint64(1)).Return(errors.New("something went wrong"))
// 		h := NewMotherServiceHandler(mockSvc)

// 		app := fiber.New(fiber.Config{})
// 		app.Post("/mother-services/:id/delete", h.Delete())

// 		req := httptest.NewRequest(http.MethodPost, "/mother-services/1/delete", nil)

// 		resp, _ := app.Test(req)
// 		defer resp.Body.Close()

// 		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
// 	})

// 	t.Run("error_item_not_found", func(t *testing.T) {
// 		mockSvc := new(mocks.MockMotherService)
// 		mockSvc.On("Delete", mock.Anything, uint64(1)).Return(pkg.ErrTestScenarioNotFound)
// 		h := NewMotherServiceHandler(mockSvc)

// 		app := fiber.New(fiber.Config{})
// 		app.Post("/mother-services/:id/delete", h.Delete())

// 		req := httptest.NewRequest(http.MethodPost, "/mother-services/1/delete", nil)

// 		resp, _ := app.Test(req)
// 		defer resp.Body.Close()

// 		bts, err := io.ReadAll(resp.Body)
// 		assert.Nil(t, err)

// 		var response response.ErrorResponse
// 		err = json.Unmarshal(bts, &response)
// 		assert.Nil(t, err)

// 		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
// 		assert.Equal(t, response.Error, pkg.TestScenarioNotFound)
// 	})

// 	t.Run("success_case", func(t *testing.T) {
// 		mockSvc := new(mocks.MockMotherService)
// 		mockSvc.On("Delete", mock.Anything, uint64(1)).Return(nil)
// 		h := NewMotherServiceHandler(mockSvc)

// 		app := fiber.New(fiber.Config{})
// 		app.Post("/mother-services/:id/delete", h.Delete())

// 		req := httptest.NewRequest(http.MethodPost, "/mother-services/1/delete", nil)

// 		resp, _ := app.Test(req)
// 		defer resp.Body.Close()

// 		bts, err := io.ReadAll(resp.Body)
// 		assert.Nil(t, err)

// 		var response response.SuccessResponse
// 		err = json.Unmarshal(bts, &response)
// 		assert.Nil(t, err)

// 		assert.Equal(t, http.StatusOK, resp.StatusCode)
// 		assert.Equal(t, response.Message, pkg.MotherServiceDelete)
// 	})

// }
