## Building the project from scratch

* The project was initiated by following the structure of [Sample Controller](https://github.com/kubernetes/sample-controller) provided by Kubernetes.
  
* The workflow of the sample controller is shown below.
  ![sample controller workflow](../images/controller-workflow.jpeg)

* More information regarding this can be found [here](https://medium.com/speechmatics/how-to-write-kubernetes-custom-controllers-in-go-8014c4a04235).

* Introduced a new type called APIManager in [`types.go`](https://github.com/wso2/K8s-wso2am-operator/blob/master/pkg/apis/apim/v1alpha1/types.go).

* Modified the [`hack/update-codegen.sh`](https://github.com/wso2/K8s-wso2am-operator/blob/master/hack/update-codegen.sh) file to generate the **generated** directory with clientset, informers and listers.

* Then developed the go code to create deployments, services and ingresses for wso2 apim pattern-1 artifacts.

* [`main.go`](https://github.com/wso2/K8s-wso2am-operator/blob/master/cmd/controller/main.go) has the command to run the controller. It will create the binary image of the controller.

* [`controller.go`](https://github.com/wso2/K8s-wso2am-operator/blob/master/pkg/controller/controller.go) has the actual logic to run once the custom resource triggers it. The components are separated into their own go files for easy access and invoking in the relevant places inside the controller.go file.

* [`Makefile`](https://github.com/wso2/K8s-wso2am-operator/blob/master/Makefile) has the commands to build and run the controller. 

* The `make all` command make all will run the controller even if we don’t pull the controller from Docker Hub.


## Building the Docker image

* `Dockerfile` is created with relevant commands to run the controller executable. So if we push that to our own Docker Hub account, then simply we can pull and run it on any platform.

* For that you need to have a Docker Hub account. https://hub.docker.com/

* The commands to build the docker image is as follows. USERNAME is the DockerHub account username, and TAGNAME can be anything and we build the docker image of wso2am-controller.

```
    docker build -t USERNAME/wso2am-controller:TAGNAME .
    sudo docker push USERNAME/wso2am-controller:TAGNAME
```

* The following is the output of the above commands.

```
    Sending build context to Docker daemon  491.7MB
    Step 1/3 : FROM scratch
    ---> 
    Step 2/3 : COPY wso2am-controller /
    ---> cd81974b50cb
    Step 3/3 : ENTRYPOINT ["/wso2am-controller","-logtostderr=true"]
    ---> Running in 9deb5f5ab0da
    Removing intermediate container 9deb5f5ab0da
    ---> f6083fe77956
    Successfully built f6083fe77956
    Successfully tagged wso2/wso2am-controller:latest

```
```
    The push refers to repository [docker.io/wso2/wso2am-controller]
    affab6de580c: Pushed 

```

