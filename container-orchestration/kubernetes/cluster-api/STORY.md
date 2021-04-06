# Story

I'm trying out the Cluster API. I'm starting from the official doc by following
the quick start guide here

https://cluster-api.sigs.k8s.io/user/quick-start.html

```bash

$ kind create cluster

$ kubectl cluster-info

$ curl -L https://github.com/kubernetes-sigs/cluster-api/releases/download/v0.3.15/clusterctl-darwin-amd64 -o clusterctl

$ chmod +x ./clusterctl

$ sudo mv ./clusterctl /usr/local/bin/clusterctl

$ clusterctl version
clusterctl version: &version.Info{Major:"0", Minor:"3", GitVersion:"v0.3.15", GitCommit:"b900c6f89f3d433c32db1aa269f77f004a83cc4f", GitTreeState:"clean", BuildDate:"2021-03-30T16:14:03Z", GoVersion:"go1.13.15", Compiler:"gc", Platform:"darwin/amd64"}
```

```bash
$ clusterctl init
Fetching providers
Installing cert-manager Version="v1.1.0"
Waiting for cert-manager to be available...
Installing Provider="cluster-api" Version="v0.3.15" TargetNamespace="capi-system"
Installing Provider="bootstrap-kubeadm" Version="v0.3.15" TargetNamespace="capi-kubeadm-bootstrap-system"
Installing Provider="control-plane-kubeadm" Version="v0.3.15" TargetNamespace="capi-kubeadm-control-plane-system"

Your management cluster has been initialized successfully!

You can now create your first workload cluster by running the following:

  clusterctl config cluster [name] --kubernetes-version [version] | kubectl apply -f -
```

```bash
$ clusterctl init --infrastructure docker
Fetching providers
Skipping installing cert-manager as it is already installed
Installing Provider="infrastructure-docker" Version="v0.3.15" TargetNamespace="capd-system"
```

We need to create a workload cluster now

https://cluster-api.sigs.k8s.io/reference/glossary.html#workload-cluster

```bash
$ clusterctl config cluster capi-quickstart --flavor development \
  --kubernetes-version v1.18.16 \
  --control-plane-machine-count=3 \
  --worker-machine-count=3 \
  > capi-quickstart.yaml
```


```bash
$ kubectl apply -f capi-quickstart.yaml

cluster.cluster.x-k8s.io/capi-quickstart created
dockercluster.infrastructure.cluster.x-k8s.io/capi-quickstart created
dockermachinetemplate.infrastructure.cluster.x-k8s.io/capi-quickstart-control-plane created
kubeadmcontrolplane.controlplane.cluster.x-k8s.io/capi-quickstart-control-plane created
dockermachinetemplate.infrastructure.cluster.x-k8s.io/capi-quickstart-md-0 created
kubeadmconfigtemplate.bootstrap.cluster.x-k8s.io/capi-quickstart-md-0 created
machinedeployment.cluster.x-k8s.io/capi-quickstart-md-0 created
```

```bash
$ kubectl get cluster --all-namespaces
NAMESPACE   NAME              PHASE
default     capi-quickstart   Provisioning
```

```bash
$ clusterctl describe cluster capi-quickstart

NAME                                                                READY  SEVERITY  REASON                           SINCE  MESSAGE                                                                            
/capi-quickstart                                                    False  Info      WaitingForControlPlane           51s                                                                                       
├─ClusterInfrastructure - DockerCluster/capi-quickstart                                                                                                                                                         
├─ControlPlane - KubeadmControlPlane/capi-quickstart-control-plane                                                                                                                                              
└─Workers                                                                                                                                                                                                       
  └─MachineDeployment/capi-quickstart-md-0                                                                                                                                                                      
    └─3 Machines...                                                 False  Info      WaitingForClusterInfrastructure  51s    See capi-quickstart-md-0-7cc765486-77mrb, capi-quickstart-md-0-7cc765486-f46dm, ...
```

