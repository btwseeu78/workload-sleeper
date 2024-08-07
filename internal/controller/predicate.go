package controller

import (
	greenworkloadv1beta1 "github.com/btwseeu78/workload-sleeper/api/v1beta1"
	"github.com/google/go-cmp/cmp"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

type StatusUpdatePredicate struct{}

func (p *StatusUpdatePredicate) Create(e event.CreateEvent) bool {
	return true // We want to handle creation events if needed
}

func (p *StatusUpdatePredicate) Delete(e event.DeleteEvent) bool {
	return false // Ignore delete events
}

func (p *StatusUpdatePredicate) Update(e event.UpdateEvent) bool {
	oldDeployment, ok1 := e.ObjectOld.(*greenworkloadv1beta1.SleepSchedule)
	newDeployment, ok2 := e.ObjectNew.(*greenworkloadv1beta1.SleepSchedule)

	if !ok1 || !ok2 {
		return false
	}
	if cmp.Equal(oldDeployment.Status.LastTriggered, newDeployment.Status.LastTriggered) {
		return true
	}
	// Compare only the status fields
	return false
}

func (p *StatusUpdatePredicate) Generic(e event.GenericEvent) bool {
	return false // Ignore generic events
}
