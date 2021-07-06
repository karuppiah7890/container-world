https://duckduckgo.com/?t=ffab&q=lightweight+k8s&ia=web

https://microk8s.io/

```bash
$ brew install ubuntu/microk8s/microk8s
Updating Homebrew...
==> Auto-updated Homebrew!
Updated 1 tap (homebrew/cask).
==> Updated Casks
Updated 1 cask.

==> Tapping ubuntu/microk8s
Cloning into '/usr/local/Homebrew/Library/Taps/ubuntu/homebrew-microk8s'...
remote: Enumerating objects: 110, done.
remote: Counting objects: 100% (110/110), done.
remote: Compressing objects: 100% (74/74), done.
remote: Total 110 (delta 26), reused 39 (delta 9), pack-reused 0
Receiving objects: 100% (110/110), 29.42 KiB | 7.35 MiB/s, done.
Resolving deltas: 100% (26/26), done.
Tapped 1 formula (15 files, 52.4KB).
==> Installing microk8s from ubuntu/microk8s
==> Downloading https://files.pythonhosted.org/packages/06/a9/cd1fd8ee13f73a4d4f491ee219deeeae20afefa914dfb4c130cfc9dc397a/cer
######################################################################## 100.0%
==> Downloading https://files.pythonhosted.org/packages/fc/bb/a5768c230f9ddb03acc9ef3f0d4a3cf93462473795d18e9535498c8f929d/cha
######################################################################## 100.0%
==> Downloading https://files.pythonhosted.org/packages/27/6f/be940c8b1f1d69daceeb0032fee6c34d7bd70e3e649ccac0951500b4720e/cli
######################################################################## 100.0%
==> Downloading https://files.pythonhosted.org/packages/65/c4/80f97e9c9628f3cac9b98bfca0402ede54e0563b56482e3e6e45c43c4935/idn
######################################################################## 100.0%
==> Downloading https://files.pythonhosted.org/packages/71/fc/7c8e01f41a6e671d7b11be470eeb3d15339c75ce5559935f3f55890eec6b/pro
######################################################################## 100.0%
==> Downloading https://files.pythonhosted.org/packages/97/10/92d25b93e9c266c94b76a5548f020f3f1dd0eb40649cb1993532c0af8f4c/req
######################################################################## 100.0%
==> Downloading https://files.pythonhosted.org/packages/f3/94/67d781fb32afbee0fffa0ad9e16ad0491f1a9c303e14790ae4e18f11be19/req
######################################################################## 100.0%
==> Downloading https://files.pythonhosted.org/packages/f0/07/26b519e6ebb03c2a74989f7571e6ae6b82e9d7d81b8de6fcdbfc643c7b58/sim
######################################################################## 100.0%
==> Downloading https://files.pythonhosted.org/packages/b9/19/5cbd78eac8b1783671c40e34bb0fa83133a06d340a38b55c645076d40094/tom
######################################################################## 100.0%
==> Downloading https://files.pythonhosted.org/packages/fd/fa/b21f4f03176463a6cccdb612a5ff71b927e5224e83483012747c12fc5d62/url
######################################################################## 100.0%
==> Downloading https://github.com/ubuntu/microk8s/archive/installer-v2.1.0.tar.gz
==> Downloading from https://codeload.github.com/ubuntu/microk8s/tar.gz/installer-v2.1.0
                     -=O=-      #     #     #     #                           
==> python3 -m venv --system-site-packages /usr/local/Cellar/microk8s/2.1.0/libexec
==> /usr/local/Cellar/microk8s/2.1.0/libexec/bin/pip install -v --no-deps --no-binary :all: --ignore-installed /private/tmp/mi
==> /usr/local/Cellar/microk8s/2.1.0/libexec/bin/pip install -v --no-deps --no-binary :all: --ignore-installed /private/tmp/mi
==> /usr/local/Cellar/microk8s/2.1.0/libexec/bin/pip install -v --no-deps --no-binary :all: --ignore-installed /private/tmp/mi
==> /usr/local/Cellar/microk8s/2.1.0/libexec/bin/pip install -v --no-deps --no-binary :all: --ignore-installed /private/tmp/mi
==> /usr/local/Cellar/microk8s/2.1.0/libexec/bin/pip install -v --no-deps --no-binary :all: --ignore-installed /private/tmp/mi
==> /usr/local/Cellar/microk8s/2.1.0/libexec/bin/pip install -v --no-deps --no-binary :all: --ignore-installed /private/tmp/mi
==> /usr/local/Cellar/microk8s/2.1.0/libexec/bin/pip install -v --no-deps --no-binary :all: --ignore-installed /private/tmp/mi
==> /usr/local/Cellar/microk8s/2.1.0/libexec/bin/pip install -v --no-deps --no-binary :all: --ignore-installed /private/tmp/mi
==> /usr/local/Cellar/microk8s/2.1.0/libexec/bin/pip install -v --no-deps --no-binary :all: --ignore-installed /private/tmp/mi
==> /usr/local/Cellar/microk8s/2.1.0/libexec/bin/pip install -v --no-deps --no-binary :all: --ignore-installed /private/tmp/mi
==> /usr/local/Cellar/microk8s/2.1.0/libexec/bin/pip install -v --no-deps --no-binary :all: --ignore-installed /private/tmp/mi
==> Caveats
Run `microk8s install` to start with MicroK8s
==> Summary
🍺  /usr/local/Cellar/microk8s/2.1.0: 1,564 files, 15.9MB, built in 31 seconds
$ 
```

