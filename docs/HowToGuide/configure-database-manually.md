### Configure Database Manually

- The user has the option to configure a database deployed somewhere else rather than using the default mysql pod deployment. 

- Set the following **Spec** in **APIManager** Resource to disable the default mysql deployment.
    ```yaml
    useMysql: "false"
    ```
  
    Sample usage:
    ```yaml
    apiVersion: apim.wso2.com/v1alpha1
    kind: APIManager
    metadata:
      name: custom-pattern-1
    spec:
      pattern: Pattern-X
      service:
        type: NodePort
      useMysql: "false"
      profiles:
        - name: all-in-one-api-manager
          type: api-manager
          service:
            name: wso2apim
          deployment:
            configMaps:
              deploymentConfigMap: apim-conf
    ```