```bash
$ kubectl describe cluster
Name:         capi-quickstart
Namespace:    default
Labels:       <none>
Annotations:  <none>
API Version:  cluster.x-k8s.io/v1alpha3
Kind:         Cluster
Metadata:
  Creation Timestamp:  2021-04-06T15:47:48Z
  Finalizers:
    cluster.cluster.x-k8s.io
  Generation:  1
  Managed Fields:
    API Version:  cluster.x-k8s.io/v1alpha3
    Fields Type:  FieldsV1
    fieldsV1:
      f:metadata:
        f:annotations:
          .:
          f:kubectl.kubernetes.io/last-applied-configuration:
      f:spec:
        .:
        f:clusterNetwork:
          .:
          f:pods:
            .:
            f:cidrBlocks:
          f:serviceDomain:
          f:services:
            .:
            f:cidrBlocks:
        f:controlPlaneRef:
          .:
          f:apiVersion:
          f:kind:
          f:name:
          f:namespace:
        f:infrastructureRef:
          .:
          f:apiVersion:
          f:kind:
          f:name:
          f:namespace:
    Manager:      kubectl-client-side-apply
    Operation:    Update
    Time:         2021-04-06T15:47:48Z
    API Version:  cluster.x-k8s.io/v1alpha3
    Fields Type:  FieldsV1
    fieldsV1:
      f:metadata:
        f:finalizers:
          .:
          v:"cluster.cluster.x-k8s.io":
      f:status:
        .:
        f:conditions:
        f:observedGeneration:
        f:phase:
    Manager:         manager
    Operation:       Update
    Time:            2021-04-06T15:47:50Z
  Resource Version:  12172
  UID:               8a1487d3-5807-4dd3-b4f5-3713affdd05f
Spec:
  Cluster Network:
    Pods:
      Cidr Blocks:
        192.168.0.0/16
    Service Domain:  cluster.local
    Services:
      Cidr Blocks:
        10.128.0.0/12
  Control Plane Endpoint:
    Host:  
    Port:  0
  Control Plane Ref:
    API Version:  controlplane.cluster.x-k8s.io/v1alpha3
    Kind:         KubeadmControlPlane
    Name:         capi-quickstart-control-plane
    Namespace:    default
  Infrastructure Ref:
    API Version:  infrastructure.cluster.x-k8s.io/v1alpha3
    Kind:         DockerCluster
    Name:         capi-quickstart
    Namespace:    default
Status:
  Conditions:
    Last Transition Time:  2021-04-06T15:47:50Z
    Reason:                WaitingForControlPlane
    Severity:              Info
    Status:                False
    Type:                  Ready
    Last Transition Time:  2021-04-06T15:47:50Z
    Reason:                WaitingForControlPlane
    Severity:              Info
    Status:                False
    Type:                  ControlPlaneReady
    Last Transition Time:  2021-04-06T15:47:50Z
    Reason:                WaitingForInfrastructure
    Severity:              Info
    Status:                False
    Type:                  InfrastructureReady
  Observed Generation:     1
  Phase:                   Provisioning
Events:                    <none>
```

For some reason it's waiting for some things to work. I'm checking why and what

```bash
$ kubectl get kubeadmcontrolplane
NAME                            INITIALIZED   API SERVER AVAILABLE   VERSION    REPLICAS   READY   UPDATED   UNAVAILABLE
capi-quickstart-control-plane                                        v1.18.16                                
```

```bash
$ kubectl get machinedeployment
NAME                   PHASE       REPLICAS   READY   UPDATED   UNAVAILABLE
capi-quickstart-md-0   ScalingUp   3                  3         3
```

