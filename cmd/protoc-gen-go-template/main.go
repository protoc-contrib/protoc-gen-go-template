// Command protoc-gen-go-template is a protoc plugin that renders arbitrary
// files from Go text/template sources, driven by the parsed proto AST.
package main

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/protoc-contrib/protoc-gen-go-template/internal/generator"
)

func main() {
	opts := &generator.Options{}
	protogen.Options{
		ParamFunc: opts.Set,
	}.Run(func(plugin *protogen.Plugin) error {
		plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		return generator.Generate(plugin, opts)
	})
}
