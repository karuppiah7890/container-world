# Kind

https://github.com/kubernetes-sigs/kind/

https://kind.sigs.k8s.io/

https://kind.sigs.k8s.io/docs/user/quick-start/

Kind is short for Kubernetes in Docker.

I mean, that's like wow and crazy!!! Why? Well, Kubernetes has a lot of running
components as it's a distributed system. Imagine running Kubernetes in Docker.
I mean, I'm guessing it's a single Docker Container. But let's see! :)

I'm reading the https://kind.sigs.k8s.io/ page.

Apparently kind was designed for testing Kubernetes itself. But it can also
be used to run Kubernetes locally or in a Continuous Integration environment.
Nice!

I tried

```bash
$ go get sigs.k8s.io/kind@v0.10.0
```

https://github.com/kubernetes-sigs/kind/releases

And it worked but I went with brew. It's easier to upgrade. I mean, I can
upgrade with `go` too. Still. Anyways :)

```bash
$ brew install kind
```

But damn `brew` is slow usually. I could have used `gofish`. Hmm

```bash
$ gofish search kind
NAME    RIG                             VERSION
kind    github.com/fishworks/fish-food  0.10.0
```

Anyways, done. With `brew` I don't have to wait for version updates / maintain
the package in the package manager like in `gofish`. `brew` has lot of
contributors :)

```bash
$ kind --help
kind creates and manages local Kubernetes clusters using Docker container 'nodes'

Usage:
  kind [command]

Available Commands:
  build       Build one of [node-image]
  completion  Output shell completion code for the specified shell (bash, zsh or fish)
  create      Creates one of [cluster]
  delete      Deletes one of [cluster]
  export      Exports one of [kubeconfig, logs]
  get         Gets one of [clusters, nodes, kubeconfig]
  help        Help about any command
  load        Loads images into nodes
  version     Prints the kind CLI version

Flags:
  -h, --help              help for kind
      --loglevel string   DEPRECATED: see -v instead
  -q, --quiet             silence all stderr output
  -v, --verbosity int32   info log verbosity
      --version           version for kind

Use "kind [command] --help" for more information about a command.
```

```
$ kind --version
kind version 0.10.0

$ kind version
kind v0.10.0 go1.15.7 darwin/amd64
```

Creating a cluster is too simple it sesms

```bash
$ kind create cluster
```

Apparently kind uses a node image - https://kind.sigs.k8s.io/docs/design/node-image , which in turn is based on a base image - https://kind.sigs.k8s.io/docs/design/base-image

Before jumping more into things. I just wanted to stop and understand that I'm
now running Kubernetes locally, using Docker containers.

I think I now know two ways to easily run Kubernetes clusters in my local. One
is using minikube, now another is using kind. minikube uses VMs and kind uses
Docker containers.

I can only see one Docker container running.

```bash
$ docker ps
CONTAINER ID   IMAGE                  COMMAND                  CREATED         STATUS         PORTS                       NAMES
1556f88bead9   kindest/node:v1.20.2   "/usr/local/bin/entr…"   3 minutes ago   Up 3 minutes   127.0.0.1:52237->6443/tcp   kind-control-plane
```

So, like I thought, full kubernetes cluster inside ONE docker container. Wow.
Hmm.

Final output from kind after running the cluster

```bash
$ kind create cluster
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.20.2) 🖼
 ✓ Preparing nodes 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind

Thanks for using kind! 😊
```

It also added a user and cluster and a context in my kube config I think. And
has also changed my current kubernetes context in my kube config (file).

```bash
$ kubectl cluster-info
Kubernetes master is running at https://127.0.0.1:52237
KubeDNS is running at https://127.0.0.1:52237/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.

$ # same output when I use specific context flag
$ kubectl cluster-info --context kind-kind
Kubernetes master is running at https://127.0.0.1:52237
KubeDNS is running at https://127.0.0.1:52237/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
```

