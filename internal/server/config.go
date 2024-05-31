package server

import (
	"fmt"

	"github.com/goto/shield/pkg/telemetry"
)

type GRPCConfig struct {
	Port           int `mapstructure:"port" default:"8081"`
	MaxRecvMsgSize int `mapstructure:"max_recv_msg_size" default:"33554432"`
	MaxSendMsgSize int `mapstructure:"max_send_msg_size" default:"33554432"`
}

type ServiceDataConfig struct {
	BootstrapServiceDataKey bool `yaml:"bootstrap_service_data_key" mapstructure:"bootstrap_service_data_key" default:"true"`
	MaxUpsert               int  `yaml:"max_upsert" mapstructure:"max_upsert" default:"1"`
}

func (cfg Config) grpcAddr() string { return fmt.Sprintf("%s:%d", cfg.Host, cfg.GRPC.Port) }

type Config struct {
	// port to listen HTTP requests on
	Port int `yaml:"port" mapstructure:"port" default:"8080"`

	// GRPC Config
	GRPC GRPCConfig `mapstructure:"grpc"`

	//metrics port
	MetricsPort int `yaml:"metrics_port" mapstructure:"metrics_port" default:"9000"`

	// the network interface to listen on
	Host string `yaml:"host" mapstructure:"host" default:"127.0.0.1"`

	Name string

	// RulesPath is a directory path where ruleset is defined
	// that this service should implement
	RulesPath string `yaml:"ruleset" mapstructure:"ruleset"`
	// RulesPathSecret could be a env name, file path or actual value required
	// to access RulesPath files
	RulesPathSecret string `yaml:"ruleset_secret" mapstructure:"ruleset_secret"`

	// TODO might not suitable here because it is also being used by proxy
	// Headers which will have user's email id
	IdentityProxyHeader string `yaml:"identity_proxy_header" mapstructure:"identity_proxy_header" default:"X-Shield-Email"`

	// Header which will have user_id
	UserIDHeader string `yaml:"user_id_header" mapstructure:"user_id_header" default:"X-Shield-User-Id"`

	// ResourcesPath is a directory path where resources is defined
	// that this service should implement
	ResourcesConfigPath string `yaml:"resources_config_path" mapstructure:"resources_config_path"`

	// ResourcesPathSecretSecret could be a env name, file path or actual value required
	// to access ResourcesPathSecretPath files
	ResourcesConfigPathSecret string `yaml:"resources_config_path_secret" mapstructure:"resources_config_path_secret"`

	TelemetryConfig telemetry.Config `yaml:"telemetry_config" mapstructure:"telemetry_config"`

	// CheckAPILimit will have the maximum number of resource permissions that can be included
	// in the resource permission check API. Default: 5
	CheckAPILimit int `yaml:"check_api_limit" mapstructure:"check_api_limit" default:"5"`

	DefaultSystemEmail string `yaml:"default_system_email" mapstructure:"default_system_email"  default:"shield-service@gotocompany.com"`

	ServiceData ServiceDataConfig `yaml:"service_data" mapstructure:"service_data"`
}
