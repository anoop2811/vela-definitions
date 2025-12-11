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

package components

import (
	"github.com/oam-dev/kubevela/pkg/definition/defkit"
)

// SharedHelpers provides reusable template helpers for common patterns
// across component definitions. These helpers generate consistent CUE
// code for volume mounts, image pull secrets, ports, etc.

// --- Volume Mount Helpers ---

// volumeMountSources are the standard volume mount source fields.
var volumeMountSources = []string{"pvc", "configMap", "secret", "emptyDir", "hostPath"}

// ContainerMountsHelper creates a helper for container volumeMounts.
// Transforms volumeMounts from multiple sources (pvc, configMap, secret, emptyDir, hostPath)
// into the container volumeMounts format: [{name, mountPath, subPath?}]
//
// Usage:
//
//	containerMounts := ContainerMountsHelper(tpl, volumeMounts)
//	deployment.SetIf(volumeMounts.IsSet(), "spec...volumeMounts", containerMounts)
func ContainerMountsHelper(tpl *defkit.Template, volumeMounts defkit.Value) *defkit.HelperVar {
	return tpl.Helper("containerMountsArray").
		FromFields(volumeMounts, volumeMountSources...).
		Pick("name", "mountPath").
		PickIf(defkit.ItemFieldIsSet("subPath"), "subPath").
		Build()
}

// ContainerMountsDedupedHelper creates a helper for deduplicated container volumeMounts.
// Same as ContainerMountsHelper but removes duplicate entries by name.
//
// Deprecated: Use ContainerMountsHelper instead. Container volumeMounts should NOT be
// deduplicated because it's valid to mount the same volume at multiple paths.
// Only pod volumes need deduplication (use PodVolumesDedupedHelper for that).
func ContainerMountsDedupedHelper(tpl *defkit.Template, volumeMounts defkit.Value) *defkit.HelperVar {
	mountsArray := tpl.Helper("mountsArray").
		FromFields(volumeMounts, volumeMountSources...).
		Pick("name", "mountPath").
		PickIf(defkit.ItemFieldIsSet("subPath"), "subPath").
		Build()

	return tpl.Helper("deDupMountsArray").
		FromHelper(mountsArray).
		Dedupe("name").
		Build()
}

// PodVolumesHelper creates a helper for pod volumes.
// Transforms volumeMounts from multiple sources into Kubernetes volume specs.
// Each source type maps to its corresponding volume spec format:
//   - pvc -> persistentVolumeClaim
//   - configMap -> configMap
//   - secret -> secret
//   - emptyDir -> emptyDir
//   - hostPath -> hostPath
//
// Usage:
//
//	podVolumes := PodVolumesHelper(tpl, volumeMounts)
//	deployment.SetIf(volumeMounts.IsSet(), "spec...volumes", podVolumes)
func PodVolumesHelper(tpl *defkit.Template, volumeMounts defkit.Value) *defkit.HelperVar {
	return tpl.Helper("volumesList").
		FromFields(volumeMounts, volumeMountSources...).
		MapBySource(podVolumeMappings()).
		Build()
}

// PodVolumesDedupedHelper creates a helper for deduplicated pod volumes.
// Same as PodVolumesHelper but removes duplicate entries by name.
func PodVolumesDedupedHelper(tpl *defkit.Template, volumeMounts defkit.Value) *defkit.HelperVar {
	volumesList := tpl.Helper("volumesList").
		FromFields(volumeMounts, volumeMountSources...).
		MapBySource(podVolumeMappings()).
		Build()

	return tpl.Helper("deDupVolumesList").
		FromHelper(volumesList).
		Dedupe("name").
		Build()
}

