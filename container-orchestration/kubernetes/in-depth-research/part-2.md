# Part 2

In this part, let's run the kubernetes scheduler so that the scheduling can
happen and pods can be scheduled and run a particular worker node. Let's not
worry that worker nodes are not present at the moment

```bash
$ kube-scheduler
I1229 11:07:48.909167    1149 serving.go:331] Generated self-signed cert in-memory
W1229 11:07:49.658318    1149 authentication.go:303] No authentication-kubeconfig provided in order to lookup client-ca-file in configmap/extension-apiserver-authentication in kube-system, so client certificate authentication won't work.
W1229 11:07:49.662852    1149 authentication.go:327] No authentication-kubeconfig provided in order to lookup requestheader-client-ca-file in configmap/extension-apiserver-authentication in kube-system, so request-header client certificate authentication won't work.
W1229 11:07:49.663710    1149 authorization.go:173] No authorization-kubeconfig provided, so SubjectAccessReview of authorization tokens won't work.
W1229 11:07:49.664201    1149 options.go:332] Neither --kubeconfig nor --master was specified. Using default API client. This might not work.
invalid configuration: no configuration has been provided, try setting KUBERNETES_MASTER environment variable
```

There's a mention about some config map

```bash
$ kubectl get cm -A
NAMESPACE     NAME                                 DATA   AGE
kube-system   extension-apiserver-authentication   1      37h

$ kubectl get cm -n kube-system extension-apiserver-authentication -o yaml
apiVersion: v1
data:
  client-ca-file: |
    -----BEGIN CERTIFICATE-----
    MIIDsjCCApqgAwIBAgIUX/L3lloZ1YkIWeBo4q6y5lFI+tYwDQYJKoZIhvcNAQEL
    BQAwcTELMAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcT
    DVNhbiBGcmFuY2lzY28xEzARBgNVBAoTCkt1YmVybmV0ZXMxCzAJBgNVBAsTAkNB
    MRMwEQYDVQQDEwpLdWJlcm5ldGVzMB4XDTIwMTIyNzA4NTEwMFoXDTI1MTIyNjA4
    NTEwMFowcTELMAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNV
    BAcTDVNhbiBGcmFuY2lzY28xEzARBgNVBAoTCkt1YmVybmV0ZXMxCzAJBgNVBAsT
    AkNBMRMwEQYDVQQDEwpLdWJlcm5ldGVzMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A
    MIIBCgKCAQEA6rT3qK+RmUU3023P3e7j/IxTJxnHl8epIAXDjmNiDSzD4jeMg8+b
    /QiBZfoZG489cmW1iwauRbM7MiuGjo1nRiTVnJI5puOSda007me88HQQJtPQbKVQ
    2O9W+w8epR9MODhs6StyxjXtaVd0iiSKt5qiq7vz0x4wss4TXdzwh0JR54uVKixN
    xQ54uCnE5DiWwbS9/wsax3NtEk+actXS6eFpv7qHXGsUReCyukU1rh0pvHv0aOd2
    gkG87YGKqdT82Yj5lE1NuD06ueTTJM4NJL2WaB3j69CNhc9g3B7yM6SMdJ5Y5zJJ
    F93U4GxZ9SEyOZpcsExAvV5igf6hZx0egQIDAQABo0IwQDAOBgNVHQ8BAf8EBAMC
    AQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUx31bkRIQTMwWVUiPu+zW/jtu
    qnUwDQYJKoZIhvcNAQELBQADggEBAM8AByMwBJVFxQpf6vOY99+lpULbNXGVRhGN
    RzlJlEF/ZIRctw8aFQs7ivLckIytd/hIPJRWiGIiF4Vx0VqrQUo3ZtT3pjFDAfiQ
    wDQTNtEXOzaMXLQbQUkHrMiXo+3pssNWFJW/4/YTAczMcf/EAhRZAtiqWmyCkEPl
    t765Tqjajw+M+PPPH0QSP50cJ9M1khy1x0xV1hn1uzvJt0eyfNsT1RnkHEm/zA8o
    iLc4o7krFdQc26PoTjro2WEfuuzNvDPLswRE32NWnV+J7irs2zzPffAxMa8v+qwA
    MkFl8pqyDG1DoZACxaoldeTY/alUur5Ib7rRrbFxcHtKKUmnZZQ=
    -----END CERTIFICATE-----
kind: ConfigMap
metadata:
  creationTimestamp: "2020-12-27T15:50:04Z"
  managedFields:
  - apiVersion: v1
    fieldsType: FieldsV1
    fieldsV1:
      f:data:
        .: {}
        f:client-ca-file: {}
    manager: kube-apiserver
    operation: Update
    time: "2020-12-27T15:50:04Z"
  name: extension-apiserver-authentication
  namespace: kube-system
  resourceVersion: "419"
  uid: 56cbb77c-936a-4807-8226-fb5468f1f4a9
```

