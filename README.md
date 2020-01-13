## WSO2 APIM Operator for Kubernetes

Deploying WSO2 APIM Patterns in kubernetes through a simple command. Introducing a new Custom Resource Definition called APIManager to efficiently and easily deploy all 4 patterns, and the custom pattern in Kubernetes.

![K8S CRD workflow](https://lh3.googleusercontent.com/-wqlc7Sgs72s/XhbIAHVK36I/AAAAAAAADOM/_9lEe_RtNks9fj9j87zaB65dWI1bw2ONgCK8BGAsYHg/s0/pic1.png "K8S CRD workflow")

## Quick Start Guide

In this document, we will walk through the following.
* Deploy a default pattern in Kubernetes
* Deploy other patterns
* Deploy Custom pattern
* Expose Service via NodePort (default-LoadBalancer)
* Override Configurations (configmap, pvc, deploy-config)
* Add new Configurations (configmap, pvc)
* Invoke the API

### Installation Prerequisites
* [Golang](https://golang.org/doc/install) v1.12+ 
* [Kubernetes cluster](https://kubernetes.io/docs/setup/) and client v1.12 or above
* [Docker](https://docs.docker.com/install/) & [DockerHub](https://hub.docker.com/) / private docker registry account
* [Helm](http://docs.shippable.com/deploy/tutorial/deploy-to-gcp-gke-helm/)

## Steps:
1. Clone the wso2am-k8s-operator repository

``` 
    git clone https://github.com/wso2-incubator/wso2am-k8s-operator.git 
```

2. Setup the required Storage
```
    helm install stable/nfs-server-provisioner
```
   
<details>
<summary>Advanced</summary>
<br>
<p>
GCP Users:
    
- External NFS setup.
  [doc](https://docs.google.com/document/d/1oLLbz5q53_vN9fXN-byXuCifdobT-_jXAno7zc87Gnk/edit?ts=5e16c0ca)
    
    
Minikube Users:

- HostPath setup.
</p>
</details>

  
    
3. Apply the command to create the controller-artifacts (in wso2-system namespace)

``` 
    kubectl apply -f artifacts/install/controller-artifacts/ 

    Output: 

    namespace/wso2-system created
    serviceaccount/wso2am-pattern-1-svc-account created
    clusterrole.rbac.authorization.k8s.io/wso2am-controller-role created
    clusterrolebinding.rbac.authorization.k8s.io/wso2am-controller-role-binding created
    customresourcedefinition.apiextensions.k8s.io/apimanagers.apim.wso2.com created
    deployment.apps/wso2am-controller created
```
4. Apply the command below to create controller-configs (in wso2-system namespace)
```
    kubectl apply -f artifacts/install/controller-configs/
    
    Output:
    
    configmap/controller-config created
    configmap/pvc-config created
```

5. Apply the command below to create pattern-spceific api manager artifacts
```
    kubectl apply -f artifacts/install/api-manager-artifacts/pattern-1/
    
    Output:
    
    configmap/wso2am-pattern-1-am-analytics-dashboard-bin created
    configmap/dash-conf created
    configmap/worker-conf created
    configmap/wso2am-pattern-1-am-1-conf created
    configmap/wso2am-pattern-1-am-2-conf created
    configmap/mysql-dbscripts created
```

**Sample Scenarios**
[Scenario-1](https://github.com/wso2-incubator/wso2am-k8s-operator/edit/master/README.md) 
