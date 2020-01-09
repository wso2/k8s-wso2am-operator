## WSO2 APIM Operator for Kubernetes

Deploying WSO2 APIM Patterns in kubernetes through a simple command. Introducing a new Custom Resource Definition called APIManager to efficiently and easily deploy all 4 patterns, and the custom pattern in Kubernetes.

![alt text](https://lh3.googleusercontent.com/-SDgiNAZzsD4/XhbH_LXDBRI/AAAAAAAADOI/Ani2mUSrfMI6yqJcIlYoXoQxmtKMdyxtwCLcBGAsYHQ/s0/pic1.png "K8S CRD workflow")

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
* [GCP](https://cloud.google.com/) / [Minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/)

## Steps:
1. Clone the wso2am-k8s-operator repository

``` 
git clone https://github.com/wso2-incubator/wso2am-k8s-operator.git 
```

2. Setup the required Storage

    [GCP users](https://docs.google.com/document/d/1oLLbz5q53_vN9fXN-byXuCifdobT-_jXAno7zc87Gnk/edit?ts=5e16c0ca)
    
    [Minikube users](https://docs.google.com/document/d/1ILIQKGqZ53y2cMhS731RRZMKsdbY3C-OSi4M10g7i8Q/edit?usp=sharing)
   
3. Go inside wso2am-k8s-operator/artifacts/install/controller-artifacts/ folder. Open the below set of files and replace <USER-NAMESPACE> with any name you like.
    1-namespace.yaml 
    2-service-account.yaml
    4-cluster-role-binding.yaml
    6-wso2-apim.yaml
    7-controller.yaml
    
4. After replacing the namespace, execute the following command

``` 
kubectl apply -f artifacts/install/controller-artifacts/ 

Output: 

namespace/<USER-NAMESPACE> created
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
kubectl logs apim-1-deploy-7fcd974f8-m7ghq -n <USER-NAMESPACE>
```