```bash
$ microk8s install
VM disk size requested exceeds free space on host.
Launched: microk8s-vm                                                           
2021-07-06T22:33:46+05:30 INFO Waiting for automatic snapd restart...
microk8s v1.21.1 from Canonical✓ installed
microk8s-integrator-macos 0.1 from Canonical✓ installed
MicroK8s is up and running. See the available commands with `microk8s --help`.
```

```bash
$ microk8s status --wait-ready
microk8s is running
high-availability: no
  datastore master nodes: 127.0.0.1:19001
  datastore standby nodes: none
addons:
  enabled:
    ha-cluster           # Configure high availability on the current node
  disabled:
    ambassador           # Ambassador API Gateway and Ingress
    cilium               # SDN, fast with full network policy
    dashboard            # The Kubernetes dashboard
    dns                  # CoreDNS
    fluentd              # Elasticsearch-Fluentd-Kibana logging and monitoring
    gpu                  # Automatic enablement of Nvidia CUDA
    helm                 # Helm 2 - the package manager for Kubernetes
    helm3                # Helm 3 - Kubernetes package manager
    host-access          # Allow Pods connecting to Host services smoothly
    ingress              # Ingress controller for external access
    istio                # Core Istio service mesh services
    jaeger               # Kubernetes Jaeger operator with its simple config
    keda                 # Kubernetes-based Event Driven Autoscaling
    knative              # The Knative framework on Kubernetes.
    kubeflow             # Kubeflow for easy ML deployments
    linkerd              # Linkerd is a service mesh for Kubernetes and other frameworks
    metallb              # Loadbalancer for your Kubernetes cluster
    metrics-server       # K8s Metrics Server for API access to service metrics
    multus               # Multus CNI enables attaching multiple network interfaces to pods
    openebs              # OpenEBS is the open-source storage solution for Kubernetes
    openfaas             # openfaas serverless framework
    portainer            # Portainer UI for your Kubernetes cluster
    prometheus           # Prometheus operator for monitoring and logging
    rbac                 # Role-Based Access Control for authorisation
    registry             # Private image registry exposed on localhost:32000
    storage              # Storage class; allocates storage from host directory
    traefik              # traefik Ingress controller for external access
$ microk8s enable --help
Usage: microk8s enable [OPTIONS] ADDONS...

  Enables a MicroK8s addon.

  For a list of available addons, run `microk8s status`.

  To see help for individual addons, run:

      microk8s enable ADDON -- --help

Options:
  --help  Show this message and exit.

$ microk8s kubectl get all --all-namespaces
NAMESPACE     NAME                                          READY   STATUS    RESTARTS   AGE
kube-system   pod/calico-node-75pw9                         1/1     Running   0          114s
kube-system   pod/calico-kube-controllers-f7868dd95-8mqrg   1/1     Running   0          114s

NAMESPACE   NAME                 TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)   AGE
default     service/kubernetes   ClusterIP   10.152.183.1   <none>        443/TCP   2m4s

NAMESPACE     NAME                         DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR            AGE
kube-system   daemonset.apps/calico-node   1         1         1       1            1           kubernetes.io/os=linux   2m2s

NAMESPACE     NAME                                      READY   UP-TO-DATE   AVAILABLE   AGE
kube-system   deployment.apps/calico-kube-controllers   1/1     1            1           2m

NAMESPACE     NAME                                                DESIRED   CURRENT   READY   AGE
kube-system   replicaset.apps/calico-kube-controllers-f7868dd95   1         1         1       114s

$ microk8s kubectl 
kubectl controls the Kubernetes cluster manager.

 Find more information at: https://kubernetes.io/docs/reference/kubectl/overview/

Basic Commands (Beginner):
  create        Create a resource from a file or from stdin.
  expose        Take a replication controller, service, deployment or pod and expose it as a new Kubernetes Service
  run           Run a particular image on the cluster
  set           Set specific features on objects

Basic Commands (Intermediate):
  explain       Documentation of resources
  get           Display one or many resources
  edit          Edit a resource on the server
  delete        Delete resources by filenames, stdin, resources and names, or by resources and label selector

Deploy Commands:
  rollout       Manage the rollout of a resource
  scale         Set a new size for a Deployment, ReplicaSet or Replication Controller
  autoscale     Auto-scale a Deployment, ReplicaSet, StatefulSet, or ReplicationController

Cluster Management Commands:
  certificate   Modify certificate resources.
  cluster-info  Display cluster info
  top           Display Resource (CPU/Memory) usage.
  cordon        Mark node as unschedulable
  uncordon      Mark node as schedulable
  drain         Drain node in preparation for maintenance
  taint         Update the taints on one or more nodes

Troubleshooting and Debugging Commands:
  describe      Show details of a specific resource or group of resources
  logs          Print the logs for a container in a pod
  attach        Attach to a running container
  exec          Execute a command in a container
  port-forward  Forward one or more local ports to a pod
  proxy         Run a proxy to the Kubernetes API server
  cp            Copy files and directories to and from containers.
  auth          Inspect authorization
  debug         Create debugging sessions for troubleshooting workloads and nodes

Advanced Commands:
  diff          Diff live version against would-be applied version
  apply         Apply a configuration to a resource by filename or stdin
  patch         Update field(s) of a resource
  replace       Replace a resource by filename or stdin
  wait          Experimental: Wait for a specific condition on one or many resources.
  kustomize     Build a kustomization target from a directory or URL.

Settings Commands:
  label         Update the labels on a resource
  annotate      Update the annotations on a resource
  completion    Output shell completion code for the specified shell (bash or zsh)

Other Commands:
  api-resources Print the supported API resources on the server
  api-versions  Print the supported API versions on the server, in the form of "group/version"
  config        Modify kubeconfig files
  plugin        Provides utilities for interacting with plugins.
  version       Print the client and server version information

Usage:
  kubectl [flags] [options]

Use "kubectl <command> --help" for more information about a given command.
Use "kubectl options" for a list of global command-line options (applies to all commands).

$ microk8s kubectl version
Client Version: version.Info{Major:"1", Minor:"21", GitVersion:"v1.21.2", GitCommit:"092fbfbf53427de67cac1e9fa54aaa09a28371d7", GitTreeState:"clean", BuildDate:"2021-06-16T12:52:14Z", GoVersion:"go1.16.5", Compiler:"gc", Platform:"darwin/amd64"}
Server Version: version.Info{Major:"1", Minor:"21+", GitVersion:"v1.21.1-3+1f02fea99e2268", GitCommit:"1f02fea99e226878ded82f3e159f28a6706ce1fc", GitTreeState:"clean", BuildDate:"2021-05-12T21:02:46Z", GoVersion:"go1.16.3", Compiler:"gc", Platform:"linux/amd64"}

$ kubectl version
Client Version: version.Info{Major:"1", Minor:"21", GitVersion:"v1.21.2", GitCommit:"092fbfbf53427de67cac1e9fa54aaa09a28371d7", GitTreeState:"clean", BuildDate:"2021-06-16T12:52:14Z", GoVersion:"go1.16.5", Compiler:"gc", Platform:"darwin/amd64"}
The connection to the server localhost:8080 was refused - did you specify the right host or port?

$ kubectl version --client
Client Version: version.Info{Major:"1", Minor:"21", GitVersion:"v1.21.2", GitCommit:"092fbfbf53427de67cac1e9fa54aaa09a28371d7", GitTreeState:"clean", BuildDate:"2021-06-16T12:52:14Z", GoVersion:"go1.16.5", Compiler:"gc", Platform:"darwin/amd64"}
```

