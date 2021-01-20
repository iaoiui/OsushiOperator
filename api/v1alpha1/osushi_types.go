/*


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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OsushiSpec defines the desired state of Osushi
type OsushiSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Osushi. Edit Osushi_types.go to remove/update
	// Size is the size of the osushi deployment
	Size  int32  `json:"size"`
	Emoji string `json:"emoji,omitempty"`
	// EndressOsushi          bool   `json:"endressOsushi,omitempty"`
	// TraditionalKaitenSushi bool   `json:"traditionalKaitenSushi,omitempty"`
	// Modes: endressOsushi, traditionalKaitenSushi
	Mode               string `json:"mode,omitempty"`
	OsushiSpeed        int32  `json:"osushiSpeed,omitempty"`
	LengthOfOsushiLane int32  `json:"lengthOfOsushiLane,omitempty"`
}

// OsushiStatus defines the observed state of Osushi
type OsushiStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Freshness string `json:"freshness,omitempty"`
	Reachable bool   `json:"reacheable,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Size",type=integer,JSONPath=`.spec.size`
// +kubebuilder:printcolumn:name="Emoji",type=string,JSONPath=`.spec.emoji`
// Osushi is the Schema for the osushis API
type Osushi struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OsushiSpec   `json:"spec,omitempty"`
	Status OsushiStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// OsushiList contains a list of Osushi
type OsushiList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Osushi `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Osushi{}, &OsushiList{})
}
