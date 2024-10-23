package app

import (
	"os"

	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Github githubapp.Config `yaml:"github"`
	App    AppConfig        `yaml:"app"`
}

type AppConfig struct {
	PullRequestPreamble string `yaml:"pull_request_preamble"`
}

func ReadConfig(path string) (*Config, error) {
	var c Config

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed reading server config file: %s", path)
	}

	if err := yaml.UnmarshalStrict(bytes, &c); err != nil {
		return nil, errors.Wrap(err, "failed parsing configuration file")
	}

	return &c, nil
}
