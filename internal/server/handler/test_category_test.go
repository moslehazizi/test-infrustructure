package handler

import (
	"control-panel-service/config"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/server/dto/response"
	"control-panel-service/internal/usecase/mocks"
	"control-panel-service/pkg"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
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

func TestGetAll(t *testing.T) {
	cfg, err := config.LoadConfig()
	assert.Nil(t, err)

	t.Run("error on getting data from service layer", func(t *testing.T) {
		srv := new(mocks.MockTestCategoryService)
		var items []entity.TestCategory
		srv.On("GetAll", mock.Anything).Return(items, errors.New("something went wrong"))
		h := NewTestCategoryHandler(&cfg, srv)

		app := fiber.New(fiber.Config{})
		app.Get("/test-categories", h.GetAll())

		req := httptest.NewRequest(http.MethodGet, "/test-categories", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
	t.Run("success case", func(t *testing.T) {
		var theTime time.Time // use nil time to avoid reflect Deep equal issue while having a json decoding
		srv := new(mocks.MockTestCategoryService)
		var items []entity.TestCategory = []entity.TestCategory{
			{
				ID:                     1,
				CreatedAt:              theTime,
				UpdatedAt:              theTime,
				Name:                   "load",
				Label:                  "Load Test",
				HasMaxTestServiceCount: true,
				HasExecutionDuration:   true,
				HasAutoStepChangeRate:  false,
			},
			{
				ID:                     2,
				CreatedAt:              theTime,
				UpdatedAt:              theTime,
				Name:                   "smoke",
				Label:                  "Smoke Test",
				HasMaxTestServiceCount: true,
				HasExecutionDuration:   true,
				HasAutoStepChangeRate:  false,
			},
		}
		srv.On("GetAll", mock.Anything).Return(items, nil)
		h := NewTestCategoryHandler(&cfg, srv)

		app := fiber.New(fiber.Config{})
		app.Get("/test-categories", h.GetAll())

		req := httptest.NewRequest(http.MethodGet, "/test-categories", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		bts, _ := io.ReadAll(resp.Body)

		want := []response.TestCategory{
			{
				ID:                     1,
				CreatedAt:              theTime,
				UpdatedAt:              theTime,
				Name:                   "load",
				Label:                  "Load Test",
				HasMaxTestServiceCount: true,
				HasExecutionDuration:   true,
				HasAutoStepChangeRate:  false,
			},
			{
				ID:                     2,
				CreatedAt:              theTime,
				UpdatedAt:              theTime,
				Name:                   "smoke",
				Label:                  "Smoke Test",
				HasMaxTestServiceCount: true,
				HasExecutionDuration:   true,
				HasAutoStepChangeRate:  false,
			},
		}

		var got []response.TestCategory
		err = json.Unmarshal(bts, &got)
		assert.NoError(t, err)
		assert.True(t, reflect.DeepEqual(got, want))
	})
}

func TestGetByID(t *testing.T) {
	cfg, err := config.LoadConfig()
	assert.Nil(t, err)

	t.Run("error: missing id in param", func(t *testing.T) {
		srv := new(mocks.MockTestCategoryService)
		var want *entity.TestCategory
		srv.On("GetByID", mock.Anything, uint64(1)).Return(want, errors.New("something went wrong"))
		h := NewTestCategoryHandler(&cfg, srv)

		app := fiber.New(fiber.Config{})
		app.Get("/test-categories", h.GetByID())

		req := httptest.NewRequest(http.MethodGet, "/test-categories", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
	t.Run("error: invalid id data in param", func(t *testing.T) {
		srv := new(mocks.MockTestCategoryService)
		var want *entity.TestCategory
		srv.On("GetByID", mock.Anything, uint64(1)).Return(want, errors.New("something went wrong"))
		h := NewTestCategoryHandler(&cfg, srv)

		app := fiber.New(fiber.Config{})
		app.Get("/test-categories/:id", h.GetByID())

		req := httptest.NewRequest(http.MethodGet, "/test-categories/invalid", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		var response struct {
			Error string `json:"error"`
		}
		err = json.Unmarshal(bts, &response)
		assert.Nil(t, err)

		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
		assert.Equal(t, response.Error, pkg.InvalidIDInParams)
	})

	t.Run("error on getting data from service layer", func(t *testing.T) {
		srv := new(mocks.MockTestCategoryService)
		var want *entity.TestCategory
		srv.On("GetByID", mock.Anything, uint64(1)).Return(want, errors.New("something went wrong"))
		h := NewTestCategoryHandler(&cfg, srv)

		app := fiber.New(fiber.Config{})
		app.Get("/test-categories/:id", h.GetByID())

		req := httptest.NewRequest(http.MethodGet, "/test-categories/1", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
	t.Run("error item not found", func(t *testing.T) {
		srv := new(mocks.MockTestCategoryService)
		var want *entity.TestCategory
		srv.On("GetByID", mock.Anything, uint64(1)).Return(want, pkg.ErrTestCategoryNotFound)
		h := NewTestCategoryHandler(&cfg, srv)

		app := fiber.New(fiber.Config{})
		app.Get("/test-categories/:id", h.GetByID())

		req := httptest.NewRequest(http.MethodGet, "/test-categories/1", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		var response struct {
			Error string `json:"error"`
		}
		err = json.Unmarshal(bts, &response)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.Equal(t, response.Error, pkg.TestCategoryNotFound)
	})
	t.Run("success case", func(t *testing.T) {
		theTime := time.Now()
		srv := new(mocks.MockTestCategoryService)
		want := &entity.TestCategory{
			ID:                      1,
			CreatedAt:               theTime,
			UpdatedAt:               theTime,
			Name:                    "load",
			Label:                   "Load Test",
			HasMaxTestServiceCount:  true,
			HasExecutionDuration:    true,
			HasAutoStepIncreaseRate: false,
		}
		srv.On("GetByID", mock.Anything, uint64(1)).Return(want, nil)
		h := NewTestCategoryHandler(&cfg, srv)

		app := fiber.New(fiber.Config{})
		app.Get("/test-categories/:id", h.GetByID())

		req := httptest.NewRequest(http.MethodGet, "/test-categories/1", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		bts, _ := io.ReadAll(resp.Body)

		var got response.TestCategory
		err = json.Unmarshal(bts, &got)
		assert.NoError(t, err)
		assert.Equal(t, want.ID, got.ID)
		assert.Equal(t, want.Name, got.Name)
		assert.Equal(t, want.Label, got.Label)
		assert.Equal(t, want.CreatedAt.Format("2006-01-02"), got.CreatedAt.Format("2006-01-02"))
		assert.Equal(t, want.UpdatedAt.Format("2006-01-02"), got.UpdatedAt.Format("2006-01-02"))
		assert.Equal(t, want.HasAutoStepIncreaseRate, got.HasAutoStepIncreaseRate)
		assert.Equal(t, want.HasExecutionDuration, got.HasExecutionDuration)
		assert.Equal(t, want.HasMaxTestServiceCount, got.HasMaxTestServiceCount)
	})
}
