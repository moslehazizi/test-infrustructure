package handler

import (
	"control-panel-service/config"
	"control-panel-service/internal/server/dto/response"
	"control-panel-service/internal/usecase/mocks"
	"control-panel-service/pkg"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTestScenario_Start(t *testing.T) {
	_, err := config.LoadConfig()
	assert.Nil(t, err)

	t.Run("error_missing_id_in_param", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		h := NewTestScenarioOperationHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/:id/start", h.Start())

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
	t.Run("error_invalid_id_data_in_param", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)

		h := NewTestScenarioOperationHandler(srv)

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

	t.Run("error_on_getting_data_from_service_layer", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Start", mock.Anything, uint64(1)).Return(errors.New("something went wrong"))
		h := NewTestScenarioOperationHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/:id/start", h.Start())

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/start", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
	t.Run("error_item_not_found", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Start", mock.Anything, uint64(1)).Return(pkg.ErrTestScenarioNotFound)
		h := NewTestScenarioOperationHandler(srv)

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

	t.Run("success_case", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Start", mock.Anything, uint64(1)).Return(nil)
		h := NewTestScenarioOperationHandler(srv)

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

func TestTestScenario_Pause(t *testing.T) {
	_, err := config.LoadConfig()
	assert.Nil(t, err)

	t.Run("error_missing_id_in_param", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		h := NewTestScenarioOperationHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/:id/pause", h.Pause())

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("error_invalid_id_data_in_param", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)

		h := NewTestScenarioOperationHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/:id/pause", h.Pause())

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/invalid/pause", nil)

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

	t.Run("error_on_getting_data_from_service_layer", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Pause", mock.Anything, uint64(1)).Return(errors.New("something went wrong"))
		h := NewTestScenarioOperationHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/:id/pause", h.Pause())

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/pause", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("error_item_not_found", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Pause", mock.Anything, uint64(1)).Return(pkg.ErrTestScenarioNotFound)
		h := NewTestScenarioOperationHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/:id/pause", h.Pause())

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/pause", nil)

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

	t.Run("success_case", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Pause", mock.Anything, uint64(1)).Return(nil)
		h := NewTestScenarioOperationHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/:id/pause", h.Pause())

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/pause", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		var response response.SuccessResponse
		err = json.Unmarshal(bts, &response)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, response.Message, pkg.TestScenarioPaused)
	})

}

func TestTestScenario_Resume(t *testing.T) {
	_, err := config.LoadConfig()
	assert.Nil(t, err)

	t.Run("error_missing_id_in_param", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		h := NewTestScenarioOperationHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/:id/resume", h.Resume())

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("error_invalid_id_data_in_param", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)

		h := NewTestScenarioOperationHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/:id/resume", h.Resume())

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/invalid/resume", nil)

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

	t.Run("error_on_getting_data_from_service_layer", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Resume", mock.Anything, uint64(1)).Return(errors.New("something went wrong"))
		h := NewTestScenarioOperationHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/:id/resume", h.Resume())

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/resume", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("error_item_not_found", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Resume", mock.Anything, uint64(1)).Return(pkg.ErrTestScenarioNotFound)
		h := NewTestScenarioOperationHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/:id/resume", h.Resume())

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/resume", nil)

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

	t.Run("success_case", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Resume", mock.Anything, uint64(1)).Return(nil)
		h := NewTestScenarioOperationHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/:id/resume", h.Resume())

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/resume", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		var response response.SuccessResponse
		err = json.Unmarshal(bts, &response)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, response.Message, pkg.TestScenarioResumed)
	})

}

func TestTestScenario_Stop(t *testing.T) {
	_, err := config.LoadConfig()
	assert.Nil(t, err)

	t.Run("error_missing_id_in_param", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		h := NewTestScenarioOperationHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/:id/stop", h.Stop())

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("error_invalid_id_data_in_param", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)

		h := NewTestScenarioOperationHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/:id/stop", h.Stop())

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/invalid/stop", nil)

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

	t.Run("error_on_getting_data_from_service_layer", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Stop", mock.Anything, uint64(1)).Return(errors.New("something went wrong"))
		h := NewTestScenarioOperationHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/:id/stop", h.Stop())

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/stop", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("error_item_not_found", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Stop", mock.Anything, uint64(1)).Return(pkg.ErrTestScenarioNotFound)
		h := NewTestScenarioOperationHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/:id/stop", h.Stop())

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/stop", nil)

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

	t.Run("success_case", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Stop", mock.Anything, uint64(1)).Return(nil)
		h := NewTestScenarioOperationHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/:id/stop", h.Stop())

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/stop", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		var response response.SuccessResponse
		err = json.Unmarshal(bts, &response)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, response.Message, pkg.TestScenarioStop)
	})

}

func TestTestScenario_Delete(t *testing.T) {
	_, err := config.LoadConfig()
	assert.Nil(t, err)

	t.Run("error_missing_id_in_param", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		h := NewTestScenarioOperationHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/:id/delete", h.Delete())

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("error_invalid_id_data_in_param", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)

		h := NewTestScenarioOperationHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/:id/delete", h.Delete())

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1sdf/delete", nil)

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

	t.Run("error_on_getting_data_from_service_layer", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Delete", mock.Anything, uint64(1)).Return(errors.New("something went wrong"))
		h := NewTestScenarioOperationHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/:id/delete", h.Delete())

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/delete", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("error_item_not_found", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Delete", mock.Anything, uint64(1)).Return(pkg.ErrTestScenarioNotFound)
		h := NewTestScenarioOperationHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/:id/delete", h.Delete())

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/delete", nil)

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

	t.Run("success_case", func(t *testing.T) {
		srv := new(mocks.MockTestScenario)
		srv.On("Delete", mock.Anything, uint64(1)).Return(nil)
		h := NewTestScenarioOperationHandler(srv)

		app := fiber.New(fiber.Config{})
		app.Post("/test-scenarios/:id/delete", h.Delete())

		req := httptest.NewRequest(http.MethodPost, "/test-scenarios/1/delete", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		var response response.SuccessResponse
		err = json.Unmarshal(bts, &response)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, response.Message, pkg.TestScenarioDelete)
	})

}
