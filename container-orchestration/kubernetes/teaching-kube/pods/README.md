# Pods

Kubernetes is a platform. Platform to run things, run your applications.

How does one run things in Kubernetes? Kubernetes only has support for
containers. For example Docker containers. So, how does on run containers in
Kubernetes?

You need to create a Kubernetes object to tell Kubernetes to run containers.
So, what kind of object is this? Pod. But wait, why is it called pod and
not a container? We are trying to run containers here right?

The name pod comes from - pod of whales , pea pod 

https://www.whalefacts.org/what-is-a-group-of-whales-called/

https://en.wikipedia.org/wiki/Pea

Whale kind of corresponds to the Docker container. Docker's mascot is a cartoon
whale!

So, as you can see, pod refers to a group of things. What things in the case of
kubernetes? Pod refers to a group of containers?

Does that mean we can never run a single container in Kubernetes? No! We can.
It's just that, pod can represent one or multiple containers.

What does it mean to run a pod in kubernetes? When you tell kubernetes to create
a pod, you tell how many containers are in the pod. For simplicity let's say
it has just one container and then you tell the container's image, and then
give information like name of the pod, name of the container and then the
command to run in the container. And there's actually more that you can define,
but let's just keep it simple for now.

So, what happens then? after we tell kubernetes about the pod that it needs to
create. This is what happens - Kubernetes creates the pod, and schedules it to a node, among the
cluster of nodes. Meaning, it gives the pod (kind of like a task) to a node, and
the node has to run the  pod (that is, do the task). Running the pod translates
to - running the containers in the pod using the container runtime/engine,
like Docker.

So, that's how you run things in Kubernetes. Using a pod, which in turn is about
running containers run in a node.

Why not directly tell Kubernetes about the containers to run? Why another
abstraction called pods on top of containers? Why the new term? Like we spoke
before, pod can contain one or **more** containers. Usually people just run one,
but it can still support **more** containers, and for those cases, the group
has been defined as a pod. There's nothing smaller than a pod. Pod is the
smallest unit of work.

So, why and when would someone run multiple containers in a pod? What does it
even mean to run multiple containers in a pod? Pod is just a logical thing.
Container is the actual thing that's running processes (applications). 
Kubernetes assigns pods (tasks) to a node, and it does it only at pod level,
because pod is the smallest unit of work. So, when a pod is assigned to a node,
it means that all the containers in the pod need to be executed in the node.
One would run multiple containers in a pod when they want to co-locate a
certain set of containers, meaning they want a certain set of containers to
always run together, in the same node, and not on separate nodes. This could be
due to multiple reasons. For example, intra node container network communication
might be faster than inter node container network communication. Another
reason to run multiple containers in a pod could be to scale a set of containers
proportionately. 

Also, usually, when a pod contains multiple containers, the best practise is to
have one of the containers as the main container and the other containers
augment the main container's functionality or provide some feature on the side.
The containers other than the main containers are called side cars. It's like
the side car in the vehicles.

Some examples for sidecars are - light weight proxy running as sidecar for
backend services and scaling along with the backend services, config reloader
sidecar which reloads the configuration file in the main container dynamically
by fetching it from a remote source and the configuration file is used by the
main container's process.


