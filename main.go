package main

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"graph/controllers"
	"graph/dao"
	"graph/services"
	"log"
	"os"
)

type queryBody struct {
	Index string `json:"index"`
}

func main() {
	// setting up logger and database
	mainLogger := log.New(os.Stdout, "main:", log.LstdFlags)
	es, err := GetElasticsearchTypedClient()
	if err != nil {
		mainLogger.Println("Error getting response from elastic client: %s", err)
	}
	log.Println(elasticsearch.Version)
	log.Println(es.Info())

	// setting up controller, services, and DAO layers
	var graphController controllers.GraphController
	graphController = *SetupController(es)

	router := gin.Default()
	router.POST("/searchGraph", graphController.GraphQuery)
	router.POST("/searchNode", graphController.NodeAttributeSearch)
	err = router.Run("localhost:8080")
	if err != nil {
		return
	}
}

func SetupController(es *elasticsearch.TypedClient) *controllers.GraphController {
	// setting up controller, services, and DAO layers
	var graphQueryDAO dao.ElasticGraphQueryDAO
	graphQueryDAO.Db = es
	var seedQueryDAO dao.ElasticSeedQueryDAO
	seedQueryDAO.Db = es
	var seedQueryService services.SeedQueryService
	seedQueryService.SeedQueryDAO = &seedQueryDAO
	var graphService services.GraphService
	graphService.SeedQueryService = &seedQueryService
	graphService.GraphQueryDAO = &graphQueryDAO
	var graphController controllers.GraphController
	graphController.GraphService = &graphService
	return &graphController
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