https://microk8s.io/docs/working-with-kubectl#heading--kubectl-macos

```bash
$ microk8s config
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUREekNDQWZlZ0F3SUJBZ0lVR3VyOVBTd25IaURXcDROS1JhQkZCaWMvR0ZJd0RRWUpLb1pJaHZjTkFRRUwKQlFBd0Z6RVZNQk1HQTFVRUF3d01NVEF1TVRVeUxqRTRNeTR4TUI0WERUSXhNRGN3TmpFM01EUTBOMW9YRFRNeApNRGN3TkRFM01EUTBOMW93RnpFVk1CTUdBMVVFQXd3TU1UQXVNVFV5TGpFNE15NHhNSUlCSWpBTkJna3Foa2lHCjl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUFyT0V0VlJSY1NaTys5ZDBFYld3Sk5iZGdHQW4wMzRFK3RiUjAKaUdRejlCNmxXdWYyT2xaQnN2NUcyUVNKSkdqY2o2WWd4anR6NU1zZVYxWnJaL08ydDZsU0dYVWtZbzk3eGR4Vgp6V2gvY29YRTBXbTRkc1JTYldoVXNXbzFhSWo2ZkV0ek80eitMZ01JWk5MczgzZHR5ZEpWUWg3Qy9YYXJaeHNyCjdHeUZqbG9pN0VaTzZuUm9nWHR5cWpRNFFURGlTNFpFdmxmWTR3c3oxM09iY1VYWFRtRC9ZT1I1ZnRMRG9TWmkKM0s4bzA1Q0tldi9GR2YvL2t1VFRQZ0ZOS09pRFNXM2xFL09HMHBRYldISlNJM3R4dG5COHMxOVdkUTB2WTY1SwpleEl4Z25wak5BdDJSdzRKaENFWlN0aU84dVB5aFMyWmtlRU55emZST3laL1RZWm1EUUlEQVFBQm8xTXdVVEFkCkJnTlZIUTRFRmdRVXRsMkhFZkhDM3NKeDc4c2IzSVJzWk00WHlLc3dId1lEVlIwakJCZ3dGb0FVdGwySEVmSEMKM3NKeDc4c2IzSVJzWk00WHlLc3dEd1lEVlIwVEFRSC9CQVV3QXdFQi96QU5CZ2txaGtpRzl3MEJBUXNGQUFPQwpBUUVBWHpaYUdqL244UmpzcVBrQ1EveFRZN3lMVWdLLzUwaUZHTEV6ZWNZMWxwWlhRZUtHZkp2dmd0MWp6TVRYClF1M1greGlkcGN3TkZhM0dDclZyRWdSd2U4djNLc0tFdnFjMUcrYldRNUlvc1RoZTdLS1hWaE9CWFR4ei8yZSsKWEJNSHoybHgrNW5PdXB5SkRUaWZsY3dFbzFoK1VWSUtkcWkvTW1JM2QwSVRYVG1jejZ3dlJWWGJVZ2JFOGJ4OQo4OVZlKzJLdXlPajFaMTRTTDZpM1NZSVdQZnhuNzNOWEtWZFJkZ1VUek5abWpQU1c5NkFxRlpncU1idC90clFHCnY5eE1BOEZnQ1VXNUp0ZlAxWWdnMXJuL0ZmV0VmalQ0V0g2aEFxaEQ5RzQ2U3lNUTdYakdSQ1FIN0xZUVhnaG8KLzBTVEg5M3FUeVU3NXhtNFEzU1BXMXlWOUE9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
    server: https://192.168.64.41:16443
  name: microk8s-cluster
contexts:
- context:
    cluster: microk8s-cluster
    user: admin
  name: microk8s
current-context: microk8s
kind: Config
preferences: {}
users:
- name: admin
  user:
    token: VVh6YkNjQ1ZkVDJ1V2h4dTg2SzFqa0xQREhaUGNzbzQwQkRqeVNsZU5iUT0K
```

