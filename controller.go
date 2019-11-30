/*
 * Copyright (c) 2019 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http:www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

import (
	"fmt"
	//"k8s.io/apimachinery/pkg/api/resource"
	"time"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	apps2informers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/apimachinery/pkg/util/intstr"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appslisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
	app2listers "k8s.io/client-go/listers/core/v1"
	apimv1alpha1 "github.com/keshiha96/wso2am-k8s-controller/pkg/apis/apim/v1alpha1"
	clientset "github.com/keshiha96/wso2am-k8s-controller/pkg/generated/clientset/versioned"
	samplescheme "github.com/keshiha96/wso2am-k8s-controller/pkg/generated/clientset/versioned/scheme"
	informers "github.com/keshiha96/wso2am-k8s-controller/pkg/generated/informers/externalversions/apim/v1alpha1"
	listers "github.com/keshiha96/wso2am-k8s-controller/pkg/generated/listers/apim/v1alpha1"
)

const controllerAgentName = "wso2am-controller"

const (
	// SuccessSynced is used as part of the Event 'reason' when a Apimanager is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a Apimanager fails
	// to sync due to a Deployment of the same name already existing.
	ErrResourceExists = "ErrResourceExists"
	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Deployment already existing
	MessageResourceExists = "Resource %q already exists and is not managed by Apimanager"
	// MessageResourceSynced is the message used for an Event fired when a Apimanager
	// is synced successfully
	MessageResourceSynced = "Apimanager synced successfully"

)

// Controller is the controller implementation for Apimanager resources
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// sampleclientset is a clientset for our own API group
	sampleclientset clientset.Interface

	deploymentsLister appslisters.DeploymentLister
	servicesLister    app2listers.ServiceLister
	deploymentsSynced cache.InformerSynced
	servicesSynced    cache.InformerSynced
	apimanagerslister listers.ApimanagerLister
	apimanagersSynced cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

// NewController returns a new wso2-apim controller
func NewController(
	kubeclientset kubernetes.Interface,
	sampleclientset clientset.Interface,
	deploymentInformer appsinformers.DeploymentInformer,
	serviceInformer apps2informers.ServiceInformer,
	apimanagerInformer informers.ApimanagerInformer) *Controller {

	// Create event broadcaster
	// Add apim-controller types to the default Kubernetes Scheme so Events can be
	// logged for apim-controller types.
	utilruntime.Must(samplescheme.AddToScheme(scheme.Scheme))
	klog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeclientset:     kubeclientset,
		sampleclientset:   sampleclientset,
		deploymentsLister: deploymentInformer.Lister(),
		deploymentsSynced: deploymentInformer.Informer().HasSynced,
		servicesLister:    serviceInformer.Lister(),
		servicesSynced:    serviceInformer.Informer().HasSynced,
		apimanagerslister: apimanagerInformer.Lister(),
		apimanagersSynced: apimanagerInformer.Informer().HasSynced,
		workqueue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Apimanagers"),
		recorder:          recorder,
	}

	klog.Info("Setting up event handlers")
	// Set up an event handler for when Apimanager resources change
	apimanagerInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueApimanager,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueApimanager(new)
		},
	})
	// Set up an event handler for when Deployment resources change. This
	// handler will lookup the owner of the given Deployment, and if it is
	// owned by a Apimanager resource will enqueue that Apimanager resource for
	// processing. This way, we don't need to implement custom logic for
	// handling Deployment resources. More info on this pattern:
	// https://github.com/kubernetes/community/blob/8cafef897a22026d42f5e5bb3f104febe7e29830/contributors/devel/controllers.md
	deploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.handleObject,
		UpdateFunc: func(old, new interface{}) {
			newDepl := new.(*appsv1.Deployment)
			oldDepl := old.(*appsv1.Deployment)
			if newDepl.ResourceVersion == oldDepl.ResourceVersion {
				// Periodic resync will send update events for all known Deployments.
				// Two different versions of the same Deployment will always have different RVs.
				return
			}
			controller.handleObject(new)
		},
		DeleteFunc: controller.handleObject,
	})

	serviceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.handleObject,
		UpdateFunc: func(old, new interface{}) {
			newServ := new.(*corev1.Service)
			oldServ := old.(*corev1.Service)
			if newServ.ResourceVersion == oldServ.ResourceVersion {
				// Periodic resync will send update events for all known Services.
				// Two different versions of the same Service will always have different RVs.
				return
			}
			controller.handleObject(new)
		},
		DeleteFunc: controller.handleObject,
	})

	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	klog.Info("Starting Apimanager controller")

	// Wait for the caches to be synced before starting workers
	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.deploymentsSynced, c.apimanagersSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("Starting workers")
	// Launch two workers to process Apimanager resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer/postpone c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// Apimanager resource to be synced.
		if err := c.syncHandler(key); err != nil {
			// Put the item back on the workqueue to handle any transient/non-permanent errors.
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

// c is the Controller object type pointer as a parameter
// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Apimanager resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the Apimanager resource with this namespace/name
	apimanager, err := c.apimanagerslister.Apimanagers(namespace).Get(name)
	if err != nil {
		// The Apimanager resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("apimanager '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}

	//////////////////////////////////////////////////////////////


	apim1deploymentName := "apim-1-deploy"
	apim2deploymentName := "apim-2-deploy"
	apim1serviceName := "apim-1-svc"
	apim2serviceName := "apim-2-svc"
	mysqldeploymentName :="mysql-deploy"
	mysqlserviceName := "wso2am-mysql-db-service"
	//dashboardDeploymentName := "analytics-dash-deploy"
	//dashboardServiceName := "analytics-dash-svc"
	//workerDeploymentName := "analytics-worker-deploy"
	//workerServiceName := "analytics-worker-svc"



	/////////check whether resourecs already exits, else create one

	// Get the deployment using hardcoded deployment name wso2-apim-1-deploy
	deployment, err := c.deploymentsLister.Deployments(apimanager.Namespace).Get(apim1deploymentName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		deployment, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Create(apim1Deployment(apimanager))
		if err != nil {
			return err
		}
	}

	// Get the deployment using hardcoded deployment name wso2-apim-1-deploy2
	deployment2, err := c.deploymentsLister.Deployments(apimanager.Namespace).Get(apim2deploymentName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		deployment2, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Create(apim2Deployment(apimanager))
		if err != nil {
			return err
		}
	}

	//// Get the dashboard deployment using hardcoded deployment name wso2-apim-1-deploy
	//dashdeployment, err := c.deploymentsLister.Deployments(apimanager.Namespace).Get(dashboardDeploymentName)
	//// If the resource doesn't exist, we'll create it
	//if errors.IsNotFound(err) {
	//	dashdeployment, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Create(dashboardDeployment(apimanager))
	//	if err != nil {
	//		return err
	//	}
	//}
	//
	//// Get the dashboard deployment using hardcoded deployment name wso2-apim-1-deploy
	//workerdeployment, err := c.deploymentsLister.Deployments(apimanager.Namespace).Get(workerDeploymentName)
	//// If the resource doesn't exist, we'll create it
	//if errors.IsNotFound(err) {
	//	workerdeployment, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Create(workerDeployment(apimanager))
	//	if err != nil {
	//		return err
	//	}
	//}

	// Get the deployment using hardcoded deployment name mysqldeployment
	mysqldeployment, err := c.deploymentsLister.Deployments(apimanager.Namespace).Get(mysqldeploymentName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		mysqldeployment, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Create(mysqlDeployment(apimanager))
		if err != nil {
			return err
		}
	}



	// Get the service with the name specified in wso2-apim spec
	service, err := c.servicesLister.Services(apimanager.Namespace).Get(apim1serviceName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		service, err = c.kubeclientset.CoreV1().Services(apimanager.Namespace).Create(apim1Service(apimanager))
	}

	// Get the service with the name specified in wso2-apim spec
	service2, err := c.servicesLister.Services(apimanager.Namespace).Get(apim2serviceName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		service2, err = c.kubeclientset.CoreV1().Services(apimanager.Namespace).Create(apim2Service(apimanager))
	}

	//// Get the service with the name specified in wso2-apim spec
	//dashservice, err := c.servicesLister.Services(apimanager.Namespace).Get(dashboardServiceName)
	//// If the resource doesn't exist, we'll create it
	//if errors.IsNotFound(err) {
	//	dashservice, err = c.kubeclientset.CoreV1().Services(apimanager.Namespace).Create(dashboardService(apimanager))
	//}
	//
	//// Get the service with the name specified in wso2-apim spec
	//workerservice, err := c.servicesLister.Services(apimanager.Namespace).Get(workerServiceName)
	//// If the resource doesn't exist, we'll create it
	//if errors.IsNotFound(err) {
	//	workerservice, err = c.kubeclientset.CoreV1().Services(apimanager.Namespace).Create(workerService(apimanager))
	//}

	// Get the service with the name specified in wso2-apim spec
	mysqlservice, err := c.servicesLister.Services(apimanager.Namespace).Get(mysqlserviceName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		mysqlservice, err = c.kubeclientset.CoreV1().Services(apimanager.Namespace).Create(mysqlService(apimanager))
	}

	// If an error occurs during Get/Create, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	/////////////check whether resources are controlled by apimanager with same owner reference

	// If the Deployment is not controlled by this Apimanager resource, we should log
	// a warning to the event recorder and ret
	if !metav1.IsControlledBy(deployment, apimanager) {
		msg := fmt.Sprintf("Deployment1 %q already exists and is not managed by Apimanager", deployment.Name)
		c.recorder.Event(apimanager, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}

	//for instance 2
	// If the Deployment2 is not controlled by this Apimanager resource, we should log
	// a warning to the event recorder and ret
	if !metav1.IsControlledBy(deployment2, apimanager) {
		msg := fmt.Sprintf("Deployment2 %q already exists and is not managed by Apimanager", deployment2.Name)
		c.recorder.Event(apimanager, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}

	////for analytics dashboard
	//// If the Deployment2 is not controlled by this Apimanager resource, we should log
	//// a warning to the event recorder and ret
	//if !metav1.IsControlledBy(dashdeployment, apimanager) {
	//	msg := fmt.Sprintf("Analytics Dashboard Deployment %q already exists and is not managed by Apimanager", dashdeployment.Name)
	//	c.recorder.Event(apimanager, corev1.EventTypeWarning, ErrResourceExists, msg)
	//	return fmt.Errorf(msg)
	//}
	//
	////for analytics worker
	//// If the Deployment2 is not controlled by this Apimanager resource, we should log
	//// a warning to the event recorder and ret
	//if !metav1.IsControlledBy(workerdeployment, apimanager) {
	//	msg := fmt.Sprintf("Analytics Dashboard Deployment %q already exists and is not managed by Apimanager", workerdeployment.Name)
	//	c.recorder.Event(apimanager, corev1.EventTypeWarning, ErrResourceExists, msg)
	//	return fmt.Errorf(msg)
	//}

	//for mysql deployment
	// If the mysql Deployment is not controlled by this Apimanager resource, we should log
	// a warning to the event recorder and ret
	if !metav1.IsControlledBy(mysqldeployment, apimanager) {
		msg := fmt.Sprintf("mysql deployment %q already exists and is not managed by Apimanager", mysqldeployment.Name)
		c.recorder.Event(apimanager, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}

	// If the instance 1 Service is not controlled by this Apimanager resource, we should log
	// a warning to the event recorder and ret
	if !metav1.IsControlledBy(service, apimanager) {
		msg := fmt.Sprintf("Service %q already exists and is not managed by Apimanager", service.Name)
		c.recorder.Event(apimanager, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}

	//for instance 2 service
	// If the Service is not controlled by this Apimanager resource, we should log
	// a warning to the event recorder and ret
	if !metav1.IsControlledBy(service2, apimanager) {
		msg := fmt.Sprintf("Service %q already exists and is not managed by Apimanager", service2.Name)
		c.recorder.Event(apimanager, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}

	////for analytics dashboard service
	//// If the Service is not controlled by this Apimanager resource, we should log
	//// a warning to the event recorder and ret
	//if !metav1.IsControlledBy(dashservice, apimanager) {
	//	msg := fmt.Sprintf("Service %q already exists and is not managed by Apimanager", dashservice.Name)
	//	c.recorder.Event(apimanager, corev1.EventTypeWarning, ErrResourceExists, msg)
	//	return fmt.Errorf(msg)
	//}
	//
	////for analytics dashboard service
	//// If the Service is not controlled by this Apimanager resource, we should log
	//// a warning to the event recorder and ret
	//if !metav1.IsControlledBy(workerservice, apimanager) {
	//	msg := fmt.Sprintf("Service %q already exists and is not managed by Apimanager", workerservice.Name)
	//	c.recorder.Event(apimanager, corev1.EventTypeWarning, ErrResourceExists, msg)
	//	return fmt.Errorf(msg)
	//}

	//for mysql service
	// If the Service is not controlled by this Apimanager resource, we should log
	// a warning to the event recorder and ret
	if !metav1.IsControlledBy(mysqlservice, apimanager) {
		msg := fmt.Sprintf("Service %q already exists and is not managed by Apimanager", mysqlservice.Name)
		c.recorder.Event(apimanager, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}

	///////////check replicas are same as defined for deployments

	// If this number of the replicas on the Apimanager resource is specified, and the
	// number does not equal the current desired replicas on the Deployment, we
	// should update the Deployment resource.
	if apimanager.Spec.Replicas != nil && *apimanager.Spec.Replicas != *deployment.Spec.Replicas {
		klog.V(4).Infof("Apimanager %s replicas: %d, deployment replicas: %d", name, *apimanager.Spec.Replicas, *deployment.Spec.Replicas)
		deployment, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Update(apim1Deployment(apimanager))
	}

	//for instance 2 also
	// If this number of the replicas on the Apimanager resource is specified, and the
	// number does not equal the current desired replicas on the Deployment2, we
	// should update the Deployment2 resource.
	if apimanager.Spec.Replicas != nil && *apimanager.Spec.Replicas != *deployment2.Spec.Replicas {
		klog.V(4).Infof("Apimanager %s replicas: %d, deployment2 replicas: %d", name, *apimanager.Spec.Replicas, *deployment2.Spec.Replicas)
		deployment2, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Update(apim2Deployment(apimanager))
	}

	////for analytics dashboard deployment
	//// If this number of the replicas on the Apimanager resource is specified, and the
	//// number does not equal the current desired replicas on the Deployment2, we
	//// should update the Deployment2 resource.
	//if apimanager.Spec.Replicas != nil && *apimanager.Spec.Replicas != *dashdeployment.Spec.Replicas {
	//	klog.V(4).Infof("Apimanager %s replicas: %d, deployment2 replicas: %d", name, *apimanager.Spec.Replicas, *dashdeployment.Spec.Replicas)
	//	dashdeployment, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Update(dashboardDeployment(apimanager))
	//}
	//
	////for analytics worker deployment
	//// If this number of the replicas on the Apimanager resource is specified, and the
	//// number does not equal the current desired replicas on the Deployment2, we
	//// should update the Deployment2 resource.
	//if apimanager.Spec.Replicas != nil && *apimanager.Spec.Replicas != *workerdeployment.Spec.Replicas {
	//	klog.V(4).Infof("Apimanager %s replicas: %d, deployment2 replicas: %d", name, *apimanager.Spec.Replicas, *workerdeployment.Spec.Replicas)
	//	dashdeployment, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Update(workerDeployment(apimanager))
	//}

	//for instance mysql deployment
	// If this number of the replicas on the Apimanager resource is specified, and the
	// number does not equal the current desired replicas on the mysql deployment, we
	// should update the mysql Deployment resource.
	if apimanager.Spec.Replicas != nil && *apimanager.Spec.Replicas != *mysqldeployment.Spec.Replicas {
		klog.V(4).Infof("Apimanager %s replicas: %d, deployment2 replicas: %d", name, *apimanager.Spec.Replicas, *mysqldeployment.Spec.Replicas)
		mysqldeployment, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Update(mysqlDeployment(apimanager))
	}

	// If an error occurs during Update, we'll requeue the item so we can
	// attempt processing again later. THis could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}


	//////////finally update the deployment resources after done checking

	// Finally, we update the status block of the Apimanager resource to reflect the
	// current state of the world
	err = c.updateApimanagerStatus(apimanager, deployment)
	if err != nil {
		return err
	}

	//for instance 2 also
	// Finally, we update the status block of the Apimanager resource to reflect the
	// current state of the world
	err = c.updateApimanagerStatus(apimanager, deployment2)
	if err != nil {
		return err
	}

	////for analytics dashboard deployment
	//// Finally, we update the status block of the Apimanager resource to reflect the
	//// current state of the world
	//err = c.updateApimanagerStatus(apimanager, dashdeployment)
	//if err != nil {
	//	return err
	//}
	//
	////for analytics worker deployment
	//// Finally, we update the status block of the Apimanager resource to reflect the
	//// current state of the world
	//err = c.updateApimanagerStatus(apimanager, workerdeployment)
	//if err != nil {
	//	return err
	//}

	//for mysql deployment
	// Finally, we update the status block of the Apimanager resource to reflect the
	// current state of the world
	err = c.updateApimanagerStatus(apimanager, mysqldeployment)
	if err != nil {
		return err
	}

	c.recorder.Event(apimanager, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

func (c *Controller) updateApimanagerStatus(apimanager *apimv1alpha1.Apimanager, deployment *appsv1.Deployment) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	apimanagerCopy := apimanager.DeepCopy()
	apimanagerCopy.Status.AvailableReplicas = deployment.Status.AvailableReplicas
	// If the CustomResourceSubresources feature gate is not enabled,
	// we must use Update instead of UpdateStatus to update the Status block of the Apimanager resource.
	// UpdateStatus will not allow changes to the Spec of the resource,
	// which is ideal for ensuring nothing other than resource status has been updated.
	_, err := c.sampleclientset.ApimV1alpha1().Apimanagers(apimanager.Namespace).Update(apimanagerCopy)
	return err
}

// enqueueApimanager takes a Apimanager resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than Apimanager.
func (c *Controller) enqueueApimanager(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}

// handleObject will take any resource implementing metav1.Object and attempt
// to find the Apimanager resource that 'owns' it. It does this by looking at the
// objects metadata.ownerReferences field for an appropriate OwnerReference.
// It then enqueues that Apimanager resource to be processed. If the object does not
// have an appropriate OwnerReference, it will simply be skipped.
func (c *Controller) handleObject(obj interface{}) {
	var object metav1.Object

	var ok bool
	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object, invalid type"))
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
			return
		}
		klog.V(4).Infof("Recovered deleted object '%s' from tombstone", object.GetName())
	}
	klog.V(4).Infof("Processing object: %s", object.GetName())
	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
		// If this object is not owned by a Apimanager, we should not do anything more
		// with it.
		if ownerRef.Kind != "Apimanager" {
			return
		}

		apimanager, err := c.apimanagerslister.Apimanagers(object.GetNamespace()).Get(ownerRef.Name)
		if err != nil {
			klog.V(4).Infof("ignoring orphaned object '%s' of apimanager '%s'", object.GetSelfLink(), ownerRef.Name)
			return
		}

		c.enqueueApimanager(apimanager)
		return
	}
}

// apim1Deployment creates a new Deployment for a Apimanager instance 1 resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the Apimanager resource that 'owns' it.
func apim1Deployment(apimanager *apimv1alpha1.Apimanager) *appsv1.Deployment {
	//apim1cpu, _ :=resource.ParseQuantity("2000m")
	//apim1mem, _ := resource.ParseQuantity("2Gi")
	//apim1cpu2, _ :=resource.ParseQuantity("3000m")
	//apim1mem2, _ := resource.ParseQuantity("3Gi")
	labels := map[string]string{
		"deployment": "wso2am-pattern-1-am",
		"node": "wso2am-pattern-1-am-1",
	}

	//defaultmode := int32(0407)
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			// Name: apimanager.Spec.DeploymentName,
			Name:      "apim-1-deploy",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: apimanager.Spec.Replicas,
			MinReadySeconds:240,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.DeploymentStrategyType(appsv1.RollingUpdateDaemonSetStrategyType),
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 1,
					},
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 0,
					},
				},
			},

			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					HostAliases: []corev1.HostAlias{
						{
							IP: "127.0.0.1",
							Hostnames: []string{
								"wso2-am",
								"wso2-gateway",
							},
						},
					},
					//InitContainers: []corev1.Container{
					//	{
					//		Name: "init-apim-analytics-db",
					//		Image: "busybox:1.31",
					//		Command: []string {
					//			//"sh", "-c", "echo -e \"Checking for the availability of MySQL Server deployment\" ; while ! nc -z \"while ! nc -z \"wso2am-mysql-db-service \"3306; do sleep 1; printf \"-\" ;done; echo -e \" >> MySQL Server has started \"; ",
					//			"sh",
					//			"-c",
					//			"echo -e \"Checking for the availability of MySQL Server deployment\"; while ! nc -z \"wso2am-mysql-db-service\" 3306; do sleep 1; printf \"-\"; done; echo -e \"  >> MySQL Server has started\";",
					//		},
					//
					//	},
					//	{
					//		Name: "init-am-analytics-worker",
					//		Image: "busybox:1.31",
					//		Command: []string {
					//			"sh", "-c", "echo -e \"Checking for the availability of WSO2 API Manager Analytics Worker deployment\" ; while ! nc -z \"while ! nc -z \"wso2am-pattern-1-analytics-worker-service 7712; do sleep 1; printf \"-\" ;done; echo -e \" >> WSO2 API Manager Analytics Worker has started \"; ",
					//		},
					//
					//	},
					//},
					Containers: []corev1.Container{
						{
							Name:  "wso2-pattern-1-am",
							Image: "wso2/wso2am:3.0.0",
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec:&corev1.ExecAction{
										Command:[]string{
											"/bin/sh",
											"-c",
											"nc -z localhost 9443",
										},
									},
								},
								InitialDelaySeconds: 240,
								PeriodSeconds:       10,
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec:&corev1.ExecAction{
										Command:[]string{
											"/bin/sh",
											"-c",
											"nc -z localhost 9443",
										},
									},
								},

								InitialDelaySeconds: 240,
								PeriodSeconds:       10,

							},

							Lifecycle: &corev1.Lifecycle{
								PreStop:&corev1.Handler{
									Exec:&corev1.ExecAction{
										Command:[]string{
											"sh",
											"-c",
											"${WSO2_SERVER_HOME}/bin/worker.sh stop",
										},
									},
								},
							},

							//Resources:corev1.ResourceRequirements{
							//	Requests:corev1.ResourceList{
							//		corev1.ResourceCPU:apim1cpu,
							//		corev1.ResourceMemory:apim1mem,
							//	},
							//	Limits:corev1.ResourceList{
							//		corev1.ResourceCPU:apim1cpu2,
							//		corev1.ResourceMemory:apim1mem2,
							//	},
							//},

							ImagePullPolicy: "Always",

							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8280,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 8243,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 9763,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 9443,
									Protocol:      "TCP",
								},


							},
							Env: []corev1.EnvVar{
								// {
								// 	Name:  "HOST_NAME",
								// 	Value: "foo-am",
								// },
								{
									Name: "NODE_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name: "wso2am-pattern-1-am-volume-claim-synapse-configs",
									MountPath: "/home/wso2carbon/wso2-artifact-volume/repository/deployment/server/synapse-configs",
								},
								{
									Name: "wso2am-pattern-1-am-volume-claim-executionplans",
									MountPath:"/home/wso2carbon/wso2-artifact-volume/repository/deployment/server/executionplans",
								},
								{
									Name: "wso2am-pattern-1-am-1-conf",
									MountPath: "/home/wso2carbon/wso2-config-volume/repository/conf/deployment.toml",
									SubPath:"deployment.toml",
								},
								{
									Name: "mysql-jdbc-driver",
									MountPath: "/home/wso2carbon/wso2-artifact-volume/repository/components/lib",
								},
								//{
								//	Name: "wso2am-pattern-1-am-conf-entrypoint",
								//	MountPath: "/home/wso2carbon/docker-entrypoint.sh",
								//	SubPath:"docker-entrypoint.sh",
								//},
							},
						},
					},

					ServiceAccountName: "wso2am-pattern-1-svc-account",
					ImagePullSecrets:[]corev1.LocalObjectReference{
						{
							Name:"wso2am-pattern-1-creds",
						},
					},

					Volumes: []corev1.Volume{
						{
							Name: "wso2am-pattern-1-am-volume-claim-synapse-configs",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName:"pvc-synapse-configs",
								},
							},
						},
						{
							Name: "wso2am-pattern-1-am-volume-claim-executionplans",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "pvc-execution-plans",
								},
							},
						},
						{
							Name: "wso2am-pattern-1-am-1-conf",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "wso2am-pattern-1-am-1-conf",
									},
								},
							},
						},
						//{
						//	Name: "wso2am-pattern-1-am-conf-entrypoint",
						//	VolumeSource: corev1.VolumeSource{
						//		ConfigMap: &corev1.ConfigMapVolumeSource{
						//			LocalObjectReference: corev1.LocalObjectReference{
						//				Name: "wso2am-pattern-1-am-conf-entrypoint",
						//			},
						//			DefaultMode:&defaultmode,
						//		},
						//	},
						//},
						{
							Name: "mysql-jdbc-driver",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "mysql-jdbc-driver-cm",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}


// newService creates a new Service for a Apimanager resource.
// It expose the service with Nodeport type with minikube ip as the externel ip.
func apim1Service(apimanager *apimv1alpha1.Apimanager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2am-pattern-1-am",
		"node": "wso2am-pattern-1-am-1",
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			// Name: apimanager.Spec.ServiceName,
			Name:      "apim-1-svc",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     "NodePort",
			// values are fetched from wso2-apim.yaml file
			// Type: apimanager.Spec.ServType,
			ExternalIPs: []string{"192.168.99.101"},
			// ExternalIPs: apimanager.Spec.ExternalIps,
			Ports: []corev1.ServicePort{
				{
					Name:       "binary",
					Protocol:   corev1.ProtocolTCP,
					Port:       9611,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9611},
					NodePort:   32002,
				},
				{
					Name:       "binary-secure",
					Protocol:   corev1.ProtocolTCP,
					Port:       9711,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9711},
					NodePort:   32003,
				},
				{
					Name:       "jms-tcp",
					Protocol:   corev1.ProtocolTCP,
					Port:       5672,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 5672},
					NodePort:   32004,
				},
				{
					Name:       "servlet-https",
					Protocol:   corev1.ProtocolTCP,
					Port:       9443,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9443},
					NodePort:   32001,
				},
			},
		},
	}
}


// apim1Deployment creates a new Deployment for a Apimanager instance 1 resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the Apimanager resource that 'owns' it.
func apim2Deployment(apimanager *apimv1alpha1.Apimanager) *appsv1.Deployment {
	////apim2cpu, _ :=resource.ParseQuantity("2000m")
	////apim2mem, _ := resource.ParseQuantity("2Gi")
	////apim2cpu2, _ :=resource.ParseQuantity("3000m")
	////apim2mem2, _ := resource.ParseQuantity("3Gi")
	//
	//labels := map[string]string{
	//	//"app":        "wso2am",
	//	//"controller": apimanager.Name,
	//	"deployment": "wso2am-pattern-1-am",
	//	"node": "wso2am-pattern-1-am-2",
	//}
	//return &appsv1.Deployment{
	//	ObjectMeta: metav1.ObjectMeta{
	//		// Name: apimanager.Spec.DeploymentName,
	//		Name:      "apim-2-deploy",
	//		Namespace: apimanager.Namespace,
	//		OwnerReferences: []metav1.OwnerReference{
	//			*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
	//		},
	//	},
	//	Spec: appsv1.DeploymentSpec{
	//		Replicas: apimanager.Spec.Replicas,
	//		MinReadySeconds:240,
	//		Strategy: appsv1.DeploymentStrategy{
	//			Type: appsv1.DeploymentStrategyType(appsv1.RollingUpdateDaemonSetStrategyType),
	//			RollingUpdate: &appsv1.RollingUpdateDeployment{
	//				MaxSurge: &intstr.IntOrString{
	//					Type:   intstr.Int,
	//					IntVal: 1,
	//				},
	//				MaxUnavailable: &intstr.IntOrString{
	//					Type:   intstr.Int,
	//					IntVal: 0,
	//				},
	//			},
	//		},
	//
	//		Selector: &metav1.LabelSelector{
	//			MatchLabels: labels,
	//		},
	//		Template: corev1.PodTemplateSpec{
	//			ObjectMeta: metav1.ObjectMeta{
	//				Labels: labels,
	//			},
	//			Spec: corev1.PodSpec{
	//				HostAliases: []corev1.HostAlias{
	//					{
	//						IP: "127.0.0.1",
	//						Hostnames: []string{
	//							"wso2-am",
	//							"wso2-gateway",
	//						},
	//					},
	//				},
	//				//InitContainers: []corev1.Container{
	//				//	{
	//				//		Name: "init-apim-analytics-db",
	//				//		Image: "busybox:1.31",
	//				//		Command: []string {
	//				//			//"sh", "-c", "echo -e \"Checking for the availability of MySQL Server deployment\" ; while ! nc -z \"while ! nc -z \"wso2am-mysql-db-service \"3306; do sleep 1; printf \"-\" ;done; echo -e \" >> MySQL Server has started \"; ",
	//				//			"sh",
	//				//			"-c",
	//				//			"echo -e \"Checking for the availability of MySQL Server deployment\"; while ! nc -z \"wso2am-mysql-db-service\" 3306; do sleep 1; printf \"-\"; done; echo -e \"  >> MySQL Server has started\";",
	//				//		},
	//				//
	//				//	},
	//				//	{
	//				//		Name: "init-am-analytics-worker",
	//				//		Image: "busybox:1.31",
	//				//		Command: []string {
	//				//			"sh", "-c", "echo -e \"Checking for the availability of WSO2 API Manager Analytics Worker deployment\" ; while ! nc -z \"while ! nc -z \"wso2am-pattern-1-analytics-worker-service 7712; do sleep 1; printf \"-\" ;done; echo -e \" >> WSO2 API Manager Analytics Worker has started \"; ",
	//				//		},
	//				//
	//				//	},
	//				//},
	//				Containers: []corev1.Container{
	//					{
	//						Name:  "wso2-pattern-1-am",
	//						Image: "wso2/wso2am:3.0.0",
	//						LivenessProbe: &corev1.Probe{
	//							Handler: corev1.Handler{
	//								Exec:&corev1.ExecAction{
	//									Command:[]string{
	//										"/bin/sh",
	//										"-c",
	//										"nc -z localhost 9443",
	//									},
	//								},
	//							},
	//							InitialDelaySeconds: 240,
	//							PeriodSeconds:       10,
	//						},
	//						ReadinessProbe: &corev1.Probe{
	//							Handler: corev1.Handler{
	//								Exec:&corev1.ExecAction{
	//									Command:[]string{
	//										"/bin/sh",
	//										"-c",
	//										"nc -z localhost 9443",
	//									},
	//								},
	//							},
	//
	//							InitialDelaySeconds: 240,
	//							PeriodSeconds:       10,
	//
	//						},
	//
	//						Lifecycle: &corev1.Lifecycle{
	//							PreStop:&corev1.Handler{
	//								Exec:&corev1.ExecAction{
	//									Command:[]string{
	//										"sh",
	//										"-c",
	//										"${WSO2_SERVER_HOME}/bin/wso2server.sh stop",
	//									},
	//								},
	//							},
	//						},
	//
	//						//Resources:corev1.ResourceRequirements{
	//						//	Requests:corev1.ResourceList{
	//						//		corev1.ResourceCPU:apim2cpu,
	//						//		corev1.ResourceMemory:apim2mem,
	//						//	},
	//						//	Limits:corev1.ResourceList{
	//						//		corev1.ResourceCPU:apim2cpu2,
	//						//		corev1.ResourceMemory:apim2mem2,
	//						//	},
	//						//},
	//
	//						ImagePullPolicy: "Always",
	//
	//						Ports: []corev1.ContainerPort{
	//							{
	//								ContainerPort: 8280,
	//								Protocol:      "TCP",
	//							},
	//							{
	//								ContainerPort: 8243,
	//								Protocol:      "TCP",
	//							},
	//							{
	//								ContainerPort: 9763,
	//								Protocol:      "TCP",
	//							},
	//							{
	//								ContainerPort: 9443,
	//								Protocol:      "TCP",
	//							},
	//
	//
	//						},
	//						Env: []corev1.EnvVar{
	//							// {
	//							// 	Name:  "HOST_NAME",
	//							// 	Value: "foo-am",
	//							// },
	//							{
	//								Name: "NODE_IP",
	//								ValueFrom: &corev1.EnvVarSource{
	//									FieldRef: &corev1.ObjectFieldSelector{
	//										FieldPath: "status.podIP",
	//									},
	//								},
	//							},
	//						},
	//						VolumeMounts: []corev1.VolumeMount{
	//							//{
	//							//	Name: "wso2am-pattern-1-am-volume-claim-synapse-configs",
	//							//	MountPath: "/home/wso2carbon/wso2am-3.0.0/repository/deployment/server/synapse-configs",
	//							//},
	//							//{
	//							//	Name: "wso2am-pattern-1-am-volume-claim-executionplans",
	//							//	MountPath: "/home/wso2carbon/wso2am-3.0.0/repository/deployment/server/executionplans",
	//							//},
	//							//{
	//							//	Name: "wso2am-pattern-1-am-2-conf",
	//							//	MountPath: "/home/wso2carbon/wso2-config-volume/repository/conf/deployment.toml",
	//							//	SubPath:"deployment.toml",
	//							//},
	//							//{
	//							//	Name: "mysql-jdbc-driver",
	//							//	MountPath: "/home/wso2carbon/wso2am-3.0.0/repository/components/lib",
	//							//},
	//							//{
	//							//	Name: "wso2am-pattern-1-am-conf-entrypoint",
	//							//	MountPath: "/home/wso2carbon/docker-entrypoint.sh",
	//							//	SubPath:"docker-entrypoint.sh",
	//							//},
	//						},
	//					},
	//				},
	//
	//				ServiceAccountName: "wso2am-pattern-1-svc-account",
	//				ImagePullSecrets:[]corev1.LocalObjectReference{
	//					{
	//						Name:"wso2am-pattern-1-creds",
	//					},
	//				},
	//
	//				Volumes: []corev1.Volume{
	//					//{
	//					//	Name: "wso2am-pattern-1-am-volume-claim-synapse-configs",
	//					//	VolumeSource: corev1.VolumeSource{
	//					//		PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
	//					//			ClaimName: "wso2am-pattern-1-am-volume-claim-synapse-configs",
	//					//		},
	//					//	},
	//					//},
	//					//{
	//					//	Name: "wso2am-pattern-1-am-volume-claim-executionplans",
	//					//	VolumeSource: corev1.VolumeSource{
	//					//		PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
	//					//			ClaimName: "wso2am-pattern-1-am-volume-claim-executionplans",
	//					//		},
	//					//	},
	//					//},
	//					//{
	//					//	Name: "wso2am-pattern-1-am-2-conf",
	//					//	VolumeSource: corev1.VolumeSource{
	//					//		ConfigMap: &corev1.ConfigMapVolumeSource{
	//					//			LocalObjectReference: corev1.LocalObjectReference{
	//					//				Name: "wso2am-pattern-1-am-2-conf",
	//					//			},
	//					//		},
	//					//	},
	//					//},
	//					//{
	//					//	Name: "mysql-jdbc-driver",
	//					//	VolumeSource: corev1.VolumeSource{
	//					//		ConfigMap: &corev1.ConfigMapVolumeSource{
	//					//			LocalObjectReference: corev1.LocalObjectReference{
	//					//				Name: "mysql-jdbc-driver-cm",
	//					//			},
	//					//		},
	//					//	},
	//					//},
	//					//{
	//					//	Name: "wso2am-pattern-1-am-conf-entrypoint",
	//					//	VolumeSource: corev1.VolumeSource{
	//					//		ConfigMap: &corev1.ConfigMapVolumeSource{
	//					//			LocalObjectReference: corev1.LocalObjectReference{
	//					//				Name: "wso2am-pattern-1-am-conf-entrypoint",
	//					//			},
	//					//			//DefaultMode:0407,
	//					//		},
	//					//	},
	//					//},
	//				},
	//			},
	//		},
	//	},
	//}
	labels := map[string]string{
		//"app":        "wso2am2",
		//"controller": apimanager.Name,
		"deployment": "wso2am-pattern-1-am",
		"node": "wso2am-pattern-1-am-2",
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			// Name:      apimanager.Spec.DeploymentName,
			Name:      "apim-2-deploy",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: apimanager.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					//InitContainers: []corev1.Container{
					//	{
					//		Name: "init-apim-analytics-db",
					//		Image: "busybox:1.31",
					//		Command: []string {
					//			//"sh", "-c", "echo -e \"Checking for the availability of MySQL Server deployment\" ; while ! nc -z \"while ! nc -z \"wso2am-mysql-db-service \"3306; do sleep 1; printf \"-\" ;done; echo -e \" >> MySQL Server has started \"; ",
					//			"sh",
					//			"-c",
					//			"echo -e \"Checking for the availability of MySQL Server deployment\"; while ! nc -z \"wso2am-mysql-db-service\" 3306; do sleep 1; printf \"-\"; done; echo -e \"  >> MySQL Server has started\";",
					//		},
					//
					//	},
					//	//{
					//	//	Name: "init-am-analytics-worker",
					//	//	Image: "busybox:1.31",
					//	//	Command: []string {
					//	//		"sh", "-c", "echo -e \"Checking for the availability of WSO2 API Manager Analytics Worker deployment\" ; while ! nc -z \"while ! nc -z \"wso2am-pattern-1-analytics-worker-service 7712; do sleep 1; printf \"-\" ;done; echo -e \" >> WSO2 API Manager Analytics Worker has started \"; ",
					//	//	},
					//	//
					//	//},
					//},
					Containers: []corev1.Container{
						{
							Name:  "wso2am2",
							Image: "wso2/wso2am:3.0.0",
							VolumeMounts: []corev1.VolumeMount{
								{
									Name: "wso2am-pattern-1-am-volume-claim-synapse-configs",
									MountPath: "/home/wso2carbon/wso2-artifact-volume/repository/deployment/server/synapse-configs",
								},
								{
									Name: "wso2am-pattern-1-am-volume-claim-executionplans",
									MountPath:"/home/wso2carbon/wso2-artifact-volume/repository/deployment/server/executionplans",
								},
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8280,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 8243,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 9443,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 9763,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 5672,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 9711,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 9611,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 7711,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 7611,
									Protocol:      "TCP",
								},
							},
						},
					},

					Volumes: []corev1.Volume{
						{
							Name: "wso2am-pattern-1-am-volume-claim-synapse-configs",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName:"pvc-synapse-configs",
								},
							},
						},
						{
							Name: "wso2am-pattern-1-am-volume-claim-executionplans",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "pvc-execution-plans",
								},
							},
						},
					},
				},
			},
		},
	}
}


// newService creates a new Service for a Apimanager resource.
// It expose the service with Nodeport type with minikube ip as the externel ip.
func apim2Service(apimanager *apimv1alpha1.Apimanager) *corev1.Service {
	labels := map[string]string{
		//"app":        "wso2am",
		//"controller": apimanager.Name,
		"deployment": "wso2am-pattern-1-am",
		"node": "wso2am-pattern-1-am-2",
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			// Name: apimanager.Spec.ServiceName,
			Name:      "apim-2-svc",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     "NodePort",
			// values are fetched from wso2-apim.yaml file
			// Type: apimanager.Spec.ServType,
			ExternalIPs: []string{"192.168.99.101"},
			// ExternalIPs: apimanager.Spec.ExternalIps,
			Ports: []corev1.ServicePort{
				{
					Name:       "servlet-https",
					Protocol:   corev1.ProtocolTCP,
					Port:       9443,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9443},
					NodePort:   32005,
				},
				{
					Name:       "binary",
					Protocol:   corev1.ProtocolTCP,
					Port:       9611,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9611},
					NodePort:   32006,
				},
				{
					Name:       "binary-secure",
					Protocol:   corev1.ProtocolTCP,
					Port:       9711,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9711},
					NodePort:   32007,
				},
				{
					Name:       "jms-tcp",
					Protocol:   corev1.ProtocolTCP,
					Port:       5672,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 5672},
					NodePort:   32008,
				},
			},
		},
	}
}


//// dashboardDeployment creates a new Deployment for a Apimanager instance 1 resource. It also sets
//// the appropriate OwnerReferences on the resource so handleObject can discover
//// the Apimanager resource that 'owns' it.
//func dashboardDeployment(apimanager *apimv1alpha1.Apimanager) *appsv1.Deployment {
//	labels := map[string]string{
//		"deployment": "wso2am-pattern-1-analytics-dashboard",
//	}
//	//defaultMode := int32(0407)
//	runasuser := int64(802)
//	dashcpu, _ :=resource.ParseQuantity("2000m")
//	dashmem, _ := resource.ParseQuantity("4Gi")
//
//	return &appsv1.Deployment{
//		ObjectMeta: metav1.ObjectMeta{
//			// Name: apimanager.Spec.DeploymentName,
//			Name:      "analytics-dash-deploy",
//			Namespace: apimanager.Namespace,
//			OwnerReferences: []metav1.OwnerReference{
//				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
//			},
//		},
//		Spec: appsv1.DeploymentSpec{
//			Replicas: apimanager.Spec.Replicas,
//			MinReadySeconds:30,
//			Strategy: appsv1.DeploymentStrategy{
//				Type: appsv1.DeploymentStrategyType(appsv1.RollingUpdateDaemonSetStrategyType),
//				RollingUpdate: &appsv1.RollingUpdateDeployment{
//					MaxSurge: &intstr.IntOrString{
//						Type:   intstr.Int,
//						IntVal: 1,
//					},
//					MaxUnavailable: &intstr.IntOrString{
//						Type:   intstr.Int,
//						IntVal: 0,
//					},
//				},
//			},
//
//			Selector: &metav1.LabelSelector{
//				MatchLabels: labels,
//			},
//			Template: corev1.PodTemplateSpec{
//				ObjectMeta: metav1.ObjectMeta{
//					Labels: labels,
//				},
//				Spec: corev1.PodSpec{
//
//					//InitContainers: []corev1.Container{
//					//	{
//					//		Name: "init-apim-analytics-db",
//					//		Image: "busybox:1.31",
//					//		Command: []string {
//					//			//"sh", "-c", "echo -e \"Checking for the availability of MySQL Server deployment\" ; while ! nc -z \"while ! nc -z \"wso2am-mysql-db-service \"3306; do sleep 1; printf \"-\" ;done; echo -e \" >> MySQL Server has started \"; ",
//					//			"sh",
//					//			"-c",
//					//			"echo -e \"Checking for the availability of MySQL Server deployment\"; while ! nc -z \"wso2am-mysql-db-service\" 3306; do sleep 1; printf \"-\"; done; echo -e \"  >> MySQL Server has started\";",
//					//		},
//					//
//					//	},
//					//	{
//					//		Name: "init-am-analytics-worker",
//					//		Image: "busybox:1.31",
//					//		Command: []string {
//					//			"sh", "-c", "echo -e \"Checking for the availability of WSO2 API Manager Analytics Worker deployment\" ; while ! nc -z \"while ! nc -z \"wso2am-pattern-1-analytics-worker-service 7712; do sleep 1; printf \"-\" ;done; echo -e \" >> WSO2 API Manager Analytics Worker has started \"; ",
//					//		},
//					//
//					//	},
//					//},
//					Containers: []corev1.Container{
//						{
//							Name:  "wso2am-pattern-1-analytics-dashboard",
//							Image: "wso2/wso2am-analytics-dashboard:3.0.0",
//							LivenessProbe: &corev1.Probe{
//								Handler: corev1.Handler{
//									Exec:&corev1.ExecAction{
//										Command:[]string{
//											"/bin/sh",
//											"-c",
//											"nc -z localhost 9643",
//										},
//									},
//								},
//								InitialDelaySeconds: 20,
//								PeriodSeconds:       10,
//							},
//							ReadinessProbe: &corev1.Probe{
//								Handler: corev1.Handler{
//									Exec:&corev1.ExecAction{
//										Command:[]string{
//											"/bin/sh",
//											"-c",
//											"nc -z localhost 9643",
//										},
//									},
//								},
//
//								InitialDelaySeconds: 20,
//								PeriodSeconds:       10,
//
//							},
//
//							Lifecycle: &corev1.Lifecycle{
//								PreStop:&corev1.Handler{
//									Exec:&corev1.ExecAction{
//										Command:[]string{
//											"sh",
//											"-c",
//											"${WSO2_SERVER_HOME}/bin/dashboard.sh stop",
//										},
//									},
//								},
//							},
//
//							Resources:corev1.ResourceRequirements{
//								Requests:corev1.ResourceList{
//									corev1.ResourceCPU:dashcpu,
//									corev1.ResourceMemory:dashmem,
//								},
//								Limits:corev1.ResourceList{
//									corev1.ResourceCPU:dashcpu,
//									corev1.ResourceMemory:dashmem,
//								},
//							},
//
//							ImagePullPolicy: "Always",
//
//							SecurityContext: &corev1.SecurityContext{
//								RunAsUser:&runasuser,
//							},
//
//							Ports: []corev1.ContainerPort{
//								{
//									ContainerPort: 9713,
//									Protocol:      "TCP",
//								},
//								{
//									ContainerPort: 9643,
//									Protocol:      "TCP",
//								},
//								{
//									ContainerPort: 9613,
//									Protocol:      "TCP",
//								},
//								{
//									ContainerPort: 7713,
//									Protocol:      "TCP",
//								},
//								{
//									ContainerPort: 9091,
//									Protocol:      "TCP",
//								},
//								{
//									ContainerPort: 7613,
//									Protocol:      "TCP",
//								},
//
//							},
//
//							VolumeMounts: []corev1.VolumeMount{
//								//{
//								//	Name: "wso2am-pattern-1-am-analytics-dashboard-conf",
//								//	MountPath: "/home/wso2carbon/wso2-config-volume/conf/dashboard/deployment.yaml",
//								//	SubPath:"deployment.yaml",
//								//},
//								//{
//								//	Name: "wso2am-pattern-1-am-analytics-dashboard-bin",
//								//	MountPath: "/home/wso2carbon/wso2-config-volume/wso2/dashboard/bin/carbon.sh",
//								//	SubPath:"carbon.sh",
//								//},
//								//{
//								//	Name: "wso2am-pattern-1-am-analytics-dashboard-plugins",
//								//	MountPath: "/home/wso2carbon/wso2-patch-volume/wso2/lib/plugins/org.wso2.analytics.apim.rest.api.proxy_3.0.0.jar",
//								//	SubPath:"org.wso2.analytics.apim.rest.api.proxy_3.0.0.jar",
//								//},
//								//{
//								//	Name: "wso2am-pattern-1-am-analytics-dashboard-conf-entrypoint",
//								//	MountPath: "/home/wso2carbon/docker-entrypoint.sh",
//								//	SubPath:"docker-entrypoint.sh",
//								//},
//							},
//						},
//					},
//
//					ServiceAccountName: "wso2am-pattern-1-svc-account",
//					ImagePullSecrets:[]corev1.LocalObjectReference{
//						{
//							Name:"wso2am-pattern-1-creds",
//						},
//					},
//
//					Volumes: []corev1.Volume{
//						//{
//						//	Name: "wso2am-pattern-1-am-analytics-dashboard-conf",
//						//	VolumeSource: corev1.VolumeSource{
//						//		ConfigMap: &corev1.ConfigMapVolumeSource{
//						//			LocalObjectReference: corev1.LocalObjectReference{
//						//				Name: "wso2am-pattern-1-am-analytics-dashboard-conf",
//						//			},
//						//			DefaultMode:&defaultMode,
//						//		},
//						//	},
//						//},
//						//{
//						//	Name: "wso2am-pattern-1-am-analytics-dashboard-bin",
//						//	VolumeSource: corev1.VolumeSource{
//						//		ConfigMap: &corev1.ConfigMapVolumeSource{
//						//			LocalObjectReference: corev1.LocalObjectReference{
//						//				Name: "wso2am-pattern-1-am-analytics-dashboard-bin",
//						//			},
//						//			DefaultMode:&defaultMode,
//						//		},
//						//	},
//						//},
//						//{
//						//	Name: "wso2am-pattern-1-am-analytics-dashboard-plugins",
//						//	VolumeSource: corev1.VolumeSource{
//						//		ConfigMap: &corev1.ConfigMapVolumeSource{
//						//			LocalObjectReference: corev1.LocalObjectReference{
//						//				Name: "wso2am-pattern-1-am-analytics-dashboard-plugins",
//						//			},
//						//			DefaultMode:&defaultMode,
//						//		},
//						//	},
//						//},
//						//{
//						//	Name: "wso2am-pattern-1-am-analytics-dashboard-conf-entrypoint",
//						//	VolumeSource: corev1.VolumeSource{
//						//		ConfigMap: &corev1.ConfigMapVolumeSource{
//						//			LocalObjectReference: corev1.LocalObjectReference{
//						//				Name: "wso2am-pattern-1-am-analytics-dashboard-conf-entrypoint",
//						//			},
//						//			DefaultMode:&defaultMode,
//						//		},
//						//	},
//						//},
//
//					},
//				},
//			},
//		},
//	}
//}
//
//
//// dashboardDeployment creates a new Deployment for a Apimanager instance 1 resource. It also sets
//// the appropriate OwnerReferences on the resource so handleObject can discover
//// the Apimanager resource that 'owns' it.
//func workerDeployment(apimanager *apimv1alpha1.Apimanager) *appsv1.Deployment {
//	workercpu,_ := resource.ParseQuantity("2000m")
//	workermem,_ := resource.ParseQuantity("4Gi")
//	labels := map[string]string{
//		"deployment": "wso2am-pattern-1-analytics-worker",
//	}
//	//defaultMode := int32(0407)
//	runasuser := int64(802)
//	return &appsv1.Deployment{
//		ObjectMeta: metav1.ObjectMeta{
//			// Name: apimanager.Spec.DeploymentName,
//			Name:      "analytics-worker-deploy",
//			Namespace: apimanager.Namespace,
//			OwnerReferences: []metav1.OwnerReference{
//				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
//			},
//		},
//		Spec: appsv1.DeploymentSpec{
//			Replicas: apimanager.Spec.Replicas,
//			MinReadySeconds:30,
//			Strategy: appsv1.DeploymentStrategy{
//				Type: appsv1.DeploymentStrategyType(appsv1.RollingUpdateDaemonSetStrategyType),
//				RollingUpdate: &appsv1.RollingUpdateDeployment{
//					MaxSurge: &intstr.IntOrString{
//						Type:   intstr.Int,
//						IntVal: 1,
//					},
//					MaxUnavailable: &intstr.IntOrString{
//						Type:   intstr.Int,
//						IntVal: 0,
//					},
//				},
//			},
//
//			Selector: &metav1.LabelSelector{
//				MatchLabels: labels,
//			},
//			Template: corev1.PodTemplateSpec{
//				ObjectMeta: metav1.ObjectMeta{
//					Labels: labels,
//				},
//				Spec: corev1.PodSpec{
//
//					//InitContainers: []corev1.Container{
//					//	{
//					//		Name: "init-apim-analytics-db",
//					//		Image: "busybox:1.31",
//					//		Command: []string {
//					//			//"sh", "-c", "echo -e \"Checking for the availability of MySQL Server deployment\" ; while ! nc -z \"while ! nc -z \"wso2am-mysql-db-service \"3306; do sleep 1; printf \"-\" ;done; echo -e \" >> MySQL Server has started \"; ",
//					//			"sh",
//					//			"-c",
//					//			"echo -e \"Checking for the availability of MySQL Server deployment\"; while ! nc -z \"wso2am-mysql-db-service\" 3306; do sleep 1; printf \"-\"; done; echo -e \"  >> MySQL Server has started\";",
//					//		},
//					//
//					//	},
//					//	{
//					//		Name: "init-am-analytics-worker",
//					//		Image: "busybox:1.31",
//					//		Command: []string {
//					//			"sh", "-c", "echo -e \"Checking for the availability of WSO2 API Manager Analytics Worker deployment\" ; while ! nc -z \"while ! nc -z \"wso2am-pattern-1-analytics-worker-service 7712; do sleep 1; printf \"-\" ;done; echo -e \" >> WSO2 API Manager Analytics Worker has started \"; ",
//					//		},
//					//
//					//	},
//					//},
//					Containers: []corev1.Container{
//						{
//							Name:  "wso2am-pattern-1-analytics-worker",
//							Image: "wso2/wso2am-analytics-worker:3.0.0",
//							LivenessProbe: &corev1.Probe{
//								Handler: corev1.Handler{
//									Exec:&corev1.ExecAction{
//										Command:[]string{
//											"/bin/sh",
//											"-c",
//											"nc -z localhost 9444",
//										},
//									},
//								},
//								InitialDelaySeconds: 20,
//								PeriodSeconds:       10,
//							},
//							ReadinessProbe: &corev1.Probe{
//								Handler: corev1.Handler{
//									Exec:&corev1.ExecAction{
//										Command:[]string{
//											"/bin/sh",
//											"-c",
//											"nc -z localhost 9444",
//										},
//									},
//								},
//
//								InitialDelaySeconds: 20,
//								PeriodSeconds:       10,
//
//							},
//
//							Lifecycle: &corev1.Lifecycle{
//								PreStop:&corev1.Handler{
//									Exec:&corev1.ExecAction{
//										Command:[]string{
//											"sh",
//											"-c",
//											"${WSO2_SERVER_HOME}/bin/worker.sh stop",
//										},
//									},
//								},
//							},
//
//							Resources:corev1.ResourceRequirements{
//								Requests:corev1.ResourceList{
//									corev1.ResourceCPU:workercpu,
//									corev1.ResourceMemory:workermem,
//								},
//								Limits:corev1.ResourceList{
//									corev1.ResourceCPU:workercpu,
//									corev1.ResourceMemory:workermem,
//								},
//							},
//
//							ImagePullPolicy: "Always",
//
//							SecurityContext: &corev1.SecurityContext{
//								RunAsUser:&runasuser,
//							},
//
//							Ports: []corev1.ContainerPort{
//								{
//									ContainerPort: 9764,
//									Protocol:      "TCP",
//								},
//								{
//									ContainerPort: 9444,
//									Protocol:      "TCP",
//								},
//								{
//									ContainerPort: 7612,
//									Protocol:      "TCP",
//								},
//								{
//									ContainerPort: 7712,
//									Protocol:      "TCP",
//								},
//								{
//									ContainerPort: 9091,
//									Protocol:      "TCP",
//								},
//								{
//									ContainerPort: 7071,
//									Protocol:      "TCP",
//								},
//								{
//									ContainerPort: 7444,
//									Protocol:      "TCP",
//								},
//							},
//
//							VolumeMounts: []corev1.VolumeMount{
//								//{
//								//	Name: "wso2am-pattern-1-am-analytics-worker-conf",
//								//	MountPath: "/home/wso2carbon/wso2-config-volume/conf/worker/deployment.yaml",
//								//	SubPath:"deployment.yaml",
//								//},
//
//							},
//						},
//					},
//
//					ServiceAccountName: "wso2am-pattern-1-svc-account",
//					ImagePullSecrets:[]corev1.LocalObjectReference{
//						{
//							Name:"wso2am-pattern-1-creds",
//						},
//					},
//
//					Volumes: []corev1.Volume{
//						//{
//						//	Name: "wso2am-pattern-1-am-analytics-worker-conf",
//						//	VolumeSource: corev1.VolumeSource{
//						//		ConfigMap: &corev1.ConfigMapVolumeSource{
//						//			LocalObjectReference: corev1.LocalObjectReference{
//						//				Name: "wso2am-pattern-1-am-analytics-worker-conf",
//						//			},
//						//		},
//						//	},
//						//},
//					},
//				},
//			},
//		},
//	}
//}



//
//// newService creates a new Service for a Apimanager resource.
//// It expose the service with Nodeport type with minikube ip as the externel ip.
//func workerService(apimanager *apimv1alpha1.Apimanager) *corev1.Service {
//	labels := map[string]string{
//		"deployment": "wso2am-pattern-1-analytics-worker",
//	}
//	return &corev1.Service{
//		ObjectMeta: metav1.ObjectMeta{
//			// Name: apimanager.Spec.ServiceName,
//			Name:      "analytics-worker-svc",
//			Namespace: apimanager.Namespace,
//			OwnerReferences: []metav1.OwnerReference{
//				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
//			},
//		},
//		Spec: corev1.ServiceSpec{
//			Selector: labels,
//			Type:     "NodePort",
//			// values are fetched from wso2-apim.yaml file
//			// Type: apimanager.Spec.ServType,
//			ExternalIPs: []string{"192.168.99.101"},
//			// ExternalIPs: apimanager.Spec.ExternalIps,
//			Ports: []corev1.ServicePort{
//				{
//					Name:       "thrift",
//					Protocol:   corev1.ProtocolTCP,
//					Port:       7612,
//					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7612},
//					NodePort:   32010,
//				},
//				{
//					Name:       "thrift-ssl",
//					Protocol:   corev1.ProtocolTCP,
//					Port:       7712,
//					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7712},
//					NodePort:   32011,
//				},
//				{
//					Name:       "rest-api-port-1",
//					Protocol:   corev1.ProtocolTCP,
//					Port:       9444,
//					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9444},
//					NodePort:   32012,
//				},
//				{
//					Name:       "rest-api-port-2",
//					Protocol:   corev1.ProtocolTCP,
//					Port:       9091,
//					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9091},
//					NodePort:   32013,
//				},
//				{
//					Name:       "rest-api-port-3",
//					Protocol:   corev1.ProtocolTCP,
//					Port:       7071,
//					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7071},
//					NodePort:   32014,
//				},
//				{
//					Name:       "rest-api-port-4",
//					Protocol:   corev1.ProtocolTCP,
//					Port:       7444,
//					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 7444},
//					NodePort:   32015,
//				},
//			},
//		},
//	}
//}
//
//
//// newService creates a new Service for a Apimanager resource.
//// It expose the service with Nodeport type with minikube ip as the externel ip.
//func dashboardService(apimanager *apimv1alpha1.Apimanager) *corev1.Service {
//	labels := map[string]string{
//		"deployment": "wso2am-pattern-1-analytics-dashboard",
//	}
//	return &corev1.Service{
//		ObjectMeta: metav1.ObjectMeta{
//			// Name: apimanager.Spec.ServiceName,
//			Name:      "analytics-dash-svc",
//			Namespace: apimanager.Namespace,
//			OwnerReferences: []metav1.OwnerReference{
//				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
//			},
//		},
//		Spec: corev1.ServiceSpec{
//			Selector: labels,
//			Type:     "NodePort",
//			// values are fetched from wso2-apim.yaml file
//			// Type: apimanager.Spec.ServType,
//			ExternalIPs: []string{"192.168.99.101"},
//			// ExternalIPs: apimanager.Spec.ExternalIps,
//			Ports: []corev1.ServicePort{
//				{
//					Name:       "analytics-dashboard",
//					Protocol:   corev1.ProtocolTCP,
//					Port:       9643,
//					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9643},
//					NodePort:   32009,
//				},
//			},
//		},
//	}
//}
//


func mysqlDeployment(apimanager *apimv1alpha1.Apimanager) *appsv1.Deployment {
	labels := map[string]string{
		"deployment": "wso2apim-with-analytics-mysql",
	}
	runasuser := int64(999)
	mysqlreplics := int32(1)

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mysql-deploy",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &mysqlreplics,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "wso2apim-with-analytics-mysql",
							Image: "mysql:5.7",
							ImagePullPolicy: "IfNotPresent",
							SecurityContext: &corev1.SecurityContext{
								RunAsUser: &runasuser,
							},
							Env: []corev1.EnvVar{
								{
									Name:  "MYSQL_ROOT_PASSWORD",
									Value: "root",
								},
								{
									Name:  "MYSQL_USER",
									Value: "wso2carbon",
								},
								{
									Name:  "MYSQL_PASSWORD",
									Value: "wso2carbon",
								},

							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 3306,
									Protocol:      "TCP",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name: "mysql-dbscripts",
									MountPath: "/docker-entrypoint-initdb.d",
								},

								{
									Name: "apim-rdbms-persistent-storage",
									MountPath: "/var/lib/mysql",
								},
							},
							Args: []string{
								"--max-connections",
								"10000",
							},

						},
					},

					Volumes: []corev1.Volume{
						{
							Name: "mysql-dbscripts",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "mysql-dbscripts",
									},
								},
							},
						},
						{
							Name: "apim-rdbms-persistent-storage",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName:"pvc-mysql",
								},
							},
						},
					},
					ServiceAccountName: "wso2am-pattern-1-svc-account",
				},
			},
		},
	}
}

func mysqlService(apimanager *apimv1alpha1.Apimanager) *corev1.Service {
	labels := map[string]string{
		"deployment": "wso2apim-with-analytics-mysql",
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wso2am-mysql-db-service",
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     "ClusterIP",
			Ports: []corev1.ServicePort{
				{
					Name:       "mysql-port",
					Protocol:   corev1.ProtocolTCP,
					Port:       3306,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 3306},
				},
			},
		},
	}
}
