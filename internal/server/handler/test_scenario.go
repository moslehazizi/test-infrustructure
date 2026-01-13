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
			MaxRequests: req.Config.MaxRequests,
			MaxDuration: req.Config.MaxDuration,
			RequestDelayDuration: func() *int {
				if req.Config.RequestDelayDuration != nil {
					return req.Config.RequestDelayDuration
				}

				return nil
			}(),
			RandomRequestDelayMin: func() *int {
				if req.Config.RandomRequestDelayMin != nil {
					return req.Config.RandomRequestDelayMin
				}

				return nil
			}(),
			RandomRequestDelayMax: func() *int {
				if req.Config.RandomRequestDelayMax != nil {
					return req.Config.RandomRequestDelayMax
				}

				return nil
			}(),
			FixedTestNumber: func() *int {
				if req.Config.FixedTestNumber != nil {
					return req.Config.FixedTestNumber
				}

				return nil
			}(),
			RandomTestNumberMin: func() *int {
				if req.Config.RandomTestNumberMin != nil {
					return req.Config.RandomTestNumberMin
				}

				return nil
			}(),
			RandomTestNumberMax: func() *int {
				if req.Config.RandomTestNumberMax != nil {
					return req.Config.RandomTestNumberMax
				}

				return nil
			}(),
			BadValueRate:        req.Config.BadValueRate,
			NegativeValueRate:   req.Config.NegativeValueRate,
			ZeroValueRate:       req.Config.ZeroValueRate,
			StringValueRate:     req.Config.StringValueRate,
			RealValueRate:       req.Config.RealValueRate,
			LongStringValueRate: req.Config.LongStringValueRate,
			NullValueRate:       req.Config.NullValueRate,
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
		req := new(request.TestScenarioPaginationRequest)
		if err := ctx.BodyParser(req); err != nil {
			return pkg.ToHTTPError(pkg.ErrBadRequest).AsFiber(ctx)
		}

		items, err := handler.testScenario.GetPaginated(ctx.Context(), entity.TestScenarioPaginationRequest{
			Page:    req.Page,
			PerPage: req.PerPage,
		})
		if err != nil {
			return pkg.ToHTTPError(err).AsFiber(ctx)
		}

		var responses []response.TestScenario
		for _, item := range items {
			responses = append(responses, response.TestScenario{
				ID:                  item.ID,
				Name:                item.Name,
				CreatedAt:           item.CreatedAt,
				UpdatedAt:           item.UpdatedAt,
				Status:              item.Status,
				MaxTestServiceCount: item.MaxTestServiceCount,
				AutoStepChangeRate:  item.AutoStepChangeRate,
				ExecutionDuration:   item.ExecutionDuration,
				TestCategory: func() *response.TestCategory {
					if item.TestCategory == nil {
						return nil
					}

					return &response.TestCategory{
						ID:                     item.TestCategory.ID,
						Name:                   item.TestCategory.Name,
						Label:                  item.TestCategory.Label,
						HasMaxTestServiceCount: item.TestCategory.HasMaxTestServiceCount,
						HasExecutionDuration:   item.TestCategory.HasExecutionDuration,
						HasAutoStepChangeRate:  item.TestCategory.HasAutoStepChangeRate,
						CreatedAt:              item.TestCategory.CreatedAt,
						UpdatedAt:              item.TestCategory.UpdatedAt,
					}
				}(),
				MotherService: func() *response.MotherService {
					if item.MotherService == nil {
						return nil
					}

					return &response.MotherService{
						ID:                       item.MotherService.ID,
						CreatedAt:                item.MotherService.CreatedAt,
						UpdatedAt:                item.MotherService.UpdatedAt,
						Name:                     item.MotherService.Name,
						ExceptionRate:            item.MotherService.ExceptionRate,
						ResponseDelayRate:        item.MotherService.ResponseDelayRate,
						ResponseDelayDuration:    item.MotherService.ResponseDelayDuration,
						RandomResponseDelayMin:   item.MotherService.RandomResponseDelayMin,
						RandomResponseDelayMax:   item.MotherService.RandomResponseDelayMax,
						Status:                   item.MotherService.Status,
						ServiceDeploymentAddress: item.MotherService.ServiceDeploymentAddress,
						DatabaseName:             item.MotherService.DatabaseName,
						DatabaseTableName:        item.MotherService.DatabaseTableName,
					}
				}(),
			})
		}

		return ctx.Status(http.StatusOK).JSON(&response.Paginated[[]response.TestScenario]{
			Data:    responses,
			Page:    req.Page,
			PerPage: req.PerPage,
		})
	}
}
