package handler

import (
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/server/dto/request"
	"control-panel-service/internal/server/dto/response"
	"control-panel-service/pkg"
	"control-panel-service/pkg/logger"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type MotherServiceHandler struct {
	motherService MotherService
}

func NewMotherServiceHandler(motherService MotherService) *MotherServiceHandler {
	return &MotherServiceHandler{
		motherService: motherService,
	}
}

// Create godoc
//
//	@Summary		Create a mother service
//	@Description	Create a new mother service
//	@Tags			mother-services
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.MotherService	true	"Request body"
//	@Success		200		{object}	response.SuccessResponse
//	@Failure		400		{object}	response.ErrorResponse
//	@Failure		422		{object}	response.ErrorResponse
//	@Failure		500		{object}	response.ErrorResponse
//	@Router			/api/v1/mother-services [post]
func (handler *MotherServiceHandler) Create() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tracer := otel.Tracer("mother-service-handler")
		traceCtx, span := tracer.Start(ctx.Context(), "create-mother-service-handler")
		defer span.End()

		requestID := logger.GetRequestID(ctx.Context())
		span.SetAttributes(attribute.String("request_id", requestID))

		req := new(request.MotherService)

		if err := ctx.BodyParser(req); err != nil {
			span.SetAttributes(attribute.String("error.type", "bad_request"))
			return pkg.ToHTTPError(pkg.ErrBadRequest).AsFiber(ctx)
		}

		span.SetAttributes(
			attribute.String("service.name", req.Name),
			attribute.Float64("service.exception_rate", float64(req.ExceptionRate)),
			attribute.Float64("service.response_delay_rate", float64(req.ResponseDelayRate)),
		)

		reqService := req.ToMotherServiceEntity()

		err := handler.motherService.Create(traceCtx, reqService)
		if err != nil {
			span.SetAttributes(
				attribute.String("error.type", "create_error"),
				attribute.String("error.message", err.Error()),
			)
			return pkg.ToHTTPError(err).AsFiber(ctx)
		}

		span.SetAttributes(attribute.String("status", "success"))

		return ctx.Status(http.StatusOK).JSON(&response.SuccessResponse{
			Message: pkg.CreateMotherServiceSuccessfully,
		})
	}
}

// GetByID godoc
//
//	@Summary		Get mother service by ID
//	@Description	Retrieve a specific mother service by its ID
//	@Tags			mother-services
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Mother service ID"
//	@Success		200	{object}	response.MotherServiceResponseByID
//	@Failure		400	{object}	response.ErrorResponse
//	@Failure		404	{object}	response.ErrorResponse
//	@Failure		500	{object}	response.ErrorResponse
//	@Router			/api/v1/mother-services/{id} [get]
func (handler *MotherServiceHandler) GetByID() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tracer := otel.Tracer("mother-service-handler")
		traceCtx, span := tracer.Start(ctx.Context(), "get-mother-service-by-id-handler")
		defer span.End()

		requestID := logger.GetRequestID(ctx.Context())
		span.SetAttributes(attribute.String("request_id", requestID))

		strID := strings.TrimSpace(ctx.Params("id"))

		id, err := strconv.ParseUint(strID, 10, 64)
		if err != nil {
			span.SetAttributes(attribute.String("error.type", "invalid_id_in_params"))
			return pkg.ToHTTPError(pkg.ErrInvalidIDInParams).AsFiber(ctx)
		}

		span.SetAttributes(attribute.String("service.id", strconv.FormatUint(id, 10)))

		svcResult, err := handler.motherService.GetByID(traceCtx, id)
		if err != nil {
			if errors.Is(err, pkg.ErrMotherServiceNotFound) {
				return pkg.ToHTTPError(pkg.ErrMotherServiceNotFound).AsFiber(ctx)
			}

			return pkg.ToHTTPError(err).AsFiber(ctx)
		}
		var result response.MotherService
		result.FromMotherServiceEntity(svcResult)

		return ctx.Status(http.StatusOK).JSON(&response.MotherServiceResponseByID{
			Data: result,
		})
	}
}

// GetPaginated godoc
//
//	@Summary		Get paginated mother services
//	@Description	Get mother services with pagination support
//	@Tags			mother-services
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.PaginationRequest	true	"Pagination request with page and per_page"
//	@Success		200		{object}	response.PaginatedMotherServices
//	@Failure		400		{object}	response.ErrorResponse
//	@Failure		500		{object}	response.ErrorResponse
//	@Router			/api/v1/mother-services/search [post]
func (handler *MotherServiceHandler) GetPaginated() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tracer := otel.Tracer("mother-service-handler")
		traceCtx, span := tracer.Start(ctx.Context(), "get-paginated-mother-services-handler")
		defer span.End()

		requestID := logger.GetRequestID(ctx.Context())
		span.SetAttributes(attribute.String("request_id", requestID))

		req := new(request.PaginationRequest)

		if err := ctx.BodyParser(req); err != nil {
			span.SetAttributes(attribute.String("error.type", "bad_request"))

			return pkg.ToHTTPError(pkg.ErrBadRequest).AsFiber(ctx)
		}

		if req.Page < 0 || req.PerPage < 0 {
			span.SetAttributes(attribute.String("error.type", "invalid_pagination"))

			return pkg.ToHTTPError(pkg.ErrBadRequest).AsFiber(ctx)
		}

		span.SetAttributes(attribute.String("pagination.page", strconv.Itoa(req.Page)), attribute.String("pagination.per_page", strconv.Itoa(req.PerPage)))

		reqSvc := entity.PaginationRequest{
			Page:    req.Page,
			PerPage: req.PerPage,
		}

		svcResults, count, err := handler.motherService.GetPaginated(traceCtx, reqSvc)
		if err != nil {
			span.SetAttributes(attribute.String("error.type", "get_paginated_error"), attribute.String("error.message", err.Error()))

			return pkg.ToHTTPError(err).AsFiber(ctx)
		}

		var responses []response.MotherService
		for _, svcResult := range svcResults {
			var result response.MotherService
			result.FromMotherServiceEntity(svcResult)

			responses = append(responses, result)
		}

		return ctx.Status(http.StatusOK).JSON(&response.PaginatedMotherServices{
			Data:    responses,
			Page:    req.Page,
			PerPage: req.PerPage,
			Total:   count,
		})
	}
}

// Delete godoc
//
//	@Summary		Delete a mother service.
//	@Description	Delete a specific mother service by its ID.
//	@Tags			mother-services
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Test scenario ID"
//	@Success		200	{object}	response.SuccessResponse
//	@Failure		400	{object}	response.ErrorResponse
//	@Failure		404	{object}	response.ErrorResponse
//	@Failure		500	{object}	response.ErrorResponse
//	@Router			/api/v1/mother-services/{id}/delete [post]
func (handler *MotherServiceHandler) Delete() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tracer := otel.Tracer("mother-service-handler")
		traceCtx, span := tracer.Start(ctx.Context(), "delete-mother-service-handler")
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

		span.SetAttributes(attribute.String("mother_service.id", strconv.FormatUint(id, 10)))

		err = handler.motherService.Delete(traceCtx, id)
		if err != nil {
			span.SetAttributes(attribute.String("error.type", "delete_error"), attribute.String("error.message", err.Error()))
			return pkg.ToHTTPError(err).AsFiber(ctx)
		}

		return ctx.Status(http.StatusOK).JSON(&response.SuccessResponse{
			Message: pkg.MotherServiceDelete,
		})
	}
}
