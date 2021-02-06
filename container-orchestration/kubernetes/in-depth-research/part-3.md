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

Hmm. I can see that the path to some device manager stuff and device plugins is
actually hard coded. Hence the issue...

```bash
$ sudo mkdir -p /var/lib/kubelet
$ sudo chown ubuntu -R /var/lib/kubelet/
```

```bash
$ kubelet --root-dir /opt/kubelet/ --tls-cert-file kubelet-server.pem --tls-private-key-file kubelet-server-key.pem --seccomp-profile-root /opt/kubelet/seccomp
Flag --tls-cert-file has been deprecated, This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.
Flag --tls-private-key-file has been deprecated, This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.
Flag --seccomp-profile-root has been deprecated, will be removed in 1.23, in favor of using the `<root-dir>/seccomp` directory
I0206 15:18:38.895248    2011 server.go:416] Version: v1.20.1
W0206 15:18:38.896091    2011 server.go:558] standalone mode, no API client
W0206 15:18:38.896899    2011 server.go:614] failed to get the kubelet's cgroup: cpu and memory cgroup hierarchy not unified.  cpu: /, memory: /user.slice/user-1000.slice/session-5.scope.  Kubelet system container metrics may be missing.


W0206 15:18:43.949208    2011 server.go:473] No api server defined - no events will be sent to API server.
I0206 15:18:43.949688    2011 server.go:645] --cgroups-per-qos enabled, but --cgroup-root was not specified.  defaulting to /
I0206 15:18:43.951156    2011 container_manager_linux.go:274] container manager verified user specified cgroup-root exists: []
I0206 15:18:43.951583    2011 container_manager_linux.go:279] Creating Container Manager object based on Node Config: {RuntimeCgroupsName: SystemCgroupsName: KubeletCgroupsName: ContainerRuntime:docker CgroupsPerQOS:true CgroupRoot:/ CgroupDriver:cgroupfs KubeletRootDir:/opt/kubelet/ ProtectKernelDefaults:false NodeAllocatableConfig:{KubeReservedCgroupName: SystemReservedCgroupName: ReservedSystemCPUs: EnforceNodeAllocatable:map[pods:{}] KubeReserved:map[] SystemReserved:map[] HardEvictionThresholds:[{Signal:memory.available Operator:LessThan Value:{Quantity:100Mi Percentage:0} GracePeriod:0s MinReclaim:<nil>} {Signal:nodefs.available Operator:LessThan Value:{Quantity:<nil> Percentage:0.1} GracePeriod:0s MinReclaim:<nil>} {Signal:nodefs.inodesFree Operator:LessThan Value:{Quantity:<nil> Percentage:0.05} GracePeriod:0s MinReclaim:<nil>} {Signal:imagefs.available Operator:LessThan Value:{Quantity:<nil> Percentage:0.15} GracePeriod:0s MinReclaim:<nil>}]} QOSReserved:map[] ExperimentalCPUManagerPolicy:none ExperimentalTopologyManagerScope:container ExperimentalCPUManagerReconcilePeriod:10s ExperimentalPodPidsLimit:-1 EnforceCPULimits:true CPUCFSQuotaPeriod:100ms ExperimentalTopologyManagerPolicy:none}
I0206 15:18:43.952867    2011 topology_manager.go:120] [topologymanager] Creating topology manager with none policy per container scope
I0206 15:18:43.953431    2011 container_manager_linux.go:310] [topologymanager] Initializing Topology Manager with none policy and container-level scope
I0206 15:18:43.953857    2011 container_manager_linux.go:315] Creating device plugin manager: true
E0206 15:18:43.954508    2011 server.go:754] kubelet needs to run as uid `0`. It is being run as 1000
W0206 15:18:44.171438    2011 server.go:762] write /proc/self/oom_score_adj: permission denied
W0206 15:18:44.174160    2011 kubelet.go:297] Using dockershim is deprecated, please consider using a full-fledged CRI implementation
I0206 15:18:44.177841    2011 client.go:77] Connecting to docker on unix:///var/run/docker.sock
I0206 15:18:44.178480    2011 client.go:94] Start docker client with request timeout=2m0s
F0206 15:18:44.179542    2011 server.go:269] failed to run Kubelet: mkdir /var/lib/dockershim: permission denied
goroutine 1 [running]:
k8s.io/kubernetes/vendor/k8s.io/klog/v2.stacks(0xc000010001, 0xc0001e24e0, 0x71, 0xd0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:1026 +0xb9
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).output(0x70c9460, 0xc000000003, 0x0, 0x0, 0xc0007be380, 0x6f34162, 0x9, 0x10d, 0x411b00)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:975 +0x19b
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).printDepth(0x70c9460, 0xc000000003, 0x0, 0x0, 0x0, 0x0, 0x1, 0xc0007975c0, 0x1, 0x1)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:732 +0x16f
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).print(...)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:714
k8s.io/kubernetes/vendor/k8s.io/klog/v2.Fatal(...)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:1482
k8s.io/kubernetes/cmd/kubelet/app.NewKubeletCommand.func1(0xc000227b80, 0xc00004e0a0, 0x8, 0x8)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/cmd/kubelet/app/server.go:269 +0x845
k8s.io/kubernetes/vendor/github.com/spf13/cobra.(*Command).execute(0xc000227b80, 0xc00004e0a0, 0x8, 0x8, 0xc000227b80, 0xc00004e0a0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/github.com/spf13/cobra/command.go:854 +0x2c2
k8s.io/kubernetes/vendor/github.com/spf13/cobra.(*Command).ExecuteC(0xc000227b80, 0x16611fa138bbcf13, 0x70c9020, 0x409b25)
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
k8s.io/kubernetes/vendor/go.opencensus.io/stats/view.(*worker).start(0xc00009b220)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/go.opencensus.io/stats/view/worker.go:154 +0x105
created by k8s.io/kubernetes/vendor/go.opencensus.io/stats/view.init.0
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/go.opencensus.io/stats/view/worker.go:32 +0x57

goroutine 97 [syscall]:
os/signal.signal_recv(0x0)
        /usr/local/go/src/runtime/sigqueue.go:147 +0x9d
os/signal.loop()
        /usr/local/go/src/os/signal/signal_unix.go:23 +0x25
created by os/signal.Notify.func1.1
        /usr/local/go/src/os/signal/signal.go:150 +0x45

goroutine 94 [select]:
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.BackoffUntil(0x4a721f0, 0x4f0bde0, 0xc0004c6990, 0x1, 0xc0000a40c0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:167 +0x149
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.JitterUntil(0x4a721f0, 0x12a05f200, 0x0, 0xc00001ad01, 0xc0000a40c0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:133 +0x98
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.Until(...)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:90
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.Forever(0x4a721f0, 0x12a05f200)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:81 +0x4f
created by k8s.io/kubernetes/vendor/k8s.io/component-base/logs.InitLogs
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/component-base/logs/logs.go:58 +0x8a

goroutine 98 [chan receive]:
k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/server.SetupSignalContext.func1(0xc0007f2c20)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/server/signal.go:48 +0x36
created by k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/server.SetupSignalContext
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/server/signal.go:47 +0xf3

goroutine 109 [chan receive]:
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/watch.(*Broadcaster).loop(0xc000b31b00)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/watch/mux.go:219 +0x66
created by k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/watch.NewBroadcaster
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/watch/mux.go:73 +0xf7

goroutine 110 [chan receive]:
k8s.io/kubernetes/vendor/k8s.io/client-go/tools/record.(*eventBroadcasterImpl).StartEventWatcher.func1(0x4f1c9a0, 0xc0007992f0, 0xc000845f50)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/client-go/tools/record/event.go:301 +0xaa
created by k8s.io/kubernetes/vendor/k8s.io/client-go/tools/record.(*eventBroadcasterImpl).StartEventWatcher
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/client-go/tools/record/event.go:299 +0x6e

```

```bash
$ sudo mkdir -p /var/lib/dockershim
$ sudo chown ubuntu -R /var/lib/dockershim
```

