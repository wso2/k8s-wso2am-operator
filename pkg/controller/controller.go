/*
 *
 *  * Copyright (c) 2019 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
 *  *
 *  * WSO2 Inc. licenses this file to you under the Apache License,
 *  * Version 2.0 (the "License"); you may not use this file except
 *  * in compliance with the License.
 *  * You may obtain a copy of the License at
 *  *
 *  * http:www.apache.org/licenses/LICENSE-2.0
 *  *
 *  * Unless required by applicable law or agreed to in writing,
 *  * software distributed under the License is distributed on an
 *  * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 *  * KIND, either express or implied. See the License for the
 *  * specific language governing permissions and limitations
 *  * under the License.
 *
 */

package controller

import (
	"fmt"
	"strconv"
	"time"

	apimv1alpha1 "github.com/wso2/k8s-wso2am-operator/pkg/apis/apim/v1alpha1"
	"github.com/wso2/k8s-wso2am-operator/pkg/controller/mysql"
	"github.com/wso2/k8s-wso2am-operator/pkg/controller/pattern1"
	"github.com/wso2/k8s-wso2am-operator/pkg/controller/patternX"
	clientset "github.com/wso2/k8s-wso2am-operator/pkg/generated/clientset/versioned"
	samplescheme "github.com/wso2/k8s-wso2am-operator/pkg/generated/clientset/versioned/scheme"
	informers "github.com/wso2/k8s-wso2am-operator/pkg/generated/informers/externalversions/apim/v1alpha1"
	listers "github.com/wso2/k8s-wso2am-operator/pkg/generated/listers/apim/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	coreinformers "k8s.io/client-go/informers/core/v1"
	externsionsv1beta1informers "k8s.io/client-go/informers/extensions/v1beta1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appslisters "k8s.io/client-go/listers/apps/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	corev1listers "k8s.io/client-go/listers/core/v1"
	extensionsv1beta1listers "k8s.io/client-go/listers/extensions/v1beta1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
)

const controllerAgentName = "wso2am-controller"

// Controller is the controller implementation for Apimanager resources
type Controller struct {
	kubeclientset                kubernetes.Interface // kubeclientset is a standard kubernetes clientset
	sampleclientset              clientset.Interface  // sampleclientset is a clientset for our own API group
	deploymentsLister            appslisters.DeploymentLister
	statefulSetsLister           appslisters.StatefulSetLister
	servicesLister               corelisters.ServiceLister
	ingressLister                extensionsv1beta1listers.IngressLister //for ingress lister
	deploymentsSynced            cache.InformerSynced
	statefulSetsSynced           cache.InformerSynced
	servicesSynced               cache.InformerSynced
	apimanagerslister            listers.APIManagerLister
	apimanagersSynced            cache.InformerSynced
	configMapLister              corev1listers.ConfigMapLister
	persistentVolumeClaimsLister corelisters.PersistentVolumeClaimLister
	recorder                     record.EventRecorder // recorder is an event recorder for recording Event resources to the Kubernetes API.
	workqueue                    workqueue.RateLimitingInterface
	// workqueue is a rate limited work queue. This is used to queue work to be processed instead of performing it as
	// soon as a change happens. This means we can ensure we only process a fixed amount of resources at a time, and
	// makes it easy to ensure we are never processing the same item simultaneously in two different workers.
}

// NewController returns a new wso2am controller
func NewController(
	kubeclientset kubernetes.Interface,
	sampleclientset clientset.Interface,
	deploymentInformer appsinformers.DeploymentInformer,
	statefulsetInformer appsinformers.StatefulSetInformer,
	ingressinformer externsionsv1beta1informers.IngressInformer,
	serviceInformer coreinformers.ServiceInformer,
	configmapInformer coreinformers.ConfigMapInformer,
	persistentVolumeClaimInformer coreinformers.PersistentVolumeClaimInformer,
	apimanagerInformer informers.APIManagerInformer) *Controller {

	// Create event broadcaster.
	// Add apim-controller types to the default Kubernetes Scheme so Events can be logged for apim-controller types.
	utilruntime.Must(samplescheme.AddToScheme(scheme.Scheme))
	klog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeclientset:                kubeclientset,
		sampleclientset:              sampleclientset,
		deploymentsLister:            deploymentInformer.Lister(),
		deploymentsSynced:            deploymentInformer.Informer().HasSynced,
		statefulSetsLister:           statefulsetInformer.Lister(),
		statefulSetsSynced:           statefulsetInformer.Informer().HasSynced,
		servicesLister:               serviceInformer.Lister(),
		ingressLister:                ingressinformer.Lister(),
		servicesSynced:               serviceInformer.Informer().HasSynced,
		configMapLister:              configmapInformer.Lister(),
		persistentVolumeClaimsLister: persistentVolumeClaimInformer.Lister(),
		apimanagerslister:            apimanagerInformer.Lister(),
		apimanagersSynced:            apimanagerInformer.Informer().HasSynced,
		workqueue:                    workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Apimanagers"),
		recorder:                     recorder,
	}

	klog.Info("Setting up event handlers")
	// Set up an event handler for when Apimanager resources change
	apimanagerInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueApimanager,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueApimanager(new)
		},
	})
	// Set up an event handler for when Deployment resources change. This handler will lookup the owner of the given
	// Deployment, and if it is owned by a Apimanager resource will enqueue that Apimanager resource for processing.
	// This way, we don't need to implement custom logic for handling Deployment resources. More info on this pattern:
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
	statefulsetInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.handleObject,
		UpdateFunc: func(old, new interface{}) {
			newDepl := new.(*appsv1.StatefulSet)
			oldDel := old.(*appsv1.StatefulSet)
			if newDepl.ResourceVersion == oldDel.ResourceVersion {
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

	persistentVolumeClaimInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.handleObject,
		DeleteFunc: controller.handleObject,
	})

	configmapInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.handleObject,
		UpdateFunc: func(old, new interface{}) {
			newConf := new.(*corev1.ConfigMap)
			oldConf := old.(*corev1.ConfigMap)
			newConf.ResourceVersion = ""
			oldConf.ResourceVersion = ""

			controller.handleObject(new)
		},
		DeleteFunc: controller.handleObject,
	})

	ingressinformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.handleObject,
		DeleteFunc: controller.handleObject,
	})

	return controller
}

