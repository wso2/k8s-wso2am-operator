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

// NewController returns a new sample controller
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

	deploymentName := apimanager.Spec.DeploymentName
	if deploymentName == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		utilruntime.HandleError(fmt.Errorf("%s: deployment name must be specified", key))
		return nil
	}
	//deploymentName ="wso2-am-deploy-1"
	//deploymentName2 = "wso2-am-deploy-2"

	serviceName := apimanager.Spec.ServiceName
	if serviceName == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		utilruntime.HandleError(fmt.Errorf("%s: service name must be specified", key))
		return nil
	}

	// Get the deployment using hardcoded deployment name wso2-am-deploy-1
	deployment, err := c.deploymentsLister.Deployments(apimanager.Namespace).Get(deploymentName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		deployment, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Create(newDeployment(apimanager))
		if err != nil {
			return err
		}
	}

	//// Get the deployment using hardcoded deployment name wso2-am-deploy-2
	//deployment2, err := c.deploymentsLister.Deployments(apimanager.Namespace).Get(deploymentName)
	//// If the resource doesn't exist, we'll create it
	//if errors.IsNotFound(err) {
	//	deployment2, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Create(newDeployment2(apimanager))
	//	if err != nil {
	//		return err
	//	}
	//}

	// Get the service with the name specified in wso2-apim spec
	service, err := c.servicesLister.Services(apimanager.Namespace).Get(serviceName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		service, err = c.kubeclientset.CoreV1().Services(apimanager.Namespace).Create(newService(apimanager))
	}

	// If an error occurs during Get/Create, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	// If the Deployment is not controlled by this Apimanager resource, we should log
	// a warning to the event recorder and ret
	if !metav1.IsControlledBy(deployment, apimanager) {
		msg := fmt.Sprintf("Deployment1 %q already exists and is not managed by Apimanager", deployment.Name)
		c.recorder.Event(apimanager, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}

	////for instance 2
	//// If the Deployment2 is not controlled by this Apimanager resource, we should log
	//// a warning to the event recorder and ret
	//if !metav1.IsControlledBy(deployment2, apimanager) {
	//	msg := fmt.Sprintf("Deployment2 %q already exists and is not managed by Apimanager", deployment2.Name)
	//	c.recorder.Event(apimanager, corev1.EventTypeWarning, ErrResourceExists, msg)
	//	return fmt.Errorf(msg)
	//}

	// If the Service is not controlled by this Apimanager resource, we should log
	// a warning to the event recorder and ret
	if !metav1.IsControlledBy(service, apimanager) {
		msg := fmt.Sprintf("Service %q already exists and is not managed by Apimanager", service.Name)
		c.recorder.Event(apimanager, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}

	// If this number of the replicas on the Apimanager resource is specified, and the
	// number does not equal the current desired replicas on the Deployment, we
	// should update the Deployment resource.
	if apimanager.Spec.Replicas != nil && *apimanager.Spec.Replicas != *deployment.Spec.Replicas {
		klog.V(4).Infof("Apimanager %s replicas: %d, deployment replicas: %d", name, *apimanager.Spec.Replicas, *deployment.Spec.Replicas)
		deployment, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Update(newDeployment(apimanager))
	}

	////for instance 2 also
	//// If this number of the replicas on the Apimanager resource is specified, and the
	//// number does not equal the current desired replicas on the Deployment2, we
	//// should update the Deployment2 resource.
	//if apimanager.Spec.Replicas != nil && *apimanager.Spec.Replicas != *deployment2.Spec.Replicas {
	//	klog.V(4).Infof("Apimanager %s replicas: %d, deployment2 replicas: %d", name, *apimanager.Spec.Replicas, *deployment2.Spec.Replicas)
	//	deployment2, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Update(newDeployment2(apimanager))
	//}

	// If an error occurs during Update, we'll requeue the item so we can
	// attempt processing again later. THis could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	// Finally, we update the status block of the Apimanager resource to reflect the
	// current state of the world
	err = c.updateApimanagerStatus(apimanager, deployment)
	if err != nil {
		return err
	}

	////for instance 2 also
	//// Finally, we update the status block of the Apimanager resource to reflect the
	//// current state of the world
	//err = c.updateApimanagerStatus(apimanager, deployment2)
	//if err != nil {
	//	return err
	//}

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

// newDeployment creates a new Deployment for a Apimanager resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the Apimanager resource that 'owns' it.
func newDeployment(apimanager *apimv1alpha1.Apimanager) *appsv1.Deployment {
	labels := map[string]string{
		"app":        "wso2am",
		"controller": apimanager.Name,
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      apimanager.Spec.DeploymentName,
			//Name:      "wso2-am-deploy-1",
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
					Containers: []corev1.Container{
						{
							Name:  "wso2am1",
							Image: "wso2/wso2am:3.0.0",
							VolumeMounts: []corev1.VolumeMount{
								//{
								//	Name:      "wso2am-configmap-instance1",
								//	MountPath: "/home/wso2carbon/wso2-config-volume/repository/conf/deployment.toml",
								//	SubPath:   "deployment.toml",
								//},admin--Apple_v1.0.0.xml
								{
									Name: "pvclaimvol",
									//MountPath:"/home/wso2carbon/wso2am-3.0.0/repository/deployment/server/executionplans",
									//ReadOnly:true,
									MountPath: "/home/wso2carbon/wso2am-3.0.0/repository/deployment/server/synapse-configs",
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
							//Env: []corev1.EnvVar{
							//	{
							//		Name:  "HOST_NAME",
							//		Value: "foo-am",
							//	},
							//	{
							//		Name: "NODE_IP",
							//		ValueFrom: &corev1.EnvVarSource{
							//			FieldRef: &corev1.ObjectFieldSelector{
							//				FieldPath: "status.podIP",
							//			},
							//		},
							//	},
							//},
						},
					},
					//HostAliases: []corev1.HostAlias{
					//	{
					//		IP: "127.0.0.1",
					//		Hostnames: []string{
					//			"foo-am",
					//			"foo-gateway",
					//		},
					//	},
					//},

					Volumes: []corev1.Volume{
						//{
						//	Name: "wso2am-configmap-instance1",
						//	VolumeSource: corev1.VolumeSource{
						//		ConfigMap: &corev1.ConfigMapVolumeSource{
						//			LocalObjectReference: corev1.LocalObjectReference{
						//				Name: "newinstance1",
						//			},
						//		},
						//	},
						//},
						{
							Name: "pvclaimvol",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "wso2pvclaim2",
								},
							},
						},
					},
				},
			},
		},
	}
}
//
////this is for the API-M instance 2
//// newDeployment2 creates a new Deployment for a Apimanager resource. It also sets
//// the appropriate OwnerReferences on the resource so handleObject can discover
//// the Apimanager resource that 'owns' it.
//func newDeployment2(apimanager *apimv1alpha1.Apimanager) *appsv1.Deployment {
//	labels := map[string]string{
//		"app":        "wso2am",
//		"controller": apimanager.Name,
//	}
//	return &appsv1.Deployment{
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      apimanager.Spec.DeploymentName,
//			//Name:      "wso2-am-deploy-2",
//			Namespace: apimanager.Namespace,
//			OwnerReferences: []metav1.OwnerReference{
//				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
//			},
//		},
//		Spec: appsv1.DeploymentSpec{
//			Replicas: apimanager.Spec.Replicas,
//			Selector: &metav1.LabelSelector{
//				MatchLabels: labels,
//			},
//			Template: corev1.PodTemplateSpec{
//				ObjectMeta: metav1.ObjectMeta{
//					Labels: labels,
//				},
//				Spec: corev1.PodSpec{
//					Containers: []corev1.Container{
//						{
//							Name:  "wso2am1",
//							Image: "wso2/wso2am:3.0.0",
//							VolumeMounts: []corev1.VolumeMount{
//								//{
//								//	Name:      "wso2am-configmap-instance1",
//								//	MountPath: "/home/wso2carbon/wso2-config-volume/repository/conf/deployment.toml",
//								//	SubPath:    "deployment.toml",
//								//
//								//},
//								{
//									Name: "pvclaimvol",
//									//MountPath:"/home/wso2carbon/wso2am-3.0.0/repository/deployment/server/executionplans",
//									//ReadOnly:true,
//									MountPath: "/home/wso2carbon/wso2am-3.0.0/repository/deployment/server/synapse-configs",
//								},
//							},
//							Ports: []corev1.ContainerPort{
//								{
//									ContainerPort: 8280,
//									Protocol:      "TCP",
//								},
//								{
//									ContainerPort: 8243,
//									Protocol:      "TCP",
//								},
//								{
//									ContainerPort: 9443,
//									Protocol:      "TCP",
//								},
//								{
//									ContainerPort: 9763,
//									Protocol:      "TCP",
//								},
//								{
//									ContainerPort: 5672,
//									Protocol:      "TCP",
//								},
//								{
//									ContainerPort: 9711,
//									Protocol:      "TCP",
//								},
//								{
//									ContainerPort: 9611,
//									Protocol:      "TCP",
//								},
//								{
//									ContainerPort: 7711,
//									Protocol:      "TCP",
//								},
//								{
//									ContainerPort: 7611,
//									Protocol:      "TCP",
//								},
//							},
//						},
//					},
//
//					Volumes: []corev1.Volume{
//						//{
//						//	Name: "wso2am-configmap-instance1",
//						//	VolumeSource: corev1.VolumeSource{
//						//		ConfigMap: &corev1.ConfigMapVolumeSource{
//						//			LocalObjectReference: corev1.LocalObjectReference{
//						//				Name: "newinstance1",
//						//
//						//			},
//						//		},
//						//	},
//						//},
//						{
//							Name: "pvclaimvol",
//							VolumeSource: corev1.VolumeSource{
//								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
//									ClaimName: "wso2pvclaim2",
//								},
//							},
//						},
//					},
//				},
//			},
//		},
//	}
//}

// newService creates a new Service for a Apimanager resource.
// It expose the service with Nodeport type with minikube ip as the externel ip.
func newService(apimanager *apimv1alpha1.Apimanager) *corev1.Service {
	labels := map[string]string{
		"app":        "wso2am",
		"controller": apimanager.Name,
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      apimanager.Spec.ServiceName,
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			//Type:     "NodePort",    values are fetched from wso2-apim.yaml file
			Type: apimanager.Spec.ServType,
			//ExternalIPs: []string{"192.168.99.101"},
			ExternalIPs: apimanager.Spec.ExternalIps,
			Ports: []corev1.ServicePort{
				{
					Name:       "binary",
					Protocol:   corev1.ProtocolTCP,
					Port:       9443,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9443},
					NodePort:   32001,
				},
				{
					Name:       "port2",
					Protocol:   corev1.ProtocolTCP,
					Port:       8280,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8280},
					NodePort:   32002,
				},
				{
					Name:       "port3",
					Protocol:   corev1.ProtocolTCP,
					Port:       8243,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8243},
					NodePort:   32003,
				},
				{
					Name:       "port4",
					Protocol:   corev1.ProtocolTCP,
					Port:       9763,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9763},
					NodePort:   32004,
				},
			},
		},
	}
}