I'm guessing this is the CA certificate. Let's cross check

```bash
$ sudo apt install jq
```

```bash
$ kubectl get cm -n kube-system extension-apiserver-authentication -o json | jq '.data | ."client-ca-file"' -r | diff - ca.pem
23d22
<
```

So yes, it's the same thing. The CA certificate.

So, it keeps saying that authentication/authorization is not happening. Cool.

Okay, so, we need to create a certificate for the kube-scheduler I guess? Or
provide a kube-config, or both? Cuz client certificate is only the credentials,
but it still needs details on where the API server is - host and port. Kube
Config is a combination of cluster info, user info (credentials/certs etc).

Let's see what the kube-scheduler binary has in the help

```bash
$ kube-scheduler --help | grep -i auth
      --secure-port int                        The port on which to serve HTTPS with authentication and authorization. If 0, don't serve HTTPS at all. (default 10259)
      --port int         DEPRECATED: the port on which to serve HTTP insecurely without authentication and authorization. If 0, don't serve plain HTTP at all. See --secure-port instead. (default 10251)
Authentication flags:
      --authentication-kubeconfig string                  kubeconfig file pointing at the 'core' kubernetes server with enough rights to create tokenreviews.authentication.k8s.io. This is optional. If empty, all token requests are considered to be anonymous and no client CA is looked up in the cluster.
      --authentication-skip-lookup                        If false, the authentication-kubeconfig will be used to lookup missing authentication configuration from the cluster.
      --authentication-token-webhook-cache-ttl duration   The duration to cache responses from the webhook token authenticator. (default 10s)
      --authentication-tolerate-lookup-failure            If true, failures to look up missing authentication configuration from the cluster are not considered fatal. Note that this can result in authentication that treats all requests as anonymous. (default true)
      --client-ca-file string                             If set, any request presenting a client certificate signed by one of the authorities in the client-ca-file is authenticated with an identity corresponding to the CommonName of the client certificate.
      --requestheader-allowed-names strings               List of client certificate common names to allow to provide usernames in headers specified by --requestheader-username-headers. If empty, any client certificate validated by the authorities in --requestheader-client-ca-file is allowed.
      --requestheader-client-ca-file string               Root certificate bundle to use to verify client certificates on incoming requests before trusting usernames in headers specified by --requestheader-username-headers. WARNING: generally do not depend on authorization being already done for incoming requests.
Authorization flags:
      --authorization-always-allow-paths strings                A list of HTTP paths to skip during authorization, i.e. these are authorized without contacting the 'core' kubernetes server. (default [/healthz])
      --authorization-kubeconfig string                         kubeconfig file pointing at the 'core' kubernetes server with enough rights to create subjectaccessreviews.authorization.k8s.io. This is optional. If empty, all requests not skipped by authorization are forbidden.
      --authorization-webhook-cache-authorized-ttl duration     The duration to cache 'authorized' responses from the webhook authorizer. (default 10s)
      --authorization-webhook-cache-unauthorized-ttl duration   The duration to cache 'unauthorized' responses from the webhook authorizer. (default 10s)
      --kubeconfig string                          DEPRECATED: path to kubeconfig file with authorization and master location information.
```

