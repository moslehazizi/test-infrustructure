package kubernetese

import (
	"context"
	"control-panel-service/pkg/kubernetese/domain/entity"
	"fmt"
	"os"
	"path/filepath"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Kubernetese interface {
	ApplyDeployment(ctx context.Context, spec entity.DeploymentSpec) error
	ApplyService(ctx context.Context, spec entity.ServiceSpec) error
	WaitForDeployment(ctx context.Context, name string, timeout time.Duration) error
}

type KubernConfig struct {
	NameSpace  string
	kubeConfig string
	ConfigMap  map[string]string
	Secret     map[string]string
}

type Kuber struct {
	cfg       *KubernConfig
	Clientset *kubernetes.Clientset
}

func NewKubernetes(ctx context.Context, config *KubernConfig) (Kubernetese, error) {
	if config.kubeConfig == "" {
		config.kubeConfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}

	cfg, err := clientcmd.BuildConfigFromFlags("", config.kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	k := &Kuber{
		cfg:       config,
		Clientset: clientset,
	}

	// Apply ConfigMap
	if err := k.applyConfigMap(ctx); err != nil {
		return nil, fmt.Errorf("apply configmap: %w", err)
	}

	// Apply Secret
	if err := k.applySecret(ctx); err != nil {
		return nil, fmt.Errorf("apply secret: %w", err)
	}

	return k, nil
}

// applyConfigMap creates or updates the ConfigMap
func (k *Kuber) applyConfigMap(ctx context.Context) error {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: k.cfg.NameSpace + "-config",
		},
		Data: k.cfg.ConfigMap,
	}

	_, err := k.Clientset.CoreV1().ConfigMaps(k.cfg.NameSpace).Get(ctx, cm.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		fmt.Println("creating configmap ...")
		_, err = k.Clientset.CoreV1().ConfigMaps(k.cfg.NameSpace).Create(ctx, cm, metav1.CreateOptions{})
		return err
	}
	if err != nil {
		return err
	}

	fmt.Println("updating configmap ...")
	_, err = k.Clientset.CoreV1().ConfigMaps(k.cfg.NameSpace).Update(ctx, cm, metav1.UpdateOptions{})
	return err
}

// applySecret creates or updates the Secret
func (k *Kuber) applySecret(ctx context.Context) error {
	sec := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: k.cfg.NameSpace + "-secret",
		},
		StringData: k.cfg.Secret,
	}

	_, err := k.Clientset.CoreV1().Secrets(k.cfg.NameSpace).Get(ctx, sec.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		fmt.Println("creating secret ...")
		_, err = k.Clientset.CoreV1().Secrets(k.cfg.NameSpace).Create(ctx, sec, metav1.CreateOptions{})
		return err
	}
	if err != nil {
		return err
	}

	fmt.Println("updating secret ...")
	_, err = k.Clientset.CoreV1().Secrets(k.cfg.NameSpace).Update(ctx, sec, metav1.UpdateOptions{})
	return err
}

// applyDeployment creates a deployment if it doesn't exist, otherwise updates it
func (k *Kuber) ApplyDeployment(ctx context.Context, spec entity.DeploymentSpec) error {
	existing, err := k.Clientset.AppsV1().Deployments(k.cfg.NameSpace).Get(ctx, spec.Name, metav1.GetOptions{})

	if apierrors.IsNotFound(err) {
		fmt.Printf("creating deployment %s...\n", spec.Name)
		_, err = k.Clientset.AppsV1().Deployments(k.cfg.NameSpace).Create(ctx, spec.Deployment, metav1.CreateOptions{})
		return err
	}

	if err != nil {
		return err
	}

	fmt.Printf("updating deployment %s...\n", spec.Name)
	spec.Deployment.ResourceVersion = existing.ResourceVersion
	_, err = k.Clientset.AppsV1().Deployments(k.cfg.NameSpace).Update(ctx, spec.Deployment, metav1.UpdateOptions{})
	return err
}

// applyService creates a service if it doesn't exist, otherwise updates it
func (k *Kuber) ApplyService(ctx context.Context, spec entity.ServiceSpec) error {
	existing, err := k.Clientset.CoreV1().Services(k.cfg.NameSpace).Get(ctx, spec.Name, metav1.GetOptions{})

	if apierrors.IsNotFound(err) {
		fmt.Printf("creating service %s...\n", spec.Name)
		_, err = k.Clientset.CoreV1().Services(k.cfg.NameSpace).Create(ctx, spec.Service, metav1.CreateOptions{})
		return err
	}

	if err != nil {
		return err
	}

	fmt.Printf("updating service %s...\n", spec.Name)
	spec.Service.ResourceVersion = existing.ResourceVersion
	spec.Service.Spec.ClusterIP = existing.Spec.ClusterIP
	_, err = k.Clientset.CoreV1().Services(k.cfg.NameSpace).Update(ctx, spec.Service, metav1.UpdateOptions{})
	return err
}

// waitForDeploymentReady waits for a deployment to become ready within the timeout
func (k *Kuber) WaitForDeployment(ctx context.Context, name string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		dep, err := k.Clientset.AppsV1().Deployments(k.cfg.NameSpace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return err
		}

		fmt.Printf("%s dep status - ready replica: %v, spec replica: %v\n",
			name, dep.Status.ReadyReplicas, *dep.Spec.Replicas)

		if dep.Status.ReadyReplicas > 0 && dep.Status.ReadyReplicas == *dep.Spec.Replicas {
			fmt.Printf("%s is ready!\n", name)
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}

	return fmt.Errorf("%s did not become ready within %v", name, timeout)
}
