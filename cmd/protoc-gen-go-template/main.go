// Command protoc-gen-go-template is a protoc plugin that renders arbitrary
// files from Go text/template sources, driven by the parsed proto AST.
package main

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/protoc-contrib/protoc-gen-go-template/internal/generator"
)

func main() {
	opts := &generator.Options{}
	protogen.Options{
		ParamFunc: opts.Set,
	}.Run(func(plugin *protogen.Plugin) error {
		plugin.SupportedFeatures = uint64(
			pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL |
				pluginpb.CodeGeneratorResponse_FEATURE_SUPPORTS_EDITIONS,
		)
		plugin.SupportedEditionsMinimum = descriptorpb.Edition_EDITION_2023
		plugin.SupportedEditionsMaximum = descriptorpb.Edition_EDITION_2023
		return generator.Generate(plugin, opts)
	})
}
