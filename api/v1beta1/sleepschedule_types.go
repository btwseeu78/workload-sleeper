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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type SleepStatus string

const (
	SleepStatusPending SleepStatus = "Pending"
	SleepStatusPaused  SleepStatus = "Paused"
	SleepStatusResumed SleepStatus = "Resumed"
	SleepStatusAbandon SleepStatus = "Abandoned"
)

// SleepScheduleSpec defines the desired state of SleepSchedule
type SleepScheduleSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// +kubebuilder:validation:Optional
	Schedule *scheduleInfo `json:"schedule"`
	//Target            targetInfo            `json:"target"`
	NamespaceSelector *metav1.LabelSelector `json:"namespaceSelector,omitempty"`
}

// SleepScheduleStatus defines the observed state of SleepSchedule
type SleepScheduleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Last triggered name
	LastTriggered string `json:"lastTriggered"`
	//// Triggered bool
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="Pending"
	CurrStatus SleepStatus `json:"currStatus,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="CurrStatus",type="string",JSONPath=".status.currStatus"
// +kubebuilder:printcolumn:name="LastTriggered",type="string",JSONPath=".status.lastTriggered"
// SleepSchedule is the Schema for the sleepschedules API
type SleepSchedule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// Describes Details About how the object would be selected for sync
	Spec   SleepScheduleSpec   `json:"spec,omitempty"`
	Status SleepScheduleStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SleepScheduleList contains a list of SleepSchedule
type SleepScheduleList struct {
	metav1.TypeMeta `json:",inline"`

	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SleepSchedule `json:"items"`
}

type scheduleInfo struct {

	// +kubebuilder:default="UTC"
	// +kubebuilder:validation:Optional
	TimeZone string `json:"timeZone,omitempty"`
	// Used to give Pause In Case of Custom Needs. Once This NeedTo Unset Manually
	// +kubebuilder:validation:Optional
	PauseScheduled bool `json:"pauseScheduled,omitempty"`
	// Start date in format dd/mm/yyyy, by default we will consider Today.
	// +kubebuilder:validation:Optional
	SleepStartDate string `json:"sleepStartDate,omitempty"`
	// One Year from Start Date
	// +kubebuilder:validation:Optional
	SleepEndDate string `json:"sleepEndDate,omitempty"`
	// Pause start time
	// +kubebuilder:validation:Optional
	SleepStartTime string `json:"sleepStartTime,omitempty"`
	// End Time or Resume time
	// +kubebuilder:validation:Optional
	SleepEndTime string `json:"sleepEndTime,omitempty"`
}

func init() {
	SchemeBuilder.Register(&SleepSchedule{}, &SleepScheduleList{})
}
