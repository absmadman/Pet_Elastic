package main

import (
	"Pet_Elastic/config"
	esfunc "Pet_Elastic/elasticsearch"
	"Pet_Elastic/helpers"
	"Pet_Elastic/server"
	"github.com/elastic/go-elasticsearch/v8"
)

func main() {

	es, err := elasticsearch.NewClient(*config.NewCfg())

	helpers.Check(err)

	esfunc.DefaultInit(es)

	server.LaunchHttpServer(es)

}
