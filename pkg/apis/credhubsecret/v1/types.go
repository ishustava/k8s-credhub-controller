package v1

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const ResourceKind = "CredhubSecret"

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CredhubSecret describes a CredhubSecret resource
type CredhubSecret struct {
	// TypeMeta is the metadata for the resource, like kind and apiversion
	meta_v1.TypeMeta `json:",inline"`
	// ObjectMeta contains the metadata for the particular object, including
	// things like...
	//  - name
	//  - namespace
	//  - self link
	//  - labels
	//  - ... etc ...
	meta_v1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the custom resource spec
	Spec CredhubSecretSpec `json:"spec"`
}

// CredhubSecretSpec is the spec for a CredhubSecret resource
type CredhubSecretSpec struct {
	// Message and SomeValue are example custom spec fields
	//
	// this is where you would put your custom resource data
	Type string `json:"type"`
	Parameters *CredhubParameters `json:"parameters"`
}

type CredhubParameters struct {

}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CredhubSecretList is a list of CredhubSecret resources
type CredhubSecretList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`

	Items []CredhubSecret `json:"items"`
}
