### Using Different Volume types

User has the flexibility to choose between different kinds of Volume Types.

They are,
- Internal NFS
- External NFS
- HostPath in Minikube

Using Internal NFS is the simplest thing user can try without much effort.
As user just need to apply the below command with preinstalled Helm,

```
helm install stable/nfs-server-provisioner
```  

Secondly, Configuring the External NFS can be done by following the **Advanced** steps in [Quick-Start-Guide](https://github.com/wso2-incubator/wso2am-k8s-operator)
or by following the [scenario-7](https://github.com/wso2-incubator/wso2am-k8s-operator/tree/master/scenarios/scenario-7)

Finally, for Minikube users, it much more appropriate to use HostPaths as a volume type for the required persistent volumes.
Configuring HostPath also given in the **Advanced** steps in [Quick-Start-Guide](https://github.com/wso2-incubator/wso2am-k8s-operator).

