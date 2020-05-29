# Thoughts

We are going to build a container based on three
concepts in Linux:
- Namespaces
- Chroot - change root
- Cgroups - control groups

What happens when we run a docker container?

```
$ docker run --rm -it ubuntu /bin/bash
root@5c72fd7d80f5:/#
```

We are able to run a container based on an image, run an arbitrary command
inside the container

We have hostname command

```
root@5c72fd7d80f5:/# hostname
5c72fd7d80f5
```

It's sort of the ID of the container. Some kind of random ID that's been
allocated

When we check the processes running using `ps` command inside the container

```
root@5c72fd7d80f5:/# ps
  PID TTY          TIME CMD
    1 pts/0    00:00:00 bash
   11 pts/0    00:00:00 ps
```

We can only see processes running inside the container. And the process IDs are
numbered starting from `1`

But from the host machine, if we check `ps` command

```
$ ps
  PID TTY           TIME CMD
  699 ttys000    0:00.14 -bash
  902 ttys000    0:00.01 tmux
 1382 ttys001    0:00.01 bash -c exec /bin/bash... 2> /dev/null & reattach-to-user-namespace -l /bin/bash
 1389 ttys001    0:01.40 -bash
38684 ttys001    0:00.08 docker run --rm -it ubuntu /bin/bash
51309 ttys003    0:00.00 bash -c exec /bin/bash... 2> /dev/null & reattach-to-user-namespace -l /bin/bash
51312 ttys003    0:00.17 -bash
```

We have got might higher process IDs

And even the hostname is different in the host machine

```
$ hostname
Karuppiah-N.local
```

We are going to try and recreate something like `docker run`

---

Namespaces

What you can **see**

Limit what a process can see

A container running, it could only see a few of the processes, which are
processes running inside the container, but it could not see processes in the
host. This is because of the namespace for process IDs

The container could only see it's own hostname. Again it's because of
namespacing.

We setup these namespaces using syscalls. Depending on the Linux kernel version
there are different number of syscalls present

Examples:
- Unix Timesharing System
- Process IDs
- Mounts
- Network
- User IDs
- InterProcess Comms

Namespace is a big part of what makes a container into a container. It's
restricting the view of the process, the view of the things that are going on in
the host machine.

Let's start building our dopple ganger for Docker 😉 starting with namespaces

---

We write a small and simple program to run an arbitrary command, given the
command/program name and it's arguments. For example, if you want to run

```
$ echo "ok" "wow"
```

`echo` would be the command name and `"ok"` and `"wow"` and arguments to it.

In golang, we can run this using the `exec` golang standard library and it's
`Cmd` struct.

After doing that, we can containerize this program now. We are going to start
with namespaces, for which we need to create namespaces.

You can apply the namespaces using the `SysProcAttr` field in the
`Cmd` struct

We use the `Cloneflags` field in `SysProcAttr`, which creates the new process we
want to run our arbitrary command in. And since this whole thing works only
in Linux, we have to run the program in Linux for us to be able to see it
working. 

I'm writing the code in VS Code editor in Mac and I was not able to see the
`Cloneflags` field in `SysProcAttr` as a suggestion or in the package, this
was because the `exec` package with which the suggest was given is a `darwin`
specific package.

Going to use `multipass` to run an `ubuntu` instance. I'll mount the source
code in it and then run it

```
$ multipass launch -n containers
$ multipass exec containers bash
```

Using [`gofish`](https://gofi.sh) to install `go`.

```bash
ubuntu@containers:~$ curl -fsSL https://raw.githubusercontent.com/fishworks/gofish/master/scripts/install.sh | bash

Downloading https://gofi.sh/releases/gofish-v0.11.0-linux-amd64.tar.gz
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100 6621k  100 6621k    0     0  1611k      0  0:00:04  0:00:04 --:--:-- 1611k
Preparing to install into /usr/local/bin
gofish installed into /usr/local/bin/gofish
Run 'gofish init' to get started!

ubuntu@containers:~$ gofish init
The following new directories will be created:
/usr/local/gofish
/usr/local/gofish/Barrel
/usr/local/gofish/Rigs
/usr/local/bin
/home/ubuntu/.gofish
==> Installing default fish food...
🐠  rig constructed in 3.032279488s

ubuntu@containers:~$ gofish install go
==> Installing go...
🐠  go 1.14.3: installed in 22.357973514s
```

Mounted the source code from host machine

```bash
$ multipass mount doppledocker/ containers:/doppledocker
```

```bash
ubuntu@containers:~$ cd /doppledocker/

ubuntu@containers:/doppledocker$ ls
README.md  doppledocker  go.mod  main.go
```

Awesome! :)

