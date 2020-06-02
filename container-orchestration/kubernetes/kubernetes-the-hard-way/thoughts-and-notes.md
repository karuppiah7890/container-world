# Thoughts and Notes

I'm going to try to setup a Kubernetes Cluster, by using VMs in my local.
I'll be using `multipass` to create `ubuntu` VMs

I'm going through the pre-requisites now.
https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/01-prerequisites.md

Since I don't need GCP because I'm going to try everything in local, I don't
need most of these.

I noticed something really cool about tmux! So, apparently, we can run commands
on multiple panes in `tmux`, and I'm using `tmux` so I think it's going to be
really useful! :)

https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/01-prerequisites.md#running-commands-in-parallel-with-tmux

Next https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/02-client-tools.md

I have got the client tools.

Next https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/03-compute-resources.md

There are some cloud networking stuff, and public IP and stuff, which I'm
skipping. I don't think I need it for my local setup. If I hit some networking
issue, I'll get back here. From what I tested in my local, I'm able to connect
from one `multipass` VM to another `multipass`. I started a mock server using
`netcat -l 8080` and in the other machine ran `telnet <other-vm-ip> 8080` and
it worked, which ensured the connectivity. I also tested it with `ping`.
I followed this to some extent - https://ubidots.com/blog/how-to-simulate-a-tcpudp-client-using-netcat/

Now, I'm directly jumping to this
https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/03-compute-resources.md#compute-instances

They are providing fixed IPs and that's done by `multipass` for me.

---
Side note:
Also, there are some more repos which does the whole thing with local machine

https://github.com/mmumshad/kubernetes-the-hard-way
https://github.com/MOrlassino/kubernetes-the-hard-way

I'll refer to them I guess, when I get stuck.

---

Now, for provisioning instances, I'll just start off 3 `ubuntu` instances for
the Kubernetes Control Plane

```bash
$ multipass launch -n controller-0
$ multipass launch -n controller-1
$ multipass launch -n controller-2
$ multipass ls
Name                    State             IPv4             Image
containers              Stopped           --               Ubuntu 18.04 LTS
controller-0            Running           192.168.64.30    Ubuntu 18.04 LTS
controller-1            Running           192.168.64.28    Ubuntu 18.04 LTS
controller-2            Running           192.168.64.29    Ubuntu 18.04 LTS
faas                    Stopped           --               Ubuntu 18.04 LTS
konfigadm               Stopped           --               Ubuntu 18.04 LTS
stolon                  Stopped           --               Ubuntu 18.04 LTS
```

Now, for the worker instances, it says

```
Each worker instance requires a pod subnet allocation from the Kubernetes cluster CIDR range. The pod subnet allocation will be used to configure container networking in a later exercise. The pod-cidr instance metadata will be used to expose pod subnet allocations to compute instances at runtime.

The Kubernetes cluster CIDR range is defined by the Controller Manager's --cluster-cidr flag. In this tutorial the cluster CIDR range will be set to 10.200.0.0/16, which supports 254 subnets.
```

CIDR - Classless Inter-Domain Routing. 

I'm assuming this is all networking which is not specific to the cloud and I
can see it's a private network based on the numbers I'm seeing for the network
addresses. https://en.wikipedia.org/wiki/Private_network ,
https://en.wikipedia.org/wiki/Private_network#Private_IPv4_addresses

The problem now is, unlike the cloud, I'm not sure if we can set metadata for
a `multipass` VM.

Anyways, I'm going to skip the metadata part for now, and I can also see some
extra networking stuff that I'm skipping. For example, something about IP
Forwarding. I guess I'll do all that stuff when I get stuck because of it, like,
when it doesn't work :P I know where to look, the other repos which do it in
local. So, I'm gonna head to creating the worker instances.

```bash
$ multipass launch -n worker-0
$ multipass launch -n worker-1
$ multipass launch -n worker-2
$ multipass ls
Name                    State             IPv4             Image
containers              Stopped           --               Ubuntu 18.04 LTS
controller-0            Running           192.168.64.30    Ubuntu 18.04 LTS
controller-1            Running           192.168.64.28    Ubuntu 18.04 LTS
controller-2            Running           192.168.64.29    Ubuntu 18.04 LTS
faas                    Stopped           --               Ubuntu 18.04 LTS
konfigadm               Stopped           --               Ubuntu 18.04 LTS
stolon                  Stopped           --               Ubuntu 18.04 LTS
worker-0                Running           192.168.64.32    Ubuntu 18.04 LTS
worker-1                Running           192.168.64.33    Ubuntu 18.04 LTS
worker-2                Running           192.168.64.34    Ubuntu 18.04 LTS
```

And I can ssh into these machines, using

```bash
$ multipass shell <machine-name>
$ # or
$ multipass exec <machine-name> bash
```

Now, next is
https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/04-certificate-authority.md

I need to bootstrap a Certificate Authority (CA), and then use it to generate 
TLS certificates (issue certificate) for the following components: etcd,
kube-apiserver, kube-controller-manager, kube-scheduler, kubelet, and kube-proxy

For the [CA](https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/04-certificate-authority.md#certificate-authority)
I created the `ca-config.json` for the Certificate Authority configuration,
and then the `ca-csr.json` Certificate Authority's (CA) Certificate Signing
Request (CSR) and then used `cfssl` and `cfssljson` and created the public
key and private key for the CA.

For the [client and server certificate](https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/04-certificate-authority.md#client-and-server-certificates)

I have created the public and private key for the admin user

Going to create a certificate for each worker node that meets the Kubernetes Node
Authorizer's requirements using this

```bash
$ for instance in worker-0 worker-1 worker-2; do
cat > ${instance}-csr.json <<EOF
{
  "CN": "system:node:${instance}",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "US",
      "L": "Portland",
      "O": "system:nodes",
      "OU": "Kubernetes The Hard Way",
      "ST": "Oregon"
    }
  ]
}
EOF

INTERNAL_IP=$(multipass info ${instance} --format json | jq -r ".info[\"$instance\"].ipv4[0]")

cfssl gencert \
  -ca=ca.pem \
  -ca-key=ca-key.pem \
  -config=ca-config.json \
  -hostname=${instance},${INTERNAL_IP} \
  -profile=kubernetes \
  ${instance}-csr.json | cfssljson -bare ${instance}
done
```

Next gotta generate the kube-controller-manager client certificate and private key

Next is kube-proxy client certificate and private key

Next is the kube-scheduler client certificate and private key

Next is the Kubernetes API Server certificate and private key

For this, I think we need something extra. In this step here
https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/03-compute-resources.md#kubernetes-public-ip-address

They created a public IP adress, which belonged to an external load balancer
(LB). In our case, since we don't have the cloud and we are running local
instances, we need to run our own Load Balancer in front of the API servers,
and this has also been done by the other repos that do this in the local.
I have to see how to configure this Load Balancer. Just the usualy nginx or
some similar software I guess. For now I'll just spin up an instance use the
IP, but I'll have to configure it later.

Other options for load balancers are HA Proxies. Also, I might have to properly
configure the load balancer, or else when doing API server uprgades, it could
go wrong, and send traffic to unhealthy API servers, which I have experienced
before in GKE cloud service due to some issue in their Load Balancer. And many
other issues could occur if the Load Balancer is not working correctly 😅

Anyways. 

```bash
$ multipass launch -n kube-api-loadbalancer
$ multipass ls
Name                    State             IPv4             Image
containers              Stopped           --               Ubuntu 18.04 LTS
controller-0            Running           192.168.64.30    Ubuntu 18.04 LTS
controller-1            Running           192.168.64.28    Ubuntu 18.04 LTS
controller-2            Running           192.168.64.29    Ubuntu 18.04 LTS
faas                    Stopped           --               Ubuntu 18.04 LTS
konfigadm               Stopped           --               Ubuntu 18.04 LTS
kube-api-loadbalancer   Running           192.168.64.35    Ubuntu 18.04 LTS
stolon                  Stopped           --               Ubuntu 18.04 LTS
worker-0                Running           192.168.64.32    Ubuntu 18.04 LTS
worker-1                Running           192.168.64.33    Ubuntu 18.04 LTS
worker-2                Running           192.168.64.34    Ubuntu 18.04 LTS
```

And now, continuing with the kube api server certificate and private key
generation with this

```bash
$ KUBERNETES_API_LOADBALANCER_IP=$(multipass info kube-api-loadbalancer --format json | jq -r '.info."kube-api-loadbalancer".ipv4[0]')
$ KUBERNETES_MASTER_IP_0=$(multipass info controller-0 --format json | jq -r '.info."controller-0".ipv4[0]')
$ KUBERNETES_MASTER_IP_1=$(multipass info controller-1 --format json | jq -r '.info."controller-1".ipv4[0]')
$ KUBERNETES_MASTER_IP_2=$(multipass info controller-2 --format json | jq -r '.info."controller-2".ipv4[0]')

$ KUBERNETES_HOSTNAMES=kubernetes,kubernetes.default,kubernetes.default.svc,kubernetes.default.svc.cluster,kubernetes.svc.cluster.local

$ cat > kubernetes-csr.json <<EOF
{
  "CN": "kubernetes",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "US",
      "L": "Portland",
      "O": "Kubernetes",
      "OU": "Kubernetes The Hard Way",
      "ST": "Oregon"
    }
  ]
}
EOF

$ cfssl gencert \
  -ca=ca.pem \
  -ca-key=ca-key.pem \
  -config=ca-config.json \
  -hostname=10.32.0.1,${KUBERNETES_MASTER_IP_0},${KUBERNETES_MASTER_IP_1},${KUBERNETES_MASTER_IP_2},${KUBERNETES_API_LOADBALANCER_IP},127.0.0.1,${KUBERNETES_HOSTNAMES} \
  -profile=kubernetes \
  kubernetes-csr.json | cfssljson -bare kubernetes
```

It seems

```
The Kubernetes API server is automatically assigned the kubernetes internal dns name, which will be linked to the first IP address (10.32.0.1) from the address range (10.32.0.0/24) reserved for internal cluster services during the control plane bootstrapping lab.
```

So, that's what the `10.32.0.1` means.

Next is service account key pair - the service-account certificate and private key

Now I need to distribute all these client and server certificates
https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/04-certificate-authority.md#distribute-the-client-and-server-certificates

For worker instances, I'm going to use this

```bash
for instance in worker-0 worker-1 worker-2; do
  multipass transfer ca.pem ${instance}-key.pem ${instance}.pem ${instance}:~/
done
```

and for controller instances this

```bash
for instance in controller-0 controller-1 controller-2; do
  multipass transfer ca.pem ca-key.pem kubernetes-key.pem kubernetes.pem \
    service-account-key.pem service-account.pem ${instance}:~/
done
```

For some weird reason, `transfer` command doesn't work, so I'm taking a short
cut, not a very good looking one!

This one is for worker instances

```bash
$ for instance in worker-0 worker-1 worker-2; do
  multipass mount certificate-stuff ${instance}:/certificate-stuff
  multipass exec ${instance} -- cp /certificate-stuff/ca.pem /certificate-stuff/${instance}-key.pem /certificate-stuff/${instance}.pem /home/ubuntu
  multipass umount ${instance}:/certificate-stuff
done
```

This one is for controller instances

```bash
$ for instance in controller-0 controller-1 controller-2; do
  multipass mount certificate-stuff ${instance}:/certificate-stuff
  multipass exec ${instance} -- cp /certificate-stuff/ca.pem /certificate-stuff/ca-key.pem \
/certificate-stuff/kubernetes-key.pem /certificate-stuff/kubernetes.pem \
/certificate-stuff/service-account-key.pem /certificate-stuff/service-account.pem \
/home/ubuntu
  multipass umount ${instance}:/certificate-stuff
done
```

Apparently the kube-proxy, kube-controller-manager, kube-scheduler, and kubelet
client certificates will be used to generate client authentication configuration
files

So, like we saw, the next is
https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/05-kubernetes-configuration-files.md

So, we need to generate kubeernetes configurations - kubeconfigs, which enable
Kubernetes clients to locate and authenticate to the Kubernetes API Servers.

And we need to generate kubeconfig files for the controller manager, kubelet,
kube-proxy, and scheduler clients and the admin user.

Each kubeconfig requires a Kubernetes API Server to connect to. Since there
are many API server instances, we will use the Load Balancer to connect to the
API server and hence use the Load Balancer IP in the kubeconfig

```bash
$ KUBERNETES_API_LOADBALANCER_IP=$(multipass info kube-api-loadbalancer --format json | jq -r '.info."kube-api-loadbalancer".ipv4[0]')
```

I'll be using the below to generate the kubeconfig for the kubelets (worker
instances)

```bash
$ for instance in worker-0 worker-1 worker-2; do
  kubectl config set-cluster kubernetes-the-hard-way \
    --certificate-authority=ca.pem \
    --embed-certs=true \
    --server=https://${KUBERNETES_API_LOADBALANCER_IP}:6443 \
    --kubeconfig=${instance}.kubeconfig

  kubectl config set-credentials system:node:${instance} \
    --client-certificate=${instance}.pem \
    --client-key=${instance}-key.pem \
    --embed-certs=true \
    --kubeconfig=${instance}.kubeconfig

  kubectl config set-context default \
    --cluster=kubernetes-the-hard-way \
    --user=system:node:${instance} \
    --kubeconfig=${instance}.kubeconfig

  kubectl config use-context default --kubeconfig=${instance}.kubeconfig
done
```

And it gave the output as

```bash
Cluster "kubernetes-the-hard-way" set.
User "system:node:worker-0" set.
Context "default" created.
Switched to context "default".
Cluster "kubernetes-the-hard-way" set.
User "system:node:worker-1" set.
Context "default" created.
Switched to context "default".
Cluster "kubernetes-the-hard-way" set.
User "system:node:worker-2" set.
Context "default" created.
Switched to context "default".
```

For kube-proxy, I'm going to use this

```bash
  kubectl config set-cluster kubernetes-the-hard-way \
    --certificate-authority=ca.pem \
    --embed-certs=true \
    --server=https://${KUBERNETES_API_LOADBALANCER_IP}:6443 \
    --kubeconfig=kube-proxy.kubeconfig

  kubectl config set-credentials system:kube-proxy \
    --client-certificate=kube-proxy.pem \
    --client-key=kube-proxy-key.pem \
    --embed-certs=true \
    --kubeconfig=kube-proxy.kubeconfig

  kubectl config set-context default \
    --cluster=kubernetes-the-hard-way \
    --user=system:kube-proxy \
    --kubeconfig=kube-proxy.kubeconfig

  kubectl config use-context default --kubeconfig=kube-proxy.kubeconfig
```

I'm using some stuff different from the tutorial because of how I named some
things differently. For example `KUBERNETES_API_LOADBALANCER_IP` instead of
`KUBERNETES_PUBLIC_ADDRESS` in the tutorial. I think the former is a better
name in this case, though I could have just used any name and put the correct
value and it would have worked. Anyways.

Next, I need to create kubeconfig for kube-controller-manager. That's done.

Next is kube-scheduler kubeconfig. That's also done!

Next is admin user kubeconfig. Final one. That's also done!

Also, I just realized that, the braces `{` `}` in the shell code, was actually
valid. Without the braces, they run weirdly in my `bash` shell. But, with it,
they work really well. I think it denotes a block of commands! :)