// Run will set up the event handlers for types we are interested in, as well as syncing informer caches and starting
// workers. It will block until stopCh is closed, at which point it will shutdown the workqueue and wait for workers
// to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	klog.Info("Starting APIManager controller")

	// Wait for the caches to be synced before starting workers
	//if controller is down and then once its up, controller query api server for objects, so it has to wait for objects in cache to sync
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

// runWorker is a long-running function that will continually call the processNextWorkItem function in order to read
// and process a message on the workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and attempt to process it, by calling the syncHandler.
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
		klog.Infof("Current object in workqueue: %s", obj)
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
			klog.Warningf("Invalid item in workqueue: %s", key)
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

// syncHandler compares the actual state with the desired, and attempts to converge the two.
// It then updates the Status block of the Apimanager resource with the current status of the resource.
// c is the Controller object type pointer as a parameter
func (c *Controller) syncHandler(key string) error {
	// Split the key into a namespace & object name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the Apimanager resource with this namespace/name from the lister
	// Use a Lister to find the object in the API server
	apimanager, err := c.apimanagerslister.APIManagers(namespace).Get(name)
	if err != nil {
		// The Apimanager resource may no longer exist, in which case we stop processing.
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("apimanager '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}

	configMapName := "wso2am-operator-controller-config"
	configmap, err := c.configMapLister.ConfigMaps("wso2-system").Get(configMapName)

	// UseMysql - default to true
	useMysqlPod := true
	if apimanager.Spec.UseMysql != "" {
		useMysqlPod, err = strconv.ParseBool(apimanager.Spec.UseMysql)
		if err != nil {
			return err
		}
	}

	// Enablenalytics - default to true
	enableAnalytics := true
	if apimanager.Spec.EnableAnalytics != "" {
		enableAnalytics, err = strconv.ParseBool(apimanager.Spec.EnableAnalytics)
		if err != nil {
			return err
		}
	}

	if apimanager.Spec.Pattern == "Pattern-1" {

		apim1deploymentName := "wso2-am-1-" + apimanager.Name
		apim2deploymentName := "wso2-am-2-" + apimanager.Name
		apim1serviceName := "wso2-am-1-svc"
		apim2serviceName := "wso2-am-2-svc"
		apimcommonservice := "wso2-am-svc"
		mysqldeploymentName := "mysql-" + apimanager.Name
		mysqlserviceName := "mysql-svc"
		dashboardDeploymentName := "wso2-am-analytics-dashboard-" + apimanager.Name
		dashboardServiceName := "wso2-am-analytics-dashboard-svc"
		workerDeploymentName := "wso2-am-analytics-worker-statefulset"
		workerServiceName := "wso2-am-analytics-worker-svc"
		workerHlServiceName := "wso2-am-analytics-worker-headless-svc"

		if enableAnalytics {
			dashConfName := "wso2am-p1-analytics-dash-conf"
			dashConfWso2, err := c.configMapLister.ConfigMaps("wso2-system").Get(dashConfName)
			dashConfUserName := "wso2am-p1-analytics-dash-conf-" + apimanager.Name
			dashConfUser, err := c.configMapLister.ConfigMaps(apimanager.Namespace).Get(dashConfUserName)
			if errors.IsNotFound(err) {
				dashConfUser, err = c.kubeclientset.CoreV1().ConfigMaps(apimanager.Namespace).Create(pattern1.MakeConfigMap(apimanager, dashConfWso2))
				if err != nil {
					fmt.Println("Creating dashboard configmap in user specified ns", dashConfUser)
				}
			}

			workerConfName := "wso2am-p1-analytics-worker-conf"
			workerConfWso2, err := c.configMapLister.ConfigMaps("wso2-system").Get(workerConfName)
			workerConfUserName := "wso2am-p1-analytics-worker-conf-" + apimanager.Name
			workerConfUser, err := c.configMapLister.ConfigMaps(apimanager.Namespace).Get(workerConfUserName)
			if errors.IsNotFound(err) {
				workerConfUser, err = c.kubeclientset.CoreV1().ConfigMaps(apimanager.Namespace).Create(pattern1.MakeConfigMap(apimanager, workerConfWso2))
				if err != nil {
					fmt.Println("Creating worker configmap in user specified ns", workerConfUser)

				}
			}
		}

		mysqlDbConfName := "wso2am-p1-mysql-dbscripts"
		mysqlDbConfWso2, err := c.configMapLister.ConfigMaps("wso2-system").Get(mysqlDbConfName)
		mysqlDbConfUserName := "wso2am-p1-mysql-dbscripts-" + apimanager.Name
		mysqlDbConfUser, err := c.configMapLister.ConfigMaps(apimanager.Namespace).Get(mysqlDbConfUserName)
		if errors.IsNotFound(err) {
			mysqlDbConfUser, err = c.kubeclientset.CoreV1().ConfigMaps(apimanager.Namespace).Create(pattern1.MakeConfigMap(apimanager, mysqlDbConfWso2))
			if err != nil {
				fmt.Println("Creating mysql dbscripts configmap in user specified ns", mysqlDbConfUser)
			}
		}

		am1ConfName := "wso2am-p1-apim-1-conf"
		am1ConfWso2, err := c.configMapLister.ConfigMaps("wso2-system").Get(am1ConfName)
		am1ConfUserName := "wso2am-p1-apim-1-conf-" + apimanager.Name
		am1ConfUser, err := c.configMapLister.ConfigMaps(apimanager.Namespace).Get(am1ConfUserName)
		if errors.IsNotFound(err) {
			am1ConfUser, err = c.kubeclientset.CoreV1().ConfigMaps(apimanager.Namespace).Create(pattern1.MakeConfigMap(apimanager, am1ConfWso2))
			if err != nil {
				fmt.Println("Creating am1 configmap in user specified ns", am1ConfUser)

			}
		}

		am2ConfName := "wso2am-p1-apim-2-conf"
		am2ConfWso2, err := c.configMapLister.ConfigMaps("wso2-system").Get(am2ConfName)
		am2ConfUserName := "wso2am-p1-apim-2-conf-" + apimanager.Name
		am2ConfUser, err := c.configMapLister.ConfigMaps(apimanager.Namespace).Get(am2ConfUserName)
		if errors.IsNotFound(err) {
			am2ConfUser, err = c.kubeclientset.CoreV1().ConfigMaps(apimanager.Namespace).Create(pattern1.MakeConfigMap(apimanager, am2ConfWso2))
			if err != nil {
				fmt.Println("Creating am2 configmap in user specified ns", am2ConfUser)
			}
		}

		analyticsBinConfName := "wso2am-p1-analytics-bin"
		analyticsBinConfWso2, err := c.configMapLister.ConfigMaps("wso2-system").Get(analyticsBinConfName)
		analyticsBinConfUserName := "wso2am-p1-analytics-bin-" + apimanager.Name
		analyticsBinConfUser, err := c.configMapLister.ConfigMaps(apimanager.Namespace).Get(analyticsBinConfUserName)
		if errors.IsNotFound(err) {
			analyticsBinConfUser, err = c.kubeclientset.CoreV1().ConfigMaps(apimanager.Namespace).Create(pattern1.MakeConfigMap(apimanager, analyticsBinConfWso2))
			klog.Error("Dash Configs 1 Error: ", err)
			if err != nil {
				fmt.Println("Creating analytics bin configmap in user specified ns", analyticsBinConfUser)
			}
		}

		// Parse the object and look for it’s deployment
		// Use a Lister to find the deployment object referred to in the Apimanager resource
		// Get apim instance 1 deployment name using hardcoded value

		am1num := 0
		am2num := 0
		dashnum := 0
		worknum := 0

		totalProfiles := len(apimanager.Spec.Profiles)

		i := 0

		if totalProfiles > 0 {
			for i = 0; i < totalProfiles; i++ {
				if apimanager.Spec.Profiles[i].Name == "api-manager-1" {
					am1num = i
				}
				if apimanager.Spec.Profiles[i].Name == "api-manager-2" {
					am2num = i
				}
				if apimanager.Spec.Profiles[i].Name == "analytics-dashboard" {
					dashnum = i
				}
				if apimanager.Spec.Profiles[i].Name == "analytics-worker" {
					worknum = i
				}
			}
		}

		// Get mysql deployment name using hardcoded value
		mysqldeployment, err := c.deploymentsLister.Deployments(apimanager.Namespace).Get(mysqldeploymentName)

		if useMysqlPod {
			// If the resource doesn't exist, we'll create it
			if errors.IsNotFound(err) {
				//y:= pattern1.AssignMysqlConfigMapValues(apimanager,configmap)
				mysqldeployment, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Create(mysql.MysqlDeployment(apimanager))
				if err != nil {
					return err
				}
			}

			// Get mysql service name using hardcoded value
			mysqlservice, err := c.servicesLister.Services(apimanager.Namespace).Get(mysqlserviceName)

			// If the resource doesn't exist, we'll create it
			if errors.IsNotFound(err) {
				mysqlservice, err = c.kubeclientset.CoreV1().Services(apimanager.Namespace).Create(mysql.MysqlService(apimanager))
			} else {
				fmt.Println("Mysql Service is already available. [Service name] ,", mysqlservice)
			}

			for mysqldeployment.Status.ReadyReplicas == 0 {
				time.Sleep(5 * time.Second)
				mysqldeployment, err = c.deploymentsLister.Deployments(apimanager.Namespace).Get(mysqldeploymentName)
			}
		}

		// Get analytics dashboard deployment name using hardcoded value
		dashdeployment, err := c.deploymentsLister.Deployments(apimanager.Namespace).Get(dashboardDeploymentName)
		if enableAnalytics {
			// If the dash resource doesn't exist, we'll create it
			if errors.IsNotFound(err) {
				y := pattern1.AssignApimAnalyticsDashboardConfigMapValues(apimanager, configmap, dashnum)

				dashdeployment, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Create(pattern1.DashboardDeployment(apimanager, y, dashnum))
				if err != nil {
					return err
				}
			}

			// Get analytics dashboard service name using hardcoded value
			dashservice, err := c.servicesLister.Services(apimanager.Namespace).Get(dashboardServiceName)
			// If the resource doesn't exist, we'll create it
			if errors.IsNotFound(err) {
				dashservice, err = c.kubeclientset.CoreV1().Services(apimanager.Namespace).Create(pattern1.DashboardService(apimanager))
			} else {
				fmt.Println("Dash Service is already available. [Service name] ,", dashservice)
			}
		}

		// Get analytics worker deployment name using hardcoded value
		workerdeployment, err1 := c.statefulSetsLister.StatefulSets(apimanager.Namespace).Get(workerDeploymentName)
		if enableAnalytics {
			// If the worker resource doesn't exist, we'll create it
			if errors.IsNotFound(err1) {
				y := pattern1.AssignApimAnalyticsWorkerConfigMapValues(apimanager, configmap, worknum)

				// workerdeployment, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Create(pattern1.WorkerDeployment(apimanager, y, worknum))
				workerdeployment, err = c.kubeclientset.AppsV1().StatefulSets(apimanager.Namespace).Create(pattern1.WorkerDeployment(apimanager, y, worknum))
				if err != nil {
					return err
				}
			}
			// Get analytics worker service name using hardcoded value
			workerservice, err := c.servicesLister.Services(apimanager.Namespace).Get(workerServiceName)
			// If the resource doesn't exist, we'll create it
			if errors.IsNotFound(err) {
				workerservice, err = c.kubeclientset.CoreV1().Services(apimanager.Namespace).Create(pattern1.WorkerService(apimanager))
			} else {
				fmt.Println("Worker Service is already available. [Service name] ,", workerservice)
			}

			//Get worker-analytics headless service
			workerhlservice, err := c.servicesLister.Services(apimanager.Namespace).Get(workerHlServiceName)

			if errors.IsNotFound(err) {
				workerhlservice, err = c.kubeclientset.CoreV1().Services(apimanager.Namespace).Create(pattern1.WorkerHeadlessService(apimanager))
			} else {
				fmt.Println("Worker Headless Service is already available. [Service name] ,", workerhlservice)
			}
		}

		deployment, err := c.deploymentsLister.Deployments(apimanager.Namespace).Get(apim1deploymentName)
		// If the resource doesn't exist, we'll create it
		if errors.IsNotFound(err) {
			x := pattern1.AssignApimConfigMapValues(apimanager, configmap, am1num)

			deployment, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Create(pattern1.Apim1Deployment(apimanager, x, am1num))
			if err != nil {
				return err
			}
		}

		// Get apim instance 1 service name using hardcoded value
		service, err := c.servicesLister.Services(apimanager.Namespace).Get(apim1serviceName)
		// If the resource doesn't exist, we'll create it
		if errors.IsNotFound(err) {
			service, err = c.kubeclientset.CoreV1().Services(apimanager.Namespace).Create(pattern1.Apim1Service(apimanager))
		}

		// Get apim instance 2 deployment name using hardcoded value
		deployment2, err := c.deploymentsLister.Deployments(apimanager.Namespace).Get(apim2deploymentName)
		// If the resource doesn't exist, we'll create it
		if errors.IsNotFound(err) {
			z := pattern1.AssignApimConfigMapValues(apimanager, configmap, am2num)

			deployment2, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Create(pattern1.Apim2Deployment(apimanager, z, am2num))
			if err != nil {
				return err
			}
		}

		// Get apim instance 2 service name using hardcoded value
		service2, err := c.servicesLister.Services(apimanager.Namespace).Get(apim2serviceName)
		// If the resource doesn't exist, we'll create it
		if errors.IsNotFound(err) {
			service2, err = c.kubeclientset.CoreV1().Services(apimanager.Namespace).Create(pattern1.Apim2Service(apimanager))
		}

		// Get apim common service name using hardcoded value
		commonservice, err := c.servicesLister.Services(apimanager.Namespace).Get(apimcommonservice)
		// If the resource doesn't exist, we'll create it
		if errors.IsNotFound(err) {
			commonservice, err = c.kubeclientset.CoreV1().Services(apimanager.Namespace).Create(pattern1.ApimCommonService(apimanager))
		}

		// If an error occurs during Get/Create, we'll requeue the item so we can
		// attempt processing again later. This could have been caused by a
		// temporary network failure, or any other transient reason.
		if err != nil {
			return err
		}

		if apimanager.Spec.Expose == "Ingress" {
			// Get apim instance 1 service name using hardcoded value
			apimingressname := "wso2-am-p1-ingress"
			amingress, err := c.ingressLister.Ingresses(apimanager.Namespace).Get(apimingressname)
			// If the resource doesn't exist, we'll create it
			if errors.IsNotFound(err) {
				amingress, err = c.kubeclientset.ExtensionsV1beta1().Ingresses(apimanager.Namespace).Create(pattern1.ApimIngress(apimanager))
				if err != nil {
					return err
				}
			}
			// Get apim instance 1 service name using hardcoded value
			gatewayingressname := "wso2-am-gateway-p1-ingress"
			gatewayingress, err := c.ingressLister.Ingresses(apimanager.Namespace).Get(gatewayingressname)
			// If the resource doesn't exist, we'll create it
			if errors.IsNotFound(err) {
				gatewayingress, err = c.kubeclientset.ExtensionsV1beta1().Ingresses(apimanager.Namespace).Create(pattern1.GatewayIngress(apimanager))
				if err != nil {
					return err
				}
			}

			// Get apim instance 1 service name using hardcoded value
			dashingressname := "wso2-am-analytics-dashboard-p1-ingress"
			dashingress, err1 := c.ingressLister.Ingresses(apimanager.Namespace).Get(dashingressname)

			if enableAnalytics {
				// If the resource doesn't exist, we'll create it
				if errors.IsNotFound(err1) {
					dashingress, err = c.kubeclientset.ExtensionsV1beta1().Ingresses(apimanager.Namespace).Create(pattern1.DashboardIngress(apimanager))
					if err != nil {
						return err
					}
				}
			}

			// If the apim ingress is not controlled by this Apimanager resource, we should log a warning to the event recorder and return
			if !metav1.IsControlledBy(amingress, apimanager) {
				msg := fmt.Sprintf("am ingress %q already exists and is not managed by APIManager", amingress.Name)
				c.recorder.Event(apimanager, corev1.EventTypeWarning, "ErrResourceExists", msg)
				return fmt.Errorf(msg)
			}
			// If the apim ingress is not controlled by this Apimanager resource, we should log a warning to the event recorder and return
			if !metav1.IsControlledBy(gatewayingress, apimanager) {
				msg := fmt.Sprintf("gateway ingress %q already exists and is not managed by APIManager", gatewayingress.Name)
				c.recorder.Event(apimanager, corev1.EventTypeWarning, "ErrResourceExists", msg)
				return fmt.Errorf(msg)
			}

			if enableAnalytics {
				// If the apim ingress is not controlled by this Apimanager resource, we should log a warning to the event recorder and return
				if !metav1.IsControlledBy(dashingress, apimanager) {
					msg := fmt.Sprintf("dashboard ingress %q already exists and is not managed by APIManager", dashingress.Name)
					c.recorder.Event(apimanager, corev1.EventTypeWarning, "ErrResourceExists", msg)
					return fmt.Errorf(msg)
				}
			}
		}

		/////////////checking whether resources are controlled by apimanager with same owner reference

		// If the apim instance 1 Deployment is not controlled by this Apimanager resource, we should log a warning to the event recorder and return
		if !metav1.IsControlledBy(deployment, apimanager) {
			msg := fmt.Sprintf("Deployment1 %q already exists and is not managed by APIManager", deployment.Name)
			c.recorder.Event(apimanager, corev1.EventTypeWarning, "ErrResourceExists", msg)
			return fmt.Errorf(msg)
		}

		// If the apim instance 2 Deployment is not controlled by this Apimanager resource, we should log a warning to the event recorder and return
		if !metav1.IsControlledBy(deployment2, apimanager) {
			msg := fmt.Sprintf("Deployment2 %q already exists and is not managed by APIManager", deployment2.Name)
			c.recorder.Event(apimanager, corev1.EventTypeWarning, "ErrResourceExists", msg)
			return fmt.Errorf(msg)
		}

		if enableAnalytics {
			// If the analytics dashboard Deployment is not controlled by this Apimanager resource, we should log a warning to the event recorder and return
			if !metav1.IsControlledBy(dashdeployment, apimanager) {
				msg := fmt.Sprintf("Analytics Dashboard Deployment %q already exists and is not managed by APIManager", dashdeployment.Name)
				c.recorder.Event(apimanager, corev1.EventTypeWarning, "ErrResourceExists", msg)
				return fmt.Errorf(msg)
			}

			// If the analytics worker Deployment is not controlled by this Apimanager resource, we should log a warning to the event recorder and return
			if !metav1.IsControlledBy(workerdeployment, apimanager) {
				msg := fmt.Sprintf("Analytics Worker Deployment %q already exists and is not managed by APIManager", workerdeployment.Name)
				c.recorder.Event(apimanager, corev1.EventTypeWarning, "ErrResourceExists", msg)
				return fmt.Errorf(msg)
			}
		}

		if useMysqlPod {
			//// If the mysql Deployment is not controlled by this Apimanager resource, we should log a warning to the event recorder and return
			if !metav1.IsControlledBy(mysqldeployment, apimanager) {
				msg := fmt.Sprintf("mysql deployment %q already exists and is not managed by APIManager", mysqldeployment.Name)
				c.recorder.Event(apimanager, corev1.EventTypeWarning, "ErrResourceExists", msg)
				return fmt.Errorf(msg)
			}
		}

		// If the apim instance 1 Service is not controlled by this Apimanager resource, we should log a warning to the event recorder and return
		if !metav1.IsControlledBy(service, apimanager) {
			msg := fmt.Sprintf("service1 %q already exists and is not managed by APIManager", service.Name)
			c.recorder.Event(apimanager, corev1.EventTypeWarning, "ErrResourceExists", msg)
			return fmt.Errorf(msg)
		}

		// If the apim instance 2 Service is not controlled by this Apimanager resource, we should log a warning to the event recorder and return
		if !metav1.IsControlledBy(service2, apimanager) {
			msg := fmt.Sprintf("service2 %q already exists and is not managed by APIManager", service2.Name)
			c.recorder.Event(apimanager, corev1.EventTypeWarning, "ErrResourceExists", msg)
			return fmt.Errorf(msg)
		}

		if enableAnalytics {

			dashservice, _ := c.servicesLister.Services(apimanager.Namespace).Get(dashboardServiceName)
			// If the analytics dashboard Service is not controlled by this Apimanager resource, we should log a warning to the event recorder and return
			if !metav1.IsControlledBy(dashservice, apimanager) {
				msg := fmt.Sprintf("dashboard Service %q already exists and is not managed by APIManager", dashservice.Name)
				c.recorder.Event(apimanager, corev1.EventTypeWarning, "ErrResourceExists", msg)
				return fmt.Errorf(msg)
			}

			workerservice, _ := c.servicesLister.Services(apimanager.Namespace).Get(workerServiceName)
			// If the analytics worker Service is not controlled by this Apimanager resource, we should log a warning to the event recorder and return
			if !metav1.IsControlledBy(workerservice, apimanager) {
				msg := fmt.Sprintf("worker Service %q already exists and is not managed by APIManager", workerservice.Name)
				c.recorder.Event(apimanager, corev1.EventTypeWarning, "ErrResourceExists", msg)
				return fmt.Errorf(msg)
			}

			workerhlservice, _ := c.servicesLister.Services(apimanager.Namespace).Get(workerHlServiceName)
			// If the analytics worker Headless Service is not controlled by this Apimanager resource, we should log a warning to the event recorder and return
			if !metav1.IsControlledBy(workerhlservice, apimanager) {
				msg := fmt.Sprintf("worker headless Service %q already exists and is not managed by APIManager", workerhlservice.Name)
				c.recorder.Event(apimanager, corev1.EventTypeWarning, "ErrResourceExists", msg)
				return fmt.Errorf(msg)
			}
		}

		// If the analytics worker Service is not controlled by this Apimanager resource, we should log a warning to the event recorder and return
		if !metav1.IsControlledBy(commonservice, apimanager) {
			msg := fmt.Sprintf("common Service %q already exists and is not managed by APIManager", commonservice.Name)
			c.recorder.Event(apimanager, corev1.EventTypeWarning, "ErrResourceExists", msg)
			return fmt.Errorf(msg)
		}

		if useMysqlPod {
			// If the mysql Service is not controlled by this Apimanager resource, we should log a warning to the event recorder and return

			mysqlservice, _ := c.servicesLister.Services(apimanager.Namespace).Get(mysqlserviceName)
			if !metav1.IsControlledBy(mysqlservice, apimanager) {
				msg := fmt.Sprintf("mysql service %q already exists and is not managed by APIManager", mysqlservice.Name)
				c.recorder.Event(apimanager, corev1.EventTypeWarning, "ErrResourceExists", msg)
				return fmt.Errorf(msg)
			}
		}

		///////////check replicas are same as defined for deployments

		// If the Apimanager resource has changed update the deployment
		// If this number of the replicas on the Apimanager resource is specified, and the number does not equal the
		// current desired replicas on the Deployment, we should update the Deployment resource.
		if apimanager.Spec.Replicas != nil && *apimanager.Spec.Replicas != *deployment.Spec.Replicas {
			x := pattern1.AssignApimConfigMapValues(apimanager, configmap, am1num)
			klog.V(4).Infof("APIManager %s replicas: %d, deployment replicas: %d", name, *apimanager.Spec.Replicas, *deployment.Spec.Replicas)
			deployment, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Update(pattern1.Apim1Deployment(apimanager, x, am1num))
		}

		//for apim instance 2 also
		if apimanager.Spec.Replicas != nil && *apimanager.Spec.Replicas != *deployment2.Spec.Replicas {
			z := pattern1.AssignApimConfigMapValues(apimanager, configmap, am2num)
			klog.V(4).Infof("APIManager %s replicas: %d, deployment2 replicas: %d", name, *apimanager.Spec.Replicas, *deployment2.Spec.Replicas)
			deployment2, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Update(pattern1.Apim2Deployment(apimanager, z, am2num))
		}

		if enableAnalytics {
			//for analytics dashboard deployment
			if apimanager.Spec.Replicas != nil && *apimanager.Spec.Replicas != *dashdeployment.Spec.Replicas {
				y := pattern1.AssignApimAnalyticsDashboardConfigMapValues(apimanager, configmap, dashnum)
				klog.V(4).Infof("APIManager %s replicas: %d, deployment2 replicas: %d", name, *apimanager.Spec.Replicas, *dashdeployment.Spec.Replicas)
				dashdeployment, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Update(pattern1.DashboardDeployment(apimanager, y, dashnum))
			}

			//for analytics worker deployment
			if apimanager.Spec.Replicas != nil && *apimanager.Spec.Replicas != *workerdeployment.Spec.Replicas {
				y := pattern1.AssignApimAnalyticsWorkerConfigMapValues(apimanager, configmap, worknum)
				klog.V(4).Infof("APIManager %s replicas: %d, deployment2 replicas: %d", name, *apimanager.Spec.Replicas, *workerdeployment.Spec.Replicas)
				workerdeployment, err = c.kubeclientset.AppsV1().StatefulSets(apimanager.Namespace).Update(pattern1.WorkerDeployment(apimanager, y, worknum))
			}
		}

		if useMysqlPod {
			//for instance mysql deployment
			if apimanager.Spec.Replicas != nil && *apimanager.Spec.Replicas != *mysqldeployment.Spec.Replicas {
				//y:= pattern1.AssignMysqlConfigMapValues(apimanager,configmap)
				klog.V(4).Infof("APIManager %s replicas: %d, deployment2 replicas: %d", name, *apimanager.Spec.Replicas, *mysqldeployment.Spec.Replicas)
				mysqldeployment, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Update(mysql.MysqlDeployment(apimanager))
			}
		}

		// If an error occurs during Update, we'll requeue the item so we can attempt processing again later.
		// This could have been caused by a temporary network failure, or any other transient reason.
		if err != nil {
			return err
		}

		//////////finally update the deployment resources after done checking

		// Finally, we update the status block of the Apimanager resource to reflect the current state of the world
		err = c.updateApimanagerStatus(apimanager, deployment)
		if err != nil {
			return err
		}

		//for instance 2 also
		err = c.updateApimanagerStatus(apimanager, deployment2)
		if err != nil {
			return err
		}

		if enableAnalytics {
			//for analytics dashboard deployment
			err = c.updateApimanagerStatus(apimanager, dashdeployment)
			if err != nil {
				return err
			}

			//for analytics worker deployment
			err = c.updateApiMangerStatusForStatefulSet(apimanager, workerdeployment)
			if err != nil {
				return err
			}
		}

		if useMysqlPod {
			//for mysql deployment
			err = c.updateApimanagerStatus(apimanager, mysqldeployment)
			if err != nil {
				return err
			}
		}

		c.recorder.Event(apimanager, corev1.EventTypeNormal, "synced", "APIManager synced successfully")
		return nil

	}

	if apimanager.Spec.Pattern == "Pattern-X" {

		configMapName := "wso2am-operator-controller-config"
		configmap, err := c.configMapLister.ConfigMaps("wso2-system").Get(configMapName)

		if errors.IsNotFound(err) {
			fmt.Println("Configmap not found!")
		}

		var apimVolDefined bool = false
		var dashVolDefined bool = false
		var workerVolDefined bool = false

		for _, r := range apimanager.Spec.Profiles {

			if r.Type == "api-manager" {

				apim1deploymentName := r.Name

				deployment, err := c.deploymentsLister.Deployments(apimanager.Namespace).Get(apim1deploymentName)
				// If the resource doesn't exist, we'll create it
				x := patternX.AssignApimXConfigMapValues(apimanager, configmap, r)
				if errors.IsNotFound(err) {
					deployment, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Create(patternX.ApimXDeployment(apimanager, &r, x))
					if err != nil {
						return err
					}
				}
				// Get apim instance 1 service name using hardcoded value
				apimXserviceName := r.Service.Name
				service, err := c.servicesLister.Services(apimanager.Namespace).Get(apimXserviceName)
				// If the resource doesn't exist, we'll create it
				if errors.IsNotFound(err) {
					service, err = c.kubeclientset.CoreV1().Services(apimanager.Namespace).Create(patternX.ApimXService(apimanager, &r))
				}

				if err != nil {
					return err
				}

				/////////////checking whether resources are controlled by apimanager with same owner reference

				// If the apim instance 1 Deployment is not controlled by this Apimanager resource, we should log a warning to the event recorder and return
				if !metav1.IsControlledBy(deployment, apimanager) {
					msg := fmt.Sprintf("Deployment1 %q already exists and is not managed by APIManager", deployment.Name)
					c.recorder.Event(apimanager, corev1.EventTypeWarning, "ErrResourceExists", msg)
					return fmt.Errorf(msg)
				}
				// If the apim instance 1 Service is not controlled by this Apimanager resource, we should log a warning to the event recorder and return
				if !metav1.IsControlledBy(service, apimanager) {
					msg := fmt.Sprintf("apimananger service %q already exists and is not managed by APIManager", service.Name)
					c.recorder.Event(apimanager, corev1.EventTypeWarning, "ErrResourceExists", msg)
					return fmt.Errorf(msg)
				}

				if apimanager.Spec.Replicas != nil && *apimanager.Spec.Replicas != *deployment.Spec.Replicas {
					x := patternX.AssignApimXConfigMapValues(apimanager, configmap, r)
					klog.V(4).Infof("APIManager %s replicas: %d, deployment replicas: %d", name, *apimanager.Spec.Replicas, *deployment.Spec.Replicas)
					deployment, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Update(patternX.ApimXDeployment(apimanager, &r, x))
				}

				// If an error occurs during Update, we'll requeue the item so we can attempt processing again later.
				// This could have been caused by a temporary network failure, or any other transient reason.
				if err != nil {
					return err
				}

				// Finally, we update the status block of the Apimanager resource to reflect the current state of the world
				err = c.updateApimanagerStatus(apimanager, deployment)
				if err != nil {
					return err
				}

			}

			if r.Type == "analytics-dashboard" {

				dashdeploymentName := r.Name

				deployment, err := c.deploymentsLister.Deployments(apimanager.Namespace).Get(dashdeploymentName)
				// If the resource doesn't exist, we'll create it
				x := patternX.AssignApimAnalyticsConfigMapValues(apimanager, configmap, r)

				if errors.IsNotFound(err) {
					deployment, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Create(patternX.DashboardXDeployment(apimanager, &r, x))
					if err != nil {
						return err
					}
				}
				// Get apim instance 1 service name using hardcoded value
				apimXserviceName := r.Service.Name
				service, err := c.servicesLister.Services(apimanager.Namespace).Get(apimXserviceName)
				// If the resource doesn't exist, we'll create it
				if errors.IsNotFound(err) {
					service, err = c.kubeclientset.CoreV1().Services(apimanager.Namespace).Create(patternX.DashboardXService(apimanager, &r))
				}

				if err != nil {
					return err
				}

				/////////////checking whether resources are controlled by apimanager with same owner reference

				// If the apim instance 1 Deployment is not controlled by this Apimanager resource, we should log a warning to the event recorder and return
				if !metav1.IsControlledBy(deployment, apimanager) {
					msg := fmt.Sprintf("dashboard deployment %q already exists and is not managed by APIManager", deployment.Name)
					c.recorder.Event(apimanager, corev1.EventTypeWarning, "ErrResourceExists", msg)
					return fmt.Errorf(msg)
				}
				// If the apim instance 1 Service is not controlled by this Apimanager resource, we should log a warning to the event recorder and return
				if !metav1.IsControlledBy(service, apimanager) {
					msg := fmt.Sprintf("dashboard service %q already exists and is not managed by APIManager", service.Name)
					c.recorder.Event(apimanager, corev1.EventTypeWarning, "ErrResourceExists", msg)
					return fmt.Errorf(msg)
				}

				if apimanager.Spec.Replicas != nil && *apimanager.Spec.Replicas != *deployment.Spec.Replicas {
					x := patternX.AssignApimAnalyticsConfigMapValues(apimanager, configmap, r)
					klog.V(4).Infof("APIManager %s replicas: %d, deployment replicas: %d", name, *apimanager.Spec.Replicas, *deployment.Spec.Replicas)
					deployment, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Update(patternX.ApimXDeployment(apimanager, &r, x))
				}

				// If an error occurs during Update, we'll requeue the item so we can attempt processing again later.
				// This could have been caused by a temporary network failure, or any other transient reason.
				if err != nil {
					return err
				}

				// Finally, we update the status block of the Apimanager resource to reflect the current state of the world
				err = c.updateApimanagerStatus(apimanager, deployment)
				if err != nil {
					return err
				}

			}

			if r.Type == "analytics-worker" {
				x := patternX.AssignApimAnalyticsConfigMapValues(apimanager, configmap, r)

				workerdeploymentName := r.Name

				deployment, err := c.deploymentsLister.Deployments(apimanager.Namespace).Get(workerdeploymentName)
				// If the resource doesn't exist, we'll create it

				if errors.IsNotFound(err) {
					deployment, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Create(patternX.WorkerXDeployment(apimanager, &r, x))
					if err != nil {
						return err
					}
				}
				// Get apim instance 1 service name using hardcoded value
				apimXserviceName := r.Service.Name
				service, err := c.servicesLister.Services(apimanager.Namespace).Get(apimXserviceName)
				// If the resource doesn't exist, we'll create it
				if errors.IsNotFound(err) {
					service, err = c.kubeclientset.CoreV1().Services(apimanager.Namespace).Create(patternX.WorkerXService(apimanager, &r))
				}

				if err != nil {
					return err
				}

				/////////////checking whether resources are controlled by apimanager with same owner reference

				// If the apim instance 1 Deployment is not controlled by this Apimanager resource, we should log a warning to the event recorder and return
				if !metav1.IsControlledBy(deployment, apimanager) {
					msg := fmt.Sprintf("worker deployment %q already exists and is not managed by APIManager", deployment.Name)
					c.recorder.Event(apimanager, corev1.EventTypeWarning, "ErrResourceExists", msg)
					return fmt.Errorf(msg)
				}
				// If the apim instance 1 Service is not controlled by this Apimanager resource, we should log a warning to the event recorder and return
				if !metav1.IsControlledBy(service, apimanager) {
					msg := fmt.Sprintf("worker service %q already exists and is not managed by APIManager", service.Name)
					c.recorder.Event(apimanager, corev1.EventTypeWarning, "ErrResourceExists", msg)
					return fmt.Errorf(msg)
				}

				if apimanager.Spec.Replicas != nil && *apimanager.Spec.Replicas != *deployment.Spec.Replicas {

					x := patternX.AssignApimAnalyticsConfigMapValues(apimanager, configmap, r)

					klog.V(4).Infof("APIManager %s replicas: %d, deployment replicas: %d", name, *apimanager.Spec.Replicas, *deployment.Spec.Replicas)
					deployment, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Update(patternX.WorkerXDeployment(apimanager, &r, x))
				}

				// If an error occurs during Update, we'll requeue the item so we can attempt processing again later.
				// This could have been caused by a temporary network failure, or any other transient reason.
				if err != nil {
					return err
				}

				// Finally, we update the status block of the Apimanager resource to reflect the current state of the world
				err = c.updateApimanagerStatus(apimanager, deployment)
				if err != nil {
					return err
				}

			}
			//else {
			//	fmt.Println("sorry NO matching type found, so no deployments & services are made")
			//	}

		}

		if apimVolDefined && dashVolDefined && workerVolDefined {

			mysqldeploymentName := "mysql-" + apimanager.Name
			mysqlserviceName := "mysql-svc"

			mysqlDbConfName := "wso2am-p1-mysql-dbscripts"
			mysqlDbConfWso2, err := c.configMapLister.ConfigMaps("wso2-system").Get(mysqlDbConfName)
			mysqlDbConfUserName := "wso2am-p1-mysql-dbscripts-" + apimanager.Name
			mysqlDbConfUser, err := c.configMapLister.ConfigMaps(apimanager.Namespace).Get(mysqlDbConfUserName)
			if errors.IsNotFound(err) {
				mysqlDbConfUser, err = c.kubeclientset.CoreV1().ConfigMaps(apimanager.Namespace).Create(pattern1.MakeConfigMap(apimanager, mysqlDbConfWso2))
				if err != nil {
					fmt.Println("Creating mysql dbscripts configmap in user specified ns", mysqlDbConfUser)
				}
			}

			if useMysqlPod {

				// Get mysql deployment name using hardcoded value
				mysqldeployment, err := c.deploymentsLister.Deployments(apimanager.Namespace).Get(mysqldeploymentName)
				// If the resource doesn't exist, we'll create it
				if errors.IsNotFound(err) {
					//y:= pattern1.AssignMysqlConfigMapValues(apimanager,configmap)
					mysqldeployment, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Create(mysql.MysqlDeployment(apimanager))
					if err != nil {
						return err
					}
				}
				//
				// Get mysql service name using hardcoded value
				mysqlservice, err := c.servicesLister.Services(apimanager.Namespace).Get(mysqlserviceName)
				// If the resource doesn't exist, we'll create it
				if errors.IsNotFound(err) {
					mysqlservice, err = c.kubeclientset.CoreV1().Services(apimanager.Namespace).Create(mysql.MysqlService(apimanager))
					if err != nil {
						return err
					}
				}

				//// If the mysql Deployment is not controlled by this Apimanager resource, we should log a warning to the event recorder and return
				if !metav1.IsControlledBy(mysqldeployment, apimanager) {
					msg := fmt.Sprintf("mysql deployment %q already exists and is not managed by APIManager", mysqldeployment.Name)
					c.recorder.Event(apimanager, corev1.EventTypeWarning, "ErrResourceExists", msg)
					return fmt.Errorf(msg)
				}

				// If the mysql Service is not controlled by this Apimanager resource, we should log a warning to the event recorder and return
				if !metav1.IsControlledBy(mysqlservice, apimanager) {
					msg := fmt.Sprintf("mysql service %q already exists and is not managed by APIManager", mysqlservice.Name)
					c.recorder.Event(apimanager, corev1.EventTypeWarning, "ErrResourceExists", msg)
					return fmt.Errorf(msg)
				}

				//for instance mysql deployment
				if apimanager.Spec.Replicas != nil && *apimanager.Spec.Replicas != *mysqldeployment.Spec.Replicas {
					//y:= pattern1.AssignMysqlConfigMapValues(apimanager,configmap)
					klog.V(4).Infof("APIManager %s replicas: %d, deployment2 replicas: %d", name, *apimanager.Spec.Replicas, *mysqldeployment.Spec.Replicas)
					mysqldeployment, err = c.kubeclientset.AppsV1().Deployments(apimanager.Namespace).Update(mysql.MysqlDeployment(apimanager))
				}

				//for mysql deployment
				err = c.updateApimanagerStatus(apimanager, mysqldeployment)
				if err != nil {
					return err
				}
			}

		}

		//////////finally update the deployment resources after done checking

		c.recorder.Event(apimanager, corev1.EventTypeNormal, "synced", "APIManager synced successfully")
		return nil

	}

	return nil
}

func (c *Controller) updateApiMangerStatusForStatefulSet(apimanager *apimv1alpha1.APIManager, statefulSet *appsv1.StatefulSet) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	apimanagerCopy := apimanager.DeepCopy()
	apimanagerCopy.Status.AvailableReplicas = statefulSet.Status.CurrentReplicas
	// If the CustomResourceSubresources feature gate is not enabled,
	// we must use Update instead of UpdateStatus to update the Status block of the Apimanager resource.
	// UpdateStatus will not allow changes to the Spec of the resource,
	// which is ideal for ensuring nothing other than resource status has been updated.
	_, err := c.sampleclientset.ApimV1alpha1().APIManagers(apimanager.Namespace).Update(apimanagerCopy)
	return err
}

func (c *Controller) updateApimanagerStatus(apimanager *apimv1alpha1.APIManager, deployment *appsv1.Deployment) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	apimanagerCopy := apimanager.DeepCopy()
	apimanagerCopy.Status.AvailableReplicas = deployment.Status.AvailableReplicas
	// If the CustomResourceSubresources feature gate is not enabled,
	// we must use Update instead of UpdateStatus to update the Status block of the Apimanager resource.
	// UpdateStatus will not allow changes to the Spec of the resource,
	// which is ideal for ensuring nothing other than resource status has been updated.
	_, err := c.sampleclientset.ApimV1alpha1().APIManagers(apimanager.Namespace).Update(apimanagerCopy)
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
		// If this object is not owned by a Apimanager, we should not do anything more with it.
		if ownerRef.Kind != "APIManager" {
			return
		}

		apimanager, err := c.apimanagerslister.APIManagers(object.GetNamespace()).Get(ownerRef.Name)
		if err != nil {
			klog.V(4).Infof("ignoring orphaned object '%s' of apimanager '%s'", object.GetSelfLink(), ownerRef.Name)
			return
		}

		c.enqueueApimanager(apimanager)
		return
	}
}
