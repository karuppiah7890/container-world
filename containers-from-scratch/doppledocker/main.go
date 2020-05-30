package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

// docker run <image> <cmd> <params>
// doppledocker run <cmd> <params>

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("bad command!")
	}
}

func run() {
	fmt.Printf("Running %v with process ID %d\n", os.Args[2:], os.Getpid())

	args := append([]string{"child"}, os.Args[2:]...)

	cmd := exec.Command("/proc/self/exe", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func child() {
	fmt.Printf("Running %v with process ID %d\n", os.Args[2:], os.Getpid())

	controlGroup()
	// This sets the hostname in the new namespace
	// that we created in the run command.
	syscall.Sethostname([]byte("dopplecontainer"))
	syscall.Chroot("/alpine-fs")
	syscall.Chdir("/")
	syscall.Mount("proc", "/proc", "proc", 0, "")

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// And we don't need to use SysProcAttr
	// again this time as run takes care
	// of creating the new namespace;
	// child command never creates a
	// new namespace. It just runs the command! :)

	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	syscall.Unmount("/proc", 0)
}

func controlGroup() {
	controlGroupPath := filepath.Join("/", "sys", "fs", "cgroup")
	pidsControlGroupPath := filepath.Join(controlGroupPath, "pids")
	doppleDockerPathForPIDs := filepath.Join(pidsControlGroupPath, "doppledocker")
	err := os.MkdirAll(doppleDockerPathForPIDs, 0700)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	pidsMaxPath := filepath.Join(doppleDockerPathForPIDs, "pids.max")
	must(ioutil.WriteFile(pidsMaxPath, []byte("20"), 0700))

	notifyOnReleasePath := filepath.Join(doppleDockerPathForPIDs, "notify_on_release")
	must(ioutil.WriteFile(notifyOnReleasePath, []byte("1"), 0700))

	cGroupProcsPath := filepath.Join(doppleDockerPathForPIDs, "cgroup.procs")
	currentProcessID := os.Getpid()
	must(ioutil.WriteFile(cGroupProcsPath, []byte(strconv.Itoa(currentProcessID)), 0700))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