Now, we need to distribute the kubeconfig files to the appropriate machines!

Copy the appropriate kubelet and kube-proxy kube config files to each worker
instance. Unfortunately, can't use `transfer` command in `multipass` as it's not
working.

```bash
for instance in worker-0 worker-1 worker-2; do
  multipass mount kubeconfigs ${instance}:/kubeconfigs
  multipass exec ${instance} -- cp /kubeconfigs/${instance}.kubeconfig /kubeconfigs/kube-proxy.kubeconfig /home/ubuntu
  multipass umount ${instance}:/kubeconfigs
done
```

Copy the appropriate admin, kube-controller-manager and kube-scheduler
kubeconfig files to each controller instance:

```bash
for instance in controller-0 controller-1 controller-2; do
  multipass mount kubeconfigs ${instance}:/kubeconfigs
  multipass exec ${instance} -- cp /kubeconfigs/admin.kubeconfig /kubeconfigs/kube-controller-manager.kubeconfig /kubeconfigs/kube-scheduler.kubeconfig /home/ubuntu
  multipass umount ${instance}:/kubeconfigs
done
```

Now, the next is
https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/06-data-encryption-keys.md

I need to create a kubernetes resource called `EncryptionConfig`, which will
help with encrypting data at rest given an encryption key

I need to copy this encryption config into the controller instances.

Unfortunately, can't use `transfer` command in `multipass` as it's not
working.

```bash
for instance in controller-0 controller-1 controller-2; do
  multipass mount kube-resources ${instance}:/kube-resources
  multipass exec ${instance} -- cp /kube-resources/encryption-config.yaml /home/ubuntu
  multipass umount ${instance}:/kube-resources
done
```

Now, that's done!

Next is
https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/07-bootstrapping-etcd.md

`etcd` stuff is going to be fun I think! ;) :D

So, we are going to run our `etcd` members in each of the `controller`
instances. And we are going to use the tmux's synchronize panes cool feature! ;)

So, in between I needed to get the IP address of the linux machine. The private
IP address. Apparently `hostname` can help with that

https://stackoverflow.com/questions/21336126/linux-bash-script-to-extract-ip-address

```bash
$ hostname -I
```

That gives all the IP addresses of the host. It works on linux. Doesn't work on
unix, in my MacOS.

Anyways, back to `etcd` stuff. I'm going to use this for the internal IP address
thing

```bash
$ INTERNAL_IP=$(hostname -I)
```

Okay, so I got an error now, while trying to start the `etcd` service

```bash
invalid value "https://192.168.64.29" for flag -initial-advertise-peer-urls: URL address does not have the form "host:port": https://192.168.64.29
```

I'm going to change it and see what happens.

Oh wait. I didn't read the error properly or see the command properly. The
command came out like this

```bash
$ /usr/local/bin/etcd \
  --name controller-2 \
  --cert-file=/etc/etcd/kubernetes.pem \
  --key-file=/etc/etcd/kubernetes-key.pem \
  --peer-cert-file=/etc/etcd/kubernetes.pem \
  --peer-key-file=/etc/etcd/kubernetes-key.pem \
  --trusted-ca-file=/etc/etcd/ca.pem \
  --peer-trusted-ca-file=/etc/etcd/ca.pem \
  --peer-client-cert-auth \
  --client-cert-auth \
  --initial-advertise-peer-urls 192.168.64.29 :2380 \
  --listen-peer-urls https://192.168.64.29 :2380 \
  --listen-client-urls https://192.168.64.29 :2379,https://127.0.0.1:2379 \
  --advertise-client-urls https://192.168.64.29 :2379 \
  --initial-cluster-token etcd-cluster-0 \
  --initial-cluster controller-0=https://10.240.0.10:2380,controller-1=https://10.240.0.11:2380,controller-2=https://10.240.0.12:2380 \
  --initial-cluster-state new \
  --data-dir=/var/lib/etcd
```

The value for the flag `--initial-advertise-peer-urls` and some more flags,
have a space in it. It's because of an issue in the `INTERNAL_IP` variable's
value.

```bash
ubuntu@controller-2:~$ echo ${INTERNAL_IP}ok
192.168.64.29 ok
```

Right. Not something I had expected. Let's fix that.

This looks better now

```bash
ubuntu@controller-2:~$ INTERNAL_IP=$(hostname -I | sed 's/ //')
ubuntu@controller-2:~$ echo ${INTERNAL_IP}ok
192.168.64.29ok
```

I fixed that issue, but seems like I didn't check the service file properly.
Let me read it this time 🙈 The error -

```
etcdmain: --initial-cluster has controller-2=https://10.240.0.12:2380 but missing from --initial-advertise-peer-urls=https://192.168.64.29:2380 ("https://192.168.64.29:2380"(resolved from "https://192.168.64.29:2380") != "https://10.240.0.12:2380"(resolved from "https://10.240.0.12:2380"))
```

Seems like I need to give the controller instance IP addresses here, which are
hardcoded here. In my case, I need to use this

```bash
$ cat <<EOF | sudo tee /etc/systemd/system/etcd.service
[Unit]
Description=etcd
Documentation=https://github.com/coreos

[Service]
Type=notify
ExecStart=/usr/local/bin/etcd \\
  --name ${ETCD_NAME} \\
  --cert-file=/etc/etcd/kubernetes.pem \\
  --key-file=/etc/etcd/kubernetes-key.pem \\
  --peer-cert-file=/etc/etcd/kubernetes.pem \\
  --peer-key-file=/etc/etcd/kubernetes-key.pem \\
  --trusted-ca-file=/etc/etcd/ca.pem \\
  --peer-trusted-ca-file=/etc/etcd/ca.pem \\
  --peer-client-cert-auth \\
  --client-cert-auth \\
  --initial-advertise-peer-urls https://${INTERNAL_IP}:2380 \\
  --listen-peer-urls https://${INTERNAL_IP}:2380 \\
  --listen-client-urls https://${INTERNAL_IP}:2379,https://127.0.0.1:2379 \\
  --advertise-client-urls https://${INTERNAL_IP}:2379 \\
  --initial-cluster-token etcd-cluster-0 \\
  --initial-cluster controller-0=https://192.168.64.30:2380,controller-1=https://192.168.64.28:2380,controller-2=https://192.168.64.29:2380 \\
  --initial-cluster-state new \\
  --data-dir=/var/lib/etcd
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
```

Okay, it still didn't work. It said this error

```bash
May 31 21:58:55 controller-0 etcd[27426]: rejected connection from "192.168.64.28:49768" (error "remote error: tls: bad certificate", ServerName "")
May 31 21:58:55 controller-0 etcd[27426]: rejected connection from "192.168.64.28:49764" (error "remote error: tls: bad certificate", ServerName "")
May 31 21:58:55 controller-0 etcd[27426]: rejected connection from "192.168.64.29:41072" (error "remote error: tls: bad certificate", ServerName "")
```

Hmm. I wonder why it's a bad certificate. I think I need to check the
certificates 🙈

Yeah. Like I thought. I made a mistake. I wrote down stuff. But I didn't
run them. The issue is mentioned here 

https://github.com/kelseyhightower/kubernetes-the-hard-way/issues/570

And I was checking `etcd` issues 🤦‍♂ I mean, there are some `etcd` issues, but
they were not the cause of my problem.

Cool. I made it right and recreated the `kubernetes.pem` and `kubernetes-key.pem`.
I can't imagine what other issues are going to come 🙈

```bash
$ openssl x509 -in kubernetes.pem -text -noout | rg Address
```

There really must be checks done like the above. Or else it's just 🤷‍♂ Hmm...

Now I need to put the keys in the controller instances again! And do stuff! All
over again!

With that, I just had to copy the cert and key to an etcd directory and then
had to restart the service, so I just ran these again

```bash
{
  sudo systemctl daemon-reload
  sudo systemctl enable etcd
  sudo systemctl start etcd
}
```

And everything worked! :D Yay!

Logs:

```
ubuntu@controller-2:~$ systemctl status etcd
● etcd.service - etcd
   Loaded: loaded (/etc/systemd/system/etcd.service; enabled; vendor preset: enabled)
   Active: active (running) since Sun 2020-05-31 22:25:18 IST; 2min 8s ago
     Docs: https://github.com/coreos
 Main PID: 19731 (etcd)
    Tasks: 8 (limit: 1152)
   CGroup: /system.slice/etcd.service
           └─19731 /usr/local/bin/etcd --name controller-2 --cert-file=/etc/etcd/kubernetes.pem --key-file=/etc/etcd/kubernetes-key.pem -

May 31 22:25:18 controller-2 etcd[19731]: raft2020/05/31 22:25:18 INFO: 263824555f9f3950 [logterm: 1, index: 3, vote: 0] cast MsgVote for
May 31 22:25:18 controller-2 etcd[19731]: raft2020/05/31 22:25:18 INFO: raft.node: 263824555f9f3950 elected leader 1f1ba19f7a762c3a at te
May 31 22:25:18 controller-2 etcd[19731]: published {Name:controller-2 ClientURLs:[https://192.168.64.29:2379]} to cluster e78997be5ad126
May 31 22:25:18 controller-2 systemd[1]: Started etcd.
May 31 22:25:18 controller-2 etcd[19731]: ready to serve client requests
```

```
ubuntu@controller-1:~$ systemctl status etcd
● etcd.service - etcd
   Loaded: loaded (/etc/systemd/system/etcd.service; enabled; vendor preset: enabled)
   Active: active (running) since Sun 2020-05-31 22:28:18 IST; 52s ago
     Docs: https://github.com/coreos
 Main PID: 27630 (etcd)
    Tasks: 8 (limit: 1152)
   CGroup: /system.slice/etcd.service
           └─27630 /usr/local/bin/etcd --name controller-1 --cert-file=/etc/etcd/kubernetes.pem --key-file=/etc/etcd/kubernetes-key.pem -

May 31 22:28:18 controller-1 etcd[27630]: established a TCP streaming connection with peer 1f1ba19f7a762c3a (stream Message writer)
May 31 22:28:18 controller-1 etcd[27630]: established a TCP streaming connection with peer 263824555f9f3950 (stream Message reader)
May 31 22:28:18 controller-1 etcd[27630]: established a TCP streaming connection with peer 263824555f9f3950 (stream MsgApp v2 reader)
May 31 22:28:18 controller-1 etcd[27630]: published {Name:controller-1 ClientURLs:[https://192.168.64.28:2379]} to cluster e78997be5ad126
May 31 22:28:18 controller-1 systemd[1]: Started etcd.
May 31 22:28:18 controller-1 etcd[27630]: ready to serve client requests
May 31 22:28:18 controller-1 etcd[27630]: serving client requests on 127.0.0.1:2379
May 31 22:28:18 controller-1 etcd[27630]: ready to serve client requests
May 31 22:28:18 controller-1 etcd[27630]: serving client requests on 192.168.64.28:2379
May 31 22:28:18 controller-1 etcd[27630]: b2a5547e041c2f9b initialized peer connection; fast-forwarding 8 ticks (election ticks 10) with
```

