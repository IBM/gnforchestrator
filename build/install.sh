#! /usr/bin/env bash

# (C) Copyright IBM Corporation 2021.
# LICENSE: GPL-3.0 https://opensource.org/licenses/GPL-3.0

NAMESPACE=${1:-default}

kubectl apply -f deploy/service_account.yaml -n $NAMESPACE
kubectl apply -f deploy/role.yaml -n $NAMESPACE
kubectl apply -f deploy/role_binding.yaml -n $NAMESPACE

kubectl apply -f deploy/crds/gnforchestrator.ibm.com_networkservices_crd.yaml
kubectl apply -f deploy/cluster_role.yaml -n $NAMESPACE


export NUMBER_OF_OPERATORS=$(kubectl get deployment -l name=gnforchestrator --all-namespaces -o jsonpath='{.items}' | wc -l)

if [ "$NUMBER_OF_OPERATORS" -lt 1 ] ; then
envsubst < deploy/cluster_role_binding.yaml.template | kubectl apply -f -
else
cat <<EOF
**************************************************************************************************
* Update cluster role binding manually, perform kubectl edit clusterrolebinding gnforchestrator  *
* add your service account as a subject                                                          *
**************************************************************************************************
EOF
fi

envsubst < deploy/operator.yaml.template | kubectl apply -n ${NAMESPACE} -f -
