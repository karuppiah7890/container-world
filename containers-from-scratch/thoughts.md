# Thoughts And Notes

The below are based on the slides and contents of the
[YouTube talk](https://www.youtube.com/watch?v=8fi7uSYlOdc) - a lot of times the
exact same thing Liz Rice explained, and some of my own thoughts too in between,
while trying out things

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

So, I need a file system now. I could take up any OS's file system at any stage.
A better thing would be to get the clean slate file system of a machine where
the OS was just installed. I tried to check how to get it using `docker` and
found this method -

```bash
$ docker pull alpine
$ docker save alpine -o alpine.tar
$ # tldr to the rescue to find out
$ # how to use tar command!
$ tldr tar
$ mkdir alpine
$ tar xvf alpine.tar -C alpine
x 485d7306187faf0cc9b77fc210e8def9d67b1953f0669e675877e90e6542cb6d/
x 485d7306187faf0cc9b77fc210e8def9d67b1953f0669e675877e90e6542cb6d/VERSION
x 485d7306187faf0cc9b77fc210e8def9d67b1953f0669e675877e90e6542cb6d/json
x 485d7306187faf0cc9b77fc210e8def9d67b1953f0669e675877e90e6542cb6d/layer.tar
x a187dde48cd289ac374ad8539930628314bc581a481cdb41409c9289419ddb72.json
x manifest.json
x repositories
$ cd alpine
$ ls
485d7306187faf0cc9b77fc210e8def9d67b1953f0669e675877e90e6542cb6d
a187dde48cd289ac374ad8539930628314bc581a481cdb41409c9289419ddb72.json
manifest.json
repositories
$ rm -rfv repositories manifest.json a187dde48cd289ac374ad8539930628314bc581a481cdb41409c9289419ddb72.json
repositories
manifest.json
a187dde48cd289ac374ad8539930628314bc581a481cdb41409c9289419ddb72.json
$ ls
485d7306187faf0cc9b77fc210e8def9d67b1953f0669e675877e90e6542cb6d
$ cd 485d7306187faf0cc9b77fc210e8def9d67b1953f0669e675877e90e6542cb6d/
$ ls
VERSION         json            layer.tar
$ rm -rfv VERSION json
VERSION
json
$ ls
layer.tar
$ mv layer.tar ..
$ ls
$ cd ..
$ ls
485d7306187faf0cc9b77fc210e8def9d67b1953f0669e675877e90e6542cb6d
layer.tar
$ rm -rfv 485d7306187faf0cc9b77fc210e8def9d67b1953f0669e675877e90e6542cb6d/
485d7306187faf0cc9b77fc210e8def9d67b1953f0669e675877e90e6542cb6d/
$ ls
layer.tar
$ tar xvf layer.tar
x bin/
x bin/arch
x bin/ash
x bin/base64
x bin/bbconfig
x bin/busybox
x bin/cat
...
...
$ ls
bin             home            media           proc            sbin            tmp
dev             layer.tar       mnt             root            srv             usr
etc             lib             opt             run             sys             var
$ rm -rfv layer.tar
layer.tar
$ ls
bin     etc     lib     mnt     proc    run     srv     tmp     var
dev     home    media   opt     root    sbin    sys     usr
$ touch DUMMY_FILE_TO_MARK_ROOT_OF_CONTAINER
$ ls
DUMMY_FILE_TO_MARK_ROOT_OF_CONTAINER    proc
bin                                     root
dev                                     run
etc                                     sbin
home                                    srv
lib                                     sys
media                                   tmp
mnt                                     usr
opt                                     var
```

That was a bit tedious. I think there's a better way to it. I was actually
thinking of just running the container and copying all the files to the host
machine. I think even that would have just worked. I don't know 🙈 I actually
want a more simpler way. Anyways, I'll get back to that later! Let me mount
this `alpine` directory into the Linux VM that I have

```bash
$ multipass mount alpine containers:/alpine-fs
```

Done!

```bash
ubuntu@containers:/doppledocker$ ls /alpine-fs/
DUMMY_FILE_TO_MARK_ROOT_OF_CONTAINER  dev  home  media  opt   root  sbin  sys  usr
bin                                   etc  lib   mnt    proc  run   srv   tmp  var
```

Now, let's write the code for this! :) Doing `chroot` using `syscall.Chroot()`
;) I'm not going to use `syscall.Chdir()` for now and see what weird behavior
happens if that's done!

This is what I get

```bash
ubuntu@containers:/doppledocker$ sudo go run main.go run /bin/bash
Running [/bin/bash] with process ID 1543
Running [/bin/bash] with process ID 1
panic: fork/exec /bin/bash: no such file or directory

goroutine 1 [running]:
main.child()
        /doppledocker/main.go:66 +0x30a
main.main()
        /doppledocker/main.go:18 +0x78
panic: exit status 2

goroutine 1 [running]:
main.run()
        /doppledocker/main.go:40 +0x433
main.main()
        /doppledocker/main.go:16 +0x55
exit status 2
```

Hmm. But if I remove the `Chroot` code, I can see it properly running

```bash
ubuntu@containers:/doppledocker$ sudo go run main.go run /bin/bash
Running [/bin/bash] with process ID 1578
Running [/bin/bash] with process ID 1
root@dopplecontainer:/doppledocker# exit
```

I think I'll add both `Chroot` and `Chdir` and try!

Still does not work! 🙈

```bash
ubuntu@containers:/doppledocker$ sudo go run main.go run /bin/bash
Running [/bin/bash] with process ID 1623
Running [/bin/bash] with process ID 1
panic: fork/exec /bin/bash: no such file or directory

goroutine 1 [running]:
main.child()
        /doppledocker/main.go:66 +0x323
main.main()
        /doppledocker/main.go:18 +0x78
panic: exit status 2

goroutine 1 [running]:
main.run()
        /doppledocker/main.go:40 +0x433
main.main()
        /doppledocker/main.go:16 +0x55
exit status 2
```

Ohhhhh. Sorry. I forgot. 🙈 `/bin/bash` is not present in `alpine`. Damn `alpine`
lol. 😅😛 So, that's why it said `/bin/bash` not found. I was thinking the
error is not complete and is possibly because I didn't put the `Chdir` code but
then it didn't work even with `Chdir` code. But that wasn't the issue. Hmm.

Let me try `/bin/sh` but without `Chdir`! :)

Okay, so this is what happened!

