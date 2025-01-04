package dao

import (
	"encoding/base64"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/gin-gonic/gin"
	"graph/models"
	"log"
	"net/http"
	"strings"
)

func (e *models.Env) ValidateQuery(c *gin.Context) {
	var request models.GraphParam
	if err := c.BindJSON(&request); err != nil {
		log.Println("Invalid json request: %s", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	parsedQuery := strings.Join([]string(request.Queries), ",")
	boolQuery := fmt.Sprintf("{\"query\": {\"bool\": {\"filter\":[%s]}}}", parsedQuery)
	res, err := e.db.Indices.ValidateQuery(
		e.db.Indices.ValidateQuery.WithQuery(boolQuery),
		e.db.Indices.ValidateQuery.WithIndex(request.Datasource),
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
	nodeQuery := models.NodeQueryModel{
		FromField:          request.Vertices[0].FromField,
		ToField:            request.Vertices[0].ToField,
		Value:              "tom",
		Constraints:        nil,
		Datasource:         request.Datasource,
		NumberOfNeighbours: request.NumberOfNeighbours,
		QuerySize:          request.DocCount,
		HopLeft:            request.Hop,
		CommonFieldName:    request.Vertices[0].CommonFieldName,
		Reverse:            false,
	}
	e.BidirectionalQuery(nodeQuery)
	c.IndentedJSON(http.StatusOK, r)
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
