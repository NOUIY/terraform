// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package genconfig

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terraform/internal/addrs"
	"github.com/hashicorp/terraform/internal/configs/configschema"
	"github.com/hashicorp/terraform/internal/lang/marks"
)

func TestConfigGeneration(t *testing.T) {
	tcs := map[string]struct {
		schema   *configschema.Block
		addr     addrs.AbsResourceInstance
		provider addrs.LocalProviderConfig
		value    cty.Value
		expected string
	}{
		"simple_resource": {
			schema: &configschema.Block{
				BlockTypes: map[string]*configschema.NestedBlock{
					"list_block": {
						Block: configschema.Block{
							Attributes: map[string]*configschema.Attribute{
								"nested_value": {
									Type:     cty.String,
									Optional: true,
								},
							},
						},
						Nesting: configschema.NestingSingle,
					},
				},
				Attributes: map[string]*configschema.Attribute{
					"id": {
						Type:     cty.String,
						Computed: true,
					},
					"value": {
						Type:     cty.String,
						Optional: true,
					},
				},
			},
			addr: addrs.AbsResourceInstance{
				Module: nil,
				Resource: addrs.ResourceInstance{
					Resource: addrs.Resource{
						Mode: addrs.ManagedResourceMode,
						Type: "tfcoremock_simple_resource",
						Name: "empty",
					},
					Key: nil,
				},
			},
			provider: addrs.LocalProviderConfig{
				LocalName: "tfcoremock",
			},
			value: cty.NilVal,
			expected: `
resource "tfcoremock_simple_resource" "empty" {
  value = null          # OPTIONAL string
  list_block {          # OPTIONAL block
    nested_value = null # OPTIONAL string
  }
}`,
		},
		"simple_resource_with_state": {
			schema: &configschema.Block{
				BlockTypes: map[string]*configschema.NestedBlock{
					"list_block": {
						Block: configschema.Block{
							Attributes: map[string]*configschema.Attribute{
								"nested_value": {
									Type:     cty.String,
									Optional: true,
								},
							},
						},
						Nesting: configschema.NestingSingle,
					},
				},
				Attributes: map[string]*configschema.Attribute{
					"id": {
						Type:     cty.String,
						Computed: true,
					},
					"value": {
						Type:     cty.String,
						Optional: true,
					},
				},
			},
			addr: addrs.AbsResourceInstance{
				Module: nil,
				Resource: addrs.ResourceInstance{
					Resource: addrs.Resource{
						Mode: addrs.ManagedResourceMode,
						Type: "tfcoremock_simple_resource",
						Name: "empty",
					},
					Key: nil,
				},
			},
			provider: addrs.LocalProviderConfig{
				LocalName: "tfcoremock",
			},
			value: cty.ObjectVal(map[string]cty.Value{
				"id":    cty.StringVal("D2320658"),
				"value": cty.StringVal("Hello, world!"),
				"list_block": cty.ObjectVal(map[string]cty.Value{
					"nested_value": cty.StringVal("Hello, solar system!"),
				}),
			}),
			expected: `
resource "tfcoremock_simple_resource" "empty" {
  value = "Hello, world!"
  list_block {
    nested_value = "Hello, solar system!"
  }
}`,
		},
		"simple_resource_with_partial_state": {
			schema: &configschema.Block{
				BlockTypes: map[string]*configschema.NestedBlock{
					"list_block": {
						Block: configschema.Block{
							Attributes: map[string]*configschema.Attribute{
								"nested_value": {
									Type:     cty.String,
									Optional: true,
								},
							},
						},
						Nesting: configschema.NestingSingle,
					},
				},
				Attributes: map[string]*configschema.Attribute{
					"id": {
						Type:     cty.String,
						Computed: true,
					},
					"value": {
						Type:     cty.String,
						Optional: true,
					},
				},
			},
			addr: addrs.AbsResourceInstance{
				Module: nil,
				Resource: addrs.ResourceInstance{
					Resource: addrs.Resource{
						Mode: addrs.ManagedResourceMode,
						Type: "tfcoremock_simple_resource",
						Name: "empty",
					},
					Key: nil,
				},
			},
			provider: addrs.LocalProviderConfig{
				LocalName: "tfcoremock",
			},
			value: cty.ObjectVal(map[string]cty.Value{
				"id": cty.StringVal("D2320658"),
				"list_block": cty.ObjectVal(map[string]cty.Value{
					"nested_value": cty.StringVal("Hello, solar system!"),
				}),
			}),
			expected: `
resource "tfcoremock_simple_resource" "empty" {
  value = null
  list_block {
    nested_value = "Hello, solar system!"
  }
}`,
		},
		"simple_resource_with_alternate_provider": {
			schema: &configschema.Block{
				BlockTypes: map[string]*configschema.NestedBlock{
					"list_block": {
						Block: configschema.Block{
							Attributes: map[string]*configschema.Attribute{
								"nested_value": {
									Type:     cty.String,
									Optional: true,
								},
							},
						},
						Nesting: configschema.NestingSingle,
					},
				},
				Attributes: map[string]*configschema.Attribute{
					"id": {
						Type:     cty.String,
						Computed: true,
					},
					"value": {
						Type:     cty.String,
						Optional: true,
					},
				},
			},
			addr: addrs.AbsResourceInstance{
				Module: nil,
				Resource: addrs.ResourceInstance{
					Resource: addrs.Resource{
						Mode: addrs.ManagedResourceMode,
						Type: "tfcoremock_simple_resource",
						Name: "empty",
					},
					Key: nil,
				},
			},
			provider: addrs.LocalProviderConfig{
				LocalName: "mock",
			},
			value: cty.ObjectVal(map[string]cty.Value{
				"id":    cty.StringVal("D2320658"),
				"value": cty.StringVal("Hello, world!"),
				"list_block": cty.ObjectVal(map[string]cty.Value{
					"nested_value": cty.StringVal("Hello, solar system!"),
				}),
			}),
			expected: `
resource "tfcoremock_simple_resource" "empty" {
  provider = mock
  value    = "Hello, world!"
  list_block {
    nested_value = "Hello, solar system!"
  }
}`,
		},
		"simple_resource_with_aliased_provider": {
			schema: &configschema.Block{
				BlockTypes: map[string]*configschema.NestedBlock{
					"list_block": {
						Block: configschema.Block{
							Attributes: map[string]*configschema.Attribute{
								"nested_value": {
									Type:     cty.String,
									Optional: true,
								},
							},
						},
						Nesting: configschema.NestingSingle,
					},
				},
				Attributes: map[string]*configschema.Attribute{
					"id": {
						Type:     cty.String,
						Computed: true,
					},
					"value": {
						Type:     cty.String,
						Optional: true,
					},
				},
			},
			addr: addrs.AbsResourceInstance{
				Module: nil,
				Resource: addrs.ResourceInstance{
					Resource: addrs.Resource{
						Mode: addrs.ManagedResourceMode,
						Type: "tfcoremock_simple_resource",
						Name: "empty",
					},
					Key: nil,
				},
			},
			provider: addrs.LocalProviderConfig{
				LocalName: "tfcoremock",
				Alias:     "alternate",
			},
			value: cty.ObjectVal(map[string]cty.Value{
				"id":    cty.StringVal("D2320658"),
				"value": cty.StringVal("Hello, world!"),
				"list_block": cty.ObjectVal(map[string]cty.Value{
					"nested_value": cty.StringVal("Hello, solar system!"),
				}),
			}),
			expected: `
resource "tfcoremock_simple_resource" "empty" {
  provider = tfcoremock.alternate
  value    = "Hello, world!"
  list_block {
    nested_value = "Hello, solar system!"
  }
}`,
		},
		"resource_with_nulls": {
			schema: &configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"id": {
						Type:     cty.String,
						Computed: true,
					},
					"single": {
						NestedType: &configschema.Object{
							Attributes: map[string]*configschema.Attribute{},
							Nesting:    configschema.NestingSingle,
						},
						Required: true,
					},
					"list": {
						NestedType: &configschema.Object{
							Attributes: map[string]*configschema.Attribute{
								"nested_id": {
									Type:     cty.String,
									Optional: true,
								},
							},
							Nesting: configschema.NestingList,
						},
						Required: true,
					},
					"map": {
						NestedType: &configschema.Object{
							Attributes: map[string]*configschema.Attribute{
								"nested_id": {
									Type:     cty.String,
									Optional: true,
								},
							},
							Nesting: configschema.NestingMap,
						},
						Required: true,
					},
				},
				BlockTypes: map[string]*configschema.NestedBlock{
					"nested_single": {
						Nesting: configschema.NestingSingle,
						Block: configschema.Block{
							Attributes: map[string]*configschema.Attribute{
								"nested_id": {
									Type:     cty.String,
									Optional: true,
								},
							},
						},
					},
					// No configschema.NestingGroup example for this test, because this block type can never be null in state.
					"nested_list": {
						Nesting: configschema.NestingList,
						Block: configschema.Block{
							Attributes: map[string]*configschema.Attribute{
								"nested_id": {
									Type:     cty.String,
									Optional: true,
								},
							},
						},
					},
					"nested_set": {
						Nesting: configschema.NestingSet,
						Block: configschema.Block{
							Attributes: map[string]*configschema.Attribute{
								"nested_id": {
									Type:     cty.String,
									Optional: true,
								},
							},
						},
					},
					"nested_map": {
						Nesting: configschema.NestingMap,
						Block: configschema.Block{
							Attributes: map[string]*configschema.Attribute{
								"nested_id": {
									Type:     cty.String,
									Optional: true,
								},
							},
						},
					},
				},
			},
			addr: addrs.AbsResourceInstance{
				Module: nil,
				Resource: addrs.ResourceInstance{
					Resource: addrs.Resource{
						Mode: addrs.ManagedResourceMode,
						Type: "tfcoremock_simple_resource",
						Name: "empty",
					},
					Key: nil,
				},
			},
			provider: addrs.LocalProviderConfig{
				LocalName: "tfcoremock",
			},
			value: cty.ObjectVal(map[string]cty.Value{
				"id":     cty.StringVal("D2320658"),
				"single": cty.NullVal(cty.Object(map[string]cty.Type{})),
				"list": cty.NullVal(cty.List(cty.Object(map[string]cty.Type{
					"nested_id": cty.String,
				}))),
				"map": cty.NullVal(cty.Map(cty.Object(map[string]cty.Type{
					"nested_id": cty.String,
				}))),
				"nested_single": cty.NullVal(cty.Object(map[string]cty.Type{
					"nested_id": cty.String,
				})),
				"nested_list": cty.ListValEmpty(cty.Object(map[string]cty.Type{
					"nested_id": cty.String,
				})),
				"nested_set": cty.SetValEmpty(cty.Object(map[string]cty.Type{
					"nested_id": cty.String,
				})),
				"nested_map": cty.MapValEmpty(cty.Object(map[string]cty.Type{
					"nested_id": cty.String,
				})),
			}),
			expected: `
resource "tfcoremock_simple_resource" "empty" {
  list   = null
  map    = null
  single = null
}`,
		},
		"simple_resource_with_stringified_json_object": {
			schema: &configschema.Block{
				// BlockTypes: map[string]*configschema.NestedBlock{},
				Attributes: map[string]*configschema.Attribute{
					"id": {
						Type:     cty.String,
						Computed: true,
					},
					"value": {
						Type:     cty.String,
						Optional: true,
					},
				},
			},
			addr: addrs.AbsResourceInstance{
				Module: nil,
				Resource: addrs.ResourceInstance{
					Resource: addrs.Resource{
						Mode: addrs.ManagedResourceMode,
						Type: "tfcoremock_simple_resource",
						Name: "empty",
					},
					Key: nil,
				},
			},
			provider: addrs.LocalProviderConfig{
				LocalName: "tfcoremock",
			},
			value: cty.ObjectVal(map[string]cty.Value{
				"id":    cty.StringVal("D2320658"),
				"value": cty.StringVal(`{ "0Hello": "World", "And": ["Solar", "System"], "ready": true }`),
			}),
			expected: `
resource "tfcoremock_simple_resource" "empty" {
  value = jsonencode({
    "0Hello" = "World"
    And      = ["Solar", "System"]
    ready    = true
  })
}`,
		},
		"simple_resource_with_stringified_json_array": {
			schema: &configschema.Block{
				// BlockTypes: map[string]*configschema.NestedBlock{},
				Attributes: map[string]*configschema.Attribute{
					"id": {
						Type:     cty.String,
						Computed: true,
					},
					"value": {
						Type:     cty.String,
						Optional: true,
					},
				},
			},
			addr: addrs.AbsResourceInstance{
				Module: nil,
				Resource: addrs.ResourceInstance{
					Resource: addrs.Resource{
						Mode: addrs.ManagedResourceMode,
						Type: "tfcoremock_simple_resource",
						Name: "empty",
					},
					Key: nil,
				},
			},
			provider: addrs.LocalProviderConfig{
				LocalName: "tfcoremock",
			},
			value: cty.ObjectVal(map[string]cty.Value{
				"id":    cty.StringVal("D2320658"),
				"value": cty.StringVal(`["Hello", "World"]`),
			}),
			expected: `
resource "tfcoremock_simple_resource" "empty" {
  value = jsonencode(["Hello", "World"])
}`,
		},
		"simple_resource_with_json_primitive_strings": {
			schema: &configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"id": {
						Type:     cty.String,
						Computed: true,
					},
					"value_string_number": {
						Type:     cty.String,
						Optional: true,
					},
					"value_string_bool": {
						Type:     cty.String,
						Optional: true,
					},
				},
			},
			addr: addrs.AbsResourceInstance{
				Module: nil,
				Resource: addrs.ResourceInstance{
					Resource: addrs.Resource{
						Mode: addrs.ManagedResourceMode,
						Type: "tfcoremock_simple_resource",
						Name: "empty",
					},
					Key: nil,
				},
			},
			provider: addrs.LocalProviderConfig{
				LocalName: "tfcoremock",
			},
			value: cty.ObjectVal(map[string]cty.Value{
				"id":                  cty.StringVal("D2320658"),
				"value_string_number": cty.StringVal("42"),
				"value_string_bool":   cty.StringVal("true"),
			}),
			expected: `
resource "tfcoremock_simple_resource" "empty" {
  value_string_bool   = "true"
  value_string_number = "42"
}`,
		},
		"simple_resource_with_malformed_json": {
			schema: &configschema.Block{
				// BlockTypes: map[string]*configschema.NestedBlock{},
				Attributes: map[string]*configschema.Attribute{
					"id": {
						Type:     cty.String,
						Computed: true,
					},
					"value": {
						Type:     cty.String,
						Optional: true,
					},
				},
			},
			addr: addrs.AbsResourceInstance{
				Module: nil,
				Resource: addrs.ResourceInstance{
					Resource: addrs.Resource{
						Mode: addrs.ManagedResourceMode,
						Type: "tfcoremock_simple_resource",
						Name: "empty",
					},
					Key: nil,
				},
			},
			provider: addrs.LocalProviderConfig{
				LocalName: "tfcoremock",
			},
			value: cty.ObjectVal(map[string]cty.Value{
				"id":    cty.StringVal("D2320658"),
				"value": cty.StringVal(`["Hello", "World"`),
			}),
			expected: `
resource "tfcoremock_simple_resource" "empty" {
  value = "[\"Hello\", \"World\""
}`,
		},
		// Just try all the simple values with sensitive marks.
		"sensitive_values": {
			schema: &configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"string":       sensitiveAttribute(cty.String),
					"empty_string": sensitiveAttribute(cty.String),
					"number":       sensitiveAttribute(cty.Number),
					"bool":         sensitiveAttribute(cty.Bool),
					"object": sensitiveAttribute(cty.Object(map[string]cty.Type{
						"nested": cty.String,
					})),
					"list": sensitiveAttribute(cty.List(cty.String)),
					"map":  sensitiveAttribute(cty.Map(cty.String)),
					"set":  sensitiveAttribute(cty.Set(cty.String)),
				},
			},
			addr: addrs.AbsResourceInstance{
				Module: addrs.RootModuleInstance,
				Resource: addrs.ResourceInstance{
					Resource: addrs.Resource{
						Mode: addrs.ManagedResourceMode,
						Type: "tfcoremock_sensitive_values",
						Name: "values",
					},
					Key: addrs.NoKey,
				},
			},
			provider: addrs.LocalProviderConfig{
				LocalName: "tfcoremock",
			},
			value: cty.ObjectVal(map[string]cty.Value{
				// Values that are sensitive will now be marked as such
				"string":       cty.StringVal("Hello, world!").Mark(marks.Sensitive),
				"empty_string": cty.StringVal("").Mark(marks.Sensitive),
				"number":       cty.NumberIntVal(42).Mark(marks.Sensitive),
				"bool":         cty.True.Mark(marks.Sensitive),
				"object": cty.ObjectVal(map[string]cty.Value{
					"nested": cty.StringVal("Hello, solar system!"),
				}).Mark(marks.Sensitive),
				"list": cty.ListVal([]cty.Value{
					cty.StringVal("Hello, world!"),
				}).Mark(marks.Sensitive),
				"map": cty.MapVal(map[string]cty.Value{
					"key": cty.StringVal("Hello, world!"),
				}).Mark(marks.Sensitive),
				"set": cty.SetVal([]cty.Value{
					cty.StringVal("Hello, world!"),
				}).Mark(marks.Sensitive),
			}),
			expected: `
resource "tfcoremock_sensitive_values" "values" {
  bool         = null # sensitive
  empty_string = null # sensitive
  list         = null # sensitive
  map          = null # sensitive
  number       = null # sensitive
  object       = null # sensitive
  set          = null # sensitive
  string       = null # sensitive
}`,
		},
		"simple_map_with_whitespace_in_keys": {
			schema: &configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"map": {
						Type:     cty.Map(cty.String),
						Optional: true,
					},
				},
			},
			addr: addrs.AbsResourceInstance{
				Module: addrs.RootModuleInstance,
				Resource: addrs.ResourceInstance{
					Resource: addrs.Resource{
						Mode: addrs.ManagedResourceMode,
						Type: "testing_resource",
						Name: "resource",
					},
					Key: addrs.NoKey,
				},
			},
			provider: addrs.LocalProviderConfig{
				LocalName: "testing",
			},
			value: cty.ObjectVal(map[string]cty.Value{
				"map": cty.MapVal(map[string]cty.Value{
					"key with spaces":      cty.StringVal("spaces"),
					"key_with_underscores": cty.StringVal("underscores"),
				}),
			}),
			expected: `resource "testing_resource" "resource" {
  map = {
    "key with spaces"    = "spaces"
    key_with_underscores = "underscores"
  }
}`,
		},
		"nested_map_with_whitespace_in_keys": {
			schema: &configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"map": {
						NestedType: &configschema.Object{
							Attributes: map[string]*configschema.Attribute{
								"value": {
									Type:     cty.String,
									Optional: true,
								},
							},
							Nesting: configschema.NestingMap,
						},
						Optional: true,
					},
				},
			},
			addr: addrs.AbsResourceInstance{
				Module: addrs.RootModuleInstance,
				Resource: addrs.ResourceInstance{
					Resource: addrs.Resource{
						Mode: addrs.ManagedResourceMode,
						Type: "testing_resource",
						Name: "resource",
					},
					Key: addrs.NoKey,
				},
			},
			provider: addrs.LocalProviderConfig{
				LocalName: "testing",
			},
			value: cty.ObjectVal(map[string]cty.Value{
				"map": cty.MapVal(map[string]cty.Value{
					"key with spaces": cty.ObjectVal(map[string]cty.Value{
						"value": cty.StringVal("spaces"),
					}),
					"key_with_underscores": cty.ObjectVal(map[string]cty.Value{
						"value": cty.StringVal("underscores"),
					}),
				}),
			}),
			expected: `resource "testing_resource" "resource" {
  map = {
    "key with spaces" = {
      value = "spaces"
    }
    key_with_underscores = {
      value = "underscores"
    }
  }
}`,
		},
		"simple_map_with_periods_in_keys": {
			schema: &configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"map": {
						Type:     cty.Map(cty.String),
						Optional: true,
					},
				},
			},
			addr: addrs.AbsResourceInstance{
				Module: addrs.RootModuleInstance,
				Resource: addrs.ResourceInstance{
					Resource: addrs.Resource{
						Mode: addrs.ManagedResourceMode,
						Type: "testing_resource",
						Name: "resource",
					},
					Key: addrs.NoKey,
				},
			},
			provider: addrs.LocalProviderConfig{
				LocalName: "testing",
			},
			value: cty.ObjectVal(map[string]cty.Value{
				"map": cty.MapVal(map[string]cty.Value{
					"key.with.periods":     cty.StringVal("periods"),
					"key_with_underscores": cty.StringVal("underscores"),
				}),
			}),
			expected: `resource "testing_resource" "resource" {
  map = {
    "key.with.periods"   = "periods"
    key_with_underscores = "underscores"
  }
}`,
		},
		"nested_map_with_periods_in_keys": {
			schema: &configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"map": {
						NestedType: &configschema.Object{
							Attributes: map[string]*configschema.Attribute{
								"value": {
									Type:     cty.String,
									Optional: true,
								},
							},
							Nesting: configschema.NestingMap,
						},
						Optional: true,
					},
				},
			},
			addr: addrs.AbsResourceInstance{
				Module: addrs.RootModuleInstance,
				Resource: addrs.ResourceInstance{
					Resource: addrs.Resource{
						Mode: addrs.ManagedResourceMode,
						Type: "testing_resource",
						Name: "resource",
					},
					Key: addrs.NoKey,
				},
			},
			provider: addrs.LocalProviderConfig{
				LocalName: "testing",
			},
			value: cty.ObjectVal(map[string]cty.Value{
				"map": cty.MapVal(map[string]cty.Value{
					"key.with.periods": cty.ObjectVal(map[string]cty.Value{
						"value": cty.StringVal("periods"),
					}),
					"key_with_underscores": cty.ObjectVal(map[string]cty.Value{
						"value": cty.StringVal("underscores"),
					}),
				}),
			}),
			expected: `resource "testing_resource" "resource" {
  map = {
    "key.with.periods" = {
      value = "periods"
    }
    key_with_underscores = {
      value = "underscores"
    }
  }
}`,
		},
	}
	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			err := tc.schema.InternalValidate()
			if err != nil {
				t.Fatalf("schema failed InternalValidate: %s", err)
			}
			contents, diags := GenerateResourceContents(tc.addr, tc.schema, tc.provider, tc.value, false)
			if len(diags) > 0 {
				t.Errorf("expected no diagnostics but found %s", diags)
			}

			got := contents.String()
			want := strings.TrimSpace(tc.expected)
			if diff := cmp.Diff(got, want); len(diff) > 0 {
				t.Errorf("got:\n%s\nwant:\n%s\ndiff:\n%s", got, want, diff)
			}
		})
	}
}

