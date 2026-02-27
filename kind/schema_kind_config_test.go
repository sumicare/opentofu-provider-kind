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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// Test data constants for schema validation.
const (
	// Schema block names.
	kindConfigBlockName = "kind_config"
	nodeBlockName       = "node"
	networkingBlockName = "networking"

	// Schema field names.
	kindFieldName                    = "kind"
	apiVersionFieldName              = "api_version"
	containerdConfigPatchesFieldName = "containerd_config_patches"
	runtimeConfigFieldName           = "runtime_config"
	featureGatesFieldName            = "feature_gates"
)

func TestKindConfigBlocks(t *testing.T) {
	tests := []struct {
		name        string
		expectedKey string
		description string
	}{
		{
			name:        "has kind_config block",
			expectedKey: kindConfigBlockName,
			description: "blocks should have kind_config key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blocks := kindConfigBlocks()
			assert.NotNil(t, blocks, "blocks should not be nil")
			assert.Contains(t, blocks, tt.expectedKey, tt.description)
		})
	}
}

func TestKindConfigFieldsFramework(t *testing.T) {
	tests := []struct {
		name        string
		expectedKey string
		description string
	}{
		{
			name:        "has kind field",
			expectedKey: kindFieldName,
			description: "fields should have kind key",
		},
		{
			name:        "has api_version field",
			expectedKey: apiVersionFieldName,
			description: "fields should have api_version key",
		},
		{
			name:        "has containerd_config_patches field",
			expectedKey: containerdConfigPatchesFieldName,
			description: "fields should have containerd_config_patches key",
		},
		{
			name:        "has runtime_config field",
			expectedKey: runtimeConfigFieldName,
			description: "fields should have runtime_config key",
		},
		{
			name:        "has feature_gates field",
			expectedKey: featureGatesFieldName,
			description: "fields should have feature_gates key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := kindConfigFieldsFramework()
			assert.NotNil(t, fields, "fields should not be nil")
			assert.Contains(t, fields, tt.expectedKey, tt.description)
		})
	}
}

func TestKindConfigNestedBlocks(t *testing.T) {
	tests := []struct {
		validate    func(t *testing.T, block *schema.Resource)
		name        string
		expectedKey string
		description string
	}{
		{
			name:        "has node block",
			expectedKey: nodeBlockName,
			description: "blocks should have node key",
		},
		{
			name:        "has networking block",
			expectedKey: networkingBlockName,
			description: "blocks should have networking key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blocks := kindConfigNestedBlocks()
			assert.NotNil(t, blocks, "blocks should not be nil")
			assert.Contains(t, blocks, tt.expectedKey, tt.description)

			// Validate individual block schema
			block := blocks[tt.expectedKey]
			assert.NotNil(t, block, "%s block should not be nil", tt.expectedKey)
		})
	}
}
