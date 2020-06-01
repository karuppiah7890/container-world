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

