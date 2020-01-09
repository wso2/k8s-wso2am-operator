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

    [setting up nfs](https://docs.google.com/document/d/1oLLbz5q53_vN9fXN-byXuCifdobT-_jXAno7zc87Gnk/edit?ts=5e16c0ca)
   
3. Go inside wso2am-k8s-operator/ folder
4. Execute the below command

``` kubectl apply -f artifacts/install/controller-artifacts/ ```

and you will get the following output
```
namespace/test created
serviceaccount/wso2am-pattern-1-svc-account created
clusterrole.rbac.authorization.k8s.io/wso2am-controller-role created
clusterrolebinding.rbac.authorization.k8s.io/wso2am-controller-role-binding created
customresourcedefinition.apiextensions.k8s.io/apimanagers.apim.wso2.com created
apimanager.apim.wso2.com/cluster-1 created
deployment.apps/wso2am-controller created
```





