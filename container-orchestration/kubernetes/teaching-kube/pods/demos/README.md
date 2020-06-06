# Demos

We can start with a simple pod yaml.


```yaml
# simple-pod.yaml
apiVersion: v1
kind: Pod
metadata:
  name: dobby
spec:
  containers:
  - name: dobby
    image: thecasualcoder/dobby
```

We are running the `dobby` application in this pod.
You can find `dobby`'s source code here - https://github.com/thecasualcoder/dobby/

To create this pod, do this

```bash
$ kubectl apply -f simple-pod.yaml
pod/dobby created
```

To check the pods that are running, you can use this

```
$ kubectl get pod -o wide
NAME                      READY   STATUS    RESTARTS   AGE    IP            NODE       NOMINATED NODE   READINESS GATES
dobby                     1/1     Running   0          63s    10.200.1.15   worker-1   <none>           <none>
```

With `-o wide` you can also see the IP address of the pod, also the node in which
it's running.

To check the pod's logs, use this

```bash
$ kubectl logs -f dobby
[GIN-debug] [WARNING] Now Gin requires Go 1.6 or later and Go 1.7 will be required soon.

[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /health                   --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).Health-fm (3 handlers)
[GIN-debug] GET    /readiness                --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).Ready-fm (3 handlers)
[GIN-debug] GET    /version                  --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).Version-fm (3 handlers)
[GIN-debug] GET    /meta                     --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).Meta-fm (3 handlers)
[GIN-debug] PUT    /control/health/perfect   --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).MakeHealthPerfect-fm (3 handlers)
[GIN-debug] PUT    /control/health/sick      --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).MakeHealthSick-fm (3 handlers)
[GIN-debug] PUT    /control/ready/perfect    --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).MakeReadyPerfect-fm (3 handlers)
[GIN-debug] PUT    /control/ready/sick       --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).MakeReadySick-fm (3 handlers)
[GIN-debug] PUT    /control/goturbo/memory   --> github.com/thecasualcoder/dobby/pkg/handler.GoTurboMemory (3 handlers)
[GIN-debug] PUT    /control/goturbo/cpu      --> github.com/thecasualcoder/dobby/pkg/handler.GoTurboCPU (3 handlers)
[GIN-debug] PUT    /control/crash            --> github.com/thecasualcoder/dobby/pkg/handler.Crash.func1 (3 handlers)
```

The `-f` is to follow / stream the logs.

To port foward the pod, to be able to access the pod's port from your machine,
use this

```bash
$ kubectl port-forward dobby 4444
Forwarding from 127.0.0.1:4444 -> 4444
Forwarding from [::1]:4444 -> 4444
```

It will keep running and blocking. So, you gotta use another terminal to check
if it works. Now in another terminal you can access this port! :)

```bash
$ curl -i localhost:4444
HTTP/1.1 404 Not Found
Content-Type: text/plain
Date: Sat, 06 Jun 2020 06:57:14 GMT
Content-Length: 18

404 page not found

$ curl -i localhost:4444/health
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Sat, 06 Jun 2020 06:57:15 GMT
Content-Length: 16

{"healthy":true}

$ curl -i localhost:4444/version
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Sat, 06 Jun 2020 06:57:55 GMT
Content-Length: 39

{"version":"1.0.v1.0-4-g7218d24-dirty"}
```

You can notice in the port forwarded terminal these logs

```bash
$ kubectl port-forward dobby 4444
Forwarding from 127.0.0.1:4444 -> 4444
Forwarding from [::1]:4444 -> 4444
Handling connection for 4444
Handling connection for 4444
Handling connection for 4444
Handling connection for 4444
Handling connection for 4444
Handling connection for 4444
Handling connection for 4444
```

Now, you can stop this port forwarding. Now, let's get inside the pod, that is,
execute into the pod. This is similar to `docker exec`, or in the VM world, ssh
into the VM and run commands

You can run any command inside the pod's container using `kubectl exec`, and
if you run a shell in an interactive manner with stdin (standard input) attached
to it and with `tty` option enabled, it's kind of like SSHing into a machine