```sh
ubuntu@containers:/doppledocker$ sudo go run main.go run /bin/sh
Running [/bin/sh] with process ID 1737
Running [/bin/sh] with process ID 1
sh: getcwd: No such file or directory
(unknown) # pwd

sh: getcwd: No such file or directory
(unknown) # ls
README.md  go.mod     main.go
sh: getcwd: No such file or directory
(unknown) # cd ..
sh: getcwd: No such file or directory
(unknown) # ls
alpine-fs       doppledocker    initrd.img.old  media           root            srv             var
bin             etc             lib             mnt             run             sys             vmlinuz
boot            home            lib64           opt             sbin            tmp             vmlinuz.old
dev             initrd.img      lost+found      proc            snap            usr
sh: getcwd: No such file or directory
(unknown) # cd ..
sh: getcwd: No such file or directory
(unknown) # ls
alpine-fs       doppledocker    initrd.img.old  media           root            srv             var
bin             etc             lib             mnt             run             sys             vmlinuz
boot            home            lib64           opt             sbin            tmp             vmlinuz.old
dev             initrd.img      lost+found      proc            snap            usr
sh: getcwd: No such file or directory
(unknown) # ls proc/
1                  162                209                425                bus                meminfo
10                 1639               21                 431                cgroups            misc
109                164                210                437                cmdline            modules
11                 166                216                557                consoles           mounts
1123               167                22                 6                  cpuinfo            mtrr
1129               168                23                 695                crypto             net
1134               17                 24                 7                  devices            pagetypeinfo
12                 1715               25                 716                diskstats          partitions
1282               1716               26                 78                 dma                sched_debug
1293               1737               27                 79                 driver             schedstat
1294               1743               28                 8                  execdomains        scsi
...
...
sh: getcwd: No such file or directory
(unknown) #
```

So, I was able to navigate the host machine file system with `cd` and check
stuff and I had access to everything from inside the container! 😅

But let's see what happens if I do this

```sh
(unknown) # cd /
/ # pwd
/
/ # ls
DUMMY_FILE_TO_MARK_ROOT_OF_CONTAINER  media                                 sbin
bin                                   mnt                                   srv
dev                                   opt                                   sys
etc                                   proc                                  tmp
home                                  root                                  usr
lib                                   run                                   var
/ # cd ..
/ # ls
DUMMY_FILE_TO_MARK_ROOT_OF_CONTAINER  media                                 sbin
bin                                   mnt                                   srv
dev                                   opt                                   sys
etc                                   proc                                  tmp
home                                  root                                  usr
lib                                   run                                   var
/ # cd ..
/ # ls
DUMMY_FILE_TO_MARK_ROOT_OF_CONTAINER  media                                 sbin
bin                                   mnt                                   srv
dev                                   opt                                   sys
etc                                   proc                                  tmp
home                                  root                                  usr
lib                                   run                                   var
/ #
```

So, once I `cd` into `/`, then it's not able to see the host machine file system!
Seems like if we don't do the `Chdir()`, it's more of a possible security issue
in our case :P So, let's do it in the code!

Okay, all good now :)

```sh
ubuntu@containers:/doppledocker$ sudo go run main.go run /bin/sh
Running [/bin/sh] with process ID 1778
Running [/bin/sh] with process ID 1
/ # ls
DUMMY_FILE_TO_MARK_ROOT_OF_CONTAINER  media                                 sbin
bin                                   mnt                                   srv
dev                                   opt                                   sys
etc                                   proc                                  tmp
home                                  root                                  usr
lib                                   run                                   var
/ # cd /
/ # cd ..
/ # ls
DUMMY_FILE_TO_MARK_ROOT_OF_CONTAINER  media                                 sbin
bin                                   mnt                                   srv
dev                                   opt                                   sys
etc                                   proc                                  tmp
home                                  root                                  usr
lib                                   run                                   var
/ #
```

In the host machine, this shows as

```bash
ubuntu@containers:~$ ps aux | grep /bin/sh
root      1792  0.0  0.4  63968  4192 pts/0    S    22:54   0:00 sudo go run main.go run /bin/sh
root      1793  0.1  1.8 1163096 18496 pts/0   Sl   22:54   0:00 go run main.go run /bin/sh
root      1815  0.0  0.1 703228  1248 pts/0    Sl   22:54   0:00 /tmp/go-build860519042/b001/exe/main run /bin/sh
root      1819  0.0  0.1 703228  1244 pts/0    Sl   22:54   0:00 /proc/self/exe child /bin/sh
root      1824  0.0  0.0   1644     4 pts/0    S+   22:54   0:00 /bin/sh
ubuntu    1922  0.0  0.1  14856  1060 pts/1    S+   22:54   0:00 grep --color=auto /bin/sh
```

Now, checking the `ps aux` inside the container, weirdly I see this

```sh
/ # ps
PID   USER     TIME  COMMAND
/ # ps aux
PID   USER     TIME  COMMAND
/ # sleep 10 &
/ # ps
PID   USER     TIME  COMMAND
/ # ps aux
PID   USER     TIME  COMMAND
/ # sh
/ # ps
PID   USER     TIME  COMMAND
/ # ps aux
PID   USER     TIME  COMMAND
```

Okay, we will get back to that. For now, let's try this

```sh
ubuntu@containers:/doppledocker$ sudo go run main.go run /bin/sh
Running [/bin/sh] with process ID 2649
Running [/bin/sh] with process ID 1
/ # sleep 100
```

In the host we see this

```bash
root@containers:/# ps -C sleep
  PID TTY          TIME CMD
 2659 pts/1    00:00:00 sleep
root@containers:/# ls -l /proc/2659
total 0
dr-xr-xr-x 2 root root 0 May 30 08:14 attr
-rw-r--r-- 1 root root 0 May 30 08:14 autogroup
-r-------- 1 root root 0 May 30 08:14 auxv
-r--r--r-- 1 root root 0 May 30 08:14 cgroup
--w------- 1 root root 0 May 30 08:14 clear_refs
-r--r--r-- 1 root root 0 May 30 08:13 cmdline
-rw-r--r-- 1 root root 0 May 30 08:14 comm
-rw-r--r-- 1 root root 0 May 30 08:14 coredump_filter
-r--r--r-- 1 root root 0 May 30 08:14 cpuset
lrwxrwxrwx 1 root root 0 May 30 08:14 cwd -> /alpine-fs
-r-------- 1 root root 0 May 30 08:14 environ
lrwxrwxrwx 1 root root 0 May 30 08:14 exe -> /alpine-fs/bin/busybox
...
...
root@containers:/# ls -l /proc/2659/root
lrwxrwxrwx 1 root root 0 May 30 08:14 /proc/2659/root -> /alpine-fs
```

