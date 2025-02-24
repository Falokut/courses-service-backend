// nolint:gochecknoglobals
package conf

import (
	"github.com/Falokut/go-kit/config"
)

type LocalConfig struct {
	HealthcheckPort uint32          `yaml:"healthcheck_port" env:"HEALTHCHECK_PORT" env-default:"8081"`
	DB              config.Database `yaml:"db"`
	Listen          config.Listen   `yaml:"listen"`
	Images          Images          `yaml:"images"`
	Auth            Auth            `yaml:"auth"`
}

type Images struct {
	BaseServiceUrl string `yaml:"base_service_url" env:"IMAGES_BASE_SERVICE_URL"`
	BaseImagePath  string `yaml:"base_image_path" env:"BASE_IMAGE_PATH"`
}

type Auth struct {
	BcryptCost int        `yaml:"bcrypt_cost" validate:"min=3,max=31"`
	InitAdmin  *InitAdmin `yaml:"init_admin"`
}

type InitAdmin struct {
	Username string `yaml:"username" validate:"min=3,max=40"`
	Password string `yaml:"password" validate:"min=3,max=40"`
}
