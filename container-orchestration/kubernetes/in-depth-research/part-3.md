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

It's still trying to create a directory. Hmm. Is it possible that it maybe for
the seccomp profile?

```bash
$ kubelet --root-dir /opt/kubelet/ --seccomp-profile-root /opt/kubelet/seccomp
Flag --seccomp-profile-root has been deprecated, will be removed in 1.23, in favor of using the `<root-dir>/seccomp` directory
F1230 18:43:29.757550    1292 server.go:257] mkdir /var/lib/kubelet: permission denied
goroutine 1 [running]:
k8s.io/kubernetes/vendor/k8s.io/klog/v2.stacks(0xc000010001, 0xc000a2f340, 0x57, 0xa9)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:1026 +0xb9
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).output(0x70c9460, 0xc000000003, 0x0, 0x0, 0xc0007be000, 0x6f34162, 0x9, 0x101, 0x411b00)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:975 +0x19b
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).printDepth(0x70c9460, 0xc000000003, 0x0, 0x0, 0x0, 0x0, 0x1, 0xc0007f26f0, 0x1, 0x1)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:732 +0x16f
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).print(...)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:714
k8s.io/kubernetes/vendor/k8s.io/klog/v2.Fatal(...)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:1482
k8s.io/kubernetes/cmd/kubelet/app.NewKubeletCommand.func1(0xc00048b8c0, 0xc00004e0b0, 0x4, 0x4)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/cmd/kubelet/app/server.go:257 +0x62b
k8s.io/kubernetes/vendor/github.com/spf13/cobra.(*Command).execute(0xc00048b8c0, 0xc00004e0b0, 0x4, 0x4, 0xc00048b8c0, 0xc00004e0b0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/github.com/spf13/cobra/command.go:854 +0x2c2
k8s.io/kubernetes/vendor/github.com/spf13/cobra.(*Command).ExecuteC(0xc00048b8c0, 0x165580c13ee99bc9, 0x70c9020, 0x409b25)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/github.com/spf13/cobra/command.go:958 +0x375
k8s.io/kubernetes/vendor/github.com/spf13/cobra.(*Command).Execute(...)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/github.com/spf13/cobra/command.go:895
main.main()
        _output/dockerized/go/src/k8s.io/kubernetes/cmd/kubelet/kubelet.go:41 +0xe5

goroutine 6 [chan receive]:
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).flushDaemon(0x70c9460)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:1169 +0x8b
created by k8s.io/kubernetes/vendor/k8s.io/klog/v2.init.0
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:417 +0xdf

goroutine 79 [select]:
k8s.io/kubernetes/vendor/go.opencensus.io/stats/view.(*worker).start(0xc00004f360)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/go.opencensus.io/stats/view/worker.go:154 +0x105
created by k8s.io/kubernetes/vendor/go.opencensus.io/stats/view.init.0
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/go.opencensus.io/stats/view/worker.go:32 +0x57

goroutine 94 [select]:
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.BackoffUntil(0x4a721f0, 0x4f0bde0, 0xc00038e150, 0x1, 0xc00009a0c0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:167 +0x149
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.JitterUntil(0x4a721f0, 0x12a05f200, 0x0, 0xc000b44d01, 0xc00009a0c0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:133 +0x98
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.Until(...)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:90
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.Forever(0x4a721f0, 0x12a05f200)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:81 +0x4f
created by k8s.io/kubernetes/vendor/k8s.io/component-base/logs.InitLogs
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/component-base/logs/logs.go:58 +0x8a
```

Still doesn't work. Hmm. Right. It's the TLS related directory to create certs

```bash
$ kubelet --root-dir /opt/kubelet/ --help | grep '\/var\/lib\/kubelet'
      --cert-dir string                                          The directory where the TLS certs are located. If --tls-cert-file and --tls-private-key-file are provided, this flag will be ignored. (default "/var/lib/kubelet/pki")
      --root-dir string                                          Directory path for managing kubelet files (volume mounts,etc). (default "/var/lib/kubelet")
      --seccomp-profile-root string                              <Warning: Alpha feature> Directory path for seccomp profiles. (default "/var/lib/kubelet/seccomp") (DEPRECATED: will be removed in 1.23, in favor of using the `<root-dir>/seccomp` directory)
```

