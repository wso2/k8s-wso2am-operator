## Scenario-6 : Override ConfigMaps and PersistentVolumeClaims

You can override the default deployemnt configmaps and persistent volume claims present for each profile.

1. Create a new deployment configmap using the deployment.toml/ deployment.yaml files. And apply the command.

    ```
      kubectl create configmap <NEW-PROFILE-CONFIGMAP> --from-file=<PATH_TO_DEPLOYMENT.TOML>
    ```

2. Update the custom resources with created configmap.

    Replace, the **deploymentConfigmap** with configmap name.

3. Deploy updated pattern-1.

    ```
      kubectl apply -f wso2-apim.yaml
    ```

Similarly, you can create the persistent volume claims for synapse configs and execution plans. Add the newly created pvc names in relevant places of **synapseConfigs** and **executionPlans**.

    persistentVolumeClaim:
      synapseConfigs: <NEW_PROFILE_SYNAPSE_CONFIGS_PVC>
      executionPlans: <NEW_PROFILE_EXECUTION_PLANS_PVC>
