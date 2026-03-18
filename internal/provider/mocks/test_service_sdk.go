package mocks

import (
	"context"
	"control-panel-service/internal/provider/dto/request"
	"control-panel-service/internal/provider/dto/response"

	"github.com/stretchr/testify/mock"
)

type MockTestServiceSDK struct {
	mock.Mock
}

func (m *MockTestServiceSDK) Health(ctx context.Context, url request.URL) (response.HealthResponse, error) {
	args := m.Called(ctx, url)

	return args.Get(0).(response.HealthResponse), args.Error(1)
}

func (m *MockTestServiceSDK) GetMetrics(ctx context.Context, url request.URL) (response.MetricsSnapshot, error) {
	args := m.Called(ctx, url)

	return args.Get(0).(response.MetricsSnapshot), args.Error(1)
}

func (m *MockTestServiceSDK) Live(ctx context.Context, url request.URL) (response.HealthResponse, error) {
	args := m.Called(ctx, url)

	return args.Get(0).(response.HealthResponse), args.Error(1)
}

func (m *MockTestServiceSDK) RunExecute(ctx context.Context, url request.URL, runReq request.RunRequest) (response.FactorialExecutionResult, error) {
	args := m.Called(ctx, url, runReq)

	return args.Get(0).(response.FactorialExecutionResult), args.Error(1)
}

func (m *MockTestServiceSDK) ReadyForTest(ctx context.Context, url request.URL) (response.HealthResponse, error) {
	args := m.Called(ctx, url)

	return args.Get(0).(response.HealthResponse), args.Error(1)
}

func (m *MockTestServiceSDK) Pause(ctx context.Context, url request.URL) (*response.PauseResponse, error) {
	args := m.Called(ctx, url)

	var result *response.PauseResponse
	if args.Get(0) != nil {
		result = args.Get(0).(*response.PauseResponse)
	}

	return result, args.Error(1)
}

func (m *MockTestServiceSDK) Resume(ctx context.Context, url request.URL) (*response.ResumeResponse, error) {
	args := m.Called(ctx, url)

	var result *response.ResumeResponse
	if args.Get(0) != nil {
		result = args.Get(0).(*response.ResumeResponse)
	}

	return result, args.Error(1)
}

func (m *MockTestServiceSDK) Stop(ctx context.Context, url request.URL) (*response.StopResponse, error) {
	args := m.Called(ctx, url)

	var result *response.StopResponse
	if args.Get(0) != nil {
		result = args.Get(0).(*response.StopResponse)
	}

	return result, args.Error(1)
}