```bash
$ kubectl exec dobby -- uname -a
Linux dobby 4.15.0-101-generic #102-Ubuntu SMP Mon May 11 10:07:26 UTC 2020 x86_64 x86_64 x86_64 GNU/Linux
```

Now let's run a shell, like `bash`, which is present in this container and get
into the shell too! ;)

```bash
$ kubectl exec dobby -it -- bash
root@dobby:/# ls
bin  boot  dev  etc  home  lib  lib64  media  mnt  opt  proc  root  run  sbin  srv  sys  tmp  usr  var

root@dobby:/# apt update
Get:1 http://archive.ubuntu.com/ubuntu bionic InRelease [242 kB]
Get:2 http://security.ubuntu.com/ubuntu bionic-security InRelease [88.7 kB]
Get:3 http://archive.ubuntu.com/ubuntu bionic-updates InRelease [88.7 kB]
Get:4 http://archive.ubuntu.com/ubuntu bionic-backports InRelease [74.6 kB]
Get:5 http://archive.ubuntu.com/ubuntu bionic/main amd64 Packages [1344 kB]
...
...

root@dobby:/# apt install curl
Reading package lists... Done
Building dependency tree
Reading state information... Done
The following additional packages will be installed:
  ca-certificates krb5-locales libasn1-8-heimdal libcurl4 libgssapi-krb5-2 libgssapi3-heimdal libhcrypto4-heimdal libheimbase1-heimdal
  libheimntlm0-heimdal libhx509-5-heimdal libk5crypto3 libkeyutils1 libkrb5-26-heimdal libkrb5-3 libkrb5support0 libldap-2.4-2 libldap-common
  libnghttp2-14 libpsl5 libroken18-heimdal librtmp1 libsasl2-2 libsasl2-modules libsasl2-modules-db libsqlite3-0 libssl1.1 libwind0-heimdal openssl
  publicsuffix
...
...
Setting up curl (7.58.0-2ubuntu3.8) ...
Processing triggers for libc-bin (2.27-3ubuntu1) ...
Processing triggers for ca-certificates (20190110~18.04.1) ...
Updating certificates in /etc/ssl/certs...
0 added, 0 removed; done.
Running hooks in /etc/ca-certificates/update.d...
done.

root@dobby:/# curl -i localhost:4444/health
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Sat, 06 Jun 2020 07:15:13 GMT
Content-Length: 16

{"healthy":true}

root@dobby:/# curl -i localhost:4444/version
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Sat, 06 Jun 2020 07:15:16 GMT
Content-Length: 39

{"version":"1.0.v1.0-4-g7218d24-dirty"}

root@dobby:/#
```

So, we got into the pod and tried to hit the API! :)

You can look at the logs again to see the request logs

