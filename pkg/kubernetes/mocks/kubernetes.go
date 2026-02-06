package mocks

import (
	"context"
	"control-panel-service/pkg/kubernetes/domain/entity"
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"
	"k8s.io/client-go/kubernetes"
)

type KuberneteseMock struct {
	mock.Mock
}

func (kuberneteseMock *KuberneteseMock) ApplyDeployment(ctx context.Context, spec entity.DeploymentSpec, ConfigMap, Secret map[string]string) error {
	args := kuberneteseMock.Called(ctx, spec, ConfigMap, Secret)

	if args.Error(0) != nil {
		return fmt.Errorf("%w", args.Error(0))
	}

	return nil
}

func (kuberneteseMock *KuberneteseMock) ApplyService(ctx context.Context, spec entity.ServiceSpec, ConfigMap, Secret map[string]string) error {
	args := kuberneteseMock.Called(ctx, spec, ConfigMap, Secret)

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

func (kuberneteseMock *KuberneteseMock) Client() *kubernetes.Clientset {
	args := kuberneteseMock.Called()

	result := args.Get(0).(*kubernetes.Clientset)

	return result
}