Looks like kube config is deprecated, hmm

```bash
--kubeconfig string                          DEPRECATED: path to kubeconfig file with authorization and master location information.
```

I just noticed in the `--help` that the flags are properly categorized! That's
great! I was too lazy to look at the whole thing 😅 🙈 Anyways, there are so
many categories and it also has a category named as "deprecated flags!" that's
pretty useful! :) :D

```bash
$ kube-scheduler --master https://localhost:6443
I1230 08:37:24.605936    1077 serving.go:331] Generated self-signed cert in-memory
W1230 08:37:25.009897    1077 authentication.go:303] No authentication-kubeconfig provided in order to lookup client-ca-file in configmap/extension-apiserver-authentication in kube-system, so client certificate authentication won't work.
W1230 08:37:25.010212    1077 authentication.go:327] No authentication-kubeconfig provided in order to lookup requestheader-client-ca-file in configmap/extension-apiserver-authentication in kube-system, so request-header client certificate authentication won't work.
W1230 08:37:25.010600    1077 authorization.go:173] No authorization-kubeconfig provided, so SubjectAccessReview of authorization tokens won't work.
W1230 08:37:25.091768    1077 authorization.go:47] Authorization is disabled
W1230 08:37:25.091971    1077 authentication.go:40] Authentication is disabled
I1230 08:37:25.092375    1077 deprecated_insecure_serving.go:51] Serving healthz insecurely on [::]:10251
I1230 08:37:25.094574    1077 secure_serving.go:197] Serving securely on [::]:10259
I1230 08:37:25.095244    1077 tlsconfig.go:240] Starting DynamicServingCertificateController
E1230 08:37:25.108803    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.Node: failed to list *v1.Node: Get "https://localhost:6443/api/v1/nodes?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:25.114284    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.Pod: failed to list *v1.Pod: Get "https://localhost:6443/api/v1/pods?fieldSelector=status.phase%21%3DSucceeded%2Cstatus.phase%21%3DFailed&limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:25.118183    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.PersistentVolumeClaim: failed to list *v1.PersistentVolumeClaim: Get "https://localhost:6443/api/v1/persistentvolumeclaims?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:25.121597    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.StorageClass: failed to list *v1.StorageClass: Get "https://localhost:6443/apis/storage.k8s.io/v1/storageclasses?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:25.126640    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.CSINode: failed to list *v1.CSINode: Get "https://localhost:6443/apis/storage.k8s.io/v1/csinodes?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:25.131111    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.ReplicationController: failed to list *v1.ReplicationController: Get "https://localhost:6443/api/v1/replicationcontrollers?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:25.135222    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.StatefulSet: failed to list *v1.StatefulSet: Get "https://localhost:6443/apis/apps/v1/statefulsets?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:25.139974    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.PersistentVolume: failed to list *v1.PersistentVolume: Get "https://localhost:6443/api/v1/persistentvolumes?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:25.144153    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1beta1.PodDisruptionBudget: failed to list *v1beta1.PodDisruptionBudget: Get "https://localhost:6443/apis/policy/v1beta1/poddisruptionbudgets?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:25.148220    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.Service: failed to list *v1.Service: Get "https://localhost:6443/api/v1/services?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:25.151946    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.ReplicaSet: failed to list *v1.ReplicaSet: Get "https://localhost:6443/apis/apps/v1/replicasets?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:25.940815    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.StorageClass: failed to list *v1.StorageClass: Get "https://localhost:6443/apis/storage.k8s.io/v1/storageclasses?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:25.970125    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1beta1.PodDisruptionBudget: failed to list *v1beta1.PodDisruptionBudget: Get "https://localhost:6443/apis/policy/v1beta1/poddisruptionbudgets?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:26.243434    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.ReplicaSet: failed to list *v1.ReplicaSet: Get "https://localhost:6443/apis/apps/v1/replicasets?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:26.347918    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.Node: failed to list *v1.Node: Get "https://localhost:6443/api/v1/nodes?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:26.351907    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.ReplicationController: failed to list *v1.ReplicationController: Get "https://localhost:6443/api/v1/replicationcontrollers?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:26.356410    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.CSINode: failed to list *v1.CSINode: Get "https://localhost:6443/apis/storage.k8s.io/v1/csinodes?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:26.403381    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.Pod: failed to list *v1.Pod: Get "https://localhost:6443/api/v1/pods?fieldSelector=status.phase%21%3DSucceeded%2Cstatus.phase%21%3DFailed&limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:26.487675    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.PersistentVolumeClaim: failed to list *v1.PersistentVolumeClaim: Get "https://localhost:6443/api/v1/persistentvolumeclaims?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:26.492118    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.Service: failed to list *v1.Service: Get "https://localhost:6443/api/v1/services?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:26.702162    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.StatefulSet: failed to list *v1.StatefulSet: Get "https://localhost:6443/apis/apps/v1/statefulsets?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:26.733535    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.PersistentVolume: failed to list *v1.PersistentVolume: Get "https://localhost:6443/api/v1/persistentvolumes?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:27.695348    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.StorageClass: failed to list *v1.StorageClass: Get "https://localhost:6443/apis/storage.k8s.io/v1/storageclasses?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:27.878905    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.ReplicaSet: failed to list *v1.ReplicaSet: Get "https://localhost:6443/apis/apps/v1/replicasets?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:28.030297    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.ReplicationController: failed to list *v1.ReplicationController: Get "https://localhost:6443/api/v1/replicationcontrollers?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:28.505897    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.Pod: failed to list *v1.Pod: Get "https://localhost:6443/api/v1/pods?fieldSelector=status.phase%21%3DSucceeded%2Cstatus.phase%21%3DFailed&limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:28.535531    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.CSINode: failed to list *v1.CSINode: Get "https://localhost:6443/apis/storage.k8s.io/v1/csinodes?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:28.651284    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.Node: failed to list *v1.Node: Get "https://localhost:6443/api/v1/nodes?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:28.836342    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1beta1.PodDisruptionBudget: failed to list *v1beta1.PodDisruptionBudget: Get "https://localhost:6443/apis/policy/v1beta1/poddisruptionbudgets?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:28.930991    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.Service: failed to list *v1.Service: Get "https://localhost:6443/api/v1/services?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:28.984235    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.StatefulSet: failed to list *v1.StatefulSet: Get "https://localhost:6443/apis/apps/v1/statefulsets?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:29.490950    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.PersistentVolume: failed to list *v1.PersistentVolume: Get "https://localhost:6443/api/v1/persistentvolumes?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:29.687897    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.PersistentVolumeClaim: failed to list *v1.PersistentVolumeClaim: Get "https://localhost:6443/api/v1/persistentvolumeclaims?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:32.137838    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.CSINode: failed to list *v1.CSINode: Get "https://localhost:6443/apis/storage.k8s.io/v1/csinodes?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:32.223412    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.ReplicationController: failed to list *v1.ReplicationController: Get "https://localhost:6443/api/v1/replicationcontrollers?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:32.353292    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.ReplicaSet: failed to list *v1.ReplicaSet: Get "https://localhost:6443/apis/apps/v1/replicasets?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:33.751824    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.Pod: failed to list *v1.Pod: Get "https://localhost:6443/api/v1/pods?fieldSelector=status.phase%21%3DSucceeded%2Cstatus.phase%21%3DFailed&limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:33.823744    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.StorageClass: failed to list *v1.StorageClass: Get "https://localhost:6443/apis/storage.k8s.io/v1/storageclasses?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:33.840737    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.Service: failed to list *v1.Service: Get "https://localhost:6443/api/v1/services?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:34.468467    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.StatefulSet: failed to list *v1.StatefulSet: Get "https://localhost:6443/apis/apps/v1/statefulsets?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:34.740967    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.PersistentVolume: failed to list *v1.PersistentVolume: Get "https://localhost:6443/api/v1/persistentvolumes?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:34.742091    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.PersistentVolumeClaim: failed to list *v1.PersistentVolumeClaim: Get "https://localhost:6443/api/v1/persistentvolumeclaims?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:35.057780    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.Node: failed to list *v1.Node: Get "https://localhost:6443/api/v1/nodes?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
E1230 08:37:35.124524    1077 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1beta1.PodDisruptionBudget: failed to list *v1beta1.PodDisruptionBudget: Get "https://localhost:6443/apis/policy/v1beta1/poddisruptionbudgets?limit=500&resourceVersion=0": x509: certificate is valid for kubernetes, kubernetes.default, kubernetes.default.svc, kubernetes.default.svc.cluster, kubernetes.svc.cluster.local, not localhost
```

