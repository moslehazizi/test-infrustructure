package provider

import (
	"context"
	"control-panel-service/internal/provider/dto/request"
	"control-panel-service/internal/provider/mocks"
	"control-panel-service/pkg"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewSDKTestService(t *testing.T) {
	mockClient := new(mocks.MockHTTPClient)
	sdk := NewSDKTestService(mockClient)

	assert.NotNil(t, sdk)
}

func TestSDKTestService_Health(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		mockClient := new(mocks.MockHTTPClient)

		body := io.NopCloser(strings.NewReader(`{"ok": true}`))

		mockResp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       body,
		}

		mockClient.
			On("Do", mock.AnythingOfType("*http.Request")).
			Return(mockResp, nil)

		sdk := NewSDKTestService(mockClient)

		resp, err := sdk.Health(context.Background(), "http://localhost:8085")

		require.NoError(t, err)
		assert.True(t, resp.OK)

		mockClient.AssertExpectations(t)
	})

	t.Run("failed case - invalid host", func(t *testing.T) {
		mockClient := new(mocks.MockHTTPClient)

		mockClient.
			On("Do", mock.Anything).
			Return(nil, errors.New("network error"))

		sdk := NewSDKTestService(mockClient)

		_, err := sdk.Health(context.Background(), "http://localhost:8085")

		assert.Error(t, err)

		mockClient.AssertExpectations(t)
	})

	t.Run("failed case - status code not ok", func(t *testing.T) {
		mockClient := new(mocks.MockHTTPClient)

		body := io.NopCloser(strings.NewReader(`{"ok": false}`))

		mockResp := &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       body,
		}

		mockClient.
			On("Do", mock.AnythingOfType("*http.Request")).
			Return(mockResp, nil)

		sdk := NewSDKTestService(mockClient)

		resp, err := sdk.Health(context.Background(), "http://localhost:8085")

		assert.Error(t, err)
		assert.False(t, resp.OK)

		mockClient.AssertExpectations(t)
	})
}

func TestSDKTestService_GetMetric(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		mockClient := new(mocks.MockHTTPClient)

		body := io.NopCloser(strings.NewReader(`{
		"requests": 10,
		"success": 8,
		"failed": 2,
		"min_duration": 5,
		"max_duration": 100,
		"avg_duration": 50
	}`))

		mockResp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       body,
		}

		mockClient.
			On("Do", mock.MatchedBy(func(req *http.Request) bool {
				return req.Method == http.MethodGet &&
					req.URL.Path == "/api/v1/metrics"
			})).
			Return(mockResp, nil)

		sdk := NewSDKTestService(mockClient)

		resp, err := sdk.GetMetrics(context.Background(), "http://localhost:8085")

		require.NoError(t, err)
		assert.Equal(t, int64(10), resp.Requests)
		assert.Equal(t, int64(8), resp.Success)

		mockClient.AssertExpectations(t)
	})

	t.Run("failed case - network error", func(t *testing.T) {
		mockClient := new(mocks.MockHTTPClient)

		mockClient.
			On("Do", mock.Anything).
			Return(nil, errors.New("network error"))

		sdk := NewSDKTestService(mockClient)

		_, err := sdk.GetMetrics(context.Background(), "http://localhost:8085")

		assert.Error(t, err)

		mockClient.AssertExpectations(t)
	})

	t.Run("failed case - status code not ok", func(t *testing.T) {
		mockClient := new(mocks.MockHTTPClient)

		body := io.NopCloser(strings.NewReader(`{
		"requests": 0
	}`))

		mockResp := &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       body,
		}

		mockClient.
			On("Do", mock.Anything).
			Return(mockResp, nil)

		sdk := NewSDKTestService(mockClient)

		resp, err := sdk.GetMetrics(context.Background(), "http://localhost:8085")

		assert.Error(t, err)
		assert.Equal(t, int64(0), resp.Requests)

		mockClient.AssertExpectations(t)
	})
}

func TestSDKTestService_Live(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		mockClient := new(mocks.MockHTTPClient)

		body := io.NopCloser(strings.NewReader(`{"ok": true}`))

		mockResp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       body,
		}

		mockClient.
			On("Do", mock.MatchedBy(func(req *http.Request) bool {
				return req.Method == http.MethodGet &&
					req.URL.Path == "/api/v1/live"
			})).
			Return(mockResp, nil)

		sdk := NewSDKTestService(mockClient)

		resp, err := sdk.Live(context.Background(), "http://localhost:8085")

		require.NoError(t, err)
		assert.True(t, resp.OK)

		mockClient.AssertExpectations(t)
	})

	t.Run("live failed - network error", func(t *testing.T) {
		mockClient := new(mocks.MockHTTPClient)

		mockClient.
			On("Do", mock.Anything).
			Return(nil, errors.New("network error"))

		sdk := NewSDKTestService(mockClient)

		resp, err := sdk.Live(context.Background(), "http://localhost:8085")

		assert.Error(t, err)
		assert.False(t, resp.OK)

		mockClient.AssertExpectations(t)
	})

	t.Run("live failed - status not ok", func(t *testing.T) {
		mockClient := new(mocks.MockHTTPClient)

		body := io.NopCloser(strings.NewReader(`{"ok": false}`))

		mockResp := &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       body,
		}

		mockClient.
			On("Do", mock.Anything).
			Return(mockResp, nil)

		sdk := NewSDKTestService(mockClient)

		resp, err := sdk.Live(context.Background(), "http://localhost:8085")

		assert.Error(t, err)
		assert.False(t, resp.OK)

		mockClient.AssertExpectations(t)
	})
}

