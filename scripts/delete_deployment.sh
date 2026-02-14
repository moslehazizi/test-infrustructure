#!/usr/bin/bash

KUBECONFIG=/home/appuser/.kube/config

deployment=${1:-"challenge-mother-service-10"}
endpoint=apis/apps/v1/namespaces/$KUBERNETES_NAMESPACE/deployments/$deployment
ca_cert=/var/run/secrets/kubernetes.io/serviceaccount/ca.crt
token=$(grep token "$KUBECONFIG" | cut -d ":" -f 2 | tr -d "[:space:]")
api_server="https://kubernetes.default.svc"

curl -X DELETE $api_server/$endpoint -H "Authorization: Bearer $token" --cacert $ca_cert