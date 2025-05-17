package config

import (
	"os"
)

func GetEnvSecretKey() string {
	if envSecretKey == nil {
		envSecretKey = new(string)
		*envSecretKey = os.Getenv(envar_SLV_ENV_SECRET_KEY)
	}
	return *envSecretKey
}

func GetEnvSecretBinding() string {
	if envSecretBinding == nil {
		envSecretBinding = new(string)
		*envSecretBinding = os.Getenv(envar_SLV_ENV_SECRET_BINDING)
	}
	return *envSecretBinding
}

func GetEnvSecretPassword() string {
	if envSecretPassword == nil {
		envSecretPassword = new(string)
		*envSecretPassword = os.Getenv(envar_SLV_ENV_SECRET_PASSWORD)
	}
	return *envSecretPassword
}
