## Troubleshooting Guide

#### Check logs of wso2am-controller

The controller will be running in "wso2-system" namespace
   
```
>> kubectl get pods -n wso2-system

Output:
NAME                               READY   STATUS    RESTARTS   AGE
wso2am-controller-75c5b84c-vsp4x   1/1     Running   0          76m

---

>> kubectl logs wso2am-controller-75c5b84c-vsp4x -n wso2-system

Output:
W0113 09:00:45.694404       1 client_config.go:543] Neither --kubeconfig nor --master was specified.  Using the inClusterConfig.  This might not work.
I0113 09:00:45.698364       1 controller.go:128] Setting up event handlers
I0113 09:00:45.698636       1 controller.go:194] Starting Apimanager controller
I0113 09:00:45.698834       1 controller.go:198] Waiting for informer caches to sync
I0113 09:00:45.799345       1 controller.go:203] Starting workers
I0113 09:00:45.799736       1 controller.go:209] Started workers
I0113 09:02:23.393306       1 controller.go:265] Successfully synced 'default/cluster-1'
I0113 09:02:23.393748       1 event.go:281] Event(v1.ObjectReference{Kind:"APIManager", Namespace:"default", Name:"cluster-1", UID:"4a2ea188-374e-481e-99e4-497db9472916", APIVersion:"apim
.wso2.com/v1alpha1", ResourceVersion:"961222", FieldPath:""}): type: 'Normal' reason: 'synced' Apimanager synced successfully
I0113 09:02:23.739969       1 controller.go:265] Successfully synced 'default/cluster-1'

```


#### If your pods keep on restarting, and if you see the following output,

```
>> kubectl get pods 

Output:
NAME                                                      READY   STATUS    RESTARTS   AGE
mysql-cluster-1-5998bf6df7-jvkq6                          1/1     Running   0          6m
wso2-am-1-cluster-1-66cb576bd5-n6szh                      0/1     Running   2          6m
wso2-am-2-cluster-1-5fc65bbf7c-65bp4                      0/1     Running   2          6m
wso2-am-analytics-dashboard-cluster-1-577fbd446d-7tpm2    1/1     Running   0          6m
wso2-am-analytics-worker-cluster-1-5fb86bf777-dwgkw       1/1     Running   0          6m
```
Then,
- Check whether the cluster is created with enough resources specified
- Try to delete the pod using the following command
   ```
   >> kubectl delete pod <POD_NAME>
   ```
 
- See the logs of the pods, if they throw any errors, using the command,
   ```
   >> kubectl logs <POD_NAME>
   ```