```bash
$ kubectl logs -f dobby
[GIN-debug] [WARNING] Now Gin requires Go 1.6 or later and Go 1.7 will be required soon.

[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /health                   --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).Health-fm (3 handlers)
[GIN-debug] GET    /readiness                --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).Ready-fm (3 handlers)
[GIN-debug] GET    /version                  --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).Version-fm (3 handlers)
[GIN-debug] GET    /meta                     --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).Meta-fm (3 handlers)
[GIN-debug] PUT    /control/health/perfect   --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).MakeHealthPerfect-fm (3 handlers)
[GIN-debug] PUT    /control/health/sick      --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).MakeHealthSick-fm (3 handlers)
[GIN-debug] PUT    /control/ready/perfect    --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).MakeReadyPerfect-fm (3 handlers)
[GIN-debug] PUT    /control/ready/sick       --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).MakeReadySick-fm (3 handlers)
[GIN-debug] PUT    /control/goturbo/memory   --> github.com/thecasualcoder/dobby/pkg/handler.GoTurboMemory (3 handlers)
[GIN-debug] PUT    /control/goturbo/cpu      --> github.com/thecasualcoder/dobby/pkg/handler.GoTurboCPU (3 handlers)
[GIN-debug] PUT    /control/crash            --> github.com/thecasualcoder/dobby/pkg/handler.Crash.func1 (3 handlers)
[GIN] 2020/06/06 - 06:57:03 | 404 |      18.155µs |       127.0.0.1 | GET      /
[GIN] 2020/06/06 - 06:57:07 | 404 |      22.549µs |       127.0.0.1 | GET      /ping
[GIN] 2020/06/06 - 06:57:11 | 200 |     103.589µs |       127.0.0.1 | GET      /health
[GIN] 2020/06/06 - 06:57:14 | 404 |      23.056µs |       127.0.0.1 | GET      /
[GIN] 2020/06/06 - 06:57:15 | 200 |     153.315µs |       127.0.0.1 | GET      /health
[GIN] 2020/06/06 - 06:57:53 | 200 |      94.314µs |       127.0.0.1 | GET      /version
[GIN] 2020/06/06 - 06:57:55 | 200 |     200.024µs |       127.0.0.1 | GET      /version
[GIN] 2020/06/06 - 07:15:09 | 200 |     281.659µs |             ::1 | GET      /health
[GIN] 2020/06/06 - 07:15:13 | 200 |      83.631µs |             ::1 | GET      /health
[GIN] 2020/06/06 - 07:15:16 | 200 |     235.315µs |             ::1 | GET      /version
[GIN] 2020/06/06 - 07:15:58 | 404 |      24.186µs |             ::1 | GET      /help
```

How about we crash the pod using the `/control/crash` endpoint? ;) Let's do it!

I'm port forwarding port `4444` and doing a http request with `curl`

```bash
$ kubectl port-forward dobby 4444

Forwarding from 127.0.0.1:4444 -> 4444
Forwarding from [::1]:4444 -> 4444
```

In another terminal, this

```bash
$ curl -X PUT localhost:4444/control/crash
curl: (52) Empty reply from server

$ kubectl get pod dummy
NAME                      READY   STATUS    RESTARTS   AGE
dobby                     0/1     Error     1          5h34m

$ kubectl get pod dummy
NAME                      READY   STATUS    RESTARTS   AGE
dobby                     0/1     Error     1          5h34m

$ kubectl get pod dummy
NAME                      READY   STATUS    RESTARTS   AGE
dobby                     0/1     Error     1          5h34m
```

You can see the restarts and I did a lot of them. I checked the logs to see
how it looks, and this is what I saw -

```
$ kubectl logs -f dobby
[GIN-debug] [WARNING] Now Gin requires Go 1.6 or later and Go 1.7 will be required soon.

[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /health                   --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).Health-fm (3 handlers)
[GIN-debug] GET    /readiness                --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).Ready-fm (3 handlers)
[GIN-debug] GET    /version                  --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).Version-fm (3 handlers)
[GIN-debug] GET    /meta                     --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).Meta-fm (3 handlers)
[GIN-debug] PUT    /control/health/perfect   --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).MakeHealthPerfect-fm (3 handlers)
[GIN-debug] PUT    /control/health/sick      --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).MakeHealthSick-fm (3 handlers)
[GIN-debug] PUT    /control/ready/perfect    --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).MakeReadyPerfect-fm (3 handlers)
[GIN-debug] PUT    /control/ready/sick       --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).MakeReadySick-fm (3 handlers)
[GIN-debug] PUT    /control/goturbo/memory   --> github.com/thecasualcoder/dobby/pkg/handler.GoTurboMemory (3 handlers)
[GIN-debug] PUT    /control/goturbo/cpu      --> github.com/thecasualcoder/dobby/pkg/handler.GoTurboCPU (3 handlers)
[GIN-debug] PUT    /control/crash            --> github.com/thecasualcoder/dobby/pkg/handler.Crash.func1 (3 handlers)
you asked me do so, killing myself :-)

```

Lol. That's some dark humor :P `you asked me do so, killing myself :-)`

I started doing some more crashes to see it like this

```bash
$ kubectl get pod dobby
NAME                      READY   STATUS             RESTARTS   AGE
dobby                     0/1     CrashLoopBackOff   2          5h38m
```

