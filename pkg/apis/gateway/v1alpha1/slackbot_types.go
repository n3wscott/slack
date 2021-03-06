/*
Copyright 2020 The Knative Authors.

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

package v1alpha1

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/apis/duck"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/kmeta"
	"knative.dev/pkg/webhook/resourcesemantics"
)

// +genclient
// +genreconciler
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type Slackbot struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec holds the desired state of the Slackbot (from the client).
	Spec SlackbotSpec `json:"spec"`

	// Status communicates the observed state of the Slackbot (from the controller).
	// +optional
	Status SlackbotStatus `json:"status,omitempty"`
}

// GetGroupVersionKind returns the GroupVersionKind.
func (s *Slackbot) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("Slackbot")
}

var _ resourcesemantics.GenericCRD = (*Slackbot)(nil)

// Check that Slackbot is a runtime.Object.
var _ runtime.Object = (*Slackbot)(nil)

// Check that we can create OwnerReferences to a Slackbot.
var _ kmeta.OwnerRefable = (*Slackbot)(nil)

// Check that Slackbot implements the Conditions duck type.
var _ = duck.VerifyType(&Slackbot{}, &duckv1.Conditions{})

// SlackbotSpec holds the desired state of the Slackbot (from the client).
type SlackbotSpec struct {
	// inherits duck/v1 SourceSpec, which currently provides:
	// * Sink - a reference to an object that will resolve to a domain name or
	//   a URI directly to use as the sink.
	// * CloudEventOverrides - defines overrides to control the output format
	//   and modifications of the event sent to the sink.
	duckv1.SourceSpec `json:",inline"`
}

// SlackbotStatus communicates the observed state of the Slackbot (from the controller).
type SlackbotStatus struct {
	// inherits duck/v1 SourceStatus, which currently provides:
	// * ObservedGeneration - the 'Generation' of the Service that was last
	//   processed by the controller.
	// * Conditions - the latest available observations of a resource's current
	//   state.
	// * SinkURI - the current active sink URI that has been configured for the
	//   Source.
	duckv1.SourceStatus `json:",inline"`

	// AddressStatus is the part where the Slackbot fulfills the Addressable contract.
	duckv1.AddressStatus `json:",inline"`

	Team     *SlackTeamInfo `json:"team,omitempty"`
	Channels SlackChannels  `json:"channels,omitempty"`
	IMs      SlackIMs       `json:"ims,omitempty"`
}

// +k8s:deepcopy-gen=true
type SlackTeamInfo struct {
	ID   string    `json:"id,omitempty"`
	Name string    `json:"name,omitempty"`
	URL  *apis.URL `json:"domain,omitempty"`
}

// +k8s:deepcopy-gen=true
type SlackChannel struct {
	Name     string `json:"name,omitempty"`
	ID       string `json:"id,omitempty"`
	IsMember bool   `json:"isMember,omitempty"`
}

// +k8s:deepcopy-gen=true
type SlackChannels []SlackChannel

// +k8s:deepcopy-gen=true
type SlackIM struct {
	ID   string `json:"id,omitempty"`
	With string `json:"with,omitempty"`
}

// +k8s:deepcopy-gen=true
type SlackIMs []SlackIM

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SlackbotList is a list of Slackbot resources
type SlackbotList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Slackbot `json:"items"`
}

// TODO: implement
func (in *Slackbot) SetDefaults(context.Context) {
	// TODO: add default instance annotations.
}

// TODO: implement
func (in *Slackbot) Validate(context.Context) *apis.FieldError {
	return nil
}
