package mocks

import (
	"context"
	"control-panel-service/pkg/kubernetese/domain/entity"
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"
)

type KuberneteseMock struct {
	mock.Mock
}

func (kuberneteseMock *KuberneteseMock) ApplyDeployment(ctx context.Context, spec entity.DeploymentSpec) error {
	args := kuberneteseMock.Called(ctx, spec)

	if args.Error(0) != nil {
		return fmt.Errorf("%w", args.Error(0))
	}

	return nil
}

func (kuberneteseMock *KuberneteseMock) ApplyService(ctx context.Context, spec entity.ServiceSpec) error {
	args := kuberneteseMock.Called(ctx, spec)

	if args.Error(0) != nil {
		return fmt.Errorf("%w", args.Error(0))
	}

	return nil
}

func (kuberneteseMock *KuberneteseMock) WaitForDeployment(ctx context.Context, name string, timeout time.Duration) error {
	args := kuberneteseMock.Called(ctx, name, timeout)

	if args.Error(0) != nil {
		return fmt.Errorf("%w", args.Error(0))
	}

	return nil
}