```bash
$ kubelet --root-dir /opt/kubelet/ --help | grep tls
      --allowed-unsafe-sysctls strings                           Comma-separated whitelist of unsafe sysctls or unsafe sysctl patterns (ending in *). Use these at your own risk. (DEPRECATED: This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.)
      --cert-dir string                                          The directory where the TLS certs are located. If --tls-cert-file and --tls-private-key-file are provided, this flag will be ignored. (default "/var/lib/kubelet/pki")
      --docker-tls                                               use TLS to connect to docker (DEPRECATED: This is a cadvisor flag that was mistakenly registered with the Kubelet. Due to legacy concerns, it will follow the standard CLI deprecation timeline before being removed.)
      --docker-tls-ca string                                     path to trusted CA (default "ca.pem") (DEPRECATED: This is a cadvisor flag that was mistakenly registered with the Kubelet. Due to legacy concerns, it will follow the standard CLI deprecation timeline before being removed.)
      --docker-tls-cert string                                   path to client certificate (default "cert.pem") (DEPRECATED: This is a cadvisor flag that was mistakenly registered with the Kubelet. Due to legacy concerns, it will follow the standard CLI deprecation timeline before being removed.)
      --docker-tls-key string                                    path to private key (default "key.pem") (DEPRECATED: This is a cadvisor flag that was mistakenly registered with the Kubelet. Due to legacy concerns, it will follow the standard CLI deprecation timeline before being removed.)
                Sysctls=true|false (BETA - default=true)
      --tls-cert-file string                                     File containing x509 Certificate used for serving HTTPS (with intermediate certs, if any, concatenated after server cert). If --tls-cert-file and --tls-private-key-file are not provided, a self-signed certificate and key are generated for the public address and saved to the directory passed to --cert-dir. (DEPRECATED: This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.)
      --tls-cipher-suites strings                                Comma-separated list of cipher suites for the server. If omitted, the default Go cipher suites will be used.
      --tls-min-version string                                   Minimum TLS version supported. Possible values: VersionTLS10, VersionTLS11, VersionTLS12, VersionTLS13 (DEPRECATED: This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.)
      --tls-private-key-file string                              File containing x509 private key matching --tls-cert-file. (DEPRECATED: This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.)
```

Let's get to it then! :)

```bash
$ cat > kubelet-server-csr.json <<EOF
{
    "CN": "kubelet-server",
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "names": [
        {
            "C": "US",
            "L": "San Francisco",
            "O": "Kubernetes",
            "OU": "Kubelet Server",
            "ST": "California"
        }
    ]
}
EOF

$ cfssl gencert \
  -ca=ca.pem \
  -ca-key=ca-key.pem \
  -config=ca-config.json \
  -profile=kubernetes \
  kubelet-server-csr.json | cfssljson -bare kubelet-server
```

Now it worked for a bit :P

