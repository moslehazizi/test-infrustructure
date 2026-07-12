package handler

// import (
// 	"control-panel-service/internal/domain/entity"
// 	"control-panel-service/internal/server/dto/request"
// 	"control-panel-service/internal/server/dto/response"
// 	"control-panel-service/pkg"
// 	"control-panel-service/pkg/logger"
// 	"net/http"
// 	"strconv"

// 	"github.com/gofiber/fiber/v2"
// 	"go.opentelemetry.io/otel"
// 	"go.opentelemetry.io/otel/attribute"
// )

// type TestScenarioHandler struct {
// 	testScenario TestScenario
// }

// func NewTestScenarioHandler(testScenario TestScenario) *TestScenarioHandler {
// 	return &TestScenarioHandler{
// 		testScenario: testScenario,
// 	}
// }

// // Create godoc
// //
// //	@Summary		Create a test scenario
// //	@Description	Create a new test scenario with configuration
// //	@Tags			test-scenarios
// //	@Accept			json
// //	@Produce		json
// //	@Param			request	body		request.TestScenario	true	"Request body"
// //	@Success		200		{object}	response.SuccessResponse
// //	@Failure		400		{object}	response.ErrorResponse
// //	@Failure		422		{object}	response.ErrorResponse
// //	@Failure		500		{object}	response.ErrorResponse
// //	@Router			/api/v1/test-scenarios [post]
// func (handler *TestScenarioHandler) Create() fiber.Handler {
// 	return func(ctx *fiber.Ctx) error {
// 		requestID := logger.GetRequestID(ctx.Context())
// 		tracer := otel.Tracer("test-scenario-handler")
// 		traceCtx, span := tracer.Start(ctx.Context(), "create-test-scenario-handler")
// 		defer span.End()

// 		span.SetAttributes(attribute.String("request_id", requestID))

// 		req := new(request.TestScenario)

// 		if err := ctx.BodyParser(req); err != nil {
// 			span.SetAttributes(attribute.String("error.type", "bad_request"))
// 			return pkg.ToHTTPError(pkg.ErrBadRequest).AsFiber(ctx)
// 		}

// 		span.SetAttributes(
// 			attribute.String("test_scenario.name", req.Name),
// 			attribute.String("test_category.id", strconv.FormatUint(req.TestCategoryID, 10)),
// 			attribute.String("mother_service.id", strconv.FormatUint(req.MotherServiceID, 10)),
// 		)

// 		testScenario := new(entity.TestScenario)
// 		req.ToTestScenarioEntity(testScenario)

// 		err := handler.testScenario.Create(traceCtx, testScenario)
// 		if err != nil {
// 			span.SetAttributes(attribute.String("error.type", "create_error"), attribute.String("error.message", err.Error()))
// 			return pkg.ToHTTPError(err).AsFiber(ctx)
// 		}

// 		return ctx.Status(http.StatusOK).JSON(&response.SuccessResponse{
// 			Message: pkg.CreateTestScenarioSuccessfully,
// 		})
// 	}
// }

// // GetPaginated godoc
// //
// //	@Summary		Get paginated test scenarios
// //	@Description	Get test scenarios with pagination support
// //	@Tags			test-scenarios
// //	@Accept			json
// //	@Produce		json
// //	@Param			request	body		request.TestScenarioPaginationRequest	true	"Pagination request with page and per_page"
// //	@Success		200		{object}	response.PaginatedTestScenario
// //	@Failure		400		{object}	response.ErrorResponse
// //	@Failure		422		{object}	response.ErrorResponse
// //	@Failure		500		{object}	response.ErrorResponse
// //	@Router			/api/v1/test-scenarios/search [post]
// func (handler *TestScenarioHandler) GetPaginated() fiber.Handler {
// 	return func(ctx *fiber.Ctx) error {
// 		tracer := otel.Tracer("test-scenario-handler")
// 		traceCtx, span := tracer.Start(ctx.Context(), "get-paginated-test-scenarios-handler")
// 		defer span.End()

// 		requestID := logger.GetRequestID(ctx.Context())
// 		span.SetAttributes(attribute.String("request_id", requestID))

// 		req := new(request.TestScenarioPaginationRequest)
// 		if err := ctx.BodyParser(req); err != nil {
// 			span.SetAttributes(attribute.String("error.type", "bad_request"))
// 			return pkg.ToHTTPError(pkg.ErrBadRequest).AsFiber(ctx)
// 		}

// 		span.SetAttributes(attribute.String("pagination.page", strconv.Itoa(req.Page)), attribute.String("pagination.per_page", strconv.Itoa(req.PerPage)))

// 		items, count, err := handler.testScenario.GetPaginated(traceCtx, entity.TestScenarioPaginationRequest{
// 			Page:    req.Page,
// 			PerPage: req.PerPage,
// 		})
// 		if err != nil {
// 			span.SetAttributes(attribute.String("error.type", "get_paginated_error"), attribute.String("error.message", err.Error()))
// 			return pkg.ToHTTPError(err).AsFiber(ctx)
// 		}

