### Configure Database Manually

- The user has the option to configure a database deployed somewhere else rather than using the default mysql pod deployment. 

- Set the following flag in [controller configurations](../../artifacts/operator-configs/controller-conf.yaml) to disable the default mysql deployment.

    ```
    use-mysql-pod: "false"
    ```
  
- Update the [configmaps](../../artifacts/api-manager-artifacts) with the database connections or create new configmaps for the deployment.toml or deployment.yaml.
    