```bash
$ kubelet --root-dir /opt/kubelet/ --tls-cert-file kubelet-server.pem --tls-private-key-file kubelet-server-key.pem
Flag --tls-cert-file has been deprecated, This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.
Flag --tls-private-key-file has been deprecated, This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.
I1230 18:50:30.503906    1395 server.go:416] Version: v1.20.1
W1230 18:50:30.505273    1395 server.go:558] standalone mode, no API client
W1230 18:50:30.505923    1395 server.go:614] failed to get the kubelet's cgroup: cpu and memory cgroup hierarchy not unified.  cpu: /, memory: /user.slice/user-1000.slice/session-5.scope.  Kubelet system container metrics may be missing.
W1230 18:50:36.447707    1395 server.go:473] No api server defined - no events will be sent to API server.
I1230 18:50:36.448411    1395 server.go:645] --cgroups-per-qos enabled, but --cgroup-root was not specified.  defaulting to /
I1230 18:50:36.449642    1395 container_manager_linux.go:274] container manager verified user specified cgroup-root exists: []
I1230 18:50:36.450950    1395 container_manager_linux.go:279] Creating Container Manager object based on Node Config: {RuntimeCgroupsName: SystemCgroupsName: KubeletCgroupsName: ContainerRuntime:docker CgroupsPerQOS:true CgroupRoot:/ CgroupDriver:cgroupfs KubeletRootDir:/opt/kubelet/ ProtectKernelDefaults:false NodeAllocatableConfig:{KubeReservedCgroupName: SystemReservedCgroupName: ReservedSystemCPUs: EnforceNodeAllocatable:map[pods:{}] KubeReserved:map[] SystemReserved:map[] HardEvictionThresholds:[{Signal:memory.available Operator:LessThan Value:{Quantity:100Mi Percentage:0} GracePeriod:0s MinReclaim:<nil>} {Signal:nodefs.available Operator:LessThan Value:{Quantity:<nil> Percentage:0.1} GracePeriod:0s MinReclaim:<nil>} {Signal:nodefs.inodesFree Operator:LessThan Value:{Quantity:<nil> Percentage:0.05} GracePeriod:0s MinReclaim:<nil>} {Signal:imagefs.available Operator:LessThan Value:{Quantity:<nil> Percentage:0.15} GracePeriod:0s MinReclaim:<nil>}]} QOSReserved:map[] ExperimentalCPUManagerPolicy:none ExperimentalTopologyManagerScope:container ExperimentalCPUManagerReconcilePeriod:10s ExperimentalPodPidsLimit:-1 EnforceCPULimits:true CPUCFSQuotaPeriod:100ms ExperimentalTopologyManagerPolicy:none}
I1230 18:50:36.451759    1395 topology_manager.go:120] [topologymanager] Creating topology manager with none policy per container scope
I1230 18:50:36.452423    1395 container_manager_linux.go:310] [topologymanager] Initializing Topology Manager with none policy and container-level scope
I1230 18:50:36.453021    1395 container_manager_linux.go:315] Creating device plugin manager: true
F1230 18:50:36.453729    1395 server.go:269] failed to run Kubelet: failed to initialize checkpoint manager: mkdir /var/lib/kubelet: permission denied
goroutine 1 [running]:
k8s.io/kubernetes/vendor/k8s.io/klog/v2.stacks(0xc000010001, 0xc0001da4e0, 0x97, 0xd0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:1026 +0xb9
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).output(0x70c9460, 0xc000000003, 0x0, 0x0, 0xc0007c4700, 0x6f34162, 0x9, 0x10d, 0x411b00)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:975 +0x19b
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).printDepth(0x70c9460, 0xc000000003, 0x0, 0x0, 0x0, 0x0, 0x1, 0xc0007993a0, 0x1, 0x1)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:732 +0x16f
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).print(...)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:714
k8s.io/kubernetes/vendor/k8s.io/klog/v2.Fatal(...)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:1482
k8s.io/kubernetes/cmd/kubelet/app.NewKubeletCommand.func1(0xc0004858c0, 0xc00004e080, 0x6, 0x6)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/cmd/kubelet/app/server.go:269 +0x845
k8s.io/kubernetes/vendor/github.com/spf13/cobra.(*Command).execute(0xc0004858c0, 0xc00004e080, 0x6, 0x6, 0xc0004858c0, 0xc00004e080)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/github.com/spf13/cobra/command.go:854 +0x2c2
k8s.io/kubernetes/vendor/github.com/spf13/cobra.(*Command).ExecuteC(0xc0004858c0, 0x165581234dde2676, 0x70c9020, 0x409b25)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/github.com/spf13/cobra/command.go:958 +0x375
k8s.io/kubernetes/vendor/github.com/spf13/cobra.(*Command).Execute(...)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/github.com/spf13/cobra/command.go:895
main.main()
        _output/dockerized/go/src/k8s.io/kubernetes/cmd/kubelet/kubelet.go:41 +0xe5

goroutine 6 [chan receive]:
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).flushDaemon(0x70c9460)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:1169 +0x8b
created by k8s.io/kubernetes/vendor/k8s.io/klog/v2.init.0
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:417 +0xdf

goroutine 110 [chan receive]:
k8s.io/kubernetes/vendor/k8s.io/client-go/tools/record.(*eventBroadcasterImpl).StartEventWatcher.func1(0x4f1c9a0, 0xc0007ecc60, 0xc00084ded0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/client-go/tools/record/event.go:301 +0xaa
created by k8s.io/kubernetes/vendor/k8s.io/client-go/tools/record.(*eventBroadcasterImpl).StartEventWatcher
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/client-go/tools/record/event.go:299 +0x6e

goroutine 109 [chan receive]:
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/watch.(*Broadcaster).loop(0xc000b3bac0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/watch/mux.go:219 +0x66
created by k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/watch.NewBroadcaster
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/watch/mux.go:73 +0xf7

goroutine 98 [chan receive]:
k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/server.SetupSignalContext.func1(0xc0007f8bf0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/server/signal.go:48 +0x36
created by k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/server.SetupSignalContext
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/server/signal.go:47 +0xf3

goroutine 79 [select]:
k8s.io/kubernetes/vendor/go.opencensus.io/stats/view.(*worker).start(0xc00009b220)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/go.opencensus.io/stats/view/worker.go:154 +0x105
created by k8s.io/kubernetes/vendor/go.opencensus.io/stats/view.init.0
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/go.opencensus.io/stats/view/worker.go:32 +0x57

goroutine 94 [select]:
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.BackoffUntil(0x4a721f0, 0x4f0bde0, 0xc00038e7b0, 0x1, 0xc00009c0c0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:167 +0x149
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.JitterUntil(0x4a721f0, 0x12a05f200, 0x0, 0xc000b44d01, 0xc00009c0c0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:133 +0x98
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.Until(...)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:90
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.Forever(0x4a721f0, 0x12a05f200)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:81 +0x4f
created by k8s.io/kubernetes/vendor/k8s.io/component-base/logs.InitLogs
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/component-base/logs/logs.go:58 +0x8a

goroutine 97 [syscall]:
os/signal.signal_recv(0xc00056b730)
        /usr/local/go/src/runtime/sigqueue.go:147 +0x9d
os/signal.loop()
        /usr/local/go/src/os/signal/signal_unix.go:23 +0x25
created by os/signal.Notify.func1.1
        /usr/local/go/src/os/signal/signal.go:150 +0x45
```