So, the certificate is not valid for `localhost` hostname / domain name. So, I
gotta use the IP address maybe, I guess. 😅

```bash
$ kube-scheduler --master https://192.168.64.39:6443
I1230 08:40:08.040772    1094 serving.go:331] Generated self-signed cert in-memory
W1230 08:40:08.872340    1094 authentication.go:303] No authentication-kubeconfig provided in order to lookup client-ca-file in configmap/extension-apiserver-authentication in kube-system, so client certificate authentication won't work.
W1230 08:40:08.872536    1094 authentication.go:327] No authentication-kubeconfig provided in order to lookup requestheader-client-ca-file in configmap/extension-apiserver-authentication in kube-system, so request-header client certificate authentication won't work.
W1230 08:40:08.872903    1094 authorization.go:173] No authorization-kubeconfig provided, so SubjectAccessReview of authorization tokens won't work.
W1230 08:40:08.919123    1094 authorization.go:47] Authorization is disabled
W1230 08:40:08.919628    1094 authentication.go:40] Authentication is disabled
I1230 08:40:08.919971    1094 deprecated_insecure_serving.go:51] Serving healthz insecurely on [::]:10251
I1230 08:40:08.922119    1094 secure_serving.go:197] Serving securely on [::]:10259
I1230 08:40:08.922641    1094 tlsconfig.go:240] Starting DynamicServingCertificateController
E1230 08:40:08.929468    1094 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.StatefulSet: failed to list *v1.StatefulSet: Get "https://192.168.64.39:6443/apis/apps/v1/statefulsets?limit=500&resourceVersion=0": x509: certificate signed by unknown authority
E1230 08:40:08.935335    1094 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.ReplicaSet: failed to list *v1.ReplicaSet: Get "https://192.168.64.39:6443/apis/apps/v1/replicasets?limit=500&resourceVersion=0": x509: certificate signed by unknown authority
E1230 08:40:08.939144    1094 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.Pod: failed to list *v1.Pod: Get "https://192.168.64.39:6443/api/v1/pods?fieldSelector=status.phase%21%3DSucceeded%2Cstatus.phase%21%3DFailed&limit=500&resourceVersion=0": x509: certificate signed by unknown authority
E1230 08:40:08.942787    1094 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.StorageClass: failed to list *v1.StorageClass: Get "https://192.168.64.39:6443/apis/storage.k8s.io/v1/storageclasses?limit=500&resourceVersion=0": x509: certificate signed by unknown authority
E1230 08:40:08.946002    1094 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.CSINode: failed to list *v1.CSINode: Get "https://192.168.64.39:6443/apis/storage.k8s.io/v1/csinodes?limit=500&resourceVersion=0": x509: certificate signed by unknown authority
E1230 08:40:08.950535    1094 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1beta1.PodDisruptionBudget: failed to list *v1beta1.PodDisruptionBudget: Get "https://192.168.64.39:6443/apis/policy/v1beta1/poddisruptionbudgets?limit=500&resourceVersion=0": x509: certificate signed by unknown authority
E1230 08:40:08.954393    1094 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.Node: failed to list *v1.Node: Get "https://192.168.64.39:6443/api/v1/nodes?limit=500&resourceVersion=0": x509: certificate signed by unknown authority
E1230 08:40:08.958179    1094 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.Service: failed to list *v1.Service: Get "https://192.168.64.39:6443/api/v1/services?limit=500&resourceVersion=0": x509: certificate signed by unknown authority
E1230 08:40:08.961257    1094 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.PersistentVolume: failed to list *v1.PersistentVolume: Get "https://192.168.64.39:6443/api/v1/persistentvolumes?limit=500&resourceVersion=0": x509: certificate signed by unknown authority
E1230 08:40:08.965322    1094 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.PersistentVolumeClaim: failed to list *v1.PersistentVolumeClaim: Get "https://192.168.64.39:6443/api/v1/persistentvolumeclaims?limit=500&resourceVersion=0": x509: certificate signed by unknown authority
E1230 08:40:08.969225    1094 reflector.go:138] k8s.io/client-go/informers/factory.go:134: Failed to watch *v1.ReplicationController: failed to list *v1.ReplicationController: Get "https://192.168.64.39:6443/api/v1/replicationcontrollers?limit=500&resourceVersion=0": x509: certificate signed by unknown authority
```

