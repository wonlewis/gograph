package dao

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/validatequery"
	"github.com/elastic/go-elasticsearch/v8/typedapi/ml/validate"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/gin-gonic/gin"
	"graph/models"
	"log"
	"net/http"
	"strings"
)

type ElasticSeedQueryDAO struct {
	Db *elasticsearch.TypedClient
}

type IElasticSeedQueryDAO interface {
	ValidateQueries(queries []string, datasource string) models.ValidationResponse
	ValidateFields(fields []models.FieldModel, datasource string) models.ValidationResponse
	GetSeedQueries(query models.GraphParam) []models.NodeQueryModel
}

func (e ElasticSeedQueryDAO) ValidateQueries(queries []string, datasource string) models.ValidationResponse {
	if len(queries) == 0 {
		return models.ValidationResponse{
			Validity:     true,
			InvalidField: "",
			ErrorMessage: models.OK,
		}
	}
	var ListOfQueries []types.Query
	for _, query := range queries {
		encodedQuery := base64.StdEncoding.EncodeToString([]byte(query))
		wrapperQuery := &types.WrapperQuery{
			Query: encodedQuery,
		}
		queryWrapped := &types.Query{
			Wrapper: wrapperQuery,
		}
		ListOfQueries = append(ListOfQueries, *queryWrapped)
	}
	boolQuery := &types.BoolQuery{
		Must: ListOfQueries,
	}
	res, err := e.Db.Indices.ValidateQuery().
		Index(datasource).
		Request(&validatequery.Request{
			Query: &types.Query{
				Bool: boolQuery,
			},
		}).Do(context.Background())
	if err != nil {
		log.Println("Error getting response: %s", err)
		return models.ValidationResponse{
			Validity:     false,
			InvalidField: "",
			ErrorMessage: models.ERR4,
		}
	}
	var r map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		log.Println("Error decoding response: %s", err)
		return models.ValidationResponse{
			Validity:     false,
			InvalidField: "",
			ErrorMessage: models.ERR1,
		}
	}
	return models.ValidationResponse{
		Validity:     true,
		InvalidField: "",
		ErrorMessage: models.OK,
	}
}

func (e *Env) KeywordSearchTyped(c *gin.Context) {
	var request models.GraphParam
	if err := c.BindJSON(&request); err != nil {
		log.Println("Invalid json request: %s", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := e.dbTyped.Search().
		Index(request.Datasource).
		Request(&search.Request{
			Query: &types.Query{
				Match: map[string]types.MatchQuery{
					"sender": {Query: "tom"},
				},
			},
		}).Do(context.Background())
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

func (e *Env) WrapperQuery(c *gin.Context) {
	var request models.GraphParam
	if err := c.BindJSON(&request); err != nil {
		log.Println("Invalid json request: %s", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	query := `{
	  "bool": {
		  "filter": [
			{
			  "match": {
				"sender": "tom"
			  }
			}
		  ]
		}
	}`
	encodedString := base64.StdEncoding.EncodeToString([]byte(query))
	res, err := e.dbTyped.Search().
		Index(request.Datasource).
		Request(&search.Request{
			Query: &types.Query{
				Wrapper: &types.WrapperQuery{
					Query: encodedString,
				},
			},
		}).Do(context.Background())
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
