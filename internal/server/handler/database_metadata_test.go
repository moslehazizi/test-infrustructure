package handler

import (
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/server/dto/response"
	"control-panel-service/internal/usecase/mocks"
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
	"github.com/stretchr/testify/require"
)

func TestDatabaseMetadata_Initialization(t *testing.T) {
	mockSrv := new(mocks.MockDatabaseMetadataUsecase)

	handler := NewDatabaseMetadataHandler(mockSrv)

	assert.NotNil(t, handler)

	assert.NotNil(t, handler.databaseMetaDataService)
}

func TestStorageHandler_GetDatabases(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		mockUC := new(mocks.MockDatabaseMetadataUsecase)
		handler := NewDatabaseMetadataHandler(mockUC)

		app := fiber.New(fiber.Config{})
		app.Get("/databases", handler.GetAll())

		expectedDatabases := []string{"db1", "db2"}
		mockUC.On("GetAll", mock.Anything).Return(expectedDatabases, nil)

		req := httptest.NewRequest(http.MethodGet, "/databases", nil)
		resp, _ := app.Test(req)
		defer resp.Body.Close()

		var result response.DatabaseMetadataDatabasesResponse
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expectedDatabases, result.Data)

		mockUC.AssertExpectations(t)
	})

	t.Run("failure case - usecase error", func(t *testing.T) {
		mockUC := new(mocks.MockDatabaseMetadataUsecase)
		handler := NewDatabaseMetadataHandler(mockUC)

		app := fiber.New(fiber.Config{})
		app.Get("/databases", handler.GetAll())

		mockUC.On("GetAll", mock.Anything).Return([]string{}, errors.New("db error"))

		req := httptest.NewRequest(http.MethodGet, "/databases", nil)
		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockUC.AssertExpectations(t)
	})
}

func TestStorageHandler_GetTablesByDBNamePost(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		mockUC := new(mocks.MockDatabaseMetadataUsecase)
		handler := NewDatabaseMetadataHandler(mockUC)
		app := fiber.New()
		app.Post("/databases/tables", handler.GetTablesByDBNamePost())

		reqBody := `{"database_name": "testdb"}`
		expectedMotherTables := []string{"mother11", "mother2"}
		expectedTestTables := []string{"test11", "test2"}

		mockUC.On("GetTablesByDBName", mock.Anything, "testdb").Return(&entity.TablesByType{
			MotherTables: expectedMotherTables,
			TestTables:   expectedTestTables,
		}, nil)

		req := httptest.NewRequest(http.MethodPost, "/databases/tables", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		var result response.TablesByType
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expectedMotherTables, result.MotherTables)
		assert.Equal(t, expectedTestTables, result.TestTables)

		mockUC.AssertExpectations(t)
	})

	t.Run("success case - empty result", func(t *testing.T) {
		mockUC := new(mocks.MockDatabaseMetadataUsecase)
		handler := NewDatabaseMetadataHandler(mockUC)
		app := fiber.New()
		app.Post("/databases/tables", handler.GetTablesByDBNamePost())

		reqBody := `{"database_name": "testdb"}`
		expectedMotherTables := []string{}
		expectedTestTables := []string{}

		mockUC.On("GetTablesByDBName", mock.Anything, "testdb").Return(nil, nil)

		req := httptest.NewRequest(http.MethodPost, "/databases/tables", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		var result response.TablesByType
		bts, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bts, &result)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expectedMotherTables, result.MotherTables)
		assert.Equal(t, expectedTestTables, result.TestTables)

		mockUC.AssertExpectations(t)
	})

	t.Run("failure case - bad request body", func(t *testing.T) {
		mockUC := new(mocks.MockDatabaseMetadataUsecase)
		handler := NewDatabaseMetadataHandler(mockUC)
		app := fiber.New()
		app.Post("/databases/tables", handler.GetTablesByDBNamePost())

		req := httptest.NewRequest(http.MethodPost, "/databases/tables", strings.NewReader("invalid json"))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("failure case - missing database_name", func(t *testing.T) {
		mockUC := new(mocks.MockDatabaseMetadataUsecase)
		handler := NewDatabaseMetadataHandler(mockUC)
		app := fiber.New()
		app.Post("/databases/tables", handler.GetTablesByDBNamePost())

		req := httptest.NewRequest(http.MethodPost, "/databases/tables", strings.NewReader(`{"database_name": ""}`))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("failure case - usecase error", func(t *testing.T) {
		mockUC := new(mocks.MockDatabaseMetadataUsecase)
		handler := NewDatabaseMetadataHandler(mockUC)
		app := fiber.New()
		app.Post("/databases/tables", handler.GetTablesByDBNamePost())

		mockUC.On("GetTablesByDBName", mock.Anything, "testdb").Return(nil, errors.New("query failed"))

		req := httptest.NewRequest(http.MethodPost, "/databases/tables", strings.NewReader(`{"database_name": "testdb"}`))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
