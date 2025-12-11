/*
Copyright 2025 The KubeVela Authors.

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

// Package components provides Go implementations of KubeVela's built-in component definitions.
package components

// Env represents an environment variable configuration.
type Env struct {
	// Name is the environment variable name.
	Name string
	// Value is the direct value of the environment variable.
	Value *string
	// ValueFrom specifies a source for the environment variable value.
	ValueFrom *EnvValueFrom
}

// EnvValueFrom specifies a source for an environment variable value.
type EnvValueFrom struct {
	// SecretKeyRef selects a key of a secret.
	SecretKeyRef *KeyRef
	// ConfigMapKeyRef selects a key of a config map.
	ConfigMapKeyRef *KeyRef
}

// KeyRef references a key from a Secret or ConfigMap.
type KeyRef struct {
	// Name is the name of the Secret or ConfigMap.
	Name string
	// Key is the key to select from.
	Key string
}

// Port represents a port configuration.
type Port struct {
	// Port is the port number to expose on the pod's IP address.
	Port int
	// ContainerPort is the container port to connect to (defaults to Port if not specified).
	ContainerPort *int
	// Name is the name of the port.
	Name *string
	// Protocol is the protocol for the port (TCP, UDP, or SCTP). Defaults to TCP.
	Protocol string
	// Expose indicates if the port should be exposed as a Service.
	Expose bool
	// NodePort is the exposed node port (only valid when ExposeType is NodePort).
	NodePort *int
}

// VolumeMounts contains different types of volume mounts.
type VolumeMounts struct {
	// PVC mounts PersistentVolumeClaim volumes.
	PVC []PVCMount
	// ConfigMap mounts ConfigMap volumes.
	ConfigMap []ConfigMapMount
	// Secret mounts Secret volumes.
	Secret []SecretMount
	// EmptyDir mounts EmptyDir volumes.
	EmptyDir []EmptyDirMount
	// HostPath mounts HostPath volumes.
	HostPath []HostPathMount
}

// PVCMount represents a PVC volume mount.
type PVCMount struct {
	Name      string
	MountPath string
	SubPath   *string
	ClaimName string
}

// ConfigMapMount represents a ConfigMap volume mount.
type ConfigMapMount struct {
	Name        string
	MountPath   string
	SubPath     *string
	DefaultMode int
	CMName      string
	Items       []VolumeItem
}

// SecretMount represents a Secret volume mount.
type SecretMount struct {
	Name        string
	MountPath   string
	SubPath     *string
	DefaultMode int
	SecretName  string
	Items       []VolumeItem
}

// EmptyDirMount represents an EmptyDir volume mount.
type EmptyDirMount struct {
	Name      string
	MountPath string
	SubPath   *string
	Medium    string
}

// HostPathMount represents a HostPath volume mount.
type HostPathMount struct {
	Name      string
	MountPath string
	SubPath   *string
	Path      string
}

// VolumeItem represents a key-to-path mapping for ConfigMap or Secret volumes.
type VolumeItem struct {
	Key  string
	Path string
	Mode int
}

// HealthProbe represents container health probe configuration.
type HealthProbe struct {
	// Exec specifies a command-based health check.
	Exec *ExecProbe
	// HTTPGet specifies an HTTP GET-based health check.
	HTTPGet *HTTPGetProbe
	// TCPSocket specifies a TCP socket-based health check.
	TCPSocket *TCPSocketProbe
	// InitialDelaySeconds is the delay before the first probe.
	InitialDelaySeconds int
	// PeriodSeconds is how often to perform the probe.
	PeriodSeconds int
	// TimeoutSeconds is the probe timeout.
	TimeoutSeconds int
	// SuccessThreshold is the minimum consecutive successes.
	SuccessThreshold int
	// FailureThreshold is the minimum consecutive failures.
	FailureThreshold int
}

// ExecProbe specifies a command-based health check.
type ExecProbe struct {
	Command []string
}

// HTTPGetProbe specifies an HTTP GET-based health check.
type HTTPGetProbe struct {
	Path        string
	Port        int
	Host        *string
	Scheme      string
	HTTPHeaders []HTTPHeader
}

// HTTPHeader represents an HTTP header for health probes.
type HTTPHeader struct {
	Name  string
	Value string
}

// TCPSocketProbe specifies a TCP socket-based health check.
type TCPSocketProbe struct {
	Port int
}

// HostAlias represents host alias configuration.
type HostAlias struct {
	IP        string
	Hostnames []string
}

// ResourceLimit represents resource limits.
type ResourceLimit struct {
	CPU    *string
	Memory *string
}

// NewDefaultHealthProbe creates a HealthProbe with default values.
func NewDefaultHealthProbe() *HealthProbe {
	return &HealthProbe{
		InitialDelaySeconds: 0,
		PeriodSeconds:       10,
		TimeoutSeconds:      1,
		SuccessThreshold:    1,
		FailureThreshold:    3,
	}
}

// Helper function to create a string pointer
func StringPtr(s string) *string {
	return &s
}

// Helper function to create an int pointer
func IntPtr(i int) *int {
	return &i
}