```bash
$ k version
Client Version: version.Info{Major:"1", Minor:"19", GitVersion:"v1.19.5", GitCommit:"e338cf2c6d297aa603b50ad3a301f761b4173aa6", GitTreeState:"clean", BuildDate:"2020-12-09T11:18:51Z", GoVersion:"go1.15.2", Compiler:"gc", Platform:"darwin/amd64"}
Server Version: version.Info{Major:"1", Minor:"20", GitVersion:"v1.20.2", GitCommit:"faecb196815e248d3ecfb03c680a4507229c2a56", GitTreeState:"clean", BuildDate:"2021-01-21T01:11:42Z", GoVersion:"go1.15.5", Compiler:"gc", Platform:"linux/amd64"}

$ k get nodes
NAME                 STATUS   ROLES                  AGE     VERSION
kind-control-plane   Ready    control-plane,master   5m24s   v1.20.2
```

Now I checked the different resources like pods and the config maps and all

```bash
$ k get pods,cm,svc,roles,rolebinding -A
NAMESPACE            NAME                                             READY   STATUS    RESTARTS   AGE
kube-system          pod/coredns-74ff55c5b-7cntz                      1/1     Running   0          6m27s
kube-system          pod/coredns-74ff55c5b-kvpxc                      1/1     Running   0          6m27s
kube-system          pod/etcd-kind-control-plane                      1/1     Running   0          6m35s
kube-system          pod/kindnet-g52m8                                1/1     Running   0          6m27s
kube-system          pod/kube-apiserver-kind-control-plane            1/1     Running   0          6m35s
kube-system          pod/kube-controller-manager-kind-control-plane   1/1     Running   0          6m35s
kube-system          pod/kube-proxy-gtj97                             1/1     Running   0          6m27s
kube-system          pod/kube-scheduler-kind-control-plane            1/1     Running   0          6m35s
local-path-storage   pod/local-path-provisioner-78776bfc44-nf8w5      1/1     Running   0          6m27s

NAMESPACE            NAME                                           DATA   AGE
default              configmap/kube-root-ca.crt                     1      6m28s
kube-node-lease      configmap/kube-root-ca.crt                     1      6m28s
kube-public          configmap/cluster-info                         2      6m42s
kube-public          configmap/kube-root-ca.crt                     1      6m28s
kube-system          configmap/coredns                              1      6m41s
kube-system          configmap/extension-apiserver-authentication   6      6m47s
kube-system          configmap/kube-proxy                           2      6m41s
kube-system          configmap/kube-root-ca.crt                     1      6m28s
kube-system          configmap/kubeadm-config                       2      6m44s
kube-system          configmap/kubelet-config-1.20                  1      6m44s
local-path-storage   configmap/kube-root-ca.crt                     1      6m28s
local-path-storage   configmap/local-path-config                    1      6m37s

NAMESPACE     NAME                 TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)                  AGE
default       service/kubernetes   ClusterIP   10.96.0.1    <none>        443/TCP                  6m45s
kube-system   service/kube-dns     ClusterIP   10.96.0.10   <none>        53/UDP,53/TCP,9153/TCP   6m41s

NAMESPACE     NAME                                                                            CREATED AT
kube-public   role.rbac.authorization.k8s.io/kubeadm:bootstrap-signer-clusterinfo             2021-03-06T17:37:38Z
kube-public   role.rbac.authorization.k8s.io/system:controller:bootstrap-signer               2021-03-06T17:37:35Z
kube-system   role.rbac.authorization.k8s.io/extension-apiserver-authentication-reader        2021-03-06T17:37:35Z
kube-system   role.rbac.authorization.k8s.io/kube-proxy                                       2021-03-06T17:37:39Z
kube-system   role.rbac.authorization.k8s.io/kubeadm:kubelet-config-1.20                      2021-03-06T17:37:36Z
kube-system   role.rbac.authorization.k8s.io/kubeadm:nodes-kubeadm-config                     2021-03-06T17:37:36Z
kube-system   role.rbac.authorization.k8s.io/system::leader-locking-kube-controller-manager   2021-03-06T17:37:35Z
kube-system   role.rbac.authorization.k8s.io/system::leader-locking-kube-scheduler            2021-03-06T17:37:35Z
kube-system   role.rbac.authorization.k8s.io/system:controller:bootstrap-signer               2021-03-06T17:37:35Z
kube-system   role.rbac.authorization.k8s.io/system:controller:cloud-provider                 2021-03-06T17:37:35Z
kube-system   role.rbac.authorization.k8s.io/system:controller:token-cleaner                  2021-03-06T17:37:35Z

NAMESPACE     NAME                                                                                      ROLE                                                  AGE
kube-public   rolebinding.rbac.authorization.k8s.io/kubeadm:bootstrap-signer-clusterinfo                Role/kubeadm:bootstrap-signer-clusterinfo             6m42s
kube-public   rolebinding.rbac.authorization.k8s.io/system:controller:bootstrap-signer                  Role/system:controller:bootstrap-signer               6m45s
kube-system   rolebinding.rbac.authorization.k8s.io/kube-proxy                                          Role/kube-proxy                                       6m41s
kube-system   rolebinding.rbac.authorization.k8s.io/kubeadm:kubelet-config-1.20                         Role/kubeadm:kubelet-config-1.20                      6m44s
kube-system   rolebinding.rbac.authorization.k8s.io/kubeadm:nodes-kubeadm-config                        Role/kubeadm:nodes-kubeadm-config                     6m44s
kube-system   rolebinding.rbac.authorization.k8s.io/system::extension-apiserver-authentication-reader   Role/extension-apiserver-authentication-reader        6m45s
kube-system   rolebinding.rbac.authorization.k8s.io/system::leader-locking-kube-controller-manager      Role/system::leader-locking-kube-controller-manager   6m45s
kube-system   rolebinding.rbac.authorization.k8s.io/system::leader-locking-kube-scheduler               Role/system::leader-locking-kube-scheduler            6m45s
kube-system   rolebinding.rbac.authorization.k8s.io/system:controller:bootstrap-signer                  Role/system:controller:bootstrap-signer               6m45s
kube-system   rolebinding.rbac.authorization.k8s.io/system:controller:cloud-provider                    Role/system:controller:cloud-provider                 6m45s
kube-system   rolebinding.rbac.authorization.k8s.io/system:controller:token-cleaner                     Role/system:controller:token-cleaner                  6m45s
```

