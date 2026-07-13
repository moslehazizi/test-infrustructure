package handler

import (
	"control-panel-service/config"
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

func NewTestCategoryHandler(cfg *config.Config, testCategoryService TestCategoryService) *TestCategoryHandler {
	return &TestCategoryHandler{
		cfg,
		testCategoryService,
	}
}

type TestCategoryHandler struct {
	cfg                 *config.Config
	testCategoryService TestCategoryService
}

// GetAll godoc
//
//	@Summary		Get all test categories
//	@Description	Retrieve all available test categories
//	@Tags			test-categories
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		response.TestCategory
//	@Failure		500	{object}	response.ErrorResponse
//	@Router			/api/v1/test-categories [get]
func (handler *TestCategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("test-category-handler")
	traceCtx, span := tracer.Start(ctx, "get-test-categories-handler")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	svcResults, err := handler.testCategoryService.GetAll(traceCtx)
	if err != nil {
		span.SetAttributes(
			attribute.String("error.type", "get_all_error"),
			attribute.String("error.message", err.Error()),
		)

		httpError := pkg.ToHTTPError(err)
		httpError.WriteError(w)

		return
	}

	var responses []response.TestCategory
	for _, svcResult := range svcResults {
		var result response.TestCategory
		result.FromTestCategoryEntity(&svcResult)

		responses = append(responses, result)
	}

	responseWriter.WriteJSON(w, http.StatusOK, responses)
}

// GetByID godoc
//
//	@Summary		Get test category by ID
//	@Description	Retrieve a specific test category by its ID
//	@Tags			test-categories
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Test category ID"
//	@Success		200	{object}	response.TestCategory
//	@Failure		400	{object}	response.ErrorResponse
//	@Failure		404	{object}	response.ErrorResponse
//	@Failure		500	{object}	response.ErrorResponse
//	@Router			/api/v1/test-categories/{id} [get]
func (handler *TestCategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("test-category-handler")
	traceCtx, span := tracer.Start(ctx, "get-test-category-by-id-handler")
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

	span.SetAttributes(attribute.String("test_category.id", strconv.FormatUint(id, 10)))

	svcResult, err := handler.testCategoryService.GetByID(traceCtx, id)
	if err != nil {
		span.SetAttributes(
			attribute.String("error.type", "get_by_id_error"),
			attribute.String("error.message", err.Error()),
		)

		httpError := pkg.ToHTTPError(err)
		httpError.WriteError(w)

		return
	}

	var result response.TestCategory
	result.FromTestCategoryEntity(svcResult)

	responseWriter.WriteJSON(w, http.StatusOK, result)
}
