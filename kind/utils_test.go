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
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test data constants for reusability across tests.
const (
	testKey         = "k"
	testStringValue = "test"
	testIntValue    = 42
	testFloatValue  = 3.14
)

var emptyMap = make(map[string]any)

// mustParseBigFloat parses a string to big.Float or panics on error.
//
//nolint:revive // Test helper function
func mustParseBigFloat(s string) *big.Float {
	f, _, err := big.ParseFloat(s, 10, 64, big.ToNearestEven)
	if err != nil {
		panic(err)
	}

	return f
}

func TestGetStringSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]any
		key      string
		expected []string
		isNil    bool
	}{
		{
			name:     "extracts all strings",
			input:    map[string]any{testKey: []any{"a", "b", "c"}},
			key:      testKey,
			expected: []string{"a", "b", "c"},
			isNil:    false,
		},
		{
			name:     "filters out non-strings",
			input:    map[string]any{testKey: []any{"a", 123, "b"}},
			key:      testKey,
			expected: []string{"a", "b"},
			isNil:    false,
		},
		{
			name:     "handles empty slice",
			input:    map[string]any{testKey: make([]any, 0)},
			key:      testKey,
			expected: []string{},
			isNil:    false,
		},
		{
			name:     "missing key returns nil",
			input:    emptyMap,
			key:      testKey,
			expected: nil,
			isNil:    true,
		},
		{
			name:     "nil value returns nil",
			input:    map[string]any{testKey: nil},
			key:      testKey,
			expected: nil,
			isNil:    true,
		},
		{
			name:     "wrong type returns nil",
			input:    map[string]any{testKey: testStringValue},
			key:      testKey,
			expected: nil,
			isNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getStringSlice(tt.input, tt.key)
			if tt.isNil {
				assert.Nil(t, result, "getStringSlice should return nil for invalid inputs")
			} else {
				assert.Equal(
					t,
					tt.expected,
					result,
					"getStringSlice should handle all input types correctly",
				)
			}
		})
	}
}

func TestGetMapSlice(t *testing.T) {
	tests := []struct {
		name        string
		input       map[string]any
		key         string
		expectedLen int
		isNil       bool
	}{
		{
			name: "extracts all maps",
			input: map[string]any{
				testKey: []any{map[string]any{"a": 1}, map[string]any{"b": 2}},
			},
			key:         testKey,
			expectedLen: 2,
			isNil:       false,
		},
		{
			name: "filters out non-maps",
			input: map[string]any{
				testKey: []any{map[string]any{"a": 1}, testStringValue, map[string]any{"b": 2}},
			},
			key:         testKey,
			expectedLen: 2,
			isNil:       false,
		},
		{
			name:        "handles empty slice",
			input:       map[string]any{testKey: make([]any, 0)},
			key:         testKey,
			expectedLen: 0,
			isNil:       false,
		},
		{
			name:        "missing key returns nil",
			input:       emptyMap,
			key:         testKey,
			expectedLen: 0,
			isNil:       true,
		},
		{
			name:        "wrong type returns nil",
			input:       map[string]any{testKey: testStringValue},
			key:         testKey,
			expectedLen: 0,
			isNil:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getMapSlice(tt.input, tt.key)
			if tt.isNil {
				assert.Nil(t, result, "getMapSlice should return nil for invalid inputs")
			} else {
				assert.Len(
					t,
					result,
					tt.expectedLen,
					"getMapSlice should return correct slice length",
				)
			}
		})
	}
}

