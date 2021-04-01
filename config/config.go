package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Postgres struct {
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Host         string `yaml:"host"`
	Name         string `yaml:"name"`
	MaxIdleConns int    `yaml:"maxIdleConns"`
	MaxOpenConns int    `yaml:"maxOpenConns"`
	DisableTLS   bool   `yaml:"disableTLS"`
}

type Redis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Config struct {
	Port       int      `yaml:"port"`
	MarvelHost string   `yaml:"marvelHost"`
	PublicKey  string   `yaml:"publicKey"`
	PrivateKey string   `yaml:"privateKey"`
	DB         Postgres `yaml:"db"`
	Caching    Redis    `yaml:"caching"`
}

func Read(filename string) (Config, error) {

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, errors.Wrapf(err, "read config file: %s", filename)
	}

	var config Config

	if err := yaml.Unmarshal(b, &config); err != nil {
		return Config{}, errors.Wrapf(err, "unmarshal config file: %s", filename)
	}

	return config, nil
}