Interesting! Hmm

```bash
$ microk8s kubectl get pods -A
NAMESPACE     NAME                                      READY   STATUS    RESTARTS   AGE
kube-system   calico-node-75pw9                         1/1     Running   0          4m29s
kube-system   calico-kube-controllers-f7868dd95-8mqrg   1/1     Running   0          4m29s

$ microk8s kubectl top
Display Resource (CPU/Memory) usage.

 The top command allows you to see the resource consumption for nodes or pods.

 This command requires Metrics Server to be correctly configured and working on the server.

Available Commands:
  node        Display Resource (CPU/Memory) usage of nodes
  pod         Display Resource (CPU/Memory) usage of pods

Usage:
  kubectl top [flags] [options]

Use "kubectl <command> --help" for more information about a given command.
Use "kubectl options" for a list of global command-line options (applies to all commands).

$ microk8s kubectl top node
W0706 22:40:03.021694   25351 top_node.go:119] Using json format to get metrics. Next release will switch to protocol-buffers, switch early by passing --use-protocol-buffers flag
error: Metrics API not available
```

```bash
$ microk8s enable
Usage: microk8s enable [OPTIONS] ADDONS...

Error: Missing argument "addons".
An error occurred when trying to execute 'sudo microk8s.enable' with 'multipass': returned exit code 2.

$ microk8s enable -h
Usage: microk8s enable [OPTIONS] ADDONS...

  Enables a MicroK8s addon.

  For a list of available addons, run `microk8s status`.

  To see help for individual addons, run:

      microk8s enable ADDON -- --help

Options:
  --help  Show this message and exit.

$ microk8s enable --help
Usage: microk8s enable [OPTIONS] ADDONS...

  Enables a MicroK8s addon.

  For a list of available addons, run `microk8s status`.

  To see help for individual addons, run:

      microk8s enable ADDON -- --help

Options:
  --help  Show this message and exit.
```

