package internal

import (
	"github.com/google/uuid"
	"strings"
)

type Config struct {
	AuthSecret string
	AutoSecret bool
	DBFileName string
}

func NewDefaultConfig() *Config {
	maxPasswordLength := 16

	return &Config{
		AuthSecret: strings.Replace(uuid.New().String()[:maxPasswordLength], "-", "", -1),
		AutoSecret: true,
		DBFileName: "minimi.db",
	}
}
