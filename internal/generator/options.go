package generator

import (
	"fmt"
	"strconv"
)

// Options controls how the plugin discovers templates and emits output.
type Options struct {
	// TemplateDir is the root directory containing .tmpl files to render.
	TemplateDir string
	// DestinationDir is the prefix written into each output file's name.
	DestinationDir string
	// Debug enables verbose logging from the template engine.
	Debug bool
	// All, when true, runs the template set against every proto file in the
	// request rather than only those that declare a service.
	All bool
	// SinglePackageMode loads the full request into a grpc-gateway registry
	// so templates can walk cross-file message references.
	SinglePackageMode bool
	// FileMode, when true, runs the template set once per proto file that
	// declares at least one service (rather than once per service).
	FileMode bool
}

// Set applies a single `name=value` plugin parameter to the options. The
// signature matches what protogen.Options.ParamFunc expects.
func (o *Options) Set(name, value string) error {
	switch name {
	case "template_dir":
		o.TemplateDir = value
		return nil
	case "destination_dir":
		o.DestinationDir = value
		return nil
	case "debug":
		return setBool(&o.Debug, name, value)
	case "all":
		return setBool(&o.All, name, value)
	case "single-package-mode":
		return setBool(&o.SinglePackageMode, name, value)
	case "file-mode":
		return setBool(&o.FileMode, name, value)
	default:
		return fmt.Errorf("unknown plugin option %q", name)
	}
}

func setBool(dst *bool, name, value string) error {
	v, err := strconv.ParseBool(value)
	if err != nil {
		return fmt.Errorf("invalid value for %q: %w", name, err)
	}
	*dst = v
	return nil
}