For getting suggestions in my VS Code editor based on linux `exec` standard
library, I put this settings for the workspace specific settings json file

```json
"go.toolsEnvVars": {
    "GOOS": "linux"
}
```

This is picked up by the Golang extension and the tools and I'm getting
proper suggestions now for `Cloneflags` field ! :)

And this is only possible because Golang has the code for all the platforms
that it supports - this is to do cross compilation for multiple platforms by
being on just one platform. So, using Mac, I can compile for Mac, Linux, Windows
and for different architectures toos! :) I usually primarily work with Mac or
Linux with `amd64` architecture

We are starting with the flag for `Unix Timesharing System`. It's not much, it's
just the namespace that contains only the `hostname`.

Currently, without any namespaces, this is what the program shows the `hostname`
as

```bash
ubuntu@containers:/doppledocker$ go run main.go run hostname
Running [hostname]
containers
```

As you can see, it's inheriting (?) the host machine's `hostname`, which is
`containers` and shows that when it runs. It knows about the host machine
`hostname`. The `syscall.CLONE_NEWUTS` is going to let us have our own `hostname`
inside the container, so it can see it's own `hostname` but can't see what's
happening on the host.

After inserting the code and runnning, it gives me this

```bash
ubuntu@containers:/doppledocker$ go run main.go run hostname
Running [hostname]
panic: fork/exec /bin/hostname: operation not permitted

goroutine 1 [running]:
main.run()
        /doppledocker/main.go:36 +0x255
main.main()
        /doppledocker/main.go:16 +0x4e
exit status 2
```

And then I used `sudo` and it worked

```bash
ubuntu@containers:/doppledocker$ sudo go run main.go run hostname
Running [hostname]
containers
```

But it still says the `hostname` as `containers` which is again the host machine
`hostname`!

And I can actually get into a shell inside the container and then keep using it!
Like how we do in docker containers! Like this

```bash
ubuntu@containers:/doppledocker$ go run main.go run /bin/bash
Running [/bin/bash]
panic: fork/exec /bin/bash: operation not permitted

goroutine 1 [running]:
main.run()
        /doppledocker/main.go:36 +0x255
main.main()
        /doppledocker/main.go:16 +0x4e
exit status 2
```

But it doesn't work without root. So, let's use `sudo`

```
ubuntu@containers:/doppledocker$ sudo go run main.go run /bin/bash
Running [/bin/bash]
root@containers:/doppledocker# hostname
containers
```

Notice how I got into a new shell? Usually you cannot see it as the username
is usually the same, but this time, since we used `sudo`, the username of the
shell shows as `root` and yes, I'm inside a shell inside a shell! :P and
`hostname` is still `containers`, like before. Like we saw, it was inherited
from the host machine

But we can change it, inside the container

```
root@containers:/doppledocker# hostname blah
root@containers:/doppledocker# hostname
blah
```

And in a separate tab, I check the `hostname` of the host machine

```
ubuntu@containers:~$ hostname
containers
```

and it has not been affected! :) But if you remove the code with the
`SysProcAttr` and `Cloneflags: syscall.CLONE_NEWUTS`, and run and try changing
`hostname` inside the container, it affects the host machine `hostname` too!

But as we can see, the prompt is not showing the `hostname` correctly. So, it's
hard to tell the `hostname` by checking the prompt and instead we need to run
the `hostname` command to see it. It's easier if the `hostname` is set up
before spawning the shell in the container, then it would show up properly in
the prompt. And it will also help us understand if we are inside the container
or not. 

Now, we can set the `hostname` using this