// podVolumeMappings returns the standard field mappings for pod volumes.
func podVolumeMappings() map[string]defkit.FieldMap {
	return map[string]defkit.FieldMap{
		"pvc": {
			"name":                  defkit.FieldRef("name"),
			"persistentVolumeClaim": defkit.Nested(defkit.FieldMap{"claimName": defkit.FieldRef("claimName")}),
		},
		"configMap": {
			"name": defkit.FieldRef("name"),
			"configMap": defkit.Nested(defkit.FieldMap{
				"name":        defkit.FieldRef("cmName"),
				"defaultMode": defkit.FieldRef("defaultMode"),
				"items":       defkit.Optional("items"),
			}),
		},
		"secret": {
			"name": defkit.FieldRef("name"),
			"secret": defkit.Nested(defkit.FieldMap{
				"secretName":  defkit.FieldRef("secretName"),
				"defaultMode": defkit.FieldRef("defaultMode"),
				"items":       defkit.Optional("items"),
			}),
		},
		"emptyDir": {
			"name":     defkit.FieldRef("name"),
			"emptyDir": defkit.Nested(defkit.FieldMap{"medium": defkit.FieldRef("medium")}),
		},
		"hostPath": {
			"name":     defkit.FieldRef("name"),
			"hostPath": defkit.Nested(defkit.FieldMap{"path": defkit.FieldRef("path")}),
		},
	}
}

// --- Image Pull Secrets Helper ---

// ImagePullSecretsTransform transforms a string array of secret names
// into the Kubernetes imagePullSecrets format: [{name: "secret1"}, ...]
//
// Usage:
//
//	pullSecrets := ImagePullSecretsTransform(imagePullSecrets)
//	deployment.SetIf(imagePullSecrets.IsSet(), "spec...imagePullSecrets", pullSecrets)
func ImagePullSecretsTransform(imagePullSecrets defkit.Value) *defkit.CollectionOp {
	return defkit.Each(imagePullSecrets).Wrap("name")
}

// --- Port Helpers ---

// ContainerPortsTransform transforms port definitions to container port format.
// Maps: {port, name?, protocol} -> {containerPort, name, protocol}
// Name defaults to "port-{port}" if not specified.
//
// Usage:
//
//	containerPorts := ContainerPortsTransform(ports)
//	deployment.SetIf(ports.IsSet(), "spec...ports", containerPorts)
func ContainerPortsTransform(ports defkit.Value) *defkit.CollectionOp {
	return defkit.Each(ports).Map(defkit.FieldMap{
		"containerPort": defkit.FieldRef("port"),
		"name":          defkit.FieldRef("name").Or(defkit.Format("port-%v", defkit.FieldRef("port"))),
		"protocol":      defkit.FieldRef("protocol"),
	})
}

// ServicePortsTransform transforms port definitions to Service port format.
// Maps: {port, name?, protocol} -> {port, targetPort, name, protocol}
// Name defaults to "port-{port}" if not specified.
//
// Usage:
//
//	servicePorts := ServicePortsTransform(ports)
//	service.SetIf(ports.IsSet(), "spec.ports", servicePorts)
func ServicePortsTransform(ports defkit.Value) *defkit.CollectionOp {
	return defkit.Each(ports).Map(defkit.FieldMap{
		"port":       defkit.FieldRef("port").Or(defkit.FieldRef("containerPort")),
		"targetPort": defkit.FieldRef("port").Or(defkit.FieldRef("containerPort")),
		"name":       defkit.FieldRef("name").Or(defkit.Format("port-%v", defkit.FieldRef("port").Or(defkit.FieldRef("containerPort")))),
		"protocol":   defkit.FieldRef("protocol"),
	})
}

// --- Common Parameter Definitions ---

// CommonVolumeParams returns the standard volumeMounts parameter definition.
func CommonVolumeParams() defkit.Param {
	return defkit.Object("volumeMounts").Description("Volume mount configurations")
}

// CommonImagePullSecretsParam returns the standard imagePullSecrets parameter.
func CommonImagePullSecretsParam() defkit.Param {
	return defkit.StringList("imagePullSecrets").Description("Specify image pull secrets for your service")
}

// CommonProbeParams returns liveness and readiness probe parameters.
func CommonProbeParams() (livenessProbe, readinessProbe defkit.Param) {
	livenessProbe = defkit.Object("livenessProbe").
		Description("Instructions for assessing whether the container is alive")
	readinessProbe = defkit.Object("readinessProbe").
		Description("Instructions for assessing whether the container is in a suitable state to serve traffic")
	return
}

// CommonResourceParams returns cpu and memory parameters.
func CommonResourceParams() (cpu, memory defkit.Param) {
	cpu = defkit.String("cpu").
		Description("Number of CPU units for the service, like `0.5` (0.5 CPU core), `1` (1 CPU core)")
	memory = defkit.String("memory").
		Description("Specifies the attributes of the memory resource required for the container.")
	return
}