```bash
$ kubectl describe machinedeployment
Name:         capi-quickstart-md-0
Namespace:    default
Labels:       cluster.x-k8s.io/cluster-name=capi-quickstart
Annotations:  machinedeployment.clusters.x-k8s.io/revision: 1
API Version:  cluster.x-k8s.io/v1alpha3
Kind:         MachineDeployment
Metadata:
  Creation Timestamp:  2021-04-06T15:47:48Z
  Generation:          1
  Managed Fields:
    API Version:  cluster.x-k8s.io/v1alpha3
    Fields Type:  FieldsV1
    fieldsV1:
      f:metadata:
        f:annotations:
          .:
          f:kubectl.kubernetes.io/last-applied-configuration:
      f:spec:
        .:
        f:clusterName:
        f:replicas:
        f:selector:
        f:template:
          .:
          f:spec:
            .:
            f:bootstrap:
              .:
              f:configRef:
                .:
                f:apiVersion:
                f:kind:
                f:name:
                f:namespace:
            f:clusterName:
            f:infrastructureRef:
              .:
              f:apiVersion:
              f:kind:
              f:name:
              f:namespace:
            f:version:
    Manager:      kubectl-client-side-apply
    Operation:    Update
    Time:         2021-04-06T15:47:48Z
    API Version:  cluster.x-k8s.io/v1alpha3
    Fields Type:  FieldsV1
    fieldsV1:
      f:metadata:
        f:annotations:
          f:machinedeployment.clusters.x-k8s.io/revision:
        f:ownerReferences:
      f:status:
        .:
        f:observedGeneration:
        f:phase:
        f:replicas:
        f:selector:
        f:unavailableReplicas:
        f:updatedReplicas:
    Manager:    manager
    Operation:  Update
    Time:       2021-04-06T15:47:50Z
  Owner References:
    API Version:     cluster.x-k8s.io/v1alpha3
    Kind:            Cluster
    Name:            capi-quickstart
    UID:             8a1487d3-5807-4dd3-b4f5-3713affdd05f
  Resource Version:  12206
  UID:               c45ab242-63f1-49ad-b563-3832fe091848
Spec:
  Cluster Name:               capi-quickstart
  Min Ready Seconds:          0
  Progress Deadline Seconds:  600
  Replicas:                   3
  Revision History Limit:     1
  Selector:
    Match Labels:
      cluster.x-k8s.io/cluster-name:     capi-quickstart
      cluster.x-k8s.io/deployment-name:  capi-quickstart-md-0
  Strategy:
    Rolling Update:
      Max Surge:        1
      Max Unavailable:  0
    Type:               RollingUpdate
  Template:
    Metadata:
      Labels:
        cluster.x-k8s.io/cluster-name:     capi-quickstart
        cluster.x-k8s.io/deployment-name:  capi-quickstart-md-0
    Spec:
      Bootstrap:
        Config Ref:
          API Version:  bootstrap.cluster.x-k8s.io/v1alpha3
          Kind:         KubeadmConfigTemplate
          Name:         capi-quickstart-md-0
          Namespace:    default
      Cluster Name:     capi-quickstart
      Infrastructure Ref:
        API Version:  infrastructure.cluster.x-k8s.io/v1alpha3
        Kind:         DockerMachineTemplate
        Name:         capi-quickstart-md-0
        Namespace:    default
      Version:        v1.18.16
Status:
  Observed Generation:   1
  Phase:                 ScalingUp
  Replicas:              3
  Selector:              cluster.x-k8s.io/cluster-name=capi-quickstart,cluster.x-k8s.io/deployment-name=capi-quickstart-md-0
  Unavailable Replicas:  3
  Updated Replicas:      3
Events:
  Type    Reason            Age    From                          Message
  ----    ------            ----   ----                          -------
  Normal  SuccessfulCreate  7m41s  machinedeployment-controller  Created MachineSet "capi-quickstart-md-0-7cc765486"
```

```bash
$ kubectl get dockermachine
NAME                         AGE
capi-quickstart-md-0-kzdmw   11m
capi-quickstart-md-0-nmt7d   11m
capi-quickstart-md-0-zrchc   11m
```

```bash
$ kubectl get machine
NAME                                   PROVIDERID   PHASE     VERSION
capi-quickstart-md-0-7cc765486-77mrb                Pending   v1.18.16
capi-quickstart-md-0-7cc765486-f46dm                Pending   v1.18.16
capi-quickstart-md-0-7cc765486-qqk89                Pending   v1.18.16
```

