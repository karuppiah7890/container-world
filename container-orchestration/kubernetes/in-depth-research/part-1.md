# Part 1

To start off, I want to manually do everything that's needed to setup a cluster.
As manually as possible. This will help me understand some of the in-depth
stuff, more like the nitty gritties.

I'm planning to go from top to bottom approach. So, I'll look at the high level
architecture or components and then dig in deep.

I also plan to use the latest released stable version of Kubernetes. It's
v1.20.1 as of this writing.

https://kubernetes.io/docs/concepts/overview/what-is-kubernetes/

https://kubernetes.io/docs/concepts/overview/components/

Looking at the high level components, I can see that there's master node
components (control plane components) and then worker node components.

I plan to run all of them in one single machine. :)

Now, after reading https://kubernetes.io/docs/concepts/overview/components/ , I
have a basic idea of what stuff I need to run. Maybe I could start with the
api server? I'm not sure though. But it's okay. I'm not going to go and check
what kubernetes the hard way tutorial did. I wanna try things on my own and see
what happens. Given this is a top down approach, I can start with just the
components I have just seen.

I feel that api server and etcd is a good place to start. Let's see :)

https://github.com/kubernetes/kubernetes/releases

https://groups.google.com/g/kubernetes-announce

https://groups.google.com/g/kubernetes-announce/c/qdt2OTuuFsc v1.20.1

I'm going to get v1.20.1

Let me start by getting the api server and also the client (kubectl) binaries.

I plan to use multipass for running VMs

https://github.com/canonical/multipass

It's cross platform, so, anyone can try it out :) My other alternative option
was vagrant with virtual box. Anyways :)

```bash
$ multipass launch --name my-own-k8s-cluster
Creating my-own-k8s-cluster \
Retrieving image: 2%
Retrieving kernel image:  /
Retrieving initrd image:  /
launch failed: The following errors occurred:
Instance stopped while starting

$ multipass start my-own-k8s-cluster
start failed: The following errors occurred:
Instance stopped while starting
```

According to multipass help, the machine has all defaults - CPU, RAM and disk.
CPU - 1. RAM - 1G, disk - 5GB

Weird. The whole thing doesn't work

```bash
$ multipass ls
Name                    State             IPv4             Image
my-own-k8s-cluster      Stopped           --               Ubuntu 20.04 LTS

$ multipass delete my-own-k8s-cluster

$ multipass purge

$ multipass launch --name my-own-k8s-cluster
launch failed: The following errors occurred:
Instance stopped while starting
```

I'm planning to get the latest version of multipass before getting started :)

```bash
$ multipass version
multipass  1.1.0+mac
multipassd 1.1.0+mac
```

https://github.com/canonical/multipass/releases/tag/v1.5.0

https://github.com/canonical/multipass/releases/download/v1.5.0/multipass-1.5.0+mac-Darwin.pkg

```bash
$ multipass version
multipass  1.5.0+mac
multipassd 1.5.0+mac
```

Now it all works! :)

```bash
$ multipass launch --name my-own-k8s-cluster
Launched: my-own-k8s-cluster
```

```bash
$ multipass exec my-own-k8s-cluster bash
To run a command as administrator (user "root"), use "sudo <command>".
See "man sudo_root" for details.

ubuntu@my-own-k8s-cluster:~$
```

Now I'm getting the binaries

https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.20.md

https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.20.md#downloads-for-v1201

https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.20.md#client-binaries

https://dl.k8s.io/v1.20.1/kubernetes-client-linux-amd64.tar.gz

https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.20.md#server-binaries

https://dl.k8s.io/v1.20.1/kubernetes-server-linux-amd64.tar.gz

https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.20.md#node-binaries-1

https://dl.k8s.io/v1.20.0/kubernetes-node-linux-amd64.tar.gz

I also got links for the node binary. I'll get them all ;) :D

**STEP 1**

```bash
$ wget https://dl.k8s.io/v1.20.1/kubernetes-client-linux-amd64.tar.gz https://dl.k8s.io/v1.20.1/kubernetes-server-linux-amd64.tar.gz https://dl.k8s.io/v1.20.0/kubernetes-node-linux-amd64.tar.gz

$ ls -lh
total 408M
-rw-rw-r-- 1 ubuntu ubuntu  12M Dec 18 20:52 kubernetes-client-linux-amd64.tar.gz
-rw-rw-r-- 1 ubuntu ubuntu  95M Dec  9 00:59 kubernetes-node-linux-amd64.tar.gz
-rw-rw-r-- 1 ubuntu ubuntu 302M Dec 18 20:52 kubernetes-server-linux-amd64.tar.gz
```

Now, let's get started with the API server :)

**STEP 2**

```bash
$ tar -xvzf kubernetes-server-linux-amd64.tar.gz
$ tar -xvzf kubernetes-client-linux-amd64.tar.gz
$ tar -xvzf kubernetes-node-linux-amd64.tar.gz
```

All the contents are put into a single directory called `kubernetes` with
multiple LICENSE files and three important sub directories based on the tar
balls - `server`, `client`, `node`

Let's start with `server` :)

