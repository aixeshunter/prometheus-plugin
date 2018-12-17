## Prometheus Plugin

Deletint the pods of statefulset in k8s prometheus cluster when node power down.

The plugin is to resolve the issue [Resource Prometheus can be created by Deployment? #2214](https://github.com/coreos/prometheus-operator/issues/2214).

### Prerequisites

- Kubernetes version = 1.11.0

### Make
```
build:  Go build
docker: build and run in docker container
gotest: run go tests and reformats
format: formatting code
vet:    vetting code
```

**build**: runs go build for nfs_exporter

**docker**: runs docker build and copy new built nfs_exporter


### Parameters

| Option                    | Default             | Description
| ------------------------- | ------------------- | -----------------
| -h, --help                | -                   | Displays usage.
| --namespace     | "kube-system"             | The namespace that pods in
| --prometheus-name     | -             | prometheus pod name
| --alertmanager-name    | -             | alertmanager pod name
| --period   | "1m"             | The loop time to detect


### deploy

[Yaml file](manifests)


### Points

The Statefulset pods need to be deleted forcely.[Reference](https://kubernetes.io/docs/tasks/run-application/force-delete-stateful-set-pod/)

Deleting the prometheus pod needs `force delete` and `patch finalizers`:

```bash
kubectl delete po -n monitoring  prometheus-k8s-0  --grace-period=0 --force (the command will hang)

kubectl patch pod prometheus-k8s-0 -p '{"metadata":{"finalizers":null}}' -n monitoring (delete command will return)
```

So, I need to delete it by channel:

```goalng
done := make(chan bool)
var gracePeriodSeconds int64 = 0
if isPodUnknown(pod) {
    go DeletePod(m.client, &metav1.DeleteOptions{GracePeriodSeconds: &gracePeriodSeconds}, p, m.Namespace, done)
}

select {
case <-done:
    glog.Infof("Delete the pod %s will return, no timeout.", p)
case <-time.After(constants.DeletePodTimeout):                          // wait the delete command timeout
    glog.Infof("Delete the pod %s timeout, need to patch finalizers.", p)
    PatchPod(m.client, p, m.Namespace)
}
```

### Godep

[godep](https://github.com/tools/godep) is an older dependency management tool, which is
used by the main Kubernetes repo and `client-go` to manage dependencies.

Before proceeding with the below instructions, you should ensure that your
$GOPATH is empty except for containing your own package and its dependencies,
and you have a copy of godep somewhere in your $PATH.

To install `client-go` and place its dependencies in your `$GOPATH`:

```sh
go get k8s.io/client-go/...
cd $GOPATH/src/k8s.io/client-go
git checkout v9.0.0 # replace v9.0.0 with the required version
# cd 1.5 # only necessary with 1.5 and 1.4 clients.
godep restore ./...
```

At this point, `client-go`'s dependencies have been placed in your $GOPATH, but
if you were to build, `client-go` would still see its own copy of its
dependencies in its `vendor` directory. You have two options at this point.

If you would like to keep dependencies in your own project's vendor directory,
then you can continue like this:

```sh
cd $GOPATH/src/<my-pkg>
godep save ./...
```

Alternatively, if you want to build using the dependencies in your `$GOPATH`,
then `rm -rf vendor/` to remove `client-go`'s copy of its dependencies.