```bash
$ microk8s status 
microk8s is running
high-availability: no
  datastore master nodes: 127.0.0.1:19001
  datastore standby nodes: none
addons:
  enabled:
    ha-cluster           # Configure high availability on the current node
  disabled:
    ambassador           # Ambassador API Gateway and Ingress
    cilium               # SDN, fast with full network policy
    dashboard            # The Kubernetes dashboard
    dns                  # CoreDNS
    fluentd              # Elasticsearch-Fluentd-Kibana logging and monitoring
    gpu                  # Automatic enablement of Nvidia CUDA
    helm                 # Helm 2 - the package manager for Kubernetes
    helm3                # Helm 3 - Kubernetes package manager
    host-access          # Allow Pods connecting to Host services smoothly
    ingress              # Ingress controller for external access
    istio                # Core Istio service mesh services
    jaeger               # Kubernetes Jaeger operator with its simple config
    keda                 # Kubernetes-based Event Driven Autoscaling
    knative              # The Knative framework on Kubernetes.
    kubeflow             # Kubeflow for easy ML deployments
    linkerd              # Linkerd is a service mesh for Kubernetes and other frameworks
    metallb              # Loadbalancer for your Kubernetes cluster
    metrics-server       # K8s Metrics Server for API access to service metrics
    multus               # Multus CNI enables attaching multiple network interfaces to pods
    openebs              # OpenEBS is the open-source storage solution for Kubernetes
    openfaas             # openfaas serverless framework
    portainer            # Portainer UI for your Kubernetes cluster
    prometheus           # Prometheus operator for monitoring and logging
    rbac                 # Role-Based Access Control for authorisation
    registry             # Private image registry exposed on localhost:32000
    storage              # Storage class; allocates storage from host directory
    traefik              # traefik Ingress controller for external access
```

Let me enable the metrics server, hmm

```bash
$ microk8s enable metrics-server
Enabling Metrics-Server
clusterrole.rbac.authorization.k8s.io/system:aggregated-metrics-reader created
clusterrolebinding.rbac.authorization.k8s.io/metrics-server:system:auth-delegator created
rolebinding.rbac.authorization.k8s.io/metrics-server-auth-reader created
Warning: apiregistration.k8s.io/v1beta1 APIService is deprecated in v1.19+, unavailable in v1.22+; use apiregistration.k8s.io/v1 APIService
apiservice.apiregistration.k8s.io/v1beta1.metrics.k8s.io created
serviceaccount/metrics-server created
deployment.apps/metrics-server created
service/metrics-server created
clusterrole.rbac.authorization.k8s.io/system:metrics-server created
clusterrolebinding.rbac.authorization.k8s.io/system:metrics-server created
clusterrolebinding.rbac.authorization.k8s.io/microk8s-admin created
Metrics-Server is enabled
```

```bash
$ microk8s kubectl get pods -A
NAMESPACE     NAME                                      READY   STATUS    RESTARTS   AGE
kube-system   metrics-server-8bbfb4bdb-wggnw            1/1     Running   0          45s
kube-system   calico-node-75pw9                         1/1     Running   0          8m27s
kube-system   calico-kube-controllers-f7868dd95-8mqrg   1/1     Running   0          8m27s
```

