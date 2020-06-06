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

