### Using Different Persistent Volumes

- User has the flexibility to choose between different kinds of Volume Types.

    - Internal NFS
    
        Using Internal NFS is the simplest approach and user can use Helm v3 to deploy NFS as follows.
        
        ```
        helm repo add stable https://kubernetes-charts.storage.googleapis.com
        helm install nfs-provisioner stable/nfs-server-provisioner
        ```  
    
    - External NFS

        Checkout [Scenario 8](../../scenarios/scenario-8) for running external NFS.
