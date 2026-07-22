package handler

import (
	"control-panel-service/internal/server/dto/response"
	"control-panel-service/internal/usecase/mocks"
	"control-panel-service/pkg"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTestScenarioOperation_Start(t *testing.T) {
	t.Run("error_invalid_id_data_in_param", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		h := NewTestScenarioOperationHandler(srv)

		r := chi.NewRouter()
		r.Post("/test-scenarios/{id}/start", h.Start)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/invalid/start", nil)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var result response.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, pkg.InvalidIDInParams, result.Error)
	})

	t.Run("error_on_getting_data_from_service_layer", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Start", mock.Anything, uint64(1)).Return(errors.New("something went wrong"))
		h := NewTestScenarioOperationHandler(srv)

		r := chi.NewRouter()
		r.Post("/test-scenarios/{id}/start", h.Start)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/start", nil)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("error_item_not_found", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Start", mock.Anything, uint64(1)).Return(pkg.ErrTestScenarioNotFound)
		h := NewTestScenarioOperationHandler(srv)

		r := chi.NewRouter()
		r.Post("/test-scenarios/{id}/start", h.Start)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/start", nil)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var result response.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, pkg.TestScenarioNotFound, result.Error)
	})

	t.Run("success_case", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Start", mock.Anything, uint64(1)).Return(nil)
		h := NewTestScenarioOperationHandler(srv)

		r := chi.NewRouter()
		r.Post("/test-scenarios/{id}/start", h.Start)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/start", nil)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var result response.SuccessResponse
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, pkg.TestScenarioStarted, result.Message)
	})
}

func TestTestScenarioOperation_Pause(t *testing.T) {
	t.Run("error_invalid_id_data_in_param", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		h := NewTestScenarioOperationHandler(srv)

		r := chi.NewRouter()
		r.Post("/test-scenarios/{id}/pause", h.Pause)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/invalid/pause", nil)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var result response.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, pkg.InvalidIDInParams, result.Error)
	})

	t.Run("error_on_getting_data_from_service_layer", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Pause", mock.Anything, uint64(1)).Return(errors.New("something went wrong"))
		h := NewTestScenarioOperationHandler(srv)

		r := chi.NewRouter()
		r.Post("/test-scenarios/{id}/pause", h.Pause)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/pause", nil)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("error_item_not_found", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Pause", mock.Anything, uint64(1)).Return(pkg.ErrTestScenarioNotFound)
		h := NewTestScenarioOperationHandler(srv)

		r := chi.NewRouter()
		r.Post("/test-scenarios/{id}/pause", h.Pause)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/pause", nil)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var result response.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, pkg.TestScenarioNotFound, result.Error)
	})

	t.Run("success_case", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Pause", mock.Anything, uint64(1)).Return(nil)
		h := NewTestScenarioOperationHandler(srv)

		r := chi.NewRouter()
		r.Post("/test-scenarios/{id}/pause", h.Pause)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/pause", nil)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var result response.SuccessResponse
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, pkg.TestScenarioPaused, result.Message)
	})
}

func TestTestScenarioOperation_Resume(t *testing.T) {
	t.Run("error_invalid_id_data_in_param", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		h := NewTestScenarioOperationHandler(srv)

		r := chi.NewRouter()
		r.Post("/test-scenarios/{id}/resume", h.Resume)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/invalid/resume", nil)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var result response.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, pkg.InvalidIDInParams, result.Error)
	})

	t.Run("error_on_getting_data_from_service_layer", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Resume", mock.Anything, uint64(1)).Return(errors.New("something went wrong"))
		h := NewTestScenarioOperationHandler(srv)

		r := chi.NewRouter()
		r.Post("/test-scenarios/{id}/resume", h.Resume)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/resume", nil)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("error_item_not_found", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Resume", mock.Anything, uint64(1)).Return(pkg.ErrTestScenarioNotFound)
		h := NewTestScenarioOperationHandler(srv)

		r := chi.NewRouter()
		r.Post("/test-scenarios/{id}/resume", h.Resume)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/resume", nil)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var result response.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, pkg.TestScenarioNotFound, result.Error)
	})

	t.Run("success_case", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Resume", mock.Anything, uint64(1)).Return(nil)
		h := NewTestScenarioOperationHandler(srv)

		r := chi.NewRouter()
		r.Post("/test-scenarios/{id}/resume", h.Resume)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/resume", nil)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var result response.SuccessResponse
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, pkg.TestScenarioResumed, result.Message)
	})
}

func TestTestScenarioOperation_Stop(t *testing.T) {
	t.Run("error_invalid_id_data_in_param", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		h := NewTestScenarioOperationHandler(srv)

		r := chi.NewRouter()
		r.Post("/test-scenarios/{id}/stop", h.Stop)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/invalid/stop", nil)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var result response.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, pkg.InvalidIDInParams, result.Error)
	})

	t.Run("error_on_getting_data_from_service_layer", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Stop", mock.Anything, uint64(1)).Return(errors.New("something went wrong"))
		h := NewTestScenarioOperationHandler(srv)

		r := chi.NewRouter()
		r.Post("/test-scenarios/{id}/stop", h.Stop)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/stop", nil)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("error_item_not_found", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Stop", mock.Anything, uint64(1)).Return(pkg.ErrTestScenarioNotFound)
		h := NewTestScenarioOperationHandler(srv)

		r := chi.NewRouter()
		r.Post("/test-scenarios/{id}/stop", h.Stop)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/stop", nil)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var result response.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, pkg.TestScenarioNotFound, result.Error)
	})

	t.Run("success_case", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Stop", mock.Anything, uint64(1)).Return(nil)
		h := NewTestScenarioOperationHandler(srv)

		r := chi.NewRouter()
		r.Post("/test-scenarios/{id}/stop", h.Stop)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/stop", nil)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var result response.SuccessResponse
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, pkg.TestScenarioStop, result.Message)
	})
}

func TestTestScenarioOperation_Delete(t *testing.T) {
	t.Run("error_invalid_id_data_in_param", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		h := NewTestScenarioOperationHandler(srv)

		r := chi.NewRouter()
		r.Post("/test-scenarios/{id}/delete", h.Delete)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1sdf/delete", nil)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var result response.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, pkg.InvalidIDInParams, result.Error)
	})

	t.Run("error_on_getting_data_from_service_layer", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Delete", mock.Anything, uint64(1)).Return(errors.New("something went wrong"))
		h := NewTestScenarioOperationHandler(srv)

		r := chi.NewRouter()
		r.Post("/test-scenarios/{id}/delete", h.Delete)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/delete", nil)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("error_item_not_found", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Delete", mock.Anything, uint64(1)).Return(pkg.ErrTestScenarioNotFound)
		h := NewTestScenarioOperationHandler(srv)

		r := chi.NewRouter()
		r.Post("/test-scenarios/{id}/delete", h.Delete)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/delete", nil)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var result response.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, pkg.TestScenarioNotFound, result.Error)
	})

	t.Run("success_case", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Delete", mock.Anything, uint64(1)).Return(nil)
		h := NewTestScenarioOperationHandler(srv)

		r := chi.NewRouter()
		r.Post("/test-scenarios/{id}/delete", h.Delete)

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/delete", nil)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var result response.SuccessResponse
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, pkg.TestScenarioDelete, result.Message)
	})
}