```bash
$ kubelet --root-dir /opt/kubelet/ --tls-cert-file kubelet-server.pem --tls-private-key-file kubelet-server-key.pem --seccomp-profile-root /opt/kubelet/seccomp
Flag --tls-cert-file has been deprecated, This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.
Flag --tls-private-key-file has been deprecated, This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.
Flag --seccomp-profile-root has been deprecated, will be removed in 1.23, in favor of using the `<root-dir>/seccomp` directory
I0206 15:23:38.591595    2027 server.go:416] Version: v1.20.1
W0206 15:23:38.592556    2027 server.go:558] standalone mode, no API client
W0206 15:23:38.593434    2027 server.go:614] failed to get the kubelet's cgroup: cpu and memory cgroup hierarchy not unified.  cpu: /, memory: /user.slice/user-1000.slice/session-5.scope.  Kubelet system container metrics may be missing.
W0206 15:23:43.645614    2027 server.go:473] No api server defined - no events will be sent to API server.
I0206 15:23:43.646431    2027 server.go:645] --cgroups-per-qos enabled, but --cgroup-root was not specified.  defaulting to /
I0206 15:23:43.647997    2027 container_manager_linux.go:274] container manager verified user specified cgroup-root exists: []
I0206 15:23:43.648888    2027 container_manager_linux.go:279] Creating Container Manager object based on Node Config: {RuntimeCgroupsName: SystemCgroupsName: KubeletCgroupsName: ContainerRuntime:docker CgroupsPerQOS:true CgroupRoot:/ CgroupDriver:cgroupfs KubeletRootDir:/opt/kubelet/ ProtectKernelDefaults:false NodeAllocatableConfig:{KubeReservedCgroupName: SystemReservedCgroupName: ReservedSystemCPUs: EnforceNodeAllocatable:map[pods:{}] KubeReserved:map[] SystemReserved:map[] HardEvictionThresholds:[{Signal:memory.available Operator:LessThan Value:{Quantity:100Mi Percentage:0} GracePeriod:0s MinReclaim:<nil>} {Signal:nodefs.available Operator:LessThan Value:{Quantity:<nil> Percentage:0.1} GracePeriod:0s MinReclaim:<nil>} {Signal:nodefs.inodesFree Operator:LessThan Value:{Quantity:<nil> Percentage:0.05} GracePeriod:0s MinReclaim:<nil>} {Signal:imagefs.available Operator:LessThan Value:{Quantity:<nil> Percentage:0.15} GracePeriod:0s MinReclaim:<nil>}]} QOSReserved:map[] ExperimentalCPUManagerPolicy:none ExperimentalTopologyManagerScope:container ExperimentalCPUManagerReconcilePeriod:10s ExperimentalPodPidsLimit:-1 EnforceCPULimits:true CPUCFSQuotaPeriod:100ms ExperimentalTopologyManagerPolicy:none}
I0206 15:23:43.649663    2027 topology_manager.go:120] [topologymanager] Creating topology manager with none policy per container scope
I0206 15:23:43.650386    2027 container_manager_linux.go:310] [topologymanager] Initializing Topology Manager with none policy and container-level scope
I0206 15:23:43.651279    2027 container_manager_linux.go:315] Creating device plugin manager: true
E0206 15:23:43.652029    2027 server.go:754] kubelet needs to run as uid `0`. It is being run as 1000
W0206 15:23:43.908401    2027 server.go:762] write /proc/self/oom_score_adj: permission denied
W0206 15:23:43.908907    2027 kubelet.go:297] Using dockershim is deprecated, please consider using a full-fledged CRI implementation
I0206 15:23:43.909862    2027 client.go:77] Connecting to docker on unix:///var/run/docker.sock
I0206 15:23:43.911592    2027 client.go:94] Start docker client with request timeout=2m0s
F0206 15:23:43.915153    2027 server.go:269] failed to run Kubelet: failed to get docker version: Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?
goroutine 1 [running]:
k8s.io/kubernetes/vendor/k8s.io/klog/v2.stacks(0xc000010001, 0xc0001d84e0, 0xc4, 0xd0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:1026 +0xb9
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).output(0x70c9460, 0xc000000003, 0x0, 0x0, 0xc0007f65b0, 0x6f34162, 0x9, 0x10d, 0x411b00)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:975 +0x19b
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).printDepth(0x70c9460, 0xc000000003, 0x0, 0x0, 0x0, 0x0, 0x1, 0xc000793720, 0x1, 0x1)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:732 +0x16f
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).print(...)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:714
k8s.io/kubernetes/vendor/k8s.io/klog/v2.Fatal(...)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:1482
k8s.io/kubernetes/cmd/kubelet/app.NewKubeletCommand.func1(0xc0004858c0, 0xc00004e0a0, 0x8, 0x8)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/cmd/kubelet/app/server.go:269 +0x845
k8s.io/kubernetes/vendor/github.com/spf13/cobra.(*Command).execute(0xc0004858c0, 0xc00004e0a0, 0x8, 0x8, 0xc0004858c0, 0xc00004e0a0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/github.com/spf13/cobra/command.go:854 +0x2c2
k8s.io/kubernetes/vendor/github.com/spf13/cobra.(*Command).ExecuteC(0xc0004858c0, 0x16611fe6fe4ae5b3, 0x70c9020, 0x409b25)
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
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/watch.(*Broadcaster).loop(0xc000b31ac0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/watch/mux.go:219 +0x66
created by k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/watch.NewBroadcaster
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/watch/mux.go:73 +0xf7

goroutine 98 [chan receive]:
k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/server.SetupSignalContext.func1(0xc0007f0c20)
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
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.BackoffUntil(0x4a721f0, 0x4f0bde0, 0xc00038e630, 0x1, 0xc00009c0c0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:167 +0x149
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.JitterUntil(0x4a721f0, 0x12a05f200, 0x0, 0xc000b3ad01, 0xc00009c0c0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:133 +0x98
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.Until(...)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:90
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.Forever(0x4a721f0, 0x12a05f200)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:81 +0x4f
created by k8s.io/kubernetes/vendor/k8s.io/component-base/logs.InitLogs
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/component-base/logs/logs.go:58 +0x8a

goroutine 110 [chan receive]:
k8s.io/kubernetes/vendor/k8s.io/client-go/tools/record.(*eventBroadcasterImpl).StartEventWatcher.func1(0x4f1c9a0, 0xc000450d20, 0xc000843f30)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/client-go/tools/record/event.go:301 +0xaa
created by k8s.io/kubernetes/vendor/k8s.io/client-go/tools/record.(*eventBroadcasterImpl).StartEventWatcher
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/client-go/tools/record/event.go:299 +0x6e
```

```bash
$ docker

Command 'docker' not found, but can be installed with:

sudo snap install docker     # version 19.03.11, or
sudo apt  install docker.io  # version 19.03.8-0ubuntu1.20.04.1

See 'snap info docker' for additional versions.
```

```bash
$ sudo snap install docker
```

In the future I could try to use something like rkt, podman or cri-o for
container runtime :)

```bash
$ docker ps
Got permission denied while trying to connect to the Docker daemon socket at unix:///var/run/docker.sock: Get http://%2Fvar%2Frun%2Fdocker.sock/v1.40/containers/json: dial unix /var/run/docker.sock: connect: permission denied
$ sudo docker ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
```

