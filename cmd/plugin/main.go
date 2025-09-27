package main

import (
	hplugin "github.com/hashicorp/go-plugin"

	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"

	"github.com/jpower432/opensource-securitycon-2025-oscal-in-action/cmd/plugin/server"
)

func main() {
	opaPlugin := server.NewPlugin()
	plugins := map[string]hplugin.Plugin{
		plugin.PVPPluginName: &plugin.PVPPlugin{Impl: opaPlugin},
	}
	config := plugin.ServeConfig{
		PluginSet: plugins,
		Logger:    server.Logger(),
	}
	plugin.Register(config)
}
