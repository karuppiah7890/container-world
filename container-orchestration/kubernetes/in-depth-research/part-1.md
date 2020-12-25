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