```bash
$ kubelet --root-dir /opt/kubelet/ --tls-cert-file kubelet-server.pem --tls-private-key-file kubelet-server-key.pem --seccomp-profile-root /opt/kubelet/seccomp
Flag --tls-cert-file has been deprecated, This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.
Flag --tls-private-key-file has been deprecated, This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.
Flag --seccomp-profile-root has been deprecated, will be removed in 1.23, in favor of using the `<root-dir>/seccomp` directory
I0206 15:34:50.214506    3178 server.go:416] Version: v1.20.1
W0206 15:34:50.215710    3178 server.go:558] standalone mode, no API client
W0206 15:34:50.216799    3178 server.go:614] failed to get the kubelet's cgroup: cpu and memory cgroup hierarchy not unified.  cpu: /, memory: /user.slice/user-1000.slice/session-5.scope.  Kubelet system container metrics may be missing.
W0206 15:34:50.222476    3178 server.go:621] failed to get the container runtime's cgroup: failed to get container name for docker process: cpu and memory cgroup hierarchy not unified.  cpu: /, memory: /system.slice/snap.docker.dockerd.service. Runtime system container metrics may be missing.
W0206 15:34:55.427838    3178 server.go:473] No api server defined - no events will be sent to API server.
I0206 15:34:55.427884    3178 server.go:645] --cgroups-per-qos enabled, but --cgroup-root was not specified.  defaulting to /
I0206 15:34:55.428942    3178 container_manager_linux.go:274] container manager verified user specified cgroup-root exists: []
I0206 15:34:55.428989    3178 container_manager_linux.go:279] Creating Container Manager object based on Node Config: {RuntimeCgroupsName: SystemCgroupsName: KubeletCgroupsName: ContainerRuntime:docker CgroupsPerQOS:true CgroupRoot:/ CgroupDriver:cgroupfs KubeletRootDir:/opt/kubelet/ ProtectKernelDefaults:false NodeAllocatableConfig:{KubeReservedCgroupName: SystemReservedCgroupName: ReservedSystemCPUs: EnforceNodeAllocatable:map[pods:{}] KubeReserved:map[] SystemReserved:map[] HardEvictionThresholds:[{Signal:memory.available Operator:LessThan Value:{Quantity:100Mi Percentage:0} GracePeriod:0s MinReclaim:<nil>} {Signal:nodefs.available Operator:LessThan Value:{Quantity:<nil> Percentage:0.1} GracePeriod:0s MinReclaim:<nil>} {Signal:nodefs.inodesFree Operator:LessThan Value:{Quantity:<nil> Percentage:0.05} GracePeriod:0s MinReclaim:<nil>} {Signal:imagefs.available Operator:LessThan Value:{Quantity:<nil> Percentage:0.15} GracePeriod:0s MinReclaim:<nil>}]} QOSReserved:map[] ExperimentalCPUManagerPolicy:none ExperimentalTopologyManagerScope:container ExperimentalCPUManagerReconcilePeriod:10s ExperimentalPodPidsLimit:-1 EnforceCPULimits:true CPUCFSQuotaPeriod:100ms ExperimentalTopologyManagerPolicy:none}
I0206 15:34:55.430258    3178 topology_manager.go:120] [topologymanager] Creating topology manager with none policy per container scope
I0206 15:34:55.430969    3178 container_manager_linux.go:310] [topologymanager] Initializing Topology Manager with none policy and container-level scope
I0206 15:34:55.431538    3178 container_manager_linux.go:315] Creating device plugin manager: true
E0206 15:34:55.432304    3178 server.go:754] kubelet needs to run as uid `0`. It is being run as 1000
W0206 15:34:55.647479    3178 server.go:762] write /proc/self/oom_score_adj: permission denied
W0206 15:34:55.648565    3178 kubelet.go:297] Using dockershim is deprecated, please consider using a full-fledged CRI implementation
I0206 15:34:55.650041    3178 client.go:77] Connecting to docker on unix:///var/run/docker.sock
I0206 15:34:55.651206    3178 client.go:94] Start docker client with request timeout=2m0s
F0206 15:34:55.654897    3178 server.go:269] failed to run Kubelet: failed to get docker version: Got permission denied while trying to connect to the Docker daemon socket at unix:///var/run/docker.sock: Get "http://%2Fvar%2Frun%2Fdocker.sock/v1.24/version": dial unix /var/run/docker.sock: connect: permission denied
goroutine 1 [running]:
k8s.io/kubernetes/vendor/k8s.io/klog/v2.stacks(0xc000010001, 0xc000910000, 0x13e, 0x2b0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:1026 +0xb9
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).output(0x70c9460, 0xc000000003, 0x0, 0x0, 0xc0007bfb20, 0x6f34162, 0x9, 0x10d, 0x411b00)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:975 +0x19b
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).printDepth(0x70c9460, 0xc000000003, 0x0, 0x0, 0x0, 0x0, 0x1, 0xc0007f2cc0, 0x1, 0x1)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:732 +0x16f
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).print(...)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:714
k8s.io/kubernetes/vendor/k8s.io/klog/v2.Fatal(...)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:1482
k8s.io/kubernetes/cmd/kubelet/app.NewKubeletCommand.func1(0xc0003278c0, 0xc00004e0a0, 0x8, 0x8)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/cmd/kubelet/app/server.go:269 +0x845
k8s.io/kubernetes/vendor/github.com/spf13/cobra.(*Command).execute(0xc0003278c0, 0xc00004e0a0, 0x8, 0x8, 0xc0003278c0, 0xc00004e0a0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/github.com/spf13/cobra/command.go:854 +0x2c2
k8s.io/kubernetes/vendor/github.com/spf13/cobra.(*Command).ExecuteC(0xc0003278c0, 0x166120835de496c4, 0x70c9020, 0x409b25)
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
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/watch.(*Broadcaster).loop(0xc000b39b80)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/watch/mux.go:219 +0x66
created by k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/watch.NewBroadcaster
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/watch/mux.go:73 +0xf7

goroutine 98 [chan receive]:
k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/server.SetupSignalContext.func1(0xc000714930)
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
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.BackoffUntil(0x4a721f0, 0x4f0bde0, 0xc000396780, 0x1, 0xc0000a40c0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:167 +0x149
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.JitterUntil(0x4a721f0, 0x12a05f200, 0x0, 0xc000b42d01, 0xc0000a40c0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:133 +0x98
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.Until(...)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:90
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.Forever(0x4a721f0, 0x12a05f200)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:81 +0x4f
created by k8s.io/kubernetes/vendor/k8s.io/component-base/logs.InitLogs
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/component-base/logs/logs.go:58 +0x8a

goroutine 111 [chan receive]:
k8s.io/kubernetes/vendor/k8s.io/client-go/tools/record.(*eventBroadcasterImpl).StartEventWatcher.func1(0x4f1c9a0, 0xc000793260, 0xc000845ef0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/client-go/tools/record/event.go:301 +0xaa
created by k8s.io/kubernetes/vendor/k8s.io/client-go/tools/record.(*eventBroadcasterImpl).StartEventWatcher
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/client-go/tools/record/event.go:299 +0x6e
```

```bash
$ ls -al /var/run/docker.sock
srw-rw---- 1 root root 0 Feb  6 15:25 /var/run/docker.sock
$ sudo chown ubuntu /var/run/docker.sock
```

