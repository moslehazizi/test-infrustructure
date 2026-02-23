package handler

import (
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/server/dto/request"
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

type TestScenario struct {
	testScenario usecase.TestScenario
}

func NewTestScenarioHandler(testScenario usecase.TestScenario) *TestScenario {
	return &TestScenario{
		testScenario: testScenario,
	}
}

// Create godoc
//
//	@Summary		Create a test scenario
//	@Description	Create a new test scenario with configuration
//	@Tags			test-scenarios
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.TestScenario	true	"Request body"
//	@Success		200		{object}	response.SuccessResponse
//	@Failure		400		{object}	response.ErrorResponse
//	@Failure		422		{object}	response.ErrorResponse
//	@Failure		500		{object}	response.ErrorResponse
//	@Router			/api/v1/test-scenarios [post]
func (handler *TestScenario) Create() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		requestID := logger.GetRequestID(ctx.Context())
		tracer := otel.Tracer("test-scenario-handler")
		traceCtx, span := tracer.Start(ctx.Context(), "create_test_scenario")
		defer span.End()

		span.SetAttributes(attribute.String("request_id", requestID))

		req := new(request.TestScenario)

		if err := ctx.BodyParser(req); err != nil {
			span.SetAttributes(attribute.String("error.type", "bad_request"))
			return pkg.ToHTTPError(pkg.ErrBadRequest).AsFiber(ctx)
		}

		span.SetAttributes(
			attribute.String("test_scenario.name", req.Name),
			attribute.String("test_category.id", strconv.FormatUint(req.TestCategoryID, 10)),
			attribute.String("mother_service.id", strconv.FormatUint(req.MotherServiceID, 10)))

		testScenario := &entity.TestScenario{
			Name:                req.Name,
			TestCategoryID:      req.TestCategoryID,
			MotherServiceID:     req.MotherServiceID,
			MaxTestServiceCount: req.MaxTestServiceCount,
			ExecutionDuration:   req.ExecutionDuration,
			AutoStepChangeRate:  req.AutoStepChangeRate,
			TestServiceConfig: func() *entity.TestServiceConfig {
				if req.Config == nil {
					return nil
				}

				return &entity.TestServiceConfig{
					MaxRequests:           req.Config.MaxRequests,
					MaxDuration:           int64(req.Config.MaxDuration),
					RequestDelayDuration:  req.Config.RequestDelayDuration,
					RandomRequestDelayMin: req.Config.RandomRequestDelayMin,
					RandomRequestDelayMax: req.Config.RandomRequestDelayMax,
					FixedTestNumber:       req.Config.FixedTestNumber,
					RandomTestNumberMin:   req.Config.RandomTestNumberMin,
					RandomTestNumberMax:   req.Config.RandomTestNumberMax,
					BadValueRate:          req.Config.BadValueRate,
					NegativeValueRate:     req.Config.NegativeValueRate,
					ZeroValueRate:         req.Config.ZeroValueRate,
					StringValueRate:       req.Config.StringValueRate,
					RealValueRate:         req.Config.RealValueRate,
					LongStringValueRate:   req.Config.LongStringValueRate,
					NullValueRate:         req.Config.NullValueRate,
					DatabaseName:          req.Config.DatabaseName,
					DatabaseTableName:     req.Config.DatabaseTableName,
				}
			}(),
		}

		err := handler.testScenario.Create(traceCtx, testScenario)
		if err != nil {
			span.SetAttributes(attribute.String("error.type", "create_error"), attribute.String("error.message", err.Error()))
			return pkg.ToHTTPError(err).AsFiber(ctx)
		}

		return ctx.Status(http.StatusOK).JSON(&response.SuccessResponse{
			Message: pkg.CreateTestScenarioSuccessfully,
		})
	}
}

