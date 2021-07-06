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
$ kubectl api-resources
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
$ 
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

---

I just noticed that there's an extra step that I missed to notice 🙈 I was so impatient that I missed it!

One of the lines in the docs says

```
if you are planning to use the docker infrastructure provider, please follow the additional instructions in the dedicated tab:
```

And asks to do this

```bash
$ cat > kind-cluster-with-extramounts.yaml <<EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraMounts:
    - hostPath: /var/run/docker.sock
      containerPath: /var/run/docker.sock
EOF

$ kind create cluster --config kind-cluster-with-extramounts.yaml
```

```bash
clusterctl init --infrastructure docker
```

```bash
$ clusterctl init --infrastructure docker
Fetching providers
Installing cert-manager Version="v1.1.0"
Waiting for cert-manager to be available...
Installing Provider="cluster-api" Version="v0.4.0" TargetNamespace="capi-system"
Installing Provider="bootstrap-kubeadm" Version="v0.4.0" TargetNamespace="capi-kubeadm-bootstrap-system"
Installing Provider="control-plane-kubeadm" Version="v0.4.0" TargetNamespace="capi-kubeadm-control-plane-system"
Installing Provider="infrastructure-docker" Version="v0.4.0" TargetNamespace="capd-system"

Your management cluster has been initialized successfully!

You can now create your first workload cluster by running the following:

  clusterctl generate cluster [name] --kubernetes-version [version] | kubectl apply -f -

Error: unable to verify clusterctl version: unable to semver parse clusterctl GitVersion: strconv.ParseUint: parsing "": invalid syntax

$ clusterctl version
clusterctl version: &version.Info{Major:"", Minor:"", GitVersion:"", GitCommit:"", GitTreeState:"", BuildDate:"", GoVersion:"go1.16.5", Compiler:"gc", Platform:"darwin/amd64"}
Error: unable to verify clusterctl version: unable to semver parse clusterctl GitVersion: strconv.ParseUint: parsing "": invalid syntax

$ brew info clusterctl
clusterctl: stable 0.4.0 (bottled)
Home for the Cluster Management API work, a subproject of sig-cluster-lifecycle
https://cluster-api.sigs.k8s.io
/usr/local/Cellar/clusterctl/0.4.0 (5 files, 50.6MB) *
  Poured from bottle on 2021-07-06 at 21:22:30
From: https://github.com/Homebrew/homebrew-core/blob/HEAD/Formula/clusterctl.rb
License: Apache-2.0
==> Dependencies
Build: go ✔
==> Analytics
install: 105 (30 days), 247 (90 days), 292 (365 days)
install-on-request: 105 (30 days), 247 (90 days), 292 (365 days)
build-error: 0 (30 days)

$ clusterctl version
clusterctl version: &version.Info{Major:"", Minor:"", GitVersion:"", GitCommit:"", GitTreeState:"", BuildDate:"", GoVersion:"go1.16.5", Compiler:"gc", Platform:"darwin/amd64"}
Error: unable to verify clusterctl version: unable to semver parse clusterctl GitVersion: strconv.ParseUint: parsing "": invalid syntax

$ echo $?
1

$ clusterctl generate cluster capi-quickstart --flavor development \
>   --kubernetes-version v1.19.7 \
>   --control-plane-machine-count=3 \
>   --worker-machine-count=3 \
>   > capi-quickstart.yaml

Error: unable to verify clusterctl version: unable to semver parse clusterctl GitVersion: strconv.ParseUint: parsing "": invalid syntax
$ 
$ 
```

Some weird error when doing `brew` installation of `clusterctl`. Hmm

https://github.com/kubernetes-sigs/cluster-api/issues/4802

Anyways, things work though! :)

