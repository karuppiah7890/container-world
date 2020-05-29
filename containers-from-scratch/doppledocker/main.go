package main

import (
	"fmt"
	"os"
	"os/exec"
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
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func child() {
	fmt.Printf("Running %v with process ID %d\n", os.Args[2:], os.Getpid())

	// This sets the hostname in the new namespace
	// that we created in the run command.
	syscall.Sethostname([]byte("dopplecontainer"))

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
}