// 		var responses []response.TestScenario
// 		for _, item := range items {
// 			response := response.TestScenario{}
// 			response.FromTestScenarioEntity(item)
// 			responses = append(responses, response)
// 		}

// 		return ctx.Status(http.StatusOK).JSON(&response.PaginatedTestScenario{
// 			Page:    req.Page,
// 			PerPage: req.PerPage,
// 			Data:    responses,
// 			Total:   count,
// 		})
// 	}
// }

// // GetByID godoc
// //
// //	@Summary		Get test scenario by ID
// //	@Description	Retrieve a specific test scenario by its ID
// //	@Tags			test-scenarios
// //	@Accept			json
// //	@Produce		json
// //	@Param			id	path		int	true	"Test scenario ID"
// //	@Success		200	{object}	response.TestScenarioResponseByID
// //	@Failure		400	{object}	response.ErrorResponse
// //	@Failure		404	{object}	response.ErrorResponse
// //	@Failure		500	{object}	response.ErrorResponse
// //	@Router			/api/v1/test-scenarios/{id} [get]
// func (handler *TestScenarioHandler) GetByID() fiber.Handler {
// 	return func(ctx *fiber.Ctx) error {
// 		tracer := otel.Tracer("test-scenario-handler")
// 		traceCtx, span := tracer.Start(ctx.Context(), "get-test-scenario-by-id-handler")
// 		defer span.End()

// 		requestID := logger.GetRequestID(ctx.Context())
// 		span.SetAttributes(attribute.String("request_id", requestID))

// 		idParam := ctx.Params("id")
// 		if idParam == "" {
// 			span.SetAttributes(attribute.String("error.type", "missing_id"))
// 			return pkg.ToHTTPError(pkg.ErrPageNotFound).AsFiber(ctx)
// 		}

// 		id, err := strconv.ParseUint(idParam, 10, 64)
// 		if err != nil {
// 			span.SetAttributes(attribute.String("error.type", "invalid_id_in_params"))
// 			return pkg.ToHTTPError(pkg.ErrInvalidIDInParams).AsFiber(ctx)
// 		}

// 		span.SetAttributes(attribute.String("test_scenario.id", strconv.FormatUint(id, 10)))

// 		svcResult, err := handler.testScenario.GetByID(traceCtx, id)
// 		if err != nil {
// 			span.SetAttributes(attribute.String("error.type", "get_by_id_error"), attribute.String("error.message", err.Error()))
// 			return pkg.ToHTTPError(err).AsFiber(ctx)
// 		}

// 		result := response.TestScenario{}
// 		result.FromTestScenarioEntity(svcResult)

// 		return ctx.Status(http.StatusOK).JSON(&response.TestScenarioResponseByID{
// 			Data: result,
// 		})
// 	}
// }

// // Update godoc
// //
// //	@Summary		Update a test scenario
// //	@Description	Update test scenario with configuration
// //	@Tags			test-scenarios
// //	@Accept			json
// //	@Produce		json
// //	@Param			request	body		request.TestScenarioUpdateRequest	true	"Request body"
// //	@Success		200		{object}	response.SuccessResponse
// //	@Failure		400		{object}	response.ErrorResponse
// //	@Failure		404		{object}	response.ErrorResponse
// //	@Failure		422		{object}	response.ErrorResponse
// //	@Failure		500		{object}	response.ErrorResponse
// //	@Router			/api/v1/test-scenarios/update [post]
// func (handler *TestScenarioHandler) Update() fiber.Handler {
// 	return func(ctx *fiber.Ctx) error {
// 		requestID := logger.GetRequestID(ctx.Context())
// 		tracer := otel.Tracer("test-scenario-handler")
// 		traceCtx, span := tracer.Start(ctx.Context(), "update-test-scenario-handler")
// 		defer span.End()

// 		span.SetAttributes(attribute.String("request_id", requestID))

// 		req := new(request.TestScenarioUpdateRequest)

// 		if err := ctx.BodyParser(req); err != nil {
// 			span.SetAttributes(attribute.String("error.type", "bad_request"))
// 			return pkg.ToHTTPError(pkg.ErrBadRequest).AsFiber(ctx)
// 		}

// 		span.SetAttributes(
// 			attribute.String("test_scenario.id", strconv.FormatUint(req.ID, 10)),
// 			attribute.String("test_scenario.name", req.Name),
// 			attribute.String("mother_service.id", strconv.FormatUint(req.MotherServiceID, 10)))

// 		err := handler.testScenario.Update(traceCtx, req)
// 		if err != nil {
// 			span.SetAttributes(attribute.String("error.type", "update_error"), attribute.String("error.message", err.Error()))
// 			return pkg.ToHTTPError(err).AsFiber(ctx)
// 		}

// 		return ctx.Status(http.StatusOK).JSON(&response.SuccessResponse{
// 			Message: pkg.UpdateTestScenarioSuccessfully,
// 		})
// 	}
// }
