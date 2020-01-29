### Configure Own Custom Pattern

The user has the flexibility to create a completely new Pattern they wish.

They must provide the type for each of the profiles. And then based on the types, a template for that type of profile will be created.
But the user can specify whatever the fields they want to override without taking the default configuratiosn specific to that type.
Also they can add the Newly created configmaps and persistent volume claims. 

Finally via applying that custom resource YAML file, it will automatically start to create all the artifacts for the newly created pattern.

The scenario for Customizing the new pattern can be seen in [Scenario-6](https://github.com/wso2-incubator/wso2am-k8s-operator/tree/master/scenarios/scenario-6)

