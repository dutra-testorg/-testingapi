package config

import (
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v2"
)

// ServiceConfig ...
type ServiceConfig struct {
	Server          serverInfo    `yaml:"server" json:"server"`
	Cors            corsInfo      `yaml:"cors" json:"cors"`
	Datadog         datadogInfo   `yaml:"datadog" json:"datadog"`
	RestAPI         restAPIInfo   `yaml:"rest_api" json:"rest_api"`
	Environment     string        `envconfig:"DD_ENV" yaml:"environment"`
	CursorKey       string        `envconfig:"CURSOR_KEY" yaml:"cursor_key" json:"cursor_key" split_words:"true"`
	ServiceName     string        `envconfig:"SERVICE_NAME" yaml:"service_name" json:"service_name" split_words:"true"`
	LogLevel        zapcore.Level `envconfig:"LOG_LEVEL" yaml:"log_level" json:"log_level" split_words:"true"`
	LogDump         bool          `envconfig:"LOG_DUMP" yaml:"log_dump" json:"log_dump" split_words:"true"`
	ProfilerEnabled bool          `envconfig:"PROFILER_ENABLED" yaml:"profiler_enabled" json:"profiler_enabled" split_words:"true"`
	SwaggerEnabled  bool          `envconfig:"SWAGGER_ENABLED" yaml:"swagger_enabled" json:"swagger_enabled" split_words:"true"`
}

type serverInfo struct {
	Address         string        `envconfig:"SERVER_ADDRESS" yaml:"address" json:"address"`
	WriteTimeout    time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" yaml:"write_timeout" json:"write_timeout" split_words:"true"`
	ReadTimeout     time.Duration `envconfig:"SERVER_READ_TIMEOUT" yaml:"read_timeout" json:"read_timeout" split_words:"true"`
	IdleTimeout     time.Duration `envconfig:"SERVER_IDLE_TIMEOUT" yaml:"idle_timeout" json:"idle_timeout" split_words:"true"`
	ShutdownTimeout time.Duration `envconfig:"SERVER_SHUTDOWN_TIMEOUT" yaml:"shutdown_timeout" json:"shutdown_timeout" split_words:"true"`
}

type corsInfo struct {
	AllowedHeaders []string `envconfig:"CORS_ALLOWED_HEADERS" yaml:"allowed_headers" json:"allowed_headers" split_words:"true"`
	AllowedMethods []string `envconfig:"CORS_ALLOWED_METHODS" yaml:"allowed_methods" json:"allowed_methods" split_words:"true"`
	AllowedOrigins []string `envconfig:"CORS_ALLOWED_ORIGINS" yaml:"allowed_origins" json:"allowed_origins" split_words:"true"`
	ExposedHeaders []string `envconfig:"CORS_EXPOSED_HEADERS" yaml:"exposed_headers" json:"exposed_headers" split_words:"true"`
	MaxAge         int      `envconfig:"CORS_MAX_AGE" yaml:"max_age" json:"max_age" split_words:"true"`
}

type datadogInfo struct {
	Host    string `envconfig:"DATADOG_HOST" yaml:"host" json:"host"`
	Port    string `envconfig:"DATADOG_PORT" yaml:"port" json:"port"`
	Enabled bool   `envconfig:"DATADOG_ENABLED" yaml:"enabled" json:"enabled"`
}

type restAPIInfo struct {
	Address        string        `envconfig:"REST_SERVER_ADDRESS" yaml:"address" json:"address"`
	RequestTimeout time.Duration `envconfig:"REST_SERVER_REQUEST_TIMEOUT" yaml:"request_timeout" json:"request_timeout" split_words:"true"`
}

// LoadServiceConfig ...
func LoadServiceConfig(configFile string) (*ServiceConfig, error) {
	var cfg ServiceConfig

	if err := loadServiceConfigFromFile(configFile, &cfg); err != nil {
		return nil, err
	}

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func loadServiceConfigFromFile(configFile string, cfg *ServiceConfig) error {
	_, err := os.Stat(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(yamlFile, &cfg)
}
