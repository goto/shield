package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/goto/salt/config"
	"github.com/goto/salt/telemetry"
	"github.com/goto/shield/internal/proxy"
	"github.com/goto/shield/internal/server"
	"github.com/goto/shield/internal/store/spicedb"
	"github.com/goto/shield/pkg/db"
)

type Log struct {
	// log level - debug, info, warning, error, fatal
	Level string `yaml:"level" mapstructure:"level" default:"info"`

	// format strategy - plain, json
	Format string `yaml:"format" mapstructure:"format" default:"json"`

	Activity ActivityLogConfig `yaml:"activity" mapstructure:"activity"`
}

type ActivityLogConfig struct {
	// activity sink strategy - db, stdout, none
	Sink string `yaml:"sink" mapstructure:"sink" default:"none"`
}

type Shield struct {
	// configuration version
	Version   int                  `yaml:"version"`
	Proxy     proxy.ServicesConfig `yaml:"proxy"`
	Log       Log                  `yaml:"log"`
	NewRelic  NewRelic             `yaml:"new_relic"`
	App       server.Config        `yaml:"app"`
	DB        db.Config            `yaml:"db"`
	SpiceDB   spicedb.Config       `yaml:"spicedb"`
	Telemetry telemetry.Config     `yaml:"telemetry" mapstructure:"telemetry"`
}

type NewRelic struct {
	AppName string `yaml:"app_name" mapstructure:"app_name"`
	License string `yaml:"license" mapstructure:"license"`
	Enabled bool   `yaml:"enabled" mapstructure:"enabled"`
}

func Load(serverConfigFileFromFlag string) (*Shield, error) {
	conf := &Shield{}

	var options []config.LoaderOption
	options = append(options, config.WithName("config"))
	options = append(options, config.WithEnvKeyReplacer(".", "_"))
	options = append(options, config.WithEnvPrefix("SHIELD"))
	if p, err := os.Getwd(); err == nil {
		options = append(options, config.WithPath(p))
	}
	if execPath, err := os.Executable(); err == nil {
		options = append(options, config.WithPath(filepath.Dir(execPath)))
	}
	if currentHomeDir, err := os.UserHomeDir(); err == nil {
		options = append(options, config.WithPath(currentHomeDir))
		options = append(options, config.WithPath(filepath.Join(currentHomeDir, ".config")))
	}

	// override all config sources and prioritize one from file
	if serverConfigFileFromFlag != "" {
		options = []config.LoaderOption{config.WithFile(serverConfigFileFromFlag)}
	}

	l := config.NewLoader(options...)
	if err := l.Load(conf); err != nil {
		if !errors.As(err, &config.ConfigFileNotFoundError{}) {
			return nil, err
		}
	}
	return conf, nil
}