```bash
$ ls kubernetes/server/bin/
$ cat kubernetes/server/bin/kube-apiserver.docker_tag
v1.20.1

$ ls kubernetes/server/bin/kube-apiserver.tar
kubernetes/server/bin/kube-apiserver.tar

$ cd kubernetes/server/bin/

$ tar -xvf kube-apiserver.tar
75c7f711208082c548b935ab31e681ea30acccdce6b7abeecabae5bbfd326627.json
8ed3da63de1c24a56a0a57e010f19ac8433d9785fab3a9f878ffc3e61d6474aa/
8ed3da63de1c24a56a0a57e010f19ac8433d9785fab3a9f878ffc3e61d6474aa/VERSION
8ed3da63de1c24a56a0a57e010f19ac8433d9785fab3a9f878ffc3e61d6474aa/json
8ed3da63de1c24a56a0a57e010f19ac8433d9785fab3a9f878ffc3e61d6474aa/layer.tar
97253bc52fe39e9e1d68f79a3e936e039a2dfb384cebbaf11021119a15072c13/
97253bc52fe39e9e1d68f79a3e936e039a2dfb384cebbaf11021119a15072c13/VERSION
...
```

OOPS. I think that was the tar ball format of the docker image. Not sure though.
Hmm. Anyways, the `kube-apiserver` is present as a proper binary itself. Instead
I assumed it's inside the tar ball. 🤦

```bash
$ ./kube-apiserver
W1225 14:43:27.624208    3001 services.go:37] No CIDR for service cluster IPs specified. Default value which was 10.0.0.0/24 is deprecated and will be removed in future releases. Please specify it using --service-cluster-ip-range on kube-apiserver.
Error: error creating self-signed certificates: mkdir /var/run/kubernetes: permission denied

$ whoami
ubuntu

$ ls -l /var
total 44
drwxr-xr-x  2 root root   4096 Apr 15  2020 backups
drwxr-xr-x 12 root root   4096 Dec 25 14:30 cache
drwxrwxrwt  2 root root   4096 Dec 10 19:03 crash
drwxr-xr-x 38 root root   4096 Dec 25 14:30 lib
drwxrwsr-x  2 root staff  4096 Apr 15  2020 local
lrwxrwxrwx  1 root root      9 Dec 10 19:00 lock -> /run/lock
drwxrwxr-x  8 root syslog 4096 Dec 25 14:30 log
drwxrwsr-x  2 root mail   4096 Dec 10 19:00 mail
drwxr-xr-x  2 root root   4096 Dec 10 19:00 opt
lrwxrwxrwx  1 root root      4 Dec 10 19:00 run -> /run
drwxr-xr-x  5 root root   4096 Dec 10 19:04 snap
drwxr-xr-x  4 root root   4096 Dec 10 19:00 spool
drwxrwxrwt  5 root root   4096 Dec 25 14:30 tmp

$ ls -l /var | grep run
lrwxrwxrwx  1 root root      9 Dec 10 19:00 lock -> /run/lock
lrwxrwxrwx  1 root root      4 Dec 10 19:00 run -> /run
```

Only the `root` can create directories in `/var/run`, hmm.

**STEP**

```bash
$ sudo mkdir -p /var/run/kubernetes

$ ls -l /var/run/ | grep kubernetes
drwxr-xr-x  2 root root   40 Dec 25 14:45 kubernetes

$ sudo chown ubuntu -R /var/run/kubernetes

$ ls -l /var/run/ | grep kuber
drwxr-xr-x  2 ubuntu root   40 Dec 25 14:45 kubernetes
```

```
$ ./kube-apiserver
W1225 14:47:23.676350    3090 services.go:37] No CIDR for service cluster IPs specified. Default value which was 10.0.0.0/24 is deprecated and will be removed in future releases. Please specify it using --service-cluster-ip-range on kube-apiserver.
I1225 14:47:24.589921    3090 serving.go:325] Generated self-signed cert (/var/run/kubernetes/apiserver.crt, /var/run/kubernetes/apiserver.key)
I1225 14:47:24.591055    3090 server.go:632] external host was not specified, using 192.168.64.39
W1225 14:47:24.592177    3090 authentication.go:519] AnonymousAuth is not allowed with the AlwaysAllow authorizer. Resetting AnonymousAuth to false. You should use a different authorizer
Error: [--etcd-servers must be specified, service-account-issuer is a required flag, --service-account-signing-key-file and --service-account-issuer are required flags]
```

So, it's creating a self signed certificate. Hmm. And it's also using a default
value for the range of service's cluster IPs.

There's some log about Anonymous Auth. There's some sort of Authorizer for the
API server, hmm. Makes sense.

There are some required flags. Hmm.

`--etcd-servers` for ETCD Server URLs I guess

`--service-account-issuer` - Not sure about this. Gotta check what service
account this is. Kubernetes, or if it means cloud service provider service
accounts.

There's also `--service-account-signing-key-file`, hmm.

Let's start with etcd maybe :)

https://etcd.io/

https://etcd.io/docs/v3.4.0/

I plan to run just one instance of etcd :)

https://etcd.io/docs/v3.4.0/demo/#set-up-a-cluster

https://etcd.io/docs/v3.4.0/dl-build/

https://github.com/etcd-io/etcd/releases/

https://github.com/etcd-io/etcd/releases/tag/v3.4.14

http://play.etcd.io/install