Of course it's an unknown Certificate Authority. Let's try to fix that first. :)

---

I'm going to check how kubernetes-the-hard-way does it. Btw, I also noticed that
it currently uses 1.18.6 actually

https://github.com/kelseyhightower/kubernetes-the-hard-way/commit/ca96371e4d2d2176e8b2c3f5b656b5d92973479e

Good thing :) ;) :D

Okay, so, kube-config is deprecated. So, what I'm going to use is Authentication
Kube Config and see if it works out :)

Kubernetes The Hard Way uses kubernetes config yaml files for configuration.
That's okay, I'm just going to use flags and also, it uses the deprecated
kube config I think. At least that's how it looks like to me. Not sure.

Now, I need to create some certificates. I was wondering if I should really
put the `system:kube-scheduler` in the username or in the group as part of
organization. Hmm.

According to the latest docs

https://kubernetes.io/docs/setup/best-practices/certificates/#configure-certificates-for-user-accounts

Nothing much really. Hmm. Anyways, now, let's get to work and create two
certificates. One with the group, and another without :)

```bash
$ cat > normal-kube-scheduler-csr.json <<EOF
{
    "CN": "system:kube-scheduler",
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "names": [
        {
            "C": "US",
            "L": "San Francisco",
            "O": "Kubernetes",
            "OU": "Kubernetes Scheduler",
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
  normal-kube-scheduler-csr.json | cfssljson -bare normal-kube-scheduler
```

