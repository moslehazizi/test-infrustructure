package entity

import networkingv1 "k8s.io/api/networking/v1"

type IngressSpec struct {
	Name    string
	Ingress *networkingv1.Ingress
}
