package handler

import (
	"control-panel-service/internal/server/dto/response"
	"control-panel-service/pkg"
	"control-panel-service/pkg/logger"
	responseWriter "control-panel-service/pkg/responsewriter"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type TestScenarioOperationHandler struct {
	testScenarioOperation TestScenarioOperation
}

func NewTestScenarioOperationHandler(testScenarioOperation TestScenarioOperation) *TestScenarioOperationHandler {
	return &TestScenarioOperationHandler{
		testScenarioOperation: testScenarioOperation,
	}
}

// Start godoc
//
//	@Summary		Start a test scenario.
//	@Description	Retrieve a specific test scenario by its ID and start the scenario.
//	@Tags			test-scenarios
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Test scenario ID"
//	@Success		200	{object}	response.SuccessResponse
//	@Failure		400	{object}	response.ErrorResponse
//	@Failure		404	{object}	response.ErrorResponse
//	@Failure		500	{object}	response.ErrorResponse
//	@Router			/api/v1/test-scenarios/{id}/start [post]
func (handler *TestScenarioOperationHandler) Start(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("test-scenario-handler")
	traceCtx, span := tracer.Start(ctx, "start-test-scenario-handler")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	strID := strings.TrimSpace(chi.URLParam(r, "id"))

	id, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "invalid_id_in_params"))

		httpError := pkg.ToHTTPError(pkg.ErrInvalidIDInParams)
		httpError.WriteError(w)

		return
	}

	span.SetAttributes(attribute.String("test_scenario.id", strconv.FormatUint(id, 10)))

	err = handler.testScenarioOperation.Start(traceCtx, id)
	if err != nil {
		span.SetAttributes(
			attribute.String("error.type", "start_error"),
			attribute.String("error.message", err.Error()),
		)

		httpError := pkg.ToHTTPError(err)
		httpError.WriteError(w)

		return
	}

	responseWriter.WriteJSON(w, http.StatusOK, response.SuccessResponse{
		Message: pkg.TestScenarioStarted,
	})
}

// Pause godoc
//
//	@Summary		Pause a test scenario.
//	@Description	Retrieve a specific test scenario by its ID and pause the scenario.
//	@Tags			test-scenarios
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Test scenario ID"
//	@Success		200	{object}	response.SuccessResponse
//	@Failure		400	{object}	response.ErrorResponse
//	@Failure		404	{object}	response.ErrorResponse
//	@Failure		500	{object}	response.ErrorResponse
//	@Router			/api/v1/test-scenarios/{id}/pause [post]
func (handler *TestScenarioOperationHandler) Pause(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("test-scenario-handler")
	traceCtx, span := tracer.Start(ctx, "pause-test-scenario-handler")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	strID := strings.TrimSpace(chi.URLParam(r, "id"))

	id, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "invalid_id_in_params"))

		httpError := pkg.ToHTTPError(pkg.ErrInvalidIDInParams)
		httpError.WriteError(w)

		return
	}

	span.SetAttributes(attribute.String("test_scenario.id", strconv.FormatUint(id, 10)))

	err = handler.testScenarioOperation.Pause(traceCtx, id)
	if err != nil {
		span.SetAttributes(
			attribute.String("error.type", "pause_error"),
			attribute.String("error.message", err.Error()),
		)

		httpError := pkg.ToHTTPError(err)
		httpError.WriteError(w)

		return
	}

	responseWriter.WriteJSON(w, http.StatusOK, response.SuccessResponse{
		Message: pkg.TestScenarioPaused,
	})
}

