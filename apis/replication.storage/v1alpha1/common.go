/*
Copyright 2022 The Kubernetes-CSI-Addons Authors.

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

// State captures the latest state of the replication operation.
type State string

const (
	// PrimaryState represents the Primary replication state.
	PrimaryState State = "Primary"

	// SecondaryState represents the Secondary replication state.
	SecondaryState State = "Secondary"

	// UnknownState represents the Unknown replication state.
	UnknownState State = "Unknown"
)
