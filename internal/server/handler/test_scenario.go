package handler

import (
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/server/dto/request"
	"control-panel-service/internal/server/dto/response"
	"control-panel-service/internal/usecase"
	"control-panel-service/pkg"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type TestScenario struct {
	testScenario usecase.TestScenario
}

func NewTestScenarioHandler(testScenario usecase.TestScenario) *TestScenario {
	return &TestScenario{
		testScenario: testScenario,
	}
}

func (handler *TestScenario) Create() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		req := new(request.TestScenario)

		if err := ctx.BodyParser(req); err != nil {
			return pkg.ToHTTPError(pkg.ErrBadRequest).AsFiber(ctx)
		}

		testScenario := &entity.TestScenario{
			Name:            req.Name,
			TestCategoryID:  req.TestCategoryID,
			MotherServiceID: req.MotherServiceID,
			MaxTestServiceCount: func() *int {
				if req.MaxTestServiceCount != nil {
					return req.MaxTestServiceCount
				}

				return nil
			}(),
			ExecutionDuration: func() *int {
				if req.ExecutionDuration != nil {
					return req.ExecutionDuration
				}

				return nil
			}(),
			AutoStepChangeRate: func() *int {
				if req.AutoStepChangeRate != nil {
					return req.AutoStepChangeRate
				}

				return nil
			}(),
		}
		testSvcCfg := &entity.TestServiceConfig{
			MaxRequests: req.MaxRequests,
			MaxDuration: req.MaxDuration,
			RequestDelayDuration: func() *int {
				if req.RequestDelayDuration != nil {
					return req.RequestDelayDuration
				}

				return nil
			}(),
			RandomRequestDelayMin: func() *int {
				if req.RandomRequestDelayMin != nil {
					return req.RandomRequestDelayMin
				}

				return nil
			}(),
			RandomRequestDelayMax: func() *int {
				if req.RandomRequestDelayMax != nil {
					return req.RandomRequestDelayMax
				}

				return nil
			}(),
			FixedTestNumber: func() *int {
				if req.FixedTestNumber != nil {
					return req.FixedTestNumber
				}

				return nil
			}(),
			RandomTestNumberMin: func() *int {
				if req.RandomTestNumberMin != nil {
					return req.RandomTestNumberMin
				}

				return nil
			}(),
			RandomTestNumberMax: func() *int {
				if req.RandomTestNumberMax != nil {
					return req.RandomTestNumberMax
				}

				return nil
			}(),
			BadValueRate:        req.BadValueRate,
			NegativeValueRate:   req.NegativeValueRate,
			ZeroValueRate:       req.ZeroValueRate,
			StringValueRate:     req.StringValueRate,
			RealValueRate:       req.RealValueRate,
			LongStringValueRate: req.LongStringValueRate,
			NullValueRate:       req.NullValueRate,
		}

		err := handler.testScenario.Create(ctx.Context(), testScenario, testSvcCfg)
		if err != nil {
			return pkg.ToHTTPError(err).AsFiber(ctx)
		}

		return ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"message": pkg.CreateTestScenarioSuccessfully,
		})
	}
}

func (handler *TestScenario) GetPaginated() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		req := new(request.PaginationRequest)
		if err := ctx.BodyParser(req); err != nil {
			return pkg.ToHTTPError(pkg.ErrBadRequest).AsFiber(ctx)
		}

		items, err := handler.testScenario.GetPaginated(ctx.Context(), entity.TestScenarioPaginationRequest{
			Page:    req.Page,
			PerPage: req.PerPage,
		})
		_ = err

		var responses []response.TestScenario
		for _, item := range items {
			responses = append(responses, response.TestScenario{
				ID:                  item.ID,
				Name:                item.Name,
				CreatedAt:           item.CreatedAt,
				UpdatedAt:           item.UpdatedAt,
				TestCategoryID:      item.TestCategoryID,
				MotherServiceID:     item.MotherServiceID,
				Status:              item.Status,
				MaxTestServiceCount: item.MaxTestServiceCount,
				ExecutionDuration:   item.ExecutionDuration,
				AutoStepChangeRate:  item.AutoStepChangeRate,
			})
		}

		return ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"data":     responses,
			"page":     req.Page,
			"per_page": req.PerPage,
		})
	}
}
