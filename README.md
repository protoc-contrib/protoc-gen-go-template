# protoc-gen-go-template

A generic protoc plugin that renders arbitrary files from Go
[`text/template`](https://pkg.go.dev/text/template) sources driven by the
parsed proto AST.

The plugin walks a user-provided template directory, parses every `.tmpl`
file, and writes the rendered output into the protoc response. The template
engine is extended with helpers from
[Masterminds/sprig](https://github.com/Masterminds/sprig) plus a proto-aware
funcmap (type manipulation, HTTP annotation accessors, naming helpers, etc.).

## Philosophy

- protobuf-first — every output is derived from a parsed proto AST
- no built-in templates — the plugin is a runner for user-owned templates
- keep it stupid simple

## Install

```console
go install github.com/protoc-contrib/protoc-gen-go-template/cmd/protoc-gen-go-template@latest
```

Or build with Nix:

```console
nix build .
```

## Usage

Every file ending in `.tmpl` under `template_dir` is rendered and written to
the destination folder, preserving the directory layout under `template_dir`
and stripping the `.tmpl` suffix.

```console
$ ls -R
input.proto  templates/doc.txt.tmpl  templates/config.json.tmpl
$ protoc --go-template_out=. input.proto
$ ls -R
input.proto  templates/doc.txt.tmpl  templates/config.json.tmpl
doc.txt      config.json
```

### Options

```console
$ protoc --go-template_out=debug=true,template_dir=/path/to/templates:. input.proto
```

| Option                | Default      | Values           | Description                                                                 |
|-----------------------|--------------|------------------|-----------------------------------------------------------------------------|
| `template_dir`        | `./template` | path             | root directory containing `.tmpl` files                                     |
| `destination_dir`     | `.`          | path             | base path written into each output file name                                |
| `debug`               | `false`      | bool             | verbose template-engine logging                                             |
| `all`                 | `false`      | bool             | run templates against every proto file, not only those that define services |
| `single-package-mode` | `false`      | bool             | load the full request into a grpc-gateway registry for cross-file lookups   |
| `file-mode`           | `false`      | bool             | run templates once per file (that declares a service) rather than per service |

## Funcmap

The template engine is loaded with all [sprig](https://github.com/Masterminds/sprig)
helpers plus a proto-aware set: naming (`camelCase`, `lowerCamelCase`,
`snakeCase`, `kebabCase`, `goNormalize`, `shortType`, ...), arithmetic
(`add`, `subtract`, `multiply`, `divide`), message/field walkers
(`getMessageType`, `isFieldMessage`, `fieldMapKeyType`, ...), HTTP annotation
accessors (`httpPath`, `httpVerb`, `httpBody`, ...), and extension readers
(`stringFieldExtension`, `boolFieldExtension`, ...). See
[`internal/generator/helpers.go`](./internal/generator/helpers.go) for the
complete list.

## Repository layout

```
cmd/protoc-gen-go-template/   binary entry point (protogen plugin)
internal/generator/           template engine, funcmap, encoders
.github/                      CI workflows, release-please config, dependabot
flake.nix                     Nix build (buildGoModule)
```

## Credits

This project is a continuation of
[moul/protoc-gen-gotemplate](https://github.com/moul/protoc-gen-gotemplate)
by Manfred Touron and contributors, republished under the `protoc-contrib`
org and renamed to `protoc-gen-go-template` to align with common protoc
plugin naming conventions. Original git history is preserved.

### Migrating from `protoc-gen-gotemplate`

- Binary renamed: `protoc-gen-gotemplate` → `protoc-gen-go-template`
- Protoc flag renamed: `--gotemplate_out` → `--go-template_out`
- Install path: `go install github.com/protoc-contrib/protoc-gen-go-template/cmd/protoc-gen-go-template@latest`

## License

MIT — see [`LICENSE.md`](./LICENSE.md).