We notice the sleep process and what the root of the sleep process too. Now,
that's pretty much the equivalent of a container image. When you specify the
image it takes a copy of the file system that's packed up in that image, unpacks
it somewhere on your host machine and chroots (changes the root) of the container
to see just that new file system, so, we have kind of done the equivalent.

Again back to `ps` and `/proc` directory. We were trying to get `ps` to show
just the processes running inside the container.

We noticed that we don't get anything when we do `ps` in our alpine dopple
container.

Apparently if we do this in `ubuntu` `bash` in a completely new file system,
similar to how we did it in `alpine` with `sh`

```bash
ubuntu@ubuntu:/$ ps
Error, do this: mount -t proc proc /proc
```

Turns out, `/proc` is a pseudo file system, it's a mechanism for the kernel
and the user space to share information. And at the moment, `/proc` in our
container, in the chrooted file system, has nothing in it. And we need to
mount that directory as a proc pseudo file system, so that the kernel knows
we are going to populate that we are gonna populate that with all the
information about these running processes.

So, let's do the mounting and unmounting too in the code, for `/proc`, using
`syscall.Mount()` and `syscall.Unmount()`. We unmount when we finish using
it, that is when we finish running the command that the user provided us.

And now, let's try `ps`

```bash
ubuntu@containers:/doppledocker$ sudo go run main.go run /bin/sh
Running [/bin/sh] with process ID 2891
Running [/bin/sh] with process ID 1
/ # ps
PID   USER     TIME  COMMAND
    1 root      0:00 /proc/self/exe child /bin/sh
    5 root      0:00 /bin/sh
    6 root      0:00 ps
/ # mount
:/Users/karuppiahn/alpine on / type fuse.sshfs (rw,nosuid,nodev,relatime,user_id=0,group_id=0,allow_other)
proc on /proc type proc (rw,relatime)
proc on /proc type proc (rw,relatime)
/ #
```

And it works! :) We can see process IDs starting from `1`! :) And we can see
the mount inside the container. We can also see what directory of the host
machine is our root directory. In our case it's `/Users/karuppiahn/alpine`.

In our host machine, if we check proc related mounts, we get this

```bash
root@containers:/# mount | grep proc
proc on /proc type proc (rw,nosuid,nodev,noexec,relatime)
systemd-1 on /proc/sys/fs/binfmt_misc type autofs (rw,relatime,fd=37,pgrp=1,timeout=0,minproto=5,maxproto=5,direct,pipe_ino=11661)
proc on /alpine-fs/proc type proc (rw,relatime)
proc on /alpine-fs/proc type proc (rw,relatime)
```

You can see the mounts we did in the container in the host machine.

Now, there's a namespace for mounts, it's called `syscall.CLONE_NEWNS`, where
`NS` stands for `Namespace`. Apparently this was the first of the namespaces
to be invented and added to the Kernel. And probably at the time they didn't
really think there would ever be a need for any other namespaces, so they called
it namespace, but it's really for mounts.

By default, under systemd, mounts get this recursively shared property, and at
the moment, our root directory on the host machine recursively shares between
all namespaces, any mounts, and we need to deliberately turn that off with this
thing called `Unshareflags` in the code, along with `sycall.CLONE_NEWNS`. 

So, we have written the code to tell that

```
I have got this new namespace in my container and I don't want you to share it
with the host
```

Because by default it would have shared it, back with the host.

Now we see this in container

```bash
ubuntu@containers:/doppledocker$ sudo go run main.go run /bin/sh
Running [/bin/sh] with process ID 3102
Running [/bin/sh] with process ID 1
/ # mount
:/Users/karuppiahn/alpine on / type fuse.sshfs (rw,nosuid,nodev,relatime,user_id=0,group_id=0,allow_other)
proc on /proc type proc (rw,relatime)
/ #
```

And this in host

```bash
root@containers:/# mount | grep proc
proc on /proc type proc (rw,nosuid,nodev,noexec,relatime)
systemd-1 on /proc/sys/fs/binfmt_misc type autofs (rw,relatime,fd=37,pgrp=1,timeout=0,minproto=5,maxproto=5,direct,pipe_ino=11661)
proc on /alpine-fs/proc type proc (rw,relatime)
```

Hmm. Previously, the weird thing was, I was seeing two mounts for the `proc`
with respect to the `/alpine-fs/proc`. Hmm. I think this was due to the sharing
of mounts. 

I had to do this to unmount for now

```bash
root@containers:/# umount /alpine-fs/proc
root@containers:/# mount | grep proc
proc on /proc type proc (rw,nosuid,nodev,noexec,relatime)
systemd-1 on /proc/sys/fs/binfmt_misc type autofs (rw,relatime,fd=37,pgrp=1,timeout=0,minproto=5,maxproto=5,direct,pipe_ino=11661)
```

And I tried running the container again. This is what we see in the container

```bash
ubuntu@containers:/doppledocker$ sudo go run main.go run /bin/sh
Running [/bin/sh] with process ID 3102
Running [/bin/sh] with process ID 1
/ # mount
:/Users/karuppiahn/alpine on / type fuse.sshfs (rw,nosuid,nodev,relatime,user_id=0,group_id=0,allow_other)
proc on /proc type proc (rw,relatime)
/ #
```

And this is what we see in the host

```bash
root@containers:/# mount | grep proc
proc on /proc type proc (rw,nosuid,nodev,noexec,relatime)
systemd-1 on /proc/sys/fs/binfmt_misc type autofs (rw,relatime,fd=37,pgrp=1,timeout=0,minproto=5,maxproto=5,direct,pipe_ino=11661)
root@containers:/#
```

Okay, cool, so now we are not cluttered with the mounts inside the containers.
And the host does not really need to know about the container's mounts. Cool!

We could look at the mount information of the containers from the `/proc`
directory though. Let's see how to do that! In the container

```bash
ubuntu@containers:/doppledocker$ sudo go run main.go run /bin/sh
Running [/bin/sh] with process ID 3142
Running [/bin/sh] with process ID 1
/ # sleep 100
```

In the host

```bash
root@containers:/# ps -C sleep
  PID TTY          TIME CMD
 3153 pts/1    00:00:00 sleep
root@containers:/# cat /proc/3153/mounts
:/Users/karuppiahn/alpine / fuse.sshfs rw,nosuid,nodev,relatime,user_id=0,group_id=0,allow_other 0 0
proc /proc proc rw,relatime 0 0
root@containers:/# cat /proc/3153/mountinfo
432 399 0:50 / / rw,nosuid,nodev,relatime - fuse.sshfs :/Users/karuppiahn/alpine rw,user_id=0,group_id=0,allow_other
361 432 0:52 / /proc rw,relatime - proc proc rw
root@containers:/# cat /proc/3153/mountstats
device :/Users/karuppiahn/alpine mounted on / with fstype fuse.sshfs
device proc mounted on /proc with fstype proc
```

