package main

import (
	"github.com/didiyudha/marvel/business/data/character"
	"github.com/didiyudha/marvel/client"
	"github.com/didiyudha/marvel/config"
	"github.com/didiyudha/marvel/foundation/database"
	"github.com/didiyudha/marvel/worker/internal"
	"github.com/go-resty/resty/v2"
	"log"
)

func main() {

	configFile := "./../config.yaml"
	conf, err := config.Read(configFile)
	if err != nil {
		log.Fatal(err)
	}


	cred := client.CredentialRequest{
		Host:    conf.MarvelHost,
		PublicKey:  conf.PublicKey,
		PrivateKey: conf.PrivateKey,
	}

	restyClient := resty.New()
	marvelClient := client.NewMarvelClient(cred, restyClient)

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
		log.Fatal(err)
	}
	store := character.NewStore(db)


	worker := internal.NewWorker()
	marvelWorker := internal.NewMarvelWorker(store, marvelClient)
	worker.AddFunc(marvelWorker)

	worker.Run()
}
