package handler

import (
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/server/dto/response"
	"control-panel-service/internal/usecase/mocks"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
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
	t.Run("success_case", func(t *testing.T) {
		mockUC := new(mocks.MockDatabaseMetadataUsecase)
		handler := NewDatabaseMetadataHandler(mockUC)

		r := chi.NewRouter()
		r.Get("/databases", handler.GetAll)

		expectedDatabases := []string{"db1", "db2"}
		mockUC.On("GetAll", mock.Anything).Return(expectedDatabases, nil)

		req := httptest.NewRequest(http.MethodGet, "/databases", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var result response.DatabaseMetadataDatabasesResponse
		err := json.NewDecoder(rec.Body).Decode(&result)
		assert.NoError(t, err)

		assert.Equal(t, expectedDatabases, result.Data)

		mockUC.AssertExpectations(t)
	})

	t.Run("failure_case_usecase_error", func(t *testing.T) {
		mockUC := new(mocks.MockDatabaseMetadataUsecase)
		handler := NewDatabaseMetadataHandler(mockUC)

		r := chi.NewRouter()
		r.Get("/databases", handler.GetAll)

		mockUC.On("GetAll", mock.Anything).Return([]string{}, errors.New("db error"))

		req := httptest.NewRequest(http.MethodGet, "/databases", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockUC.AssertExpectations(t)
	})
}

func TestStorageHandler_GetTablesByDBNamePost(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		mockUC := new(mocks.MockDatabaseMetadataUsecase)
		handler := NewDatabaseMetadataHandler(mockUC)

		r := chi.NewRouter()
		r.Post("/databases/tables", handler.GetTablesByDBNamePost)

		reqBody := `{"database_name": "testdb"}`

		expectedMotherTables := []string{"mother11", "mother2"}
		expectedTestTables := []string{"test11", "test2"}

		mockUC.On("GetTablesByDBName", mock.Anything, "testdb").
			Return(&entity.TablesByType{
				MotherTables: expectedMotherTables,
				TestTables:   expectedTestTables,
			}, nil)

		req := httptest.NewRequest(
			http.MethodPost,
			"/databases/tables",
			strings.NewReader(reqBody),
		)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var result response.TablesByType
		err := json.NewDecoder(rec.Body).Decode(&result)
		assert.NoError(t, err)

		assert.Equal(t, expectedMotherTables, result.MotherTables)
		assert.Equal(t, expectedTestTables, result.TestTables)

		mockUC.AssertExpectations(t)
	})

	t.Run("success_case_empty_result", func(t *testing.T) {
		mockUC := new(mocks.MockDatabaseMetadataUsecase)
		handler := NewDatabaseMetadataHandler(mockUC)

		r := chi.NewRouter()
		r.Post("/databases/tables", handler.GetTablesByDBNamePost)

		reqBody := `{"database_name": "testdb"}`

		expectedMotherTables := []string{}
		expectedTestTables := []string{}

		// mockUC.
		// 	On("GetTablesByDBName", mock.Anything, "testdb").
		// 	Return([]string{}, []string{}, nil)
		// Return(nil, nil)

		mockUC.
			On("GetTablesByDBName", mock.Anything, "testdb").
			Return(&entity.TablesByType{
				MotherTables: []string{},
				TestTables:   []string{},
			}, nil)

		req := httptest.NewRequest(
			http.MethodPost,
			"/databases/tables",
			strings.NewReader(reqBody),
		)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var result response.TablesByType
		err := json.NewDecoder(rec.Body).Decode(&result)
		require.NoError(t, err)

		assert.Equal(t, expectedMotherTables, result.MotherTables)
		assert.Equal(t, expectedTestTables, result.TestTables)

		mockUC.AssertExpectations(t)
	})

	t.Run("failure_case_bad_request_body", func(t *testing.T) {
		mockUC := new(mocks.MockDatabaseMetadataUsecase)
		handler := NewDatabaseMetadataHandler(mockUC)

		r := chi.NewRouter()
		r.Post("/databases/tables", handler.GetTablesByDBNamePost)

		req := httptest.NewRequest(http.MethodPost, "/databases/tables", strings.NewReader("invalid json"))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("failure_case_missing_database_name", func(t *testing.T) {
		mockUC := new(mocks.MockDatabaseMetadataUsecase)
		handler := NewDatabaseMetadataHandler(mockUC)
		r := chi.NewRouter()
		r.Post("/databases/tables", handler.GetTablesByDBNamePost)

		req := httptest.NewRequest(http.MethodPost, "/databases/tables", strings.NewReader(`{"database_name": ""}`))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("failure_case_usecase_error", func(t *testing.T) {
		mockUC := new(mocks.MockDatabaseMetadataUsecase)
		handler := NewDatabaseMetadataHandler(mockUC)
		r := chi.NewRouter()
		r.Post("/databases/tables", handler.GetTablesByDBNamePost)

		mockUC.On("GetTablesByDBName", mock.Anything, "testdb").Return(nil, errors.New("query failed"))

		req := httptest.NewRequest(http.MethodPost, "/databases/tables", strings.NewReader(`{"database_name": "testdb"}`))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