```bash
 cluster-api  master ✘  $ kubectl api-resources
NAME                              SHORTNAMES   APIVERSION                                     NAMESPACED   KIND
bindings                                       v1                                             true         Binding
componentstatuses                 cs           v1                                             false        ComponentStatus
configmaps                        cm           v1                                             true         ConfigMap
endpoints                         ep           v1                                             true         Endpoints
events                            ev           v1                                             true         Event
limitranges                       limits       v1                                             true         LimitRange
namespaces                        ns           v1                                             false        Namespace
nodes                             no           v1                                             false        Node
persistentvolumeclaims            pvc          v1                                             true         PersistentVolumeClaim
persistentvolumes                 pv           v1                                             false        PersistentVolume
pods                              po           v1                                             true         Pod
podtemplates                                   v1                                             true         PodTemplate
replicationcontrollers            rc           v1                                             true         ReplicationController
resourcequotas                    quota        v1                                             true         ResourceQuota
secrets                                        v1                                             true         Secret
serviceaccounts                   sa           v1                                             true         ServiceAccount
services                          svc          v1                                             true         Service
challenges                                     acme.cert-manager.io/v1                        true         Challenge
orders                                         acme.cert-manager.io/v1                        true         Order
clusterresourcesetbindings                     addons.cluster.x-k8s.io/v1alpha3               true         ClusterResourceSetBinding
clusterresourcesets                            addons.cluster.x-k8s.io/v1alpha3               true         ClusterResourceSet
mutatingwebhookconfigurations                  admissionregistration.k8s.io/v1                false        MutatingWebhookConfiguration
validatingwebhookconfigurations                admissionregistration.k8s.io/v1                false        ValidatingWebhookConfiguration
customresourcedefinitions         crd,crds     apiextensions.k8s.io/v1                        false        CustomResourceDefinition
apiservices                                    apiregistration.k8s.io/v1                      false        APIService
controllerrevisions                            apps/v1                                        true         ControllerRevision
daemonsets                        ds           apps/v1                                        true         DaemonSet
deployments                       deploy       apps/v1                                        true         Deployment
replicasets                       rs           apps/v1                                        true         ReplicaSet
statefulsets                      sts          apps/v1                                        true         StatefulSet
tokenreviews                                   authentication.k8s.io/v1                       false        TokenReview
localsubjectaccessreviews                      authorization.k8s.io/v1                        true         LocalSubjectAccessReview
selfsubjectaccessreviews                       authorization.k8s.io/v1                        false        SelfSubjectAccessReview
selfsubjectrulesreviews                        authorization.k8s.io/v1                        false        SelfSubjectRulesReview
subjectaccessreviews                           authorization.k8s.io/v1                        false        SubjectAccessReview
horizontalpodautoscalers          hpa          autoscaling/v1                                 true         HorizontalPodAutoscaler
cronjobs                          cj           batch/v1beta1                                  true         CronJob
jobs                                           batch/v1                                       true         Job
kubeadmconfigs                                 bootstrap.cluster.x-k8s.io/v1alpha3            true         KubeadmConfig
kubeadmconfigtemplates                         bootstrap.cluster.x-k8s.io/v1alpha3            true         KubeadmConfigTemplate
certificaterequests               cr,crs       cert-manager.io/v1                             true         CertificateRequest
certificates                      cert,certs   cert-manager.io/v1                             true         Certificate
clusterissuers                                 cert-manager.io/v1                             false        ClusterIssuer
issuers                                        cert-manager.io/v1                             true         Issuer
certificatesigningrequests        csr          certificates.k8s.io/v1                         false        CertificateSigningRequest
clusters                          cl           cluster.x-k8s.io/v1alpha3                      true         Cluster
machinedeployments                md           cluster.x-k8s.io/v1alpha3                      true         MachineDeployment
machinehealthchecks               mhc,mhcs     cluster.x-k8s.io/v1alpha3                      true         MachineHealthCheck
machines                          ma           cluster.x-k8s.io/v1alpha3                      true         Machine
machinesets                       ms           cluster.x-k8s.io/v1alpha3                      true         MachineSet
providers                                      clusterctl.cluster.x-k8s.io/v1alpha3           true         Provider
kubeadmcontrolplanes              kcp          controlplane.cluster.x-k8s.io/v1alpha3         true         KubeadmControlPlane
leases                                         coordination.k8s.io/v1                         true         Lease
endpointslices                                 discovery.k8s.io/v1beta1                       true         EndpointSlice
events                            ev           events.k8s.io/v1                               true         Event
machinepools                      mp           exp.cluster.x-k8s.io/v1alpha3                  true         MachinePool
dockermachinepools                             exp.infrastructure.cluster.x-k8s.io/v1alpha3   true         DockerMachinePool
ingresses                         ing          extensions/v1beta1                             true         Ingress
flowschemas                                    flowcontrol.apiserver.k8s.io/v1beta1           false        FlowSchema
prioritylevelconfigurations                    flowcontrol.apiserver.k8s.io/v1beta1           false        PriorityLevelConfiguration
dockerclusters                                 infrastructure.cluster.x-k8s.io/v1alpha3       true         DockerCluster
dockermachines                                 infrastructure.cluster.x-k8s.io/v1alpha3       true         DockerMachine
dockermachinetemplates                         infrastructure.cluster.x-k8s.io/v1alpha3       true         DockerMachineTemplate
ingressclasses                                 networking.k8s.io/v1                           false        IngressClass
ingresses                         ing          networking.k8s.io/v1                           true         Ingress
networkpolicies                   netpol       networking.k8s.io/v1                           true         NetworkPolicy
runtimeclasses                                 node.k8s.io/v1                                 false        RuntimeClass
poddisruptionbudgets              pdb          policy/v1beta1                                 true         PodDisruptionBudget
podsecuritypolicies               psp          policy/v1beta1                                 false        PodSecurityPolicy
clusterrolebindings                            rbac.authorization.k8s.io/v1                   false        ClusterRoleBinding
clusterroles                                   rbac.authorization.k8s.io/v1                   false        ClusterRole
rolebindings                                   rbac.authorization.k8s.io/v1                   true         RoleBinding
roles                                          rbac.authorization.k8s.io/v1                   true         Role
priorityclasses                   pc           scheduling.k8s.io/v1                           false        PriorityClass
csidrivers                                     storage.k8s.io/v1                              false        CSIDriver
csinodes                                       storage.k8s.io/v1                              false        CSINode
storageclasses                    sc           storage.k8s.io/v1                              false        StorageClass
volumeattachments                              storage.k8s.io/v1                              false        VolumeAttachment
 cluster-api  master ✘  $ 
```


