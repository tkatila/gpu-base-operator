/*
Copyright 2025 Intel Corporation. All Rights Reserved.

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
	v1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterPolicySpec defines the desired state of ClusterPolicy.
type ClusterPolicySpec struct {
	// Which type of resource registration is to be used: Device plugin (dp) or Dynamic Resource Allocation (dra).
	// +kubebuilder:validation:Enum=dp;dra
	ResourceRegistration string `json:"resourceRegistration"`

	// To enable resource monitoring via XPU or not. Deploys GPU Plugin or DRA with monitoring enabled and
	// XPU Manager DaemonSet if true.
	ResourceMonitoring bool `json:"resourceMonitoring,omitempty"`

	// Use NFD rule to label nodes.
	UseNFDLabeling bool `json:"useNFDLabeling,omitempty"`

	// Deploy Kubernetes components to integrate with Prometheus.
	PrometheusIntegration bool `json:"prometheusIntegration,omitempty"`

	// Set up Kueue queues for node resources
	// +optional
	EnableKueue bool `json:"enableKueue,omitempty"`

	// Define Kueue queues
	Kueue *KueueQueueSpec `json:"kueue,omitempty"`

	// Enable health monitoring in DP/DRA
	// These values are applied to all the Intel GPU devices in the cluster.
	// Mechanism to monitor the values differ between DP and DRA. DP uses LevelZero API
	// directly, while DRA relies on the health status provided by XPU Manager.
	HealthinessSpec *HealthinessSpec `json:"health,omitempty"`

	// +optional
	DynamicResourceAllocationSpec DynamicResourceAllocationSpec `json:"dra"`
	// +optional
	DevicePluginSpec DevicePluginSpec `json:"dp"`
	// +optional
	XpuManagerSpec XpuManagerSpec `json:"xpu"`

	// DriverXeSpec configures out-of-tree Xe KMD deployment via the Kernel Module Management (KMM) operator.
	// +optional
	DriverXeSpec *DriverXeSpec `json:"driverXe,omitempty"`

	// Pull secret is shared with all the deployments.
	// +optional
	PullSecret *v1.LocalObjectReference `json:"pullSecret,omitempty"`

	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	Tolerations  []v1.Toleration   `json:"tolerations,omitempty"`

	// LogLevel to overwrite the default log level of the components.
	// +kubebuilder:validation:Range=0:4
	// +kubebuilder:validation:Default=2
	LogLevel int32 `json:"logLevel,omitempty"`
}

// DynamicResourceAllocationSpec defines the desired state of DynamicResourceAllocation.
type DynamicResourceAllocationSpec struct {
	Image string `json:"image,omitempty"`

	LogLevel int32 `json:"logLevel,omitempty"`

	// Enable DRA Pod's health check.
	// +kubebuilder:default=true
	PodHealthCheck bool `json:"podHealthCheck,omitempty"`

	// Allow DRA plugin to bind/unbind devices from xe/i915 driver to vfio/xe-vfio driver and back.
	// Needed if cluster is supposed to support dynamic switching from drivers. Not needed, if hosts are
	// preconfigured to either target KubeVirt or normal workloads.
	ManageBinding bool `json:"manageBinding,omitempty"`
}

// HealthinessSpec defines the thresholds for health monitoring.
type HealthinessSpec struct {
	// Not supported by Device Plugin
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=3600
	// +kubebuilder:default:=5
	CheckIntervalSeconds int32 `json:"checkIntervalSeconds,omitempty"`

	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=130
	// +kubebuilder:default:=100
	CoreTemperatureThreshold int32 `json:"coreTemperatureThreshold,omitempty"`

	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=130
	// +kubebuilder:default:=100
	MemoryTemperatureThreshold int32 `json:"memoryTemperatureThreshold,omitempty"`
}

// DevicePluginSpec defines the desired state of DevicePlugin.
type DevicePluginSpec struct {
	// Container image for the GPU plugin
	PluginImage string `json:"plugin,omitempty"`
	// Container image for the Level Zero companion container
	// Deprecated: LevelzeroImage is no longer used and this configuration item will be removed in the future.
	LevelzeroImage string `json:"levelzero,omitempty"`

	// AllowIDs and DenyIDs are used to control which devices are registered as resources.
	// Allow or deny certain PCI Device IDs. Both cannot be used together. Format is '0xabcd'.
	AllowIDs []string `json:"allowIDs,omitempty"`
	DenyIDs  []string `json:"denyIDs,omitempty"`

	// ByPathMode controls DRI by-path entries are exposed by the plugin.
	// +kubebuilder:validation:Enum=single;all;none
	ByPathMode string `json:"byPathMode,omitempty"`

	// +kubebuilder:validation:Range=0:4
	// +kubebuilder:validation:Default=1
	LogLevel int32 `json:"logLevel,omitempty"`
}

// XpuManagerSpec defines the desired state of XpuManager.
type XpuManagerSpec struct {
	Image string `json:"image,omitempty"`

	// +kubebuilder:validation:Range=0:3
	// +kubebuilder:validation:Default=2
	LogLevel int32 `json:"logLevel,omitempty"`

	// ConfigMapOverride allows overriding the default OpenTelemetry Collector configuration used by XPU Manager.
	// Configmap has to be in the same namespace as the operator and contain a key "config.yaml" with the configuration content.
	// The value should be a YAML string containing the configuration. If not set, a default configuration will be used.
	ConfigMapOverride string `json:"configMapOverride,omitempty"`

	// Set monitoring resource name for Device Plugin use. If not set, the default resource
	// name "gpu.intel.com/monitoring" will be used.
	// +kubebuilder:validation:Enum=i915_monitoring;xe_monitoring;monitoring
	MonitoringResource string `json:"monitoringResource,omitempty"`
}

// DriverXeSpec configures out-of-tree Xe KMD deployment via the Kernel Module Management (KMM) operator.
// Nodes targeted by the OoT KMD are determined by the nodeSelector at the root of the ClusterPolicy.
type DriverXeSpec struct {
	// Enable creates the KMM Module CR and the accompanying resources to build the out-of-tree Xe KMD.
	// Requires KMM operator to be pre-installed in the cluster.
	Enable bool `json:"enable,omitempty"`

	// ImageRepoSecret is the name of the Secret used to pull/push module images.
	// Needed when working with private registries. Not necessary if deployment in OCP
	// and using ImageStream for module images.
	// +optional
	ImageRepoSecret *v1.LocalObjectReference `json:"imageRepoSecret,omitempty"`

	// SkipTLSVerify disables TLS certificate verification for the image registry.
	// Typically required for self-signed registries.
	SkipTLSVerify bool `json:"skipTLSVerify,omitempty"`

	// ContainerImageBase specifies where the image is or will be located.
	// If the image does not exist, the operator will attempt to build it and
	// push it to the registry specified by ContainerImageBase.
	// Operator adds a tag to differentiate between OS/kernel targets and Xe KMD versions.
	ContainerImageBase string `json:"containerImageBase,omitempty"`

	// NodeSelector defines which nodes the out-of-tree Xe KMD will be deployed to.
	// If clusterpolicy uses useNFDLabeling, operator will deploy a rule that labels nodes and
	// updates the nodeSelector accordingly.
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// KernelMappings defines one entry per OS/kernel target.
	// If empty, operator tries to detect the OS and kernel from the nodes via NFD.
	// +optional
	KernelMappings []XeKernelMappingSpec `json:"kernelMappings,omitempty"`
}

// XeOSTarget enumerates the supported OS targets for the out-of-tree Xe KMD.
// +kubebuilder:validation:Enum=ubuntu24.04;ubuntu26.04;rhel9.8;rhel10.1;rhel10.2
type XeOSTarget string

const (
	XeOSUbuntu2404 XeOSTarget = "ubuntu24.04"
	XeOSUbuntu2604 XeOSTarget = "ubuntu26.04"
	XeOSRHEL9_8    XeOSTarget = "rhel9.8"
	XeOSRHEL10_1   XeOSTarget = "rhel10.1"
	XeOSRHEL10_2   XeOSTarget = "rhel10.2"
)

// XeKernelMappingSpec defines the out-of-tree Xe KMD configuration for one OS/kernel target.
type XeKernelMappingSpec struct {
	// OSTarget identifies the operating system this mapping applies to.
	// The operator derives defaults for KernelRegexp, ContainerImage,
	// module loading order, and build toolchain from this value.
	OSTarget XeOSTarget `json:"osTarget"`

	// Build describes how to build the module image at cluster time via KMM.
	// Only needed to pin a specific xe backport version, change kernel pattern etc.
	// +optional
	Build *XeKMDBuildSpec `json:"build,omitempty"`
}

// XeKMDBuildSpec defines parameters for building the out-of-tree Xe KMD image.
// Variables inside the spec are optionals meaning the operator will auto-fill them
// based on the OSTarget if not set.
type XeKMDBuildSpec struct {
	// KernelRegexp defines the kernel pattern targeted by this spec.
	// Example: '6\.8\.0-[0-9]+-generic' for Ubuntu 24.04.
	// +optional
	KernelRegexp string `json:"kernelRegexp,omitempty"`

	// InTreeModulesToRemove overrides the default list of in-tree modules that
	// must be unloaded before loading the OOT module.
	// +optional
	InTreeModulesToRemove []string `json:"inTreeModulesToRemove,omitempty"`

	// XeTag is the xekmd-backports release tag to fetch and build from.
	// Example: "xebr_v6.17.13.48_260409.3"
	XeTag string `json:"xeTag"`

	// XeSHA is the expected SHA256 checksum of the source tarball for XeTag.
	XeSHA string `json:"xeSHA"`

	// FirmwareFiles is an array of firmware filenames to embed.
	// Example: [bmg_guc_70.bin, bmg_huc.bin, fan_control_8086_e20b_8086_1100.bin]
	// +optional
	FirmwareFiles []string `json:"firmwareFiles,omitempty"`

	// ExtraModules is an array of additional in-tree .ko files
	// to copy into the module image alongside the OOT Xe driver.
	// TODO: List of modules should be able to get from "modinfo" for xe kmd.
	// +optional
	ExtraModules []string `json:"extraModules,omitempty"`
}

// ClusterPolicyStatus defines the observed state of ClusterPolicy.
type ClusterPolicyStatus struct {
	DevicePluginStatus string   `json:"devicePluginStatus,omitempty"`
	DRAStatus          string   `json:"draStatus,omitempty"`
	XPUManagerStatus   string   `json:"xpuManagerStatus,omitempty"`
	KMMStatus          string   `json:"kmmStatus,omitempty"`
	Errors             []string `json:"errors,omitempty"`
}

// KueueQueueSpec defines Kueue cluster and local queues
type KueueQueueSpec struct {
	// Cluster queues for dividing resources evenly
	EqualResources []ClusterQueueSpec `json:"equalResources"`
}

// ClusterQueueSpec defines a Kueue ClusterQueues
type ClusterQueueSpec struct {
	// Name of the cluster queue
	Name string `json:"name"`
	// List of Kueue LocalQueues to create for this ClusterQueue
	LocalQueues []LocalQueueSpec `json:"localQueues"`
}

// LocalQueueSpec defines a Kueue Local Queue
type LocalQueueSpec struct {
	// Name of the cluster queue
	Name string `json:"name"`
	// Namespace for the local queue
	Namespace string `json:"namespace"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=clusterpolicies,scope=Cluster
// +kubebuilder:printcolumn:name="DP",type=string,JSONPath=`.status.devicePluginStatus`
// +kubebuilder:printcolumn:name="DRA",type=string,JSONPath=`.status.draStatus`
// +kubebuilder:printcolumn:name="XPU",type=string,JSONPath=`.status.xpuManagerStatus`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
// +operator-sdk:csv:customresourcedefinitions:displayName="Intel GPU Cluster Policy"
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterPolicy is the Schema for the clusterpolicies API.
type ClusterPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterPolicySpec   `json:"spec,omitempty"`
	Status ClusterPolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterPolicyList contains a list of ClusterPolicy.
type ClusterPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterPolicy `json:"items"`
}
