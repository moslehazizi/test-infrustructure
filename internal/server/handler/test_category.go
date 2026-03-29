package handler

import (
	"control-panel-service/config"
	"control-panel-service/internal/server/dto/response"
	"control-panel-service/internal/usecase"
	"control-panel-service/pkg"
	"control-panel-service/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func NewTestCategoryHandler(cfg *config.Config, testCategoryService usecase.TestCategoryService) *TestCategoryHandler {
	return &TestCategoryHandler{
		cfg,
		testCategoryService,
	}
}

type TestCategoryHandler struct {
	cfg                 *config.Config
	testCategoryService usecase.TestCategoryService
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
func (handler *TestCategoryHandler) GetAll() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tracer := otel.Tracer("test-category-handler")
		traceCtx, span := tracer.Start(ctx.Context(), "get_test_categories")
		defer span.End()

		requestID := logger.GetRequestID(ctx.Context())
		span.SetAttributes(attribute.String("request_id", requestID))

		svcResults, err := handler.testCategoryService.GetAll(traceCtx)
		if err != nil {
			span.SetAttributes(attribute.String("error.type", "get_all_error"), attribute.String("error.message", err.Error()))
			return pkg.ToHTTPError(err).AsFiber(ctx)
		}

		var responses []response.TestCategory
		for _, svcResult := range svcResults {
			responses = append(responses, response.TestCategory{
				ID:                     svcResult.ID,
				CreatedAt:              svcResult.CreatedAt,
				UpdatedAt:              svcResult.UpdatedAt,
				Name:                   svcResult.Name,
				Label:                  svcResult.Label,
				HasMaxTestServiceCount: svcResult.HasMaxTestServiceCount,
				HasNumSteps:            svcResult.HasNumSteps,
				Active:                 svcResult.Active,
			})
		}

		return ctx.Status(http.StatusOK).JSON(responses)
	}
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
func (handler *TestCategoryHandler) GetByID() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tracer := otel.Tracer("test-category-handler")
		traceCtx, span := tracer.Start(ctx.Context(), "get_test_category_by_id")
		defer span.End()

		requestID := logger.GetRequestID(ctx.Context())
		span.SetAttributes(attribute.String("request_id", requestID))

		idParam := ctx.Params("id")
		if idParam == "" {
			span.SetAttributes(attribute.String("error.type", "missing_id"))
			return pkg.ToHTTPError(pkg.ErrPageNotFound).AsFiber(ctx)
		}

		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			span.SetAttributes(attribute.String("error.type", "invalid_id_in_params"))
			return pkg.ToHTTPError(pkg.ErrInvalidIDInParams).AsFiber(ctx)
		}

		span.SetAttributes(attribute.String("test_category.id", strconv.FormatUint(id, 10)))

		svcResult, err := handler.testCategoryService.GetByID(traceCtx, id)
		if err != nil {
			span.SetAttributes(attribute.String("error.type", "get_by_id_error"), attribute.String("error.message", err.Error()))
			return pkg.ToHTTPError(err).AsFiber(ctx)
		}

		responses := response.TestCategory{
			ID:                     svcResult.ID,
			CreatedAt:              svcResult.CreatedAt,
			UpdatedAt:              svcResult.UpdatedAt,
			Name:                   svcResult.Name,
			Label:                  svcResult.Label,
			HasMaxTestServiceCount: svcResult.HasMaxTestServiceCount,
			HasNumSteps:            svcResult.HasNumSteps,
			Active:                 svcResult.Active,
		}

		return ctx.Status(http.StatusOK).JSON(responses)
	}
}
