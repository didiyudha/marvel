package main

import (
	"github.com/didiyudha/marvel/business/api"
	"github.com/didiyudha/marvel/foundation/caching"
	"github.com/didiyudha/marvel/foundation/database"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	var config = struct {
		Port       int    `yaml:"port"`
		PublicKey  string `yaml:"publicKey"`
		PrivateKey string `yaml:"privateKey"`
		DB         struct {
			User         string `yaml:"user"`
			Password     string `yaml:"password"`
			Host         string `yaml:"host"`
			Name         string `yaml:"name"`
			MaxIdleConns int    `yaml:"maxIdleConns"`
			MaxOpenConns int    `yaml:"maxOpenConns"`
			DisableTLS   bool   `yaml:"disableTLS"`
		} `yaml:"db"`
		Caching struct {
			Addr     string `yaml:"addr"`
			Password string `yaml:"password"`
			DB       int    `yaml:"db"`
		} `yaml:"caching"`
	}{}

	b, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(b, &config); err != nil {
		return errors.Wrap(err, "unmarshal config file")
	}

	dbConfig := database.Config{
		User:         config.DB.User,
		Password:     config.DB.Password,
		Host:         config.DB.Host,
		Name:         config.DB.Name,
		MaxIdleConns: config.DB.MaxIdleConns,
		MaxOpenConns: config.DB.MaxOpenConns,
		DisableTLS:   config.DB.DisableTLS,
	}

	db, err := database.Open(dbConfig)
	if err != nil {
		return errors.Wrap(err, "open connection to database")
	}

	cachingConfig := caching.Config{
		Addr:     config.Caching.Addr,
		Password: config.Caching.Password,
		DB:       config.Caching.DB,
	}

	redisClient, err := caching.New(cachingConfig)
	if err != nil {
		return errors.Wrap(err, "create new redis client")
	}

	API := api.New(db, redisClient)

	return API.Serve(config.Port)
}