```bash
$ kubelet --root-dir /opt/kubelet/ --tls-cert-file kubelet-server.pem --tls-private-key-file kubelet-server-key.pem --seccomp-profile-root /opt/kubelet/seccomp
Flag --tls-cert-file has been deprecated, This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.
Flag --tls-private-key-file has been deprecated, This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.
Flag --seccomp-profile-root has been deprecated, will be removed in 1.23, in favor of using the `<root-dir>/seccomp` directory
I0206 15:37:20.110147    3209 server.go:416] Version: v1.20.1
W0206 15:37:20.111317    3209 server.go:558] standalone mode, no API client
W0206 15:37:20.112138    3209 server.go:614] failed to get the kubelet's cgroup: cpu and memory cgroup hierarchy not unified.  cpu: /, memory: /user.slice/user-1000.slice/session-5.scope.  Kubelet system container metrics may be missing.
W0206 15:37:20.115536    3209 server.go:621] failed to get the container runtime's cgroup: failed to get container name for docker process: cpu and memory cgroup hierarchy not unified.  cpu: /, memory: /system.slice/snap.docker.dockerd.service. Runtime system container metrics may be missing.

W0206 15:37:20.682711    3209 server.go:473] No api server defined - no events will be sent to API server.
I0206 15:37:20.682754    3209 server.go:645] --cgroups-per-qos enabled, but --cgroup-root was not specified.  defaulting to /
I0206 15:37:20.683086    3209 container_manager_linux.go:274] container manager verified user specified cgroup-root exists: []
I0206 15:37:20.683129    3209 container_manager_linux.go:279] Creating Container Manager object based on Node Config: {RuntimeCgroupsName: SystemCgroupsName: KubeletCgroupsName: ContainerRuntime:docker CgroupsPerQOS:true CgroupRoot:/ CgroupDriver:cgroupfs KubeletRootDir:/opt/kubelet/ ProtectKernelDefaults:false NodeAllocatableConfig:{KubeReservedCgroupName: SystemReservedCgroupName: ReservedSystemCPUs: EnforceNodeAllocatable:map[pods:{}] KubeReserved:map[] SystemReserved:map[] HardEvictionThresholds:[{Signal:nodefs.inodesFree Operator:LessThan Value:{Quantity:<nil> Percentage:0.05} GracePeriod:0s MinReclaim:<nil>} {Signal:imagefs.available Operator:LessThan Value:{Quantity:<nil> Percentage:0.15} GracePeriod:0s MinReclaim:<nil>} {Signal:memory.available Operator:LessThan Value:{Quantity:100Mi Percentage:0} GracePeriod:0s MinReclaim:<nil>} {Signal:nodefs.available Operator:LessThan Value:{Quantity:<nil> Percentage:0.1} GracePeriod:0s MinReclaim:<nil>}]} QOSReserved:map[] ExperimentalCPUManagerPolicy:none ExperimentalTopologyManagerScope:container ExperimentalCPUManagerReconcilePeriod:10s ExperimentalPodPidsLimit:-1 EnforceCPULimits:true CPUCFSQuotaPeriod:100ms ExperimentalTopologyManagerPolicy:none}
I0206 15:37:20.683226    3209 topology_manager.go:120] [topologymanager] Creating topology manager with none policy per container scope
I0206 15:37:20.683254    3209 container_manager_linux.go:310] [topologymanager] Initializing Topology Manager with none policy and container-level scope
I0206 15:37:20.683277    3209 container_manager_linux.go:315] Creating device plugin manager: true
E0206 15:37:20.683330    3209 server.go:754] kubelet needs to run as uid `0`. It is being run as 1000
W0206 15:37:20.910290    3209 server.go:762] write /proc/self/oom_score_adj: permission denied
W0206 15:37:20.910394    3209 kubelet.go:297] Using dockershim is deprecated, please consider using a full-fledged CRI implementation
I0206 15:37:20.910552    3209 client.go:77] Connecting to docker on unix:///var/run/docker.sock
I0206 15:37:20.910607    3209 client.go:94] Start docker client with request timeout=2m0s
W0206 15:37:20.934899    3209 docker_service.go:559] Hairpin mode set to "promiscuous-bridge" but kubenet is not enabled, falling back to "hairpin-veth"
I0206 15:37:20.934948    3209 docker_service.go:240] Hairpin mode set to "hairpin-veth"
W0206 15:37:20.937817    3209 cni.go:239] Unable to update cni config: no networks found in /etc/cni/net.d
W0206 15:37:20.950135    3209 hostport_manager.go:71] The binary conntrack is not installed, this can cause failures in network connection cleanup.
W0206 15:37:20.950592    3209 hostport_manager.go:71] The binary conntrack is not installed, this can cause failures in network connection cleanup.
W0206 15:37:20.960442    3209 plugins.go:195] can't set sysctl net/bridge/bridge-nf-call-iptables: open /proc/sys/net/bridge/bridge-nf-call-iptables: permission denied
I0206 15:37:20.960961    3209 docker_service.go:255] Docker cri networking managed by kubernetes.io/no-op
I0206 15:37:20.976748    3209 docker_service.go:260] Docker Info: &{ID:OUWO:QIKR:KLYA:YHJZ:5VDA:7NBA:RME2:MVHU:2QSS:2PBR:F2DA:T6BM Containers:0 ContainersRunning:0 ContainersPaused:0 ContainersStopped:0 Images:0 Driver:overlay2 DriverStatus:[[Backing Filesystem extfs] [Supports d_type true] [Native Overlay Diff true]] SystemStatus:[] Plugins:{Volume:[local] Network:[bridge host ipvlan macvlan null overlay] Authorization:[] Log:[awslogs fluentd gcplogs gelf journald json-file local logentries splunk syslog]} MemoryLimit:true SwapLimit:false KernelMemory:true KernelMemoryTCP:true CPUCfsPeriod:true CPUCfsQuota:true CPUShares:true CPUSet:true PidsLimit:true IPv4Forwarding:true BridgeNfIptables:true BridgeNfIP6tables:true Debug:false NFd:23 OomKillDisable:true NGoroutines:41 SystemTime:2021-02-06T15:37:20.963370342+05:30 LoggingDriver:json-file CgroupDriver:cgroupfs NEventsListener:0 KernelVersion:5.4.0-54-generic OperatingSystem:Ubuntu Core 16 OSType:linux Architecture:x86_64 IndexServerAddress:https://index.docker.io/v1/ RegistryConfig:0xc0007fbdc0 NCPU:1 MemTotal:1029029888 GenericResources:[] DockerRootDir:/var/snap/docker/common/var-lib-docker HTTPProxy: HTTPSProxy: NoProxy: Name:my-own-k8s-cluster Labels:[] ExperimentalBuild:false ServerVersion:19.03.11 ClusterStore: ClusterAdvertise: Runtimes:map[runc:{Path:runc Args:[]}] DefaultRuntime:runc Swarm:{NodeID: NodeAddr: LocalNodeState:inactive ControlAvailable:false Error: RemoteManagers:[] Nodes:0 Managers:0 Cluster:<nil> Warnings:[]} LiveRestoreEnabled:false Isolation: InitBinary:docker-init ContainerdCommit:{ID:7ad184331fa3e55e52b890ea95e65ba581ae3429 Expected:7ad184331fa3e55e52b890ea95e65ba581ae3429} RuncCommit:{ID: Expected:} InitCommit:{ID:fec3683 Expected:fec3683} SecurityOptions:[name=apparmor name=seccomp,profile=default] ProductLicense: Warnings:[WARNING: No swap limit support]}
I0206 15:37:20.976874    3209 docker_service.go:273] Setting cgroupDriver to cgroupfs
F0206 15:37:20.977570    3209 server.go:269] failed to run Kubelet: failed to listen on "unix:///var/run/dockershim.sock": failed to create temporary file: open /var/run/927772577: permission denied
goroutine 1 [running]:
k8s.io/kubernetes/vendor/k8s.io/klog/v2.stacks(0xc000010001, 0xc000b7e1e0, 0xc7, 0x1d9)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:1026 +0xb9
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).output(0x70c9460, 0xc000000003, 0x0, 0x0, 0xc0007fbea0, 0x6f34162, 0x9, 0x10d, 0x411b00)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:975 +0x19b
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).printDepth(0x70c9460, 0xc000000003, 0x0, 0x0, 0x0, 0x0, 0x1, 0xc00037e880, 0x1, 0x1)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:732 +0x16f
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).print(...)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:714
k8s.io/kubernetes/vendor/k8s.io/klog/v2.Fatal(...)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:1482
k8s.io/kubernetes/cmd/kubelet/app.NewKubeletCommand.func1(0xc0004918c0, 0xc00004e0a0, 0x8, 0x8)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/cmd/kubelet/app/server.go:269 +0x845
k8s.io/kubernetes/vendor/github.com/spf13/cobra.(*Command).execute(0xc0004918c0, 0xc00004e0a0, 0x8, 0x8, 0xc0004918c0, 0xc00004e0a0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/github.com/spf13/cobra/command.go:854 +0x2c2
k8s.io/kubernetes/vendor/github.com/spf13/cobra.(*Command).ExecuteC(0xc0004918c0, 0x166120a64570dc5f, 0x70c9020, 0x409b25)
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

goroutine 103 [chan receive]:
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/watch.(*Broadcaster).loop(0xc000790540)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/watch/mux.go:219 +0x66
created by k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/watch.NewBroadcaster
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/watch/mux.go:73 +0xf7

goroutine 98 [chan receive]:
k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/server.SetupSignalContext.func1(0xc0007f4c20)
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
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.BackoffUntil(0x4a721f0, 0x4f0bde0, 0xc00038e6c0, 0x1, 0xc00009c0c0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:167 +0x149
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.JitterUntil(0x4a721f0, 0x12a05f200, 0x0, 0xc000b40d01, 0xc00009c0c0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:133 +0x98
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.Until(...)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:90
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.Forever(0x4a721f0, 0x12a05f200)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:81 +0x4f
created by k8s.io/kubernetes/vendor/k8s.io/component-base/logs.InitLogs
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/component-base/logs/logs.go:58 +0x8a

goroutine 101 [IO wait]:
internal/poll.runtime_pollWait(0x7f0981e8de40, 0x72, 0x4f10000)
        /usr/local/go/src/runtime/netpoll.go:222 +0x55
internal/poll.(*pollDesc).wait(0xc00070d818, 0x72, 0x4f10000, 0x6fd7e28, 0x0)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:87 +0x45
internal/poll.(*pollDesc).waitRead(...)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:92
internal/poll.(*FD).Read(0xc00070d800, 0xc0007d2000, 0x1000, 0x1000, 0x0, 0x0, 0x0)
        /usr/local/go/src/internal/poll/fd_unix.go:159 +0x1a5
net.(*netFD).Read(0xc00070d800, 0xc0007d2000, 0x1000, 0x1000, 0x43de5c, 0xc000083b58, 0x46b2e0)
        /usr/local/go/src/net/fd_posix.go:55 +0x4f
net.(*conn).Read(0xc000acf7c8, 0xc0007d2000, 0x1000, 0x1000, 0x0, 0x0, 0x0)
        /usr/local/go/src/net/net.go:182 +0x8e
net/http.(*persistConn).Read(0xc000b7d680, 0xc0007d2000, 0x1000, 0x1000, 0xc0006f2ea0, 0xc000083c58, 0x409095)
        /usr/local/go/src/net/http/transport.go:1887 +0x77
bufio.(*Reader).fill(0xc0003fc060)
        /usr/local/go/src/bufio/bufio.go:101 +0x105
bufio.(*Reader).Peek(0xc0003fc060, 0x1, 0x0, 0x0, 0x1, 0x0, 0xc0006f39e0)
        /usr/local/go/src/bufio/bufio.go:139 +0x4f
net/http.(*persistConn).readLoop(0xc000b7d680)
        /usr/local/go/src/net/http/transport.go:2040 +0x1a8
created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1708 +0xcb7

goroutine 102 [select]:
net/http.(*persistConn).writeLoop(0xc000b7d680)
        /usr/local/go/src/net/http/transport.go:2340 +0x11c
created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1709 +0xcdc

goroutine 104 [chan receive]:
k8s.io/kubernetes/vendor/k8s.io/client-go/tools/record.(*eventBroadcasterImpl).StartEventWatcher.func1(0x4f1c9a0, 0xc0008a4960, 0xc000890850)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/client-go/tools/record/event.go:301 +0xaa
created by k8s.io/kubernetes/vendor/k8s.io/client-go/tools/record.(*eventBroadcasterImpl).StartEventWatcher
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/client-go/tools/record/event.go:299 +0x6e

goroutine 118 [runnable]:
k8s.io/kubernetes/pkg/kubelet/dockershim.(*dockerService).Start.func1(0xc0001dc3c0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/pkg/kubelet/dockershim/docker_service.go:408
created by k8s.io/kubernetes/pkg/kubelet/dockershim.(*dockerService).Start
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/pkg/kubelet/dockershim/docker_service.go:408 +0x4d

goroutine 106 [IO wait]:
internal/poll.runtime_pollWait(0x7f0981e8dd58, 0x72, 0x4f10000)
        /usr/local/go/src/runtime/netpoll.go:222 +0x55
internal/poll.(*pollDesc).wait(0xc00067ae18, 0x72, 0x4f10000, 0x6fd7e28, 0x0)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:87 +0x45
internal/poll.(*pollDesc).waitRead(...)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:92
internal/poll.(*FD).Read(0xc00067ae00, 0xc000bca000, 0x1000, 0x1000, 0x0, 0x0, 0x0)
        /usr/local/go/src/internal/poll/fd_unix.go:159 +0x1a5
net.(*netFD).Read(0xc00067ae00, 0xc000bca000, 0x1000, 0x1000, 0x43de5c, 0xc000086b58, 0x46b2e0)
        /usr/local/go/src/net/fd_posix.go:55 +0x4f
net.(*conn).Read(0xc000acffa0, 0xc000bca000, 0x1000, 0x1000, 0x0, 0x0, 0x0)
        /usr/local/go/src/net/net.go:182 +0x8e
net/http.(*persistConn).Read(0xc000a3a360, 0xc000bca000, 0x1000, 0x1000, 0xc00009d560, 0xc000086c58, 0x409095)
        /usr/local/go/src/net/http/transport.go:1887 +0x77
bufio.(*Reader).fill(0xc00088dc20)
        /usr/local/go/src/bufio/bufio.go:101 +0x105
bufio.(*Reader).Peek(0xc00088dc20, 0x1, 0x0, 0x0, 0x1, 0x0, 0xc00009db60)
        /usr/local/go/src/bufio/bufio.go:139 +0x4f
net/http.(*persistConn).readLoop(0xc000a3a360)
        /usr/local/go/src/net/http/transport.go:2040 +0x1a8
created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1708 +0xcb7

goroutine 107 [select]:
net/http.(*persistConn).writeLoop(0xc000a3a360)
        /usr/local/go/src/net/http/transport.go:2340 +0x11c
created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1709 +0xcdc

goroutine 119 [runnable]:
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.Until(0xc00037e6f0, 0x45d964b800, 0xc00009c0c0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:89
created by k8s.io/kubernetes/pkg/kubelet/dockershim/cm.(*containerManager).Start
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/pkg/kubelet/dockershim/cm/container_manager_linux.go:84 +0xa5
```

