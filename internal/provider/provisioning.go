package provider

import (
	"context"
	"control-panel-service/internal/domain/entity"
)

type ProvisioningService interface {
	ProvisionTestService(ctx context.Context, testServiceConfig *entity.TestServiceConfig, count int) error
	// TODO: it should be aligned with Nikola
	DeprovisionTestService(ctx context.Context) error
}

func NewProvisioningService() ProvisioningService {
	return &k8sProvisioningService{}
}

type k8sProvisioningService struct {
}

func (p *k8sProvisioningService) ProvisionTestService(ctx context.Context, testServiceConfig *entity.TestServiceConfig, count int) error {
	// TODO: use Nikola's work
	return nil
}

func (p *k8sProvisioningService) DeprovisionTestService(ctx context.Context) error {
	// TODO: use Nikola's work
	return nil
}