```bash
ETCD_VER=v3.4.14

# choose either URL
GOOGLE_URL=https://storage.googleapis.com/etcd
GITHUB_URL=https://github.com/etcd-io/etcd/releases/download
DOWNLOAD_URL=${GOOGLE_URL}

rm -f /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz
rm -rf /tmp/etcd-download-test && mkdir -p /tmp/etcd-download-test

curl -L ${DOWNLOAD_URL}/${ETCD_VER}/etcd-${ETCD_VER}-linux-amd64.tar.gz -o /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz
tar xzvf /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz -C /tmp/etcd-download-test --strip-components=1
rm -f /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz

/tmp/etcd-download-test/etcd --version
/tmp/etcd-download-test/etcdctl version
```

```bash
$ /tmp/etcd-download-test/etcd --version
etcd Version: 3.4.14
Git SHA: 8a03d2e96
Go Version: go1.12.17
Go OS/Arch: linux/amd64

$ /tmp/etcd-download-test/etcdctl version
etcdctl version: 3.4.14
API version: 3.4
```

**STEP**

```bash
# start a local etcd server
$ /tmp/etcd-download-test/etcd
...
```

```bash
$ # write,read to etcd
$ /tmp/etcd-download-test/etcdctl --endpoints=localhost:2379 put foo bar
OK

$ /tmp/etcd-download-test/etcdctl --endpoints=localhost:2379 get foo
foo
bar
```

Now that the etcd is running, I can run the kubernetes API server

```bash
$ ./kube-apiserver --help | grep etcd
      --encryption-provider-config string        The file containing configuration for encryption providers to be used for storing secrets in etcd
      --etcd-cafile string                       SSL Certificate Authority file used to secure etcd communication.
      --etcd-certfile string                     SSL certification file used to secure etcd communication.
      --etcd-compaction-interval duration        The interval of compaction requests. If 0, the compaction request from apiserver is disabled. (default 5m0s)
      --etcd-count-metric-poll-period duration   Frequency of polling etcd for number of resources per type. 0 disables the metric collection. (default 1m0s)
      --etcd-db-metric-poll-interval duration    The interval of requests to poll etcd and update metric. 0 disables the metric collection (default 30s)
      --etcd-healthcheck-timeout duration        The timeout to use when checking etcd health. (default 2s)
      --etcd-keyfile string                      SSL key file used to secure etcd communication.
      --etcd-prefix string                       The prefix to prepend to all resource paths in etcd. (default "/registry")
      --etcd-servers strings                     List of etcd servers to connect with (scheme://ip:port), comma separated.
      --etcd-servers-overrides strings           Per-resource etcd servers overrides, comma separated. The individual override format: group/resource#servers, where servers are URLs, semicolon separated.
      --storage-backend string                   The storage backend for persistence. Options: 'etcd3' (default).
      --service-account-lookup                            If true, validate ServiceAccount tokens exist in etcd as part of authentication. (default true)
```

```bash
$ ./kube-apiserver --etcd-servers localhost:2379
W1225 15:32:01.969743    3286 services.go:37] No CIDR for service cluster IPs specified. Default value which was 10.0.0.0/24 is deprecated and will be removed in future releases. Please specify it using --service-cluster-ip-range on kube-apiserver.
I1225 15:32:01.970354    3286 server.go:632] external host was not specified, using 192.168.64.39
W1225 15:32:01.970754    3286 authentication.go:519] AnonymousAuth is not allowed with the AlwaysAllow authorizer. Resetting AnonymousAuth to false. You should use a different authorizer
Error: [service-account-issuer is a required flag, --service-account-signing-key-file and --service-account-issuer are required flags]
```

Okay, next is the service account issuer thing.

```bash
$ ./kube-apiserver --help | grep service-account
      --api-audiences strings                             Identifiers of the API. The service account token authenticator will validate that tokens used against the API are bound to at least one of these audiences. If the --service-account-issuer flag is configured and this flag is not, this field defaults to a single element list containing the issuer URL.
      --service-account-extend-token-expiration           Turns on projected service account expiration extension during token generation, which helps safe transition from legacy token to bound service account token feature. If this flag is enabled, admission injected tokens would be extended up to 1 year to prevent unexpected failure during transition, ignoring value of service-account-max-token-expiration. (default true)
      --service-account-issuer string                     Identifier of the service account token issuer. The issuer will assert this identifier in "iss" claim of issued tokens. This value is a string or URI. If this option is not a valid URI per the OpenID Discovery 1.0 spec, the ServiceAccountIssuerDiscovery feature will remain disabled, even if the feature gate is set to true. It is highly recommended that this value comply with the OpenID spec: https://openid.net/specs/openid-connect-discovery-1_0.html. In practice, this means that service-account-issuer must be an https URL. It is also highly recommended that this URL be capable of serving OpenID discovery documents at {service-account-issuer}/.well-known/openid-configuration.
      --service-account-jwks-uri string                   Overrides the URI for the JSON Web Key Set in the discovery doc served at /.well-known/openid-configuration. This flag is useful if the discovery docand key set are served to relying parties from a URL other than the API server's external (as auto-detected or overridden with external-hostname). Only valid if the ServiceAccountIssuerDiscovery feature gate is enabled.
      --service-account-key-file stringArray              File containing PEM-encoded x509 RSA or ECDSA private or public keys, used to verify ServiceAccount tokens. The specified file can contain multiple keys, and the flag can be specified multiple times with different files. If unspecified, --tls-private-key-file is used. Must be specified when --service-account-signing-key is provided
      --service-account-lookup                            If true, validate ServiceAccount tokens exist in etcd as part of authentication. (default true)
      --service-account-max-token-expiration duration     The maximum validity duration of a token created by the service account token issuer. If an otherwise valid TokenRequest with a validity duration larger than this value is requested, a token will be issued with a validity duration of this value.
      --service-account-signing-key-file string     Path to the file that contains the current private key of the service account token issuer. The issuer will sign issued ID tokens with this private key.
```

