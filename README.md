# WSO2 API Manager Operator for Kubernetes

With WSO2 API Manager Operator, it makes easy to deploy WSO2 API Manager in Kubernetes through a simple command. Also it supports deploying recommended deployment patterns in Kubernetes. Introducing a new Custom Resource Definition called APIManager to efficiently and easily deploy patterns, and custom patterns in Kubernetes.

![K8S CRD workflow](docs/images/wso2am-operator.png "K8S CRD workflow")

## Quick Start Guide

In this document, we will walk through the following.
* Deploy a single node API Manager
* Access publisher and devportal

### Installation Prerequisites

* [Kubernetes cluster](https://kubernetes.io/docs/setup/) and client v1.12 or above
    * Minimum CPU : **4 vCPU** 
    * Minimum Memory : **6 GB** <br>
    
    **Note:** This is for running scenario-1.

* Download [k8s-wso2am-operator-1.0.1.zip](https://github.com/wso2/k8s-wso2am-operator/releases/download/v1.0.1/k8s-wso2am-operator-1.0.1.zip) and extract it. 

    This zip contains the artifacts that required to deploy in Kubernetes.

    **Note:** You need to run all commands from within the k8s-wso2am-operator-1.0.1 directory.

<br>

#### Step 1: Deploy WSO2 API Manager Operator

- Creates a namespace called wso2-system and deploy the controller artifacts.

    ``` 
    >> kubectl apply -f artifacts/operator-artifacts/ -f artifacts/operator-configs/ -f artifacts/api-manager-artifacts/pattern-1/
    
    Output: 
    namespace/wso2-system created
    serviceaccount/wso2am-pattern-1-svc-account created
    ...
    configmap/wso2am-p1-apim-2-conf created
    configmap/wso2am-p1-mysql-dbscripts created
    ```
<br>

#### Step 2: Deploy scenario 1 - Single node API Manager

- Deploys a single node API Manager.

    ```
    >> kubectl apply -f scenarios/scenario-1/
    
    Output:
    configmap/apim-conf created
    apimanager.apim.wso2.com/custom-pattern-1 created
    ```

* You can check the status of the API Manager pod and service by using the following commands.
    
    ```
    >> kubectl get pods
    
    Output:
    NAME                                     READY   STATUS    RESTARTS   AGE
    all-in-one-api-manager-d66d6c574-bnwps   0/1     Running   0          2m43s
    
    >> kubectl get svc
    
    Output:
    NAME         TYPE        CLUSTER-IP    EXTERNAL-IP   PORT(S)                                                       AGE
    kubernetes   ClusterIP   10.0.32.1     <none>        443/TCP                                                       6h28m
    wso2apim     NodePort    10.0.36.136   <none>        8280:32004/TCP,8243:32003/TCP,9763:32002/TCP,9443:32001/TCP   4m32s
    ```
<br>

#### Step 3: Access API Manager

- To access the API Manager, add a host mapping entry to the /etc/hosts file. As we have exposed the API portal service in Node Port type, you can use the IP address of any Kubernetes node.

    ```
    <Any K8s Node IP>  wso2apim
    ```

- For Docker for Mac use "127.0.0.1" for the K8s node IP
- For Minikube, use minikube ip command to get the K8s node IP
- For GKE
    ```
    (kubectl get nodes -o jsonpath='{ $.items[*].status.addresses[?(@.type=="ExternalIP")].address }')
    ```
    - This will give the external IPs of the nodes available in the cluster. Pick any IP to include in /etc/hosts file.
  
- Access URLs 
    - **API Publisher** : https://wso2apim:32001/publisher 
    - **API Devportal** : https://wso2apim:32001/devportal 
       
<br>

## Sample Scenarios

1. [Scenario-1 : Deploy A Single API Manager Instance (Custom Pattern)](scenarios/scenario-1)
2. [Scenario-2 : Deploy API Manager Pattern-1 (LoadBalancer Service Type)](scenarios/scenario-2)
3. [Scenario-3 : Deploy API Manager Pattern-1 (NodePort Service Type)](scenarios/scenario-3)
4. [Scenario-4 : Deploy A Single API Manager Instance with Analytics (Custom Pattern)](scenarios/scenario-4)
5. [Scenario-5 : Override default configuration values](scenarios/scenario-5)
6. [Scenario-6 : Override ConfigMaps and PersistentVolumeClaims](scenarios/scenario-6)
7. [Scenario-7 : Add new configmaps and Persistent Volume Claims](scenarios/scenario-7)
8. [Scenario-8 : Running External NFS](scenarios/scenario-8)
9. [Scenario-9 : Expose API Manager using Ingress](scenarios/scenario-9)

## Clean up

- Execute the following command if you wish to clean up the Kubernetes cluster by removing all the applied artifacts and configurations related to WSO2AM Operator.

    ```
    >> kubectl delete -f artifacts/ -R
    ```

### Troubleshooting Guide

You can refer [troubleshooting guide](docs/Troubleshooting/troubleshooting.md).

  