I actually just restarted one of them alone, as the logs looked a bit
suspicious 😅

```
ubuntu@controller-0:~$ systemctl status etcd
● etcd.service - etcd
   Loaded: loaded (/etc/systemd/system/etcd.service; enabled; vendor preset: enabled)
   Active: active (running) since Sun 2020-05-31 22:25:18 IST; 3min 18s ago
     Docs: https://github.com/coreos
 Main PID: 28455 (etcd)
    Tasks: 8 (limit: 1152)
   CGroup: /system.slice/etcd.service
           └─28455 /usr/local/bin/etcd --name controller-0 --cert-file=/etc/etcd/kubernetes.pem --key-file=/etc/etcd/kubernetes-key.pem -

May 31 22:28:16 controller-0 etcd[28455]: lost the TCP streaming connection with peer b2a5547e041c2f9b (stream Message reader)
May 31 22:28:16 controller-0 etcd[28455]: failed to dial b2a5547e041c2f9b on stream Message (read tcp 192.168.64.30:35792->192.168.64.28:
May 31 22:28:16 controller-0 etcd[28455]: peer b2a5547e041c2f9b became inactive (message send to peer failed)
May 31 22:28:17 controller-0 etcd[28455]: lost the TCP streaming connection with peer b2a5547e041c2f9b (stream Message writer)
May 31 22:28:18 controller-0 etcd[28455]: peer b2a5547e041c2f9b became active
May 31 22:28:18 controller-0 etcd[28455]: closed an existing TCP streaming connection with peer b2a5547e041c2f9b (stream MsgApp v2 writer
May 31 22:28:18 controller-0 etcd[28455]: established a TCP streaming connection with peer b2a5547e041c2f9b (stream MsgApp v2 writer)
May 31 22:28:18 controller-0 etcd[28455]: established a TCP streaming connection with peer b2a5547e041c2f9b (stream MsgApp v2 reader)
May 31 22:28:18 controller-0 etcd[28455]: established a TCP streaming connection with peer b2a5547e041c2f9b (stream Message writer)
May 31 22:28:18 controller-0 etcd[28455]: established a TCP streaming connection with peer b2a5547e041c2f9b (stream Message reader)
```

Finally, checked the members list too, using the `etcd` client `etcdctl`

```bash
ubuntu@controller-1:~$ sudo ETCDCTL_API=3 etcdctl member list \
>   --endpoints=https://127.0.0.1:2379 \
>   --cacert=/etc/etcd/ca.pem \
>   --cert=/etc/etcd/kubernetes.pem \
>   --key=/etc/etcd/kubernetes-key.pem
1f1ba19f7a762c3a, started, controller-0, https://192.168.64.30:2380, https://192.168.64.30:2379, false
263824555f9f3950, started, controller-2, https://192.168.64.29:2380, https://192.168.64.29:2379, false
b2a5547e041c2f9b, started, controller-1, https://192.168.64.28:2380, https://192.168.64.28:2379, false
```

Yay! :D 

Actually, till now I've just been following instructions, after I'm done with
all this, I need to go back and try to understand what's going on in the whole
thing and the big picture.

Note: Also, with respect to the load balancer, I could actually connect to the api
servers directly first and then later work on the load balancer and may be
even try upgrading the api servers one by one, and have the load balancer for
balancing the load and making sure the api servers are highly available, and
also try to upgrade the worker instance components, and also see how auto
scaling of worker nodes happens and how to do that in local, instead of cloud.
That would be fun to check out ;) :D 

Now, back to the tutorial. Next is
https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/08-bootstrapping-kubernetes-controllers.md

So, in the controller instances, we are going to install 3 components -
1. API server
2. Scheduler
3. Controller Manager

We are again going to use tmux's synchronize panes feature ;) in the controller
instances

Again, for Internal IP, I had to use this

```bash
$ INTERNAL_IP=$(hostname -I | sed 's/ //')
```

For service file, I need to use this

```bash
cat <<EOF | sudo tee /etc/systemd/system/kube-apiserver.service
[Unit]
Description=Kubernetes API Server
Documentation=https://github.com/kubernetes/kubernetes

[Service]
ExecStart=/usr/local/bin/kube-apiserver \\
  --advertise-address=${INTERNAL_IP} \\
  --allow-privileged=true \\
  --apiserver-count=3 \\
  --audit-log-maxage=30 \\
  --audit-log-maxbackup=3 \\
  --audit-log-maxsize=100 \\
  --audit-log-path=/var/log/audit.log \\
  --authorization-mode=Node,RBAC \\
  --bind-address=0.0.0.0 \\
  --client-ca-file=/var/lib/kubernetes/ca.pem \\
  --enable-admission-plugins=NamespaceLifecycle,NodeRestriction,LimitRanger,ServiceAccount,DefaultStorageClass,ResourceQuota \\
  --etcd-cafile=/var/lib/kubernetes/ca.pem \\
  --etcd-certfile=/var/lib/kubernetes/kubernetes.pem \\
  --etcd-keyfile=/var/lib/kubernetes/kubernetes-key.pem \\
  --etcd-servers=https://192.168.64.28:2379,https://192.168.64.29:2379,https://192.168.64.30:2379 \\
  --event-ttl=1h \\
  --encryption-provider-config=/var/lib/kubernetes/encryption-config.yaml \\
  --kubelet-certificate-authority=/var/lib/kubernetes/ca.pem \\
  --kubelet-client-certificate=/var/lib/kubernetes/kubernetes.pem \\
  --kubelet-client-key=/var/lib/kubernetes/kubernetes-key.pem \\
  --kubelet-https=true \\
  --runtime-config=api/all \\
  --service-account-key-file=/var/lib/kubernetes/service-account.pem \\
  --service-cluster-ip-range=10.32.0.0/24 \\
  --service-node-port-range=30000-32767 \\
  --tls-cert-file=/var/lib/kubernetes/kubernetes.pem \\
  --tls-private-key-file=/var/lib/kubernetes/kubernetes-key.pem \\
  --v=2
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
```

Those are the api server stuff, next is kube-controller-manager

That's done

Next is kube-scheduler

That's done too

I started off all of them. And now I get this error in the api server

```
clientconn.go:1251] grpc: addrConn.createTransport failed to connect to {https: 0  <nil>}. Err :connection error: desc = "transport: Error while dialing dial tcp: lookup https on 127.0.0.53:53: server misbehaving". Reconnecting...
```

when checking with

```bash
$ systemctl status kube-apiserver
```

And the other components - controller manager and the scheduler are failing
because they are not able to connect to the api server

```
dial tcp 127.0.0.1:6443: connect: connection refused
```

Looking at the api server error, the first thought in my mind was - where is that
IP coming from? And the port `53`, it's DNS server port. I doube checked it too.
So, it's trying to do some DNS resolution and it's failing

Okay, I found out where that IP came from

```bash
ubuntu@controller-2:~$ cat /etc/resolv.conf
# This file is managed by man:systemd-resolved(8). Do not edit.
#
# This is a dynamic resolv.conf file for connecting local clients to the
# internal DNS stub resolver of systemd-resolved. This file lists all
# configured search domains.
#
# Run "systemd-resolve --status" to see details about the uplink DNS servers
# currently in use.
#
# Third party programs must not access this file directly, but only through the
# symlink at /etc/resolv.conf. To manage man:resolv.conf(5) in a different way,
# replace this symlink by a static file or a different symlink.
#
# See man:systemd-resolved.service(8) for details about the supported modes of
# operation for /etc/resolv.conf.

nameserver 127.0.0.53
options edns0
```

From what I know, the api servers don't need to know each other. That's a guess,
given I hear it's stateless and the `etcd` is the state and I can see that I
have given the IP addresses of the `etcd` to the api server. So, I'm gonna
just see how to get rid of this DNS server from the config or what to do about
it.

I actually changed the nameserver to `8.8.8.8` and now api server says

```bash
Error while dialing dial tcp: lookup https on 8.8.8.8:53: no such host"
```

I gotta see what to do next. Also makes me think what was previously happening,
when it said "server misbehaving", was the server not running, which means
connection refused, or was it running but not responding? Hmmm

I also noticed another issue. I noticed this thing in the service file

```
https://https://192.168.64.28:2379,https://https://192.168.64.29:2379,https://https://192.168.64.30:2379
```

and some logs like this

```
client.go:354] parsed scheme: ""
8300    4029 client.go:354] scheme "" not registered, fallback to default scheme
```

Damn. I need to fix that first, the etcd server urls

I changed the service file and ran stuff

```bash
ubuntu@controller-0:~$ sudo systemctl restart kube-apiserver
Warning: The unit file, source configuration file or drop-ins of kube-apiserver.service changed on disk. Run 'systemctl daemon-reload' to reload units.
ubuntu@controller-0:~$ sudo systemctl daemon-reload
ubuntu@controller-0:~$ sudo systemctl restart kube-apiserver
```

Okay, so the previous error is gone. I think I need to change my nameserver
back to whatever it was? I'll do that for now. But I do remember seeing this
https://github.com/mmumshad/kubernetes-the-hard-way/blob/master/docs/02-compute-resources.md
Okay, there they mentioned it for the node to be able to access the Internet.
Not the biggest of my concerns now. I'll get back to that when I want to run
stuff, like pull docker images from docker hub and run pods, which will need
access to the Internet

There's a new error now though 🙈

```
http: TLS handshake error from 127.0.0.1:37856: remote error: tls: bad certificate
```

Again. A certificate issue. I need to see which certificate has gone wrong this
time!

Okay wait. So. I'm seeing these logs for the API server. Now, it says that the
error is **from** `127.0.0.1`. I don't know, but I'm looking at two
possibilities now, one is, the other components running in the `localhost`,
which is the controller manager and scheduler are not able to authenticate due
to some issue in their end or the server has some issue on it's end...or both?
🙈 Let's find out. I'll check the logs of other components

Okay, I checked the kube-scheduler and it said

```
x509: certificate is valid for 192.168.64.29, 10.0.0.1, not 127.0.0.1
```

So, it's an issue from the server side I think. The certificate that the server
is giving is not valid for IP / host `127.0.0.1` which the kube-scheduler is
using to access the api server

Same error for controller-manager

```
x509: certificate is valid for 192.168.64.29, 10.0.0.1, not 127.0.0.1
```

Okay, may be it's not a server side issue. I noticed different errors in other
machines. Should have checked before. Anyways, this is what I noticed on the
server side cert first

```
ubuntu@controller-1:/var/lib/kubernetes$ openssl x509 -in kubernetes.pem -text -noout | grep Address
                DNS:kubernetes, DNS:kubernetes.default, DNS:kubernetes.default.svc, DNS:kubernetes.default.svc.cluster, DNS:kubernetes.svc.cluster.local, IP Address:10.32.0.1, IP Address:192.168.64.30, IP Address:192.168.64.28, IP Address:192.168.64.29, IP Address:192.168.64.35, IP Address:127.0.0.1
```

So, `127.0.0.1` is present. And I checked this in all machines. It seems to be
all good.

In other machines I see this

```
x509: certificate is valid for 192.168.64.30, 10.0.0.1, not 127.0.0.1
```

```
x509: certificate is valid for 192.168.64.28, 10.0.0.1, not 127.0.0.1
```

So, I need to see what's going on over here. I mean, it's clearly mentioning
it's own internal IP, and then another IP `10.0.0.1` which I don't get where
it came from, but it says it's not okay for `127.0.0.1`

I guess I need to understand the error first, and replicate it

I can see `kubectl` telling the same thing

```bash
ubuntu@controller-1:~$ kubectl get componentstatuses --kubeconfig admin.kubeconfig
Unable to connect to the server: x509: certificate is valid for 192.168.64.28, 10.0.0.1, not 127.0.0.1
```

When I changed the `admin.kubeconfig` to point to `192.168.64.28`, it says this

```bash
ubuntu@controller-1:~$ kubectl get componentstatuses --kubeconfig admin.kubeconfig
Unable to connect to the server: x509: certificate signed by unknown authority
```

Clearly, the server certificate has some issues.

I'm able to do this

```bash
ubuntu@controller-1:~$ kubectl get componentstatuses --insecure-skip-tls-verify --kubeconfig admin.kubeconfig
NAME                 STATUS    MESSAGE             ERROR
controller-manager   Healthy   ok
scheduler            Healthy   ok
etcd-1               Healthy   {"health":"true"}
etcd-0               Healthy   {"health":"true"}
etcd-2               Healthy   {"health":"true"}
```