// GetPaginated godoc
//
//	@Summary		Get paginated test scenarios
//	@Description	Get test scenarios with pagination support
//	@Tags			test-scenarios
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.TestScenarioPaginationRequest	true	"Pagination request with page and per_page"
//	@Success		200		{object}	response.PaginatedTestScenario
//	@Failure		400		{object}	response.ErrorResponse
//	@Failure		422		{object}	response.ErrorResponse
//	@Failure		500		{object}	response.ErrorResponse
//	@Router			/api/v1/test-scenarios/search [post]
func (handler *TestScenario) GetPaginated() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tracer := otel.Tracer("test-scenario-handler")
		traceCtx, span := tracer.Start(ctx.Context(), "get_paginated_test_scenarios")
		defer span.End()

		requestID := logger.GetRequestID(ctx.Context())
		span.SetAttributes(attribute.String("request_id", requestID))

		req := new(request.TestScenarioPaginationRequest)
		if err := ctx.BodyParser(req); err != nil {
			span.SetAttributes(attribute.String("error.type", "bad_request"))
			return pkg.ToHTTPError(pkg.ErrBadRequest).AsFiber(ctx)
		}

		span.SetAttributes(attribute.String("pagination.page", strconv.Itoa(req.Page)), attribute.String("pagination.per_page", strconv.Itoa(req.PerPage)))

		items, count, err := handler.testScenario.GetPaginated(traceCtx, entity.TestScenarioPaginationRequest{
			Page:    req.Page,
			PerPage: req.PerPage,
		})
		if err != nil {
			span.SetAttributes(attribute.String("error.type", "get_paginated_error"), attribute.String("error.message", err.Error()))
			return pkg.ToHTTPError(err).AsFiber(ctx)
		}

		var responses []response.TestScenario
		for _, item := range items {
			responses = append(responses, response.TestScenario{
				ID:                  item.ID,
				Name:                item.Name,
				CreatedAt:           item.CreatedAt,
				UpdatedAt:           item.UpdatedAt,
				StartedAt:           item.StartedAt,
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
				Editable: item.Editable,
			})
		}

		return ctx.Status(http.StatusOK).JSON(&response.PaginatedTestScenario{
			Page:    req.Page,
			PerPage: req.PerPage,
			Data:    responses,
			Total:   count,
		})
	}
}

// GetByID godoc
//
//	@Summary		Get test scenario by ID
//	@Description	Retrieve a specific test scenario by its ID
//	@Tags			test-scenarios
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Test scenario ID"
//	@Success		200	{object}	response.TestScenarioResponseByID
//	@Failure		400	{object}	response.ErrorResponse
//	@Failure		404	{object}	response.ErrorResponse
//	@Failure		500	{object}	response.ErrorResponse
//	@Router			/api/v1/test-scenarios/{id} [get]
func (handler *TestScenario) GetByID() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tracer := otel.Tracer("test-scenario-handler")
		traceCtx, span := tracer.Start(ctx.Context(), "get_test_scenario_by_id")
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

		span.SetAttributes(attribute.String("test_scenario.id", strconv.FormatUint(id, 10)))

		svcResult, err := handler.testScenario.GetByID(traceCtx, id)
		if err != nil {
			span.SetAttributes(attribute.String("error.type", "get_by_id_error"), attribute.String("error.message", err.Error()))
			return pkg.ToHTTPError(err).AsFiber(ctx)
		}

		result := response.TestScenario{
			ID:                  svcResult.ID,
			CreatedAt:           svcResult.CreatedAt,
			UpdatedAt:           svcResult.UpdatedAt,
			StartedAt:           svcResult.StartedAt,
			Name:                svcResult.Name,
			Status:              svcResult.Status,
			MaxTestServiceCount: svcResult.MaxTestServiceCount,
			ExecutionDuration:   svcResult.ExecutionDuration,
			AutoStepChangeRate:  svcResult.AutoStepChangeRate,
			TestCategory: func() *response.TestCategory {
				if svcResult.TestCategory == nil {
					return nil
				}

				return &response.TestCategory{
					ID:                     svcResult.TestCategory.ID,
					Name:                   svcResult.TestCategory.Name,
					Label:                  svcResult.TestCategory.Label,
					HasMaxTestServiceCount: svcResult.TestCategory.HasMaxTestServiceCount,
					HasExecutionDuration:   svcResult.TestCategory.HasExecutionDuration,
					HasAutoStepChangeRate:  svcResult.TestCategory.HasAutoStepChangeRate,
					CreatedAt:              svcResult.TestCategory.CreatedAt,
					UpdatedAt:              svcResult.TestCategory.UpdatedAt,
				}
			}(),
			MotherService: func() *response.MotherService {
				if svcResult.MotherService == nil {
					return nil
				}

				return &response.MotherService{
					ID:                       svcResult.MotherService.ID,
					CreatedAt:                svcResult.MotherService.CreatedAt,
					UpdatedAt:                svcResult.MotherService.UpdatedAt,
					Name:                     svcResult.MotherService.Name,
					ExceptionRate:            svcResult.MotherService.ExceptionRate,
					ResponseDelayRate:        svcResult.MotherService.ResponseDelayRate,
					ResponseDelayDuration:    svcResult.MotherService.ResponseDelayDuration,
					RandomResponseDelayMin:   svcResult.MotherService.RandomResponseDelayMin,
					RandomResponseDelayMax:   svcResult.MotherService.RandomResponseDelayMax,
					Status:                   svcResult.MotherService.Status,
					ServiceDeploymentAddress: svcResult.MotherService.ServiceDeploymentAddress,
					DatabaseName:             svcResult.MotherService.DatabaseName,
					DatabaseTableName:        svcResult.MotherService.DatabaseTableName,
				}
			}(),
			TestServiceConfig: func() *response.TestServiceConfig {
				if svcResult.TestServiceConfig == nil {
					return nil
				}

				return &response.TestServiceConfig{
					ID:                    svcResult.TestServiceConfig.ID,
					CreatedAt:             svcResult.TestServiceConfig.CreatedAt,
					UpdatedAt:             svcResult.TestServiceConfig.UpdatedAt,
					MaxRequests:           svcResult.TestServiceConfig.MaxRequests,
					MaxDuration:           svcResult.TestServiceConfig.MaxDuration,
					RequestDelayDuration:  svcResult.TestServiceConfig.RequestDelayDuration,
					RandomRequestDelayMin: svcResult.TestServiceConfig.RandomRequestDelayMin,
					RandomRequestDelayMax: svcResult.TestServiceConfig.RandomRequestDelayMax,
					FixedTestNumber:       svcResult.TestServiceConfig.FixedTestNumber,
					RandomTestNumberMin:   svcResult.TestServiceConfig.RandomTestNumberMin,
					RandomTestNumberMax:   svcResult.TestServiceConfig.RandomTestNumberMax,
					BadValueRate:          svcResult.TestServiceConfig.BadValueRate,
					NegativeValueRate:     svcResult.TestServiceConfig.NegativeValueRate,
					RealValueRate:         svcResult.TestServiceConfig.RealValueRate,
					ZeroValueRate:         svcResult.TestServiceConfig.ZeroValueRate,
					StringValueRate:       svcResult.TestServiceConfig.StringValueRate,
					LongStringValueRate:   svcResult.TestServiceConfig.LongStringValueRate,
					NullValueRate:         svcResult.TestServiceConfig.NullValueRate,
					DatabaseName:          svcResult.TestServiceConfig.DatabaseName,
					DatabaseTableName:     svcResult.TestServiceConfig.DatabaseTableName,
				}
			}(),
			Editable: svcResult.Editable,
		}

		return ctx.Status(http.StatusOK).JSON(&response.TestScenarioResponseByID{
			Data: result,
		})
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
func (handler *TestScenario) Start() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tracer := otel.Tracer("test-scenario-handler")
		traceCtx, span := tracer.Start(ctx.Context(), "start_test_scenario")
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

		span.SetAttributes(attribute.String("test_scenario.id", strconv.FormatUint(id, 10)))

		err = handler.testScenario.Start(traceCtx, id)
		if err != nil {
			span.SetAttributes(attribute.String("error.type", "start_error"), attribute.String("error.message", err.Error()))
			return pkg.ToHTTPError(err).AsFiber(ctx)
		}

		return ctx.Status(http.StatusOK).JSON(&response.SuccessResponse{
			Message: pkg.TestScenarioStarted,
		})
	}
}