This way we can see the mounts a process is aware of, from the host perspective.
But they don't cluttering up our mount command.

So, we have checked out the namespaces for `Unix Timesharing System` which is
the hostname, `Process IDs`, `Mounts`. `Network` namespace works in a similar
way to make sure that the container sees only a specific set of network
interfaces, same for `User IDs` and `InterProcess Communications`, all working
in the same sort of way.

We also saw how `chroot` works and how that limits the containers so they can
only see a subset of the file system that the host can see.

There's one last property of containers and that's CGroups - Control Groups.
If namespaces restrict `What we can see from inside the container`, Control
Groups restrict `What we can use from inside the container`, that is, limit the
resources that we can use inside the container. And we achieve this - by
configuring another one of these pseudo file system interfaces, so it's another
set of what look like directories and files, but we can manipulate them to set
properties that we want the kernel to understand and the kernel will write
information into the file system, so that we can read it back out again.

And we could be talking about things like
- How much memory is the container allowed to use?
- How much CPU?
- How much IO bandwidth it's allowed?
- How many processes are allowed?

Let's see how the file system looks before we try out an example.

```bash
root@containers:/# cd /sys/fs/cgroup/
root@containers:/sys/fs/cgroup# ls
blkio  cpu,cpuacct  cpuset   freezer  memory   net_cls,net_prio  perf_event  rdma     unified
cpu    cpuacct      devices  hugetlb  net_cls  net_prio          pids        systemd
```

We can see a directory for each of the different types of control groups, that
we can setup.

Let's use memory as an example.

```bash
root@containers:/sys/fs/cgroup# cd memory
root@containers:/sys/fs/cgroup/memory# ls
cgroup.clone_children           memory.kmem.tcp.limit_in_bytes      memory.stat
cgroup.event_control            memory.kmem.tcp.max_usage_in_bytes  memory.swappiness
cgroup.procs                    memory.kmem.tcp.usage_in_bytes      memory.usage_in_bytes
cgroup.sane_behavior            memory.kmem.usage_in_bytes          memory.use_hierarchy
memory.failcnt                  memory.limit_in_bytes               notify_on_release
memory.force_empty              memory.max_usage_in_bytes           release_agent
memory.kmem.failcnt             memory.move_charge_at_immigrate     system.slice
memory.kmem.limit_in_bytes      memory.numa_stat                    tasks
memory.kmem.max_usage_in_bytes  memory.oom_control                  user.slice
memory.kmem.slabinfo            memory.pressure_level
memory.kmem.tcp.failcnt         memory.soft_limit_in_bytes
```

And there's actually a large number of parameters, that you can set, related to
memory. And we could look at `memory.limit_in_bytes`

```bash
root@containers:/sys/fs/cgroup/memory# cat memory.limit_in_bytes
9223372036854771712
```

That gives us a very very large number, which is telling us - by default,
processes can use all the memory in the system

I'm going to install docker in my host machine.

```bash
ubuntu@containers:~$ sudo apt install docker.io
ubuntu@containers:~$ sudo docker ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
ubuntu@containers:~$ sudo docker run --rm -it alpine /bin/sh
Unable to find image 'alpine:latest' locally
latest: Pulling from library/alpine
df20fa9351a1: Pull complete
Digest: sha256:185518070891758909c9f839cf4ca393ee977ac378609f700f60a771a2dfe321
Status: Downloaded newer image for alpine:latest
/ # exit
ubuntu@containers:~$
```

Now, in the `/sys/fs/cgroup/memory/` we can see a directory called `docker`.
Hmm. What's there in that huh? ;)

```bash
ubuntu@containers:~$ cd /sys/fs/cgroup/memory/
ubuntu@containers:/sys/fs/cgroup/memory$ ls
cgroup.clone_children  memory.kmem.limit_in_bytes          memory.limit_in_bytes            memory.swappiness
cgroup.event_control   memory.kmem.max_usage_in_bytes      memory.max_usage_in_bytes        memory.usage_in_bytes
cgroup.procs           memory.kmem.slabinfo                memory.move_charge_at_immigrate  memory.use_hierarchy
cgroup.sane_behavior   memory.kmem.tcp.failcnt             memory.numa_stat                 notify_on_release
docker                 memory.kmem.tcp.limit_in_bytes      memory.oom_control               release_agent
memory.failcnt         memory.kmem.tcp.max_usage_in_bytes  memory.pressure_level            system.slice
memory.force_empty     memory.kmem.tcp.usage_in_bytes      memory.soft_limit_in_bytes       tasks
memory.kmem.failcnt    memory.kmem.usage_in_bytes          memory.stat                      user.slice
ubuntu@containers:/sys/fs/cgroup/memory$ cd docker/
ubuntu@containers:/sys/fs/cgroup/memory/docker$ ls
cgroup.clone_children           memory.kmem.tcp.failcnt             memory.oom_control
cgroup.event_control            memory.kmem.tcp.limit_in_bytes      memory.pressure_level
cgroup.procs                    memory.kmem.tcp.max_usage_in_bytes  memory.soft_limit_in_bytes
memory.failcnt                  memory.kmem.tcp.usage_in_bytes      memory.stat
memory.force_empty              memory.kmem.usage_in_bytes          memory.swappiness
memory.kmem.failcnt             memory.limit_in_bytes               memory.usage_in_bytes
memory.kmem.limit_in_bytes      memory.max_usage_in_bytes           memory.use_hierarchy
memory.kmem.max_usage_in_bytes  memory.move_charge_at_immigrate     notify_on_release
memory.kmem.slabinfo            memory.numa_stat                    tasks
ubuntu@containers:/sys/fs/cgroup/memory/docker$
```

It another set of all the same parameters. Some of these are also statistics
being reported back into user space

Let's see what happens when we run a docker container

```bash
ubuntu@containers:~$ docker run --rm -it alpine /bin/sh
/ # hostname
c9da12d24350
/ #
```

In the host machine, we see this with regards to `cgroup` and `docker`