Okay. Now we need to check why this is happening. The `127.0.0.1` issue and
the `certificate signed by unknown authority`

Wow. I made one big mistake 🙈 while creating the api server service file,
somehow I put some extra spaces and some of the arguments were not passed to
the api server, which included the `tls-cert-file` and `tls-private-key`, which
explains the tls certificate issues.

I ran stuff and restarted it all, and again, in one of the machines alone,
the `$INTERNAL_IP` was not available, which was needed for the service file,
so, need to fix that alone

Okay, so things look good now. I restarted the controller manager and scheduler
too, even though the final status looked good

```bash
$ sudo systemctl restart kube-controller-manager
$ sudo systemctl restart kube-scheduler
```

The controller manager status does tell something about client ca file and
request header client ca file not being present, and that client cert auth and
request header client cert auth won't work

Verification works

```bash
ubuntu@controller-0:~$ kubectl get componentstatuses --kubeconfig admin.kubeconfig
NAME                 STATUS    MESSAGE             ERROR
controller-manager   Healthy   ok
scheduler            Healthy   ok
etcd-2               Healthy   {"health":"true"}
etcd-0               Healthy   {"health":"true"}
etcd-1               Healthy   {"health":"true"}
ubuntu@controller-0:~$ curl -k -i https://127.0.0.1:6443/healthz
HTTP/2 200
content-type: text/plain; charset=utf-8
x-content-type-options: nosniff
content-length: 2
date: Mon, 01 Jun 2020 07:49:47 GMT

ok
```

Next there's some RBAC stuff for kubelet
https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/08-bootstrapping-kubernetes-controllers.md#rbac-for-kubelet-authorization

I need to run some commands, but just from one controller node, using `kubectl`
to apply cluster role and cluster role binding

Next is provisioning a front end load balancer, which I'm not going to do,
so I'll just use one of the controller instance IPs and do the check from my
host machine

```bash
$ KUBERNETES_MASTER_IP_0=$(multipass info controller-0 --format json | jq -r '.info."controller-0".ipv4[0]')
$ curl --cacert certificate-stuff/ca.pem https://$KUBERNETES_MASTER_IP_0:6443/version
{
  "major": "1",
  "minor": "15",
  "gitVersion": "v1.15.3",
  "gitCommit": "2d3c76f9091b6bec110a5e63777c332469e0cba2",
  "gitTreeState": "clean",
  "buildDate": "2019-08-19T11:05:50Z",
  "goVersion": "go1.12.9",
  "compiler": "gc",
  "platform": "linux/amd64"
}
```

For now I'm moving to the next section, which is
https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/09-bootstrapping-kubernetes-workers.md

Here we are going to bootstrap the worker nodes. The following components will
be installed on each node:
1. runc https://github.com/opencontainers/runc
2. container networking plugins https://github.com/containernetworking/cni
3. containerd https://github.com/containerd/containerd
4. kubelet
5. kube-proxy

We are again using tmux's synchronize tab's feature to install all these in
all the worker instances ;)

Good that I noticed I have all the controller instance shells opened up. I
closed them and I need to get into the shells of the worker instances
using `multipass shell`

Apparently the `socat` binary enables support for the kubectl port-forward command.

That's interesting. I think I could learn how that's done :)

Anyways, I got some errors while installing OS packages

```bash
E: Could not get lock /var/lib/dpkg/lock-frontend - open (11: Resource temporarily unavailable)
E: Unable to acquire the dpkg frontend lock (/var/lib/dpkg/lock-frontend), is another process using it?
```

But they went away when I tried installing again! :)

Okay, my swap is already disabled and that's what is recommended it seems, having
swap disabled

Next I'm installing the worker binaries! :D 

I have installed the stuff now in worker instances. Next I need to put the
`$POD_CIDR` value in each instance separately as I don't have any metadata
available

In `worker-0`, it will be

```bash
POD_CIDR=10.200.0.0/24
```

For `worker-1` and `worker-2`

```bash
POD_CIDR=10.200.1.0/24
```

```bash
POD_CIDR=10.200.2.0/24
```

Okay, I did all the configurations and everything and started off the kubelet,
kube-proxy and containerd

containerd status is good. But I'm seeing errors now for kubelet and kube proxy.

kubelet error in `worker-0` node

```bash
ubuntu@worker-0:~$ systemctl status kubelet
● kubelet.service - Kubernetes Kubelet
   Loaded: loaded (/etc/systemd/system/kubelet.service; enabled; vendor preset: enabled)
   Active: active (running) since Mon 2020-06-01 15:12:57 IST; 6min ago
     Docs: https://github.com/kubernetes/kubernetes
 Main PID: 2255 (kubelet)
    Tasks: 13 (limit: 1152)
   CGroup: /system.slice/kubelet.service
           └─2255 /usr/local/bin/kubelet --config=/var/lib/kubelet/kubelet-config.yaml --container-runtime=remote --container-runtime-endpoint=unix:///var/run/containerd/containerd.sock --image-pull-progress-deadline=2m --kubeconfig=/var/lib/kubelet/kubeconfig --network-plugin=cni --register-node=true --v=2

Jun 01 15:19:06 worker-0 kubelet[2255]: E0601 15:19:06.051247    2255 kubelet.go:2248] node "worker-0" not found
Jun 01 15:19:06 worker-0 kubelet[2255]: E0601 15:19:06.202007    2255 kubelet.go:2248] node "worker-0" not found
Jun 01 15:19:06 worker-0 kubelet[2255]: E0601 15:19:06.212600    2255 reflector.go:125] k8s.io/client-go/informers/factory.go:133: Failed to list *v1beta1.RuntimeClass: Get https://192.168.64.35:6443/apis/node.k8s.io/v1beta1/runtimeclasses?limit=500&resourceVersion=0: dial tcp 192.168.64.35:6443: connect: connection refused
Jun 01 15:19:06 worker-0 kubelet[2255]: E0601 15:19:06.214465    2255 reflector.go:125] k8s.io/kubernetes/pkg/kubelet/config/apiserver.go:47: Failed to list *v1.Pod: Get https://192.168.64.35:6443/api/v1/pods?fieldSelector=spec.nodeName%3Dworker-0&limit=500&resourceVersion=0: dial tcp 192.168.64.35:6443: connect: connection refused
Jun 01 15:19:06 worker-0 kubelet[2255]: E0601 15:19:06.215098    2255 reflector.go:125] k8s.io/kubernetes/pkg/kubelet/kubelet.go:453: Failed to list *v1.Node: Get https://192.168.64.35:6443/api/v1/nodes?fieldSelector=metadata.name%3Dworker-0&limit=500&resourceVersion=0: dial tcp 192.168.64.35:6443: connect: connection refused
Jun 01 15:19:06 worker-0 kubelet[2255]: E0601 15:19:06.215534    2255 reflector.go:125] k8s.io/kubernetes/pkg/kubelet/kubelet.go:444: Failed to list *v1.Service: Get https://192.168.64.35:6443/api/v1/services?limit=500&resourceVersion=0: dial tcp 192.168.64.35:6443: connect: connection refused
Jun 01 15:19:06 worker-0 kubelet[2255]: E0601 15:19:06.215915    2255 reflector.go:125] k8s.io/client-go/informers/factory.go:133: Failed to list *v1beta1.CSIDriver: Get https://192.168.64.35:6443/apis/storage.k8s.io/v1beta1/csidrivers?limit=500&resourceVersion=0: dial tcp 192.168.64.35:6443: connect: connection refused
Jun 01 15:19:06 worker-0 kubelet[2255]: E0601 15:19:06.302920    2255 kubelet.go:2248] node "worker-0" not found
Jun 01 15:19:06 worker-0 kubelet[2255]: E0601 15:19:06.403637    2255 kubelet.go:2248] node "worker-0" not found
Jun 01 15:19:06 worker-0 kubelet[2255]: E0601 15:19:06.576939    2255 kubelet.go:2248] node "worker-0" not found
```

Something similar in kube-proxy in `worker-0` node

```
ubuntu@worker-0:~$ systemctl status kube-proxy
● kube-proxy.service - Kubernetes Kube Proxy
   Loaded: loaded (/etc/systemd/system/kube-proxy.service; enabled; vendor preset: enabled)
   Active: active (running) since Mon 2020-06-01 15:12:57 IST; 6min ago
     Docs: https://github.com/kubernetes/kubernetes
 Main PID: 2251 (kube-proxy)
    Tasks: 0 (limit: 1152)
   CGroup: /system.slice/kube-proxy.service
           └─2251 /usr/local/bin/kube-proxy --config=/var/lib/kube-proxy/kube-proxy-config.yaml

Jun 01 15:18:43 worker-0 kube-proxy[2251]: E0601 15:18:43.510441    2251 reflector.go:125] k8s.io/client-go/informers/factory.go:133: Failed to list *v1.Service: Get https://192.168.64.35:6443/api/v1/services?labelSelector=%21service.kubernetes.io%2Fservice-proxy-name&limit=500&resourceVersion=0: dial tcp 192.168.64.35:6443: connect: connection refused
Jun 01 15:18:43 worker-0 kube-proxy[2251]: E0601 15:18:43.526583    2251 reflector.go:125] k8s.io/client-go/informers/factory.go:133: Failed to list *v1.Endpoints: Get https://192.168.64.35:6443/api/v1/endpoints?labelSelector=%21service.kubernetes.io%2Fservice-proxy-name&limit=500&resourceVersion=0: dial tcp 192.168.64.35:6443: connect: connection refused
Jun 01 15:18:44 worker-0 kube-proxy[2251]: E0601 15:18:44.512033    2251 reflector.go:125] k8s.io/client-go/informers/factory.go:133: Failed to list *v1.Service: Get https://192.168.64.35:6443/api/v1/services?labelSelector=%21service.kubernetes.io%2Fservice-proxy-name&limit=500&resourceVersion=0: dial tcp 192.168.64.35:6443: connect: connection refused
Jun 01 15:18:44 worker-0 kube-proxy[2251]: E0601 15:18:44.529400    2251 reflector.go:125] k8s.io/client-go/informers/factory.go:133: Failed to list *v1.Endpoints: Get https://192.168.64.35:6443/api/v1/endpoints?labelSelector=%21service.kubernetes.io%2Fservice-proxy-name&limit=500&resourceVersion=0: dial tcp 192.168.64.35:6443: connect: connection refused
Jun 01 15:18:45 worker-0 kube-proxy[2251]: E0601 15:18:45.518728    2251 reflector.go:125] k8s.io/client-go/informers/factory.go:133: Failed to list *v1.Service: Get https://192.168.64.35:6443/api/v1/services?labelSelector=%21service.kubernetes.io%2Fservice-proxy-name&limit=500&resourceVersion=0: dial tcp 192.168.64.35:6443: connect: connection refused
Jun 01 15:18:45 worker-0 kube-proxy[2251]: E0601 15:18:45.532745    2251 reflector.go:125] k8s.io/client-go/informers/factory.go:133: Failed to list *v1.Endpoints: Get https://192.168.64.35:6443/api/v1/endpoints?labelSelector=%21service.kubernetes.io%2Fservice-proxy-name&limit=500&resourceVersion=0: dial tcp 192.168.64.35:6443: connect: connection refused
Jun 01 15:18:46 worker-0 kube-proxy[2251]: E0601 15:18:46.528911    2251 reflector.go:125] k8s.io/client-go/informers/factory.go:133: Failed to list *v1.Service: Get https://192.168.64.35:6443/api/v1/services?labelSelector=%21service.kubernetes.io%2Fservice-proxy-name&limit=500&resourceVersion=0: dial tcp 192.168.64.35:6443: connect: connection refused
Jun 01 15:18:46 worker-0 kube-proxy[2251]: E0601 15:18:46.537093    2251 reflector.go:125] k8s.io/client-go/informers/factory.go:133: Failed to list *v1.Endpoints: Get https://192.168.64.35:6443/api/v1/endpoints?labelSelector=%21service.kubernetes.io%2Fservice-proxy-name&limit=500&resourceVersion=0: dial tcp 192.168.64.35:6443: connect: connection refused
Jun 01 15:18:47 worker-0 kube-proxy[2251]: E0601 15:18:47.561436    2251 reflector.go:125] k8s.io/client-go/informers/factory.go:133: Failed to list *v1.Service: Get https://192.168.64.35:6443/api/v1/services?labelSelector=%21service.kubernetes.io%2Fservice-proxy-name&limit=500&resourceVersion=0: dial tcp 192.168.64.35:6443: connect: connection refused
Jun 01 15:18:47 worker-0 kube-proxy[2251]: E0601 15:18:47.561584    2251 reflector.go:125] k8s.io/client-go/informers/factory.go:133: Failed to list *v1.Endpoints: Get https://192.168.64.35:6443/api/v1/endpoints?labelSelector=%21service.kubernetes.io%2Fservice-proxy-name&limit=500&resourceVersion=0: dial tcp 192.168.64.35:6443: connect: connection refused

```