Seems like something I might not need. I need to check how to do authentication
actually. Hmm. As this one talks about OpenID and stuff and asks me to mention
service account issuer.

I'll go and check what are the authentication mechanisms available for the
api server :)

I also need to check authorization

https://kubernetes.io/docs/reference/access-authn-authz/authentication/

https://kubernetes.io/docs/reference/access-authn-authz/authorization/

Looks like I can do authentication with tokens or using certificates as there
is no concept of users in kubernetes, only service accounts..

---

On the side I started checking about `service-account-issuer` flag and found
some stuff!

https://duckduckgo.com/?t=ffab&q=service-account-issuer+&ia=web

https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/

https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/#service-account-issuer-discovery

`ServiceAccountIssuerDiscovery`

Looks like it's for a feature I don't even know about, and I don't think I need
it now for a simple and basic setup.

I'm trying to turn it off as it looks like it's turned on by default, though it
says "beta" but I guess it's beta only starting from 1.20.x

https://kubernetes.io/docs/reference/command-line-tools-reference/feature-gates/

```bash
$ ./kube-apiserver -h | grep feature-gate
      --feature-gates mapStringBool                  A set of key=value pairs that describe feature gates for alpha/experimental features. Options are:

$ ./kube-apiserver -h | grep ServiceAccountIssuerDiscovery
                                                     ServiceAccountIssuerDiscovery=true|false (BETA - default=true)
      --service-account-issuer string                     Identifier of the service account token issuer. The issuer will assert this identifier in "iss" claim of issued tokens. This value is a string or URI. If this option is not a valid URI per the OpenID Discovery 1.0 spec, the ServiceAccountIssuerDiscovery feature will remain disabled, even if the feature gate is set to true. It is highly recommended that this value comply with the OpenID spec: https://openid.net/specs/openid-connect-discovery-1_0.html. In practice, this means that service-account-issuer must be an https URL. It is also highly recommended that this URL be capable of serving OpenID discovery documents at {service-account-issuer}/.well-known/openid-configuration.
      --service-account-jwks-uri string                   Overrides the URI for the JSON Web Key Set in the discovery doc served at /.well-known/openid-configuration. This flag is useful if the discovery docand key set are served to relying parties from a URL other than the API server's external (as auto-detected or overridden with external-hostname). Only valid if the ServiceAccountIssuerDiscovery feature gate is enabled.
```

I just found out how to provide the feature gate as before this I was doing it
wrong, I think :P

https://duckduckgo.com/?t=ffab&q=kubernetes+feature+gate&ia=web

https://kubernetes.io/docs/reference/command-line-tools-reference/feature-gates/#overview

Let me disable it using the `feature-gates` flag then

```bash
$ ./kube-apiserver --etcd-servers localhost:2379 --feature-gates "ServiceAccountIssuerDiscovery=false"
W1226 17:40:48.785314   20871 services.go:37] No CIDR for service cluster IPs specified. Default value which was 10.0.0.0/24 is deprecated and will be removed in future releases. Please specify it using --service-cluster-ip-range on kube-apiserver.
I1226 17:40:48.786069   20871 server.go:632] external host was not specified, using 192.168.64.39
W1226 17:40:48.786184   20871 authentication.go:519] AnonymousAuth is not allowed with the AlwaysAllow authorizer. Resetting AnonymousAuth to false. You should use a different authorizer
Error: [service-account-issuer is a required flag, --service-account-signing-key-file and --service-account-issuer are required flags]
```

It still doesn't work though. It's still looking for the service account issuer
flag. Hmm.

Looks like I also need to disable another feature

https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/#service-account-token-volume-projection

It's a stable feature actually. Hmm. It's called `TokenRequestProjection`

```bash
$ ./kube-apiserver --etcd-servers localhost:2379 --feature-gates "TokenRequestProjection=false"
Error: invalid argument "TokenRequestProjection=false" for "--feature-gates" flag: cannot set feature gate TokenRequestProjection to false, feature is locked to true
```

Looks like I can't disable it even if I want to

https://kubernetes.io/docs/reference/command-line-tools-reference/feature-gates/#using-a-feature

---

https://duckduckgo.com/?t=ffab&q=service-account-issuer+is+a+required+flag%2C&ia=web

https://github.com/kubernetes/kubernetes/blob/master/pkg/kubeapiserver/options/authentication.go

https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.20.md

https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.20.md#urgent-upgrade-notes

https://openid.net/specs/openid-connect-discovery-1_0.html

