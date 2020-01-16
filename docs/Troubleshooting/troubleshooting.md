## Trouble Shooting Guide - API Operator

#### Check logs of wso2am-controller

The controller will be running in "wso2-system" namespace
   
```
kubectl get pods -n wso2-system

Output:
NAME                               READY   STATUS    RESTARTS   AGE
wso2am-controller-75c5b84c-vsp4x   1/1     Running   0          76m

---

kubectl logs wso2am-controller-75c5b84c-vsp4x -n wso2-system

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


If the status says "Err", then most possibly it can be due to configuration issue related to Docker-Hub user. 
Check you have put the appropriate DockerHub image in the _wso2am-k8s-operator/artifacts/install/controller-artifacts/7-controller.yaml_ file.
