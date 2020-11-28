# Demos

- Create a simple statefulset containing an nginx with a volume
- Create a simple service to point to the different nginx instances
- Deploy utils pod to run some commands to see if the nginx works

```bash
$ kubectl run --generator=run-pod/v1 utils --image=arunvelsriram/utils -n default --command -- sleep 36000
```

- Create html file(s) in each of the pods at `/usr/share/nginx/html/`. At least
  `index.html`, with pod name in the html document

- Run `curl` and hit the pod IP of the nginx pods

```bash
$ curl 172.17.0.3
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Nginx 0</title>
</head>
<body>
    Nginx 0
</body>
</html>
```

- Run `nslookup` or `dig` against the kubernetes service

```bash
$ nslookup nginx.default.svc.cluster.local
Server:         10.96.0.10
Address:        10.96.0.10#53

Name:   nginx.default.svc.cluster.local
Address: 172.17.0.5
Name:   nginx.default.svc.cluster.local
Address: 172.17.0.4
Name:   nginx.default.svc.cluster.local
Address: 172.17.0.3
```

```bash
$ dig nginx.default.svc.cluster.local

; <<>> DiG 9.11.3-1ubuntu1.11-Ubuntu <<>> nginx.default.svc.cluster.local
;; global options: +cmd
;; Got answer:
;; WARNING: .local is reserved for Multicast DNS
;; You are currently testing what happens when an mDNS query is leaked to DNS
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 7195
;; flags: qr aa rd; QUERY: 1, ANSWER: 3, AUTHORITY: 0, ADDITIONAL: 1
;; WARNING: recursion requested but not available

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 4096
; COOKIE: 49829580637a1d5e (echoed)
;; QUESTION SECTION:
;nginx.default.svc.cluster.local. IN    A

;; ANSWER SECTION:
nginx.default.svc.cluster.local. 28 IN  A       172.17.0.4
nginx.default.svc.cluster.local. 28 IN  A       172.17.0.3
nginx.default.svc.cluster.local. 28 IN  A       172.17.0.5

;; Query time: 0 msec
;; SERVER: 10.96.0.10#53(10.96.0.10)
;; WHEN: Sat Nov 28 15:06:26 UTC 2020
;; MSG SIZE  rcvd: 213

$ dig +short nginx.default.svc.cluster.local
172.17.0.5
172.17.0.3
172.17.0.4

```

```bash
$ nslookup web-0.nginx.default.svc.cluster.local
Server:         10.96.0.10
Address:        10.96.0.10#53

Name:   web-0.nginx.default.svc.cluster.local
Address: 172.17.0.3

$ nslookup web-1.nginx.default.svc.cluster.local
Server:         10.96.0.10
Address:        10.96.0.10#53

Name:   web-1.nginx.default.svc.cluster.local
Address: 172.17.0.4

$ nslookup web-2.nginx.default.svc.cluster.local
Server:         10.96.0.10
Address:        10.96.0.10#53

Name:   web-2.nginx.default.svc.cluster.local
Address: 172.17.0.5
```

- Delete pods and see how it comes back and hit the pod again to see what the
  content inside the pod is.
  - Check the volume attached to it.