```bash
$ sudo ln -s ~/kubernetes/server/bin/kubelet /usr/local/bin/
```

```bash
$ sudo kubelet
I0206 15:40:27.165141    3454 server.go:416] Version: v1.20.1
W0206 15:40:27.166681    3454 server.go:558] standalone mode, no API client
W0206 15:40:27.167281    3454 server.go:614] failed to get the kubelet's cgroup: cpu and memory cgroup hierarchy not unified.  cpu: /, memory: /user.slice/user-1000.slice/session-5.scope.  Kubelet system container metrics may be missing.
W0206 15:40:27.173054    3454 server.go:621] failed to get the container runtime's cgroup: failed to get container name for docker process: cpu and memory cgroup hierarchy not unified.  cpu: /, memory: /system.slice/snap.docker.dockerd.service. Runtime system container metrics may be missing.
W0206 15:40:27.360358    3454 server.go:473] No api server defined - no events will be sent to API server.
I0206 15:40:27.360403    3454 server.go:645] --cgroups-per-qos enabled, but --cgroup-root was not specified.  defaulting to /
I0206 15:40:27.360734    3454 container_manager_linux.go:274] container manager verified user specified cgroup-root exists: []
I0206 15:40:27.360773    3454 container_manager_linux.go:279] Creating Container Manager object based on Node Config: {RuntimeCgroupsName: SystemCgroupsName: KubeletCgroupsName: ContainerRuntime:docker CgroupsPerQOS:true CgroupRoot:/ CgroupDriver:cgroupfs KubeletRootDir:/var/lib/kubelet ProtectKernelDefaults:false NodeAllocatableConfig:{KubeReservedCgroupName: SystemReservedCgroupName: ReservedSystemCPUs: EnforceNodeAllocatable:map[pods:{}] KubeReserved:map[] SystemReserved:map[] HardEvictionThresholds:[{Signal:nodefs.available Operator:LessThan Value:{Quantity:<nil> Percentage:0.1} GracePeriod:0s MinReclaim:<nil>} {Signal:nodefs.inodesFree Operator:LessThan Value:{Quantity:<nil> Percentage:0.05} GracePeriod:0s MinReclaim:<nil>} {Signal:imagefs.available Operator:LessThan Value:{Quantity:<nil> Percentage:0.15} GracePeriod:0s MinReclaim:<nil>} {Signal:memory.available Operator:LessThan Value:{Quantity:100Mi Percentage:0} GracePeriod:0s MinReclaim:<nil>}]} QOSReserved:map[] ExperimentalCPUManagerPolicy:none ExperimentalTopologyManagerScope:container ExperimentalCPUManagerReconcilePeriod:10s ExperimentalPodPidsLimit:-1 EnforceCPULimits:true CPUCFSQuotaPeriod:100ms ExperimentalTopologyManagerPolicy:none}
I0206 15:40:27.360864    3454 topology_manager.go:120] [topologymanager] Creating topology manager with none policy per container scope
I0206 15:40:27.360890    3454 container_manager_linux.go:310] [topologymanager] Initializing Topology Manager with none policy and container-level scope
I0206 15:40:27.360912    3454 container_manager_linux.go:315] Creating device plugin manager: true
W0206 15:40:27.361011    3454 kubelet.go:297] Using dockershim is deprecated, please consider using a full-fledged CRI implementation
I0206 15:40:27.361045    3454 client.go:77] Connecting to docker on unix:///var/run/docker.sock
I0206 15:40:27.361072    3454 client.go:94] Start docker client with request timeout=2m0s
W0206 15:40:27.372674    3454 docker_service.go:559] Hairpin mode set to "promiscuous-bridge" but kubenet is not enabled, falling back to "hairpin-veth"
I0206 15:40:27.372731    3454 docker_service.go:240] Hairpin mode set to "hairpin-veth"
W0206 15:40:27.372836    3454 cni.go:239] Unable to update cni config: no networks found in /etc/cni/net.d
W0206 15:40:27.377264    3454 hostport_manager.go:71] The binary conntrack is not installed, this can cause failures in network connection cleanup.
W0206 15:40:27.377689    3454 hostport_manager.go:71] The binary conntrack is not installed, this can cause failures in network connection cleanup.
I0206 15:40:27.380597    3454 docker_service.go:255] Docker cri networking managed by kubernetes.io/no-op
I0206 15:40:27.393970    3454 docker_service.go:260] Docker Info: &{ID:OUWO:QIKR:KLYA:YHJZ:5VDA:7NBA:RME2:MVHU:2QSS:2PBR:F2DA:T6BM Containers:0 ContainersRunning:0 ContainersPaused:0 ContainersStopped:0 Images:0 Driver:overlay2 DriverStatus:[[Backing Filesystem extfs] [Supports d_type true] [Native Overlay Diff true]] SystemStatus:[] Plugins:{Volume:[local] Network:[bridge host ipvlan macvlan null overlay] Authorization:[] Log:[awslogs fluentd gcplogs gelf journald json-file local logentries splunk syslog]} MemoryLimit:true SwapLimit:false KernelMemory:true KernelMemoryTCP:true CPUCfsPeriod:true CPUCfsQuota:true CPUShares:true CPUSet:true PidsLimit:true IPv4Forwarding:true BridgeNfIptables:true BridgeNfIP6tables:true Debug:false NFd:23 OomKillDisable:true NGoroutines:41 SystemTime:2021-02-06T15:40:27.383349783+05:30 LoggingDriver:json-file CgroupDriver:cgroupfs NEventsListener:0 KernelVersion:5.4.0-54-generic OperatingSystem:Ubuntu Core 16 OSType:linux Architecture:x86_64 IndexServerAddress:https://index.docker.io/v1/ RegistryConfig:0xc00083b880 NCPU:1 MemTotal:1029029888 GenericResources:[] DockerRootDir:/var/snap/docker/common/var-lib-docker HTTPProxy: HTTPSProxy: NoProxy: Name:my-own-k8s-cluster Labels:[] ExperimentalBuild:false ServerVersion:19.03.11 ClusterStore: ClusterAdvertise: Runtimes:map[runc:{Path:runc Args:[]}] DefaultRuntime:runc Swarm:{NodeID: NodeAddr: LocalNodeState:inactive ControlAvailable:false Error: RemoteManagers:[] Nodes:0 Managers:0 Cluster:<nil> Warnings:[]} LiveRestoreEnabled:false Isolation: InitBinary:docker-init ContainerdCommit:{ID:7ad184331fa3e55e52b890ea95e65ba581ae3429 Expected:7ad184331fa3e55e52b890ea95e65ba581ae3429} RuncCommit:{ID: Expected:} InitCommit:{ID:fec3683 Expected:fec3683} SecurityOptions:[name=apparmor name=seccomp,profile=default] ProductLicense: Warnings:[WARNING: No swap limit support]}
I0206 15:40:27.394147    3454 docker_service.go:273] Setting cgroupDriver to cgroupfs
I0206 15:40:27.427955    3454 remote_runtime.go:62] parsed scheme: ""
I0206 15:40:27.428599    3454 remote_runtime.go:62] scheme "" not registered, fallback to default scheme
I0206 15:40:27.429400    3454 passthrough.go:48] ccResolverWrapper: sending update to cc: {[{/var/run/dockershim.sock  <nil> 0 <nil>}] <nil> <nil>}
I0206 15:40:27.429904    3454 clientconn.go:948] ClientConn switching balancer to "pick_first"
I0206 15:40:27.430642    3454 remote_image.go:50] parsed scheme: ""
I0206 15:40:27.432337    3454 remote_image.go:50] scheme "" not registered, fallback to default scheme
I0206 15:40:27.432938    3454 passthrough.go:48] ccResolverWrapper: sending update to cc: {[{/var/run/dockershim.sock  <nil> 0 <nil>}] <nil> <nil>}
I0206 15:40:27.433630    3454 clientconn.go:948] ClientConn switching balancer to "pick_first"
I0206 15:40:27.473058    3454 kuberuntime_manager.go:216] Container runtime docker initialized, version: 19.03.11, apiVersion: 1.40.0


E0206 15:40:30.788188    3454 aws_credentials.go:77] while getting AWS credentials NoCredentialProviders: no valid providers in chain. Deprecated.
        For verbose messaging see aws.Config.CredentialsChainVerboseErrors
W0206 15:40:30.790341    3454 volume_host.go:75] kubeClient is nil. Skip initialization of CSIDriverLister
W0206 15:40:30.793695    3454 probe.go:268] Flexvolume plugin directory at /usr/libexec/kubernetes/kubelet-plugins/volume/exec/ does not exist. Recreating.
W0206 15:40:30.796436    3454 csi_plugin.go:190] kubernetes.io/csi: kubeclient not set, assuming standalone kubelet
W0206 15:40:30.796968    3454 csi_plugin.go:264] Skipping CSINode initialization, kubelet running in standalone mode
I0206 15:40:30.798744    3454 server.go:1176] Started kubelet
E0206 15:40:30.801629    3454 kubelet.go:1271] Image garbage collection failed once. Stats initialization may not have completed yet: failed to get imageFs info: unable to find data in memory cache
W0206 15:40:30.801701    3454 kubelet.go:1376] No api server defined - no node status update will be sent.
I0206 15:40:30.801773    3454 server.go:148] Starting to listen on 0.0.0.0:10250
I0206 15:40:30.817256    3454 server.go:409] Adding debug handlers to kubelet server.
I0206 15:40:30.819833    3454 fs_resource_analyzer.go:64] Starting FS ResourceAnalyzer
I0206 15:40:30.831375    3454 volume_manager.go:271] Starting Kubelet Volume Manager
I0206 15:40:30.832385    3454 desired_state_of_world_populator.go:142] Desired state populator starts to run
I0206 15:40:30.933158    3454 reconciler.go:157] Reconciler: start to sync state
I0206 15:40:31.122800    3454 kubelet_network_linux.go:56] Initialized IPv4 iptables rules.
I0206 15:40:31.123565    3454 status_manager.go:154] Kubernetes client is nil, not starting status manager.
I0206 15:40:31.124512    3454 kubelet.go:1799] Starting kubelet main sync loop.
E0206 15:40:31.128483    3454 kubelet.go:1823] skipping pod synchronization - [container runtime status check may not have completed yet, PLEG is not healthy: pleg has yet to be successful]
I0206 15:40:31.171152    3454 cpu_manager.go:193] [cpumanager] starting with none policy
I0206 15:40:31.172245    3454 cpu_manager.go:194] [cpumanager] reconciling every 10s
I0206 15:40:31.172810    3454 state_mem.go:36] [cpumanager] initializing new in-memory state store
I0206 15:40:31.176131    3454 policy_none.go:43] [cpumanager] none policy: Start
W0206 15:40:31.193034    3454 manager.go:594] Failed to retrieve checkpoint for "kubelet_internal_checkpoint": checkpoint is not found
I0206 15:40:31.194315    3454 plugin_manager.go:114] Starting Kubelet Plugin Manager
E0206 15:40:31.198973    3454 container_manager_linux.go:487] cpu and memory cgroup hierarchy not unified.  cpu: /, memory: /system.slice/snap.docker.dockerd.service
E0206 15:40:31.210836    3454 container_manager_linux.go:533] failed to find cgroups of kubelet - cpu and memory cgroup hierarchy not unified.  cpu: /, memory: /user.slice/user-1000.slice/session-5.scope
W0206 15:40:31.262283    3454 eviction_manager.go:344] eviction manager: attempting to reclaim ephemeral-storage
I0206 15:40:31.262926    3454 container_gc.go:85] attempting to delete unused containers
I0206 15:40:31.268759    3454 image_gc_manager.go:321] attempting to delete unused images
I0206 15:40:31.347812    3454 eviction_manager.go:355] eviction manager: must evict pod(s) to reclaim ephemeral-storage
E0206 15:40:31.348415    3454 eviction_manager.go:366] eviction manager: eviction thresholds have been met, but no pods are active to evict
W0206 15:40:41.400317    3454 eviction_manager.go:344] eviction manager: attempting to reclaim ephemeral-storage
I0206 15:40:41.401292    3454 container_gc.go:85] attempting to delete unused containers
I0206 15:40:41.407370    3454 image_gc_manager.go:321] attempting to delete unused images
I0206 15:40:41.478255    3454 eviction_manager.go:355] eviction manager: must evict pod(s) to reclaim ephemeral-storage
E0206 15:40:41.478788    3454 eviction_manager.go:366] eviction manager: eviction thresholds have been met, but no pods are active to evict
W0206 15:40:51.550415    3454 eviction_manager.go:344] eviction manager: attempting to reclaim ephemeral-storage
I0206 15:40:51.551534    3454 container_gc.go:85] attempting to delete unused containers
I0206 15:40:51.556923    3454 image_gc_manager.go:321] attempting to delete unused images
I0206 15:40:51.627249    3454 eviction_manager.go:355] eviction manager: must evict pod(s) to reclaim ephemeral-storage
E0206 15:40:51.627784    3454 eviction_manager.go:366] eviction manager: eviction thresholds have been met, but no pods are active to evict
W0206 15:41:01.669749    3454 eviction_manager.go:344] eviction manager: attempting to reclaim ephemeral-storage
I0206 15:41:01.670325    3454 container_gc.go:85] attempting to delete unused containers
I0206 15:41:01.675869    3454 image_gc_manager.go:321] attempting to delete unused images
I0206 15:41:01.764347    3454 eviction_manager.go:355] eviction manager: must evict pod(s) to reclaim ephemeral-storage
E0206 15:41:01.764946    3454 eviction_manager.go:366] eviction manager: eviction thresholds have been met, but no pods are active to evict
```

