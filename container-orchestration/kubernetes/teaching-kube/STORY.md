# Story

November 29th 2020

Tasks:

- Have a copy of main slides [DONE]
- Read through the slides
- Ask questions on each slide
- Think about how to introduce the concepts to the audience in the format of a
  story
- Keep asking "Why that's needed?", "Why more ways to do stuff?", "Why this
  extra feature?", "Why are there more ways to do something?" and more
  questions
- Think about the why for each of the features and introduce that as a problem
  to the audience and then introduce the solution that Kubernetes has to offer
  and also mention any external solutions
- Have kubernetes yaml examples for each of the resources (Jobs, CronJobs,
  DaemonSets, StatefulSets)
  - Prepare for the demo and write down the features to show during the demo

---

Jobs:

- To run something and finish it to completion
- Run and forget
- Want logs of old jobs? or the latest job?
- Alternatives?
  - Serverless functions like AWS Lambda
- Retry / Rerun job if it fails due to some reason? Max retries?
- How does it run? It still uses Pods to run the job
- Use case?
  - CI/CD pipelines
  - Data pipelines which need to run some data processing jobs
  - Running database migrations

CronJobs

- It's a abstraction on top of Jobs. It also runs something and finishes it
  to completion
- Run something based on a cron schedule, similar to cron and crontab in linux
- How does it run? It creates a kubernetes Job based on the cron schedule.
  The Job in turn runs a pod. So, it's again using Pods only to run the work.
  Everything use Pods to run stuff in Kubernetes.
- How to get notifications about the job succeeding or failing?
  - You can use custom code / tools to send notifications to yourself. For
    example slack notifications, Google Chat notifications, using API Keys
    and incoming webhooks
  - You can also rely on monitoring and altering infrastructure like Prometheus
    and Alertmanager and write rules and configuration to get notifications when
    jobs succeed and fail. Alertmanager has integrations for multiple services!
    :)
- Use case?
  - Run nightly jobs to do some processing
  - Run recurring backup jobs every day once, or every month once or similar

DaemonSets

- Run some long running process in all the nodes in the cluster
- Use Case?
  - To run some node specific work in all the nodes
  - Collect all the logs of all the containers running in all the nodes in the
    cluster and ship it to some remote log aggregation service for storage,
    processing and analysis
  - Get the IP addresses of all the pods in each of the nodes to use for service
    discovery. Run agents / processes on all nodes to do this. This would mean
    that you are trying to use an alternative to Kubernetes Service. For
    example, to use Consul service discovery.
- How to do it? You can manually create pods. Usually pods can run on any node.
  Also, more than one pod can run on a single node. But you want to run exactly
  one pod on each of the nodes of the cluster. You can manually assign the pods
  you create to the particular nodes. You can also try to run a deployment or
  replicaset. Ensure there's pod affinity so that only one instance of the
  process runs in one node. What about when the cluster is scaled down? That is,
  the number of nodes becomes less. What about scaling up the cluster? Your
  deployment / replicaset has to have a dynamic replica count. There's something
  called Horizontal Pod Autoscaler. This is a lot of work. Instead, Kubernetes
  provides a simple solution to this. It's called a DaemonSet. :)

---

DaemonSet

https://duckduckgo.com/?t=ffab&q=kubernetes+scheduling&ia=web

https://kubernetes.io/docs/reference/scheduling/

https://kubernetes.io/docs/concepts/scheduling-eviction/kube-scheduler/

https://docs.fluentd.org

https://docs.fluentd.org/installation

https://docs.fluentd.org/container-deployment

https://docs.fluentd.org/container-deployment/kubernetes

https://carbon.now.sh/?bg=rgba%2874%2C144%2C226%2C1%29&t=material&wt=none&l=auto&ds=true&dsyoff=20px&dsblur=68px&wc=true&wa=false&pv=56px&ph=56px&ln=false&fl=1&fm=Fira+Code&fs=18px&lh=152%25&si=false&es=4x&wm=false&code=apiVersion%253A%2520v1%250Akind%253A%2520Pod%250Ametadata%253A%250A%2520%2520name%253A%2520simple-task%250A%2520%2520labels%253A%250A%2520%2520%2520%2520app%253A%2520simple-task%250Aspec%253A%250A%2520%2520containers%253A%250A%2520%2520%2520%2520-%2520name%253A%2520echo-task%250A%2520%2520%2520%2520%2520%2520image%253A%2520busybox%250A%2520%2520%2520%2520%2520%2520command%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520-%2520%2522echo%2522%250A%2520%2520%2520%2520%2520%2520args%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520-%2520%2522ok%2522%250A

https://duckduckgo.com/?t=ffab&q=fluentd+elasticsearch&ia=web

https://docs.fluentd.org/output/elasticsearch

https://www.fluentd.org/guides/recipes/elasticsearch-and-s3

https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#inter-pod-affinity-and-anti-affinity