```golang
syscall.Sethostname([]byte{"dopplecontainer"})
```

But where do we call this method? If we call this after `cmd.Run()` then in our
case, it will run the arbitrary command, that is `bash` shell for us, and then
when we exit, then only the `hostname` will be set, but we want to set it such
that we can see the `hostname` in our `bash` prompt. 

If we set it before the `cmd.Run()`, the problem is, before running the command,
we are still in the same process, no new process has been created with the
namespaces that we have in our configuration, that is the `SysProcAttr` and the
values in it. We just have the configuration and it will be used and new process
will be created based on that only when the `cmd.Run()` runs. So, if we do it
before `cmd.Run()`, it will actually set the `hostname` for the host machine
and then that will be inherited by the container

So, what we are going to do for this is, make our process clone a new process
with a new namespace and then we are going to create another process in which
we are going to run our arbitrary command, like `bash`.

So, let's create a function called `child` now. Before doing that, let's check
how our processes list looks like from the view of the host machine, to see
what's going on behind scenes now and then see what happens when we do the above
mentioned changes. When running sleep

```bash
ubuntu@containers:/doppledocker$ sudo go run main.go run /bin/bash
Running [/bin/bash]
root@containers:/doppledocker# sleep 10
```

Below is what we see in the host machine

```bash
ubuntu@containers:~$ ps aux | grep bash
ubuntu    1496  0.0  0.5  23196  5104 pts/0    Ss   12:01   0:00 bash
ubuntu    3134  0.0  0.5  23124  5272 pts/1    Ss   13:09   0:00 bash
root      3481  0.0  0.4  63968  4216 pts/0    S    13:37   0:00 sudo go run main.go run /bin/bash
root      3482  0.0  1.8 1163096 18864 pts/0   Sl   13:37   0:00 go run main.go run /bin/bash
root      3503  0.0  0.1 703224  1252 pts/0    Sl   13:37   0:00 /tmp/go-build373887408/b001/exe/main run /bin/bash
root      3509  0.0  0.5  23000  5104 pts/0    S+   13:37   0:00 /bin/bash
ubuntu    3568  0.0  0.1  14856  1088 pts/1    S+   15:26   0:00 grep --color=auto bash
ubuntu@containers:~$ ps aux | grep sleep
root      3569  0.0  0.0   7924   756 pts/0    S+   15:26   0:00 sleep 10
ubuntu    3571  0.0  0.1  14856  1100 pts/1    S+   15:26   0:00 grep --color=auto sleep
```

You can see the `sudo` command as one process. `go` command as another, which
in turn builds the code into a temporary executable in `/tmp`, which we can
see above `/tmp/go-build373887408/b001/exe/main` and then we see the `bash`
command running which is the shell that we see!

Now, let's see how this changes with the changes in our code. :)

Now we have added another command called `child`, and we use `run` to create
a namespace and then run itself again, but using `child` command and the same
set of arguments. The `child` command just runs the command given the arguments
but it's running inside a new namespace, so if it sets a `hostname`, that will
not change the host machine `hostname`. This is how it looks!

Running the container

```bash
ubuntu@containers:/doppledocker$ sudo go run main.go run /bin/bash
Running [/bin/bash]
Running [/bin/bash]
root@dopplecontainer:/doppledocker# sleep 10
root@dopplecontainer:/doppledocker# hostname
dopplecontainer
```

And you can see how the `hostname` of the container is `dopplecontainer` in the
prompt!

And in the host machine

```
ubuntu@containers:~$ ps aux | grep bash
ubuntu    3134  0.0  0.5  23124  5272 pts/1    Ss   13:09   0:00 bash
ubuntu    3678  0.0  0.5  23132  5184 pts/0    Ss   15:38   0:00 bash
root      3725  0.0  0.4  63968  4228 pts/0    S    15:39   0:00 sudo go run main.go run /bin/bash
root      3726  0.8  1.8 1163096 18788 pts/0   Sl   15:39   0:00 go run main.go run /bin/bash
root      3751  0.0  0.1 703224  1248 pts/0    Sl   15:39   0:00 /tmp/go-build996572777/b001/exe/main run /bin/bash
root      3757  0.0  0.1 703224  1248 pts/0    Sl   15:39   0:00 /proc/self/exe child /bin/bash
root      3761  0.0  0.5  23000  5076 pts/0    S    15:39   0:00 /bin/bash
ubuntu    3774  0.0  0.1  14856  1060 pts/1    R+   15:39   0:00 grep --color=auto bash
ubuntu@containers:~$ ps aux | grep sleep
root      3772  0.0  0.0   7924   772 pts/0    S+   15:39   0:00 sleep 10
ubuntu    3776  0.0  0.1  14856  1052 pts/1    S+   15:39   0:00 grep --color=auto sleep

ubuntu@containers:~$ hostname
containers
```

