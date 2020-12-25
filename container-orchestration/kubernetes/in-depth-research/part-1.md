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