I was checking how much resources this Docker container takes up to run the K8s
cluster and also all these pods!!

```bash
$ docker stats
CONTAINER ID   NAME                 CPU %     MEM USAGE / LIMIT     MEM %     NET I/O          BLOCK I/O        PIDS
1556f88bead9   kind-control-plane   22.20%    643.9MiB / 1.452GiB   43.31%    40.6kB / 359kB   102MB / 13.6MB   253
```

That was one of the snapshots as it was continuously changing.

I also tried checking the commands running inside the Docker container and
checking how much resources they take

```bash
$ docker top 1
UID                 PID                 PPID                C                   STIME               TTY                 TIME                CMD
root                1703                1677                0                   17:36               ?                   00:00:00            /sbin/init
root                3114                3073                0                   17:38               ?                   00:00:00            /pause
root                3257                3073                0                   17:38               ?                   00:00:01            local-path-provisioner --debug start --helper-image k8s.gcr.io/build-image/debian-base:v2.1.0 --config /etc/config/config.json
root                2715                2682                0                   17:37               ?                   00:00:00            /pause
root                2788                2682                0                   17:37               ?                   00:00:00            /usr/local/bin/kube-proxy --config=/var/lib/kube-proxy/config.conf --hostname-override=kind-control-plane
root                2413                2185                13                  17:37               ?                   00:01:15            kube-apiserver --advertise-address=172.19.0.2 --allow-privileged=true --authorization-mode=Node,RBAC --client-ca-file=/etc/kubernetes/pki/ca.crt --enable-admission-plugins=NodeRestriction --enable-bootstrap-token-auth=true --etcd-cafile=/etc/kubernetes/pki/etcd/ca.crt --etcd-certfile=/etc/kubernetes/pki/apiserver-etcd-client.crt --etcd-keyfile=/etc/kubernetes/pki/apiserver-etcd-client.key --etcd-servers=https://127.0.0.1:2379 --insecure-port=0 --kubelet-client-certificate=/etc/kubernetes/pki/apiserver-kubelet-client.crt --kubelet-client-key=/etc/kubernetes/pki/apiserver-kubelet-client.key --kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname --proxy-client-cert-file=/etc/kubernetes/pki/front-proxy-client.crt --proxy-client-key-file=/etc/kubernetes/pki/front-proxy-client.key --requestheader-allowed-names=front-proxy-client --requestheader-client-ca-file=/etc/kubernetes/pki/front-proxy-ca.crt --requestheader-extra-headers-prefix=X-Remote-Extra- --requestheader-group-headers=X-Remote-Group --requestheader-username-headers=X-Remote-User --runtime-config= --secure-port=6443 --service-account-issuer=https://kubernetes.default.svc.cluster.local --service-account-key-file=/etc/kubernetes/pki/sa.pub --service-account-signing-key-file=/etc/kubernetes/pki/sa.key --service-cluster-ip-range=10.96.0.0/16 --tls-cert-file=/etc/kubernetes/pki/apiserver.crt --tls-private-key-file=/etc/kubernetes/pki/apiserver.key
root                2282                2185                0                   17:37               ?                   00:00:00            /pause
root                2268                2182                0                   17:37               ?                   00:00:00            /pause
root                2363                2182                0                   17:37               ?                   00:00:05            kube-scheduler --authentication-kubeconfig=/etc/kubernetes/scheduler.conf --authorization-kubeconfig=/etc/kubernetes/scheduler.conf --bind-address=127.0.0.1 --kubeconfig=/etc/kubernetes/scheduler.conf --leader-elect=true --port=0
root                2261                2183                0                   17:37               ?                   00:00:00            /pause
root                2366                2183                2                   17:37               ?                   00:00:17            kube-controller-manager --allocate-node-cidrs=true --authentication-kubeconfig=/etc/kubernetes/controller-manager.conf --authorization-kubeconfig=/etc/kubernetes/controller-manager.conf --bind-address=127.0.0.1 --client-ca-file=/etc/kubernetes/pki/ca.crt --cluster-cidr=10.244.0.0/16 --cluster-name=kind --cluster-signing-cert-file=/etc/kubernetes/pki/ca.crt --cluster-signing-key-file=/etc/kubernetes/pki/ca.key --controllers=*,bootstrapsigner,tokencleaner --enable-hostpath-provisioner=true --kubeconfig=/etc/kubernetes/controller-manager.conf --leader-elect=true --port=0 --requestheader-client-ca-file=/etc/kubernetes/pki/front-proxy-ca.crt --root-ca-file=/etc/kubernetes/pki/ca.crt --service-account-private-key-file=/etc/kubernetes/pki/sa.key --service-cluster-ip-range=10.96.0.0/16 --use-service-account-credentials=true
root                3116                3051                0                   17:38               ?                   00:00:00            /pause
root                3194                3051                0                   17:38               ?                   00:00:02            /coredns -conf /etc/coredns/Corefile
root                2275                2181                0                   17:37               ?                   00:00:00            /pause
root                2451                2181                3                   17:37               ?                   00:00:19            etcd --advertise-client-urls=https://172.19.0.2:2379 --cert-file=/etc/kubernetes/pki/etcd/server.crt --client-cert-auth=true --data-dir=/var/lib/etcd --initial-advertise-peer-urls=https://172.19.0.2:2380 --initial-cluster=kind-control-plane=https://172.19.0.2:2380 --key-file=/etc/kubernetes/pki/etcd/server.key --listen-client-urls=https://127.0.0.1:2379,https://172.19.0.2:2379 --listen-metrics-urls=http://127.0.0.1:2381 --listen-peer-urls=https://172.19.0.2:2380 --name=kind-control-plane --peer-cert-file=/etc/kubernetes/pki/etcd/peer.crt --peer-client-cert-auth=true --peer-key-file=/etc/kubernetes/pki/etcd/peer.key --peer-trusted-ca-file=/etc/kubernetes/pki/etcd/ca.crt --snapshot-count=10000 --trusted-ca-file=/etc/kubernetes/pki/etcd/ca.crt
root                3129                3050                0                   17:38               ?                   00:00:00            /pause
root                3205                3050                0                   17:38               ?                   00:00:02            /coredns -conf /etc/coredns/Corefile
root                2779                2659                0                   17:37               ?                   00:00:00            /bin/kindnetd
root                2705                2659                0                   17:37               ?                   00:00:00            /pause
root                1946                1703                2                   17:36               ?                   00:00:13            /usr/local/bin/containerd
root                2181                1703                0                   17:37               ?                   00:00:00            /usr/local/bin/containerd-shim-runc-v2 -namespace k8s.io -id bcdd70b783068932db25e096fb099cb670d665343b985d13a36c0adf2932c9ae -address /run/containerd/containerd.sock
root                2182                1703                0                   17:37               ?                   00:00:00            /usr/local/bin/containerd-shim-runc-v2 -namespace k8s.io -id 43842c48bbbf0adf786f5646654996a715edcc6f8c07064d8f46fa8a9411384c -address /run/containerd/containerd.sock
root                2183                1703                0                   17:37               ?                   00:00:00            /usr/local/bin/containerd-shim-runc-v2 -namespace k8s.io -id 54af00fc50d7de27509da2eeb28668b7e8fdc16ecbd630a16521dc2f96c81cc0 -address /run/containerd/containerd.sock
root                2185                1703                0                   17:37               ?                   00:00:00            /usr/local/bin/containerd-shim-runc-v2 -namespace k8s.io -id d9eaba722b7c1fde2ff90ac15aa20cd6a456196fa07d5ce73acdb79de63bb5d2 -address /run/containerd/containerd.sock
root                2659                1703                0                   17:37               ?                   00:00:00            /usr/local/bin/containerd-shim-runc-v2 -namespace k8s.io -id f42b3efa15f44639f950e8a8291e6b5d6f683f2d5ec927a815932358e6e21422 -address /run/containerd/containerd.sock
root                2682                1703                0                   17:37               ?                   00:00:00            /usr/local/bin/containerd-shim-runc-v2 -namespace k8s.io -id a721c59a2a7f135007394cae62af86080b09723e8fa2937b69eb607670b7efc7 -address /run/containerd/containerd.sock
root                3050                1703                0                   17:38               ?                   00:00:00            /usr/local/bin/containerd-shim-runc-v2 -namespace k8s.io -id 4a973cd1a660b354852e0f6bc9ca38855d4e058c2bee888d9309c9e83f5d9132 -address /run/containerd/containerd.sock
root                3051                1703                0                   17:38               ?                   00:00:00            /usr/local/bin/containerd-shim-runc-v2 -namespace k8s.io -id 270cec3cae644c75a8af67feec5ac26dae027b4dcae490698079749d8b8dc65f -address /run/containerd/containerd.sock
root                3073                1703                0                   17:38               ?                   00:00:00            /usr/local/bin/containerd-shim-runc-v2 -namespace k8s.io -id 3f94f99be795cd9b50e81d6265f1362bc6bdd09d440255fd009d15f7e54583c2 -address /run/containerd/containerd.sock
root                2497                1703                5                   17:37               ?                   00:00:32            /usr/bin/kubelet --bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.conf --kubeconfig=/etc/kubernetes/kubelet.conf --config=/var/lib/kubelet/config.yaml --container-runtime=remote --container-runtime-endpoint=unix:///run/containerd/containerd.sock --fail-swap-on=false --node-ip=172.19.0.2 --provider-id=kind://docker/kind/kind-control-plane --fail-swap-on=false --cgroup-root=/kubelet
root                1935                1703                0                   17:36               ?                   00:00:00            /lib/systemd/systemd-journald
```