func sensitiveAttribute(t cty.Type) *configschema.Attribute {
	return &configschema.Attribute{
		Type:      t,
		Optional:  true,
		Sensitive: true,
	}
}

func TestGenerateResourceAndIDContents(t *testing.T) {
	schema := &configschema.Block{
		Attributes: map[string]*configschema.Attribute{
			"name": {
				Type:     cty.String,
				Optional: true,
			},
			"id": {
				Type:     cty.String,
				Computed: true,
			},
			"tags": {
				Type:     cty.Map(cty.String),
				Optional: true,
			},
		},
		BlockTypes: map[string]*configschema.NestedBlock{
			"network_interface": {
				Nesting: configschema.NestingList,
				Block: configschema.Block{
					Attributes: map[string]*configschema.Attribute{
						"subnet_id": {
							Type:     cty.String,
							Required: true,
						},
						"ip_address": {
							Type:     cty.String,
							Optional: true,
						},
					},
				},
			},
		},
	}

	// Define the identity schema
	idSchema := &configschema.Object{
		Nesting: configschema.NestingSingle,
		Attributes: map[string]*configschema.Attribute{
			"id": {
				Type:     cty.String,
				Optional: true,
			},
		},
	}

	// Create mock resource instance values
	value := cty.TupleVal([]cty.Value{
		cty.ObjectVal(map[string]cty.Value{
			"state": cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("instance-1"),
				"id":   cty.StringVal("i-abcdef"),
				"tags": cty.MapVal(map[string]cty.Value{
					"Environment": cty.StringVal("Dev"),
					"Owner":       cty.StringVal("Team1"),
				}),
				"network_interface": cty.ListVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"subnet_id":  cty.StringVal("subnet-123"),
						"ip_address": cty.StringVal("10.0.0.1"),
					}),
				}),
			}),
			"identity": cty.ObjectVal(map[string]cty.Value{
				"id": cty.StringVal("i-abcdef"),
			}),
		}),
		cty.ObjectVal(map[string]cty.Value{
			"state": cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("instance-2"),
				"id":   cty.StringVal("i-123456"),
				"tags": cty.MapVal(map[string]cty.Value{
					"Environment": cty.StringVal("Prod"),
					"Owner":       cty.StringVal("Team2"),
				}),
				"network_interface": cty.ListVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"subnet_id":  cty.StringVal("subnet-456"),
						"ip_address": cty.StringVal("10.0.0.2"),
					}),
				}),
			}),
			"identity": cty.ObjectVal(map[string]cty.Value{
				"id": cty.StringVal("i-123456"),
			}),
		}),
	})

	// Create test resource address
	addr := addrs.AbsResource{
		Module: addrs.RootModuleInstance,
		Resource: addrs.Resource{
			Mode: addrs.ListResourceMode,
			Type: "aws_instance",
			Name: "example",
		},
	}

	// Create instance addresses for each instance
	instAddr1 := addr.Instance(addrs.NoKey)

	// Create provider config
	pc := addrs.LocalProviderConfig{
		LocalName: "aws",
	}

	// Generate content
	content, diags := GenerateListResourceContents(instAddr1, schema, idSchema, pc, value)
	// Check for diagnostics
	if diags.HasErrors() {
		t.Fatalf("unexpected diagnostics: %s", diags.Err())
	}

	// Check the generated content
	expectedContent := `resource "aws_instance" "example_0" {
  provider = aws
  name     = "instance-1"
  tags = {
    Environment = "Dev"
    Owner       = "Team1"
  }
  network_interface {
    ip_address = "10.0.0.1"
    subnet_id  = "subnet-123"
  }
}
import {
  to       = aws_instance.example_0
  provider = aws
  identity = {
    id = "i-abcdef"
  }
}

resource "aws_instance" "example_1" {
  provider = aws
  name     = "instance-2"
  tags = {
    Environment = "Prod"
    Owner       = "Team2"
  }
  network_interface {
    ip_address = "10.0.0.2"
    subnet_id  = "subnet-456"
  }
}
import {
  to       = aws_instance.example_1
  provider = aws
  identity = {
    id = "i-123456"
  }
}
`
	// Normalize both strings by removing extra whitespace for comparison
	normalizeString := func(s string) string {
		// Remove spaces at the end of lines and replace multiple newlines with a single one
		lines := strings.Split(s, "\n")
		for i, line := range lines {
			lines[i] = strings.TrimRight(line, " \t")
		}
		return strings.Join(lines, "\n")
	}

	normalizedExpected := normalizeString(expectedContent)

	var merged string
	res := content.Results
	for _, addr := range res {
		merged += addr.String()
	}
	normalizedActual := normalizeString(content.String())

	if diff := cmp.Diff(normalizedExpected, normalizedActual); diff != "" {
		t.Errorf("Generated content doesn't match expected. want:\n%s\ngot:\n%s\ndiff:\n%s", normalizedExpected, normalizedActual, diff)
	}
}