```bash
ubuntu@containers:/sys/fs/cgroup/memory/docker$ sudo docker ps -q
c9da12d24350
ubuntu@containers:/sys/fs/cgroup/memory/docker$ ls
c9da12d24350e0817efa2ce5f1af413099f8ac6b154e6f6cc41c7cd6f3d7906b  memory.kmem.usage_in_bytes
cgroup.clone_children                                             memory.limit_in_bytes
cgroup.event_control                                              memory.max_usage_in_bytes
cgroup.procs                                                      memory.move_charge_at_immigrate
memory.failcnt                                                    memory.numa_stat
memory.force_empty                                                memory.oom_control
memory.kmem.failcnt                                               memory.pressure_level
memory.kmem.limit_in_bytes                                        memory.soft_limit_in_bytes
memory.kmem.max_usage_in_bytes                                    memory.stat
memory.kmem.slabinfo                                              memory.swappiness
memory.kmem.tcp.failcnt                                           memory.usage_in_bytes
memory.kmem.tcp.limit_in_bytes                                    memory.use_hierarchy
memory.kmem.tcp.max_usage_in_bytes                                notify_on_release
memory.kmem.tcp.usage_in_bytes                                    tasks
ubuntu@containers:/sys/fs/cgroup/memory/docker$
```

We see the directory `c9da12d24350e0817efa2ce5f1af413099f8ac6b154e6f6cc41c7cd6f3d7906b`
now, which corresponds to the docker container with the ID `c9da12d24350`

So, docker has basically created a control group for this container. But we
didn't ask it for any particular restrictions and if we were to look at what's
inside that container's `memory.limit_in_bytes`

```bash
ubuntu@containers:/sys/fs/cgroup/memory/docker$ cat c9da12d24350e0817efa2ce5f1af413099f8ac6b154e6f6cc41c7cd6f3d7906b/memory.limit_in_bytes
9223372036854771712
```

It's still a massive number! 

Let's see what happens when we do contraint the memory. We can do that in
`docker` like this

```bash
ubuntu@containers:~$ sudo docker run --rm -it --memory 10M alpine /bin/sh
WARNING: Your kernel does not support swap limit capabilities or the cgroup is not mounted. Memory limited without swap.
/ #
```

And now if we check some details about the docker container in the host machine

```bash
ubuntu@containers:/sys/fs/cgroup/memory/docker$ sudo docker ps -q
44a70867824e

ubuntu@containers:/sys/fs/cgroup/memory/docker$ ls
44a70867824ed487bcb493dc17629dcd105588351fde88a491c00251da3a85c2  memory.kmem.usage_in_bytes
cgroup.clone_children                                             memory.limit_in_bytes
cgroup.event_control                                              memory.max_usage_in_bytes
cgroup.procs                                                      memory.move_charge_at_immigrate
memory.failcnt                                                    memory.numa_stat
memory.force_empty                                                memory.oom_control
memory.kmem.failcnt                                               memory.pressure_level
memory.kmem.limit_in_bytes                                        memory.soft_limit_in_bytes
memory.kmem.max_usage_in_bytes                                    memory.stat
memory.kmem.slabinfo                                              memory.swappiness
memory.kmem.tcp.failcnt                                           memory.usage_in_bytes
memory.kmem.tcp.limit_in_bytes                                    memory.use_hierarchy
memory.kmem.tcp.max_usage_in_bytes                                notify_on_release
memory.kmem.tcp.usage_in_bytes                                    tasks

ubuntu@containers:/sys/fs/cgroup/memory/docker$ cat 44a70867824ed487bcb493dc17629dcd105588351fde88a491c00251da3a85c2/memory.limit_in_bytes
10485760
```

`10485760` is a small number of bytes. It corresponds to [~10 MB](https://duckduckgo.com/?q=10485760+bytes+to+megabytes&t=ffab&ia=answer)

So, `docker` wrote that number into that file, and that's how it tells the
Kernel to limit that particular container to that amount of memory.

We are gonna try to do the same kind of thing in our `doppledocker` ;) And we
are going to do it for the number of processes. So, let's just have a quick
look at what to do

```bash
ubuntu@containers:/doppledocker$ sudo docker run --rm -it alpine /bin/sh
/ #
```

In a separate tab, in the host machine

```bash
ubuntu@containers:~$ cd /sys/fs/cgroup/pids
ubuntu@containers:/sys/fs/cgroup/pids$ ls
ls
cgroup.clone_children  docker             system.slice
cgroup.procs           notify_on_release  tasks
cgroup.sane_behavior   release_agent      user.slice
ubuntu@containers:/sys/fs/cgroup/pids$ cd docker/
ubuntu@containers:/sys/fs/cgroup/pids/docker$ ls
0a8afe4f66657069b5deae0ffd826fe561976c5ac237a4894440e8eb160332a1  cgroup.procs       pids.current  pids.max
cgroup.clone_children                                             notify_on_release  pids.events   tasks
ubuntu@containers:/sys/fs/cgroup/pids/docker$ sudo docker ps -q
0a8afe4f6665
ubuntu@containers:/sys/fs/cgroup/pids/docker$ cat 0a8afe4f66657069b5deae0ffd826fe561976c5ac237a4894440e8eb160332a1/pids.max
max
ubuntu@containers:/sys/fs/cgroup/pids/docker$ cat 0a8afe4f66657069b5deae0ffd826fe561976c5ac237a4894440e8eb160332a1/pids.current
1
```

So, as you can see, by default, when you create a `docker` container, there is
no limit to the number of processes that can be spawned inside that process.
But, we are going to create a control group that does limit the number of
processes

So, we create a directory for `doppledocker` under `/sys/fs/cgroup/pids` and
create some files for limiting the max number of processes the child process
can create. `pids.max` to tell the max number of processes, `cgroup.procs` to
tell the process ID for which this control group applies, and then
`notify_on_release` is set to `1`, which removes the process ID from the file
`cgroup.procs` after the process, that is, the container exits.

So, initially kept the max number of processes as `2`, thinking the child
command will be one process and the `/bin/sh` will be another, and child
command's process ID is written to the `cgroup.procs` file. But this is what
happened when I was trying things

```bash
ubuntu@containers:/doppledocker$ sudo go run main.go run /bin/sh
Running [/bin/sh] with process ID 2824
Running [/bin/sh] with process ID 1
panic: open /sys/fs/cgroup/pids/doppledocker/cgroups.procs: permission denied

goroutine 1 [running]:
main.must(...)
        /doppledocker/main.go:100
main.controlGroup()
        /doppledocker/main.go:95 +0x54d
main.child()
        /doppledocker/main.go:51 +0x185
main.main()
        /doppledocker/main.go:21 +0x78
panic: exit status 2

goroutine 1 [running]:
main.run()
        /doppledocker/main.go:44 +0x43b
main.main()
        /doppledocker/main.go:19 +0x55
exit status 2
```

