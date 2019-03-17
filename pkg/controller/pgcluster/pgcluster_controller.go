package pgcluster

import (
	"context"
	"fmt"
	"reflect"

	pgclusterv1alpha1 "github.com/jasonodonnell/example-operator/pkg/apis/pgcluster/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_pgcluster")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new PGCluster Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcilePGCluster{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("pgcluster-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource PGCluster
	err = c.Watch(&source.Kind{Type: &pgclusterv1alpha1.PGCluster{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner PGCluster
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &pgclusterv1alpha1.PGCluster{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcilePGCluster{}

// ReconcilePGCluster reconciles a PGCluster object
type ReconcilePGCluster struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a PGCluster object and makes changes based on the state read
// and what is in the PGCluster.Spec
func (r *ReconcilePGCluster) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling PGCluster")

	// Fetch the PGCluster instance
	pgcluster := &pgclusterv1alpha1.PGCluster{}
	err := r.client.Get(context.TODO(), request.NamespacedName, pgcluster)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	fmt.Println(pgcluster)
	if err := pgcluster.Validate(); err != nil {
		reqLogger.Error(err, "PGCluster request invalid")
		return reconcile.Result{}, err
	}

	serviceFound := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: pgcluster.Name, Namespace: pgcluster.Namespace}, serviceFound)
	if err != nil && errors.IsNotFound(err) {
		// Define a new service
		service := r.newServiceForPGCluster(pgcluster)
		reqLogger.Info("Creating a new service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
		err = r.client.Create(context.TODO(), service)
		if err != nil {
			reqLogger.Error(err, "Failed to create new service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get SERVICE")
		return reconcile.Result{}, err
	}

	// Check if the deployment already exists, if not create a new one
	deploymentFound := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: pgcluster.Name, Namespace: pgcluster.Namespace}, deploymentFound)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := r.newDeploymentForPGCluster(pgcluster)
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.client.Create(context.TODO(), dep)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return reconcile.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Deployment")
		return reconcile.Result{}, err
	}

	// Ensure the deployment size is the same as the spec
	size := int32(1)
	if *deploymentFound.Spec.Replicas != size {
		deploymentFound.Spec.Replicas = &size
		err = r.client.Update(context.TODO(), deploymentFound)
		if err != nil {
			reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", deploymentFound.Namespace, "Deployment.Name", deploymentFound.Name)
			return reconcile.Result{}, err
		}
		// Spec updated - return and requeue
		return reconcile.Result{Requeue: true}, nil
	}

	// Update the pgcluster status with the pod names
	// List the pods for this pgcluster's deployment
	podList := &corev1.PodList{}
	labelSelector := labels.SelectorFromSet(labelsForPGCluster(pgcluster.Name))
	listOps := &client.ListOptions{Namespace: pgcluster.Namespace, LabelSelector: labelSelector}
	err = r.client.List(context.TODO(), listOps, podList)
	if err != nil {
		reqLogger.Error(err, "Failed to list pods", "pgcluster.Namespace", pgcluster.Namespace, "pgcluster.Name", pgcluster.Name)
		return reconcile.Result{}, err
	}
	podNames := getPodNames(podList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, pgcluster.Status.Nodes) {
		pgcluster.Status.Nodes = podNames
		err := r.client.Status().Update(context.TODO(), pgcluster)
		if err != nil {
			reqLogger.Error(err, "Failed to update pgcluster status")
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}