I'm not sure where this IP `192.168.64.35`, which is the IP I was planning to
use for the loadbalancer, came into the picture. Looks like I might need to
set it up now.

Okay, I found out where it gets it's configuration. As usual - only kubeconfigs
have the config to locate the API server and authenticate with the API server.
And I have used the Loadbalancer IP for that. Hmm. Now. I could go and change
the config, or just setup the loadbalancer. Hmm. I think I'll just setup the
loadbalancer!

I'm checking out nginx for the loadbalancer. And if it doesn't workout for some
reason, then haproxy.

So, in nginx, I want to listen at port `6443` and proxy pass all requests that
come, to one of the api servers, based on health checks. I'm not sure how to do
the health check part exactly, so, I'll skip that for now, and think about it
later.

Apparently, instead of configuring a http server, which will need SSL
certificates for https and which will do SSL termination, at least that's what
I read, we are going to blindly proxy pass the tcp connection directly to one
of the api servers.

I skimmed through some QAs, links

https://stackoverflow.com/questions/39420613/can-nginx-do-tcp-load-balance-with-ssl-termination/39421271#39421271
https://www.cyberciti.biz/faq/configure-nginx-ssltls-passthru-with-tcp-load-balancing/
https://www.alibabacloud.com/blog/how-to-use-nginx-as-an-https-forward-proxy-server_595799
http://nginx.org/en/docs/stream/ngx_stream_core_module.html

Seems like one nginx feature can help and I'm just doing a small demo like thing,
so as long as it works, it's fine and enough!

So, I'm following this
https://www.cyberciti.biz/faq/configure-nginx-ssltls-passthru-with-tcp-load-balancing/

And this is the config file I'm using `/etc/nginx/passthrough.conf` which is
included in the `nginx.conf` file at the root level

```
## tcp LB  and SSL passthrough for backend ##
stream {
    upstream kubeapiserver {
        server 192.168.64.28:6443 max_fails=3 fail_timeout=10s;
        server 192.168.64.29:6443 max_fails=3 fail_timeout=10s;
        server 192.168.64.30:6443 max_fails=3 fail_timeout=10s;
    }

    log_format basic '$remote_addr [$time_local] '
                 '$protocol $status $bytes_sent $bytes_received '
                 '$session_time "$upstream_addr" '
                 '"$upstream_bytes_sent" "$upstream_bytes_received" "$upstream_connect_time"';

    access_log /var/log/nginx/ssl_passthrough_access.log basic;
    error_log  /var/log/nginx/ssl_passthrough_error.log;

    server {
        listen 6443;
        proxy_pass kubeapiserver;
        proxy_next_upstream on;
    }
}
```

And configuration test was successful too

```bash
ubuntu@kube-api-loadbalancer:~$ sudo nginx -t
nginx: the configuration file /etc/nginx/nginx.conf syntax is ok
nginx: configuration file /etc/nginx/nginx.conf test is successful
ubuntu@kube-api-loadbalancer:~$ echo $?
0
```

```bash
ubuntu@kube-api-loadbalancer:~$ sudo systemctl reload nginx
```

And I tested it too! It works! Yay! :D 

```bash
ubuntu@kube-api-loadbalancer:~$ curl https://localhost:6443
curl: (60) SSL certificate problem: unable to get local issuer certificate
More details here: https://curl.haxx.se/docs/sslcerts.html

curl failed to verify the legitimacy of the server and therefore could not
establish a secure connection to it. To learn more about this situation and
how to fix it, please visit the web page mentioned above.
ubuntu@kube-api-loadbalancer:~$ curl -k https://localhost:6443/version
{
  "major": "1",
  "minor": "15",
  "gitVersion": "v1.15.3",
  "gitCommit": "2d3c76f9091b6bec110a5e63777c332469e0cba2",
  "gitTreeState": "clean",
  "buildDate": "2019-08-19T11:05:50Z",
  "goVersion": "go1.12.9",
  "compiler": "gc",
  "platform": "linux/amd64"
}
```

The proxy does not have any SSL certificates and it can still listen to HTTPS
requests and proxy the backend which talk HTTP protocol using the SSL
Certificates that they have!

I restarted the kube-proxy in one of the workers. But, I saw that the logs
were no longer showing errors, so I just went and checked the nodes in the
controller instance

```bash
ubuntu@controller-1:~$ kubectl get nodes --kubeconfig admin.kubeconfig
NAME       STATUS   ROLES    AGE    VERSION
worker-0   Ready    <none>   3m3s   v1.15.3
worker-1   Ready    <none>   3m3s   v1.15.3
worker-2   Ready    <none>   3m3s   v1.15.3
```

And it's all working! Yay! :D

Finally!! Next is
https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/10-configuring-kubectl.md

For remote `kubectl` access, I'm using this in my local machine in the
`certificate-stuff` directory

```bash
{
  KUBERNETES_API_LOADBALANCER_IP=$(multipass info kube-api-loadbalancer --format json | jq -r '.info."kube-api-loadbalancer".ipv4[0]')

  kubectl config set-cluster kubernetes-the-hard-way \
    --certificate-authority=ca.pem \
    --embed-certs=true \
    --server=https://${KUBERNETES_API_LOADBALANCER_IP}:6443

  kubectl config set-credentials admin \
    --client-certificate=admin.pem \
    --client-key=admin-key.pem

  kubectl config set-context kubernetes-the-hard-way \
    --cluster=kubernetes-the-hard-way \
    --user=admin

  kubectl config use-context kubernetes-the-hard-way
}
```

But it doesn't work :/

```bash
$ kubectl get componentstatuses
NAME                 AGE
controller-manager   <unknown>
scheduler            <unknown>
etcd-1               <unknown>
etcd-2               <unknown>
etcd-0               <unknown>
```

I can see `nodes` though

```bash
$ kubectl get nodes
NAME       STATUS   ROLES    AGE   VERSION
worker-0   Ready    <none>   15m   v1.15.3
worker-1   Ready    <none>   15m   v1.15.3
worker-2   Ready    <none>   15m   v1.15.3
```

Not sure what's going on. I can see the component statuses in all the controller
machines very fine. It looks like this

```bash
ubuntu@controller-0:~$ kubectl get componentstatuses
NAME                 STATUS    MESSAGE             ERROR
controller-manager   Healthy   ok
scheduler            Healthy   ok
etcd-0               Healthy   {"health":"true"}
etcd-2               Healthy   {"health":"true"}
etcd-1               Healthy   {"health":"true"}
```

I think I'll skip this issue for now. And then come back if this is an issue
later. 

I can actually run other `kubectl` command to get other resources, and it all
works. So, gonna move on to the next section

Next is
https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/11-pod-network-routes.md

So, this is about pods getting an IP address from the node's pod CIDR range,
and how pods cannot access other pods running in another node due to missing
network routes. It's a problem of routing and it's being solved by adding
routes, using `gcloud` which for GCP cloud. But I need to do something locally
for this. I need to see how to do it. For now, I'll come back to this later.
Let me see what else is left out and how much this setup works for me without
some stuff - that is, with some stuff missing.

Next up is
https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/12-dns-addon.md

The DNS, core DNS.

And I ran this

```bash
$ kubectl apply -f https://storage.googleapis.com/kubernetes-the-hard-way/coredns.yaml
serviceaccount/coredns created
clusterrole.rbac.authorization.k8s.io/system:coredns created
clusterrolebinding.rbac.authorization.k8s.io/system:coredns created
configmap/coredns created
deployment.extensions/coredns created
service/kube-dns created
```

On checking stuff, I see some errors :/

```bash

$ kubectl get pods -l k8s-app=kube-dns -n kube-system
NAME                     READY   STATUS    RESTARTS   AGE
coredns-5fb99965-mlswr   0/1     Running   0          58s
coredns-5fb99965-r8kkt   0/1     Running   0          58s
$ kubectl logs -f -n kube-system coredns-5fb99965-mlswr
Error from server: Get https://worker-1:10250/containerLogs/kube-system/coredns-5fb99965-mlswr/coredns?follow=true: dial tcp: lookup worker-1 on 127.0.0.53:53: server misbehaving
```

It's not able to resolve what's `worker-1`. And it's trying to access the local
DNS server which I don't know if it's running or not. I'm just going to see if
the node name's IP is configured correctly in `/etc/hosts` or do just that

Weirdly, I see this in `worker-0` machine

```bash
ubuntu@worker-0:~$ cat /etc/hosts
# Your system has configured 'manage_etc_hosts' as True.
# As a result, if you wish for changes to this file to persist
# then you will need to either
# a.) make changes to the master file in /etc/cloud/templates/hosts.debian.tmpl
# b.) change or remove the value of 'manage_etc_hosts' in
#     /etc/cloud/cloud.cfg or cloud-config from user-data
#
127.0.1.1 worker-0 worker-0
127.0.0.1 localhost

# The following lines are desirable for IPv6 capable hosts
::1 ip6-localhost ip6-loopback
fe00::0 ip6-localnet
ff00::0 ip6-mcastprefix
ff02::1 ip6-allnodes
ff02::2 ip6-allrouters
ff02::3 ip6-allhosts
```

```bash
ubuntu@worker-0:~$ dig worker-0

; <<>> DiG 9.11.3-1ubuntu1.12-Ubuntu <<>> worker-0
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 48169
;; flags: qr rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 65494
;; QUESTION SECTION:
;worker-0.                      IN      A

;; ANSWER SECTION:
worker-0.               0       IN      A       127.0.1.1

;; Query time: 0 msec
;; SERVER: 127.0.0.53#53(127.0.0.53)
;; WHEN: Mon Jun 01 16:59:12 IST 2020
;; MSG SIZE  rcvd: 53
```


```bash
ubuntu@worker-0:~$ nslookup worker-0
Server:         127.0.0.53
Address:        127.0.0.53#53

Non-authoritative answer:
Name:   worker-0
Address: 127.0.1.1
```

```bash
ubuntu@worker-0:~$ dig @127.0.0.53 worker-0

; <<>> DiG 9.11.3-1ubuntu1.12-Ubuntu <<>> @127.0.0.53 worker-0
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 39937
;; flags: qr rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 65494
;; QUESTION SECTION:
;worker-0.                      IN      A

;; ANSWER SECTION:
worker-0.               0       IN      A       127.0.1.1

;; Query time: 0 msec
;; SERVER: 127.0.0.53#53(127.0.0.53)
;; WHEN: Mon Jun 01 17:00:06 IST 2020
;; MSG SIZE  rcvd: 53
```

I'm not sure why it says the `server misbehaving` actually. Damn 🙈 Even this
works

```bash
ubuntu@worker-0:~$ telnet 127.0.0.53 53
Trying 127.0.0.53...
Connected to 127.0.0.53.
Escape character is '^]'.
^]
telnet> Connection closed.
```

So, some DNS server is running and working correctly I think!

Oops. So. I was checking the configuration in the `worker-0` node, but the
error is coming from the pod logs. It's not able to understand how to resolve
`worker-0` and it says that `127.0.0.1:53:53: server misbehaving`. 

Let me see how to help the container be able to resolve the node names.
Because of this, I'm not even able to execute into a pod, however the containers
are dying continuously I think. For some reason the restart count for the pods
show `0`. `execpod` error -

```bash
execpod -n kube-system

 kubectl exec --namespace='kube-system' coredns-5fb99965-mlswr -c coredns -it sh

kubectl exec [POD] [COMMAND] is DEPRECATED and will be removed in a future version. Use kubectl kubectl exec [POD] -- [COMMAND] instead.
Error from server: error dialing backend: dial tcp: lookup worker-1 on 127.0.0.53:53: server misbehaving
```

I also simply tried deleting a pod and it works ;) 

Okay, now, I need to check the DNS issue for node names.

Actually. I just realized one thing. It's not an issue for the container to
resolve the node name `worker-0`. In fact I have been getting the same error for
executing into pod, and the same for logs. I didn't read an important part of
the error. Damn! 🙈 It clearly says "Error from server". Any command I type in
my local, interacts with kube api server. And it says that the error is from the
api server and it's not a log from the container. So, I guess I need to put
some entries in my api servers for the worker node's names, let me do that
first!

Also, wow, to put these entries I was checking the IP address of the worker
nodes and I see this

```bash
ubuntu@worker-0:~$ hostname -I
192.168.64.32 10.200.0.1
```

```bash
ubuntu@worker-1:~$ hostname -I
192.168.64.33 10.200.1.1
```

```bash
ubuntu@worker-2:~$ hostname -I
192.168.64.34 10.200.2.1
```

Let's put these entries in the controller nodes, with the `192.168.64.x` IP.
The other IP is part of another network interface, related to the Container
Networking Interface. I think it will be cool to understand how that works! :)

Now, in the controller instances, I went and put this in the `/etc/hosts` file

```bash
$ sudo vi /etc/hosts
```

