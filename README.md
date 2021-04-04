# GNF Orchestrator
Generic Network Functions Orchestrator, A prototypical implementation of a Kubernetes-native NFV MANO.

The goal of this prototype is to prove that Kubernetes with Operator SDK can be used AS IS as NFV MANO.

| NFV MANO Function | Implementaion in Kubernetes/Operators |
|-------------------|:-------------------------------------:|
| Onboarding        | Operator Lifecycle Manager            |
| Instantiate       | kubectl create                        |
| Scale             | kubectl scale + horizontal pod autoscaling|
| Update configuration| kubectl apply/edit                  |
| Upgrade           | kubectl apply/edit|
| Fault tolerance   | Kubernetes Deployments + reconciliation by the operators + Kubernetes health checks|
| Terminate         | kubectl delete + owner references + finalizers|

## CRDs

```
apiVersion: gnforchestrator.ibm.com/v2alpha1
kind: NetworkService
metadata:
  name: example-pingpong
  labels:
    service: pingpong
spec:
  properties:
    message: Hello
  components:
    ping:
      template:
        apiVersion: ping.example.com/v1alpha1
        kind: Ping
        metadata:
          name: "[% meta.name %]-ping"
          namespace: "[% meta.namespace %]"
        spec:
          pingVersion: v1.0
          pongAddress: "[% pong.status.ip %]"
          pongPort: 6006
    pong:
      template:
        apiVersion: pong.example.com/v1alpha1
        kind: Pong
        metadata:
          name: "[% meta.name %]-pong"
          namespace: "[% meta.namespace %]"
        spec:
          pongVersion: v1.4
          message: "[% message %]"
          minReplicas: "[% replicas*2 %]"
          maxReplicas: "[% replicas*3 %]"
  replicas: 2
  statusTemplate:
    ip: "[% ping.status.ip %]"
status:
  ip: 1.1.1.1 # the status of ping's IP
```

## Prerequisites

1. OperatoSDK, starting from 0.16.0.


## Install of the operator

1.  Create a namespace in Kubernetes

1.  Set the `NAMESPACE` environment variable to hold the name of your namespace. The following command sets it to be
    equal to your current user id.

    ```
    $ export NAMESPACE=$USER
    ```

1.  Set the `REGISTRY` environment variable to hold the name of your docker registry. Omit for docker.com.

    ```
    $ export REGISTRY=
    ```

1.  Set the `IMAGE` environment variable to hold the image of the operator.

    ```
    $ export IMAGE=$REGISTRY/$(basename $(pwd)):v2.1
    ```

1.  Run make to install the operator:

    ```
    $ make install "IMAGE=$IMAGE" "NAMESPACE=$NAMESPACE"
    ```

    Add `docker-push` to the list of targets in the command above if you want to push the operator to a repository.

### Cleanup

```
$ make clean "NAMESPACE=$NAMESPACE"
```
