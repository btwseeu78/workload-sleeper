/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	greenworkloadv1beta1 "github.com/btwseeu78/workload-sleeper/api/v1beta1"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strconv"
)

// WorkloadScalerReconciler reconciles a WorkloadScaler object
type WorkloadScalerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=greenworkload.platform.io,resources=sleepschedules,verbs=get;list;watch
//+kubebuilder:rbac:groups=greenworkload.platform.io,resources=sleepschedules/status,verbs=get
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;patch;update;
//+kubebuilder:rbac:groups=apps,resources=deployments/scale,verbs=get;list;watch;patch;update;

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the WorkloadScaler object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.2/pkg/reconcile
func (r *WorkloadScalerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("workloadscaler", req.NamespacedName)

	// get the sleepschedule status
	sleepSchedule := &greenworkloadv1beta1.SleepSchedule{}
	err := r.Get(ctx, req.NamespacedName, sleepSchedule)

	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("WorkloadScaler resource not found. Ignoring sleep.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to fetch sleep schedule resource")
		return ctrl.Result{}, err
	}
	log.Info("Reconciling sleep schedule resource", "Namespace", req.Namespace, "Name", req.Name)

	// The main Code Block
	sleepCopy := sleepSchedule.DeepCopy()
	currStatus := sleepCopy.Status.CurrStatus
	if currStatus == greenworkloadv1beta1.SleepStatusPaused {
		lstList := sleepCopy.Spec.NamespaceSelector

		selector, err := metav1.LabelSelectorAsSelector(lstList)
		if err != nil {
			log.Error(err, "Failed to convert LabelSelector to Selector")
			return reconcile.Result{}, err
		}

		lstOptions := client.ListOptions{
			LabelSelector: selector,
			Namespace:     req.Namespace,
		}
		deployList := &v1.DeploymentList{}
		err = r.List(ctx, deployList, &lstOptions)
		if err != nil {
			log.Error(err, "Failed to list Deployments")
			return reconcile.Result{}, err
		}
		for _, deploy := range deployList.Items {
			log.Info("Deployment", "Namespace", deploy.Namespace, "Name", deploy.Name)

			tempDeploy := deploy.DeepCopy()
			var replicas int32 = 0
			tempDeploy.Spec.Replicas = &replicas
			tempDeploy.ObjectMeta.Labels["old-replica"] = strconv.Itoa(int(replicas))
			err = r.Update(ctx, tempDeploy)
			if err != nil {
				log.Error(err, "Failed to update Deployment", "Namespace", tempDeploy.Namespace, "Name", tempDeploy.Name)
				return reconcile.Result{}, err
			}
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, nil

	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WorkloadScalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&greenworkloadv1beta1.SleepSchedule{}).
		WithEventFilter(&StatusUpdatePredicate{}).
		Complete(r)
}
