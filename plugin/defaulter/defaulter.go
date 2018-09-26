// Package defaulter is a gogoproto plugin to generate Default() functions
// for types, satisfying the Defaulter interface. If a default value for
// a protobuf field is specified, the value given will be assigned to the
// field in the struct and represented so that it can be parsed by Go and
// assigned to the appropriate type. Because of complexity issues it can
// only support basic type assignments like string, bool, and integers.
package defaulter

import (
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	sensuproto "github.com/sensu/sensu-proto/protobuf"
)

func init() {
	generator.RegisterPlugin(NewDefaulter())
}

// GetDefault returns the default value for a field as a string.
func GetDefault(field *descriptor.FieldDescriptorProto) *string {
	if field == nil {
		return nil
	}

	if field.Options != nil {
		v, err := proto.GetExtension(field.Options, sensuproto.E_Default)
		if err == nil && v.(*string) != nil {
			return (v.(*string))
		}
	}

	return nil
}

// Plugin is the default plugin.
type Plugin struct {
	*generator.Generator
	generator.PluginImports
	messages []*generator.Descriptor
}

// NewDefaulter creates a new Defaulter generator
func NewDefaulter() *Plugin {
	return &Plugin{}
}

// Name returns the name of the plugin.
func (p *Plugin) Name() string {
	return "defaulter"
}

// Init initializes the plugin with the given generator.
func (p *Plugin) Init(g *generator.Generator) {
	p.Generator = g
}

// Generate the output for this plugin.
func (p *Plugin) Generate(file *generator.FileDescriptor) {
	p.PluginImports = generator.NewPluginImports(p.Generator)
	p.messages = make([]*generator.Descriptor, 0)

	for _, message := range file.Messages() {
		// Default()
		// Generates the Default() function for the type.

		// e.g. Rule
		baseTypeName := generator.CamelCaseSlice(message.TypeName())

		// "r" for use in func(r *Rule)
		typeShort := strings.ToLower(string(baseTypeName[0]))
		// Gives a newline
		p.P()

		// func (r *Rule) Default() {
		p.P(`func (`, typeShort, `*`, baseTypeName, `) Default() {`)

		// Indent
		p.In()

		//
		if meta := message.GetFieldDescriptor("metadata"); meta != nil {
			// this goes into every type.

			// r.Kind = "Role"
			p.P(typeShort, `.Kind = "`, *message.Name, `"`)

			// r.ApiVersion = SchemeGroupVersion.GroupVersionString()
			// (when executed by Go): r.ApiVersion = rbac.sensu.io/v2alpha1
			p.P(typeShort, `.ApiVersion = SchemeGroupVersion.GroupVersionString()`)
		}

		// Now do field specific defaults. -->
		for _, field := range message.Field {
			// e.g. Namespace
			fieldName := p.GetFieldName(message, field)

			// GetDefault(field) -> *"somevalue"
			if sensuDefault := GetDefault(field); sensuDefault != nil {
				defaultValue := *sensuDefault

				if field.IsString() {
					// string namespace = 1 [(sensuproto.default = "default")]
					// r.Namespace = "default";
					p.P(typeShort, `.`, fieldName, `= "`, defaultValue, `"`)
				} else if field.IsScalar() {
					// integer version = 1 [(sensuproto.default = 1)];
					// r.Version = 1
					p.P(typeShort, `.`, fieldName, `= `, defaultValue)
				}
			}
		}

		// Unindent
		p.Out()
		p.P(`}`)
		// Gives a newline
		p.P()
	}
}
