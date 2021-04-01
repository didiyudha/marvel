package main

import (
	"github.com/didiyudha/marvel/business/data/character"
	"github.com/didiyudha/marvel/client"
	"github.com/didiyudha/marvel/cmd/exec"
	"github.com/didiyudha/marvel/config"
	"github.com/didiyudha/marvel/foundation/database"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"log"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cmd := os.Args[1]

	conf, err := config.Read("./../config.yaml")
	if err != nil {
		return err
	}

	dbConfig := database.Config{
		User:         conf.DB.User,
		Password:     conf.DB.Password,
		Host:         conf.DB.Host,
		Name:         conf.DB.Name,
		MaxIdleConns: conf.DB.MaxIdleConns,
		MaxOpenConns: conf.DB.MaxOpenConns,
		DisableTLS:   conf.DB.DisableTLS,
	}

	db, err := database.Open(dbConfig)
	if err != nil {
		return errors.Wrap(err, "open connection to database")
	}

	cred := client.CredentialRequest{
		Host:       conf.MarvelHost,
		PublicKey:  conf.PublicKey,
		PrivateKey: conf.PrivateKey,
	}

	store := character.NewStore(db)
	restyClient := resty.New()
	marvelClient := client.NewMarvelClient(cred, restyClient)

	marvelCmdExecutor := exec.NewMarvelCmdExecutor(store, marvelClient)
	cmdExecutor := exec.NewCommandExecutor(marvelCmdExecutor)

	return cmdExecutor.Exec(cmd)
}