# Containers and VMs

The below are just what I know, use it only with verification!

---

Container. Many say it in a simple manner that it's kind of like a light weight VM. "Kind of like a". What is it actually? It's an isolated process.

What's an isolated process? How is it different from a Virtual Machine?

An isolated process is a single process which has an isolated view and control over the host machine. For example, this isolated process cannot see the host file system, cannot check the hostname of the host machine, cannot view the processes of the host machine and can only view processes that it created. These are just some examples, to show that there's file system level isolation, network level isolation, process level isolation. You can also isolate the resources used by the isolated process - CPU, Memory (RAM). 

How is it different from a VM? In a VM, a complete OS runs, that is, there are a lot of processes running to boot up the OS and run it. Many processes run behind the scenes, for the working of the OS in the VM. Unlike the case of a container, an isolated process, where it's the only process and yes, it can create more processes if it wants. The isolated process just runs on top of the host machine OS, but in an isolated manner. In the case of a VM, there's a host machine OS, there's a separate guest OS, that is, the OS that runs in the VM. VMs are also slow usually, since they have a lot of components (processes) to run for their startup, unlike isolated processes, which is just one single process usually and it just runs really fast, almost no startup time really required.

Containers running in a host machine share the same kernel, the kernel of the host machine. So, if there were a kernel exploit, that allowed you to escape the container (the isolation), IF there were a kernel exploit, it's a shared kernel, so there's potentially a risk that says "If I can escape the container, I can see the other containers on that machine". It's much harder to escape a VM. So, for pretty specific reasons, you might say "I want to run my containers in separate VMs", a really good application for that is multi tenancy, like, you have groups of people, for who you know, you are running code on their behalf, and they don't trust each other. So, putting each of them into their own VMs makes it more secure for each of those tenants.

How does one create such an isolated process? You can create it using some features in Linux Kernels - namepaces, cgroups - control groups and chroot - change root. So, does that mean that you are constrained to run isolated processes only in Linux? We can't do it in other OSes like Mac, Windows? Answer is, I don't know actually. It's a good exercise for you to check for similar features in Windows and Mac Kernels to natively support creating isolated processes. 

It's still possible to run isolated processes, that is, containers in Mac and Windows in one way that I know of. You can run a Linux VM in Mac or Windows and then run isolated processes in it. I know it may sound weird, but yeah, that's one solution. Also, many cloud service providers provide VMs at a cheap cost, compared to bare metal (physical) machines. So, when you run containers in the cloud, you are already running it in VMs usually, unless you got dedicated bare metal machines.

Concept is called Containerization. Very old concept.

Container Images are similar to Virtual Machine Images.

Other terms that might be interesting - Virtualization (VMs), Micro VMs, Unikernels, Micro Kernels

