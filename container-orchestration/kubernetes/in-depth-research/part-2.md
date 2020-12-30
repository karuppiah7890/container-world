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
