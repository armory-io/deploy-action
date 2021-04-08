package action

import (
	"fmt"

	"github.com/codingconcepts/env"
)

type Config struct {
	CloudProvider string   `env:"INPUT_CLOUDPROVIDER" required:"true"`
	Accounts      []string `env:"INPUT_ACCOUNTS" required:"true"`
	Application   string   `env:"INPUT_APPLICATION" required:"true"`

	Wait    bool   `env:"INPUT_WAIT"`
	BaseURL string `env:"INPUT_BASEURL" required:"true"`
	Config  string `env:"INPUT_CONFIGPATH" required:"true"`

	ManifestPath string `env:"INPUT_MANIFESTPATH"`
	Manifest     string `env:"INPUT_MANIFEST"`
}

func (c Config) Validate() error {
	if c.CloudProvider != "kubernetes" {
		return fmt.Errorf(
			"only kubernetes deployments are supported at this time. invalid cloudProvider specified: %s", c.CloudProvider)
	}
	return nil
}

func GetConfig() (*Config, error) {
	var c Config
	// TODO: not very testable due to dependence on os.Env. fix.
	if err := env.Set(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
