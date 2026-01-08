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
			return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"error": pkg.InvalidReqBody,
			})
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
			KafkaLiveFeedTopic:       req.KafkaLiveFeedTopic,
			KafkaFactorialTopic:      req.KafkaFactorialTopic,
			ProvisioningStatus:       entity.ProvisioningStatusPending,
		}

		err := handler.motherService.Create(ctx.Context(), reqService)

		if err != nil {
			switch {
			case errors.Is(err, pkg.ErrMotherServiceAlreadyExist):
				return ctx.Status(http.StatusConflict).JSON(&fiber.Map{
					"error": pkg.MotherServiceAlreadyExist,
				})
			case errors.Is(err, pkg.ErrInvalidResponseDelayRate):
				return ctx.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
					"error": pkg.InvalidResponseDelayRate,
				})
			case errors.Is(err, pkg.ErrInvalidExceptionRate):
				return ctx.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
					"error": pkg.InvalidExceptionRate,
				})
			case errors.Is(err, pkg.ErrInvalidDelayConfiguration):
				return ctx.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
					"error": pkg.InvalidDelayConfiguration,
				})
			case errors.Is(err, pkg.ErrInvalidRandomDelayRange):
				return ctx.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
					"error": pkg.InvalidRandomDelayRange,
				})

			}

			return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"error": pkg.InternalServerErrorMessage,
			})
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
			return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"error": pkg.InvalidIDInParams,
			})
		}

		svcResult, err := handler.motherService.GetByID(ctx.Context(), id)
		if err != nil {
			if errors.Is(err, pkg.ErrMotherServiceNotFound) {
				return ctx.Status(http.StatusNotFound).JSON(&fiber.Map{
					"error": pkg.MotherServiceNotFound,
				})
			}

			return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"error": pkg.InternalServerErrorMessage,
			})
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
			ProvisioningStatus:       string(svcResult.ProvisioningStatus),
			ServiceDeploymentAddress: svcResult.ServiceDeploymentAddress,
			DatabaseName:             svcResult.DatabaseName,
			DatabaseTableName:        svcResult.DatabaseTableName,
			KafkaLiveFeedTopic:       svcResult.KafkaLiveFeedTopic,
			KafkaFactorialTopic:      svcResult.KafkaFactorialTopic,
			StoppedAt:                svcResult.StoppedAt,
			RestartedAt:              svcResult.RestartedAt,
			StartedAt:                svcResult.StartedAt,
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
			return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"error": pkg.InvalidReqBody,
			})
		}

		if req.Page < 0 || req.PerPage < 0 {
			return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"error": pkg.InvalidReqBody,
			})
		}

		reqSvc := entity.PaginationRequest{
			Page:    req.Page,
			PerPage: req.PerPage,
		}

		svcResults, err := handler.motherService.GetPaginated(ctx.Context(), reqSvc)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"error": pkg.InternalServerErrorMessage,
			})
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
				ProvisioningStatus:       string(svcResult.ProvisioningStatus),
				ServiceDeploymentAddress: svcResult.ServiceDeploymentAddress,
				DatabaseName:             svcResult.DatabaseName,
				DatabaseTableName:        svcResult.DatabaseTableName,
				KafkaLiveFeedTopic:       svcResult.KafkaLiveFeedTopic,
				KafkaFactorialTopic:      svcResult.KafkaFactorialTopic,
				StoppedAt:                svcResult.StoppedAt,
				RestartedAt:              svcResult.RestartedAt,
				StartedAt:                svcResult.StartedAt,
			})
		}

		return ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"data":     responses,
			"page":     req.Page,
			"per_page": req.PerPage,
		})
	}
}
