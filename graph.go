package main

import (
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"graph/models"
	"log"
	"net/http"
	"os"
	"strings"
)

// Env Credit on using elasticsearch client as a struct: https://github.com/gin-gonic/gin/issues/932
type Env struct {
	db *elasticsearch.Client
}

type queryBody struct {
	Index string `json:"index"`
}

func main() {
	es, err := GetElasticsearchClient()
	mainLogger := log.New(os.Stdout, "main:", log.LstdFlags)
	if err != nil {
		mainLogger.Println("Error getting response: %s", err)
	}
	log.Println(elasticsearch.Version)
	log.Println(es.Info())
	env := &Env{db: es}
	router := gin.Default()
	router.POST("/data", env.KeywordSearch)
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

func (e *Env) ValidateQuery(c *gin.Context) {
	var request models.GraphParam
	if err := c.BindJSON(&request); err != nil {
		log.Println("Invalid json request: %s", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	parsedQuery := strings.Join([]string(request.Queries), ",")
	boolQuery := fmt.Sprintf("{\"query\": {\"bool\": {\"filter\":[%s]}}}", parsedQuery)
	validateQuery := e.db.Indices.ValidateQuery
	validateQuery.WithIndex(request.Datasource)
	validateQuery.WithQuery(boolQuery)
	res, err := e.db.Indices.ValidateQuery(
		e.db.Indices.ValidateQuery.WithQuery(boolQuery),
		e.db.Indices.ValidateQuery.WithIndex(boolQuery),
	)
	if err != nil {
		log.Println("Error getting response: %s", err)
		return
	}
	var r map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		log.Println("Error decoding response: %s", err)
		return
	}
	c.IndentedJSON(http.StatusOK, r)
}

func (e *Env) KeywordSearch(c *gin.Context) {
	var request models.GraphParam
	if err := c.BindJSON(&request); err != nil {
		log.Println("Invalid json request: %s", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	query := `{
	  "query": {
		"bool": {
		  "filter": [
			{
			  "match": {
				"sender": "tom"
			  }
			}
		  ]
		}
	  }
	}`
	res, err := e.db.Search(
		e.db.Search.WithIndex(request.Datasource),
		e.db.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		log.Println("Error getting response: %s", err)
		return
	}
	var r map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		log.Println("Error decoding response: %s", err)
		return
	}
	c.IndentedJSON(http.StatusOK, r)
}