Look at these -

```bash
clusterresourcesetbindings                     addons.cluster.x-k8s.io/v1alpha3               true         ClusterResourceSetBinding
clusterresourcesets                            addons.cluster.x-k8s.io/v1alpha3               true         ClusterResourceSet
kubeadmconfigs                                 bootstrap.cluster.x-k8s.io/v1alpha3            true         KubeadmConfig
kubeadmconfigtemplates                         bootstrap.cluster.x-k8s.io/v1alpha3            true         KubeadmConfigTemplate
clusters                          cl           cluster.x-k8s.io/v1alpha3                      true         Cluster
machinedeployments                md           cluster.x-k8s.io/v1alpha3                      true         MachineDeployment
machinehealthchecks               mhc,mhcs     cluster.x-k8s.io/v1alpha3                      true         MachineHealthCheck
machines                          ma           cluster.x-k8s.io/v1alpha3                      true         Machine
machinesets                       ms           cluster.x-k8s.io/v1alpha3                      true         MachineSet
providers                                      clusterctl.cluster.x-k8s.io/v1alpha3           true         Provider
kubeadmcontrolplanes              kcp          controlplane.cluster.x-k8s.io/v1alpha3         true         KubeadmControlPlane
machinepools                      mp           exp.cluster.x-k8s.io/v1alpha3                  true         MachinePool
dockermachinepools                             exp.infrastructure.cluster.x-k8s.io/v1alpha3   true         DockerMachinePool
dockerclusters                                 infrastructure.cluster.x-k8s.io/v1alpha3       true         DockerCluster
dockermachines                                 infrastructure.cluster.x-k8s.io/v1alpha3       true         DockerMachine
dockermachinetemplates                         infrastructure.cluster.x-k8s.io/v1alpha3       true         DockerMachineTemplate
```

I think it's all related to Cluster API. The `x-k8s.io` shows it, also the
`cluster.x-k8s.io`. Hmm

Interesting stuff. But weirdly somehow the Kubernetes cluster is not up and running, hmm

I'm planning to delete the cluster and try again

But it's not getting deleted

```bash
$ kubectl delete cluster capi-quickstart
cluster.cluster.x-k8s.io "capi-quickstart" deleted
...
```

It's stuck. The command's stuck

```bash
$ kubectl get cluster 
NAME              PHASE
capi-quickstart   Deleting

$ kubectl get machine
NAME                                   PROVIDERID   PHASE      VERSION
capi-quickstart-md-0-7cc765486-77mrb                Deleting   v1.18.16
capi-quickstart-md-0-7cc765486-f46dm                Deleting   v1.18.16
capi-quickstart-md-0-7cc765486-qqk89                Deleting   v1.18.16
```

Hmm. I guess I could just delete the whole kind management cluster ? Hmm

```bash
$ kind delete cluster
Deleting cluster "kind" ...
```

Done!

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

Not sure what to do next? 😅  Check out https://kind.sigs.k8s.io/docs/user/quick-start/
```

```bash
$ kubectl cluster-info --context kind-kind

Kubernetes control plane is running at https://127.0.0.1:63019
KubeDNS is running at https://127.0.0.1:63019/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.