As you can see, there are many more things going on this time! Actually just
one extra thing - like the `/proc/self/exe`. And the `hostname` of the host
machine has not changed and it's the same! :)

In the container, we can see the processes

```bash
root@dopplecontainer:/doppledocker# ps
  PID TTY          TIME CMD
 3725 pts/0    00:00:00 sudo
 3726 pts/0    00:00:00 go
 3751 pts/0    00:00:00 main
 3757 pts/0    00:00:00 exe
 3761 pts/0    00:00:00 bash
 3790 pts/0    00:00:00 ps
root@dopplecontainer:/doppledocker# ps aux | grep bash
ubuntu    3134  0.0  0.5  23124  5296 pts/1    Ss+  13:09   0:00 bash
ubuntu    3678  0.0  0.5  23132  5184 pts/0    Ss   15:38   0:00 bash
root      3725  0.0  0.4  63968  4228 pts/0    S    15:39   0:00 sudo go run main.go run /bin/bash
root      3726  0.0  1.8 1163096 18788 pts/0   Sl   15:39   0:00 go run main.go run /bin/bash
root      3751  0.0  0.1 703224  1248 pts/0    Sl   15:39   0:00 /tmp/go-build996572777/b001/exe/main run /bin/bash
root      3757  0.0  0.1 703224  1248 pts/0    Sl   15:39   0:00 /proc/self/exe child /bin/bash
root      3761  0.0  0.5  23132  5240 pts/0    S    15:39   0:00 /bin/bash
root      3796  0.0  0.1  14856  1148 pts/0    R+   15:46   0:00 grep --color=auto bash
```

In the host, we see the same

```bash
ubuntu@containers:~$ ps
  PID TTY          TIME CMD
 3134 pts/1    00:00:00 bash
 3800 pts/1    00:00:00 ps
ubuntu@containers:~$ ps aux | grep bash
ubuntu    3134  0.0  0.5  23124  5296 pts/1    Ss   13:09   0:00 bash
ubuntu    3678  0.0  0.5  23132  5184 pts/0    Ss   15:38   0:00 bash
root      3725  0.0  0.4  63968  4228 pts/0    S    15:39   0:00 sudo go run main.go run /bin/bash
root      3726  0.0  1.8 1163096 18788 pts/0   Sl   15:39   0:00 go run main.go run /bin/bash
root      3751  0.0  0.1 703224  1248 pts/0    Sl   15:39   0:00 /tmp/go-build996572777/b001/exe/main run /bin/bash
root      3757  0.0  0.1 703224  1248 pts/0    Sl   15:39   0:00 /proc/self/exe child /bin/bash
root      3761  0.0  0.5  23132  5240 pts/0    S+   15:39   0:00 /bin/bash
ubuntu    3794  0.0  0.1  14856  1148 pts/1    S+   15:46   0:00 grep --color=auto bash
```

We can see the same, high number processes in the container that we see in the
host. This is because we have not namespaced the processes (process IDs).

We want `ps` to just return the processes running inside the container. This
is also a part of isolating a process. The container cannot know what's running
inside the host machine.

We also want the process IDs of procsses inside the container to start with
process ID 1.

The namespace for this is called `syscall.CLONE_NEWPID`, where `PID` refers to
process ID.

Let's also print the process ID when running the commands to see the ID. We
can use `os.Getpid()` for this

Awesome, this is what see now!

