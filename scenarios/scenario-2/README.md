## Scenario-2 : Deploy API Manager Pattern-1 (LoadBalancer Service Type)

In this scenario we are deploying API Manager Pattern-1 by using the following simple yaml definition.

```
apiVersion: apim.wso2.com/v1alpha1
kind: APIManager
metadata:
  name: cluster-1
spec:
  pattern: Pattern-1
```

### Installation Prerequisites

* [Helm](https://helm.sh/docs/intro/install/) v3.1.0 or above

This is used to deploy the NFS Provisioner and this used for persistent volumes.

##### Setup persistent volumes

```
helm repo add stable https://kubernetes-charts.storage.googleapis.com
helm install nfs-provisioner stable/nfs-server-provisioner
```

<details><summary> Without using Helm, you can also setup NFS manually, click here</summary>
<p>

**Prerequisites**
 * A pre-configured Network File System (NFS) to be used as the persistent volume for artifact sharing and persistence. In the NFS server instance, create a Linux system user account named wso2carbon with user id 802 and a system group named wso2 with group id 802. Add the wso2carbon user to the group wso2.

```
groupadd --system -g 802 wso2
useradd --system -g 802 -u 802 wso2carbon 
```
    
1.Setup a Network File System (NFS) to be used for persistent storage.

Create and export unique directories within the NFS server instance for each Kubernetes Persistent Volume resource defined in the <KUBERNETES_HOME>/artifacts/persistent-volumes/persistent-volume-for-external-nfs.yaml file.

2.Grant ownership to wso2carbon user and wso2 group, for each of the previously created directories. 

```
sudo chown -R wso2carbon:wso2 <directory_name>
```

3.Grant read-write-execute permissions to the wso2carbon user, for each of the previously created directories.

```
chmod -R 700 <directory_name>
```

4.Update the StorageClassName in the <KUBERNETES_HOME>/artifacts/persistent-volumes/storage-class.yaml file as you want.

Then, apply the following command to create a new Storage Class,

```
kubectl create -f <KUBERNETES_HOME>/artifacts/persistent-volumes/storage-class.yaml 
```

5.Update each Kubernetes Persistent Volume resource with the corresponding Namespace (NAME_SPACE), NFS server IP (NFS_SERVER_IP) and exported, NFS server directory path (NFS_LOCATION_PATH) in the <KUBERNETES_HOME>/artifacts/persistent-volumes/persistent-volume-for-external-nfs.yaml file.
      
Then, deploy the persistent volume resource as follows,

```
kubectl create -f <KUBERNETES_HOME>/artifacts/persistent-volumes/persistent-volume-for-external-nfs.yaml -n <USER-NAMESPACE>
```

6.Update PVC Configmap with the corresponding StorageClassName in the <KUBERNETES_HOME>/artifacts/operator-configs/pvc-config.yaml file.

<h3>Using Minikube Hostpath</h3>

Minikube runs a single-node Kubernetes cluster inside a Virtual Machine. Therefore the accessmode ReadWriteMany does not support it, and only ReadWriteOnce supports it. Therefore it is preferred to use other storage methods rather than mnikube hostpath. 

1. Log into Minikube Filesystem via the command,

```
    minikube ssh
```

2. Create unique directories within the Minikube filesystem for each Kubernetes Persistent Volume resource defined in the <KUBERNETES_HOME>/artifacts/persistent-volumes/pv-hostpath.yaml file.

3. Grant permission to mysql directory using the command,
```
    sudo chown 999:999 <mysq-direactory-path>
```
4. Then, deploy the persistent volumes as follows,
```
    kubectl create -f <KUBERNETES_HOME>/artifacts/persistent-volumes/pv-hostpath.yaml -n <USER-NAMESPACE>
```
That is all, Now run the flow in order from start.

</p>
</details>

-----

##### Deploy pattern 1

```
kubectl create -f wso2-apim.yaml
```

#### Access API Manager

By default, we are exposing the services as LoadBalancer service type. Kubernetes cluster provider will create a load balancer for your service.

```
kubectl get svc

NAME                                     TYPE           CLUSTER-IP    EXTERNAL-IP    PORT(S)                                                                                     AGE

mysql-svc                                ClusterIP      10.0.13.146   <none>         3306/TCP                                                                                    5m54s
wso2-am-1-svc                            ClusterIP      10.0.5.21     <none>         9611/TCP,9711/TCP,5672/TCP,9443/TCP                                                         5m54s
wso2-am-2-svc                            ClusterIP      10.0.14.65    <none>         9611/TCP,9711/TCP,5672/TCP,9443/TCP                                                         5m54s
wso2-am-analytics-dashboard-svc          LoadBalancer   10.0.4.120    34.93.83.85    9643:32126/TCP                                                                              5m54s
wso2-am-analytics-worker-svc             LoadBalancer   10.0.1.233    35.200.145.7   7612:32328/TCP,7712:31382/TCP,9444:30857/TCP,9091:30081/TCP,7071:30080/TCP,7444:32321/TCP   5m54s
wso2-am-svc                              LoadBalancer   10.0.13.20    34.93.86.89    8280:31992/TCP,8243:30819/TCP,9763:31620/TCP,9443:32739/TCP                                 5m54s            
```

- For accessing the portals via LoadBalancer, add /etc/host entries as follows.
 
   Copy the <EXTERNAL-IP> of wso2-am-svc & wso2-am-analytics-dashboard-svc.
   Add it to **/etc/hosts** file as below,
    
    ```
    <EXTERNAL-IP-OF-WSO2-AM-SVC>                        wso2apim
    <EXTERNAL-IP-OF-WSO2-AM-ANALYTICS-DASHBOARD-SVC>    wso2apim-analytics 
    ```

- **API Publisher**             : https://wso2apim:9443/publisher 
- **API Devportal**             : https://wso2apim:9443/devportal 
- **API Analytics Dashboard**   : https://wso2apim-analytics:9643/analytics-dashboard 

