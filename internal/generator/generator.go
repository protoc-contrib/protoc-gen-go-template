// Package generator renders arbitrary files from Go text/template sources
// driven by a protoc CodeGeneratorRequest.
package generator

import (
	"fmt"

	ggdescriptor "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
	"google.golang.org/protobuf/compiler/protogen"
)

// Generate walks the request the plugin was invoked with, applies the
// template set rooted at opts.TemplateDir once per service (or per file,
// depending on the mode flags), and attaches the rendered output to the
// plugin response. Files that share a name have their content concatenated
// in arrival order.
func Generate(plugin *protogen.Plugin, opts *Options) error {
	if opts == nil {
		opts = &Options{}
	}

	req := plugin.Request

	if opts.SinglePackageMode {
		reg := ggdescriptor.NewRegistry()
		SetRegistry(reg)
		if err := reg.Load(req); err != nil {
			return fmt.Errorf("registry: failed to load the request: %w", err)
		}
	}

	emitted := map[string]*generatedFile{}
	emit := func(name, content string) {
		if gf, ok := emitted[name]; ok {
			gf.content += content
			return
		}
		emitted[name] = &generatedFile{name: name, content: content}
	}

	for _, file := range req.GetProtoFile() {
		switch {
		case opts.All:
			enc := NewGenericTemplateBasedEncoder(opts.TemplateDir, file, opts.Debug, opts.DestinationDir)
			for _, tmpl := range enc.Files() {
				emit(tmpl.GetName(), tmpl.GetContent())
			}
		case opts.FileMode:
			if len(file.GetService()) == 0 {
				continue
			}
			enc := NewGenericTemplateBasedEncoder(opts.TemplateDir, file, opts.Debug, opts.DestinationDir)
			for _, tmpl := range enc.Files() {
				emit(tmpl.GetName(), tmpl.GetContent())
			}
		default:
			for _, service := range file.GetService() {
				enc := NewGenericServiceTemplateBasedEncoder(opts.TemplateDir, service, file, opts.Debug, opts.DestinationDir)
				for _, tmpl := range enc.Files() {
					emit(tmpl.GetName(), tmpl.GetContent())
				}
			}
		}
	}

	for _, gf := range emitted {
		out := plugin.NewGeneratedFile(gf.name, "")
		if _, err := out.Write([]byte(gf.content)); err != nil {
			return fmt.Errorf("write %q: %w", gf.name, err)
		}
	}
	return nil
}

type generatedFile struct {
	name    string
	content string
}