Initially I was confused. As I was running stuff with `sudo` and was thinking
if somewhere some another access issue happened. Turns out I miss spelled the
file name. It's `cgroup.procs` and not `cgroups.procs` and it gave me permission
denied error when I tried to write info to a wrong file! Hmm. Weird computer
errors. Anyways. I fixed that, then I got this

```bash
ubuntu@containers:/doppledocker$ sudo go run main.go run /bin/sh
Running [/bin/sh] with process ID 2857
Running [/bin/sh] with process ID 1
panic: fork/exec /bin/sh: resource temporarily unavailable

goroutine 1 [running]:
main.child()
        /doppledocker/main.go:72 +0x3a4
main.main()
        /doppledocker/main.go:21 +0x78
panic: exit status 2

goroutine 1 [running]:
main.run()
        /doppledocker/main.go:44 +0x43b
main.main()
        /doppledocker/main.go:19 +0x55
exit status 2
```

I realized there was some process number issue, as in, not many processes were
able to run, and hence the resource error. I started increasing the number from
`2` to `3`, `4` and it still didn't work. Finally with `5`, this is what I got

```bash
ubuntu@containers:/doppledocker$ sudo go run main.go run /bin/sh
Running [/bin/sh] with process ID 2958
Running [/bin/sh] with process ID 1
/ # ls
/bin/sh: can't fork: Resource temporarily unavailable
/ # ps
/bin/sh: can't fork: Resource temporarily unavailable
/ # echo ok
ok
/ # ps
/bin/sh: can't fork: Resource temporarily unavailable
/ # cd ..
/ # ls
/bin/sh: can't fork: Resource temporarily unavailable
/ # cd
bin/    etc/    lib/    mnt/    proc/   run/    srv/    tmp/    var/
dev/    home/   media/  opt/    root/   sbin/   sys/    usr/
/ # cd bin/
/bin # ls
/bin/sh: can't fork: Resource temporarily unavailable
/bin # sleep 1
/bin/sh: can't fork: Resource temporarily unavailable
/bin #
```

Seems like I could do stuff like `echo`, `cd`, but anything else, as simple as
`ls` or `ps` or `sleep` or any other usual commands I used didn't work. This was
because we were limiting the number of processes that can run in the container
to `5`. I increased this to `6`

In the container

```sh
/ # ps
PID   USER     TIME  COMMAND
    1 root      0:00 /proc/self/exe child /bin/sh
    5 root      0:00 /bin/sh
    9 root      0:00 ps
/ # sleep 100 &
/ # ps
/bin/sh: can't fork: Resource temporarily unavailable
/ #
/ # ps
/bin/sh: can't fork: Resource temporarily unavailable
/ #
```

In the host machine

```bash
ubuntu@containers:/sys/fs/cgroup/pids$ ls
cgroup.clone_children  cgroup.sane_behavior  doppledocker       release_agent  tasks
cgroup.procs           docker                notify_on_release  system.slice   user.slice
ubuntu@containers:/sys/fs/cgroup/pids$ cd doppledocker/
bash: cd: doppledocker/: Permission denied
ubuntu@containers:/sys/fs/cgroup/pids$ sudo -i
root@containers:~# cd /sys/fs/cgroup/pids/doppledocker/
root@containers:/sys/fs/cgroup/pids/doppledocker# ls
cgroup.clone_children  cgroup.procs  notify_on_release  pids.current  pids.events  pids.max  tasks
root@containers:/sys/fs/cgroup/pids/doppledocker# cat pids.current
5
root@containers:/sys/fs/cgroup/pids/doppledocker# cat pids.max
6
root@containers:/sys/fs/cgroup/pids/doppledocker# cat cgroup.procs
3003
3008
root@containers:/sys/fs/cgroup/pids/doppledocker# cat cgroup.procs
3003
3008
3055
root@containers:/sys/fs/cgroup/pids/doppledocker# cat pids.max
6
root@containers:/sys/fs/cgroup/pids/doppledocker# cat pids.current
6
```

You can see how I used up the max number of processes by just running one long
running command / process. I had just enough space to run one extra process.
Somehow the container took 5 process IDs already when it started off. 

I stopped the container, that is, exited out of it and I see this on host

```bash
root@containers:/sys/fs/cgroup/pids/doppledocker# cat cgroup.procs
root@containers:/sys/fs/cgroup/pids/doppledocker#
```

I'm not able to fiddle with the `notify_on_release` to see how it works when the
value in it is `0` or when it's not even present. I changed the value to `0`
using `vi` vim. For removing it, I couldn't do it

```bash
root@containers:/sys/fs/cgroup/pids# rm -rfv doppledocker/
rm: cannot remove 'doppledocker/cgroup.procs': Operation not permitted
rm: cannot remove 'doppledocker/pids.current': Operation not permitted
rm: cannot remove 'doppledocker/pids.events': Operation not permitted
rm: cannot remove 'doppledocker/tasks': Operation not permitted
rm: cannot remove 'doppledocker/notify_on_release': Operation not permitted
rm: cannot remove 'doppledocker/pids.max': Operation not permitted
rm: cannot remove 'doppledocker/cgroup.clone_children': Operation not permitted
```

Even though I'm root! Anyways, let's meddle more with the max number of
processes. I'll make it `20` this time

And we are going to run what's called a [`Fork Bomb`](https://en.wikipedia.org/wiki/Fork_bomb)

Running the below `Fork Bomb` in a shell can make things crazy

```bash
$ :() { : | : & }; :
```

In the above fork bomb, we define a function called `:` (colon) and in it's
definition, we call itself and we pipe it's output to itself which is running
in the background. And after we are done defining it, we call the function `:`.
So, what happens is, it goes on and on and on.

```bash
ubuntu@containers:/doppledocker$ sudo go run main.go run /bin/sh
Running [/bin/sh] with process ID 3718
Running [/bin/sh] with process ID 1
/ # :() { : | : & }; :
/bin/sh: syntax error: bad function name
/ # something() { something | something & }; something
/ # /bin/sh: can't open '/dev/null': No such file or directory
/bin/sh: can't open '/dev/null': No such file or directory
/bin/sh: can't open '/dev/null': No such file or directory
/bin/sh: can't open '/dev/null': No such file or directory
/bin/sh: can't fork: Resource temporarily unavailable
/bin/sh: can't fork: Resource temporarily unavailable
/bin/sh: can't open '/dev/null': No such file or directory
/bin/sh: can't open '/dev/null': No such file or directory
/bin/sh: can't open '/dev/null': No such file or directory

[1]+  Done                       something | something
/ # something() { something | something & }; something
/bin/sh: can't fork: Resource temporarily unavailable
/bin/sh: can't fork: Resource temporarily unavailable
/ # something() { something | something & }; something
/bin/sh: can't fork: Resource temporarily unavailable
[1]+  Done(2)                    something | something
/bin/sh: can't fork: Resource temporarily unavailable
/ # /bin/sh: can't open '/dev/null': No such file or directory

/ #
```