It also works with the options I had before.

```bash
$ sudo kubelet --root-dir /opt/kubelet/ --tls-cert-file kubelet-server.pem --tls-private-key-file kubelet-server-key.pem --seccomp-profile-root /opt/kubelet/seccomp
Flag --tls-cert-file has been deprecated, This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.
Flag --tls-private-key-file has been deprecated, This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.
Flag --seccomp-profile-root has been deprecated, will be removed in 1.23, in favor of using the `<root-dir>/seccomp` directory
I0206 15:43:19.541981    4213 server.go:416] Version: v1.20.1
W0206 15:43:19.542418    4213 server.go:558] standalone mode, no API client
W0206 15:43:19.542559    4213 server.go:614] failed to get the kubelet's cgroup: cpu and memory cgroup hierarchy not unified.  cpu: /, memory: /user.slice/user-1000.slice/session-5.scope.  Kubelet system container metrics may be missing.
W0206 15:43:19.544616    4213 server.go:621] failed to get the container runtime's cgroup: failed to get container name for docker process: cpu and memory cgroup hierarchy not unified.  cpu: /, memory: /system.slice/snap.docker.dockerd.service. Runtime system container metrics may be missing.
W0206 15:43:19.656358    4213 server.go:473] No api server defined - no events will be sent to API server.
I0206 15:43:19.656401    4213 server.go:645] --cgroups-per-qos enabled, but --cgroup-root was not specified.  defaulting to /
I0206 15:43:19.656844    4213 container_manager_linux.go:274] container manager verified user specified cgroup-root exists: []
I0206 15:43:19.656994    4213 container_manager_linux.go:279] Creating Container Manager object based on Node Config: {RuntimeCgroupsName: SystemCgroupsName: KubeletCgroupsName: ContainerRuntime:docker CgroupsPerQOS:true CgroupRoot:/ CgroupDriver:cgroupfs KubeletRootDir:/opt/kubelet/ ProtectKernelDefaults:false NodeAllocatableConfig:{KubeReservedCgroupName: SystemReservedCgroupName: ReservedSystemCPUs: EnforceNodeAllocatable:map[pods:{}] KubeReserved:map[] SystemReserved:map[] HardEvictionThresholds:[{Signal:memory.available Operator:LessThan Value:{Quantity:100Mi Percentage:0} GracePeriod:0s MinReclaim:<nil>} {Signal:nodefs.available Operator:LessThan Value:{Quantity:<nil> Percentage:0.1} GracePeriod:0s MinReclaim:<nil>} {Signal:nodefs.inodesFree Operator:LessThan Value:{Quantity:<nil> Percentage:0.05} GracePeriod:0s MinReclaim:<nil>} {Signal:imagefs.available Operator:LessThan Value:{Quantity:<nil> Percentage:0.15} GracePeriod:0s MinReclaim:<nil>}]} QOSReserved:map[] ExperimentalCPUManagerPolicy:none ExperimentalTopologyManagerScope:container ExperimentalCPUManagerReconcilePeriod:10s ExperimentalPodPidsLimit:-1 EnforceCPULimits:true CPUCFSQuotaPeriod:100ms ExperimentalTopologyManagerPolicy:none}
I0206 15:43:19.658451    4213 topology_manager.go:120] [topologymanager] Creating topology manager with none policy per container scope
I0206 15:43:19.659290    4213 container_manager_linux.go:310] [topologymanager] Initializing Topology Manager with none policy and container-level scope
I0206 15:43:19.660241    4213 container_manager_linux.go:315] Creating device plugin manager: true
W0206 15:43:19.661309    4213 kubelet.go:297] Using dockershim is deprecated, please consider using a full-fledged CRI implementation
I0206 15:43:19.662312    4213 client.go:77] Connecting to docker on unix:///var/run/docker.sock
I0206 15:43:19.663215    4213 client.go:94] Start docker client with request timeout=2m0s
W0206 15:43:19.675786    4213 docker_service.go:559] Hairpin mode set to "promiscuous-bridge" but kubenet is not enabled, falling back to "hairpin-veth"
I0206 15:43:19.675877    4213 docker_service.go:240] Hairpin mode set to "hairpin-veth"
W0206 15:43:19.676098    4213 cni.go:239] Unable to update cni config: no networks found in /etc/cni/net.d
W0206 15:43:19.681418    4213 hostport_manager.go:71] The binary conntrack is not installed, this can cause failures in network connection cleanup.
W0206 15:43:19.682226    4213 hostport_manager.go:71] The binary conntrack is not installed, this can cause failures in network connection cleanup.
I0206 15:43:19.685272    4213 docker_service.go:255] Docker cri networking managed by kubernetes.io/no-op
I0206 15:43:19.700305    4213 docker_service.go:260] Docker Info: &{ID:OUWO:QIKR:KLYA:YHJZ:5VDA:7NBA:RME2:MVHU:2QSS:2PBR:F2DA:T6BM Containers:0 ContainersRunning:0 ContainersPaused:0 ContainersStopped:0 Images:0 Driver:overlay2 DriverStatus:[[Backing Filesystem extfs] [Supports d_type true] [Native Overlay Diff true]] SystemStatus:[] Plugins:{Volume:[local] Network:[bridge host ipvlan macvlan null overlay] Authorization:[] Log:[awslogs fluentd gcplogs gelf journald json-file local logentries splunk syslog]} MemoryLimit:true SwapLimit:false KernelMemory:true KernelMemoryTCP:true CPUCfsPeriod:true CPUCfsQuota:true CPUShares:true CPUSet:true PidsLimit:true IPv4Forwarding:true BridgeNfIptables:true BridgeNfIP6tables:true Debug:false NFd:23 OomKillDisable:true NGoroutines:41 SystemTime:2021-02-06T15:43:19.6870301+05:30 LoggingDriver:json-file CgroupDriver:cgroupfs NEventsListener:0 KernelVersion:5.4.0-54-generic OperatingSystem:Ubuntu Core 16 OSType:linux Architecture:x86_64 IndexServerAddress:https://index.docker.io/v1/ RegistryConfig:0xc0008449a0 NCPU:1 MemTotal:1029029888 GenericResources:[] DockerRootDir:/var/snap/docker/common/var-lib-docker HTTPProxy: HTTPSProxy: NoProxy: Name:my-own-k8s-cluster Labels:[] ExperimentalBuild:false ServerVersion:19.03.11 ClusterStore: ClusterAdvertise: Runtimes:map[runc:{Path:runc Args:[]}] DefaultRuntime:runc Swarm:{NodeID: NodeAddr: LocalNodeState:inactive ControlAvailable:false Error: RemoteManagers:[] Nodes:0 Managers:0 Cluster:<nil> Warnings:[]} LiveRestoreEnabled:false Isolation: InitBinary:docker-init ContainerdCommit:{ID:7ad184331fa3e55e52b890ea95e65ba581ae3429 Expected:7ad184331fa3e55e52b890ea95e65ba581ae3429} RuncCommit:{ID: Expected:} InitCommit:{ID:fec3683 Expected:fec3683} SecurityOptions:[name=apparmor name=seccomp,profile=default] ProductLicense: Warnings:[WARNING: No swap limit support]}
I0206 15:43:19.700428    4213 docker_service.go:273] Setting cgroupDriver to cgroupfs
I0206 15:43:19.740132    4213 remote_runtime.go:62] parsed scheme: ""
I0206 15:43:19.740827    4213 remote_runtime.go:62] scheme "" not registered, fallback to default scheme
I0206 15:43:19.741802    4213 passthrough.go:48] ccResolverWrapper: sending update to cc: {[{/var/run/dockershim.sock  <nil> 0 <nil>}] <nil> <nil>}
I0206 15:43:19.742509    4213 clientconn.go:948] ClientConn switching balancer to "pick_first"
I0206 15:43:19.743223    4213 remote_image.go:50] parsed scheme: ""
I0206 15:43:19.744902    4213 remote_image.go:50] scheme "" not registered, fallback to default scheme
I0206 15:43:19.745505    4213 passthrough.go:48] ccResolverWrapper: sending update to cc: {[{/var/run/dockershim.sock  <nil> 0 <nil>}] <nil> <nil>}
I0206 15:43:19.746014    4213 clientconn.go:948] ClientConn switching balancer to "pick_first"
I0206 15:43:19.793500    4213 kuberuntime_manager.go:216] Container runtime docker initialized, version: 19.03.11, apiVersion: 1.40.0
E0206 15:43:20.123700    4213 aws_credentials.go:77] while getting AWS credentials NoCredentialProviders: no valid providers in chain. Deprecated.
        For verbose messaging see aws.Config.CredentialsChainVerboseErrors
W0206 15:43:20.126645    4213 volume_host.go:75] kubeClient is nil. Skip initialization of CSIDriverLister
W0206 15:43:20.129181    4213 csi_plugin.go:190] kubernetes.io/csi: kubeclient not set, assuming standalone kubelet
W0206 15:43:20.130257    4213 csi_plugin.go:264] Skipping CSINode initialization, kubelet running in standalone mode
I0206 15:43:20.132007    4213 server.go:1176] Started kubelet
E0206 15:43:20.133638    4213 kubelet.go:1271] Image garbage collection failed once. Stats initialization may not have completed yet: failed to get imageFs info: unable to find data in memory cache
W0206 15:43:20.134227    4213 kubelet.go:1376] No api server defined - no node status update will be sent.
I0206 15:43:20.141608    4213 fs_resource_analyzer.go:64] Starting FS ResourceAnalyzer
I0206 15:43:20.134840    4213 server.go:148] Starting to listen on 0.0.0.0:10250
I0206 15:43:20.144712    4213 server.go:409] Adding debug handlers to kubelet server.
I0206 15:43:20.155999    4213 volume_manager.go:271] Starting Kubelet Volume Manager
I0206 15:43:20.156865    4213 desired_state_of_world_populator.go:142] Desired state populator starts to run
I0206 15:43:20.238956    4213 kubelet_network_linux.go:56] Initialized IPv4 iptables rules.
I0206 15:43:20.239492    4213 status_manager.go:154] Kubernetes client is nil, not starting status manager.
I0206 15:43:20.239559    4213 kubelet.go:1799] Starting kubelet main sync loop.
E0206 15:43:20.239645    4213 kubelet.go:1823] skipping pod synchronization - [container runtime status check may not have completed yet, PLEG is not healthy: pleg has yet to be successful]
E0206 15:43:20.345334    4213 kubelet.go:1823] skipping pod synchronization - container runtime status check may not have completed yet
I0206 15:43:20.407538    4213 reconciler.go:157] Reconciler: start to sync state
I0206 15:43:20.424400    4213 cpu_manager.go:193] [cpumanager] starting with none policy
I0206 15:43:20.425594    4213 cpu_manager.go:194] [cpumanager] reconciling every 10s
I0206 15:43:20.426700    4213 state_mem.go:36] [cpumanager] initializing new in-memory state store
I0206 15:43:20.427972    4213 state_mem.go:88] [cpumanager] updated default cpuset: ""
I0206 15:43:20.428614    4213 state_mem.go:96] [cpumanager] updated cpuset assignments: "map[]"
I0206 15:43:20.429429    4213 policy_none.go:43] [cpumanager] none policy: Start
W0206 15:43:20.431980    4213 manager.go:594] Failed to retrieve checkpoint for "kubelet_internal_checkpoint": checkpoint is not found
I0206 15:43:20.432982    4213 plugin_manager.go:114] Starting Kubelet Plugin Manager
E0206 15:43:20.436757    4213 container_manager_linux.go:487] cpu and memory cgroup hierarchy not unified.  cpu: /, memory: /system.slice/snap.docker.dockerd.service
E0206 15:43:20.447767    4213 container_manager_linux.go:533] failed to find cgroups of kubelet - cpu and memory cgroup hierarchy not unified.  cpu: /, memory: /user.slice/user-1000.slice/session-5.scope
W0206 15:43:20.494347    4213 eviction_manager.go:344] eviction manager: attempting to reclaim ephemeral-storage
I0206 15:43:20.495128    4213 container_gc.go:85] attempting to delete unused containers
I0206 15:43:20.502464    4213 image_gc_manager.go:321] attempting to delete unused images
I0206 15:43:20.602258    4213 eviction_manager.go:355] eviction manager: must evict pod(s) to reclaim ephemeral-storage
E0206 15:43:20.603227    4213 eviction_manager.go:366] eviction manager: eviction thresholds have been met, but no pods are active to evict
```

