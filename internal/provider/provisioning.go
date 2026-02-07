package provider

import (
	"context"
	"control-panel-service/internal/domain/entity"
)

type ProvisioningService interface {
	ProvisionTestService(ctx context.Context, testServiceConfig *entity.TestServiceConfig, count int) error
	DeprovisionTestService(ctx context.Context, ids []uint64) error
}
