package handler

import (
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/server/dto/request"
	"control-panel-service/internal/server/dto/response"
	"control-panel-service/internal/usecase"
	"control-panel-service/pkg"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type MotherService struct {
	motherService usecase.MotherService
}

func NewMotherServiceHandler(motherService usecase.MotherService) *MotherService {
	return &MotherService{
		motherService: motherService,
	}
}

func (handler *MotherService) Create() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		req := new(request.MotherService)

		if err := ctx.BodyParser(req); err != nil {
			return pkg.ToHTTPError(pkg.ErrBadRequest).AsFiber(ctx)
		}

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
			ServiceDeploymentAddress: req.ServiceDeploymentAddress,
			DatabaseName:             req.DatabaseName,
			DatabaseTableName:        req.DatabaseTableName,
		}

		err := handler.motherService.Create(ctx.Context(), reqService)
		if err != nil {
			return pkg.ToHTTPError(err).AsFiber(ctx)
		}

		return ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"message": pkg.CreateMotherServiceSuccessfully,
		})
	}
}

func (handler *MotherService) GetByID() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		strID := strings.TrimSpace(ctx.Params("id"))

		id, err := strconv.ParseUint(strID, 10, 64)
		if err != nil {
			return pkg.ToHTTPError(pkg.ErrInvalidIDInParams).AsFiber(ctx)
		}

		svcResult, err := handler.motherService.GetByID(ctx.Context(), id)
		if err != nil {
			if errors.Is(err, pkg.ErrMotherServiceNotFound) {
				return pkg.ToHTTPError(pkg.ErrMotherServiceNotFound).AsFiber(ctx)
			}

			return pkg.ToHTTPError(err).AsFiber(ctx)
		}

		response := response.MotherService{
			ID:                       svcResult.ID,
			CreatedAt:                svcResult.CreatedAt,
			UpdatedAt:                svcResult.UpdatedAt,
			Name:                     svcResult.Name,
			ExceptionRate:            svcResult.ExceptionRate,
			ResponseDelayRate:        svcResult.ResponseDelayRate,
			ResponseDelayDuration:    svcResult.ResponseDelayDuration,
			RandomResponseDelayMin:   svcResult.RandomResponseDelayMin,
			RandomResponseDelayMax:   svcResult.RandomResponseDelayMax,
			Status:                   string(svcResult.Status),
			ServiceDeploymentAddress: svcResult.ServiceDeploymentAddress,
			DatabaseName:             svcResult.DatabaseName,
			DatabaseTableName:        svcResult.DatabaseTableName,
		}

		return ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"data": response,
		})
	}
}

func (handler *MotherService) GetPaginated() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		req := new(request.PaginationRequest)

		if err := ctx.BodyParser(req); err != nil {
			return pkg.ToHTTPError(pkg.ErrBadRequest).AsFiber(ctx)
		}

		if req.Page < 0 || req.PerPage < 0 {
			return pkg.ToHTTPError(pkg.ErrBadRequest).AsFiber(ctx)
		}

		reqSvc := entity.PaginationRequest{
			Page:    req.Page,
			PerPage: req.PerPage,
		}

		svcResults, err := handler.motherService.GetPaginated(ctx.Context(), reqSvc)
		if err != nil {
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
				Status:                   string(svcResult.Status),
				ServiceDeploymentAddress: svcResult.ServiceDeploymentAddress,
				DatabaseName:             svcResult.DatabaseName,
				DatabaseTableName:        svcResult.DatabaseTableName,
			})
		}

		return ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"data":     responses,
			"page":     req.Page,
			"per_page": req.PerPage,
		})
	}
}