// Resume godoc
//
//	@Summary		Resume a test scenario.
//	@Description	Retrieve a specific test scenario by its ID and resume the scenario.
//	@Tags			test-scenarios
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Test scenario ID"
//	@Success		200	{object}	response.SuccessResponse
//	@Failure		400	{object}	response.ErrorResponse
//	@Failure		404	{object}	response.ErrorResponse
//	@Failure		500	{object}	response.ErrorResponse
//	@Router			/api/v1/test-scenarios/{id}/resume [post]
func (handler *TestScenarioOperationHandler) Resume(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("test-scenario-handler")
	traceCtx, span := tracer.Start(ctx, "resume-test-scenario-handler")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	strID := strings.TrimSpace(chi.URLParam(r, "id"))

	id, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "invalid_id_in_params"))

		httpError := pkg.ToHTTPError(pkg.ErrInvalidIDInParams)
		httpError.WriteError(w)

		return
	}

	span.SetAttributes(attribute.String("test_scenario.id", strconv.FormatUint(id, 10)))

	err = handler.testScenarioOperation.Resume(traceCtx, id)
	if err != nil {
		span.SetAttributes(
			attribute.String("error.type", "resume_error"),
			attribute.String("error.message", err.Error()),
		)

		httpError := pkg.ToHTTPError(err)
		httpError.WriteError(w)

		return
	}

	responseWriter.WriteJSON(w, http.StatusOK, response.SuccessResponse{
		Message: pkg.TestScenarioResumed,
	})
}

// Stop godoc
//
//	@Summary		Stop a test scenario.
//	@Description	Stop a specific test scenario by its ID.
//	@Tags			test-scenarios
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Test scenario ID"
//	@Success		200	{object}	response.SuccessResponse
//	@Failure		400	{object}	response.ErrorResponse
//	@Failure		404	{object}	response.ErrorResponse
//	@Failure		500	{object}	response.ErrorResponse
//	@Router			/api/v1/test-scenarios/{id}/stop [post]
func (handler *TestScenarioOperationHandler) Stop(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("test-scenario-handler")
	traceCtx, span := tracer.Start(ctx, "stop-test-scenario-handler")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	strID := strings.TrimSpace(chi.URLParam(r, "id"))

	id, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "invalid_id_in_params"))

		httpError := pkg.ToHTTPError(pkg.ErrInvalidIDInParams)
		httpError.WriteError(w)

		return
	}

	span.SetAttributes(attribute.String("test_scenario.id", strconv.FormatUint(id, 10)))

	err = handler.testScenarioOperation.Stop(traceCtx, id)
	if err != nil {
		span.SetAttributes(
			attribute.String("error.type", "stop_error"),
			attribute.String("error.message", err.Error()),
		)

		httpError := pkg.ToHTTPError(err)
		httpError.WriteError(w)

		return
	}

	responseWriter.WriteJSON(w, http.StatusOK, response.SuccessResponse{
		Message: pkg.TestScenarioStop,
	})
}

// Delete godoc
//
//	@Summary		Delete a test scenario.
//	@Description	Delete a specific test scenario by its ID.
//	@Tags			test-scenarios
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Test scenario ID"
//	@Success		200	{object}	response.SuccessResponse
//	@Failure		400	{object}	response.ErrorResponse
//	@Failure		404	{object}	response.ErrorResponse
//	@Failure		500	{object}	response.ErrorResponse
//	@Router			/api/v1/test-scenarios/{id}/delete [post]
func (handler *TestScenarioOperationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("test-scenario-handler")
	traceCtx, span := tracer.Start(ctx, "delete-test-scenario-handler")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	strID := strings.TrimSpace(chi.URLParam(r, "id"))

	id, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "invalid_id_in_params"))

		httpError := pkg.ToHTTPError(pkg.ErrInvalidIDInParams)
		httpError.WriteError(w)

		return
	}

	span.SetAttributes(attribute.String("test_scenario.id", strconv.FormatUint(id, 10)))

	err = handler.testScenarioOperation.Delete(traceCtx, id)
	if err != nil {
		span.SetAttributes(
			attribute.String("error.type", "delete_error"),
			attribute.String("error.message", err.Error()),
		)

		httpError := pkg.ToHTTPError(err)
		httpError.WriteError(w)

		return
	}

	responseWriter.WriteJSON(w, http.StatusOK, response.SuccessResponse{
		Message: pkg.TestScenarioDelete,
	})
}
