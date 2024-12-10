package main

import (
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
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
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	log.Println(elasticsearch.Version)
	log.Println(es.Info())
	env := &Env{db: es}
	router := gin.Default()
	router.POST("/data", env.KeywordSearch)
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
		Password: "-gN-NuCy+qmr5IowL_X_",
		CACert:   cert,
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to connect: %w", err)
	}

	return es, nil
}

func (e *Env) KeywordSearch(c *gin.Context) {
	var request queryBody
	if err := c.BindJSON(&request); err != nil {
		return
	}
	query := `{
	  "query": {
		"match": {
		  "sender": {
			"query": "tom"
		  }
		}
	  }
	}`
	res, err := e.db.Search(
		e.db.Search.WithIndex(request.Index),
		e.db.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		return
	}
	var r map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return
	}
	c.IndentedJSON(http.StatusOK, r)
}
