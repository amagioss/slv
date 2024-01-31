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

package v1

import (
	"github.com/amagimedia/slv/core/vaults"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SLVSpec defines the desired state of SLV
type SLVSpec struct {
	vaults.Vault `json:""`
}

// SLVStatus defines the observed state of SLV
type SLVStatus struct {
	Error string `json:"error,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// SLV is the Schema for the slv API
type SLV struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   SLVSpec   `json:"spec"`
	Status SLVStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SLVList contains a list of SLV
type SLVList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SLV `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SLV{}, &SLVList{})
}
