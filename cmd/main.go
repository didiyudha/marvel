package main

import (
	"github.com/pkg/errors"
	"github.com/didiyudha/marvel/business/data/character"
	"github.com/didiyudha/marvel/client"
	"github.com/didiyudha/marvel/cmd/exec"
	"github.com/didiyudha/marvel/config"
	"github.com/didiyudha/marvel/foundation/database"
	"log"
	"os"
	"github.com/go-resty/resty/v2"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	cmd := os.Args[1]
	configFile := os.Getenv("MARVEL_CONFIG")

	conf, err := config.Read(configFile)
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
	defer db.Close()

	cred := client.CredentialRequest{
		Host:       conf.MarvelHost,
		PublicKey:  conf.PublicKey,
		PrivateKey: conf.PrivateKey,
	}

	store := character.NewStore(db)
	restyClient := resty.New()
	marvelClient := client.NewMarvelClient(cred, restyClient)

	tableMigrator := exec.NewTableMigrator(db)

	marvelCmdExecutor := exec.NewMarvelCmdExecutor(store, marvelClient)
	cmdExecutor := exec.NewCommandExecutor(marvelCmdExecutor, tableMigrator)

	return cmdExecutor.Exec(cmd)
}