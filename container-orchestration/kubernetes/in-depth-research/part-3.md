# Part 3 - Worker Node

Now, for the worker node, we need quite a few things.

Let's start with the most basic thing maybe? The kubelet. We will then move on
to the container runtime, container networking interface, kube-proxy

We won't have to work on networking for pod across multiple nodes as we will
only have a single node cluster for now which has all components - control
plan and worker node components

Now, let's get started with the kubelet
