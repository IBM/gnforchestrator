# VNFs on Kubernetes

1.  Set your namespace  
    ```
    $ export NAMESPACE=<YOUR NAMESPACE>
    ```

1.  Deploy the [sleep sample](https://github.com/istio/istio/tree/master/samples/sleep) from the [istio.io](istio.io).
    You will use this sample to send `curl` commands to the PingPong application.
    (No features of [istio.io](istio.io) are used in this demo, you will only use the sample from the project's
    [repository](https://github.com/istio/istio)).

    ```
    $ kubectl apply -f https://raw.githubusercontent.com/istio/istio/master/samples/sleep/sleep.yaml -n $NAMESPACE
    ```

1.  Export the sleep pod as an environment variable:

    ```
    $ export SLEEP_POD=$(kubectl get pod -l app=sleep -n $NAMESPACE -o jsonpath={.items..metadata.name})
    ```
1.  Deploy the network service:

    ```
    $ kubectl apply -f deploy/crds/gnforchestrator.ibm.com_v1alpha1_networkservice_k8s_cr.yaml -n $NAMESPACE
    ```

1.  Watch the resources being created:

    ```
    $ watch kubectl get nsvc,ping,pong,pod -n $NAMESPACE
    ```

1.  Create aliases to get the IP address and port of the network service:

    ```
    alias get-example-pingpong-ip='kubectl get nsvc example-pingpong -n $NAMESPACE -o jsonpath={.status.ip}'
    alias get-example-pingpong-port='kubectl get nsvc example-pingpong -n $NAMESPACE -o jsonpath={.status.port}'
    ```

1.  Run `curl` to access the network service:

    ```
    $ kubectl exec -it $SLEEP_POD -n $NAMESPACE -- curl $(get-example-pingpong-ip):$(get-example-pingpong-port)/hello
    Hello from Ping VNF
    ```

1.  Run `curl` to perform ping pong 3 times:

    ```
    $ kubectl exec -it $SLEEP_POD -n $NAMESPACE -- curl $(get-example-pingpong-ip):$(get-example-pingpong-port)/ping/3
    ping version v1: pong version v2 message Hello.pong version v2 message Hello.pong version v2 message Hello.
    ```

    Notice the versions of ping and pong, and the message.

1.  Check the status of network service. Notice the ip and the port:
    ```
    $ kubectl get nsvc example-pingpong -n $NAMESPACE -o yaml
    ...
    status:
      ip: 172.30.184.192
      port: 6003
    ```

1.  Edit the network service resource and change the message to something different, e.g. to `Hi!!!`:

    ```
    $ kubectl edit nsvc example-pingpong -n $NAMESPACE
    ```

1.  Rerun `curl` to perform ping pong:

    ```
    $ while :; do  kubectl exec -it $SLEEP_POD -n $NAMESPACE -- curl $(get-example-pingpong-ip):$(get-example-pingpong-port)/ping/3; sleep 1; echo; done
    ping version v1: pong version v2 message Hi!!!.pong version v2 message Hi!!!.pong version v2 message Hi!!!
    ```

1.  Edit the network service resource and upgrade the version of pong, e.g. to `v3`:

    ```
    $ kubectl edit nsvc example-pingpong -n $NAMESPACE
    ```

1.  Watch the VNF resource of pong being updated:

    ```
    $ watch kubectl get pong example-pingpong-pong -n $NAMESPACE -o jsonpath={.spec.pongVersion}
    v3
    ```

1.  Watch the pods of pong being restarted:

    ```
    $ watch kubectl get pods -n $NAMESPACE -l vnf=pong
    ```

1.  Rerun `curl` to perform ping pong:

    ```
    $ while :; do kubectl exec -it $SLEEP_POD -n $NAMESPACE -- curl $(get-example-pingpong-ip):$(get-example-pingpong-port)/ping/3; sleep 1; echo; done
    ping version v1: pong version v3 message Hi!!!.pong version v3 message Hi!!!.pong version v3 message Hi!!!
    ```

1.  Perform some fault management, delete the service of pong:

    ```
    $ kubectl delete svc -n $NAMESPACE -l vnf=pong
    ```

1.  Watch the ip of pong and the output ip of ping being updated:

    ```
    $ watch kubectl get nsvc,ping,pong,pod -n $NAMESPACE
    ```

1.  Rerun `curl` to perform ping pong:

    ```
    $ while :; do kubectl exec -it $SLEEP_POD -n $NAMESPACE -- curl $(get-example-pingpong-ip):$(get-example-pingpong-port)/ping/3; sleep 1; echo; done
    ping version v1: pong version v3 message Hi!!!.pong version v3 message Hi!!!.pong version v3 message Hi!!!
    ```

1.  Test the healthcheck of ping, make ping unhealthy:

    ```
    $ kubectl exec -it -n $NAMESPACE $(kubectl get pod -l vnf=ping -n $NAMESPACE -o jsonpath={.items[0].metadata.name}) -c ping-vnf-manager -- curl localhost:6001/debug/unhealthy
    ok
    ```

1.  Watch the pod of ping and see it restarted (the number of `RESTARTS` increased):

    ```
    $ watch kubectl get pods -n $NAMESPACE -l vnf=ping
    NAME                                         READY   STATUS    RESTARTS   AGE
    example-pingpong-ping-ping-797c94b8f-v8rwv   2/2     Running   1          176m
    ```

1.  Check the number of the pong pods:

    ```
    $ kubectl get pod -l vnf=pong -n $NAMESPACE
    NAME                                          READY   STATUS    RESTARTS   AGE
    example-pingpong-pong-pong-666d8d6c79-f96w4   2/2     Running   0          8m40s
    example-pingpong-pong-pong-666d8d6c79-j6xk5   2/2     Running   0          8m18s
    ```

1.  Scale the number of replicas of the network service to 2:

    ```
    $ kubectl scale nsvc example-pingpong --replicas 2 -n $NAMESPACE
    ```

1.  Watch the number of pong pods:

    ```
    $ kubectl get pod -l vnf=pong -n $NAMESPACE
    NAME                                                              READY   STATUS    RESTARTS   AGE
    example-pingpong-pong-pong-666d8d6c79-9lt8j   2/2     Running   0          22s
    example-pingpong-pong-pong-666d8d6c79-f96w4   2/2     Running   0          12m
    example-pingpong-pong-pong-666d8d6c79-j6xk5   2/2     Running   0          12m
    example-pingpong-pong-pong-666d8d6c79-smq4s   2/2     Running   0          22s
    ```

1.  Cause the horizontal pod autoscaling of pong to increase the number of replicas, by increasing the value
    of a custom metric of pong pods, each pod will report metric value 10:

    ```
    $ for PONG_POD in $(kubectl get pod -l vnf=pong -n $NAMESPACE -o jsonpath='{.items[*].metadata.name}'); do kubectl exec -it -n $NAMESPACE $PONG_POD -c pong-vnf-manager -- curl localhost:6004/debug/metric/10; done
    ok
    ok
    ok
    ok
    ```

1.  Watch the horizontal pod autoscaling to take action - the number of pong pods should increase
    (it may take several minutes).

    ```
    $ watch kubectl get pod -l vnf=pong -n $NAMESPACE
    example-pingpong-pong-pong-666d8d6c79-9lt8j   2/2     Running   0          4m1s
    example-pingpong-pong-pong-666d8d6c79-f96w4   2/2     Running   0          16m
    example-pingpong-pong-pong-666d8d6c79-j6xk5   2/2     Running   0          15m
    example-pingpong-pong-pong-666d8d6c79-smq4s   2/2     Running   0          4m1s
    example-pingpong-pong-pong-666d8d6c79-vqjtv   2/2     Running   0          45s
    example-pingpong-pong-pong-666d8d6c79-w67cr   2/2     Running   0          45s
    ```

1.  Check the number of replicas of pong pods in the status of the pong CR:

    ```
    $ kubectl get pong example-pingpong-pong -n $NAMESPACE -o yaml
    ...
    status:
      ip: 172.30.144.182
      port: 6006
      replicas: 6
    ```

1.  Reduce pong custom metric value back to 0:

    ```
    $ for PONG_POD in $(kubectl get pod -l vnf=pong -n $NAMESPACE -o jsonpath='{.items[*].metadata.name}'); do kubectl exec -it -n $NAMESPACE $PONG_POD -c pong-vnf-manager -- curl localhost:6004/debug/metric/0; done
    ok
    ok
    ok
    ok
    ok
    ok
    ```

1.  Deploy another pingpong network service:

    ```
    $ kubectl apply -f deploy/crds/gnforchestrator.ibm.com_v1alpha1_networkservice_k8s_another_cr.yaml -n $NAMESPACE
    ```

1.  Watch new resources related to the second network service are being created:

    ```
    $ watch kubectl get nsvc,ping,pong,pod -n $NAMESPACE
    ```

1.  Create aliases to get the IP address and port of the network service:

    ```
    alias get-another-example-pingpong-ip='kubectl get nsvc another-example-pingpong -n $NAMESPACE -o jsonpath={.status.ip}'
    alias get-another-example-pingpong-port='kubectl get nsvc another-example-pingpong -n $NAMESPACE -o jsonpath={.status.port}'
    ```

1.  Run `curl` to perform ping pong with the second network service:

    ```
    $ kubectl exec -it $SLEEP_POD -n $NAMESPACE -- curl $(kubectl get nsvc another-example-pingpong -n $NAMESPACE -o jsonpath={.status.ip}):$(kubectl get nsvc another-example-pingpong -n $NAMESPACE -o jsonpath={.status.port})/ping/3
    ping version v1: pong version v2 message AnotherHello.pong version v2 message AnotherHello.pong version v2 message AnotherHello.
    ```

1.  Delete the network service:

    ```
    $ kubectl delete -f deploy/crds/gnforchestrator.ibm.com_v1alpha1_networkservice_k8s_cr.yaml -n $NAMESPACE
    ```

1.  Watch the resources being deleted:

    ```
    $ watch kubectl get k8svnf,pong,ingress -n $NAMESPACE
    ```

1.  Delete the second network service:

    ```
    $ kubectl delete -f deploy/crds/gnforchestrator.ibm.com_v1alpha1_networkservice_k8s_another_cr.yaml -n $NAMESPACE
    ```

1.  Delete the aliases:

    ```
    $ unalias get-example-pingpong-ip get-example-pingpong-port
    ```

1.  Delete the sleep pod:

    ```
    $ kubectl delete -f https://raw.githubusercontent.com/istio/istio/master/samples/sleep/sleep.yaml -n $NAMESPACE
    ```

1. Unset the `SLEEP_POD` environment variable:

   ```
   $ unset SLEEP_POD
   ```
