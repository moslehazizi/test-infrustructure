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

const (
	HOST = "Host"
)

type SDKTestService interface {
	// Health is for check test service and its connections to database and kafka.
	Health(ctx context.Context, url request.URL) (response.HealthResponse, error)
	// GetMetrics is for see live situation of test service. how many request sent at the moment, average duration of requests and etc.
	GetMetrics(ctx context.Context, url request.URL) (response.MetricsSnapshot, error)
	// Live is light version of Health, it is just for check test service pod is created or not.
	Live(ctx context.Context, url request.URL) (response.HealthResponse, error)
	// RunExecute is for run test service to send requests to mother service based on its params.
	RunExecute(ctx context.Context, url request.URL, req request.RunRequest) (response.FactorialExecutionResult, error)
	// ReadyForTesting is for checking test service is ready for start or no.
	ReadyForTest(ctx context.Context, url request.URL) (response.HealthResponse, error)
	// Pause is for pause running test service.
	Pause(ctx context.Context, url request.URL) (*response.PauseResponse, error)
	// Resume is for resume paused test service
	Resume(ctx context.Context, url request.URL) (*response.ResumeResponse, error)
	// Stop is for stop running or paused test service
	Stop(ctx context.Context, url request.URL) (*response.StopResponse, error)
}

type sdkTestService struct {
	client HTTPClient
}

func NewSDKTestService(client HTTPClient) *sdkTestService {
	return &sdkTestService{
		client: client,
	}
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func (s *sdkTestService) Health(ctx context.Context, url request.URL) (response.HealthResponse, error) {
	baseUrl := url.BaseURL + "/api/v1/health"
	var result response.HealthResponse

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseUrl, nil)
	if err != nil {
		result.OK = false

		return result, fmt.Errorf("failed to create request to health check: %w", err)
	}

	for head, value := range url.Header {
		req.Header.Set(head, value)
	}
	req.Host = url.Header[HOST]

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

func (s *sdkTestService) GetMetrics(ctx context.Context, url request.URL) (response.MetricsSnapshot, error) {
	baseURL := url.BaseURL + "/api/v1/metrics"

	var result response.MetricsSnapshot

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
	if err != nil {
		return result, fmt.Errorf("failed to create request to metrics: %w", err)
	}

	for head, value := range url.Header {
		req.Header.Set(head, value)
	}
	req.Host = url.Header[HOST]

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

func (s *sdkTestService) Live(ctx context.Context, url request.URL) (response.HealthResponse, error) {
	baseUrl := url.BaseURL + "/api/v1/live"

	var result response.HealthResponse

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseUrl, nil)
	if err != nil {
		result.OK = false
		return result, fmt.Errorf("failed to create request to live check: %w", err)
	}

	for head, value := range url.Header {
		req.Header.Set(head, value)
	}
	req.Host = url.Header[HOST]

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

func (s *sdkTestService) RunExecute(ctx context.Context, url request.URL, runReq request.RunRequest) (response.FactorialExecutionResult, error) {
	baseUrl := url.BaseURL + "/api/v1/run"

	var result response.FactorialExecutionResult

	bodyBytes, err := json.Marshal(runReq)
	if err != nil {
		return result, fmt.Errorf("failed to marshal run request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseUrl,
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		return result, fmt.Errorf("failed to create run execute request: %w", err)
	}

	for head, value := range url.Header {
		req.Header.Set(head, value)
	}
	req.Host = url.Header[HOST]

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

func (s *sdkTestService) ReadyForTest(ctx context.Context, url request.URL) (response.HealthResponse, error) {
	baseUrl := url.BaseURL + "/api/v1/ready-for-testing"
	var result response.HealthResponse

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseUrl, nil)
	if err != nil {
		result.OK = false

		return result, fmt.Errorf("failed to create request to check ready for testing: %w", err)
	}

	for head, value := range url.Header {
		req.Header.Set(head, value)
	}
	req.Host = url.Header[HOST]

	resp, err := s.client.Do(req)
	if err != nil {
		result.OK = false

		return result, fmt.Errorf("failed to do request to check ready for testing: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			result.OK = false

			return result, fmt.Errorf("%w: %w", pkg.ErrFailedToGetReadyForTesting, err)
		}

		return result, fmt.Errorf("%w- check ready for testing return with status code: %d", pkg.ErrFailedToGetReadyForTesting, resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		result.OK = false

		return result, fmt.Errorf("%w: %w", pkg.ErrFailedToGetReadyForTesting, err)
	}

	return result, nil
}

func (s *sdkTestService) Pause(ctx context.Context, url request.URL) (*response.PauseResponse, error) {
	baseUrl := url.BaseURL + "/api/v1/pause"
	var result response.PauseResponse
	var errorResult response.ErrorResponse

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to pause test service: %w", err)
	}

	for head, value := range url.Header {
		req.Header.Set(head, value)
	}
	req.Host = url.Header[HOST]

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request to pause test service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if err := json.NewDecoder(resp.Body).Decode(&errorResult); err != nil {
			return nil, fmt.Errorf("%w: %w", pkg.ErrFailedToPauseTestService, err)
		}

		return nil, fmt.Errorf("%w- pause test service return with status code: %d", pkg.ErrFailedToPauseTestService, resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("%w: %w", pkg.ErrFailedToPauseTestService, err)
	}

	return &result, nil
}

func (s *sdkTestService) Resume(ctx context.Context, url request.URL) (*response.ResumeResponse, error) {
	baseUrl := url.BaseURL + "/api/v1/resume"
	var result response.ResumeResponse
	var errorResult response.ErrorResponse

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to resume test service: %w", err)
	}

	for head, value := range url.Header {
		req.Header.Set(head, value)
	}
	req.Host = url.Header[HOST]

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request to resume test service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if err := json.NewDecoder(resp.Body).Decode(&errorResult); err != nil {
			return nil, fmt.Errorf("%w: %w", pkg.ErrFailedToResumeTestService, err)
		}

		return nil, fmt.Errorf("%w- resume test service return with status code: %d", pkg.ErrFailedToResumeTestService, resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("%w: %w", pkg.ErrFailedToResumeTestService, err)
	}

	return &result, nil
}

func (s *sdkTestService) Stop(ctx context.Context, url request.URL) (*response.StopResponse, error) {
	baseUrl := url.BaseURL + "/api/v1/stop"
	var result response.StopResponse
	var errorResult response.ErrorResponse

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to stop test service: %w", err)
	}

	for head, value := range url.Header {
		req.Header.Set(head, value)
	}
	req.Host = url.Header[HOST]

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request to stop test service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if err := json.NewDecoder(resp.Body).Decode(&errorResult); err != nil {
			return nil, fmt.Errorf("%w: %w", pkg.ErrFailedToStopTestService, err)
		}

		return nil, fmt.Errorf("%w- stop test service return with status code: %d", pkg.ErrFailedToStopTestService, resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("%w: %w", pkg.ErrFailedToStopTestService, err)
	}

	return &result, nil
}