and put this content

```
192.168.64.32 worker-0
192.168.64.33 worker-1
192.168.64.34 worker-2
```

And now, when I try to see logs, I get tons of errors like this

```bash
E0601 13:31:58.071259       1 reflector.go:126] pkg/mod/k8s.io/client-go@v11.0.0+incompatible/tools/cache/reflector.go:94: Failed to list *v1.Namespace: Get https://10.0.0.1:443/api/v1/namespaces?limit=500&resourceVersion=0: x509: certificate is valid for 10.32.0.1, 192.168.64.30, 192.168.64.28, 192.168.64.29, 192.168.64.35, 127.0.0.1, not 10.0.0.1
E0601 13:31:59.072181       1 reflector.go:126] pkg/mod/k8s.io/client-go@v11.0.0+incompatible/tools/cache/reflector.go:94: Failed to list *v1.Endpoints: Get https://10.0.0.1:443/api/v1/endpoints?limit=500&resourceVersion=0: x509: certificate is valid for 10.32.0.1, 192.168.64.30, 192.168.64.28, 192.168.64.29, 192.168.64.35, 127.0.0.1, not 10.0.0.1
E0601 13:31:59.072181       1 reflector.go:126] pkg/mod/k8s.io/client-go@v11.0.0+incompatible/tools/cache/reflector.go:94: Failed to list *v1.Endpoints: Get https://10.0.0.1:443/api/v1/endpoints?limit=500&resourceVersion=0: x509: certificate is valid for 10.32.0.1, 192.168.64.30, 192.168.64.28, 192.168.64.29, 192.168.64.35, 127.0.0.1, not 10.0.0.1
E0601 13:31:59.072181       1 reflector.go:126] pkg/mod/k8s.io/client-go@v11.0.0+incompatible/tools/cache/reflector.go:94: Failed to list *v1.Endpoints: Get https://10.0.0.1:443/api/v1/endpoints?limit=500&resourceVersion=0: x509: certificate is valid for 10.32.0.1, 192.168.64.30, 192.168.64.28, 192.168.64.29, 192.168.64.35, 127.0.0.1, not 10.0.0.1
E0601 13:31:59.072181       1 reflector.go:126] pkg/mod/k8s.io/client-go@v11.0.0+incompatible/tools/cache/reflector.go:94: Failed to list *v1.Endpoints: Get https://10.0.0.1:443/api/v1/endpoints?limit=500&resourceVersion=0: x509: certificate is valid for 10.32.0.1, 192.168.64.30, 192.168.64.28, 192.168.64.29, 192.168.64.35, 127.0.0.1, not 10.0.0.1
E0601 13:31:59.075359       1 reflector.go:126] pkg/mod/k8s.io/client-go@v11.0.0+incompatible/tools/cache/reflector.go:94: Failed to list *v1.Service: Get https://10.0.0.1:443/api/v1/services?limit=500&resourceVersion=0: x509: certificate is valid for 10.32.0.1, 192.168.64.30, 192.168.64.28, 192.168.64.29, 192.168.64.35, 127.0.0.1, not 10.0.0.1
E0601 13:31:59.075359       1 reflector.go:126] pkg/mod/k8s.io/client-go@v11.0.0+incompatible/tools/cache/reflector.go:94: Failed to list *v1.Service: Get https://10.0.0.1:443/api/v1/services?limit=500&resourceVersion=0: x509: certificate is valid for 10.32.0.1, 192.168.64.30, 192.168.64.28, 192.168.64.29, 192.168.64.35, 127.0.0.1, not 10.0.0.1
E0601 13:31:59.075359       1 reflector.go:126] pkg/mod/k8s.io/client-go@v11.0.0+incompatible/tools/cache/reflector.go:94: Failed to list *v1.Service: Get https://10.0.0.1:443/api/v1/services?limit=500&resourceVersion=0: x509: certificate is valid for 10.32.0.1, 192.168.64.30, 192.168.64.28, 192.168.64.29, 192.168.64.35, 127.0.0.1, not 10.0.0.1
E0601 13:31:59.075359       1 reflector.go:126] pkg/mod/k8s.io/client-go@v11.0.0+incompatible/tools/cache/reflector.go:94: Failed to list *v1.Service: Get https://10.0.0.1:443/api/v1/services?limit=500&resourceVersion=0: x509: certificate is valid for 10.32.0.1, 192.168.64.30, 192.168.64.28, 192.168.64.29, 192.168.64.35, 127.0.0.1, not 10.0.0.1
E0601 13:31:59.077532       1 reflector.go:126] pkg/mod/k8s.io/client-go@v11.0.0+incompatible/tools/cache/reflector.go:94: Failed to list *v1.Namespace: Get https://10.0.0.1:443/api/v1/namespaces?limit=500&resourceVersion=0: x509: certificate is valid for 10.32.0.1, 192.168.64.30, 192.168.64.28, 192.168.64.29, 192.168.64.35, 127.0.0.1, not 10.0.0.1
E0601 13:31:59.077532       1 reflector.go:126] pkg/mod/k8s.io/client-go@v11.0.0+incompatible/tools/cache/reflector.go:94: Failed to list *v1.Namespace: Get https://10.0.0.1:443/api/v1/namespaces?limit=500&resourceVersion=0: x509: certificate is valid for 10.32.0.1, 192.168.64.30, 192.168.64.28, 192.168.64.29, 192.168.64.35, 127.0.0.1, not 10.0.0.1
E0601 13:31:59.077532       1 reflector.go:126] pkg/mod/k8s.io/client-go@v11.0.0+incompatible/tools/cache/reflector.go:94: Failed to list *v1.Namespace: Get https://10.0.0.1:443/api/v1/namespaces?limit=500&resourceVersion=0: x509: certificate is valid for 10.32.0.1, 192.168.64.30, 192.168.64.28, 192.168.64.29, 192.168.64.35, 127.0.0.1, not 10.0.0.1
E0601 13:31:59.077532       1 reflector.go:126] pkg/mod/k8s.io/client-go@v11.0.0+incompatible/tools/cache/reflector.go:94: Failed to list *v1.Namespace: Get https://10.0.0.1:443/api/v1/namespaces?limit=500&resourceVersion=0: x509: certificate is valid for 10.32.0.1, 192.168.64.30, 192.168.64.28, 192.168.64.29, 192.168.64.35, 127.0.0.1, not 10.0.0.1
E0601 13:32:00.133156       1 reflector.go:126] pkg/mod/k8s.io/client-go@v11.0.0+incompatible/tools/cache/reflector.go:94: Failed to list *v1.Namespace: Get https://10.0.0.1:443/api/v1/namespaces?limit=500&resourceVersion=0: x509: certificate is valid for 10.32.0.1, 192.168.64.30, 192.168.64.28, 192.168.64.29, 192.168.64.35, 127.0.0.1, not 10.0.0.1
E0601 13:32:00.134007       1 reflector.go:126] pkg/mod/k8s.io/client-go@v11.0.0+incompatible/tools/cache/reflector.go:94: Failed to list *v1.Endpoints: Get https://10.0.0.1:443/api/v1/endpoints?limit=500&resourceVersion=0: x509: certificate is valid for 10.32.0.1, 192.168.64.30, 192.168.64.28, 192.168.64.29, 192.168.64.35, 127.0.0.1, not 10.0.0.1
E0601 13:32:00.133156       1 reflector.go:126] pkg/mod/k8s.io/client-go@v11.0.0+incompatible/tools/cache/reflector.go:94: Failed to list *v1.Namespace: Get https://10.0.0.1:443/api/v1/namespaces?limit=500&resourceVersion=0: x509: certificate is valid for 10.32.0.1, 192.168.64.30, 192.168.64.28, 192.168.64.29, 192.168.64.35, 127.0.0.1, not 10.0.0.1
....
....
2020-06-01T12:55:30.055Z [INFO] plugin/ready: Still waiting on: "kubernetes"
....
....
E0601 12:50:45.585740       1 reflector.go:126] pkg/mod/k8s.io/client-go@v11.0.0+incompatible/tools/cache/reflector.go:94: Failed to list *v1.Service: Get https://10.0.0.1:443/api/v1/services?limit=500&resourceVersion=0: x509: certificate is valid for 10.32.0.1, 192.168.64.30, 192.168.64.28, 192.168.64.29, 192.168.64.35, 127.0.0.1, not 10.0.0.1
E0601 12:50:46.629102       1 reflector.go:126] pkg/mod/k8s.io/client-go@v11.0.0+incompatible/tools/cache/reflector.go:94: Failed to list *v1.Namespace: Get https://10.0.0.1:443/api/v1/namespaces?limit=500&resourceVersion=0: x509: certificate is valid for 10.32.0.1, 192.168.64.30, 192.168.64.28, 192.168.64.29, 192.168.64.35, 127.0.0.1, not 10.0.0.1
E0601 12:50:46.629102       1 reflector.go:126] pkg/mod/k8s.io/client-go@v11.0.0+incompatible/tools/cache/reflector.go:94: Failed to list *v1.Namespace: Get https://10.0.0.1:443/api/v1/namespaces?limit=500&resourceVersion=0: x509: certificate is valid for 10.32.0.1, 192.168.64.30, 192.168.64.28, 192.168.64.29, 192.168.64.35, 127.0.0.1, not 10.0.0.1
E0601 12:50:46.629102       1 reflector.go:126] pkg/mod/k8s.io/client-go@v11.0.0+incompatible/tools/cache/reflector.go:94: Failed to list *v1.Namespace: Get https://10.0.0.1:443/api/v1/namespaces?limit=500&resourceVersion=0: x509: certificate is valid for 10.32.0.1, 192.168.64.30, 192.168.64.28, 192.168.64.29, 192.168.64.35, 127.0.0.1, not 10.0.0.1
E0601 12:50:46.629102       1 reflector.go:126] pkg/mod/k8s.io/client-go@v11.0.0+incompatible/tools/cache/reflector.go:94: Failed to list *v1.Namespace: Get https://10.0.0.1:443/api/v1/namespaces?limit=500&resourceVersion=0: x509: certificate is valid for 10.32.0.1, 192.168.64.30, 192.168.64.28, 192.168.64.29, 192.168.64.35, 127.0.0.1, not 10.0.0.1
```

It was just crazy to look at the overflow of logs, it felt like it's the same
old logs. I had my doubts regarding this mysterious IP `10.0.0.1`. Not sure how
this came up to be. So, there's this service in the cluster

```bash
$ kubectl get service
NAME         TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
kubernetes   ClusterIP   10.0.0.1     <none>        443/TCP   7h4m
```

I don't know how kubernetes provides IP addresses to services. Seeing all
services

```bash
$ kubectl get svc -A
NAMESPACE     NAME         TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)                  AGE
default       kubernetes   ClusterIP   10.0.0.1     <none>        443/TCP                  7h5m
kube-system   kube-dns     ClusterIP   10.32.0.10   <none>        53/UDP,53/TCP,9153/TCP   145m
```

If you notice, only `kubernetes` service has that weird IP. And `kube-dns` has
the IP `10.32.0.10` which seems to be part of the network that was provided
as input to the kube api server using the `service-cluster-ip-range`

```
--service-cluster-ip-range=10.32.0.0/24
```

Since `24` corresponds to the subnet mask `255.255.255.0`, I think the above
looks pretty good. But yeah, I'm not sure where the other IPs were used, like
`10.32.0.1`, `10.32.0.2` and so on till `10.32.0.9`.


But hey! Wow! I just tried to execute into a pod and it worked! It just worked!
No errors about node name resolution! :D

```bash
$ execpod

 kubectl exec --namespace='default' busybox -c busybox -it sh

kubectl exec [POD] [COMMAND] is DEPRECATED and will be removed in a future version. Use kubectl kubectl exec [POD] -- [COMMAND] instead.
/ # ls
bin   dev   etc   home  proc  root  sys   tmp   usr   var
/ #
```

Okay, inside this busybox pod, I tried to do a `nslookup` and of course it
failed.

```
/ # nslookup kubernetes
Server:    10.32.0.10
Address 1: 10.32.0.10
nslookup: can't resolve 'kubernetes'
/ # cat /etc/host
hostname  hosts
/ # cat /etc/hosts
# Kubernetes-managed hosts file.
127.0.0.1       localhost
::1     localhost ip6-localhost ip6-loopback
fe00::0 ip6-localnet
fe00::0 ip6-mcastprefix
fe00::1 ip6-allnodes
fe00::2 ip6-allrouters
10.200.1.4      busybox
/ # cat /etc/resolv.conf
search default.svc.cluster.local svc.cluster.local cluster.local
nameserver 10.32.0.10
options ndots:5
```

As you can see, the name server (DNS server) IP has been set to the IP address
of the kubernetes service of the core dns app, but our core dns app is not
working, due to some reason, which I don't know why. It just doesn't work
is how it looks like to me currently and I'm not able to check the logs. So,
I need to check what's going on with the logs. 

