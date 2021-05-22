# pong-operator

## Install of the operator

1.  Create a namespace in Kubernetes

1.  Set the `NAMESPACE` environment variable to hold the name of your namespace.  

    ```
    $ export NAMESPACE=<YOUR NAMESPACE>
    ```

1.  Set the `REGISTRY` environment variable to hold the name of your docker registry.  

    ```
    $ export REGISTRY=<YOUR_REGISTRY>
    ```

1.  Set the `IMAGE` environment variable to hold the image of the operator.

    ```
    $ export IMAGE=$REGISTRY/$(basename $(pwd)):v0.0.1
    ```

1.  Run make to install the operator:

    ```
    $ make install "IMAGE=$IMAGE" "NAMESPACE=$NAMESPACE"
    ```

    Add `docker-push` to the list of targets in the command above if you want to push the operator to a repository.

## Troubleshooting

Check the logs of the operator:

```
$ kubectl logs -l name=pong-operator -c operator --tail=1000 -n $NAMESPACE
```

More information is in the logs of Ansible:

```
$ kubectl logs -l name=pong-operator -c ansible --tail=1000 -n $NAMESPACE
```

### Push the docker image

```
$ make docker-push "IMAGE=$IMAGE"
```

### Cleanup

```
$ make clean "NAMESPACE=$NAMESPACE"
```