func TestGetStringMap(t *testing.T) {
	tests := []struct {
		input    map[string]any
		expected map[string]string
		name     string
		key      string
		isNil    bool
	}{
		{
			name:     "extracts all string values",
			input:    map[string]any{testKey: map[string]any{"a": "v1", "b": "v2"}},
			key:      testKey,
			expected: map[string]string{"a": "v1", "b": "v2"},
			isNil:    false,
		},
		{
			name:     "filters out non-strings",
			input:    map[string]any{testKey: map[string]any{"a": "v1", "b": 123}},
			key:      testKey,
			expected: map[string]string{"a": "v1"},
			isNil:    false,
		},
		{
			name:     "handles empty map",
			input:    map[string]any{testKey: make(map[string]any)},
			key:      testKey,
			expected: map[string]string{},
			isNil:    false,
		},
		{
			name:     "missing key returns nil",
			input:    emptyMap,
			key:      testKey,
			expected: nil,
			isNil:    true,
		},
		{
			name:     "nil value returns nil",
			input:    map[string]any{testKey: nil},
			key:      testKey,
			expected: nil,
			isNil:    true,
		},
		{
			name:     "wrong type returns nil",
			input:    map[string]any{testKey: testStringValue},
			key:      testKey,
			expected: nil,
			isNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getStringMap(tt.input, tt.key)
			if tt.isNil {
				assert.Nil(t, result, "getStringMap should return nil for invalid inputs")
			} else {
				assert.Equal(
					t,
					tt.expected,
					result,
					"getStringMap should handle all input types correctly",
				)
			}
		})
	}
}

func TestNormalizeToml(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		expected  string
		expectErr bool
		contains  bool
	}{
		{
			name:      "handles empty string",
			input:     "",
			expected:  "",
			expectErr: false,
			contains:  false,
		},
		{
			name:      "handles nil input",
			input:     nil,
			expected:  "",
			expectErr: false,
			contains:  false,
		},
		{
			name:      "handles non-string input",
			input:     123,
			expected:  "",
			expectErr: false,
			contains:  false,
		},
		{
			name:      "parses valid TOML",
			input:     `title = "Test"`,
			expected:  "title",
			expectErr: false,
			contains:  true,
		},
		{
			name:      "returns error for invalid TOML",
			input:     `invalid [[[`,
			expected:  `invalid [[[`,
			expectErr: true,
			contains:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := normalizeToml(tt.input)
			if tt.expectErr {
				require.Error(t, err, "normalizeToml should return error for invalid input")
				assert.Equal(t, tt.expected, result, "normalizeToml should return input on error")
			} else {
				assert.NoError(t, err, "normalizeToml should not return error for valid input")

				if tt.contains {
					assert.Contains(
						t,
						result,
						tt.expected,
						"normalizeToml should contain expected substring",
					)
				} else {
					assert.Equal(
						t,
						tt.expected,
						result,
						"normalizeToml should return expected result",
					)
				}
			}
		})
	}
}

func TestObjectToMap(t *testing.T) {
	tests := []struct {
		expected map[string]any
		input    types.Object
		name     string
		isNil    bool
	}{
		{
			name:     "null object returns nil",
			input:    types.ObjectNull(make(map[string]attr.Type)),
			expected: nil,
			isNil:    true,
		},
		{
			name:     "unknown object returns nil",
			input:    types.ObjectUnknown(map[string]attr.Type{testKey: types.StringType}),
			expected: nil,
			isNil:    true,
		},
		{
			name: "empty object returns empty map",
			input: types.ObjectValueMust(
				make(map[string]attr.Type),
				make(map[string]attr.Value),
			),
			expected: map[string]any{},
			isNil:    false,
		},
		{
			name: "object with fields extracts correctly",
			input: types.ObjectValueMust(
				map[string]attr.Type{"name": types.StringType, "age": types.Int64Type},
				map[string]attr.Value{
					"name": types.StringValue(testStringValue),
					"age":  types.Int64Value(int64(testIntValue)),
				},
			),
			expected: map[string]any{"name": testStringValue, "age": testIntValue},
			isNil:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := objectToMap(tt.input)
			if tt.isNil {
				assert.Nil(t, result, "objectToMap should return nil for null/unknown objects")
			} else {
				assert.Equal(
					t,
					tt.expected,
					result,
					"objectToMap should handle all object types correctly",
				)
			}
		})
	}
}

func TestListToSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    types.List
		expected []any
		isNil    bool
	}{
		{
			name:     "null list returns nil",
			input:    types.ListNull(types.StringType),
			expected: nil,
			isNil:    true,
		},
		{
			name:     "unknown list returns nil",
			input:    types.ListUnknown(types.StringType),
			expected: nil,
			isNil:    true,
		},
		{
			name:     "empty list returns empty slice",
			input:    types.ListValueMust(types.StringType, make([]attr.Value, 0)),
			expected: []any{},
			isNil:    false,
		},
		{
			name: "list with values extracts correctly",
			input: types.ListValueMust(
				types.StringType,
				[]attr.Value{types.StringValue("a"), types.StringValue("b")},
			),
			expected: []any{"a", "b"},
			isNil:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := listToSlice(tt.input)
			if tt.isNil {
				assert.Nil(t, result, "listToSlice should return nil for null/unknown lists")
			} else {
				assert.Equal(
					t,
					tt.expected,
					result,
					"listToSlice should handle all list types correctly",
				)
			}
		})
	}
}

func TestSetToSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    types.Set
		expected []any
		isNil    bool
	}{
		{
			name:     "null set returns nil",
			input:    types.SetNull(types.StringType),
			expected: nil,
			isNil:    true,
		},
		{
			name:     "unknown set returns nil",
			input:    types.SetUnknown(types.StringType),
			expected: nil,
			isNil:    true,
		},
		{
			name:     "empty set returns empty slice",
			input:    types.SetValueMust(types.StringType, make([]attr.Value, 0)),
			expected: []any{},
			isNil:    false,
		},
		{
			name:     "set with values extracts correctly",
			input:    types.SetValueMust(types.StringType, []attr.Value{types.StringValue("x")}),
			expected: []any{"x"},
			isNil:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := setToSlice(tt.input)
			if tt.isNil {
				assert.Nil(t, result, "setToSlice should return nil for null/unknown sets")
			} else {
				assert.Equal(
					t,
					tt.expected,
					result,
					"setToSlice should handle all set types correctly",
				)
			}
		})
	}
}

func TestMapToMap(t *testing.T) {
	tests := []struct {
		expected map[string]any
		input    types.Map
		name     string
		isNil    bool
	}{
		{
			name:     "null map returns nil",
			input:    types.MapNull(types.StringType),
			expected: nil,
			isNil:    true,
		},
		{
			name:     "unknown map returns nil",
			input:    types.MapUnknown(types.StringType),
			expected: nil,
			isNil:    true,
		},
		{
			name:     "empty map returns empty map",
			input:    types.MapValueMust(types.StringType, make(map[string]attr.Value)),
			expected: map[string]any{},
			isNil:    false,
		},
		{
			name: "map with values extracts correctly",
			input: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{"k1": types.StringValue("v1"), "k2": types.StringValue("v2")},
			),
			expected: map[string]any{"k1": "v1", "k2": "v2"},
			isNil:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapToMap(tt.input)
			if tt.isNil {
				assert.Nil(t, result, "mapToMap should return nil for null/unknown maps")
			} else {
				assert.Equal(
					t,
					tt.expected,
					result,
					"mapToMap should handle all map types correctly",
				)
			}
		})
	}
}

