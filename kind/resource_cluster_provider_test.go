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
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewKindProvider(t *testing.T) {
	tests := []struct {
		name        string
		provider    string
		errContains string
		wantErr     bool
	}{
		{
			name:     "empty_auto_detects",
			provider: "",
		},
		{
			name:     "docker",
			provider: "docker",
		},
		{
			name:     "podman",
			provider: "podman",
		},
		{
			name:     "nerdctl",
			provider: "nerdctl",
		},
		{
			name:        "unsupported_rejects",
			provider:    "containerd",
			wantErr:     true,
			errContains: "unsupported provider",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := newKindProvider(tt.provider)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, p)

				return
			}

			require.NoError(t, err)
			assert.NotNil(t, p)
		})
	}
}

func TestClusterResource_Schema_Golden(t *testing.T) {
	r := &ClusterResource{}

	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	r.Schema(t.Context(), req, resp)

	require.False(t, resp.Diagnostics.HasError(), "schema should not have diagnostics errors")

	// Extract attribute names and properties for golden file comparison
	type attrInfo struct {
		Description string `json:"description"`
		Required    bool   `json:"required,omitempty"`
		Optional    bool   `json:"optional,omitempty"`
		Computed    bool   `json:"computed,omitempty"`
		Sensitive   bool   `json:"sensitive,omitempty"`
	}

	attrs := make(map[string]attrInfo)

	for name, attr := range resp.Schema.Attributes {
		attrs[name] = attrInfo{
			Required:    attr.IsRequired(),
			Optional:    attr.IsOptional(),
			Computed:    attr.IsComputed(),
			Sensitive:   attr.IsSensitive(),
			Description: attr.GetDescription(),
		}
	}

	data, err := json.MarshalIndent(attrs, "", "  ")
	require.NoError(t, err)

	g := goldie.New(t, goldie.WithFixtureDir(".goldie"))
	g.Assert(t, "cluster_schema", data)
}

func TestClusterResource_Schema_HasRuntimeAttribute(t *testing.T) {
	r := &ClusterResource{}
	resp := &resource.SchemaResponse{}

	r.Schema(t.Context(), resource.SchemaRequest{}, resp)

	attr, ok := resp.Schema.Attributes["runtime"]
	require.True(t, ok, "schema must have 'runtime' attribute")
	assert.True(t, attr.IsOptional(), "runtime should be optional")
	assert.False(t, attr.IsRequired(), "runtime should not be required")
	assert.False(t, attr.IsComputed(), "runtime should not be computed")
}

func TestProviderConstants(t *testing.T) {
	assert.Equal(t, "docker", providerDocker)
	assert.Equal(t, "podman", providerPodman)
	assert.Equal(t, "nerdctl", providerNerdctl)
}
