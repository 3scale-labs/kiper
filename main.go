package main

import (
	"fmt"
	"os"

	"github.com/3scale/kiper/pkg/queries"

	"github.com/open-policy-agent/opa/runtime"

	"github.com/3scale/kiper/internal/istio_plugin"
	"github.com/open-policy-agent/opa/cmd"
	"github.com/open-policy-agent/opa/plugins"
)

// Factory defines the interface OPA uses to instantiate a plugin.
type Factory struct{}

// New returns the object initialized with a valid plugin configuration.
func (Factory) New(m *plugins.Manager, config interface{}) plugins.Plugin {
	return istio_plugin.New(m, config.(*istio_plugin.Config))
}

// Validate returns a valid configuration to instantiate the plugin.
func (Factory) Validate(m *plugins.Manager, config []byte) (interface{}, error) {
	return istio_plugin.Validate(m, config)
}

func main() {
	runtime.RegisterPlugin("envoy.ext_authz.grpc", Factory{}) // for backwards compatibility
	runtime.RegisterPlugin("envoy_ext_authz_grpc", Factory{})

	queries.RegisterThreeScaleQueries()
	queries.RegisterRateLimitQueries()

	if err := cmd.RootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