So, it didn't allow `:` (colon) to be a function name, so I renamed it to
`something` and as we can see, it just couldn't run after sometime.

I tried the function again after sometime and now, the host machine is
unresponsive actually. I could see some absurd behavior a few moments before,
like, `cgroup.procs` always showed only two processes, even if I ran multiple
`sleep` processes in the background. Also, the `pids.current` always showed as
`19`. And previously when I tried the fork bomb, it didn't affect the container
much. It just stopped after sometime. But the last time I tried, it just
threw tons of logs telling resource is not available and then I noticed host
machine is also not repsonsive when I typed `ls`, see below

```bash
ubuntu@containers:/doppledocker$ ls
bash: fork: retry: Resource temporarily unavailable
^Cbash: fork: Interrupted system call

ubuntu@containers:/doppledocker$ ls
bash: fork: retry: Resource temporarily unavailable
bash: fork: retry: Resource temporarily unavailable

bash: fork: retry: Resource temporarily unavailable
bash: fork: Interrupted system call

ubuntu@containers:/doppledocker$ ls
bash: fork: retry: Resource temporarily unavailable
bash: fork: retry: Resource temporarily unavailable
bash: fork: Interrupted system call
ubuntu@containers:/doppledocker$
```

Anyways, the host machine was actually a VM and not my physical machine. So, all
is good. I restarted the VM using `multipass`

```
$ multipass stop containers
$ multipass start containers
$ multipass exec containers bash
ubuntu@containers:~$ cd /doppledocker/
ubuntu@containers:/doppledocker$
```

I tried the fork bomb again!

```bash
ubuntu@containers:/doppledocker$ sudo go run main.go run /bin/sh
Running [/bin/sh] with process ID 1709
Running [/bin/sh] with process ID 1
/ # something() { something | something & }; something
/ # /bin/sh: can't open '/dev/null': No such file or directory
/bin/sh: can't open '/dev/null': No such file or directory
/bin/sh: can't open '/dev/null': No such file or directory
/bin/sh: can't open '/dev/null': No such file or directory
/bin/sh: can't open '/dev/null': No such file or directory
/bin/sh: can't fork: Resource temporarily unavailable
/bin/sh: can't open '/dev/null': No such file or directory
/bin/sh: can't fork: Resource temporarily unavailable
/bin/sh: can't open '/dev/null': No such file or directory

[1]+  Done                       something | something
/ # ps
PID   USER     TIME  COMMAND
    1 root      0:00 /proc/self/exe child /bin/sh
    5 root      0:00 /bin/sh
    8 root      0:00 [sh]
    9 root      0:00 [sh]
   10 root      0:00 [sh]
   11 root      0:00 [sh]
   12 root      0:00 [sh]
   13 root      0:00 [sh]
   14 root      0:00 [sh]
   15 root      0:00 [sh]
   16 root      0:00 [sh]
   17 root      0:00 [sh]
   18 root      0:00 [sh]
   19 root      0:00 [sh]
   20 root      0:00 [sh]
   23 root      0:00 ps
```

Now I can see the processes! Hmm. I tried some `sleep` commands in the background

```sh
/ # sleep 10 &
/ # ps
PID   USER     TIME  COMMAND
    1 root      0:00 /proc/self/exe child /bin/sh
    5 root      0:00 /bin/sh
    8 root      0:00 [sh]
    9 root      0:00 [sh]
   10 root      0:00 [sh]
   11 root      0:00 [sh]
   12 root      0:00 [sh]
   13 root      0:00 [sh]
   14 root      0:00 [sh]
   15 root      0:00 [sh]
   16 root      0:00 [sh]
   17 root      0:00 [sh]
   18 root      0:00 [sh]
   19 root      0:00 [sh]
   20 root      0:00 [sh]
   24 root      0:00 sleep 10
   25 root      0:00 ps
/ # sleep 1
/ # sleep 20 &
[1]-  Done                       sleep 10
/ # ps
PID   USER     TIME  COMMAND
    1 root      0:00 /proc/self/exe child /bin/sh
    5 root      0:00 /bin/sh
    8 root      0:00 [sh]
    9 root      0:00 [sh]
   10 root      0:00 [sh]
   11 root      0:00 [sh]
   12 root      0:00 [sh]
   13 root      0:00 [sh]
   14 root      0:00 [sh]
   15 root      0:00 [sh]
   16 root      0:00 [sh]
   17 root      0:00 [sh]
   18 root      0:00 [sh]
   19 root      0:00 [sh]
   20 root      0:00 [sh]
   27 root      0:00 sleep 20
   28 root      0:00 ps
/ # sleep 20 &
/ # ps
/bin/sh: can't fork: Resource temporarily unavailable
/ # ps
/bin/sh: can't fork: Resource temporarily unavailable
/ #
```

In the host I noticed that I had exhausted the max number of processes through
the max number of process IDs that can run

```bash
root@containers:/sys/fs/cgroup/pids/doppledocker# cat pids.current
20
```

Able to run stuff again now in container

```sh
/ # ps
/bin/sh: can't fork: Resource temporarily unavailable
[1]+  Done                       sleep 20
[2]+  Done                       sleep 20
/ # ps
PID   USER     TIME  COMMAND
    1 root      0:00 /proc/self/exe child /bin/sh
    5 root      0:00 /bin/sh
    8 root      0:00 [sh]
    9 root      0:00 [sh]
   10 root      0:00 [sh]
   11 root      0:00 [sh]
   12 root      0:00 [sh]
   13 root      0:00 [sh]
   14 root      0:00 [sh]
   15 root      0:00 [sh]
   16 root      0:00 [sh]
   17 root      0:00 [sh]
   18 root      0:00 [sh]
   19 root      0:00 [sh]
   20 root      0:00 [sh]
   33 root      0:00 ps
/ #
```

It's back to okayish now

```bash
root@containers:/sys/fs/cgroup/pids/doppledocker# cat pids.current
18
```

We can get rid of the multiple `sh` processes in the background by killing them.
But I'm going to just exit the container! :)

The recent incident where host machine was affected, I could probably see how
to replicate it. I don't know how it happened, but it's scary, in the sense
that it affected the host machine, given it must be isolated and not affect
any other isolated processes or the host machine.

