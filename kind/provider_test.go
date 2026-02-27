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

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/stretchr/testify/assert"
)

const testVersion = "test"

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"kind": providerserver.NewProtocol6WithError(New("test")()),
}

func TestNew(t *testing.T) {
	tests := []struct {
		validator func(t *testing.T, provider any)
		name      string
		version   string
	}{
		{
			name:    "creates provider with test version",
			version: testVersion,
			validator: func(t *testing.T, provider any) {
				t.Helper()
				// Additional validation can be added here if needed
				assert.IsType(
					t,
					&KindProvider{},
					provider,
					"provider should be of type KindProvider",
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := New(tt.version)
			assert.NotNil(t, factory, "factory should not be nil")

			provider := factory()
			assert.NotNil(t, provider, "provider should not be nil")
			tt.validator(t, provider)
		})
	}
}

func TestKindProvider(t *testing.T) {
	tests := []struct {
		validator func(t *testing.T, provider *KindProvider)
		name      string
		version   string
	}{
		{
			name:    "has correct test version",
			version: testVersion,
			validator: func(t *testing.T, provider *KindProvider) {
				t.Helper()
				assert.Equal(
					t,
					testVersion,
					provider.version,
					"version should match expected value",
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &KindProvider{version: tt.version}
			assert.NotNil(t, provider, "provider should not be nil")
			tt.validator(t, provider)
		})
	}
}
