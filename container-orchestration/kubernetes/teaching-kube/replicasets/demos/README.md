# Demos

Let's see how to create a simple replicaset yaml ;)

```yaml
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: dobby
spec:
  replicas: 3
  selector:
    matchLabels:
      app: dobby
  template:
    metadata:
      name: dobby
      labels:
        app: dobby
    spec: 
      containers:
      - name: dobby
        image: thecasualcoder/dobby
```

Let's get rid of the running pods!

```bash
$ kubectl get pod
NAME      READY   STATUS    RESTARTS   AGE
dobby-1   1/1     Running   0          29m
dobby-2   1/1     Running   0          29m
$ kubectl delete pod dobby-1 dobby-2
pod "dobby-1" deleted
pod "dobby-2" deleted
```

And now let's apply our replicaset! :D

```bash
$ kubectl apply -f simple-replicaset.yaml
replicaset.apps/dobby created

$ kubectl get pod
NAME          READY   STATUS    RESTARTS   AGE
dobby-mqfwh   1/1     Running   0          7s
dobby-q4qj5   1/1     Running   0          7s
dobby-vhhdc   1/1     Running   0          7s
```

Before going into the details of what's going on and the yaml file contents,
I also tried some changes in the yaml.

Like changing the `matchLabels` field in the `spec` and the `labels` field in
the `template` like the below to include another random label

```yaml
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: dobby
spec:
  replicas: 3
  selector:
    matchLabels:
      app: dobby
      appVersion: "1"
  template:
    metadata:
      name: dobby
      labels:
        app: dobby
        appVersion: "1"
    spec: 
      containers:
      - name: dobby
        image: thecasualcoder/dobby
```

And I got this error

```bash
$ kubectl apply -f simple-replicaset.yaml
The ReplicaSet "dobby" is invalid: spec.selector: Invalid value: v1.LabelSelector{MatchLabels:map[string]string{"app":"dobby", "appVersion":"1"}, MatchExpressions:[]v1.LabelSelectorRequirement(nil)}: field is immutable
```

It says that the `selector` field is immutable.

Then I changed it to this, to include the new label only in the `template`

```yaml
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: dobby
spec:
  replicas: 3
  selector:
    matchLabels:
      app: dobby
  template:
    metadata:
      name: dobby
      labels:
        app: dobby
        appVersion: "1"
    spec: 
      containers:
      - name: dobby
        image: thecasualcoder/dobby
```

And it worked!

```bash
$ kubectl apply -f simple-replicaset.yaml
replicaset.apps/dobby configured
```

But there was no changes in the pods

```bash
$ kubectl get pod --show-labels
NAME          READY   STATUS    RESTARTS   AGE    LABELS
dobby-mqfwh   1/1     Running   0          116s   app=dobby
dobby-q4qj5   1/1     Running   0          116s   app=dobby
dobby-vhhdc   1/1     Running   0          116s   app=dobby
```

The new label hadn't come in the pod! And now, to simply meddle with things, I
delete one of the pods! ;)

```bash
$ kubectl delete pod dobby-mqfwh
pod "dobby-mqfwh" deleted
$ kubectl get pod --show-labels
NAME          READY   STATUS    RESTARTS   AGE     LABELS
dobby-q4qj5   1/1     Running   0          2m48s   app=dobby
dobby-spw5v   1/1     Running   0          23s     app=dobby,appVersion=1
dobby-vhhdc   1/1     Running   0          2m48s   app=dobby
```

If you notice, a new pod has come up now, and the new pod has the new label too!
What's really going on??