$ kubectl cluster-info 
Kubernetes control plane is running at https://127.0.0.1:63019
KubeDNS is running at https://127.0.0.1:63019/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
```

Oops, I applied the yaml by mistake before initializing the management cluster
with the Cluster API components. Hehe

```bash
$ k apply -f capi-quickstart.yaml 
unable to recognize "capi-quickstart.yaml": no matches for kind "Cluster" in version "cluster.x-k8s.io/v1alpha3"
unable to recognize "capi-quickstart.yaml": no matches for kind "DockerCluster" in version "infrastructure.cluster.x-k8s.io/v1alpha3"
unable to recognize "capi-quickstart.yaml": no matches for kind "DockerMachineTemplate" in version "infrastructure.cluster.x-k8s.io/v1alpha3"
unable to recognize "capi-quickstart.yaml": no matches for kind "KubeadmControlPlane" in version "controlplane.cluster.x-k8s.io/v1alpha3"
unable to recognize "capi-quickstart.yaml": no matches for kind "DockerMachineTemplate" in version "infrastructure.cluster.x-k8s.io/v1alpha3"
unable to recognize "capi-quickstart.yaml": no matches for kind "KubeadmConfigTemplate" in version "bootstrap.cluster.x-k8s.io/v1alpha3"
unable to recognize "capi-quickstart.yaml": no matches for kind "MachineDeployment" in version "cluster.x-k8s.io/v1alpha3"
```

```bash
$ clusterctl init --infrastructure docker
Fetching providers
Installing cert-manager Version="v1.1.0"
Waiting for cert-manager to be available...
Installing Provider="cluster-api" Version="v0.3.15" TargetNamespace="capi-system"
Installing Provider="bootstrap-kubeadm" Version="v0.3.15" TargetNamespace="capi-kubeadm-bootstrap-system"
Installing Provider="control-plane-kubeadm" Version="v0.3.15" TargetNamespace="capi-kubeadm-control-plane-system"
Installing Provider="infrastructure-docker" Version="v0.3.15" TargetNamespace="capd-system"

Your management cluster has been initialized successfully!

You can now create your first workload cluster by running the following:

  clusterctl config cluster [name] --kubernetes-version [version] | kubectl apply -f -
```

```bash
clusterctl config cluster capi-quickstart --flavor development \
  --kubernetes-version v1.18.16 \
  --control-plane-machine-count=3 \
  --worker-machine-count=3 \
  > capi-quickstart.yaml
```

And the yaml is the same old yaml.

```bash
$ kubectl apply -f capi-quickstart.yaml
cluster.cluster.x-k8s.io/capi-quickstart created
dockercluster.infrastructure.cluster.x-k8s.io/capi-quickstart created
dockermachinetemplate.infrastructure.cluster.x-k8s.io/capi-quickstart-control-plane created
kubeadmcontrolplane.controlplane.cluster.x-k8s.io/capi-quickstart-control-plane created
dockermachinetemplate.infrastructure.cluster.x-k8s.io/capi-quickstart-md-0 created
kubeadmconfigtemplate.bootstrap.cluster.x-k8s.io/capi-quickstart-md-0 created
machinedeployment.cluster.x-k8s.io/capi-quickstart-md-0 created
```

Same old thing. Damn. Hmm

```bash
$ kubectl get cluster
NAME              PHASE
capi-quickstart   Provisioning
```

Something big going wrong I guess

```bash
$ clusterctl describe cluster capi-quickstart
NAME                                                                READY  SEVERITY  REASON                           SINCE  MESSAGE                                                                            
/capi-quickstart                                                    False  Info      WaitingForControlPlane           19s                                                                                       
├─ClusterInfrastructure - DockerCluster/capi-quickstart                                                                                                                                                         
├─ControlPlane - KubeadmControlPlane/capi-quickstart-control-plane                                                                                                                                              
└─Workers                                                                                                                                                                                                       
  └─MachineDeployment/capi-quickstart-md-0                                                                                                                                                                      
    └─3 Machines...                                                 False  Info      WaitingForClusterInfrastructure  18s    See capi-quickstart-md-0-7cc765486-2dlkz, capi-quickstart-md-0-7cc765486-r7pxz, ...
```

I have to try something differently I guess

```bash
$ kind delete cluster
Deleting cluster "kind" ...
```
