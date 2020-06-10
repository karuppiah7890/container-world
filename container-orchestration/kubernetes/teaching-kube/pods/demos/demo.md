# Demo

- Get nodes
- Get pods
- Create a pod
- Check the pod IP and node in which the pods are running in
- Describe a pod
- Get logs of a pod
- Execute into a pod
- Communicate from one pod to another pod
- Debugging issues in pod
    - Wrong image
        - image pull issue possible reasons
            - authorization - private image registry service where your stored securely
              behind an authorization wall which needs credentials to access the image
            - wrong image repository name, or version etc
- delete pods
- What happens when containers in a pod restarts?
    - restart reasons
        - kill command
        - main process in container failed due to some reason
- How to pass environment variables?
- Not all fields in a pod can be updated, only some fields can be update ( mutable fields ).
  Example: ports field in a container in a pod.
- How to mention ports?
- Multiple containers in a pod
- Can we access pod IP from outside the cluster?
- Port Forwarding a pod's ports

---

You can find the demos as videos here - 

Pods Demo 1 https://www.youtube.com/watch?v=h8q7s1CKMQ4
Pods Demo 2 https://www.youtube.com/watch?v=JvPOYIE8DWo

And you can also find the terminal recordings of the demos using
[asciinema](https://asciinema.org) in this directory. Use this to play it

```bash
$ asciinema play <cast-file>
$ asciinema play pods-demo-1.cast
$ asciinema play pods-demo-2.cast
```
