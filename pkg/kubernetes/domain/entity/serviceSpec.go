package entity

import corev1 "k8s.io/api/core/v1"

// ServiceSpec holds a service definition
type ServiceSpec struct {
	Name    string
	Service *corev1.Service
}
