// nolint:gochecknoglobals
package conf

import (
	"github.com/Falokut/go-kit/config"
)

type LocalConfig struct {
	HealthcheckPort uint32          `yaml:"healthcheck_port" env:"HEALTHCHECK_PORT" env-default:"8081"`
	DB              config.Database `yaml:"db"`
	Listen          config.Listen   `yaml:"listen"`
	FileStorage     FileStorage     `yaml:"file_storage"`
	Auth            Auth            `yaml:"auth"`
}

type FileStorage struct {
	BaseServiceUrl string `yaml:"base_service_url" env:"FILE_STORAGE_BASE_SERVICE_URL"`
}

type Auth struct {
	BcryptCost int        `yaml:"bcrypt_cost" validate:"min=3,max=31"`
	InitAdmin  *InitAdmin `yaml:"init_admin"`
}

type InitAdmin struct {
	Username string `yaml:"username" validate:"min=3,max=40"`
	Password string `yaml:"password" validate:"min=3,max=40"`
}