It finally looks like this now

```bash
$ kubectl get pod dobby
NAME                      READY   STATUS    RESTARTS   AGE
dobby                     1/1     Running   3          5h38m
```

And port forwarded terminal looks like this!

```bash
$ kubectl port-forward dobby 4444

Forwarding from 127.0.0.1:4444 -> 4444
Forwarding from [::1]:4444 -> 4444
Handling connection for 4444
Handling connection for 4444
Handling connection for 4444
Handling connection for 4444
E0606 17:58:49.372414   56854 portforward.go:400] an error occurred forwarding 4444 -> 4444: error forwarding port 4444 to pod 38ae2cb647ec928df0b12786e2494d8f3d69eb82fb11cf3c36a12522ba1cd1d3, uid : failed to execute portforward in network namespace "/var/run/netns/cni-a6f6f922-faeb-2b2f-8204-2f566da7cf8b": socat command returns error: exit status 1, stderr: "2020/06/06 17:58:49 socat[26785] E connect(5, AF=2 127.0.0.1:4444, 16): Connection refused\n"
```

There were connection refuses because the container had crashed 🤷‍♂

Anyways, we can actually do more things, like use liveness and readiness probe
and play with it. But for now, we are just going to move on. Also, by the way,
dobby server has other cool endpoints, like have a spike in memory, cpu ;)

Now, how about we run two containers inside the same pod?? ;)

This is to see how the containers in the pod share the kernel namespace. Can
one container see the processes running in the other containers within the
same pod. Let's find out! I'm going to run an nginx container along side
dobby :)

This is the pod yaml that I'm writing

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: dobby-with-nginx
spec:
  containers:
  - name: dobby
    image: thecasualcoder/dobby
  - name: nginx
    image: nginx
```

Just a simple pod yaml with an extra container running `nginx` server.

Okay, I see the pod now

```bash
$ kubectl get pod
NAME                      READY   STATUS    RESTARTS   AGE
dobby                     1/1     Running   3          5h50m
dobby-with-nginx          2/2     Running   0          43s
```

It took quite some time for the `nginx` image to be pulled, that it was in
`ContainerCreating` state for a long time and the `kubectl describe pod dobby-with-nginx`
showed in events that it was pull the nginx image.

Anyways, I wanted to check the logs of the pod, and I see this!

```bash
$ kubectl logs -f dobby-with-nginx
error: a container name must be specified for pod dobby-with-nginx, choose one of: [dobby nginx]
```

As you can see, when there are more than one container in a pod, you have to
specifically mention which container's logs you want. Last time it was only one
container, so no such issues or errors or questions! :)

And this makes sense too. Like, in docker, you see logs for one container at a
time. But hey, `docker-compose` could show the logs of all containers in an
interleaved manner. Anyways. Now, let me check logs one by one now for each
container

```
$ kubectl logs -f dobby-with-nginx -c dobby
[GIN-debug] [WARNING] Now Gin requires Go 1.6 or later and Go 1.7 will be required soon.

[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /health                   --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).Health-fm (3 handlers)
[GIN-debug] GET    /readiness                --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).Ready-fm (3 handlers)
[GIN-debug] GET    /version                  --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).Version-fm (3 handlers)
[GIN-debug] GET    /meta                     --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).Meta-fm (3 handlers)
[GIN-debug] PUT    /control/health/perfect   --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).MakeHealthPerfect-fm (3 handlers)
[GIN-debug] PUT    /control/health/sick      --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).MakeHealthSick-fm (3 handlers)
[GIN-debug] PUT    /control/ready/perfect    --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).MakeReadyPerfect-fm (3 handlers)
[GIN-debug] PUT    /control/ready/sick       --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).MakeReadySick-fm (3 handlers)
[GIN-debug] PUT    /control/goturbo/memory   --> github.com/thecasualcoder/dobby/pkg/handler.GoTurboMemory (3 handlers)
[GIN-debug] PUT    /control/goturbo/cpu      --> github.com/thecasualcoder/dobby/pkg/handler.GoTurboCPU (3 handlers)
[GIN-debug] PUT    /control/crash            --> github.com/thecasualcoder/dobby/pkg/handler.Crash.func1 (3 handlers)
^C