func TestSDKTestService_Run(t *testing.T) {
	t.Run("run execute success", func(t *testing.T) {
		mockClient := new(mocks.MockHTTPClient)

		runReq := request.RunRequest{
			StepNum:     1,
			ExecutionId: uuid.New(),
		}

		body := io.NopCloser(strings.NewReader(`{
		"success_count": 5,
		"failed_count": 1,
		"max_duration": 100,
		"min_duration": 10,
		"ave_duration": 50
	}`))

		mockResp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       body,
		}

		mockClient.
			On("Do", mock.MatchedBy(func(req *http.Request) bool {
				return req.Method == http.MethodPost &&
					req.URL.Path == "/api/v1/run"
			})).
			Return(mockResp, nil)

		sdk := NewSDKTestService(mockClient)

		resp, err := sdk.RunExecute(context.Background(), "http://localhost:8085", runReq)

		require.NoError(t, err)
		assert.Equal(t, int64(5), resp.SuccessCount)

		mockClient.AssertExpectations(t)
	})

	t.Run("run execute conflict", func(t *testing.T) {
		mockClient := new(mocks.MockHTTPClient)

		mockResp := &http.Response{
			StatusCode: http.StatusConflict,
			Body:       io.NopCloser(strings.NewReader(`{"error":"already running"}`)),
		}

		mockClient.
			On("Do", mock.Anything).
			Return(mockResp, nil)

		sdk := NewSDKTestService(mockClient)

		_, err := sdk.RunExecute(context.Background(), "http://localhost:8085", request.RunRequest{})

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrRunAlreadyInProgress)

		mockClient.AssertExpectations(t)
	})

	t.Run("run execute internal error", func(t *testing.T) {
		mockClient := new(mocks.MockHTTPClient)

		mockResp := &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(strings.NewReader(`{"error":"internal"}`)),
		}

		mockClient.
			On("Do", mock.Anything).
			Return(mockResp, nil)

		sdk := NewSDKTestService(mockClient)

		_, err := sdk.RunExecute(context.Background(), "http://localhost:8085", request.RunRequest{})

		assert.Error(t, err)

		mockClient.AssertExpectations(t)
	})

	t.Run("run execute network error", func(t *testing.T) {
		mockClient := new(mocks.MockHTTPClient)

		mockClient.
			On("Do", mock.Anything).
			Return(nil, errors.New("network error"))

		sdk := NewSDKTestService(mockClient)

		_, err := sdk.RunExecute(context.Background(), "http://localhost:8085", request.RunRequest{})

		assert.Error(t, err)

		mockClient.AssertExpectations(t)
	})
}

func TestSDKTestService_ReadyForTesting(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		mockClient := new(mocks.MockHTTPClient)

		body := io.NopCloser(strings.NewReader(`{"ok": true}`))

		mockResp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       body,
		}

		mockClient.
			On("Do", mock.AnythingOfType("*http.Request")).
			Return(mockResp, nil)

		sdk := NewSDKTestService(mockClient)

		resp, err := sdk.ReadyForTest(context.Background(), "http://localhost:8085")

		require.NoError(t, err)
		assert.True(t, resp.OK)

		mockClient.AssertExpectations(t)
	})

	t.Run("failed case - invalid host", func(t *testing.T) {
		mockClient := new(mocks.MockHTTPClient)

		mockClient.
			On("Do", mock.Anything).
			Return(nil, errors.New("network error"))

		sdk := NewSDKTestService(mockClient)

		_, err := sdk.ReadyForTest(context.Background(), "http://localhost:8085")

		assert.Error(t, err)

		mockClient.AssertExpectations(t)
	})

	t.Run("failed case - status code not ok", func(t *testing.T) {
		mockClient := new(mocks.MockHTTPClient)

		body := io.NopCloser(strings.NewReader(`{"ok": false}`))

		mockResp := &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       body,
		}

		mockClient.
			On("Do", mock.AnythingOfType("*http.Request")).
			Return(mockResp, nil)

		sdk := NewSDKTestService(mockClient)

		resp, err := sdk.ReadyForTest(context.Background(), "http://localhost:8085")

		assert.Error(t, err)
		assert.False(t, resp.OK)

		mockClient.AssertExpectations(t)
	})
}