// Update godoc
//
//	@Summary		Update a test scenario
//	@Description	Update test scenario with configuration
//	@Tags			test-scenarios
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.TestScenarioUpdateRequest	true	"Request body"
//	@Success		200		{object}	response.SuccessResponse
//	@Failure		400		{object}	response.ErrorResponse
//	@Failure		404		{object}	response.ErrorResponse
//	@Failure		422		{object}	response.ErrorResponse
//	@Failure		500		{object}	response.ErrorResponse
//	@Router			/api/v1/test-scenarios/update [post]
func (handler *TestScenario) Update() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		requestID := logger.GetRequestID(ctx.Context())
		tracer := otel.Tracer("test-scenario-handler")
		traceCtx, span := tracer.Start(ctx.Context(), "update_test_scenario")
		defer span.End()

		span.SetAttributes(attribute.String("request_id", requestID))

		req := new(request.TestScenarioUpdateRequest)

		if err := ctx.BodyParser(req); err != nil {
			span.SetAttributes(attribute.String("error.type", "bad_request"))
			return pkg.ToHTTPError(pkg.ErrBadRequest).AsFiber(ctx)
		}

		span.SetAttributes(
			attribute.String("test_scenario.id", strconv.FormatUint(req.ID, 10)),
			attribute.String("test_scenario.name", req.Name),
			attribute.String("mother_service.id", strconv.FormatUint(req.MotherServiceID, 10)))

		err := handler.testScenario.Update(traceCtx, req)
		if err != nil {
			span.SetAttributes(attribute.String("error.type", "update_error"), attribute.String("error.message", err.Error()))
			return pkg.ToHTTPError(err).AsFiber(ctx)
		}

		return ctx.Status(http.StatusOK).JSON(&response.SuccessResponse{
			Message: pkg.UpdateTestScenarioSuccessfully,
		})
	}
}
