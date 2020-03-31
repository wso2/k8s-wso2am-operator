## Scenario 3 : Deploy API Manager Pattern-1 (NodePort Service Type)

In this scenario we are deploying API Manager Pattern-1 where we expose API Manager in NodePort service type.

#### Installation Prerequisites

* [Helm](https://helm.sh/docs/intro/install/) v3.1.0 or above

This is used to deploy the NFS Provisioner and this used for persistent volumes.

##### Setup persistent volumes

```
helm repo add stable https://kubernetes-charts.storage.googleapis.com
helm install nfs-provisioner stable/nfs-server-provisioner
```

#### Deploy configmaps

Default installation of WSO2AM Operator supportes load balancer type. The configurations of WSO2 servers are configured for Load Balancer type. Hence we need to deploy relevant configuration changes for NodePort.

```
kubectl create -f configmaps/
```

#### Deploy pattern-1 in NodePort

```
kubectl create -f wso2-apim.yaml
```

#### Access API Manager

```
kubectl get svc

Output:
NAME                                TYPE          CLUSTER-IP   EXTERNAL-IP     PORT(S)                                                        AGE                                                                          

mysql-svc                         ClusterIP      10.0.45.192   <none>          3306/TCP                                                         10m                                         
wso2-am-1-svc                     ClusterIP      10.0.42.214   <none>          9611/TCP,9711/TCP,5672/TCP,9443/TCP                              10m                                          
wso2-am-2-svc                     ClusterIP      10.0.38.247   <none>          9611/TCP,9711/TCP,5672/TCP,9443/TCP                              10m                                          
wso2-am-analytics-dashboard-svc   NodePort       10.0.47.57    <none>          9643:32201/TCP                                                   10m                                         
wso2-am-analytics-worker-svc      ClusterIP      10.0.34.39    <none>          7612/TCP,7712/TCP,9444/TCP,9091/TCP,7071/TCP,7444/TCP            10m                                          
wso2-am-svc                       NodePort       10.0.47.181   <none>          8280:32004/TCP,8243:32003/TCP,9763:32002/TCP,9443:32001/TCP      10m                                             
```

To access the API Manager, add a host mapping entry to the /etc/hosts file. As we have exposed services in Node Port type, you can use the IP address of any Kubernetes node.

```
<Any K8s Node IP>  wso2apim wso2apim-analytics
```

- For Docker for Mac use "127.0.0.1" for the K8s node IP
- For Minikube, use minikube ip command to get the K8s node IP
- For GKE
    ```
    (kubectl get nodes -o jsonpath='{ $.items[*].status.addresses[?(@.type=="ExternalIP")].address }')
    ```
    - This will give the external IPs of the nodes available in the cluster. Pick any IP to include in /etc/hosts file.
  
- **API Publisher** : https://wso2apim:32001/publisher 
- **API Devportal** : https://wso2apim:32001/devportal 
- **API Analytics Dashboard**   : https://wso2apim-analytics:32201/analytics-dashboard 
