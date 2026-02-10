package kubernetese

import (
	"context"
	"control-panel-service/pkg/kubernetes/domain/entity"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	Config = "-config"
	Secret = "-secret"
	Wait   = 2 * time.Second
)

var (
	ErrDeploymentNotReady = errors.New("deployment did not become ready within timeout")
)

type Kubernetese interface {
	ApplyDeployment(ctx context.Context, spec entity.DeploymentSpec, configMap, Secret map[string]string) error
	GetDeploymentReplicas(ctx context.Context, name string) (int32, error)
	ScaleDeployment(ctx context.Context, name string, replicas int32) error
	DeleteDeployment(ctx context.Context, name string) error
	ApplyService(ctx context.Context, spec entity.ServiceSpec, configMap, Secret map[string]string) error
	DeleteService(ctx context.Context, name string) error
	DeleteConfigMap(ctx context.Context, name string) error
	DeleteSecret(ctx context.Context, name string) error
	WaitForDeployment(ctx context.Context, name string, timeout time.Duration) error
	Client() *kubernetes.Clientset
}

type KubernConfig struct {
	NameSpace  string
	kubeConfig string
}

type Kuber struct {
	cfg       *KubernConfig
	Clientset *kubernetes.Clientset
}

func New(ctx context.Context, config *KubernConfig) (Kubernetese, error) {
	if config.kubeConfig == "" {
		config.kubeConfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}

	cfg, err := clientcmd.BuildConfigFromFlags("", config.kubeConfig)
	if err != nil {
		// return nil, fmt.Errorf("failed to build kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		// return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	return &Kuber{
		cfg:       config,
		Clientset: clientset,
	}, nil
}

// applyConfigMap creates or updates the configMap.
func (k *Kuber) applyConfigMap(ctx context.Context, configMap map[string]string, app string) error {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: app + Config,
		},
		Data: configMap,
	}

	_, err := k.Clientset.CoreV1().ConfigMaps(k.cfg.NameSpace).Get(ctx, cm.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		zap.L().Info("creating configmap ...", zap.String("apllication", app))
		_, err = k.Clientset.CoreV1().ConfigMaps(k.cfg.NameSpace).Create(ctx, cm, metav1.CreateOptions{})

		return err
	}
	if err != nil {
		return err
	}

	zap.L().Info("updating configmap ...", zap.String("apllication", app))
	_, err = k.Clientset.CoreV1().ConfigMaps(k.cfg.NameSpace).Update(ctx, cm, metav1.UpdateOptions{})

	return err
}

// applySecret creates or updates the Secret.
func (k *Kuber) applySecret(ctx context.Context, secretMap map[string]string, app string) error {
	sec := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: app + Secret,
		},
		StringData: secretMap,
	}

	_, err := k.Clientset.CoreV1().Secrets(k.cfg.NameSpace).Get(ctx, sec.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		zap.L().Info("creating secret ...", zap.String("apllication", app))
		_, err = k.Clientset.CoreV1().Secrets(k.cfg.NameSpace).Create(ctx, sec, metav1.CreateOptions{})

		return err
	}
	if err != nil {
		return err
	}

	zap.L().Info("updating secret ...", zap.String("apllication", app))
	_, err = k.Clientset.CoreV1().Secrets(k.cfg.NameSpace).Update(ctx, sec, metav1.UpdateOptions{})

	return err
}

// applyDeployment creates a deployment if it doesn't exist, otherwise updates it.
func (k *Kuber) ApplyDeployment(ctx context.Context, spec entity.DeploymentSpec, configMap, secretMap map[string]string) error {
	err := k.applyConfigMap(ctx, configMap, spec.Name)
	if err != nil {
		return err
	}

	err = k.applySecret(ctx, secretMap, spec.Name)
	if err != nil {
		return err
	}

	existing, err := k.Clientset.AppsV1().Deployments(k.cfg.NameSpace).Get(ctx, spec.Name, metav1.GetOptions{})

	if apierrors.IsNotFound(err) {
		zap.L().Info("creating deployment ...", zap.String("apllication", spec.Name))
		_, err = k.Clientset.AppsV1().Deployments(k.cfg.NameSpace).Create(ctx, spec.Deployment, metav1.CreateOptions{})

		return err
	}

	if err != nil {
		return err
	}

	zap.L().Info("updating deployment ...", zap.String("apllication", spec.Name))
	spec.Deployment.ResourceVersion = existing.ResourceVersion
	_, err = k.Clientset.AppsV1().Deployments(k.cfg.NameSpace).Update(ctx, spec.Deployment, metav1.UpdateOptions{})

	return err
}

