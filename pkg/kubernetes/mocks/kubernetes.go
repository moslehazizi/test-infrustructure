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

func (kuberneteseMock *KuberneteseMock) ApplyService(ctx context.Context, spec entity.ServiceSpec) error {
	args := kuberneteseMock.Called(ctx, spec)

	if args.Error(0) != nil {
		return fmt.Errorf("%w", args.Error(0))
	}

	return nil
}

func (kuberneteseMock *KuberneteseMock) ApplyIngress(ctx context.Context, spec entity.IngressSpec) error {
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

func (kuberneteseMock *KuberneteseMock) GetDeploymentReplicas(ctx context.Context, name string) (int32, error) {
	args := kuberneteseMock.Called(ctx, name)
	return int32(args.Int(0)), args.Error(1) // #nosec G115 -- mock return
}

func (kuberneteseMock *KuberneteseMock) ScaleDeployment(ctx context.Context, name string, replicas int32) error {
	args := kuberneteseMock.Called(ctx, name, replicas)
	return args.Error(0)
}

func (kuberneteseMock *KuberneteseMock) DeleteDeployment(ctx context.Context, name string) error {
	args := kuberneteseMock.Called(ctx, name)
	return args.Error(0)
}

func (kuberneteseMock *KuberneteseMock) DeleteService(ctx context.Context, name string) error {
	args := kuberneteseMock.Called(ctx, name)
	return args.Error(0)
}

func (kuberneteseMock *KuberneteseMock) DeleteIngress(ctx context.Context, name string) error {
	args := kuberneteseMock.Called(ctx, name)
	return args.Error(0)
}

func (kuberneteseMock *KuberneteseMock) DeleteConfigMap(ctx context.Context, name string) error {
	args := kuberneteseMock.Called(ctx, name)
	return args.Error(0)
}

func (kuberneteseMock *KuberneteseMock) DeleteSecret(ctx context.Context, name string) error {
	args := kuberneteseMock.Called(ctx, name)
	return args.Error(0)
}

func (kuberneteseMock *KuberneteseMock) Client() *kubernetes.Clientset {
	args := kuberneteseMock.Called()

	result := args.Get(0).(*kubernetes.Clientset)

	return result
}
