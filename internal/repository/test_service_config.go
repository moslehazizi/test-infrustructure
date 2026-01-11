package repository

import (
	"context"
	"control-panel-service/internal/domain/entity"
)

type TestServiceConfigRepository interface {
	Create(ctx context.Context, testSvcCfg *entity.TestServiceConfig) error
}
