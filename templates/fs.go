package templates

import "embed"

var (
	//go:embed init/schema init/tools init/generate.go.tmpl init/gqlgen.yml.tmpl init/rule-config.json.tmpl init/graph/model
	InitPreGenerate embed.FS

	//go:embed init/server.go.tmpl init/graph init/cmd init/handler init/deployment
	InitPostGenerate embed.FS
)
