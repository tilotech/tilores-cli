package templates

import "embed"

var (
	//go:embed init/schema init/tools init/generate.go.tmpl init/gqlgen.yml.tmpl
	InitPreGenerate embed.FS

	//go:embed init/server.go.tmpl
	InitPostGenerate embed.FS
)
