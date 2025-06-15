Part 1

Kubernetes is the sum of all the bash scripts and best practices that most system administrators would cobble together over time, presented as a single system behind a declarative set of APIs.

A container orchestration system such as Kubernetes is often required when maintaining containerized applications. The main responsibility of an orchestration system is the starting and stopping of containers. In addition, they offer networking between containers and health monitoring

docker compose, which also takes care of the same tasks; starting and stopping, networking and health monitoring. What makes Kubernetes special is the robust feature set for automating all of it

A cluster is a group of machines, nodes, that work together - in this case, they are part of a Kubernetes cluster.
"server node" to refer to nodes with control-plane and "agent node" to refer to the nodes without that role

"load balancer" proxy, that'll redirect a connection to 6443 into the server node, and that's how we can access the contents of the cluster

kubeconfig, a file that is used to organize information about clusters, users, namespaces, and authentication mechanisms

Kubectl is the Kubernetes command-line tool and will allow us to interact with the cluster. Kubectl will read kubeconfig from the location in KUBECONFIG environment value or by default from ~/.kube/config and use the information to connect to the cluster

To deploy an image, we need the cluster to have access to the image. By default, Kubernetes is intended to be used with a registry. To deploy an application, we will need to create a deployment object with the image.

Pod is an abstraction around one or more containers. Pods provide a context for 1..N containers so that they can share storage and a network. They can be thought of as a container of containers.

A deployment resource takes care of deployment. It's a way to tell Kubernetes what container you want, how they should be running and how many of them should be running.
While we created the Deployment we also created a ReplicaSet object. ReplicaSets are used to tell how many replicas of a Pod you want

Instead of deleting the deployment, we could just apply a modified deployment on top of what we already have. Kubernetes will take care of rolling out a new version. By using tags (e.g. dwk/image:tag) in the deployments, each time we update the image we can modify and apply the new deployment yaml.

When updating anything in Kubernetes the usage of delete is actually an anti-pattern and you should use it only as the last option. As long as you don't delete the resource Kubernetes will do a rolling update, ensuring minimum (or none) downtime for the application.

As you are trying to find bugs in your configuration, start by eliminating all possibilities one by one. The key is to be systematic and to question everything.

As Deployment resources took care of deployments for us. Service resources will take care of serving the application to connections from outside (and also inside!) of the cluster.

Service is just virtual network (clusterIP). It can have ports and it can listen to the "port" that it is told to.

The port that application is exposed to is the port in Pod resource that listen for traffic and redirect requests to the application.

NodePorts are simply ports that are opened by Kubernetes to all of the nodes and the service will handle requests in that port. NodePorts are not flexible and require you to assign a different port for every application. As such NodePorts are not used in production but are helpful to know about.

Incoming Network Access resource Ingress is a completely different type of resource from Services

Ingresses are implemented by various different "controllers". This means that ingresses do not automatically work in a cluster, but give you the freedom of choosing which ingress controller works for you the best. K3s has Traefik installed already. Other options include Istio and Nginx Ingress Controller, more here.

Ingress is used to route external HTTP/S traffic to Service resources, hence they can be used together. Ingress can be used with domain names.

There are two things that are known to be difficult with Kubernetes. First is networking. Thankfully we can avoid most of the networking difficulties unless we were going to setup our own cluster. If you're interested you can watch this Webinar on "Kubernetes and Networks: Why is This So Dang Hard?" but we'll skip most of the topics discussed in the video. The other of the most difficult things is storage.

The Kubernetes volumes, in technical terms emptyDir volumes, are shared filesystems inside a pod, this means that their lifecycle is tied to a pod. When the pod is destroyed the data is lost. In addition, simply moving the pod from another node will destroy the contents of the volume as the space is reserved from the node the pod is running on. So surely you should not use emptyDir volumes e.g. for backing up a database

Persistent Volumne - PV:
In contrast to the emptyDir volumes, a Persistent Volume is something you probably had in mind when we started talking about volumes.

A Persistent Volume (PV) is a cluster-wide resource, that represents a piece of storage in the cluster that has been provisioned by the cluster administrator or is dynamically provisioned. Persistent Volumes can be backed by various types of storage such as local disk, NFS, cloud storage, etc.

PVs have a lifecycle independent of any individual pod that uses the PV. This means that the data in the PV can outlive the pod that it was attached to.

It is the cloud provider that takes care of backing storage and the Persistent Volumes that you can use. If you run your own cluster or use a local cluster such as k3s for development, you need to take care of the storage system and Persistent Volumes by yourself.

An easy option that we can use with K3s is a local PersistentVolume that uses a path in a cluster node as the storage. This solution ties the volume to a particular node and if the node becomes unavailable, the storage is not usable.

So the local Persistent Volumes are not the solution to be used in production!

Persistent Volume Claim (PVC) is a request for storage by a user.

When a user creates a PVC, Kubernetes finds an appropriate PV that satisfies the claim's requirements and binds them together. If no PV is available, depending on the configuration, the cluster might dynamically create a PV that meets the claim's needs.

Conceptually, you can think of PVs as the physical volume (the actual storage in your infrastructure), whereas PVCs are the means by which pods claim this storage for their use.

Part 2:

Networking between pods

Kubernetes includes a DNS service so communication between pods and containers in Kubernetes is pretty similar as it was with containers in Docker compose. Containers in a pod share the network. As such every other container inside a pod is accessible from localhost.

For communication between Pods a Service is used as they expose the Pods as a network service.

Organizing a cluster

Namespaces are used to keep resources separated. A company that uses one cluster but has multiple projects can use namespaces to split the cluster into virtual clusters, one for each project. Most commonly they would be used to separate environments such as production, testing, staging. DNS entry for services includes the namespace so you can still have projects communicate with each other if needed through service.namespace address.

An administrator should set a ResourceQuota(opens in a new tab) for that namespace, so that you can safely run anything there.

Labels(opens in a new tab) are used to separate an application from others inside a namespace and to group different resources together. Labels are key-value pairs and they can be modified, added or removed at any time. Labels can also be added to almost anything.

Labels can help us humans identify resources and Kubernetes can use them to act upon a group of resources. You can query resources that have a certain label. The labels are also used by selectors(opens in a new tab) to pick a set of objects.

We can use the same label on multiple namespaces and the namespace would keep them from interfering with each other.

Grouping objects with labels is simple. We either add the label into the yaml file or use command kubectl label(opens in a new tab).

With labels, we can even move pods to labeled nodes. Let's say we have a few nodes which have qualities that we wish to avoid. For example, those might have a slower network. With labels and nodeSelector(opens in a new tab) configured for deployment, we can do just that.

nodeSelector is a blunt tool. It's great when you want to define binary qualities, like "don't run this application if the node is using an HDD instead of an SSD" by labeling the nodes according to disk types. There are more sophisticated tools you should use when you have a cluster of various machines, ranging from a fighter jet(opens in a new tab) to a toaster to a supercomputer. Kubernetes can use affinity(opens in a new tab) and anti-affinity to select which nodes are prioritized for which applications and taints with tolerances(opens in a new tab) so that a pod can avoid certain nodes. For example, if a machine has a high network latency and we wouldn't want it to do some latency-critical tasks.
