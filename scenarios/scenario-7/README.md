## Scenario-7 : Add new configmaps

You can create new configmaps with new Mountpath and made it available in the cluster.

1. Create a new configmap for any of the profiles and then add under **newConfigmap** field with mountpath as specified in the wso2-apim.yaml file.
   
   ```
   kubectl create configmap <NEW_CONFIGMAP>
   ```
   
 2. Add them under newConfigmap field of wso2-apim.yaml file.
 
    ```
     configMaps:
       newConfigMap:
         - name: <NEW_CONFIGMAP_1>
           mountPath: <NEW_MOUNTPATH_FOR_CONFIGMAP_1>
    ```
        
 3. Apply the given wso2-apim.yaml file.
 
    ```
    kubectl apply -f scenarios/scenario-5/wso2-apim.yaml
    ```