I think C means CPU. Gotta check. Anyways, those are the values. Hmm.

I was also able to smoothly get inside a pod - exec into it and then check it
out. Cool stuff! Very smooth and simple kubernetes cluster! :D :)

```bash
$ execpod -a

 kubectl exec --namespace='kube-system' kindnet-g52m8 -c kindnet-cni -it -- sh

#
# bash
sh: 2: bash: not found
# pwd
/
# ls
bin  boot  dev  etc  home  lib  lib64  media  mnt  opt  proc  root  run  sbin  srv  sys  tmp  usr  var
#
```

I'm sleepy now. Gonna go sleep and resume later. I was checking how to delete
kind cluster in this case. I was checking how I can refer to the cluster I'm
running, with some name.

I noticed a `get` subcommand

```bash
$ kind get
Gets one of [clusters, nodes, kubeconfig]

Usage:
  kind get [command]

Available Commands:
  clusters    Lists existing kind clusters by their name
  kubeconfig  Prints cluster kubeconfig
  nodes       Lists existing kind nodes by their name

Flags:
  -h, --help   help for get

Global Flags:
      --loglevel string   DEPRECATED: see -v instead
  -q, --quiet             silence all stderr output
  -v, --verbosity int32   info log verbosity

Use "kind get [command] --help" for more information about a command.
```

