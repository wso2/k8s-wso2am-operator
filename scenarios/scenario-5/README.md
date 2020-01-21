## Scenario-5 : Add new configmaps and Persistent Volume Claims

You can create new configmap and persistant volume claims with new Mountpath and made it available in the cluster.


1. First create a new configmap for any of the profile
   Then add them under **newConfigmap** field with mountpath as specified in the wso2-apim.yaml file.
   Use the below command to create a configmap
   
   ```
      kubectl create configmap <NEW_CONFIGMAP>
   ```
   
 2. Add them under newConfigmap field of wso2-apim.yaml file.
 
 3. Then apply the given wso2-apim.yaml file.
 
    ```
      kubectl apply -f scenarios/scenario-5/wso2-apim.yaml
    ```
    
    
Similarly, you can create new persistent volume claims for desired profiles with mountpath and add them under **newClaim** field.


    