```bash
ubuntu@containers:/doppledocker$ sudo go run main.go run /bin/bash
Running [/bin/bash] with process ID 3876
Running [/bin/bash] with process ID 1
root@dopplecontainer:/doppledocker#
```

But. We can still see host machine processes 🙈

```bash
root@dopplecontainer:/doppledocker# ps
  PID TTY          TIME CMD
 3850 pts/0    00:00:00 sudo
 3851 pts/0    00:00:00 go
 3876 pts/0    00:00:00 main
 3880 pts/0    00:00:00 exe
 3886 pts/0    00:00:00 bash
 3903 pts/0    00:00:00 ps
root@dopplecontainer:/doppledocker# ps aux | grep bash
ubuntu    3134  0.0  0.5  23124  5296 pts/1    Ss+  13:09   0:00 bash
ubuntu    3678  0.0  0.5  23132  5184 pts/0    Ss   15:38   0:00 bash
root      3850  0.0  0.4  63968  4256 pts/0    S    15:56   0:00 sudo go run main.go run /bin/bash
root      3851  0.1  1.8 1163096 18660 pts/0   Sl   15:56   0:00 go run main.go run /bin/bash
root      3876  0.0  0.1 703228  1248 pts/0    Sl   15:56   0:00 /tmp/go-build472554747/b001/exe/main run /bin/bash
root      3880  0.0  0.1 703228  1248 pts/0    Sl   15:56   0:00 /proc/self/exe child /bin/bash
root      3886  0.0  0.4  23000  4928 pts/0    S    15:56   0:00 /bin/bash
root      3902  0.0  0.0  14856  1000 pts/0    S+   15:58   0:00 grep --color=auto bash
root@dopplecontainer:/doppledocker#
```

This is because `ps` gets it's information from the `/proc` directory. And the
container has access to the whole of the host machine file system, including
`/proc`

```bash
root@dopplecontainer:/doppledocker# ls /
bin   dev           etc   initrd.img      lib    lost+found  mnt  proc  run   snap  sys  usr  vmlinuz
boot  doppledocker  home  initrd.img.old  lib64  media       opt  root  sbin  srv   tmp  var  vmlinuz.old
root@dopplecontainer:/doppledocker# ls /doppledocker/
README.md  go.mod  main.go
```

Inside the container, I can even see the source code I mounted in the host
machine!

`/proc` has information about all of the running processes. There's a directory
for each process and the name of the directory is the process ID.

Let's try to understand `/proc/self/exe` a bit now.

```bash
ubuntu@containers:~$ ls -l /proc/self/exe
lrwxrwxrwx 1 ubuntu ubuntu 0 May 29 16:04 /proc/self/exe -> /bin/ls
ubuntu@containers:~$ ls -l /proc/self/exe
lrwxrwxrwx 1 ubuntu ubuntu 0 May 29 16:04 /proc/self/exe -> /bin/ls
ubuntu@containers:~$ ls -l /proc/self/exe
lrwxrwxrwx 1 ubuntu ubuntu 0 May 29 16:04 /proc/self/exe -> /bin/ls
```

`/proc/self/exe` always points to the program that's accessing it, in this case
the `ls` command. Also, check this

```bash
ubuntu@containers:~$ ls -l /proc/self
lrwxrwxrwx 1 root root 0 May 29 08:12 /proc/self -> 3925
ubuntu@containers:~$ ls -l /proc/self
lrwxrwxrwx 1 root root 0 May 29 08:12 /proc/self -> 3926
ubuntu@containers:~$ ls -l /proc/self
lrwxrwxrwx 1 root root 0 May 29 08:12 /proc/self -> 3927
ubuntu@containers:~$ ls -l /proc/self
lrwxrwxrwx 1 root root 0 May 29 08:12 /proc/self -> 3928
```

Each time we run the `ls` command and check the `/proc/self`, we get that `ls`
command's process's directory which has the name as process ID of the command.

And there's all sorts of interesting information in the `/proc` for each of
the processes

When we are inside the container, running `ps`, it's using the `/proc` which is
actually the host machine's `/proc`. So, we need our own version of `/proc`.

This is where the `chroot` comes in. We are going to change the root of what
the container can see.