// applyService creates a service if it doesn't exist, otherwise updates it.
func (k *Kuber) ApplyService(ctx context.Context, spec entity.ServiceSpec, configMap, secretMap map[string]string) error {
	err := k.applyConfigMap(ctx, configMap, spec.Name)
	if err != nil {
		return err
	}

	err = k.applySecret(ctx, secretMap, spec.Name)
	if err != nil {
		return err
	}

	existing, err := k.Clientset.CoreV1().Services(k.cfg.NameSpace).Get(ctx, spec.Name, metav1.GetOptions{})

	if apierrors.IsNotFound(err) {
		zap.L().Info("creating service ...", zap.String("apllication", spec.Name))
		_, err = k.Clientset.CoreV1().Services(k.cfg.NameSpace).Create(ctx, spec.Service, metav1.CreateOptions{})

		return err
	}

	if err != nil {
		return err
	}

	zap.L().Info("updating service ...", zap.String("apllication", spec.Name))
	spec.Service.ResourceVersion = existing.ResourceVersion
	spec.Service.Spec.ClusterIP = existing.Spec.ClusterIP
	_, err = k.Clientset.CoreV1().Services(k.cfg.NameSpace).Update(ctx, spec.Service, metav1.UpdateOptions{})

	return err
}

// waitForDeploymentReady waits for a deployment to become ready within the timeout.
func (k *Kuber) WaitForDeployment(ctx context.Context, name string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		dep, err := k.Clientset.AppsV1().Deployments(k.cfg.NameSpace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		zap.L().Info("waiting ...", zap.String("apllication", name), zap.Int32("ready_replica", dep.Status.ReadyReplicas), zap.Int32("spec_replica", *dep.Spec.Replicas))

		if dep.Status.ReadyReplicas > 0 && dep.Status.ReadyReplicas == *dep.Spec.Replicas {
			zap.L().Info("ready!", zap.String("apllication", name))

			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("%w", ctx.Err())
		case <-time.After(Wait):
		}
	}

	return fmt.Errorf("%w: name=%s timeout=%v", ErrDeploymentNotReady, name, timeout)
}

// GetDeploymentReplicas returns the current desired replicas of a deployment.
func (k *Kuber) GetDeploymentReplicas(ctx context.Context, name string) (int32, error) {
	dep, err := k.Clientset.AppsV1().Deployments(k.cfg.NameSpace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return 0, err
	}
	if dep.Spec.Replicas == nil {
		return 0, nil
	}

	return *dep.Spec.Replicas, nil
}

// ScaleDeployment sets the number of replicas for a deployment.
func (k *Kuber) ScaleDeployment(ctx context.Context, name string, replicas int32) error {
	dep, err := k.Clientset.AppsV1().Deployments(k.cfg.NameSpace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	dep.Spec.Replicas = &replicas
	_, err = k.Clientset.AppsV1().Deployments(k.cfg.NameSpace).Update(ctx, dep, metav1.UpdateOptions{})

	return err
}

// DeleteDeployment deletes a deployment (pods are removed by cascade).
func (k *Kuber) DeleteDeployment(ctx context.Context, name string) error {
	err := k.Clientset.AppsV1().Deployments(k.cfg.NameSpace).Delete(ctx, name, metav1.DeleteOptions{})
	if apierrors.IsNotFound(err) {
		return nil
	}

	return err
}

// DeleteService deletes a service.
func (k *Kuber) DeleteService(ctx context.Context, name string) error {
	err := k.Clientset.CoreV1().Services(k.cfg.NameSpace).Delete(ctx, name, metav1.DeleteOptions{})
	if apierrors.IsNotFound(err) {
		return nil
	}

	return err
}

// DeleteConfigMap deletes a configmap.
func (k *Kuber) DeleteConfigMap(ctx context.Context, name string) error {
	err := k.Clientset.CoreV1().ConfigMaps(k.cfg.NameSpace).Delete(ctx, name, metav1.DeleteOptions{})
	if apierrors.IsNotFound(err) {
		return nil
	}

	return err
}

// DeleteSecret deletes a secret.
func (k *Kuber) DeleteSecret(ctx context.Context, name string) error {
	err := k.Clientset.CoreV1().Secrets(k.cfg.NameSpace).Delete(ctx, name, metav1.DeleteOptions{})
	if apierrors.IsNotFound(err) {
		return nil
	}

	return err
}

func (k *Kuber) Client() *kubernetes.Clientset {
	return k.Clientset
}
