package handler

import (
	"control-panel-service/config"
	"control-panel-service/internal/server/dto/response"
	"control-panel-service/internal/usecase"
	"control-panel-service/pkg"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
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
		svcResults, err := handler.testCategoryService.GetAll(ctx.Context())
		if err != nil {
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
				HasExecutionDuration:   svcResult.HasExecutionDuration,
				HasAutoStepChangeRate:  svcResult.HasAutoStepChangeRate,
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
		idParam := ctx.Params("id")
		if idParam == "" {
			return pkg.ToHTTPError(pkg.ErrPageNotFound).AsFiber(ctx)
		}

		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			return pkg.ToHTTPError(pkg.ErrInvalidIDInParams).AsFiber(ctx)
		}

		svcResult, err := handler.testCategoryService.GetByID(ctx.Context(), id)
		if err != nil {
			return pkg.ToHTTPError(err).AsFiber(ctx)
		}

		responses := response.TestCategory{
			ID:                     svcResult.ID,
			CreatedAt:              svcResult.CreatedAt,
			UpdatedAt:              svcResult.UpdatedAt,
			Name:                   svcResult.Name,
			Label:                  svcResult.Label,
			HasMaxTestServiceCount: svcResult.HasMaxTestServiceCount,
			HasExecutionDuration:   svcResult.HasExecutionDuration,
			HasAutoStepChangeRate:  svcResult.HasAutoStepChangeRate,
		}

		return ctx.Status(http.StatusOK).JSON(responses)
	}
}
