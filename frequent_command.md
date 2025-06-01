Manage docker image in my own docker hub registry

```
docker build . -t duysmartum/<app_name>:<tag_name> && docker push duysmartum/<app_name>:<tag_name>

```

Import local image to k3d

```
k3d image import <image-name>
```

Create and manage a cluster

```
k3d cluster create -a 2
k3d cluster create --port 8082:30080@agent:0 -p 8081:80@loadbalancer --agents 2

k3d cluster stop

k3d cluster start

k3d cluster delete

kubectl get nodes -o wide
```

Pods

```
kubectl get pods -o wide

kubectl logs -f <pod-name>
kubectl logs -f <pod-name> -c <container-name>

kubectl port-forward <pod-name> <local-port>:<pod-port>
```

Handling deployment

```
kubectl cluster-info

kubectl create deployment hashgenerator-dep --image=jakousa/dwk-app1

kubectl scale deployment/hashgenerator-dep --replicas=4

kubectl set image deployment/hashgenerator-dep dwk-app1=jakousa/dwk-app1:b7fc18de2376da80ff0cfc72cf581a9f94d10e64

kubectl delete deployment hashgenerator-dep

kubectl apply -f manifests/deployment.yaml

kubectl delete -f manifests/deployment.yaml

kubectl apply -f https://raw.githubusercontent.com/kubernetes-hy/material-example/master/app1/manifests/deployment.yaml
```