So, I tried again this time

```sh
/ # something() { something | something & }; something; something;
/bin/sh: can't fork: Resource temporarily unavailable
/bin/sh: can't fork: Resource temporarily unavailable
/bin/sh: can't fork: Resource temporarily unavailable/ #

[1]+  Done(1)                    something | something
/ # something() { something | something & }; something; something;
/bin/sh: can't fork: Resource temporarily unavailable
/bin/sh: can't fork: Resource temporarily unavailable
/bin/sh: can't fork: Resource temporarily unavailable/ #

[1]+  Done(2)                    something | something
/ # something() { something | (something &) }; something; something;
/bin/sh: can't fork: Resource temporarily unavailable
/bin/sh: can't fork: Resource temporarily unavailable
/bin/sh: can't fork: Resource temporarily unavailable
/bin/sh: can't fork: Resource temporarily unavailable
/bin/sh: can't fork: Resource temporarily unavailable
/ #
```

I got something like the above and my container was stuck. I couldn't access
the terminal.

I was able to access the host machine which worked fine

```bash
root@containers:/sys/fs/cgroup/pids/doppledocker#
root@containers:/sys/fs/cgroup/pids/doppledocker# ls
cgroup.clone_children  pids.current  tasks
cgroup.procs           pids.events
notify_on_release      pids.max
root@containers:/sys/fs/cgroup/pids/doppledocker# cat pids.current
20
root@containers:/sys/fs/cgroup/pids/doppledocker# cat pids.max
20
root@containers:/sys/fs/cgroup/pids/doppledocker# cat cgroup.procs
2002
2007
root@containers:/sys/fs/cgroup/pids/doppledocker# cat tasks
2002
2004
2005
2006
2007
root@containers:/sys/fs/cgroup/pids/doppledocker# ps aux | grep 2002
root      2002  0.0  0.1 703240  1248 pts/0    Sl   20:24   0:00 /proc/self/exe child /bin/sh
root      2077  0.0  0.1  14856  1028 pts/1    S+   20:27   0:00 grep --color=auto 2002
root@containers:/sys/fs/cgroup/pids/doppledocker# ps aux | grep 2007
root      2007  0.0  0.0   1648  1000 pts/0    T    20:24   0:00 /bin/sh
root      2079  0.0  0.1  14856  1052 pts/1    S+   20:27   0:00 grep --color=auto 2007
root@containers:/sys/fs/cgroup/pids/doppledocker# ps aux | grep 2002
root      2002  0.0  0.1 703240  1248 pts/0    Sl   20:24   0:00 /proc/self/exe child /bin/sh
root      2081  0.0  0.1  14856  1076 pts/1    S+   20:27   0:00 grep --color=auto 2002
root@containers:/sys/fs/cgroup/pids/doppledocker# ps aux | grep 2004
root      2083  0.0  0.1  14856  1032 pts/1    S+   20:27   0:00 grep --color=auto 2004
root@containers:/sys/fs/cgroup/pids/doppledocker# ps aux | grep 2005
root      2085  0.0  0.1  14856  1148 pts/1    R+   20:27   0:00 grep --color=auto 2005
root@containers:/sys/fs/cgroup/pids/doppledocker# ps aux | grep 2006
root      2087  0.0  0.1  14856  1048 pts/1    S+   20:28   0:00 grep --color=auto 2006
root@containers:/sys/fs/cgroup/pids/doppledocker# ps aux | grep 2007
root      2007  0.0  0.0   1648  1000 pts/0    T    20:24   0:00 /bin/sh
root      2089  0.0  0.0  14856  1004 pts/1    R+   20:28   0:00 grep --color=auto 2007
```

Apparently the `tasks` file shows the child / forked processes and I can see
that only the processes in `cgroup.procs` are still present, through `ps aux`

Let's kill them all to again use my container! 

```bash
root@containers:/sys/fs/cgroup/pids/doppledocker# kill 2002
root@containers:/sys/fs/cgroup/pids/doppledocker# kill 2007
-bash: kill: (2007) - No such process
```

I think `2002` was a parent process? I guess. Anyways, the terminal in which
the container was running was showing some weird things! I couldn't get it
to get work! I had to finally exit the terminal! 😅

```bash
/ # panic: exit status 2

goroutine 1 [running]:
main.run()
        /doppledocker/main.go:44 +0x43b
main.main()
        /doppledocker/main.go:19 +0x55
exit status 2
ubuntu@containers:/doppledocker$ bash: cd: too many arguments
ubuntu@containers:/doppledocker$ ubuntu@containers:/doppledocker$ ubuntu@containers:/doppledocker$   PID TTY          TIME CMD
 1667 pts/0    00:00:00 bash
 2094 pts/0    00:00:00 ps
ubuntu@containers:/doppledocker$ ubuntu@containers:/doppledocker$ ubuntu@containers:/doppledocker$ ubuntu@containers:/doppledocker$ ubuntu@containers:/doppledocker$ README.md  go.mod  main.go
ubuntu@containers:/doppledocker$ ubuntu@containers:/doppledocker$ ubuntu@containers:/doppledocker$ ubuntu@containers:/doppledocker$ ubuntu@containers:/doppledocker$   PID TTY          TIME CMD
 1667 pts/0    00:00:00 bash
 2098 pts/0    00:00:00 ps
ubuntu@containers:/doppledocker$ ubuntu@containers:/doppledocker$ ^C
ubuntu@containers:/doppledocker$ ubuntu@containers:/doppledocker$ ubuntu@containers:/doppledocker$ ubuntu@containers:/doppledocker$ ^C
ubuntu@containers:/doppledocker$ ^C
ubuntu@containers:/doppledocker$ ^C
ubuntu@containers:/doppledocker$ ^C
ubuntu@containers:/doppledocker$ ^C
ubuntu@containers:/doppledocker$ ^C
ubuntu@containers:/doppledocker$ error: List of process IDs must follow q.

Usage:
 ps [options]

 Try 'ps --help <simple|list|output|threads|misc|all>'
  or 'ps --help <s|l|o|t|m|a>'
 for additional help text.

For more details see ps(1).
ubuntu@containers:/doppledocker$ ^C
ubuntu@containers:/doppledocker$ exit
```

Some weird thing happened. Anyways, that showed how number of processes can be
limited. That way, no one isolated process can bring down the system by using
up all the resources, trying a denial of service attack like fork bomb, like the
above

So, that's how you create a container on your own. An isolated process, using
namespaces, chroot and control groups.

Source code of the speaker
https://github.com/lizrice/containers-from-scratch
