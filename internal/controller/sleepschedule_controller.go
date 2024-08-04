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
	"fmt"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"time"

	greenworkloadv1beta1 "github.com/btwseeu78/workload-sleeper/api/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SleepScheduleReconciler reconciles a SleepSchedule object
type SleepScheduleReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=greenworkload.platform.io,resources=sleepschedules,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=greenworkload.platform.io,resources=sleepschedules/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=greenworkload.platform.io,resources=sleepschedules/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SleepSchedule object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.2/pkg/reconcile

func (r *SleepScheduleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("sleepschedule", req.NamespacedName)

	// TODO(user): your logic here
	sleepSchedule := &greenworkloadv1beta1.SleepSchedule{}
	err := r.Get(ctx, req.NamespacedName, sleepSchedule)
	if err != nil && errors.IsNotFound(err) {
		log.Info("The Resource probablyGetting Deleted", "Name", req.Name, "NameSpace", req.Namespace)
		return ctrl.Result{}, nil
	} else if err != nil {
		log.Error(err, "Unable To find the object", "Name", req.Name, "Namespace", req.Namespace)
		return ctrl.Result{}, err
	}
	log.Info("SleepSchedule reconciliation started", "Namespace", req.Namespace, "Name", req.Name)

	// Gather the inputs for time
	startDate := sleepSchedule.Spec.Schedule.SleepStartDate
	endDate := sleepSchedule.Spec.Schedule.SleepEndDate
	startTime := sleepSchedule.Spec.Schedule.SleepStartTime
	endTime := sleepSchedule.Spec.Schedule.SleepEndTime
	timeZone := sleepSchedule.Spec.Schedule.TimeZone
	if startDate == "" {

		startDate = time.Now().Format("2006-01-02")

	}
	if timeZone == "" {
		timeZone = "UTC"
	}

	if endDate == "" {
		endDate = time.Now().Add(24 * time.Hour).Format("2006-01-02")
	}
	if endTime == "" {
		endTime = "18:00"

	}
	if startTime == "" {
		startTime = "09:00"
	}

	startDateTime := startDate + " " + startTime
	endDateTime := endDate + " " + endTime
	layout := "2006-01-02 15:04"
	currDate := time.Now().Format("2006-01-02")
	currDateTime := time.Now().Format(layout)
	todaysStartDateTime := currDate + " " + startTime
	todaysEndDateTime := currDate + " " + endTime

	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error loading time zone: %v", err)
	}

	currFormatedStartDateTime, err := time.Parse(layout, todaysStartDateTime)
	if err != nil {
		return ctrl.Result{}, err
	}

	currFormatedStartDateTime = currFormatedStartDateTime.In(loc)

	currFormatedEndDateTime, err := time.Parse(layout, todaysEndDateTime)
	if err != nil {
		return ctrl.Result{}, err
	}
	currFormatedEndDateTime = currFormatedEndDateTime.In(loc)

	actualstartDateTime, err := time.Parse(layout, startDateTime)
	if err != nil {
		return ctrl.Result{}, err
	}
	actualstartDateTime = actualstartDateTime.In(loc)

	actualEndDateTime, err := time.Parse(layout, endDateTime)
	if err != nil {
		return ctrl.Result{}, err
	}
	actualEndDateTime = actualEndDateTime.In(loc)

	formattedCurrDateTime, err := time.Parse(layout, currDateTime)
	if err != nil {
		return ctrl.Result{}, err
	}
	formattedCurrDateTime = formattedCurrDateTime.In(loc)

	if sleepSchedule.Status.CurrStatus == greenworkloadv1beta1.SleepStatusAbandon {
		log.Info("Curr Status is Abandoned", "Curr Status", sleepSchedule.Status.CurrStatus)
		return ctrl.Result{}, nil
	}

	if formattedCurrDateTime.Before(actualstartDateTime) {
		getDiff := actualstartDateTime.Sub(formattedCurrDateTime).Seconds()
		log.Info("Start Time has not elapsed")
		return ctrl.Result{RequeueAfter: time.Duration(getDiff) * time.Second}, nil
	}

	if formattedCurrDateTime.After(actualEndDateTime) {
		log.Info("curr time outside elapsed")
		sleepSchedule.Status.CurrStatus = greenworkloadv1beta1.SleepStatusAbandon
		sleepSchedule.Status.LastTriggered = currDateTime
		err := r.Status().Update(ctx, sleepSchedule)
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if formattedCurrDateTime.Before(currFormatedStartDateTime) {
		log.Info("Start Time has not elapsed")
		sleepSchedule.Status.CurrStatus = greenworkloadv1beta1.SleepStatusResumed
		sleepSchedule.Status.LastTriggered = currDateTime
		err := r.Status().Update(ctx, sleepSchedule)
		if err != nil {
			return ctrl.Result{}, err
		}
		getDiff := formattedCurrDateTime.Sub(currFormatedStartDateTime).Seconds()

		return ctrl.Result{RequeueAfter: time.Duration(getDiff) * time.Second}, nil
	}

	if formattedCurrDateTime.After(currFormatedEndDateTime) {
		log.Info("End Time elapsed")
		sleepSchedule.Status.CurrStatus = greenworkloadv1beta1.SleepStatusResumed
		sleepSchedule.Status.LastTriggered = currDateTime
		err := r.Status().Update(ctx, sleepSchedule)
		if err != nil {
			return ctrl.Result{}, err
		}
		nextDate := currFormatedStartDateTime.Add(time.Hour * 24)

		getDiff := nextDate.Sub(formattedCurrDateTime).Seconds()
		return ctrl.Result{RequeueAfter: time.Duration(getDiff) * time.Second}, nil

	}

	if formattedCurrDateTime.After(actualstartDateTime) && formattedCurrDateTime.Before(actualEndDateTime) {
		log.Info("Curr time within  the window workload will be paused")
		sleepSchedule.Status.LastTriggered = currDateTime
		sleepSchedule.Status.CurrStatus = greenworkloadv1beta1.SleepStatusPaused
		err := r.Status().Update(ctx, sleepSchedule)
		if err != nil {
			return ctrl.Result{}, err
		}
		getDiff := currFormatedEndDateTime.Sub(formattedCurrDateTime).Seconds()

		return ctrl.Result{RequeueAfter: time.Duration(getDiff) * time.Second}, nil

	}

	log.Info("Curr time within resume window workooad will be resumed")

	sleepSchedule.Status.CurrStatus = greenworkloadv1beta1.SleepStatusResumed
	sleepSchedule.Status.LastTriggered = currDateTime
	err = r.Status().Update(ctx, sleepSchedule)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *SleepScheduleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&greenworkloadv1beta1.SleepSchedule{}).
		Complete(r)
}