I see this message

```bash
No api server defined - no node status update will be sent.
```

```bash
$ kubelet --help | grep server
node. It can register the node with the apiserver using one of: the hostname; a flag to
various mechanisms (primarily through the apiserver) and ensures that the containers
Other than from an PodSpec from the apiserver, there are three ways that a container
HTTP server: The kubelet can also listen for HTTP and respond to a simple API
      --anonymous-auth                                           Enables anonymous requests to the Kubelet server. Requests that are not rejected by another authentication method are treated as anonymous requests. Anonymous requests have a username of system:anonymous, and a group name of system:unauthenticated. (default true) (DEPRECATED: This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.)
      --authorization-mode string                                Authorization mode for Kubelet server. Valid options are AlwaysAllow or Webhook. Webhook mode uses the SubjectAccessReview API to determine authorization. (default "AlwaysAllow") (DEPRECATED: This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.)
      --bootstrap-kubeconfig string                              Path to a kubeconfig file that will be used to get client certificate for kubelet. If the file specified by --kubeconfig does not exist, the bootstrap kubeconfig is used to request a client certificate from the API server. On success, a kubeconfig file referencing the generated client certificate and key is written to the path specified by --kubeconfig. The client certificate and key file will be stored in the directory pointed by --cert-dir.
      --cluster-dns strings                                      Comma-separated list of DNS server IP address.  This value is used for containers DNS server in case of Pods with "dnsPolicy=ClusterFirst". Note: all DNS servers appearing in the list MUST serve the same set of records otherwise name resolution within the cluster may not work correctly. There is no guarantee as to which DNS server may be contacted for name resolution. (DEPRECATED: This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.)
      --enable-debugging-handlers                                Enables server endpoints for log collection and local running of containers and commands (default true) (DEPRECATED: This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.)
      --enable-server                                            Enable the Kubelet's server (default true) (DEPRECATED: This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.)
      --healthz-bind-address ip                                  The IP address for the healthz server to serve on (set to '0.0.0.0' for all IPv4 interfaces and '::' for all IPv6 interfaces) (default 127.0.0.1) (DEPRECATED: This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.)
      --kube-api-burst int32                                     Burst to use while talking with kubernetes apiserver. Doesn't cover events and node heartbeat apis which rate limiting is controlled by a different set of flags (default 10) (DEPRECATED: This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.)
      --kube-api-content-type string                             Content type of requests sent to apiserver. (default "application/vnd.kubernetes.protobuf") (DEPRECATED: This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.)
      --kube-api-qps int32                                       QPS to use while talking with kubernetes apiserver. Doesn't cover events and node heartbeat apis which rate limiting is controlled by a different set of flags (default 5) (DEPRECATED: This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.)
      --kubeconfig string                                        Path to a kubeconfig file, specifying how to connect to the API server. Providing --kubeconfig enables API server mode, omitting --kubeconfig enables standalone mode.
      --register-node                                            Register the node with the apiserver. If --kubeconfig is not provided, this flag is irrelevant, as the Kubelet won't have an apiserver to register with. (default true)
      --rotate-certificates                                      <Warning: Beta feature> Auto rotate the kubelet client certificates by requesting new certificates from the kube-apiserver when the certificate expiration approaches. (DEPRECATED: This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.)
      --rotate-server-certificates                               Auto-request and rotate the kubelet serving certificates by requesting new certificates from the kube-apiserver when the certificate expiration approaches. Requires the RotateKubeletServerCertificate feature gate to be enabled, and approval of the submitted CertificateSigningRequest objects. (DEPRECATED: This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.)
      --runonce                                                  If true, exit after spawning pods from static pod files or remote urls. Exclusive with --enable-server (DEPRECATED: This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.)
      --tls-cert-file string                                     File containing x509 Certificate used for serving HTTPS (with intermediate certs, if any, concatenated after server cert). If --tls-cert-file and --tls-private-key-file are not provided, a self-signed certificate and key are generated for the public address and saved to the directory passed to --cert-dir. (DEPRECATED: This parameter should be set via the config file specified by the Kubelet's --config flag. See https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/ for more information.)
      --tls-cipher-suites strings                                Comma-separated list of cipher suites for the server. If omitted, the default Go cipher suites will be used.
```

