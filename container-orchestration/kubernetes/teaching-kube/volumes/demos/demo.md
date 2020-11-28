# Demo

<details>
<summary>emptyDir demo</summary>

One pod with two alpine or ubuntu containers, and an emptyDir volume.
Both containers mounting the volume to two different paths.
One container writes data to a file in the mount path - writing time every 2
seconds using sleep and date command.
Another container reads data from the file using tail command and streaming it.

The demo can be anything where one container is producing data into the volume and the
other is consuming data from the volume.
</details>

<details>
<summary>Where is the emptyDir present? In the node? Is it possible to find it??</summary>

</details>

<details>
<summary>What happens to the emptyDir when the pod dies?</summary>
</details>

<details>
</details>
