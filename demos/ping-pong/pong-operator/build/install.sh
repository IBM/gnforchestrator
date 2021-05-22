#! /usr/bin/env bash

NAMESPACE=${1:-default}

kubectl apply -f deploy/crds/pong.example.com_pongs_crd.yaml

kubectl apply -f deploy/service_account.yaml -n $NAMESPACE
kubectl apply -f deploy/role.yaml -n $NAMESPACE
kubectl apply -f deploy/role_binding.yaml -n $NAMESPACE
envsubst < deploy/operator.yaml.template | kubectl apply -n ${NAMESPACE} -f -
