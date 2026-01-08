package handler

import (
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/server/dto/request"
	"control-panel-service/internal/usecase"
	"control-panel-service/pkg"
	"errors"
	"net/http"
	"time"

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
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			ResponseDelayRate: req.ResponseDelayRate,
			ResponseDelayDuration: func() *int {
				if req.ResponseDelayDuration != nil {
					return req.ResponseDelayDuration
				}
				zero := 0

				return &zero
			}(),
			RandomResponseDelayMin: func() *int {
				if req.RandomResponseDelayMin != nil {
					return req.RandomResponseDelayMin
				}
				zero := 0

				return &zero
			}(),
			RandomResponseDelayMax: func() *int {
				if req.RandomResponseDelayMax != nil {
					return req.RandomResponseDelayMax
				}
				zero := 0

				return &zero
			}(),
			ServiceDeploymentAddress: req.ServiceDeploymentAddress,
			DatabaseName:             req.DatabaseName,
			DatabaseTableName:        req.DatabaseTableName,
			KafkaLiveFeedTopic:       req.KafkaLiveFeedTopic,
			KafkaFactorialTopic:      req.KafkaFactorialTopic,
			ProvisioningStatus:       entity.ProvisioningStatus(req.ProvisioningStatus),
		}

		err := handler.motherService.Create(ctx.Context(), reqService)

		if err != nil {
			if errors.Is(err, pkg.ErrMotherServiceAlreadyExist) {
				return ctx.Status(http.StatusConflict).JSON(&fiber.Map{
					"error": pkg.MotherServiceAlreadyExist,
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