https://duckduckgo.com/?t=ffab&q=Service+Account+Token+Volume+Projection+with+vault&ia=web

https://www.alibabacloud.com/help/doc-detail/160384.htm

Looks like I can't get away without specifying service account issuer and stuff.
It's required. So, it won't work in any other case is how it sounds like. Hmm.

```bash
$ kube-apiserver --help | grep "service-account-issuer"
      --api-audiences strings                             Identifiers of the API. The service account token authenticator will validate that tokens used against the API are bound to at least one of these audiences. If the --service-account-issuer flag is configured and this flag is not, this field defaults to a single element list containing the issuer URL.
      --service-account-issuer string                     Identifier of the service account token issuer. The issuer will assert this identifier in "iss" claim of issued tokens. This value is a string or URI. If this option is not a valid URI per the OpenID Discovery 1.0 spec, the ServiceAccountIssuerDiscovery feature will remain disabled, even if the feature gate is set to true. It is highly recommended that this value comply with the OpenID spec: https://openid.net/specs/openid-connect-discovery-1_0.html. In practice, this means that service-account-issuer must be an https URL. It is also highly recommended that this URL be capable of serving OpenID discovery documents at {service-account-issuer}/.well-known/openid-configuration.
```

```bash
$ kube-apiserver --help | grep "service-account-"
      --api-audiences strings                             Identifiers of the API. The service account token authenticator will validate that tokens used against the API are bound to at least one of these audiences. If the --service-account-issuer flag is configured and this flag is not, this field defaults to a single element list containing the issuer URL.
      --service-account-extend-token-expiration           Turns on projected service account expiration extension during token generation, which helps safe transition from legacy token to bound service account token feature. If this flag is enabled, admission injected tokens would be extended up to 1 year to prevent unexpected failure during transition, ignoring value of service-account-max-token-expiration. (default true)
      --service-account-issuer string                     Identifier of the service account token issuer. The issuer will assert this identifier in "iss" claim of issued tokens. This value is a string or URI. If this option is not a valid URI per the OpenID Discovery 1.0 spec, the ServiceAccountIssuerDiscovery feature will remain disabled, even if the feature gate is set to true. It is highly recommended that this value comply with the OpenID spec: https://openid.net/specs/openid-connect-discovery-1_0.html. In practice, this means that service-account-issuer must be an https URL. It is also highly recommended that this URL be capable of serving OpenID discovery documents at {service-account-issuer}/.well-known/openid-configuration.
      --service-account-jwks-uri string                   Overrides the URI for the JSON Web Key Set in the discovery doc served at /.well-known/openid-configuration. This flag is useful if the discovery docand key set are served to relying parties from a URL other than the API server's external (as auto-detected or overridden with external-hostname). Only valid if the ServiceAccountIssuerDiscovery feature gate is enabled.
      --service-account-key-file stringArray              File containing PEM-encoded x509 RSA or ECDSA private or public keys, used to verify ServiceAccount tokens. The specified file can contain multiple keys, and the flag can be specified multiple times with different files. If unspecified, --tls-private-key-file is used. Must be specified when --service-account-signing-key is provided
      --service-account-lookup                            If true, validate ServiceAccount tokens exist in etcd as part of authentication. (default true)
      --service-account-max-token-expiration duration     The maximum validity duration of a token created by the service account token issuer. If an otherwise valid TokenRequest with a validity duration larger than this value is requested, a token will be issued with a validity duration of this value.
      --service-account-signing-key-file string     Path to the file that contains the current private key of the service account token issuer. The issuer will sign issued ID tokens with this private key.
```

Looking at https://www.alibabacloud.com/help/doc-detail/160384.htm Looks like
we can just provide

`kubernetes.default.svc` as the value for `--service-account-issuer` and for
`--api-audiences` too. `--service-account-key-file` and
`--service-account-signing-key-file` is still something we have to create but I
remember that we can reuse an existing private key that is used for something
else. I forgot what. Probably the API server's private key for the https
endpoint. There is also a mention of `--tls-private-key-file` among this. So,
I guess it makes sense.

---

Another problem I'm tackling is some machine setup issue. Damn multipass thing.
It does some magic in between that I don't know about. What happens is - it
removes the `/var/run/kubernetes` directory that I create with access to the
`ubuntu` user. So I had to create `/opt/kubernetes` as that wasn't getting
deleted. The deletion was happening on machine restart. I noticed only now.

Anyways, I'm thinking about creating a Certificate Authority (CA) certificate
first. A root CA cert. And then create / issue certificates with the CA.

For this I plan to use `cfssl` and `cfssljson` which many tend to use. I want
to see how to use it and how it helps :) Though traditionally people use openssl
and similar tools

It reminds me of the openssh alternative - https://github.com/FiloSottile/age
or something similar, for crypto stuff :)

https://github.com/cloudflare/cfssl/releases/download/v1.5.0/cfssl_1.5.0_linux_amd64

https://github.com/cloudflare/cfssl/releases/download/v1.5.0/cfssljson_1.5.0_linux_amd64

