package handler

import (
	"control-panel-service/config"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/server/dto/response"
	"control-panel-service/internal/usecase/mocks"
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
				ID:                      1,
				CreatedAt:               theTime,
				UpdatedAt:               theTime,
				Name:                    "load",
				Label:                   "Load Test",
				HasMaxTestServiceCount:  true,
				HasExecutionDuration:    true,
				HasAutoStepIncreaseRate: false,
			},
			{
				ID:                      2,
				CreatedAt:               theTime,
				UpdatedAt:               theTime,
				Name:                    "smoke",
				Label:                   "Smoke Test",
				HasMaxTestServiceCount:  true,
				HasExecutionDuration:    true,
				HasAutoStepIncreaseRate: false,
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
				ID:                      1,
				CreatedAt:               theTime,
				UpdatedAt:               theTime,
				Name:                    "load",
				Label:                   "Load Test",
				HasMaxTestServiceCount:  true,
				HasExecutionDuration:    true,
				HasAutoStepIncreaseRate: false,
			},
			{
				ID:                      2,
				CreatedAt:               theTime,
				UpdatedAt:               theTime,
				Name:                    "smoke",
				Label:                   "Smoke Test",
				HasMaxTestServiceCount:  true,
				HasExecutionDuration:    true,
				HasAutoStepIncreaseRate: false,
			},
		}

		var got []response.TestCategory
		err = json.Unmarshal(bts, &got)
		assert.NoError(t, err)
		assert.True(t, reflect.DeepEqual(got, want))
	})
}