Later, I need to check if `kubeclt port-forward` works. That would be cool to
checkout ;)

One thing I do note is, how the tutorial says that the IP address of the
`kubernetes` service is `10.32.0.1` with this `nslookup kubernetes` output

```
Server:    10.32.0.10
Address 1: 10.32.0.10 kube-dns.kube-system.svc.cluster.local

Name:      kubernetes
Address 1: 10.32.0.1 kubernetes.default.svc.cluster.local
```

Hmm. Lot of things to checkout I guess.

Also, looks like the next step is
https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/13-smoke-test.md

Which clearly does some testing of the cluster to make sure everything is
running and the final thing is cleaning up, which is pretty easy for me, given
I just need to delete all the compute instances and that's it, nothing else to
do and I won't lose money even if they run forever 😅 I mean, yeah, electricity,
true, but once my computer shuts down, the instances should just shutdown and
that's it :) they won't start until I start it off

Anyways, I'm going to try some of the smoke tests, even though I know some of
them will fail for sure :P 

Oh wow. Till now my smoke tests have been working!!!! :D :D

I'm able to port forward, I'm able to create deployments, and my secret data
is encrypted at rest apparently, and I'm able to look at logs!!!!!! :D 

So, this means that my coredns, it just doesn't work. My logging is working
fine! 

I saw the logs like this for nginx pod

```bash
$ kubectl logs $POD_NAME
127.0.0.1 - - [01/Jun/2020:15:00:24 +0000] "GET / HTTP/1.1" 200 612 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:77.0) Gecko/20100101 Firefox/77.0" "-"
2020/06/01 15:00:25 [error] 6#6: *1 open() "/usr/share/nginx/html/favicon.ico" failed (2: No such file or directory), client: 127.0.0.1, server: localhost, request: "GET /favicon.ico HTTP/1.1", host: "localhost:8080"
127.0.0.1 - - [01/Jun/2020:15:00:25 +0000] "GET /favicon.ico HTTP/1.1" 404 154 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:77.0) Gecko/20100101 Firefox/77.0" "-"
127.0.0.1 - - [01/Jun/2020:15:00:32 +0000] "HEAD / HTTP/1.1" 200 0 "-" "curl/7.64.1" "-"
127.0.0.1 - - [01/Jun/2020:15:00:43 +0000] "GET / HTTP/1.1" 200 612 "-" "curl/7.64.1" "-"
127.0.0.1 - - [01/Jun/2020:15:00:46 +0000] "HEAD / HTTP/1.1" 200 0 "-" "curl/7.64.1" "-"
```

So, I gotta see why core dns wouldn't work. Why it has so many logs giving
the error about the certificate and why it uses the IP that I saw it was using
and what's going on.

All my smoke tests passed actually!! :D :D 

I think the next steps are - to make the core dns work and see if the
networking among pods work - it shouldn't ideally, since I didn't do anything
with respect to networking - the routes and all. Let's see

---

Okay, today I started my laptop and started the machines and exec pod does not
work. I was getting the same old error regarding worker host name resolution
like yesterday. On checking I noticed that `/etc/hosts` file changes have been
removed. Checking the comments I realized it was a managed file and to
persist changes I make to it, I need to more stuff

For now, I'm changing all the `/etc/cloud/templates/hosts.debian.tmpl` template
files in all the controller instances to include the worker instance DNS
entries (IPs)

After changing it in all machines using tmux synchronize panes feature, I saved
it and restarted the machines using 

```bash
$ sudo shutdown -r now
```

The status always says `Starting` in `multipass ls` for some reason. But once I
type in `multipass shell <machine-name>` then it opens shell and when I get out
and check status, I can see it changes the status to `Running`

And cool, now the the DNS entries persist! ;) :D 

About the `kubernetes` service having IP `10.0.0.1`, I backed it up and deleted
it, and now it's good! The coredns pod is giving a different error now though.

```bash
$ kubectl logs -f coredns-74576b4776-7tddb -n kube-system
.:53
2020-06-02T02:26:14.589Z [INFO] plugin/reload: Running configuration MD5 = fbb756dad13bce75afc40db627b38529
2020-06-02T02:26:14.589Z [INFO] CoreDNS-1.6.2
2020-06-02T02:26:14.589Z [INFO] linux/amd64, go1.12.8, 795a3eb
CoreDNS-1.6.2
linux/amd64, go1.12.8, 795a3eb
2020-06-02T02:26:14.595Z [ERROR] plugin/errors: 2 9081000483047803871.5836840357274749367. HINFO: plugin/loop: no next plugin found
2020-06-02T02:29:18.565Z [ERROR] plugin/errors: 2 google.com. AAAA: plugin/loop: no next plugin found
2020-06-02T02:29:32.973Z [ERROR] plugin/errors: 2 duckduckgo.com. AAAA: plugin/loop: no next plugin found
```

Some kind of plugin issue? Idk, I'm going to check about it.

So, I found this issue
https://github.com/coredns/coredns/issues/2166

I'm going to add a proxy for now. I tried adding proxy and that failed with
an error saying `proxy` is an unknown / undefined thing.

I later realized there's no plugin or stuff like that now, in the latest
version I think it's called `forward` - https://coredns.io/plugins/forward/ .
I also checked basically what `loop` is - https://coredns.io/plugins/loop/ .
`loop` looks like a good thing, it's supposed to help find loops in DNS queries
I think? and I think the information in this page can help me. I need to read it
and understand it better. Also, I tried some stuff. I added `forward` plugin
like this `forward . 8.8.8.8` and now coredns does not give any loop errors.

And I did a lot of meddling in the pods, to finally see some results which don't
make sense yet. This is what is the final thing that I saw:

```bash
$ kubectl run --generator=run-pod/v1 utils --image=arunvelsriram/utils -n default --command -- sleep 36000
$ execpod -a

 kubectl exec --namespace='default' utils -c utils -it sh

kubectl exec [POD] [COMMAND] is DEPRECATED and will be removed in a future version. Use kubectl kubectl exec [POD] -- [COMMAND] instead.
# bash
root@utils:/# cat /etc/resolv.conf
nameserver 8.8.8.8

search default.svc.cluster.local svc.cluster.local cluster.local
nameserver 10.32.0.10
options ndots:5

root@utils:/# dig @10.32.0.10 nginx.default.svc.cluster.local
;; reply from unexpected source: 10.200.2.10#53, expected 10.32.0.10#53
;; reply from unexpected source: 10.200.2.10#53, expected 10.32.0.10#53
;; reply from unexpected source: 10.200.2.10#53, expected 10.32.0.10#53

; <<>> DiG 9.11.3-1ubuntu1.12-Ubuntu <<>> @10.32.0.10 nginx.default.svc.cluster.local
; (1 server found)
;; global options: +cmd
;; connection timed out; no servers could be reached

root@utils:/# dig @10.32.0.10 kubernetes.default.svc.cluster.local

; <<>> DiG 9.11.3-1ubuntu1.12-Ubuntu <<>> @10.32.0.10 kubernetes.default.svc.cluster.local
; (1 server found)
;; global options: +cmd
;; connection timed out; no servers could be reached
root@utils:/#
root@utils:/# dig @10.32.0.10 kubernetes.default.svc.cluster.local
;; reply from unexpected source: 10.200.2.10#53, expected 10.32.0.10#53
;; reply from unexpected source: 10.200.2.10#53, expected 10.32.0.10#53
;; reply from unexpected source: 10.200.2.10#53, expected 10.32.0.10#53

; <<>> DiG 9.11.3-1ubuntu1.12-Ubuntu <<>> @10.32.0.10 kubernetes.default.svc.cluster.local
; (1 server found)
;; global options: +cmd
;; connection timed out; no servers could be reached

root@utils:/# dig @10.200.2.10 nginx.default.svc.cluster.local

; <<>> DiG 9.11.3-1ubuntu1.12-Ubuntu <<>> @10.200.2.10 nginx.default.svc.cluster.local
; (1 server found)
;; global options: +cmd
;; Got answer:
;; WARNING: .local is reserved for Multicast DNS
;; You are currently testing what happens when an mDNS query is leaked to DNS
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 29700
;; flags: qr aa rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1
;; WARNING: recursion requested but not available

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 4096
; COOKIE: de682c09c92b5c5b (echoed)
;; QUESTION SECTION:
;nginx.default.svc.cluster.local. IN    A

;; ANSWER SECTION:
nginx.default.svc.cluster.local. 2 IN   A       10.32.0.56

;; Query time: 1 msec
;; SERVER: 10.200.2.10#53(10.200.2.10)
;; WHEN: Tue Jun 02 06:10:46 UTC 2020
;; MSG SIZE  rcvd: 119

root@utils:/# dig @10.200.2.10 kubernetes.default.svc.cluster.local

; <<>> DiG 9.11.3-1ubuntu1.12-Ubuntu <<>> @10.200.2.10 kubernetes.default.svc.cluster.local
; (1 server found)
;; global options: +cmd
;; Got answer:
;; WARNING: .local is reserved for Multicast DNS
;; You are currently testing what happens when an mDNS query is leaked to DNS
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 45717
;; flags: qr aa rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1
;; WARNING: recursion requested but not available

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 4096
; COOKIE: 0db7c915d1ca79bc (echoed)
;; QUESTION SECTION:
;kubernetes.default.svc.cluster.local. IN A

;; ANSWER SECTION:
kubernetes.default.svc.cluster.local. 5 IN A    10.32.0.1

;; Query time: 1 msec
;; SERVER: 10.200.2.10#53(10.200.2.10)
;; WHEN: Tue Jun 02 06:17:38 UTC 2020
;; MSG SIZE  rcvd: 129

root@utils:/# dig @10.200.2.10 kubernetes.svc.cluster.local

; <<>> DiG 9.11.3-1ubuntu1.12-Ubuntu <<>> @10.200.2.10 kubernetes.svc.cluster.local
; (1 server found)
;; global options: +cmd
;; Got answer:
;; WARNING: .local is reserved for Multicast DNS
;; You are currently testing what happens when an mDNS query is leaked to DNS
;; ->>HEADER<<- opcode: QUERY, status: NXDOMAIN, id: 11906
;; flags: qr aa rd; QUERY: 1, ANSWER: 0, AUTHORITY: 1, ADDITIONAL: 1
;; WARNING: recursion requested but not available

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 4096
; COOKIE: ec1b93cf41607609 (echoed)
;; QUESTION SECTION:
;kubernetes.svc.cluster.local.  IN      A

;; AUTHORITY SECTION:
cluster.local.          5       IN      SOA     ns.dns.cluster.local. hostmaster.cluster.local. 1591067821 7200 1800 86400 5

;; Query time: 2 msec
;; SERVER: 10.200.2.10#53(10.200.2.10)
;; WHEN: Tue Jun 02 06:17:51 UTC 2020
;; MSG SIZE  rcvd: 162

```

If you notice, I have actually added `8.8.8.8` name server in the pod as the
first thing. I realized that nothing was working in the pod in terms of
reaching out to the Internet, for example to fetch the repository index for
all the packages through `apt update`, as the coredns server was not working
for some reason and `nslookup` was always stuck and I had to meddle with
`/etc/resolv.conf` to get something to work. I don't know how `/etc/resolv.conf`
exactly works, so that's something to lookout for. Ideally I wouldn't wanna
change anything in it and still would want everything to work just as is. About
that, I was debugging stuff.

I installed `dig` using `apt update; apt install dnsutils`

I started checking the results of `dig` for kubernetes service. The first thing
that I saw was it wasn't giving any proper replies - that is, no replies, no
IPs. I think this was because I had set the name server to be `8.8.8.8` in the
`/etc/resolv.conf`, so I started using the `@` to mention the name server to use
and I noticed it said that it was getting a reply from an unexpected source as
you see above.

So, I started using the unexpected source mentioned in the output as the
nameserver and I got the IP addresses of the kubernetes services. But yeah,
only for Fully Qualified Domain Names (FQDNs) I think, as you can see above.
It didn't work for `kubernetes` or `nginx` which works in applications for
some reason, but not in `dig`. 

I need to answer quite some questions now actually! 🙈

To start with, yes, the kubernetes service `kube-dns` with `10.32.0.10` is
supposed be the name server that's supposed to be working and helping with
the name resolution with respect to kubernetes services and the pods behind
the kubernetes service are the coredns pods.

But what's working now is, the name server at IP `10.200.2.10`. Now, this IP
is part of the `10.200.2.0/24` network, which is the network in the `worker-2`
instance and you can see the `worker-2` instances having the first IP as one
of it's node IPs

```bash
ubuntu@worker-2:~$ hostname -I
192.168.64.34 10.200.2.1

ubuntu@worker-2:~$ ifconfig cnio0
cnio0: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
        inet 10.200.2.1  netmask 255.255.255.0  broadcast 10.200.2.255
        inet6 fe80::40b8:95ff:fe92:925b  prefixlen 64  scopeid 0x20<link>
        ether 42:b8:95:92:92:5b  txqueuelen 1000  (Ethernet)
        RX packets 64998  bytes 5078205 (5.0 MB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 52938  bytes 25455996 (25.4 MB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0
```

