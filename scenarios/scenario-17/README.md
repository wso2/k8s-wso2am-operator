## Scenario-17 : Override ConfigMaps

You can override the default deployemnt configmaps present for each profile.

1. Create a new deployment configmap using the deployment.toml/ deployment.yaml files. And apply the command.

    ```
      kubectl create configmap <NEW-PROFILE-CONFIGMAP> --from-file=<PATH_TO_DEPLOYMENT.TOML>
    ```

2. Update the custom resources with created configmap.

    Replace, the **deploymentConfigmap** with configmap name.

3. Deploy updated pattern-4.

    ```
      kubectl apply -f wso2-apim.yaml
    ```