https://kubernetes.io/docs/tasks/administer-cluster/kubelet-config-file/

I think we will need a kube config file to configure the kubelet. :) Especially
to not start it in standalone mode! I wonder if standalone mode is needed at all
but I'm not sure. Maybe it is 🤷

```bash
$ kubectl explain KubeletConfiguration
the server doesn't have a resource type "KubeletConfiguration"
```

```bash
$ sudo kubelet --root-dir /opt/kubelet/ --kubeconfig kubelet.kubeconfig
I0206 15:57:12.357364    4592 server.go:416] Version: v1.20.1
F0206 15:57:12.360060    4592 server.go:269] failed to run Kubelet: invalid kubeconfig: error loading config file "kubelet.kubeconfig": no kind "KubeletConfiguration" is registered for version "kubelet.config.k8s.io/v1beta1" in scheme "k8s.io/client-go/tools/clientcmd/api/latest/latest.go:50"
goroutine 1 [running]:
k8s.io/kubernetes/vendor/k8s.io/klog/v2.stacks(0xc000010001, 0xc00055a000, 0x126, 0x178)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:1026 +0xb9
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).output(0x70c9460, 0xc000000003, 0x0, 0x0, 0xc0007a77a0, 0x6f34162, 0x9, 0x10d, 0x411b00)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:975 +0x19b
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).printDepth(0x70c9460, 0xc000000003, 0x0, 0x0, 0x0, 0x0, 0x1, 0xc0007f91c0, 0x1, 0x1)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:732 +0x16f
k8s.io/kubernetes/vendor/k8s.io/klog/v2.(*loggingT).print(...)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:714
k8s.io/kubernetes/vendor/k8s.io/klog/v2.Fatal(...)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/klog/v2/klog.go:1482
k8s.io/kubernetes/cmd/kubelet/app.NewKubeletCommand.func1(0xc0001958c0, 0xc00004e0b0, 0x4, 0x4)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/cmd/kubelet/app/server.go:269 +0x845
k8s.io/kubernetes/vendor/github.com/spf13/cobra.(*Command).execute(0xc0001958c0, 0xc00004e0b0, 0x4, 0x4, 0xc0001958c0, 0xc00004e0b0)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/github.com/spf13/cobra/command.go:854 +0x2c2
k8s.io/kubernetes/vendor/github.com/spf13/cobra.(*Command).ExecuteC(0xc0001958c0, 0x166121bbdce5536a, 0x70c9020, 0x409b25)
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

goroutine 98 [chan receive]:
k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/server.SetupSignalContext.func1(0xc0007f8bb0)
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
k8s.io/kubernetes/vendor/go.opencensus.io/stats/view.(*worker).start(0xc00004f360)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/go.opencensus.io/stats/view/worker.go:154 +0x105
created by k8s.io/kubernetes/vendor/go.opencensus.io/stats/view.init.0
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/go.opencensus.io/stats/view/worker.go:32 +0x57

goroutine 94 [select]:
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.BackoffUntil(0x4a721f0, 0x4f0bde0, 0xc0003a06c0, 0x1, 0xc00005a120)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:167 +0x149
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.JitterUntil(0x4a721f0, 0x12a05f200, 0x0, 0xc000b46d01, 0xc00005a120)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:133 +0x98
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.Until(...)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:90
k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait.Forever(0x4a721f0, 0x12a05f200)
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:81 +0x4f
created by k8s.io/kubernetes/vendor/k8s.io/component-base/logs.InitLogs
        /workspace/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/component-base/logs/logs.go:58 +0x8a
```

https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/kubelet/config/v1beta1/types.go

```bash
$ cat kubelet.kubeconfig
apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration
evictionHard:
    memory.available:  "200Mi"
tlsCertFile: "kubelet-server.pem"
tlsPrivateKeyFile: "kubelet-server-key.pem"
```

Oops. I created a file for `--config`. Not for `--kubeconfig`. Hmm

For kube config, I need client certificates. So

Steps:
- create cfssl json config for creating client ssl certificates
  https://kubernetes.io/docs/concepts/cluster-administration/certificates/
  https://kubernetes.io/docs/concepts/cluster-administration/certificates/#cfssl
  https://kubernetes.io/docs/setup/best-practices/certificates/
  https://kubernetes.io/docs/setup/best-practices/certificates/#configure-certificates-for-user-accounts

        CN - common name - "system:node:<nodeName>"
        O - organization - "system:nodes"

        Note: The value of <nodeName> for kubelet.conf must match precisely the
        value of the node name provided by the kubelet as it registers with the
        apiserver. For further details, read the Node Authorization. -
        https://kubernetes.io/docs/reference/access-authn-authz/node/
        https://kubernetes.io/docs/reference/command-line-tools-reference/kubelet-tls-bootstrapping/

- run cfssl and create the client ssl certificates
- use kubectl config command to create kube config `kubelete.kubeconfig` and
  rename kubelet config file to `kubelet.config` as that's more appropriate.

https://kubernetes.io/docs/reference/access-authn-authz/node/

```
In order to be authorized by the Node authorizer, kubelets must use a credential
that identifies them as being in the system:nodes group, with a username of
system:node:<nodeName>. This group and user name format match the identity
created for each kubelet as part of kubelet TLS bootstrapping.

The value of <nodeName> must match precisely the name of the node as registered
by the kubelet. By default, this is the host name as provided by hostname, or
overridden via the kubelet option --hostname-override. However, when using the
--cloud-provider kubelet option, the specific hostname may be determined by the
cloud provider, ignoring the local hostname and the --hostname-override option.
For specifics about how the kubelet determines the hostname, see the kubelet
options reference.
````

I'll use "worker-1" as the node name :)