and it's the pod CIDR range in the `worker-2` instance which we setup when
configuring the CNI (Container Networking Interface) with the bridge
configuration.

Okay, I just noticed this

```bash
$ k get pod -A -o wide
NAMESPACE     NAME                      READY   STATUS    RESTARTS   AGE     IP            NODE       NOMINATED NODE   READINESS GATES
default       busybox                   1/1     Running   6          17h     10.200.1.5    worker-1   <none>           <none>
default       nginx-554b9c67f9-f684l    1/1     Running   1          15h     10.200.0.4    worker-0   <none>           <none>
default       utils                     1/1     Running   0          3h8m    10.200.2.11   worker-2   <none>           <none>
kube-system   coredns-b4dc6cf5c-8xlzf   1/1     Running   0          3h17m   10.200.0.6    worker-0   <none>           <none>
kube-system   coredns-b4dc6cf5c-hbvsd   1/1     Running   0          3h17m   10.200.2.10   worker-2   <none>           <none>
```

So, one of the coredns pods have the IP `10.200.2.10`. Right! So, we need to see
what's going on and why the `dig` didn't work with service IP but worked with
pod IP. Hmm. Also, let me also simply try another node pod IP, in this case,
from utils pod in `worker-2` to `coredns-b4dc6cf5c-8xlzf` in `worker-0`.
Yeah, it just gets stuck, this is because of the routing rules I believe.

About the `dig` issue with service IP, I think I'm going to have to check
what IP forwarding is now. Seeing this is a matter of some loadbalancing by
the kubernetes service, which is not working, and the `dig` telling stuff
about getting `reply from unexpected source` etc, I want to see if the thing
I missed out before, that is, delegated for later - IP forwarding, is biting
me just now.

Need to learn what's IP forwarding first! 😅

And I read this article to see some `dig` stuff
https://linuxize.com/post/how-to-use-dig-command-to-query-dns-in-linux/

Now I read some stuff about how to do IP forwarding. I still am yet to learn
what it exactly does.

https://www.ducea.com/2006/08/01/how-to-enable-ip-forwarding-in-linux/

Looking at this, it said, set `net.bridge.bridge-nf-call-iptables` to `1`
https://github.com/mmumshad/kubernetes-the-hard-way/blob/master/docs/02-compute-resources.md

And this helped too

https://github.com/rak8s/rak8s/issues/13

I finally did this

```bash
$ sudo modprobe br_netfilter
$ vi /etc/sysctl.conf
...
# add this line
net.bridge.bridge-nf-call-iptables=1
$  sysctl net.bridge.bridge-nf-call-iptables
net.bridge.bridge-nf-call-iptables = 1
```

and it should say the value of `net.bridge.bridge-nf-call-iptables` is `1`.

Anyways, that's not working for the `dig` with `10.32.0.10` as dns server IP.

I need to understand how to enable verbose mode for `dig` and understand what's
happening behind the scenes, and also read more about IP forwarding to see
if it's related. I think checking how services connect to pods will help too,
since in this case, pod IP works for `dig`, service IP doesn't!

Okay, so I read a bit about how Kubernetes Services work! :) Over here
https://kubernetes.io/docs/concepts/services-networking/service/

And I also saw how the kube proxy helps in the working of the Services, and
that the Cluster IP of the services are virtual IPs.

https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies

And the kube proxy I have installed, it is running in `iptables` proxy mode,
which is one among the other proxy modes, to help one pod to connect to another
pod using the Kubernetes service

https://kubernetes.io/docs/concepts/services-networking/service/#proxy-mode-iptables

In my case, I went and checked the kube-proxy in the node where my pod was
running and it had some error logs. Something about TLS handshake issues,
timeout or something. I simply restarted it to see if the error persists and
it didn't and the `dig` just worked, and I removed the `8.8.8.8` entry too.
But `dig` worked only sometimes, not all times.

Now, this was because one of the two coredns pods is in another node, so it will
have a different kind of pod IP and I haven't put any routing rules or config
for routing from one node to another and hence the issues, sometimes, but when
the coredns pod in the same node was accessed, it worked ;) :)

Next thing to do is, add routing to the worker instances to make sure that
the pods in one node can access pods in another node

I'm reading about configuring routing in linux using `ip`

https://www.cyberciti.biz/faq/howto-linux-configuring-default-route-with-ipcommand/

before this I assumed that `iptables` is the command to look at and was
checking this

https://www.digitalocean.com/community/tutorials/how-to-list-and-delete-iptables-firewall-rules

but `iptables` stuff is only going to help in understanding about kube-proxy,
which I'll come back to and check later!

one more `iptables` link - https://www.cyberciti.biz/faq/how-to-list-all-iptables-rules-in-linux/ :)

back to `ip` now

Another related article for `ip` / `route` is
https://www.cyberciti.biz/tips/configuring-static-routes-in-debian-or-red-hat-linux-systems.html

Apparently changes done `ip` or `route` command don't persist and I need to
put the config in some file probably

For now, I'm just testing temporarily with `ip` command

The pod CIDR range and node IP in each worker instance is

```
worker-0:
pod CIDR - 10.200.0.0/24
node IP  - 192.168.64.32

worker-1:
pod CIDR - 10.200.1.0/24
node IP  - 192.168.64.33

worker-2:
pod CIDR - 10.200.2.0/24
node IP  - 192.168.64.34
```

```bash
ubuntu@worker-2:~$ sudo ip route add 10.200.1.0/24 via 192.168.64.33 dev enp0s2
ubuntu@worker-2:~$ sudo ip route add 10.200.0.0/24 via 192.168.64.32 dev enp0s2

ubuntu@worker-2:~$ ping 10.200.1.1
PING 10.200.1.1 (10.200.1.1) 56(84) bytes of data.
64 bytes from 10.200.1.1: icmp_seq=1 ttl=64 time=0.542 ms
64 bytes from 10.200.1.1: icmp_seq=2 ttl=64 time=0.378 ms
64 bytes from 10.200.1.1: icmp_seq=3 ttl=64 time=0.670 ms
^C
--- 10.200.1.1 ping statistics ---
3 packets transmitted, 3 received, 0% packet loss, time 2001ms
rtt min/avg/max/mdev = 0.378/0.530/0.670/0.119 ms

ubuntu@worker-2:~$ ping 10.200.0.1
PING 10.200.0.1 (10.200.0.1) 56(84) bytes of data.
64 bytes from 10.200.0.1: icmp_seq=1 ttl=64 time=0.437 ms
64 bytes from 10.200.0.1: icmp_seq=2 ttl=64 time=0.700 ms
64 bytes from 10.200.0.1: icmp_seq=3 ttl=64 time=0.492 ms
^C
--- 10.200.0.1 ping statistics ---
3 packets transmitted, 3 received, 0% packet loss, time 2085ms
rtt min/avg/max/mdev = 0.437/0.543/0.700/0.113 ms
```

```bash
ubuntu@worker-0:~$ sudo ip route add 10.200.1.0/24 via 192.168.64.33 dev enp0s2
ubuntu@worker-0:~$ sudo ip route add 10.200.2.0/24 via 192.168.64.34 dev enp0s2

ubuntu@worker-0:~$ ping 10.200.1.1
PING 10.200.1.1 (10.200.1.1) 56(84) bytes of data.
64 bytes from 10.200.1.1: icmp_seq=1 ttl=64 time=0.466 ms
64 bytes from 10.200.1.1: icmp_seq=2 ttl=64 time=0.704 ms
64 bytes from 10.200.1.1: icmp_seq=3 ttl=64 time=0.344 ms
^C
--- 10.200.1.1 ping statistics ---
3 packets transmitted, 3 received, 0% packet loss, time 2028ms
rtt min/avg/max/mdev = 0.344/0.504/0.704/0.151 ms

ubuntu@worker-0:~$ ping 10.200.2.1
PING 10.200.2.1 (10.200.2.1) 56(84) bytes of data.
64 bytes from 10.200.2.1: icmp_seq=1 ttl=64 time=1.02 ms
64 bytes from 10.200.2.1: icmp_seq=2 ttl=64 time=1.24 ms
64 bytes from 10.200.2.1: icmp_seq=3 ttl=64 time=0.657 ms
^C
--- 10.200.2.1 ping statistics ---
3 packets transmitted, 3 received, 0% packet loss, time 2066ms
rtt min/avg/max/mdev = 0.657/0.974/1.247/0.245 ms
```

```bash
ubuntu@worker-1:~$ sudo ip route add 10.200.0.0/24 via 192.168.64.32 dev enp0s2
ubuntu@worker-1:~$ sudo ip route add 10.200.2.0/24 via 192.168.64.34 dev enp0s2

ubuntu@worker-1:~$ ping 10.200.0.1
PING 10.200.0.1 (10.200.0.1) 56(84) bytes of data.
64 bytes from 10.200.0.1: icmp_seq=1 ttl=64 time=0.506 ms
64 bytes from 10.200.0.1: icmp_seq=2 ttl=64 time=0.705 ms
64 bytes from 10.200.0.1: icmp_seq=3 ttl=64 time=0.681 ms
^C
--- 10.200.0.1 ping statistics ---
3 packets transmitted, 3 received, 0% packet loss, time 2133ms
rtt min/avg/max/mdev = 0.506/0.630/0.705/0.093 ms

ubuntu@worker-1:~$ ping 10.200.2.1
PING 10.200.2.1 (10.200.2.1) 56(84) bytes of data.
64 bytes from 10.200.2.1: icmp_seq=1 ttl=64 time=0.561 ms
64 bytes from 10.200.2.1: icmp_seq=2 ttl=64 time=0.644 ms
64 bytes from 10.200.2.1: icmp_seq=3 ttl=64 time=0.626 ms
^C
--- 10.200.2.1 ping statistics ---
3 packets transmitted, 3 received, 0% packet loss, time 2145ms
rtt min/avg/max/mdev = 0.561/0.610/0.644/0.040 ms
```

Okay! Awesome! Now I see this

```bash
$ kubectl get pod -A -o wide
NAMESPACE     NAME                      READY   STATUS    RESTARTS   AGE   IP            NODE       NOMINATED NODE   READINESS GATES
default       busybox                   1/1     Running   13         26h   10.200.1.5    worker-1   <none>           <none>
default       nginx-554b9c67f9-f684l    1/1     Running   1          24h   10.200.0.4    worker-0   <none>           <none>
default       utils                     1/1     Running   0          12h   10.200.2.11   worker-2   <none>           <none>
kube-system   coredns-b4dc6cf5c-8xlzf   1/1     Running   0          12h   10.200.0.6    worker-0   <none>           <none>
kube-system   coredns-b4dc6cf5c-hbvsd   1/1     Running   0          12h   10.200.2.10   worker-2   <none>           <none>
```

and from inside the utils pod, I can do this

```bash
$ execpod -a

 kubectl exec --namespace='default' utils -c utils -it sh

kubectl exec [POD] [COMMAND] is DEPRECATED and will be removed in a future version. Use kubectl kubectl exec [POD] -- [COMMAND] instead.
# bash
root@utils:/# nslookup kubernetes
Server:         10.32.0.10
Address:        10.32.0.10#53

Name:   kubernetes.default.svc.cluster.local
Address: 10.32.0.1

root@utils:/# nslookup nginx
Server:         10.32.0.10
Address:        10.32.0.10#53

Name:   nginx.default.svc.cluster.local
Address: 10.32.0.56

root@utils:/# curl nginx
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
    body {
        width: 35em;
        margin: 0 auto;
        font-family: Tahoma, Verdana, Arial, sans-serif;
    }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>

root@utils:/# curl -k https://kubernetes/version
{
  "major": "1",
  "minor": "15",
  "gitVersion": "v1.15.3",
  "gitCommit": "2d3c76f9091b6bec110a5e63777c332469e0cba2",
  "gitTreeState": "clean",
  "buildDate": "2019-08-19T11:05:50Z",
  "goVersion": "go1.12.9",
  "compiler": "gc",
  "platform": "linux/amd64"
}
```

In between accessing `kubernetes` service was not working, because I made
mistakes in the `curl` command 🤦‍♂ like no `https` and then putting `6443`
as the port number when none was required as service was exposed at `443`
and targeted `6443` behind it on it's own. Anyways, everything works now!

I guess I need to check how to retain the IP routes and also I think the
IP forwarding stuff that I did, may not have been needed. I mean, according
to this

https://www.ducea.com/2006/08/01/how-to-enable-ip-forwarding-in-linux/

the linux servers I have, have IP forwarding enabled for IPv4. I'll check it
out by disabling everything that I did and also may be disabling this by
default enabled thing and see how it causes issues and also understand what
it means to do IP forwarding!

And I guess I need to write down all the learnings in a crisp manner some
where ;)
