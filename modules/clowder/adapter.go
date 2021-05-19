package clowder

import (
	"bytes"
	"text/template"

	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/ghodss/yaml"
	clowdapp "github.com/redhatinsights/app-common-go/pkg/api/v1"
)

func init() {
	caddyconfig.RegisterAdapter("clowder", Adapter{})
}

// Adapter enhances a YAML config with values detected from a ClowdApp config
type Adapter struct{}

type clowderConfig struct {
	AppConfig   clowdapp.AppConfig
	Deps        map[string]map[string]clowdapp.DependencyEndpoint
	PrivateDeps map[string]map[string]clowdapp.PrivateDependencyEndpoint
}

// Adapt performs configuration substitutions based on Clowder config detected from environment
func (a Adapter) Adapt(body []byte, options map[string]interface{}) (result []byte, warnings []caddyconfig.Warning, err error) {
	// Bail early if Clowder is not enabled
	if !clowdapp.IsClowderEnabled() {
		warnings = append(warnings, caddyconfig.Warning{
			Message: "Clowder is not enabled.",
		})
		return
	}

	// Template out the provided YAML file
	var tmpl *template.Template
	tmpl, err = template.New("config").Parse(string(body))
	if err != nil {
		return
	}
	var out bytes.Buffer
	err = tmpl.Execute(&out, clowderConfig{
		AppConfig:   *clowdapp.LoadedConfig,
		Deps:        clowdapp.DependencyEndpoints,
		PrivateDeps: clowdapp.PrivateDependencyEndpoints,
	})
	if err != nil {
		return
	}

	// Finally convert templated YAML to JSON
	result, err = yaml.YAMLToJSON(out.Bytes())
	return
}

// Interface guard
var _ caddyconfig.Adapter = (*Adapter)(nil)
