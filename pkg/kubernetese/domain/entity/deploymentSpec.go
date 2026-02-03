package entity

import appsv1 "k8s.io/api/apps/v1"

// DeploymentSpec holds a deployment definition
type DeploymentSpec struct {
	Name       string
	Deployment *appsv1.Deployment
}
