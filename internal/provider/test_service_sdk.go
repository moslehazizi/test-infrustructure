package provider

import (
	"bytes"
	"context"
	"control-panel-service/internal/provider/dto/request"
	"control-panel-service/internal/provider/dto/response"
	"control-panel-service/pkg"
	"encoding/json"
	"fmt"
	"net/http"
)

type SDKTestService interface {
	// Health is for check test service and its connections to database and kafka.
	Health(ctx context.Context, baseURL string) (response.HealthResponse, error)
	// GetMetrics is for see live situation of test service. how many request sent at the moment, average duration of requests and etc.
	GetMetrics(ctx context.Context, baseURL string) (response.MetricsSnapshot, error)
	// Live is light version of Health, it is just for check test service pod is created or not.
	Live(ctx context.Context, baseURL string) (response.HealthResponse, error)
	// RunExecute is for run test service to send requests to mother service based on its params.
	RunExecute(ctx context.Context, baseURL string, req request.RunRequest) (response.FactorialExecutionResult, error)
}

type sdkTestService struct {
	client HTTPClient
}

func NewSDKTestService(client HTTPClient) SDKTestService {
	return &sdkTestService{
		client: client,
	}
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func (s *sdkTestService) Health(
	ctx context.Context,
	baseURL string,
) (response.HealthResponse, error) {
	url := fmt.Sprintf("%s/api/v1/health", baseURL)
	var result response.HealthResponse

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		result.OK = false

		return result, fmt.Errorf("failed to create request to health check: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		result.OK = false

		return result, fmt.Errorf("failed to do request to health check: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			result.OK = false

			return result, fmt.Errorf("%w: %w", pkg.ErrFailedToGetHealthCheck, err)
		}

		return result, fmt.Errorf("%w- health check return with status code: %d", pkg.ErrFailedToGetHealthCheck, resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		result.OK = false

		return result, fmt.Errorf("%w: %w", pkg.ErrFailedToGetHealthCheck, err)
	}

	return result, nil
}

func (s *sdkTestService) GetMetrics(
	ctx context.Context,
	baseURL string,
) (response.MetricsSnapshot, error) {
	url := fmt.Sprintf("%s/api/v1/metrics", baseURL)

	var result response.MetricsSnapshot

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return result, fmt.Errorf("failed to create request to metrics: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return result, fmt.Errorf("failed to do request to metrics: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result, fmt.Errorf("%w: %w", pkg.ErrFailedToGetMetrics, err)
		}

		return result, fmt.Errorf(
			"%w - metrics returned status code: %d",
			pkg.ErrFailedToGetMetrics,
			resp.StatusCode,
		)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, fmt.Errorf("%w: %w", pkg.ErrFailedToGetMetrics, err)
	}

	return result, nil
}

func (s *sdkTestService) Live(ctx context.Context, baseURL string) (response.HealthResponse, error) {
	url := fmt.Sprintf("%s/api/v1/live", baseURL)

	var result response.HealthResponse

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		result.OK = false
		return result, fmt.Errorf("failed to create request to live check: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		result.OK = false
		return result, fmt.Errorf("failed to do request to live check: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			result.OK = false
			return result, fmt.Errorf("%w: %w", pkg.ErrFailedToGetLiveCheck, err)
		}

		return result, fmt.Errorf(
			"%w - live check returned status code: %d",
			pkg.ErrFailedToGetLiveCheck,
			resp.StatusCode,
		)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		result.OK = false
		return result, fmt.Errorf("%w: %w", pkg.ErrFailedToGetLiveCheck, err)
	}

	return result, nil
}

func (s *sdkTestService) RunExecute(
	ctx context.Context,
	baseURL string,
	runReq request.RunRequest,
) (response.FactorialExecutionResult, error) {
	url := fmt.Sprintf("%s/api/v1/run-execute", baseURL)

	var result response.FactorialExecutionResult

	bodyBytes, err := json.Marshal(runReq)
	if err != nil {
		return result, fmt.Errorf("failed to marshal run request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		return result, fmt.Errorf("failed to create run execute request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return result, fmt.Errorf("failed to do run execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		return result, fmt.Errorf("%w", pkg.ErrRunAlreadyInProgress)
	}

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf(
			"%w - status code: %d",
			pkg.ErrFailedToRunTestService,
			resp.StatusCode,
		)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, fmt.Errorf("failed to decode run execute response: %w", err)
	}

	return result, nil
}