```bash
$ less capi-quickstart.yaml

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
$ clusterctl describe cluster capi-quickstart
NAME                                                                READY  SEVERITY  REASON                           SINCE  MESSAGE                                                                              
/capi-quickstart                                                    False  Warning   ScalingUp                        6s     Scaling up control plane to 3 replicas (actual 1)                                    
├─ClusterInfrastructure - DockerCluster/capi-quickstart             True                                              17s                                                                                         
├─ControlPlane - KubeadmControlPlane/capi-quickstart-control-plane  False  Warning   ScalingUp                        6s     Scaling up control plane to 3 replicas (actual 1)                                    
│ └─Machine/capi-quickstart-control-plane-7tkxm                     False  Info      WaitingForBootstrapData          16s    1 of 2 completed                                                                     
└─Workers                                                                                                                                                                                                         
  └─MachineDeployment/capi-quickstart-md-0                                                                                                                                                                        
    └─3 Machines...                                                 False  Info      WaitingForControlPlaneAvailable  17s    See capi-quickstart-md-0-665856b6bd-9kdnh, capi-quickstart-md-0-665856b6bd-9lshz, ...
Error: unable to verify clusterctl version: unable to semver parse clusterctl GitVersion: strconv.ParseUint: parsing "": invalid syntax

$ docker ps
CONTAINER ID   IMAGE                          COMMAND                  CREATED          STATUS          PORTS                                NAMES
b787ba4d8fcc   kindest/haproxy:2.1.1-alpine   "/docker-entrypoint.…"   34 seconds ago   Up 29 seconds   33083/tcp, 0.0.0.0:33083->6443/tcp   capi-quickstart-lb
836eb30ce1b9   kindest/node:v1.20.2           "/usr/local/bin/entr…"   19 minutes ago   Up 19 minutes   127.0.0.1:54864->6443/tcp            kind-control-plane
```

```bash
$ kubectl get kubeadmcontrolplane --all-namespaces
NAMESPACE   NAME                            INITIALIZED   API SERVER AVAILABLE   VERSION   REPLICAS   READY   UPDATED   UNAVAILABLE
default     capi-quickstart-control-plane
```

```bash
$ clusterctl get kubeconfig capi-quickstart > capi-quickstart.kubeconfig

$ less capi-quickstart.kubeconfig 
```

I noticed that there are three replicas. And I'm like - my machine is going to die, lol

Next time I gotta try one replica, lolol

```bash
$ kubectl --kubeconfig=./capi-quickstart.kubeconfig \
>   apply -f https://docs.projectcalico.org/v3.15/manifests/calico.yaml

Unable to connect to the server: dial tcp 172.18.0.3:6443: i/o timeout
```

```bash
$ k get nodes
NAME                 STATUS   ROLES                  AGE   VERSION
kind-control-plane   Ready    control-plane,master   25m   v1.20.2

$ kubectl --kubeconfig=./capi-quickstart.kubeconfig   apply -f https://docs.projectcalico.org/v3.15/manifests/calico.yaml
Unable to connect to the server: dial tcp 172.18.0.3:6443: i/o timeout
```

```bash
$ docker ps -a -q | xargs docker rm -f
18e58b3dd3e7
303241a2f05c
5bc8c1b8f8ab
aa8504505c9f
f7960108f931
0282c27bd19e
b787ba4d8fcc
836eb30ce1b9
```

Time to start from scratch :P

Finally after cleaning up lots of stuff, it finally worked

```bash
$ clusterctl describe cluster capi-quickstart
NAME                                                                READY  SEVERITY  REASON         SINCE  MESSAGE         
/capi-quickstart                                                    True                            22s                    
├─ClusterInfrastructure - DockerCluster/capi-quickstart             True                            81s                    
├─ControlPlane - KubeadmControlPlane/capi-quickstart-control-plane  True                            22s                    
│ └─Machine/capi-quickstart-control-plane-ss65t                     True                            28s                    
└─Workers                                                                                                                  
  └─MachineDeployment/capi-quickstart-md-0                                                                                 
    └─Machine/capi-quickstart-md-0-665856b6bd-cx6cp                 False  Info      Bootstrapping  13s    1 of 2 completed
Error: unable to verify clusterctl version: unable to semver parse clusterctl GitVersion: strconv.ParseUint: parsing "": invalid syntax
```

The workload cluster is still not able to boot and run :O Wow

I give up, lol

My Mac with 8 CPUs and 16GB RAM, out of which I provided 4 CPUs and 6 GB RAM and 16 GB disk

I'm not sure what resources wasn't enough. I think it was CPU. Not sure. But I give up, haha

I think I'll try to run some light weight thing - try out other k8s clusers with Docker providers or just use my work laptop which has more power
