#! /usr/bin/env bash

NAMESPACE=${1:-default}

kubectl delete pong example-pong --ignore-not-found=true -n $NAMESPACE
kubectl delete -f deploy/service_account.yaml -n $NAMESPACE
kubectl delete -f deploy/role.yaml -n $NAMESPACE
kubectl delete -f deploy/role_binding.yaml -n $NAMESPACE
kubectl delete deployment pong-operator -n $NAMESPACE
kubectl delete -f deploy/crds/pong.example.com_pongs_crd.yaml --ignore-not-found=true