$ kubectl logs -f dobby-with-nginx -c nginx
/docker-entrypoint.sh: /docker-entrypoint.d/ is not empty, will attempt to perform configuration
/docker-entrypoint.sh: Looking for shell scripts in /docker-entrypoint.d/
/docker-entrypoint.sh: Launching /docker-entrypoint.d/10-listen-on-ipv6-by-default.sh
10-listen-on-ipv6-by-default.sh: Getting the checksum of /etc/nginx/conf.d/default.conf
10-listen-on-ipv6-by-default.sh: Enabled listen on IPv6 in /etc/nginx/conf.d/default.conf
/docker-entrypoint.sh: Launching /docker-entrypoint.d/20-envsubst-on-templates.sh
/docker-entrypoint.sh: Configuration complete; ready for start up
```

So, the containers are good! :) Btw, the you could also use `stern` command to
check the logs at pod level, for all containers! ;) Or even for multiple pods
by just using a regex along with `stern` which matches the pod name or names ;)
:D In this case, it would look like this!

```bash
$ stern dobby-with-nginx
+ dobby-with-nginx › dobby
+ dobby-with-nginx › nginx
dobby-with-nginx nginx /docker-entrypoint.sh: /docker-entrypoint.d/ is not empty, will attempt to perform configuration
dobby-with-nginx nginx /docker-entrypoint.sh: Looking for shell scripts in /docker-entrypoint.d/
dobby-with-nginx nginx /docker-entrypoint.sh: Launching /docker-entrypoint.d/10-listen-on-ipv6-by-default.sh
dobby-with-nginx nginx 10-listen-on-ipv6-by-default.sh: Getting the checksum of /etc/nginx/conf.d/default.conf
dobby-with-nginx nginx 10-listen-on-ipv6-by-default.sh: Enabled listen on IPv6 in /etc/nginx/conf.d/default.conf
dobby-with-nginx nginx /docker-entrypoint.sh: Launching /docker-entrypoint.d/20-envsubst-on-templates.sh
dobby-with-nginx nginx /docker-entrypoint.sh: Configuration complete; ready for start up
dobby-with-nginx dobby [GIN-debug] [WARNING] Now Gin requires Go 1.6 or later and Go 1.7 will be required soon.
dobby-with-nginx dobby
dobby-with-nginx dobby [GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.
dobby-with-nginx dobby
dobby-with-nginx dobby [GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
dobby-with-nginx dobby  - using env:    export GIN_MODE=release
dobby-with-nginx dobby  - using code:   gin.SetMode(gin.ReleaseMode)
dobby-with-nginx dobby
dobby-with-nginx dobby [GIN-debug] GET    /health                   --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).Health-fm (3 handlers)
dobby-with-nginx dobby [GIN-debug] GET    /readiness                --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).Ready-fm (3 handlers)
dobby-with-nginx dobby [GIN-debug] GET    /version                  --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).Version-fm (3 handlers)
dobby-with-nginx dobby [GIN-debug] GET    /meta                     --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).Meta-fm (3 handlers)
dobby-with-nginx dobby [GIN-debug] PUT    /control/health/perfect   --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).MakeHealthPerfect-fm (3 handlers)
dobby-with-nginx dobby [GIN-debug] PUT    /control/health/sick      --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).MakeHealthSick-fm (3 handlers)
dobby-with-nginx dobby [GIN-debug] PUT    /control/ready/perfect    --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).MakeReadyPerfect-fm (3 handlers)
dobby-with-nginx dobby [GIN-debug] PUT    /control/ready/sick       --> github.com/thecasualcoder/dobby/pkg/handler.(*Handler).MakeReadySick-fm (3 handlers)
dobby-with-nginx dobby [GIN-debug] PUT    /control/goturbo/memory   --> github.com/thecasualcoder/dobby/pkg/handler.GoTurboMemory (3 handlers)
dobby-with-nginx dobby [GIN-debug] PUT    /control/goturbo/cpu      --> github.com/thecasualcoder/dobby/pkg/handler.GoTurboCPU (3 handlers)
dobby-with-nginx dobby [GIN-debug] PUT    /control/crash            --> github.com/thecasualcoder/dobby/pkg/handler.Crash.func1 (3 handlers)
^C
```

Notice how it says the pod name first and then the container name and shows
the logs ? ;) :D :) 

Now, let's check inside the containers and see what and all are shared! :)

I got into both the containers and checked if I could see the other process

```bash
$ kubectl exec dobby-with-nginx -c dobby -it -- sh