For sometime the metrics wasn't available / the metrics-server wasn't ready

Now it's ready! :)

```bash
$ microk8s kubectl top node
W0706 22:44:46.133867   25630 top_node.go:119] Using json format to get metrics. Next release will switch to protocol-buffers, switch early by passing --use-protocol-buffers flag
NAME          CPU(cores)   CPU%   MEMORY(bytes)   MEMORY%   
microk8s-vm   496m         24%    1232Mi          32%       
```

```bash
$ microk8s kubectl top pod
W0706 22:45:09.969725   25653 top_pod.go:140] Using json format to get metrics. Next release will switch to protocol-buffers, switch early by passing --use-protocol-buffers flag
No resources found in default namespace.

$ microk8s kubectl top pod -A
W0706 22:45:14.570716   25661 top_pod.go:140] Using json format to get metrics. Next release will switch to protocol-buffers, switch early by passing --use-protocol-buffers flag
NAMESPACE     NAME                                      CPU(cores)   MEMORY(bytes)   
kube-system   calico-kube-controllers-f7868dd95-8mqrg   2m           6Mi             
kube-system   calico-node-75pw9                         39m          8Mi             
kube-system   metrics-server-8bbfb4bdb-wggnw            1m           10Mi   
```


```bash
$ multipass ls
Name                    State             IPv4             Image
microk8s-vm             Running           192.168.64.41    Ubuntu 18.04 LTS

$ multipass info microk8s-vm
Name:           microk8s-vm
State:          Running
IPv4:           192.168.64.41
Release:        Ubuntu 18.04.5 LTS
Image hash:     50c38d3f7307 (Ubuntu 18.04 LTS)
Load:           0.61 0.88 0.64
Disk usage:     2.3G out of 48.3G
Memory usage:   770.4M out of 3.9G
```

Looks like internally it does use `multipass` VMs, hmm

It's interesting that `multipass` is not mentioned as a dependency in `brew` formula though, hmm

```bash
$ brew info ubuntu/microk8s/microk8s
ubuntu/microk8s/microk8s: stable 2.1.0
Small, fast, single-package Kubernetes for developers, IoT and edge
https://microk8s.io/
/usr/local/Cellar/microk8s/2.1.0 (1,564 files, 15.9MB) *
  Built from source on 2021-07-06 at 22:29:30
From: https://github.com/ubuntu/homebrew-microk8s/blob/HEAD/Formula/microk8s.rb
==> Dependencies
Required: python ✔, kubernetes-cli ✔
==> Requirements
Required: macOS >= 10.12 ✔
==> Caveats
Run `microk8s install` to start with MicroK8s
```

```bash

$ multipass ls
Name                    State             IPv4             Image
microk8s-vm             Running           192.168.64.41    Ubuntu 18.04 LTS

$ microk8s 
Usage: microk8s [OPTIONS] COMMAND [ARGS]...

Options:
  --help  Shows the available COMMANDS.

$ microk8s --help
Usage: microk8s [OPTIONS] COMMAND [ARGS]...

Options:
  --help  Show this message and exit.

Commands:
  install         Installs MicroK8s. Use --cpu, --mem, --disk and --channel to configure your setup.
  uninstall       Removes MicroK8s
  add-node        Adds a node to a cluster
  cilium          The cilium client
  config          Print the kubeconfig
  ctr             The containerd client
  dashboard-proxy Enable the Kubernetes dashboard and proxy to host
  dbctl           Backup and restore the Kubernetes datastore
  disable         Disables running add-ons
  enable          Enables useful add-ons
  helm            The helm client
  helm3           The helm3 client
  inspect         Checks the cluster and gathers logs
  istioctl        The istio client
  join            Joins this instance as a node to a cluster
  juju            The Juju client
  kubectl         The kubernetes client
  leave           Disconnects this node from any cluster it has joined
  linkerd         The linkerd client
  refresh-certs   Refresh the CA certificates in this deployment
  remove-node     Removes a node from the cluster
  reset           Cleans the cluster from all workloads
  start           Starts the kubernetes cluster
  status          Displays the status of the cluster
  stop            Stops the kubernetes cluster

$ microk8s uninstall
Stopping microk8s-vm /[2021-07-06T22:47:47.085] [error] [microk8s-vm] process error occurred Crashed

Thank you for using MicroK8s!                                                   

$ multipass ls
No instances found.
```