```bash
$ sudo apt install net-tools
$ $ ifconfig
enp0s2: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
        inet 192.168.64.39  netmask 255.255.255.0  broadcast 192.168.64.255
        inet6 fe80::70f0:8cff:fe5c:64c3  prefixlen 64  scopeid 0x20<link>
        ether 72:f0:8c:5c:64:c3  txqueuelen 1000  (Ethernet)
        RX packets 884  bytes 271578 (271.5 KB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 612  bytes 78411 (78.4 KB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

lo: flags=73<UP,LOOPBACK,RUNNING>  mtu 65536
        inet 127.0.0.1  netmask 255.0.0.0
        inet6 ::1  prefixlen 128  scopeid 0x10<host>
        loop  txqueuelen 1000  (Local Loopback)
        RX packets 90  bytes 6884 (6.8 KB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 90  bytes 6884 (6.8 KB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0
```

```bash
$ sudo curl https://github.com/cloudflare/cfssl/releases/download/v1.5.0/cfssl_1.5.0_linux_amd64 -o /usr/local/bin/cfssl
$ sudo chmod +x /usr/local/bin/cfssl

$ sudo curl https://github.com/cloudflare/cfssl/releases/download/v1.5.0/cfssljson_1.5.0_linux_amd64 -o /usr/local/bin/cfssljson
$ sudo chmod +x /usr/local/bin/cfssljson
```

Okay, I'm planning to refer the kubernetes hard way method and see how it all
works :P :P XP

https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/04-certificate-authority.md

Right, so I found a few things. I mean, it had a lot of information. But for now
I just found out how to do some of the stuff. I also was thinking about how to
authenticate myself to access the cluster, using either tokens which I put in a
token file and provide it to api server or using certificate - client
certificate and how I will have some sort of admin permission of sorts or
at least some basic permission. I just found out that the admin user client
certificate is created with the "system:masters" organization which will be
read by kubernetes as the user group. And if I belong to that group I think it
will consider me as a master / administrator. Just a guess. Gotta check if it
is still valid in v1.20.1 :)

I was also seeing about service account stuff in api server

https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/08-bootstrapping-kubernetes-controllers.md

**STEP**

```bash
$ cat > ca-csr.json <<EOF
{
    "CN": "Kubernetes",
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "names": [
        {
            "C": "US",
            "L": "San Francisco",
            "O": "Kubernetes",
            "OU": "CA",
            "ST": "California"
        }
    ]
}
EOF

$ cfssl gencert -initca ca-csr.json | cfssljson -bare ca
```

I kept getting an error

```bash
$ cfssl gencert -initca ca-csr.json | cfssljson -bare ca
/usr/local/bin/cfssljson: line 1: syntax error near unexpected token `<'
/usr/local/bin/cfssljson: line 1: `<html><body>You are being <a href="https://github-production-release-asset-2e65be.s3.amazonaws.com/21591001/2513db00-148b-11eb-82c0-e13cb9b9405f?X-Amz-Algorithm=AWS4-HMAC-SHA256&amp;X-Amz-Credential=AKIAIWNJYAX4CSVEH53A%2F20201227%2Fus-east-1%2Fs3%2Faws4_request&amp;X-Amz-Date=20201227T074138Z&amp;X-Amz-Expires=300&amp;X-Amz-Signature=8cb9de2d0b2ef0888c875c062965c4e008b4ea2a1b356738c9e5f726fc05d29f&amp;X-Amz-SignedHeaders=host&amp;actor_id=0&amp;key_id=0&amp;repo_id=21591001&amp;response-content-disposition=attachment%3B%20filename%3Dcfssljson_1.5.0_linux_amd64&amp;response-content-type=application%2Foctet-stream">redirected</a>.</body></html>'
/usr/local/bin/cfssl: line 1: syntax error near unexpected token `<'
/usr/local/bin/cfssl: line 1: `<html><body>You are being <a href="https://github-production-release-asset-2e65be.s3.amazonaws.com/21591001/b5055500-148a-11eb-9528-44972336b695?X-Amz-Algorithm=AWS4-HMAC-SHA256&amp;X-Amz-Credential=AKIAIWNJYAX4CSVEH53A%2F20201227%2Fus-east-1%2Fs3%2Faws4_request&amp;X-Amz-Date=20201227T073955Z&amp;X-Amz-Expires=300&amp;X-Amz-Signature=2fe2795e562fea205b61f75af5e46ea020733dbf2ef8c09d1a5fc7f1ab4d5bf2&amp;X-Amz-SignedHeaders=host&amp;actor_id=0&amp;key_id=0&amp;repo_id=21591001&amp;response-content-disposition=attachment%3B%20filename%3Dcfssl_1.5.0_linux_amd64&amp;response-content-type=application%2Foctet-stream">redirected</a>.</body></html>'

$ cfssl gencert -initca ca-csr.json
/usr/local/bin/cfssl: line 1: syntax error near unexpected token `<'
/usr/local/bin/cfssl: line 1: `<html><body>You are being <a href="https://github-production-release-asset-2e65be.s3.amazonaws.com/21591001/b5055500-148a-11eb-9528-44972336b695?X-Amz-Algorithm=AWS4-HMAC-SHA256&amp;X-Amz-Credential=AKIAIWNJYAX4CSVEH53A%2F20201227%2Fus-east-1%2Fs3%2Faws4_request&amp;X-Amz-Date=20201227T073955Z&amp;X-Amz-Expires=300&amp;X-Amz-Signature=2fe2795e562fea205b61f75af5e46ea020733dbf2ef8c09d1a5fc7f1ab4d5bf2&amp;X-Amz-SignedHeaders=host&amp;actor_id=0&amp;key_id=0&amp;repo_id=21591001&amp;response-content-disposition=attachment%3B%20filename%3Dcfssl_1.5.0_linux_amd64&amp;response-content-type=application%2Foctet-stream">redirected</a>.</body></html>'
```