```bash
$ kind get nodes
kind-control-plane


$ kind get clusters
kind
```

The `delete` sub command can only delete a cluster. Of course, makes sense.
Deleting that would delete node and even kubeconfig I think. Let's try!

```bash
$ k config view
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: DATA+OMITTED
    server: https://127.0.0.1:52237
  name: kind-kind
- cluster:
    certificate-authority-data: DATA+OMITTED
    server: https://192.168.199.40:6443
  name: kubernetes-the-hard-way
contexts:
- context:
    cluster: kind-kind
    user: kind-kind
  name: kind-kind
- context:
    cluster: kubernetes-the-hard-way
    namespace: default
    user: admin
  name: kubernetes-the-hard-way
current-context: kind-kind
kind: Config
preferences: {}
users:
- name: admin
  user:
    client-certificate: /Users/karuppiahn/oss/github.com/kinvolk/kubernetes-the-hard-way-vagrant/certificates/admin.pem
    client-key: /Users/karuppiahn/oss/github.com/kinvolk/kubernetes-the-hard-way-vagrant/certificates/admin-key.pem
- name: kind-kind
  user:
    client-certificate-data: REDACTED
    client-key-data: REDACTED
```

```bash
$ k config view --minify
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: DATA+OMITTED
    server: https://127.0.0.1:52237
  name: kind-kind
contexts:
- context:
    cluster: kind-kind
    user: kind-kind
  name: kind-kind
current-context: kind-kind
kind: Config
preferences: {}
users:
- name: kind-kind
  user:
    client-certificate-data: REDACTED
    client-key-data: REDACTED
```

