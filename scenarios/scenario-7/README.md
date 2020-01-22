<h3>Running External-nfs</h3>

**Prerequisites**
 * A pre-configured Network File System (NFS) to be used as the persistent volume for artifact sharing and persistence. In the NFS server instance, create a Linux system user account named wso2carbon with user id 802 and a system group named wso2 with group id 802. Add the wso2carbon user to the group wso2.

```
    groupadd --system -g 802 wso2
    useradd --system -g 802 -u 802 wso2carbon 
```
    
1.Setup a Network File System (NFS) to be used for persistent storage.
Create and export unique directories within the NFS server instance for each Kubernetes Persistent Volume resource     defined in the <KUBERNETES_HOME>/scenarios/scenario-7/pv.yaml file.

2.Grant ownership to wso2carbon user and wso2 group, for each of the previously created directories. 

```
    sudo chown -R wso2carbon:wso2 <directory_name>
```

3.Grant read-write-execute permissions to the wso2carbon user, for each of the previously created directories.

```
    chmod -R 700 <directory_name>
```

4.Update the StorageClassName in the <KUBERNETES_HOME>/scenarios/scenario-7/storage-class.yaml file as you want.

Then, apply the following command to create a new Storage Class,

```
    kubectl create -f <KUBERNETES_HOME>/scenarios/scenario-7/storage-class.yaml 
```

5.Update each Kubernetes Persistent Volume resource with the corresponding Namespace (NAME_SPACE), NFS server IP (NFS_SERVER_IP) and exported, NFS server directory path (NFS_LOCATION_PATH) in the <KUBERNETES_HOME>/scenarios/scenario-7/pv.yaml file.
      
Then, deploy the persistent volume resource as follows,

```
    kubectl create -f <KUBERNETES_HOME>/scenarios/scenario-7/pv.yaml -n <USER-NAMESPACE>
```

6.Update PVC Configmap with the corresponding StorageClassName in the <KUBERNETES_HOME>/artifacts/install/operator-configs/pvc-config.yaml file.
