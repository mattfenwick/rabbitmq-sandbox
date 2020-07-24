# Setup

```
git clone git@github.com:mattfenwick/rabbitmq-sandbox.git

cd rabbitmq-sandbox
```

[Download golang tools](https://golang.org/dl/)!

Running the code:

```
go run cmd/read/main.go

go run cmd/send/main.go
```

## Kind/rabbitmq setup

```
git clone git@github.com:rabbitmq/rabbitmq-peer-discovery-k8s.git

cd rabbitmq-peer-discovery-k8s/examples/kind

kind create cluster --config kind-cluster/kind-cluster.yaml

kubectl apply -k overlays/dev-persistence/
```

### Inspecting the kube bits of rabbitmq

```
kubectl get pvc -n rabbitmq-dev-persistence
kubectl get pv -n rabbitmq-dev-persistence
kubectl get pods -n rabbitmq-dev-persistence
```

### Shell into a rabbitmq pod

```
kubectl exec -ti rabbitmq-0 -n rabbitmq-dev-persistence bash
```

Inside of a rabbitmq pod:
```
cd /var/lib/rabbitmq/mnesia/rabbit\@rabbitmq-0.rabbitmq.rabbitmq-dev-persistence.svc.cluster.local
cat quorum/rabbit\@rabbitmq-0.rabbitmq.rabbitmq-dev-persistence.svc.cluster.local/00000001.wal
```
