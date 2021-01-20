# OsushiOperator

OshushiOperator is a Kubernetes Operator managing Osushi 🍣 as a CRDs(Custom Resource Definition).

# Quick Start

Create a cluster and launch the OsushiOperator.

```
$ kind create cluster

$ docker build . -t osushi-operator
...
Successfully tagged osushi-operator:latest

$ kind load docker-image  osushi-operator

$ k apply -f config/rbac/

$ make instal

$ make run ENABLE_WEBHOOKS=false

```

Now, you can see osushi on your terminal.

```
$ k apply -f config/samples/cache_v1alpha1_osushi.yaml

$ k get osushi
NAME    SIZE   EMOJI
syake   1                                                                                               🍣
```
