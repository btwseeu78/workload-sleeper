package controller

import (
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

	return true

}

func (p *StatusUpdatePredicate) Generic(e event.GenericEvent) bool {
	return false // Ignore generic events
}