```bash
$ cat > kube-scheduler-csr.json <<EOF
{
    "CN": "system:kube-scheduler",
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "names": [
        {
            "C": "US",
            "L": "San Francisco",
            "O": "system:kube-scheduler",
            "OU": "Kubernetes Scheduler",
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
  kube-scheduler-csr.json | cfssljson -bare kube-scheduler
```

```bash
$ export API_SERVER_NODE_IP=192.168.64.39

$ kubectl config set-cluster my-own-k8s-cluster \
    --certificate-authority=ca.pem \
    --embed-certs=true \
    --server=https://${API_SERVER_NODE_IP}:6443 \
    --kubeconfig=kube-scheduler.kubeconfig

$ kubectl config set-credentials normal-kube-scheduler \
    --client-certificate=normal-kube-scheduler.pem \
    --client-key=normal-kube-scheduler-key.pem \
    --kubeconfig=kube-scheduler.kubeconfig

$ kubectl config set-context normal-kube-scheduler \
    --cluster=my-own-k8s-cluster \
    --user=normal-kube-scheduler \
    --kubeconfig=kube-scheduler.kubeconfig

$ kubectl config set-credentials kube-scheduler \
    --client-certificate=kube-scheduler.pem \
    --client-key=kube-scheduler-key.pem \
    --kubeconfig=kube-scheduler.kubeconfig

$ kubectl config set-context kube-scheduler \
    --cluster=my-own-k8s-cluster \
    --user=kube-scheduler \
    --kubeconfig=kube-scheduler.kubeconfig

$ kubectl config use-context normal-kube-scheduler \
    --kubeconfig=kube-scheduler.kubeconfig
```