func TestAttrValueToAny(t *testing.T) {
	tests := []struct {
		input    attr.Value
		expected any
		validate func(t *testing.T, result any)
		name     string
	}{
		{
			name:     "string value converts correctly",
			input:    types.StringValue(testStringValue),
			expected: testStringValue,
			validate: func(t *testing.T, result any) {
				t.Helper()
				assert.Equal(t, testStringValue, result)
			},
		},
		{
			name:     "bool value converts correctly",
			input:    types.BoolValue(true),
			expected: true,
			validate: func(t *testing.T, result any) {
				t.Helper()
				assert.Equal(t, true, result)
			},
		},
		{
			name:     "int64 value converts correctly",
			input:    types.Int64Value(int64(testIntValue)),
			expected: testIntValue,
			validate: func(t *testing.T, result any) {
				t.Helper()
				assert.Equal(t, testIntValue, result)
			},
		},
		{
			name:     "float64 value converts correctly",
			input:    types.Float64Value(testFloatValue),
			expected: testFloatValue,
			validate: func(t *testing.T, result any) {
				t.Helper()
				assert.InDelta(t, testFloatValue, result, 0.01)
			},
		},
		{
			name:     "number value converts correctly",
			input:    types.NumberValue(mustParseBigFloat("42.5")),
			expected: 42.5,
			validate: func(t *testing.T, result any) {
				t.Helper()
				assert.InDelta(t, 42.5, result, 0.01)
			},
		},
		{
			name:     "null string returns nil",
			input:    types.StringNull(),
			expected: nil,
			validate: func(t *testing.T, result any) {
				t.Helper()
				assert.Nil(t, result)
			},
		},
		{
			name:     "unknown string returns nil",
			input:    types.StringUnknown(),
			expected: nil,
			validate: func(t *testing.T, result any) {
				t.Helper()
				assert.Nil(t, result)
			},
		},
		{
			name: "list converts to slice",
			input: types.ListValueMust(
				types.StringType,
				[]attr.Value{types.StringValue("a"), types.StringValue("b")},
			),
			validate: func(t *testing.T, result any) {
				t.Helper()

				slice, ok := result.([]any)
				require.True(t, ok, "result should be a slice")
				assert.Len(t, slice, 2)
			},
		},
		{
			name:  "set converts to slice",
			input: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("x")}),
			validate: func(t *testing.T, result any) {
				t.Helper()

				slice, ok := result.([]any)
				require.True(t, ok, "result should be a slice")
				assert.Len(t, slice, 1)
			},
		},
		{
			name: "map converts to map[string]any",
			input: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{testKey: types.StringValue("v")},
			),
			validate: func(t *testing.T, result any) {
				t.Helper()

				m, ok := result.(map[string]any)
				require.True(t, ok, "result should be a map[string]any")
				assert.Equal(t, "v", m[testKey])
			},
		},
		{
			name: "object converts to map[string]any",
			input: types.ObjectValueMust(
				map[string]attr.Type{"n": types.StringType},
				map[string]attr.Value{"n": types.StringValue("t")},
			),
			validate: func(t *testing.T, result any) {
				t.Helper()

				m, ok := result.(map[string]any)
				require.True(t, ok, "result should be a map[string]any")
				assert.Equal(t, "t", m["n"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := attrValueToAny(tt.input)
			tt.validate(t, result)
		})
	}
}

func TestParseKindConfigFromFramework(t *testing.T) {
	t.Run("handles null and empty lists correctly", func(t *testing.T) {
		ctx := t.Context()

		// Test null list
		result, err := parseKindConfigFromFramework(ctx, types.ListNull(types.ObjectType{}))
		require.NoError(t, err, "should handle null list without error")
		assert.Nil(t, result, "should return nil for null list")

		// Test empty list
		result, err = parseKindConfigFromFramework(
			ctx,
			types.ListValueMust(types.ObjectType{}, make([]attr.Value, 0)),
		)
		require.NoError(t, err, "should handle empty list without error")
		assert.Nil(t, result, "should return nil for empty list")
	})

	t.Run("parses valid kind configuration correctly", func(t *testing.T) {
		ctx := t.Context()

		// Create object type and value for Kind configuration
		objType := map[string]attr.Type{
			"kind":        types.StringType,
			"api_version": types.StringType,
		}
		obj := types.ObjectValueMust(objType, map[string]attr.Value{
			"kind":        types.StringValue("Cluster"),
			"api_version": types.StringValue("kind.x-k8s.io/v1alpha4"),
		})
		list := types.ListValueMust(types.ObjectType{AttrTypes: objType}, []attr.Value{obj})

		result, err := parseKindConfigFromFramework(ctx, list)
		require.NoError(t, err, "should parse valid configuration without error")
		require.NotNil(t, result, "should return non-nil result for valid configuration")
		assert.Equal(t, "Cluster", result.Kind, "should extract Kind field correctly")
		assert.Equal(
			t,
			"kind.x-k8s.io/v1alpha4",
			result.APIVersion,
			"should extract APIVersion field correctly",
		)
	})
}
