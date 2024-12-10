package main

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"os"
)

func main() {
	es, err := GetElasticsearchClient()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	log.Println(elasticsearch.Version)
	log.Println(es.Info())

}

func GetElasticsearchClient() (*elasticsearch.TypedClient, error) {
	username := os.Getenv("ELASTIC_USERNAME")
	password := os.Getenv("ELASTIC_PASSWORD")
	cert, _ := os.ReadFile("./http_ca.crt")
	cfg := elasticsearch.Config{
		Addresses: []string{
			"https://localhost:9200",
		},
		Username: username,
		Password: password,
		CACert:   cert,
	}
	es, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to connect: %w", err)
	}

	return es, nil
}
