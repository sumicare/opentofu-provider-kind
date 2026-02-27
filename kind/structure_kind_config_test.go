/*
   Copyright 2026 Sumicare

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

package kind

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/kind/pkg/apis/config/v1alpha4"
)

// Test data constants for reusability across tests.
const (
	testClusterKind      = "Cluster"
	testAPIVersion       = "kind.x-k8s.io/v1alpha4"
	testControlPlaneRole = "control-plane"
	testWorkerRole       = "worker"
	testNodeImage        = "kindest/node:v1.29.0"
	testAPIServerAddress = "127.0.0.1"
	testAPIServerPort    = 6443
	testHostPath         = "/host/path"
	testContainerPath    = "/container/path"
	testContainerPort    = 80
	testHostPort         = 8080
	testListenAddress    = "0.0.0.0"
	testPodSubnet        = "10.244.0.0/16"
	testServiceSubnet    = "10.96.0.0/12"
)

func TestGetString(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]any
		key      string
		expected string
	}{
		{
			name:     "extracts existing string value",
			input:    map[string]any{"key": "value"},
			key:      "key",
			expected: "value",
		},
		{
			name:     "returns empty for missing key",
			input:    map[string]any{},
			key:      "missing",
			expected: "",
		},
		{
			name:     "returns empty for nil value",
			input:    map[string]any{"key": nil},
			key:      "key",
			expected: "",
		},
		{
			name:     "returns empty for wrong type",
			input:    map[string]any{"key": 123},
			key:      "key",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getString(tt.input, tt.key)
			assert.Equal(
				t,
				tt.expected,
				result,
				"getString should handle all input types correctly",
			)
		})
	}
}

func TestGetInt(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]any
		key      string
		expected int
	}{
		{
			name:     "extracts existing int value",
			input:    map[string]any{"port": testAPIServerPort},
			key:      "port",
			expected: testAPIServerPort,
		},
		{
			name:     "returns zero for missing key",
			input:    map[string]any{},
			key:      "missing",
			expected: 0,
		},
		{
			name:     "returns zero for nil value",
			input:    map[string]any{"port": nil},
			key:      "port",
			expected: 0,
		},
		{
			name:     "returns zero for wrong type",
			input:    map[string]any{"port": "invalid"},
			key:      "port",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getInt(tt.input, tt.key)
			assert.Equal(t, tt.expected, result, "getInt should handle all input types correctly")
		})
	}
}

func TestGetBool(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]any
		key      string
		expected bool
	}{
		{
			name:     "extracts true value",
			input:    map[string]any{"enabled": true},
			key:      "enabled",
			expected: true,
		},
		{
			name:     "extracts false value",
			input:    map[string]any{"enabled": false},
			key:      "enabled",
			expected: false,
		},
		{
			name:     "returns false for missing key",
			input:    map[string]any{},
			key:      "missing",
			expected: false,
		},
		{
			name:     "returns false for nil value",
			input:    map[string]any{"enabled": nil},
			key:      "enabled",
			expected: false,
		},
		{
			name:     "returns false for wrong type",
			input:    map[string]any{"enabled": "true"},
			key:      "enabled",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getBool(tt.input, tt.key)
			assert.Equal(t, tt.expected, result, "getBool should handle all input types correctly")
		})
	}
}

func TestFlattenKindConfig(t *testing.T) {
	tests := []struct {
		input     map[string]any
		validator func(t *testing.T, result *v1alpha4.Cluster)
		name      string
	}{
		{
			name: "basic cluster config",
			input: map[string]any{
				"kind":        testClusterKind,
				"api_version": testAPIVersion,
			},
			validator: func(t *testing.T, result *v1alpha4.Cluster) {
				t.Helper()
				assert.Equal(t, testClusterKind, result.Kind, "Kind field should be set correctly")
				assert.Equal(
					t,
					testAPIVersion,
					result.APIVersion,
					"APIVersion field should be set correctly",
				)
			},
		},
		{
			name: "cluster config with nodes",
			input: map[string]any{
				"kind":        testClusterKind,
				"api_version": testAPIVersion,
				"node": []any{
					map[string]any{"role": testControlPlaneRole},
					map[string]any{"role": testWorkerRole},
				},
			},
			validator: func(t *testing.T, result *v1alpha4.Cluster) {
				t.Helper()
				assert.Len(t, result.Nodes, 2, "should have 2 nodes")
				assert.Equal(
					t,
					v1alpha4.ControlPlaneRole,
					result.Nodes[0].Role,
					"first node should be control-plane",
				)
				assert.Equal(
					t,
					v1alpha4.WorkerRole,
					result.Nodes[1].Role,
					"second node should be worker",
				)
			},
		},
		{
			name: "cluster config with networking",
			input: map[string]any{
				"kind":        testClusterKind,
				"api_version": testAPIVersion,
				"networking": []any{
					map[string]any{
						"api_server_address": testAPIServerAddress,
						"api_server_port":    testAPIServerPort,
					},
				},
			},
			validator: func(t *testing.T, result *v1alpha4.Cluster) {
				t.Helper()
				assert.Equal(
					t,
					testAPIServerAddress,
					result.Networking.APIServerAddress,
					"API server address should be set correctly",
				)
				assert.Equal(
					t,
					int32(testAPIServerPort),
					result.Networking.APIServerPort,
					"API server port should be set correctly",
				)
			},
		},
		{
			name: "cluster config with containerd patches",
			input: map[string]any{
				"kind":        testClusterKind,
				"api_version": testAPIVersion,
				"containerd_config_patches": []any{
					"[plugins.cri]\n  sandbox_image = \"test\"",
					"[plugins.cri.registry]\n  config_path = \"/etc/containerd/certs.d\"",
				},
			},
			validator: func(t *testing.T, result *v1alpha4.Cluster) {
				t.Helper()
				assert.Len(
					t,
					result.ContainerdConfigPatches,
					2,
					"should have 2 containerd config patches",
				)
				assert.Contains(
					t,
					result.ContainerdConfigPatches[0],
					"sandbox_image",
					"first patch should contain sandbox_image",
				)
				assert.Contains(
					t,
					result.ContainerdConfigPatches[1],
					"config_path",
					"second patch should contain config_path",
				)
			},
		},
		{
			name: "cluster config with runtime config",
			input: map[string]any{
				"kind":        testClusterKind,
				"api_version": testAPIVersion,
				"runtime_config": map[string]any{
					"api_alpha": "false",
					"api_beta":  "true",
				},
			},
			validator: func(t *testing.T, result *v1alpha4.Cluster) {
				t.Helper()
				assert.Len(t, result.RuntimeConfig, 2, "should have 2 runtime config entries")
				assert.Equal(
					t,
					"false",
					result.RuntimeConfig["api/alpha"],
					"api/alpha should be false",
				)
				assert.Equal(t, "true", result.RuntimeConfig["api/beta"], "api/beta should be true")
			},
		},
		{
			name: "cluster config with feature gates",
			input: map[string]any{
				"kind":        testClusterKind,
				"api_version": testAPIVersion,
				"feature_gates": map[string]any{
					"FeatureA": "true",
					"FeatureB": "false",
					"FeatureC": "True",
				},
			},
			validator: func(t *testing.T, result *v1alpha4.Cluster) {
				t.Helper()
				assert.Len(t, result.FeatureGates, 3, "should have 3 feature gates")
				assert.True(t, result.FeatureGates["FeatureA"], "FeatureA should be true")
				assert.False(t, result.FeatureGates["FeatureB"], "FeatureB should be false")
				assert.True(t, result.FeatureGates["FeatureC"], "FeatureC should be true")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := flattenKindConfig(tt.input)
			require.NoError(t, err, "flattenKindConfig should not return an error")
			require.NotNil(t, result, "flattenKindConfig should return a non-nil result")
			tt.validator(t, result)
		})
	}
}

func TestFlattenKindConfigNodes(t *testing.T) {
	tests := []struct {
		input     map[string]any
		validator func(t *testing.T, result v1alpha4.Node)
		name      string
	}{
		{
			name:  "control-plane node",
			input: map[string]any{"role": testControlPlaneRole},
			validator: func(t *testing.T, result v1alpha4.Node) {
				t.Helper()
				assert.Equal(
					t,
					v1alpha4.ControlPlaneRole,
					result.Role,
					"role should be control-plane",
				)
			},
		},
		{
			name:  "worker node",
			input: map[string]any{"role": testWorkerRole},
			validator: func(t *testing.T, result v1alpha4.Node) {
				t.Helper()
				assert.Equal(t, v1alpha4.WorkerRole, result.Role, "role should be worker")
			},
		},
		{
			name: "node with custom image",
			input: map[string]any{
				"role":  testControlPlaneRole,
				"image": testNodeImage,
			},
			validator: func(t *testing.T, result v1alpha4.Node) {
				t.Helper()
				assert.Equal(t, testNodeImage, result.Image, "image should be set correctly")
			},
		},
		{
			name: "node with custom labels",
			input: map[string]any{
				"role": testWorkerRole,
				"labels": map[string]any{
					"app":  "test",
					"tier": "backend",
				},
			},
			validator: func(t *testing.T, result v1alpha4.Node) {
				t.Helper()
				assert.Len(t, result.Labels, 2, "should have 2 labels")
				assert.Equal(t, "test", result.Labels["app"], "app label should be test")
				assert.Equal(t, "backend", result.Labels["tier"], "tier label should be backend")
			},
		},
		{
			name: "node with extra mounts",
			input: map[string]any{
				"role": testControlPlaneRole,
				"extra_mounts": []any{
					map[string]any{
						"host_path":      testHostPath,
						"container_path": testContainerPath,
					},
				},
			},
			validator: func(t *testing.T, result v1alpha4.Node) {
				t.Helper()
				assert.Len(t, result.ExtraMounts, 1, "should have 1 extra mount")
				assert.Equal(
					t,
					testHostPath,
					result.ExtraMounts[0].HostPath,
					"host path should be set correctly",
				)
				assert.Equal(
					t,
					testContainerPath,
					result.ExtraMounts[0].ContainerPath,
					"container path should be set correctly",
				)
			},
		},
		{
			name: "node with extra port mappings",
			input: map[string]any{
				"role": testControlPlaneRole,
				"extra_port_mappings": []any{
					map[string]any{
						"container_port": testContainerPort,
						"host_port":      testHostPort,
					},
				},
			},
			validator: func(t *testing.T, result v1alpha4.Node) {
				t.Helper()
				assert.Len(t, result.ExtraPortMappings, 1, "should have 1 extra port mapping")
				assert.Equal(
					t,
					int32(testContainerPort),
					result.ExtraPortMappings[0].ContainerPort,
					"container port should be set correctly",
				)
				assert.Equal(
					t,
					int32(testHostPort),
					result.ExtraPortMappings[0].HostPort,
					"host port should be set correctly",
				)
			},
		},
		{
			name: "node with kubeadm config patches",
			input: map[string]any{
				"role": testControlPlaneRole,
				"kubeadm_config_patches": []any{
					"patch1",
					"patch2",
				},
			},
			validator: func(t *testing.T, result v1alpha4.Node) {
				t.Helper()
				assert.Len(
					t,
					result.KubeadmConfigPatches,
					2,
					"should have 2 kubeadm config patches",
				)
				assert.Equal(
					t,
					"patch1",
					result.KubeadmConfigPatches[0],
					"first patch should be patch1",
				)
				assert.Equal(
					t,
					"patch2",
					result.KubeadmConfigPatches[1],
					"second patch should be patch2",
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := flattenKindConfigNodes(tt.input)
			require.NoError(t, err, "flattenKindConfigNodes should not return an error")
			tt.validator(t, result)
		})
	}
}

func TestFlattenKindConfigNetworking(t *testing.T) {
	tests := []struct {
		input     map[string]any
		validator func(t *testing.T, result v1alpha4.Networking)
		name      string
	}{
		{
			name: "networking with API server settings",
			input: map[string]any{
				"api_server_address": testAPIServerAddress,
				"api_server_port":    testAPIServerPort,
			},
			validator: func(t *testing.T, result v1alpha4.Networking) {
				t.Helper()
				assert.Equal(
					t,
					testAPIServerAddress,
					result.APIServerAddress,
					"API server address should be set correctly",
				)
				assert.Equal(
					t,
					int32(testAPIServerPort),
					result.APIServerPort,
					"API server port should be set correctly",
				)
			},
		},
		{
			name:  "networking with IPv4 family",
			input: map[string]any{"ip_family": "ipv4"},
			validator: func(t *testing.T, result v1alpha4.Networking) {
				t.Helper()
				assert.Equal(t, v1alpha4.IPv4Family, result.IPFamily, "IP family should be ipv4")
			},
		},
		{
			name:  "networking with IPv6 family",
			input: map[string]any{"ip_family": "ipv6"},
			validator: func(t *testing.T, result v1alpha4.Networking) {
				t.Helper()
				assert.Equal(t, v1alpha4.IPv6Family, result.IPFamily, "IP family should be ipv6")
			},
		},
		{
			name:  "networking with dual stack family",
			input: map[string]any{"ip_family": "dual"},
			validator: func(t *testing.T, result v1alpha4.Networking) {
				t.Helper()
				assert.Equal(
					t,
					v1alpha4.DualStackFamily,
					result.IPFamily,
					"IP family should be dual stack",
				)
			},
		},
		{
			name:  "networking with kube proxy mode",
			input: map[string]any{"kube_proxy_mode": "iptables"},
			validator: func(t *testing.T, result v1alpha4.Networking) {
				t.Helper()
				assert.Equal(
					t,
					v1alpha4.IPTablesProxyMode,
					result.KubeProxyMode,
					"kube proxy mode should be iptables",
				)
			},
		},
		{
			name:  "networking with kube proxy disabled",
			input: map[string]any{"kube_proxy_mode": "none"},
			validator: func(t *testing.T, result v1alpha4.Networking) {
				t.Helper()
				assert.Equal(
					t,
					v1alpha4.ProxyMode("none"),
					result.KubeProxyMode,
					"kube proxy mode should be none",
				)
			},
		},
		{
			name: "networking with subnets",
			input: map[string]any{
				"pod_subnet":     testPodSubnet,
				"service_subnet": testServiceSubnet,
			},
			validator: func(t *testing.T, result v1alpha4.Networking) {
				t.Helper()
				assert.Equal(
					t,
					testPodSubnet,
					result.PodSubnet,
					"pod subnet should be set correctly",
				)
				assert.Equal(
					t,
					testServiceSubnet,
					result.ServiceSubnet,
					"service subnet should be set correctly",
				)
			},
		},
		{
			name:  "networking with disable default CNI",
			input: map[string]any{"disable_default_cni": true},
			validator: func(t *testing.T, result v1alpha4.Networking) {
				t.Helper()
				assert.True(t, result.DisableDefaultCNI, "disable default CNI should be true")
			},
		},
		{
			name: "networking with DNS search",
			input: map[string]any{
				"dns_search": []any{"example.com", "test.local"},
			},
			validator: func(t *testing.T, result v1alpha4.Networking) {
				t.Helper()
				require.NotNil(t, result.DNSSearch, "DNS search should not be nil")
				assert.Len(t, *result.DNSSearch, 2, "DNS search should have 2 entries")
				assert.Equal(
					t,
					"example.com",
					(*result.DNSSearch)[0],
					"first DNS search entry should be example.com",
				)
				assert.Equal(
					t,
					"test.local",
					(*result.DNSSearch)[1],
					"second DNS search entry should be test.local",
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := flattenKindConfigNetworking(tt.input)
			require.NoError(t, err, "flattenKindConfigNetworking should not return an error")
			tt.validator(t, result)
		})
	}
}

func TestFlattenKindConfigExtraMounts(t *testing.T) {
	tests := []struct {
		input     map[string]any
		validator func(t *testing.T, result v1alpha4.Mount)
		name      string
	}{
		{
			name: "basic mount",
			input: map[string]any{
				"host_path":      testHostPath,
				"container_path": testContainerPath,
			},
			validator: func(t *testing.T, result v1alpha4.Mount) {
				t.Helper()
				assert.Equal(t, testHostPath, result.HostPath, "host path should be set correctly")
				assert.Equal(
					t,
					testContainerPath,
					result.ContainerPath,
					"container path should be set correctly",
				)
			},
		},
		{
			name: "mount with bidirectional propagation",
			input: map[string]any{
				"host_path":      testHostPath,
				"container_path": testContainerPath,
				"propagation":    "Bidirectional",
			},
			validator: func(t *testing.T, result v1alpha4.Mount) {
				t.Helper()
				assert.Equal(
					t,
					v1alpha4.MountPropagationBidirectional,
					result.Propagation,
					"propagation should be bidirectional",
				)
			},
		},
		{
			name: "mount with HostToContainer propagation",
			input: map[string]any{
				"host_path":      testHostPath,
				"container_path": testContainerPath,
				"propagation":    "HostToContainer",
			},
			validator: func(t *testing.T, result v1alpha4.Mount) {
				t.Helper()
				assert.Equal(
					t,
					v1alpha4.MountPropagationHostToContainer,
					result.Propagation,
					"propagation should be host to container",
				)
			},
		},
		{
			name: "mount with None propagation",
			input: map[string]any{
				"host_path":      testHostPath,
				"container_path": testContainerPath,
				"propagation":    "None",
			},
			validator: func(t *testing.T, result v1alpha4.Mount) {
				t.Helper()
				assert.Equal(
					t,
					v1alpha4.MountPropagationNone,
					result.Propagation,
					"propagation should be none",
				)
			},
		},
		{
			name: "mount with read only flag",
			input: map[string]any{
				"host_path":      testHostPath,
				"container_path": testContainerPath,
				"read_only":      true,
			},
			validator: func(t *testing.T, result v1alpha4.Mount) {
				t.Helper()
				assert.True(t, result.Readonly, "read only flag should be true")
			},
		},
		{
			name: "mount with selinux relabel",
			input: map[string]any{
				"host_path":       testHostPath,
				"container_path":  testContainerPath,
				"selinux_relabel": true,
			},
			validator: func(t *testing.T, result v1alpha4.Mount) {
				t.Helper()
				assert.True(t, result.SelinuxRelabel, "selinux relabel should be true")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flattenKindConfigExtraMounts(tt.input)
			tt.validator(t, result)
		})
	}
}

func TestFlattenKindConfigExtraPortMappings(t *testing.T) {
	tests := []struct {
		input     map[string]any
		validator func(t *testing.T, result v1alpha4.PortMapping)
		name      string
	}{
		{
			name: "basic port mapping",
			input: map[string]any{
				"container_port": testContainerPort,
				"host_port":      testHostPort,
			},
			validator: func(t *testing.T, result v1alpha4.PortMapping) {
				t.Helper()
				assert.Equal(
					t,
					int32(testContainerPort),
					result.ContainerPort,
					"container port should be set correctly",
				)
				assert.Equal(
					t,
					int32(testHostPort),
					result.HostPort,
					"host port should be set correctly",
				)
			},
		},
		{
			name: "port mapping with listen address",
			input: map[string]any{
				"container_port": testContainerPort,
				"host_port":      testHostPort,
				"listen_address": testListenAddress,
			},
			validator: func(t *testing.T, result v1alpha4.PortMapping) {
				t.Helper()
				assert.Equal(
					t,
					testListenAddress,
					result.ListenAddress,
					"listen address should be set correctly",
				)
			},
		},
		{
			name: "port mapping with TCP protocol",
			input: map[string]any{
				"container_port": testContainerPort,
				"host_port":      testHostPort,
				"protocol":       "TCP",
			},
			validator: func(t *testing.T, result v1alpha4.PortMapping) {
				t.Helper()
				assert.Equal(
					t,
					v1alpha4.PortMappingProtocolTCP,
					result.Protocol,
					"protocol should be TCP",
				)
			},
		},
		{
			name: "port mapping with UDP protocol",
			input: map[string]any{
				"container_port": 53,
				"host_port":      5353,
				"protocol":       "UDP",
			},
			validator: func(t *testing.T, result v1alpha4.PortMapping) {
				t.Helper()
				assert.Equal(
					t,
					v1alpha4.PortMappingProtocolUDP,
					result.Protocol,
					"protocol should be UDP",
				)
			},
		},
		{
			name: "port mapping with SCTP protocol",
			input: map[string]any{
				"container_port": 9999,
				"host_port":      9999,
				"protocol":       "SCTP",
			},
			validator: func(t *testing.T, result v1alpha4.PortMapping) {
				t.Helper()
				assert.Equal(
					t,
					v1alpha4.PortMappingProtocolSCTP,
					result.Protocol,
					"protocol should be SCTP",
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := flattenKindConfigExtraPortMappings(tt.input)
			require.NoError(t, err, "flattenKindConfigExtraPortMappings should not return an error")
			tt.validator(t, result)
		})
	}
}
