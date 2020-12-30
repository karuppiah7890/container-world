# Part 3 - Worker Node

Now, for the worker node, we need quite a few things.

Let's start with the most basic thing maybe? The kubelet. We will then move on
to the container runtime, container networking interface, kube-proxy

We won't have to work on networking for pod across multiple nodes as we will
only have a single node cluster for now which has all components - control
plan and worker node components

Now, let's get started with the kubelet

```bash
$ kubelet
F1230 10:12:00.023238    2162 server.go:257] mkdir /var/lib/kubelet: permission denied
goroutine 1 [running]:
k8s.io/kubernetes/vendor/k8s.io/klog/v2.stacks(0xc000010001, 0xc00089f600, 0x57, 0xa9)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:1026 +0xb9
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).output(0x70c9460, 0xc000000003, 0x0, 0x0, 0xc0001b8000, 0x6f34162, 0x9, 0x101, 0x411b00)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:975 +0x19b
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).printDepth(0x70c9460, 0x
```

Let's try to use a different directory for kubelet maybe? Something in `/opt`
maybe :) As I'm not sure if `/var` directories will survive restarts

```bash
$ mkdir /opt/kubelet
mkdir: cannot create directory ‘/opt/kubelet’: Permission denied
$ sudo mkdir -p /opt/kubelet
$ chown -R ubuntu:ubuntu /opt/kubelet/
chown: changing ownership of '/opt/kubelet/': Operation not permitted
$ sudo chown -R ubuntu:ubuntu /opt/kubelet/
```

Now, we need to check how to set the directory and also have some sort of
config to connect to API server

https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/

Looks like the config is this

```bash
--root-dir string                                          Directory path for managing kubelet files (volume mounts,etc). (default "/var/lib/kubelet")
```

A related thing is

```bash
--seccomp-profile-root string                              <Warning: Alpha feature> Directory path for seccomp profiles. (default "/var/lib/kubelet/seccomp") (DEPRECATED: will be removed in 1.23, in favor of using the `<root-dir>/seccomp` directory)
```

But let's just use `--root-dir` and try to see how it all works :)

```bash
$ kubelet --root-dir /opt/kubelet/
F1230 16:00:36.734501    2365 server.go:257] mkdir /var/lib/kubelet: permission denied
goroutine 1 [running]:
k8s.io/kubernetes/vendor/k8s.io/klog/v2.stacks(0xc000010001, 0xc00079aa50, 0x57, 0xa9)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:1026 +0xb9
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).output(0x70c9460, 0xc000000003, 0x0, 0x0, 0xc0007bc000, 0x6f34162, 0x9, 0x101, 0x411b00)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:975 +0x19b
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).printDepth(0x70c9460, 0xc000000003, 0x0, 0x0, 0x0, 0x0, 0x1, 0xc0007f06f0, 0x1, 0x1)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:732 +0x16f
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).print(...)
```