# bash
root@dobby-with-nginx:/# ps aux
USER       PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
root         1  0.0  0.8 413648  8996 ?        Ssl  12:40   0:00 dobby server --bind-address 0.0.0.0
root        10  0.5  0.0   4628   864 pts/0    Ss   12:48   0:00 sh
root        15  1.0  0.3  18508  3460 pts/0    S    12:48   0:00 bash
root        20  0.0  0.2  34400  2976 pts/0    R+   12:48   0:00 ps aux
root@dobby-with-nginx:/#
```

```bash
$ kubectl exec dobby-with-nginx -c nginx -it -- sh

# bash
root@dobby-with-nginx:/# ps aux
bash: ps: command not found
root@dobby-with-nginx:/#
```

So, I realized that both were proper containers - isolated processes, with the
file system, including `/proc` being not shared among the two. `/proc` is what
has all the data regarding the processes and is a virtual file system.

In both the containers, I was able to use `curl` and hit both `localhost` for
`nginx` and `localhost:4444` for `dobby`. The network space is shared in the
sense that `localhost` refers to the pod, and reaches both the containers and
for specificity, there's port, and the port space is shared too, so, if one
container uses one port, the other containers in the same pod cannot reuse that
port!

I haven't tried stuff related to volumes. But yes, pods can have volumes, which
can be shared by all the containers in the pod, and you can mount the volumes
at any paths in each container. There's a field called `volumes` in the pod spec
and another called `volumeMounts` in the containers array elements

Now, let's try to run multiple instances of our application!! :D :D But with
just `dobby`. I'm going to get rid of `nginx` container. So, back to the
simple pod yaml

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: dobby
spec:
  containers:
  - name: dobby
    image: thecasualcoder/dobby
```

Now, to scale our application, we run more pods. Another thing you can do is,
run more containers in the same pod. But you Must not run multiple containers
in the same pod like that, for scaling. 

People run multiple containers in the same pod for specific advanced use cases,
like, when they want all the containers to be colocated and present in the same
node - actually even that can be satisfied differently. But the rule of the
thumb is this - when there are multiple containers in the pod, there must be
exactly one container that's the main container - like the main application,
and the other containers must just help the main container. Check the
[pods doc](../README.md) to know more.

Also, soon we will see that Kubernetes understands how to scale only at pod
level and not at container level.

Now, let's scale our application by running more pods! How do we do that?
So, the thing is, each pod gets a unique name, and cannot be reused. 
So, we need to create more pods, running `dobby`, but with different names.
I was thinking more like `dobby-1`, `dobby-2`, `dobby-3`. Totally 3 instances!
;)

I'm going to keep changing the name in `simple-pod.yaml` and deploy it! ;)

```
$ kubectl delete pod dobby dobby-with-nginx
pod "dobby" deleted
pod "dobby-with-nginx" deleted

$ vi simple-pod.yaml
$ kubectl apply -f simple-pod.yaml
pod/dobby-1 created

$ vi simple-pod.yaml
$ kubectl apply -f simple-pod.yaml
pod/dobby-2 created

$ vi simple-pod.yaml
$ kubectl apply -f simple-pod.yaml
pod/dobby-3 created
```

And we have 3 instances of `dobby` now! ;)

```
$ kubectl get pod
NAME      READY   STATUS    RESTARTS   AGE
dobby-1   1/1     Running   0          42s
dobby-2   1/1     Running   0          37s
dobby-3   1/1     Running   0          32s
```

