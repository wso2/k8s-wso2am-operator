# WSO2 API Manager Operator for Kubernetes

With WSO2 API Manager Operator, it makes easy to deploy WSO2 API Manager in Kubernetes through a simple command. Also it supports deploying recommended deployment patterns in Kubernetes. Introducing a new Custom Resource Definition called APIManager to efficiently and easily deploy patterns, and the custom pattern in Kubernetes.

![K8S CRD workflow](docs/images/wso2am-operator.png "K8S CRD workflow")

## Quick Start Guide

In this document, we will walk through the following.
* Deploy a single node API Manager
* Access publisher and devportal

### Installation Prerequisites

* [Kubernetes cluster](https://kubernetes.io/docs/setup/) and client v1.12 or above
    * Minimum CPU : **4 vCPU** 
    * Minimum Memory : **6 GB** 

* Download [k8s-wso2am-operator-0.8.zip](https://github.com/wso2/k8s-wso2am-operator/releases/download/0.8/k8s-wso2am-operator-0.8.zip) and extract it. 

    This zip contains the artifacts that required to deploy in Kubernetes.

**Note:** You need to run all commands from within the k8s-wso2am-operator-0.8 directory.

#### Step 1: Deploy WSO2 API Manager Operator


``` 
kubectl apply -f artifacts/operator-artifacts/ -f artifacts/operator-configs/ -f artifacts/api-manager-artifacts/pattern-1/

Output: 

namespace/wso2-system created
serviceaccount/wso2am-pattern-1-svc-account created
...
configmap/wso2am-p1-apim-2-conf created
configmap/wso2am-p1-mysql-dbscripts created
```

#### Step 2: Deploy scenario 1 - Single node API Manager

```
kubectl apply -f scenarios/scenario-1/

Output:

configmap/apim-conf created
apimanager.apim.wso2.com/custom-pattern-1 created
```

* You can check the status of the API Manager pod and service by using the following commands.
    
```
kubectl get pods

Output:

NAME                                     READY   STATUS    RESTARTS   AGE
all-in-one-api-manager-d66d6c574-bnwps   0/1     Running   0          2m43s

kubectl get svc

NAME         TYPE        CLUSTER-IP    EXTERNAL-IP   PORT(S)                                                       AGE
kubernetes   ClusterIP   10.0.32.1     <none>        443/TCP                                                       6h28m
wso2apim     NodePort    10.0.36.136   <none>        8280:32004/TCP,8243:32003/TCP,9763:32002/TCP,9443:32001/TCP   4m32s
```

#### Step 3: Access API Manager

To access the API Manager, add a host mapping entry to the /etc/hosts file. As we have exposed the API portal service in Node Port type, you can use the IP address of any Kubernetes node.

```
<Any K8s Node IP>  wso2apim
```

- For Docker for Mac use "localhost" for the K8s node IP
- For Minikube, use minikube ip command to get the K8s node IP
- For GKE
    ```
    (kubectl get nodes -o jsonpath='{ $.items[*].status.addresses[?(@.type=="ExternalIP")].address }')
    ```
    - This will give the external IPs of the nodes available in the cluster. Pick any IP to include in /etc/hosts file.
  
- **API Publisher** : https://wso2apim:32001/publisher 
- **API Devportal** : https://wso2apim:32001/devportal 
   
   
## Sample Scenarios

1. [Scenario-1 : Applying Simple and shortest Custom Resource YAML](scenarios/scenario-2)
2. [Scenario-2 : Exposing via NodePort Service Type](scenarios/scenario-3)
3. [Scenario-3 : Override Deployment Configuration values](scenarios/scenario-5)
4. [Scenario-4 : Override ConfigMaps and PersistentVolumeClaims](scenarios/scenario-6)
5. [Scenario-5 : Add New Configmaps and Persistent Volume Claims](scenarios/scenario-7)
6. [Scenario-6 : Deploying Custom Pattern - Basic](scenarios/scenario-1)
7. [Scenario-7 : Deploying Custom Pattern - Advanced](scenarios/scenario-4)
8. [Scenario-8 : Running External-NFS](scenarios/scenario-8)
8. [Scenario-9 : Exposing using Ingresses](scenarios/scenario-9)

### Clean up

Execute the following command if you wish to clean up the Kubernetes cluster by removing all the applied artifacts and configurations related to WSO2AM Operator.

```
kubectl delete -f artifacts/ -R
```

### Troubleshooting Guide

You can refer [troubleshooting guide](docs/Troubleshooting/troubleshooting.md).

  
