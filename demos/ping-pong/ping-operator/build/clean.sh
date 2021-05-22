#! /usr/bin/env bash

NAMESPACE=${1:-default}

kubectl delete ping example-ping --ignore-not-found=true -n $NAMESPACE
kubectl delete -f deploy/service_account.yaml -n $NAMESPACE
kubectl delete -f deploy/role.yaml -n $NAMESPACE
kubectl delete -f deploy/role_binding.yaml -n $NAMESPACE
kubectl delete deployment ping-operator -n $NAMESPACE
kubectl delete -f deploy/crds/ping.example.com_pings_crd.yaml --ignore-not-found=true
