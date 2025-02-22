// nolint:gochecknoglobals
package conf

import (
	"time"

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
	Access     JwtToken `yaml:"access"`
	Refresh    JwtToken `yaml:"refresh"`
	BcryptCost int      `yaml:"bcrypt_cost" validate:"min=3,max=31"`
}

type JwtToken struct {
	TTL    time.Duration `yaml:"ttl" validate:"required,min=24h"`
	Secret string        `yaml:"secret" validate:"required,min=10"`
}
