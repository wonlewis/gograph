package main

import (
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"os"
)

func main() {
	cert, _ := os.ReadFile("./http_ca.crt")
	cfg := elasticsearch.Config{
		Addresses: []string{
			"https://localhost:9200",
		},
		Username: "elastic",
		Password: "-gN-NuCy+qmr5IowL_X_",
		CACert:   cert,
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	log.Println(elasticsearch.Version)
	log.Println(es.Info())

}
