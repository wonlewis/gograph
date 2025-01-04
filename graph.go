package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/gin-gonic/gin"
	"graph/models"
	"log"
	"net/http"
	"os"
	"strings"
)

// Env Credit on using elasticsearch client as a struct: https://github.com/gin-gonic/gin/issues/932
type Env struct {
	db      *elasticsearch.Client
	dbTyped *elasticsearch.TypedClient
}

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
	env := &Env{db: es, dbTyped: esTyped}
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
func (e *Env) ValidateQuery(c *gin.Context) {
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

func (e *Env) BidirectionalQuery(nodeQuery models.NodeQueryModel) (result models.QueryResultModel) {
	var boolQuery *types.BoolQuery
	if nodeQuery.HopLeft == 0 {
		return models.QueryResultModel{
			Nodes:       make([]models.NodeModel, 0),
			Edges:       make([]models.EdgeModel, 0),
			NodeQueries: make([]models.NodeQueryModel, 0),
		}
	} else {
		boolQuery = BoolQueryForBidirectional(nodeQuery.Value, nodeQuery.FromField, nodeQuery.ToField)
		if len(nodeQuery.Constraints) != 0 {
			boolQuery.Filter = nodeQuery.Constraints
		}
	}
	aggregationFrom := AggregationTerms(nodeQuery.Value, nodeQuery.FromField, nodeQuery.NumberOfNeighbours)
	aggregationTo := AggregationTerms(nodeQuery.Value, nodeQuery.ToField, nodeQuery.NumberOfNeighbours)
	for k, v := range aggregationFrom {
		aggregationTo[k] = v
	}
	var size *int
	size = new(int)
	*size = 1
	res, err := e.dbTyped.Search().
		Index(nodeQuery.Datasource).
		Request(&search.Request{
			Query: &types.Query{
				Bool: boolQuery,
			},
			Size:         size,
			Aggregations: aggregationTo,
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
	aggregations := r["aggregations"].(map[string]interface{})
	fromObjects := aggregations[nodeQuery.FromField].(map[string]interface{})["buckets"].([]interface{})
	toObjects := aggregations[nodeQuery.ToField].(map[string]interface{})["buckets"].([]interface{})
	allObjects := append(fromObjects, toObjects...)
	var graphNodes []models.NodeModel
	for _, v := range allObjects {
		graphNodes = append(graphNodes, models.NodeModel{
			FieldName:  nodeQuery.CommonFieldName,
			FieldValue: v.(map[string]interface{})["key"].(string),
			Datasource: nodeQuery.Datasource,
		})
	}
	if len(graphNodes) > 0 {
		graphNodes = append(graphNodes, models.NodeModel{
			FieldName:  nodeQuery.CommonFieldName,
			FieldValue: nodeQuery.Value,
			Datasource: nodeQuery.Datasource,
		})
	}
	var graphQueries []models.NodeQueryModel
	for _, v := range allObjects {
		graphQueries = append(graphQueries, models.NodeQueryModel{
			FromField:          nodeQuery.FromField,
			ToField:            nodeQuery.ToField,
			Value:              v.(map[string]interface{})["key"].(string),
			Constraints:        nodeQuery.Constraints,
			Datasource:         nodeQuery.Datasource,
			NumberOfNeighbours: nodeQuery.NumberOfNeighbours,
			QuerySize:          nodeQuery.QuerySize,
			HopLeft:            nodeQuery.HopLeft - 1,
			CommonFieldName:    nodeQuery.CommonFieldName,
			Reverse:            false,
		})
	}
	var graphEdges []models.EdgeModel
	for _, v := range fromObjects {
		graphEdges = append(graphEdges, models.EdgeModel{
			ToFieldName:    nodeQuery.CommonFieldName,
			ToFieldValue:   nodeQuery.Value,
			FromFieldName:  nodeQuery.CommonFieldName,
			FromFieldValue: v.(map[string]interface{})["key"].(string),
			Datasource:     nodeQuery.Datasource,
			Frequency:      int(v.(map[string]interface{})["doc_count"].(float64)),
		})
	}
	for _, v := range toObjects {
		graphEdges = append(graphEdges, models.EdgeModel{
			ToFieldName:    nodeQuery.CommonFieldName,
			ToFieldValue:   v.(map[string]interface{})["key"].(string),
			FromFieldName:  nodeQuery.CommonFieldName,
			FromFieldValue: nodeQuery.Value,
			Datasource:     nodeQuery.Datasource,
			Frequency:      int(v.(map[string]interface{})["doc_count"].(float64)),
		})
	}
	return models.QueryResultModel{
		Nodes:       graphNodes,
		Edges:       graphEdges,
		NodeQueries: graphQueries,
	}
}

func BoolQueryForBidirectional(value string, fromField string, toField string) *types.BoolQuery {
	minimumShouldMatch := new(types.MinimumShouldMatch)
	*minimumShouldMatch = 1
	boolQuery := types.NewBoolQuery()
	boolQuery.Should = []types.Query{
		types.Query{
			Match: map[string]types.MatchQuery{
				fromField: {Query: value},
			},
		},
		types.Query{
			Match: map[string]types.MatchQuery{
				toField: {Query: value},
			},
		},
	}
	boolQuery.MinimumShouldMatch = minimumShouldMatch
	return boolQuery
}

func AggregationTerms(value string, field string, numberOfNeighbours int) map[string]types.Aggregations {
	aggregations := make(map[string]types.Aggregations)
	aggregationQuery := types.Aggregations{
		Terms: &types.TermsAggregation{
			Field:   &field,
			Exclude: []string{value, ""},
			Size:    &numberOfNeighbours,
		},
	}
	aggregations[field] = aggregationQuery
	return aggregations
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
