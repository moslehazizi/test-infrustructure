package handler

import (
	"control-panel-service/config"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/server/dto/response"
	"control-panel-service/internal/usecase/mocks"
	"control-panel-service/pkg"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewTestCategoryHandler(t *testing.T) {
	cfg, err := config.LoadConfig()
	assert.Nil(t, err)

	t.Run("ok", func(t *testing.T) {
		h := NewTestCategoryHandler(&cfg, &mocks.MockTestCategoryService{})
		require.NotNil(t, h.testCategoryService)
		require.NotNil(t, h.cfg)
	})
}

func TestTestCategory_GetAll(t *testing.T) {
	cfg, err := config.LoadConfig()
	assert.Nil(t, err)

	t.Run("error_on_getting_data_from_service_layer", func(t *testing.T) {
		srv := new(mocks.MockTestCategoryService)
		var items []entity.TestCategory
		srv.On("GetAll", mock.Anything).Return(items, errors.New("something went wrong"))
		h := NewTestCategoryHandler(&cfg, srv)

		r := chi.NewRouter()
		r.Get("/test-categories", h.GetAll)

		req := httptest.NewRequest(http.MethodGet, "/test-categories", nil)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("success_case", func(t *testing.T) {
		var theTime time.Time // use nil time to avoid reflect Deep equal issue while having a json decoding
		srv := new(mocks.MockTestCategoryService)
		items := []entity.TestCategory{
			{
				ID:                     1,
				CreatedAt:              theTime,
				UpdatedAt:              theTime,
				Name:                   "load",
				Label:                  "Load Test",
				HasMaxTestServiceCount: true,
				HasNumSteps:            true,
				Active:                 true,
			},
			{
				ID:                     2,
				CreatedAt:              theTime,
				UpdatedAt:              theTime,
				Name:                   "smoke",
				Label:                  "Smoke Test",
				HasMaxTestServiceCount: true,
				HasNumSteps:            true,
				Active:                 true,
			},
		}
		srv.On("GetAll", mock.Anything).Return(items, nil)
		h := NewTestCategoryHandler(&cfg, srv)

		r := chi.NewRouter()
		r.Get("/test-categories", h.GetAll)

		req := httptest.NewRequest(http.MethodGet, "/test-categories", nil)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		want := []response.TestCategory{
			{
				ID:                     1,
				CreatedAt:              theTime,
				UpdatedAt:              theTime,
				Name:                   "load",
				Label:                  "Load Test",
				HasMaxTestServiceCount: true,
				HasNumSteps:            true,
				Active:                 true,
			},
			{
				ID:                     2,
				CreatedAt:              theTime,
				UpdatedAt:              theTime,
				Name:                   "smoke",
				Label:                  "Smoke Test",
				HasMaxTestServiceCount: true,
				HasNumSteps:            true,
				Active:                 true,
			},
		}

		var got []response.TestCategory
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		assert.NoError(t, err)
		assert.True(t, reflect.DeepEqual(got, want))
	})
}

func TestTestCategory_GetByID(t *testing.T) {
	cfg, err := config.LoadConfig()
	assert.Nil(t, err)

	t.Run("error_invalid_id_data_in_param", func(t *testing.T) {
		srv := new(mocks.MockTestCategoryService)
		h := NewTestCategoryHandler(&cfg, srv)

		r := chi.NewRouter()
		r.Get("/test-categories/{id}", h.GetByID)

		req := httptest.NewRequest(http.MethodGet, "/test-categories/invalid", nil)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var result response.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, pkg.InvalidIDInParams, result.Error)
	})

	t.Run("error_on_getting_data_from_service_layer", func(t *testing.T) {
		srv := new(mocks.MockTestCategoryService)
		var want *entity.TestCategory
		srv.On("GetByID", mock.Anything, uint64(1)).Return(want, errors.New("something went wrong"))
		h := NewTestCategoryHandler(&cfg, srv)

		r := chi.NewRouter()
		r.Get("/test-categories/{id}", h.GetByID)

		req := httptest.NewRequest(http.MethodGet, "/test-categories/1", nil)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("error_item_not_found", func(t *testing.T) {
		srv := new(mocks.MockTestCategoryService)
		var want *entity.TestCategory
		srv.On("GetByID", mock.Anything, uint64(1)).Return(want, pkg.ErrTestCategoryNotFound)
		h := NewTestCategoryHandler(&cfg, srv)

		r := chi.NewRouter()
		r.Get("/test-categories/{id}", h.GetByID)

		req := httptest.NewRequest(http.MethodGet, "/test-categories/1", nil)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var result response.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, pkg.TestCategoryNotFound, result.Error)
	})

	t.Run("success_case", func(t *testing.T) {
		theTime := time.Date(2026, 01, 13, 10, 06, 30, 0, time.UTC)
		srv := new(mocks.MockTestCategoryService)
		item := &entity.TestCategory{
			ID:                     1,
			CreatedAt:              theTime,
			UpdatedAt:              theTime,
			Name:                   "load",
			Label:                  "Load Test",
			HasMaxTestServiceCount: true,
			HasNumSteps:            true,
			Active:                 true,
		}
		srv.On("GetByID", mock.Anything, uint64(1)).Return(item, nil)
		h := NewTestCategoryHandler(&cfg, srv)

		r := chi.NewRouter()
		r.Get("/test-categories/{id}", h.GetByID)

		req := httptest.NewRequest(http.MethodGet, "/test-categories/1", nil)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		want := response.TestCategory{
			ID:                     1,
			CreatedAt:              theTime,
			UpdatedAt:              theTime,
			Name:                   "load",
			Label:                  "Load Test",
			HasMaxTestServiceCount: true,
			HasNumSteps:            true,
			Active:                 true,
		}
		var got response.TestCategory
		err := json.Unmarshal(rec.Body.Bytes(), &got)
		assert.NoError(t, err)

		assert.Equal(t, want, got)
	})
}
