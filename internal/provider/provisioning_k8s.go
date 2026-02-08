package provider

import (
	"context"
	"control-panel-service/config"
	"control-panel-service/internal/domain/entity"
)

func NewProvisioningService(cfg *config.Config) ProvisioningService {
	return &provisioningService{
		cfg: cfg,
	}
}

type provisioningService struct {
	cfg *config.Config
}

func (ps *provisioningService) ProvisionTestService(ctx context.Context, testScenario *entity.TestScenario, replica int32) error {
	return nil
}

func (ps *provisioningService) DeprovisionTestService(ctx context.Context, testScenario *entity.TestScenario, replica int32) error {
	return nil
}

func (ps *provisioningService) ProvisionMotherService(ctx context.Context, motherService *entity.MotherService) error {
	return nil
}

func (ps *provisioningService) DeprovisionMotherService(ctx context.Context, motherService *entity.MotherService) error {
	return nil
}
