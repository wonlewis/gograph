package main

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"graph/models"
	"log"
	"os"
)

type queryBody struct {
	Index string `json:"index"`
}

func main() {
	mainLogger := log.New(os.Stdout, "main:", log.LstdFlags)
	es, err := GetElasticsearchClient()
	if err != nil {
		mainLogger.Println("Error getting response: %s", err)
	}
	esTyped, err := GetElasticsearchTypedClient()
	if err != nil {
		mainLogger.Println("Error getting response from typed client: %s", err)
	}
	log.Println(elasticsearch.Version)
	log.Println(es.Info())
	env := &models.Env{Db: es, DbTyped: esTyped}
	router := gin.Default()
	router.POST("/data", env.KeywordSearch)
	router.POST("/dataTyped", env.KeywordSearchTyped)
	router.POST("/wrapper", env.WrapperQuery)
	router.POST("/validate", env.ValidateQuery)
	err = router.Run("localhost:8080")
	if err != nil {
		return
	}
}

func GetElasticsearchClient() (*elasticsearch.Client, error) {
	//username := os.Getenv("ELASTIC_USERNAME")
	//password := os.Getenv("ELASTIC_PASSWORD")
	cert, _ := os.ReadFile("./http_ca.crt")
	cfg := elasticsearch.Config{
		Addresses: []string{
			"https://localhost:9200",
		},
		Username: "elastic",
		Password: "-uephdfG+o_9H=9K0D10",
		CACert:   cert,
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to connect: %w", err)
	}

	return es, nil
}

func GetElasticsearchTypedClient() (*elasticsearch.TypedClient, error) {
	//username := os.Getenv("ELASTIC_USERNAME")
	//password := os.Getenv("ELASTIC_PASSWORD")
	cert, _ := os.ReadFile("./http_ca.crt")
	cfg := elasticsearch.Config{
		Addresses: []string{
			"https://localhost:9200",
		},
		Username: "elastic",
		Password: "-uephdfG+o_9H=9K0D10",
		CACert:   cert,
	}
	es, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to connect: %w", err)
	}

	return es, nil
}
