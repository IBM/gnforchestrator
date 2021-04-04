#! /usr/bin/env bash

# (C) Copyright IBM Corporation 2021.
# LICENSE: GPL-3.0 https://opensource.org/licenses/GPL-3.0

NAMESPACE=${1:-default}

export NUMBER_OF_OPERATORS=$(kubectl get deployment -l name=gnforchestrator --all-namespaces -o jsonpath='{.items}' | wc -l)

if [ "$NUMBER_OF_OPERATORS" -le 1 ] ; then
envsubst < deploy/cluster_role_binding.yaml.template | kubectl delete -f -
kubectl delete -f deploy/cluster_role.yaml
kubectl delete -f deploy/crds/gnforchestrator.ibm.com_networkservices_crd.yaml --ignore-not-found=true
else
cat <<EOF
***********************************************************************************************************
* Not deleting cluster role, cluster role binding and the CRD since more instances of the operator exist  *
* check the existing operators by running kubectl get deployment -l name=gnforchestrator --all-namespaces *                                                        *
***********************************************************************************************************
EOF
fi

kubectl delete -f deploy/role_binding.yaml -n $NAMESPACE
kubectl delete -f deploy/role.yaml -n $NAMESPACE
kubectl delete -f deploy/service_account.yaml -n $NAMESPACE
kubectl delete deployment gnforchestrator -n $NAMESPACE
