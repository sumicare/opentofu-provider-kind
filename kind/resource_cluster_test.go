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
	"fmt"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/kind/pkg/apis/config/defaults"
	"sigs.k8s.io/kind/pkg/cluster"
)

// testResourceName is the Terraform resource name used in acceptance tests.
const testResourceName = "kind_cluster.test"

func TestAccKindCluster_Basic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping acceptance test in short mode")
	}

	clusterName := acctest.RandomWithPrefix("tf-acc-cluster-test")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckKindClusterResourceDestroy(clusterName),
		Steps: []resource.TestStep{
			{
				Config: renderClusterConfig(ClusterConfig{
					Name:           clusterName,
					NodeImage:      defaults.Image,
					WaitForReady:   true,
					KubeconfigPath: "/tmp/kind-provider-test/new_file",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterCreate(),
					checkResourceAttr("name", clusterName),
					checkResourceAttr("node_image", defaults.Image),
					checkResourceAttr("wait_for_ready", "true"),
					checkResourceAttr("kubeconfig_path", "/tmp/kind-provider-test/new_file"),
				),
			},
		},
	})
}

func TestAccKindCluster_ConfigBase(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping acceptance test in short mode")
	}

	clusterName := acctest.RandomWithPrefix("tf-acc-config-base-test")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckKindClusterResourceDestroy(clusterName),
		Steps: []resource.TestStep{
			{
				Config: renderClusterConfig(ClusterConfig{
					Name:         clusterName,
					NodeImage:    defaults.Image,
					WaitForReady: true,
					KindConfig: &KindConfig{
						Networking: &Networking{
							APIServerAddress: "127.0.0.1",
							APIServerPort:    6443,
							KubeProxyMode:    "none",
						},
						RuntimeConfig: map[string]string{"api_alpha": "false"},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterCreate(),
					checkResourceAttr("kind_config.#", "1"),
					checkResourceAttr("kind_config.0.kind", "Cluster"),
					checkResourceAttr("kind_config.0.api_version", "kind.x-k8s.io/v1alpha4"),
					checkResourceAttr("wait_for_ready", "true"),
					checkResourceAttr("node_image", defaults.Image),
					checkResourceAttr("kind_config.0.networking.api_server_address", "127.0.0.1"),
					checkResourceAttr("kind_config.0.networking.api_server_port", "6443"),
					checkResourceAttr("kind_config.0.networking.kube_proxy_mode", "none"),
					checkResourceAttr("kind_config.0.runtime_config.%", "1"),
					checkResourceAttr("kind_config.0.runtime_config.api_alpha", "false"),
				),
			},
		},
	})
}

func TestAccKindCluster_ConfigNodes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping acceptance test in short mode")
	}

	clusterName := acctest.RandomWithPrefix("tf-acc-config-nodes-test")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckKindClusterResourceDestroy(clusterName),
		Steps: []resource.TestStep{
			{
				Config: renderClusterConfig(ClusterConfig{
					Name:         clusterName,
					NodeImage:    defaults.Image,
					WaitForReady: true,
					KindConfig: &KindConfig{
						Nodes: []Node{
							{Role: "control-plane", Labels: map[string]string{"name": "node0"}},
							{Role: "worker", Image: defaultNodeImage},
							{Role: "worker"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterCreate(),
					checkResourceAttr("kind_config.0.node.#", "3"),
					checkResourceAttr("kind_config.0.node.0.role", "control-plane"),
					checkResourceAttr("kind_config.0.node.0.labels.name", "node0"),
					checkResourceAttr("kind_config.0.node.1.role", "worker"),
					checkResourceAttr("kind_config.0.node.1.image", defaultNodeImage),
					checkResourceAttr("kind_config.0.node.2.role", "worker"),
					checkResourceAttr("wait_for_ready", "true"),
					checkResourceAttr("node_image", defaults.Image),
				),
			},
		},
	})
}

func TestAccKindCluster_ContainerdPatches(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping acceptance test in short mode")
	}

	clusterName := acctest.RandomWithPrefix("tf-acc-containerd-test")

	patch := `[plugins."io.containerd.grpc.v1.cri".registry.mirrors."localhost:5000"]
  endpoint = ["http://kind-registry:5000"]`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckKindClusterResourceDestroy(clusterName),
		Steps: []resource.TestStep{
			{
				Config: renderClusterConfig(ClusterConfig{
					Name:         clusterName,
					WaitForReady: true,
					KindConfig: &KindConfig{
						ContainerdConfigPatches: []string{patch},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterCreate(),
					checkResourceAttr("kind_config.0.containerd_config_patches.#", "1"),
				),
			},
		},
	})
}

// testAccCheckKindClusterResourceDestroy verifies the kind cluster
// has been destroyed.
func testAccCheckKindClusterResourceDestroy(clusterName string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		prov := cluster.NewProvider()

		list, err := prov.List()
		if err != nil {
			return fmt.Errorf("failed to list clusters: %w", err)
		}

		if slices.Contains(list, clusterName) {
			return fmt.Errorf("cluster %s should have been removed", clusterName)
		}

		// Verify kubeconfig context has been removed
		contextName := "kind-" + clusterName
		configAccess := clientcmd.NewDefaultPathOptions()

		config, err := configAccess.GetStartingConfig()
		if err == nil {
			if _, exists := config.Contexts[contextName]; exists {
				return fmt.Errorf("kubeconfig context %s should have been removed", contextName)
			}

			if _, exists := config.AuthInfos[contextName]; exists {
				return fmt.Errorf("kubeconfig user %s should have been removed", contextName)
			}

			if _, exists := config.Clusters[contextName]; exists {
				return fmt.Errorf("kubeconfig cluster %s should have been removed", contextName)
			}
		}

		return nil
	}
}

// testAccCheckClusterCreate verifies that a cluster resource exists in the state.
func testAccCheckClusterCreate() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[testResourceName]
		if !ok {
			return fmt.Errorf("root module should have resource %s", testResourceName)
		}

		return nil
	}
}

// checkResourceAttr verifies that a resource attribute has the expected value.
func checkResourceAttr(key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[testResourceName]
		if !ok {
			return fmt.Errorf("resource %s should exist", testResourceName)
		}

		if rs.Primary.Attributes[key] != value {
			return fmt.Errorf(
				"attribute %s should equal %s, got %s",
				key,
				value,
				rs.Primary.Attributes[key],
			)
		}

		return nil
	}
}

func TestNewClusterResource(t *testing.T) {
	clusterResource := NewClusterResource()
	assert.NotNil(t, clusterResource, "NewClusterResource should return a non-nil resource")
}

func TestClusterResource(t *testing.T) {
	clusterResource := &ClusterResource{}
	assert.NotNil(t, clusterResource, "ClusterResource should be instantiable")
}
