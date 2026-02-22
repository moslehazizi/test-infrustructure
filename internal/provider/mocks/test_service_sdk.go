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

func (m *MockTestServiceSDK) Health(ctx context.Context, baseURL string) (response.HealthResponse, error) {
	args := m.Called(ctx, baseURL)

	return args.Get(0).(response.HealthResponse), args.Error(1)
}

func (m *MockTestServiceSDK) GetMetrics(ctx context.Context, baseURL string) (response.MetricsSnapshot, error) {
	args := m.Called(ctx, baseURL)

	return args.Get(0).(response.MetricsSnapshot), args.Error(1)
}

func (m *MockTestServiceSDK) Live(ctx context.Context, baseURL string) (response.HealthResponse, error) {
	args := m.Called(ctx, baseURL)

	return args.Get(0).(response.HealthResponse), args.Error(1)
}

func (m *MockTestServiceSDK) RunExecute(ctx context.Context, baseURL string, runReq request.RunRequest) (response.FactorialExecutionResult, error) {
	args := m.Called(ctx, baseURL, runReq)

	return args.Get(0).(response.FactorialExecutionResult), args.Error(1)
}