Notice the `kind-kind` cluster and the user `kind-kind` and also the context
called `kind-kind` for `kind-kind` cluster with `kind-kind` user. :)

Let's delete the whole cluster now

```bash
$ kind delete
Deletes one of [cluster]

Usage:
  kind delete [command]

Available Commands:
  cluster     Deletes a cluster
  clusters    Deletes one or more clusters

Flags:
  -h, --help   help for delete

Global Flags:
      --loglevel string   DEPRECATED: see -v instead
  -q, --quiet             silence all stderr output
  -v, --verbosity int32   info log verbosity

Use "kind delete [command] --help" for more information about a command.
```

```bash
$ kind delete clusters -h
Deletes a resource

Usage:
  kind delete clusters [flags]

Flags:
      --all                 delete all clusters
  -h, --help                help for clusters
      --kubeconfig string   sets kubeconfig path instead of $KUBECONFIG or $HOME/.kube/config

Global Flags:
      --loglevel string   DEPRECATED: see -v instead
  -q, --quiet             silence all stderr output
  -v, --verbosity int32   info log verbosity

$ kind delete clusters
ERROR: no cluster names provided
```

Let's delete all? I mean there's only one. Hmm. Let's delete just by using the
name maybe.

Oh, I didn't have to give the name. Since there was only one, it got deleted I
think

```bash
$ kind delete cluster
Deleting cluster "kind" ...
```

```bash
$ kind get clusters
No kind clusters found.

$ kind get nodes
No kind nodes found for cluster "kind".

$ kind get kubeconfig
ERROR: could not locate any control plane nodes
```

Now no clusters are there, no nodes are there, no kube config too!

```bash
$ k config view
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: DATA+OMITTED
    server: https://192.168.199.40:6443
  name: kubernetes-the-hard-way
contexts:
- context:
    cluster: kubernetes-the-hard-way
    namespace: default
    user: admin
  name: kubernetes-the-hard-way
current-context: ""
kind: Config
preferences: {}
users:
- name: admin
  user:
    client-certificate: /Users/karuppiahn/oss/github.com/kinvolk/kubernetes-the-hard-way-vagrant/certificates/admin.pem
    client-key: /Users/karuppiahn/oss/github.com/kinvolk/kubernetes-the-hard-way-vagrant/certificates/admin-key.pem


$ k config view --minify
error: current-context must exist in order to minify
```

Notice how the current context info also has been removed from the kubeconfig?
Cool right? :) Nice tool!! :D