I checked the URL too

https://github-production-release-asset-2e65be.s3.amazonaws.com/21591001/b5055500-148a-11eb-9528-44972336b695?X-Amz-Algorithm=AWS4-HMAC-SHA256&amp;X-Amz-Credential=AKIAIWNJYAX4CSVEH53A%2F20201227%2Fus-east-1%2Fs3%2Faws4_request&amp;X-Amz-Date=20201227T073955Z&amp;X-Amz-Expires=300&amp;X-Amz-Signature=2fe2795e562fea205b61f75af5e46ea020733dbf2ef8c09d1a5fc7f1ab4d5bf2&amp;X-Amz-SignedHeaders=host&amp;actor_id=0&amp;key_id=0&amp;repo_id=21591001&amp;response-content-disposition=attachment%3B%20filename%3Dcfssl_1.5.0_linux_amd64&amp;response-content-type=application%2Foctet-stream

Finally saw this

```bash
$ cfssl gencert --help
/usr/local/bin/cfssl: line 1: syntax error near unexpected token `<'
/usr/local/bin/cfssl: line 1: `<html><body>You are being <a href="https://github-production-release-asset-2e65be.s3.amazonaws.com/21591001/b5055500-148a-11eb-9528-44972336b695?X-Amz-Algorithm=AWS4-HMAC-SHA256&amp;X-Amz-Credential=AKIAIWNJYAX4CSVEH53A%2F20201227%2Fus-east-1%2Fs3%2Faws4_request&amp;X-Amz-Date=20201227T073955Z&amp;X-Amz-Expires=300&amp;X-Amz-Signature=2fe2795e562fea205b61f75af5e46ea020733dbf2ef8c09d1a5fc7f1ab4d5bf2&amp;X-Amz-SignedHeaders=host&amp;actor_id=0&amp;key_id=0&amp;repo_id=21591001&amp;response-content-disposition=attachment%3B%20filename%3Dcfssl_1.5.0_linux_amd64&amp;response-content-type=application%2Foctet-stream">redirected</a>.</body></html>'

$ cfssl gencert
/usr/local/bin/cfssl: line 1: syntax error near unexpected token `<'
/usr/local/bin/cfssl: line 1: `<html><body>You are being <a href="https://github-production-release-asset-2e65be.s3.amazonaws.com/21591001/b5055500-148a-11eb-9528-44972336b695?X-Amz-Algorithm=AWS4-HMAC-SHA256&amp;X-Amz-Credential=AKIAIWNJYAX4CSVEH53A%2F20201227%2Fus-east-1%2Fs3%2Faws4_request&amp;X-Amz-Date=20201227T073955Z&amp;X-Amz-Expires=300&amp;X-Amz-Signature=2fe2795e562fea205b61f75af5e46ea020733dbf2ef8c09d1a5fc7f1ab4d5bf2&amp;X-Amz-SignedHeaders=host&amp;actor_id=0&amp;key_id=0&amp;repo_id=21591001&amp;response-content-disposition=attachment%3B%20filename%3Dcfssl_1.5.0_linux_amd64&amp;response-content-type=application%2Foctet-stream">redirected</a>.</body></html>'


$ cfssl --help
/usr/local/bin/cfssl: line 1: syntax error near unexpected token `<'
/usr/local/bin/cfssl: line 1: `<html><body>You are being <a href="https://github-production-release-asset-2e65be.s3.amazonaws.com/21591001/b5055500-148a-11eb-9528-44972336b695?X-Amz-Algorithm=AWS4-HMAC-SHA256&amp;X-Amz-Credential=AKIAIWNJYAX4CSVEH53A%2F20201227%2Fus-east-1%2Fs3%2Faws4_request&amp;X-Amz-Date=20201227T073955Z&amp;X-Amz-Expires=300&amp;X-Amz-Signature=2fe2795e562fea205b61f75af5e46ea020733dbf2ef8c09d1a5fc7f1ab4d5bf2&amp;X-Amz-SignedHeaders=host&amp;actor_id=0&amp;key_id=0&amp;repo_id=21591001&amp;response-content-disposition=attachment%3B%20filename%3Dcfssl_1.5.0_linux_amd64&amp;response-content-type=application%2Foctet-stream">redirected</a>.</body></html>'

$ cfssl --version
/usr/local/bin/cfssl: line 1: syntax error near unexpected token `<'
/usr/local/bin/cfssl: line 1: `<html><body>You are being <a href="https://github-production-release-asset-2e65be.s3.amazonaws.com/21591001/b5055500-148a-11eb-9528-44972336b695?X-Amz-Algorithm=AWS4-HMAC-SHA256&amp;X-Amz-Credential=AKIAIWNJYAX4CSVEH53A%2F20201227%2Fus-east-1%2Fs3%2Faws4_request&amp;X-Amz-Date=20201227T073955Z&amp;X-Amz-Expires=300&amp;X-Amz-Signature=2fe2795e562fea205b61f75af5e46ea020733dbf2ef8c09d1a5fc7f1ab4d5bf2&amp;X-Amz-SignedHeaders=host&amp;actor_id=0&amp;key_id=0&amp;repo_id=21591001&amp;response-content-disposition=attachment%3B%20filename%3Dcfssl_1.5.0_linux_amd64&amp;response-content-type=application%2Foctet-stream">redirected</a>.</body></html>'

$ cat /usr/local/bin/cfssl
<html><body>You are being <a href="https://github-production-release-asset-2e65be.s3.amazonaws.com/21591001/b5055500-148a-11eb-9528-44972336b695?X-Amz-Algorithm=AWS4-HMAC-SHA256&amp;X-Amz-Credential=AKIAIWNJYAX4CSVEH53A%2F20201227%2Fus-east-1%2Fs3%2Faws4_request&amp;X-Amz-Date=20201227T073955Z&amp;X-Amz-Expires=300&amp;X-Amz-Signature=2fe2795e562fea205b61f75af5e46ea020733dbf2ef8c09d1a5fc7f1ab4d5bf2&amp;X-Amz-SignedHeaders=host&amp;actor_id=0&amp;key_id=0&amp;repo_id=21591001&amp;response-content-disposition=attachment%3B%20filename%3Dcfssl_1.5.0_linux_amd64&amp;response-content-type=application%2Foctet-stream">redirected</a>.</body></html>

$ cat /usr/local/bin/cfssljson
<html><body>You are being <a href="https://github-production-release-asset-2e65be.s3.amazonaws.com/21591001/2513db00-148b-11eb-82c0-e13cb9b9405f?X-Amz-Algorithm=AWS4-HMAC-SHA256&amp;X-Amz-Credential=AKIAIWNJYAX4CSVEH53A%2F20201227%2Fus-east-1%2Fs3%2Faws4_request&amp;X-Amz-Date=20201227T074138Z&amp;X-Amz-Expires=300&amp;X-Amz-Signature=8cb9de2d0b2ef0888c875c062965c4e008b4ea2a1b356738c9e5f726fc05d29f&amp;X-Amz-SignedHeaders=host&amp;actor_id=0&amp;key_id=0&amp;repo_id=21591001&amp;response-content-disposition=attachment%3B%20filename%3Dcfssljson_1.5.0_linux_amd64&amp;response-content-type=application%2Foctet-stream">redirected</a>.</body></html>
```

