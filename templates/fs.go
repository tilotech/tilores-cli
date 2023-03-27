package templates

import "embed"

var (
	// InitPreGenerate contains the templates that need to be available before go generate is run the first time.
	//
	//go:embed init/schema init/tools init/generate.go.tmpl init/gqlgen.yml.tmpl init/rule-config.json.tmpl init/graph/model
	InitPreGenerate embed.FS

	// InitPostGenerate contains the templates that need to be available after go generate is run the first time.
	//
	//go:embed init/graph init/cmd init/handler init/deployment
	InitPostGenerate embed.FS

	// Upgrades contains the upgrade files and their templates to upgrade to a specific version.
	//
	//go:embed upgrades
	Upgrades embed.FS
)
