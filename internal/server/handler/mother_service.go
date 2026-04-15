package handler

import (
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/server/dto/request"
	"control-panel-service/internal/server/dto/response"
	"control-panel-service/internal/usecase"
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

type MotherService struct {
	motherService usecase.MotherService
}

func NewMotherServiceHandler(motherService usecase.MotherService) *MotherService {
	return &MotherService{
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
func (handler *MotherService) Create() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tracer := otel.Tracer("mother-service-handler")
		traceCtx, span := tracer.Start(ctx.Context(), "create_mother_service")
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

		reqService := &entity.MotherService{
			Name:              req.Name,
			ExceptionRate:     req.ExceptionRate,
			ResponseDelayRate: req.ResponseDelayRate,
			ResponseDelayDuration: func() *int {
				if req.ResponseDelayDuration != nil {
					return req.ResponseDelayDuration
				}

				return nil
			}(),
			RandomResponseDelayMin: func() *int {
				if req.RandomResponseDelayMin != nil {
					return req.RandomResponseDelayMin
				}

				return nil
			}(),
			RandomResponseDelayMax: func() *int {
				if req.RandomResponseDelayMax != nil {
					return req.RandomResponseDelayMax
				}

				return nil
			}(),
			DatabaseName:      req.DatabaseName,
			DatabaseTableName: req.DatabaseTableName,
		}

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
func (handler *MotherService) GetByID() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tracer := otel.Tracer("mother-service-handler")
		traceCtx, span := tracer.Start(ctx.Context(), "get_mother_service-by-id")
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

		result := response.MotherService{
			ID:                       svcResult.ID,
			CreatedAt:                svcResult.CreatedAt,
			UpdatedAt:                svcResult.UpdatedAt,
			Name:                     svcResult.Name,
			ExceptionRate:            svcResult.ExceptionRate,
			ResponseDelayRate:        svcResult.ResponseDelayRate,
			ResponseDelayDuration:    svcResult.ResponseDelayDuration,
			RandomResponseDelayMin:   svcResult.RandomResponseDelayMin,
			RandomResponseDelayMax:   svcResult.RandomResponseDelayMax,
			Status:                   svcResult.Status,
			ServiceDeploymentAddress: svcResult.ServiceDeploymentAddress,
			DatabaseName:             svcResult.DatabaseName,
			DatabaseTableName:        svcResult.DatabaseTableName,
		}

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
func (handler *MotherService) GetPaginated() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tracer := otel.Tracer("mother-service-handler")
		traceCtx, span := tracer.Start(ctx.Context(), "get_paginated_mother_services")
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
			responses = append(responses, response.MotherService{
				ID:                       svcResult.ID,
				CreatedAt:                svcResult.CreatedAt,
				UpdatedAt:                svcResult.UpdatedAt,
				Name:                     svcResult.Name,
				ExceptionRate:            svcResult.ExceptionRate,
				ResponseDelayRate:        svcResult.ResponseDelayRate,
				ResponseDelayDuration:    svcResult.ResponseDelayDuration,
				RandomResponseDelayMin:   svcResult.RandomResponseDelayMin,
				RandomResponseDelayMax:   svcResult.RandomResponseDelayMax,
				Status:                   svcResult.Status,
				ServiceDeploymentAddress: svcResult.ServiceDeploymentAddress,
				DatabaseName:             svcResult.DatabaseName,
				DatabaseTableName:        svcResult.DatabaseTableName,
			})
		}

		return ctx.Status(http.StatusOK).JSON(&response.PaginatedMotherServices{
			Data:    responses,
			Page:    req.Page,
			PerPage: req.PerPage,
			Total:   count,
		})
	}
}

// Abort godoc
//
//	@Summary		Abort a mother service.
//	@Description	Abort a specific mother service by its ID.
//	@Tags			mother-services
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Test scenario ID"
//	@Success		200	{object}	response.SuccessResponse
//	@Failure		400	{object}	response.ErrorResponse
//	@Failure		404	{object}	response.ErrorResponse
//	@Failure		500	{object}	response.ErrorResponse
//	@Router			/api/v1/mother-services/{id}/abort [post]
func (handler *MotherService) Abort() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tracer := otel.Tracer("mother-service-handler")
		traceCtx, span := tracer.Start(ctx.Context(), "abort_mother_service")
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

		err = handler.motherService.Abort(traceCtx, id)
		if err != nil {
			span.SetAttributes(attribute.String("error.type", "abort_error"), attribute.String("error.message", err.Error()))
			return pkg.ToHTTPError(err).AsFiber(ctx)
		}

		return ctx.Status(http.StatusOK).JSON(&response.SuccessResponse{
			Message: pkg.MotherServiceAbort,
		})
	}
}