I didn't download it correctly ... 🤦

```bash
$ sudo curl -L https://github.com/cloudflare/cfssl/releases/download/v1.5.0/cfssljson_1.
5.0_linux_amd64 -o /usr/local/bin/cfssljson

$ sudo curl -L https://github.com/cloudflare/cfssl/releases/download/v1.5.0/cfssl_1.5.0_
linux_amd64 -o /usr/local/bin/cfssl
```

I didn't use `-L` for follow redirects. Right.

```bash
$ cfssl -h
Usage:
Available commands:
        gencsr
        revoke
        bundle
        certinfo
        crl
        genkey
...

$ $ cfssljson -h
Usage of cfssljson:
  -bare
        the response from CFSSL is not wrapped in the API standard response
  -f string
        JSON input (default "-")
...
```

I should have verified if the installation worked. Hmm :)

So, things worked! :)

```bash
$ cfssl gencert -initca ca-csr.json | cfssljson -bare ca
```

Next I need to create a key pair for service account stuff. Following this

https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/04-certificate-authority.md#the-service-account-key-pair

Creating a CA config which has some profile and defaults. I see that it is being
reused for almost everything. According to the config, the certificate can be
used for many purposes I guess. For server authentication, for client
authentication, signing and key encipherment. Idk what the last one means, maybe
verifying or something? Not sure, gotta check

```bash
$ cat > ca-config.json <<EOF
{
    "signing": {
        "default": {
            "expiry": "8760h"
        },
        "profiles": {
            "kubernetes": {
                "usages": [
                    "signing",
                    "key encipherment",
                    "server auth",
                    "client auth"
                ],
                "expiry": "8760h"
            }
        }
    }
}
EOF
```

```bash
$ cat > service-account-csr.json <<EOF
{
    "CN": "service-accounts",
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "names": [
        {
            "C": "US",
            "L": "San Francisco",
            "O": "Kubernetes",
            "OU": "Service Account",
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
  service-account-csr.json | cfssljson -bare service-account

2020/12/27 14:34:28 [INFO] generate received request
2020/12/27 14:34:28 [INFO] received CSR
2020/12/27 14:34:28 [INFO] generating key: rsa-2048
2020/12/27 14:34:29 [INFO] encoded CSR
2020/12/27 14:34:29 [INFO] signed certificate with serial number 638705263800501710861529130827853022936937882506
2020/12/27 14:34:29 [WARNING] This certificate lacks a "hosts" field. This makes it unsuitable for
websites. For more information see the Baseline Requirements for the Issuance and Management
of Publicly-Trusted Certificates, v.1.1.6, from the CA/Browser Forum (https://cabforum.org);
specifically, section 10.2.3 ("Information Requirements").
```


