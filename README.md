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

   helm install stable/nfs-server-provisioner
   
   
<details>
<summary>Advanced</summary>
<br>
GCP Users:
    External NFS setup can be done
Minikube Users:
    HostPath setup can be done</details>

   
   <details><summary>For Advanced settings click here</summary>
<p>


   
   [GKE users](https://docs.google.com/document/d/1oLLbz5q53_vN9fXN-byXuCifdobT-_jXAno7zc87Gnk/edit?ts=5e16c0ca)
    
    [Minikube users](https://docs.google.com/document/d/1ILIQKGqZ53y2cMhS731RRZMKsdbY3C-OSi4M10g7i8Q/edit?usp=sharing)
   
    
3. Apply the command to get the controller-artifacts (in wso2-system namespace)

``` 
    kubectl apply -f artifacts/install/controller-artifacts/ 

    Output: 

    namespace/wso2-system created
    serviceaccount/wso2am-pattern-1-svc-account created
    clusterrole.rbac.authorization.k8s.io/wso2am-controller-role created
    clusterrolebinding.rbac.authorization.k8s.io/wso2am-controller-role-binding created
    customresourcedefinition.apiextensions.k8s.io/apimanagers.apim.wso2.com created
    apimanager.apim.wso2.com/cluster-1 created
    deployment.apps/wso2am-controller created
```
5. Now view the running pods by executing the command. Make sure to replace <USER-NAMESPACE> to the name changed in step 3
   By default pattern-1 is executing, you will see 5 pods. Each pod represents the components based on relavant patterns.
    
```
    kubectl get pods -n <USER-NAMESPACE>

    Output:

    NAME                                                       READY   STATUS    RESTARTS   AGE
    analytics-dash-deploy-54bd8d9b55-rmwnn                     1/1     Running   0          3m35s
    analytics-worker-deploy-79dc97599d-m445h                   1/1     Running   0          3m35s
    apim-1-deploy-7fcd974f8-m7ghq                              1/1     Running   0          3m35s
    apim-2-deploy-6bb4bff84-6cmz2                              1/1     Running   0          3m35s
    wso2apim-with-analytics-mysql-deployment-5fccb54d6-p29z5   1/1     Running   0          3m35s


```

6. Once the status becomes running, view the logs of each pod using following command.

```
    kubectl logs <POD-NAME> -n <USER-NAMESPACE>
```

**Note:** 
- GCP : To access the API portals, Add host mapping entries to the /etc/hosts file. As we have exposed the API portal service in Node Port type, you can use the IP address of any Kubernetes node.
```$xslt
    (kubectl get nodes -o jsonpath='{ $.items[*].status.addresses[?(@.type=="ExternalIP")].address }')
```
- Minikube: Add minikube ip to the /etc/hosts file

```
    <Any K8s Node IP / Minikube IP>  wso2apim
    <Any K8s Node IP / Minikube IP>  wso2apim-analytics
```


   _APIM Publisher_ - https://wso2apim:9443/publisher
   
   _APIM Devportal_ - https://wso2apim:9443/devportal
   

After successfully accessing the portals, Follow the below documentation and try out the complete workflow. 

[API Manager Documentation 3.0.0](https://apim.docs.wso2.com/en/latest/)
   
   

