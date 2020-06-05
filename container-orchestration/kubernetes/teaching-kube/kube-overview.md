# Kubernetes Overview

I think many of you might be running servers. 

Where do you deploy your servers? VMs? Containers? Cloud? Serverless?

How do you deploy it? Like, any tools that you use to deploy? Ansible, Chef, Puppet, Terraform?

You could be deploying to Kubernetes too. So, what's the difference?

Kubernetes is a platform for you to run your applications. But how do you run
your applications? We will come to that. But in short, Kubernetes also runs
applications on VMs, but in containers. Kubernetes only support running
containers. What's really different with Kubernetes though? To answer that,
let's talk a bit about Kubernetes.

Kubernetes has a master worker architecture, where there are master nodes and
worker nodes, each running some components. The master nodes assign work / work
loads to worker nodes. This work load is your application! The complete set of
nodes - master and worker together, is called a cluster - a kubernetes cluster.
It simply denotes that there are multiple nodes or how you say, there is a
cluster of nodes.

The assignment of workload / work to a worker node is called scheduling.
Scheduling happens based on some requirements. For example, let’s say your
application needs some X amount of RAM and y amount of CPU, you mention that
to the kubernetes master and then the kubernetes master makes sure that it
assigns this application (workload) to run on a worker node that has enough
resources based on the requirements, so that your application can run well. If
it cannot find a worker node to run the application, then the application
(workload) will be stuck and will not be executed. Usually if you are running in
the cloud, for such cases, there are auto scaling mechanisms to scale up the
number of worker nodes when the resources available in the current number of
worker nodes are not enough. It's called node auto scaling.

Now, why do people even choose Kubernetes? Kubernetes as a platform, with the
concept of scheduling has good cost savings. Since it tries to run more than one
application in a worker node (VM) and tries to utilize the resources in all the
nodes in an optimal manner, it saves a lot of cost, when compared to running
each application instance in one VM, where it’s very natural to waste some
available resources in the all the VMs put together.

Kubernetes supports containers, which also means your applications are running
in an isolated manner in containers and cannot possibly affect other
applications running in another container in the same node it’s running on. So,
all applications in the worker nodes will run smoothly.

Kubernetes also follows some good concepts based on years of experience in
infrastructure - different deployment strategies, health checks to check the
health of the services (liveness probe feature), traffic routing check to check
when network traffic can be sent to service as the service needs to be ready
before it recieves a request (readiness probe feature).

We spoke a bit about Kubernetes and stuff. Let's move on [this page for the
next topic](./pods/README.md)
