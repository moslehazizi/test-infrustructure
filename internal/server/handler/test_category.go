package handler

import (
	"control-panel-service/config"
	"control-panel-service/internal/server/dto/response"
	"control-panel-service/internal/usecase"
	"control-panel-service/pkg"
	"errors"
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

func (handler *TestCategoryHandler) GetAll() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		svcResults, err := handler.testCategoryService.GetAll(ctx.Context())
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"error": pkg.InternalServerErrorMessage,
			})
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

func (handler *TestCategoryHandler) GetByID() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		idParam := ctx.Params("id")
		if idParam == "" {
			return ctx.Status(http.StatusNotFound).JSON(&fiber.Map{
				"error": "404 not found",
			})
		}

		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"error": pkg.InvalidIDInParams,
			})
		}

		svcResult, err := handler.testCategoryService.GetByID(ctx.Context(), id)
		if err != nil {
			if errors.Is(err, pkg.ErrTestCategoryNotFound) {
				return ctx.Status(http.StatusNotFound).JSON(&fiber.Map{
					"error": pkg.TestCategoryNotFound,
				})
			}

			return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"error": pkg.InternalServerErrorMessage,
			})
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