It's still trying to create a directory under /var/lib/kubelet

I'm trying with seccomp thing

```bash
$ kubelet --root-dir /opt/kubelet/ --tls-cert-file kubelet-server.pem --tls-private-key-file kubelet-server-key.pem --seccomp-profile-root /opt/kubelet/seccomp
Flag --tls-cert-file has been deprecated, This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.
Flag --tls-private-key-file has been deprecated, This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.
Flag --seccomp-profile-root has been deprecated, will be removed in 1.23, in favor of using the `<root-dir>/seccomp` directory
I1230 18:52:39.262576    1410 server.go:416] Version: v1.20.1
W1230 18:52:39.264968    1410 server.go:558] standalone mode, no API client
W1230 18:52:39.265910    1410 server.go:614] failed to get the kubelet's cgroup: cpu and memory cgroup hierarchy not unified.  cpu: /, memory: /user.slice/user-1000.slice/session-5.scope.  Kubelet system container metrics may be missing.
W1230 18:52:44.300543    1410 server.go:473] No api server defined - no events will be sent to API server.
I1230 18:52:44.301558    1410 server.go:645] --cgroups-per-qos enabled, but --cgroup-root was not specified.  defaulting to /
I1230 18:52:44.304037    1410 container_manager_linux.go:274] container manager verified user specified cgroup-root exists: []
I1230 18:52:44.305096    1410 container_manager_linux.go:279] Creating Container Manager object based on Node Config: {RuntimeCgroupsName: SystemCgroupsName: KubeletCgroupsName: ContainerRuntime:docker CgroupsPerQOS:true CgroupRoot:/ CgroupDriver:cgroupfs KubeletRootDir:/opt/kubelet/ ProtectKernelDefaults:false NodeAllocatableConfig:{KubeReservedCgroupName: SystemReservedCgroupName: ReservedSystemCPUs: EnforceNodeAllocatable:map[pods:{}] KubeReserved:map[] SystemReserved:map[] HardEvictionThresholds:[{Signal:memory.available Operator:LessThan Value:{Quantity:100Mi Percentage:0} GracePeriod:0s MinReclaim:<nil>} {Signal:nodefs.available Operator:LessThan Value:{Quantity:<nil> Percentage:0.1} GracePeriod:0s MinReclaim:<nil>} {Signal:nodefs.inodesFree Operator:LessThan Value:{Quantity:<nil> Percentage:0.05} GracePeriod:0s MinReclaim:<nil>} {Signal:imagefs.available Operator:LessThan Value:{Quantity:<nil> Percentage:0.15} GracePeriod:0s MinReclaim:<nil>}]} QOSReserved:map[] ExperimentalCPUManagerPolicy:none ExperimentalTopologyManagerScope:container ExperimentalCPUManagerReconcilePeriod:10s ExperimentalPodPidsLimit:-1 EnforceCPULimits:true CPUCFSQuotaPeriod:100ms ExperimentalTopologyManagerPolicy:none}
I1230 18:52:44.305877    1410 topology_manager.go:120] [topologymanager] Creating topology manager with none policy per container scope
I1230 18:52:44.306800    1410 container_manager_linux.go:310] [topologymanager] Initializing Topology Manager with none policy and container-level scope
I1230 18:52:44.307383    1410 container_manager_linux.go:315] Creating device plugin manager: true
F1230 18:52:44.308305    1410 server.go:269] failed to run Kubelet: failed to initialize checkpoint manager: mkdir /var/lib/kubelet: permission denied
goroutine 1 [running]:
k8s.io/kubernetes/vendor/k8s.io/klog/v2.stacks(0xc000010001, 0xc0001e24e0, 0x97, 0xd0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:1026 +0xb9
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).output(0x70c9460, 0xc000000003, 0x0, 0x0, 0xc0007c8d20, 0x6f34162, 0x9, 0x10d, 0x411b00)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:975 +0x19b
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).printDepth(0x70c9460, 0xc000000003, 0x0, 0x0, 0x0, 0x0, 0x1, 0xc0007fc960, 0x1, 0x1)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:732 +0x16f
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).print(...)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:714
k8s.io/kubernetes/vendor/k8s.io/klog/v2.Fatal(...)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:1482
k8s.io/kubernetes/cmd/kubelet/app.NewKubeletCommand.func1(0xc0004858c0, 0xc00004e0a0, 0x8, 0x8)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/cmd/kubelet/app/server.go:269 +0x845
k8s.io/kubernetes/vendor/github.com/spf13/cobra.(*Command).execute(0xc0004858c0, 0xc00004e0a0, 0x8, 0x8, 0xc0004858c0, 0xc00004e0a0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/github.com/spf13/cobra/command.go:854 +0x2c2
k8s.io/kubernetes/vendor/github.com/spf13/cobra.(*Command).ExecuteC(0xc0004858c0, 0x165581414771f29b, 0x70c9020, 0x409b25)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/github.com/spf13/cobra/command.go:958 +0x375
k8s.io/kubernetes/vendor/github.com/spf13/cobra.(*Command).Execute(...)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/github.com/spf13/cobra/command.go:895
main.main()
        _output/dockerized/go/src/k8s.io/kubernetes/cmd/kubelet/kubelet.go:41 +0xe5

goroutine 6 [chan receive]:
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).flushDaemon(0x70c9460)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:1169 +0x8b
created by k8s.io/kubernetes/vendor/k8s.io/klog/v2.init.0
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:417 +0xdf

goroutine 109 [chan receive]:
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/watch.(*Broadcaster).loop(0xc000b45b40)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/watch/mux.go:219 +0x66
created by k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/watch.NewBroadcaster
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/watch/mux.go:73 +0xf7

goroutine 98 [chan receive]:
k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/server.SetupSignalContext.func1(0xc00071e930)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/server/signal.go:48 +0x36
created by k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/server.SetupSignalContext
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/server/signal.go:47 +0xf3

goroutine 97 [syscall]:
os/signal.signal_recv(0x0)
        /usr/local/go/src/runtime/sigqueue.go:147 +0x9d
os/signal.loop()
        /usr/local/go/src/os/signal/signal_unix.go:23 +0x25
created by os/signal.Notify.func1.1
        /usr/local/go/src/os/signal/signal.go:150 +0x45

goroutine 79 [select]:
k8s.io/kubernetes/vendor/go.opencensus.io/stats/view.(*worker).start(0xc00009b220)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/go.opencensus.io/stats/view/worker.go:154 +0x105
created by k8s.io/kubernetes/vendor/go.opencensus.io/stats/view.init.0
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/go.opencensus.io/stats/view/worker.go:32 +0x57

goroutine 94 [select]:
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.BackoffUntil(0x4a721f0, 0x4f0bde0, 0xc000394750, 0x1, 0xc0000a40c0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:167 +0x149
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.JitterUntil(0x4a721f0, 0x12a05f200, 0x0, 0xc000b50d01, 0xc0000a40c0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:133 +0x98
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.Until(...)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:90
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.Forever(0x4a721f0, 0x12a05f200)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:81 +0x4f
created by k8s.io/kubernetes/vendor/k8s.io/component-base/logs.InitLogs
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/component-base/logs/logs.go:58 +0x8a

goroutine 110 [chan receive]:
k8s.io/kubernetes/vendor/k8s.io/client-go/tools/record.(*eventBroadcasterImpl).StartEventWatcher.func1(0x4f1c9a0, 0xc000834c60, 0xc00084ff50)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/client-go/tools/record/event.go:301 +0xaa
created by k8s.io/kubernetes/vendor/k8s.io/client-go/tools/record.(*eventBroadcasterImpl).StartEventWatcher
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/client-go/tools/record/event.go:299 +0x6e
```

Still not fixed! A problem with kubernetes code? Maybe! Let's check it out :)
