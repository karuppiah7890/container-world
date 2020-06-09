# Demo

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
- How to pass environment variables?
- How to mention ports?
- Multiple containers in a pod
- Can we access pod IP from outside the cluster?
- Port Forwarding a pod
