### Configure Database Manually

The user has the option to conigure a database deployed somewhere else rather than using the default mysql pod deployment. 
Set the `use-mysql-pod` flag to **false** in [controller configurations](https://github.com/wso2/k8s-wso2am-operator/blob/master/artifacts/install/operator-configs/controller-conf.yaml) to disable the default mysql deployment.
After that, you have to manually configure your database connectivity for necessary [API Manager artifacts](https://github.com/wso2/k8s-wso2am-operator/tree/master/artifacts/install/api-manager-artifacts/pattern-1) for the patterns to work properly.