Now, let's say I want to change the metadata in the pods to add labels. Labels
help you with labelling your resources for fetching them easily later. It's
kind of like tagging things like articles or book marking web pages. Now, how
would I add labels to all the pods? I need to now change the
`simple-pod.yaml` to add the metadata and keep changing the name and apply it
for each pod! Phew!


```bash
$ vi simple-pod.yaml
$ kubectl apply -f simple-pod.yaml
pod/dobby-1 configured

$ vi simple-pod.yaml
$ kubectl apply -f simple-pod.yaml
pod/dobby-2 configured

$ vi simple-pod.yaml
$ kubectl apply -f simple-pod.yaml
pod/dobby-3 configured
```

Now, to check the pod yaml's metadata field, I used this complex thing :P

```bash
$ kubectl get pod -o yaml | yq '.items | map(.metadata | { name: .name, labels: .labels })'
[
  {
    "name": "dobby-1",
    "labels": {
      "app": "dobby"
    }
  },
  {
    "name": "dobby-2",
    "labels": {
      "app": "dobby"
    }
  },
  {
    "name": "dobby-3",
    "labels": {
      "app": "dobby"
    }
  }
]
```

That shows the pod names and their labels! 

Another, and easy way to do that would be (yeah, I forgot this one for a second)
is this -

```bash
$ kubectl get pod --show-labels
NAME      READY   STATUS    RESTARTS   AGE   LABELS
dobby-1   1/1     Running   0          10m   app=dobby
dobby-2   1/1     Running   0          10m   app=dobby
dobby-3   1/1     Running   0          10m   app=dobby
```

Now I want to reduce the number of instances to 2!

```bash
$ kubectl delete pod dobby-3
pod "dobby-3" deleted
```

Can you see what we are doing here? To create multiple instances, we used the
same single yaml file and kept changing the name metadata field and kept doing
`kubectl apply -f <yaml-file>` and it was quite cumbersome. And it was harder
to accomodatethe changes, do the whole thing again with the name, but make
the other changes too, and again use `kubectl apply`. 

And to add more instances or delete instances, I would have manually do
`kubectle apply` or `kubectl delete`

Now, in the infrastructure world, this is what people call as `imperative operations`.

You want to do something, for example in this case -
`scale up or scale down the number of instances` and what you did was, you told
the system imperatively on how to do it `add one more instance`, `delete this instance`
and so on.

Imperative meaning, you tell the how to the system. You tell how to do the thing.

There's another concept called `declarative operations`. Decalarative meaning,
you don't tell how to do something, you just say what to do to the system and
then the system takes care of doing it!

If you notice, for the label change, we just did `kubectl apply`. We didn't
tell the system `hey add labels` to this pod. Instead we gave the whole pod
yaml, and we said `apply` which means

```
this is my desired state (the pod yaml) and reach this state, somehow.
```

And `apply` will either create if the resource doesn't already exist or it will
update if the resource already exists. So, `apply` in itself is declarative,
compared to the command `kubectl create`, which is used to create a resource,
but if you use it again while updating, it will give an error saying resource
has already been created, so you can only use it for create, but not for update
and it's exactly one operation, an imperative one. Whereas `apply` understands
when to create or update, and does it accordingly, and hence it's declarative.
Also, it's idempotent - meaning, if you do it once or even hundred times, it
doesn't matter, it will always work! But you can't run `kubectl create` multiple
times with the same input, it will give error the second time itself, saying
resource already exists, if the first one successfully created the resource
of course.

So, similar to `create`, `delete` is also imperative. We don't want to be
deleting pods one by one or multiple of them, for scaling down. And it's not
idempotent.

In this case, a declarative and idempotent operation would be to say

```
hey kubernetes, just run x instances of my application
```

And it will take care of scaling up or scaling down, based on the current
state, to reach the desired state of `x instances`.

And kubernetes loves following decalartive configuration. It's main idea is
that. And to help with scaling, there's another resource! It's called a
replicaset!! :D :D
