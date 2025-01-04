package dao

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"graph/models"
	"log"
)

// Env Credit on using elasticsearch client as a struct: https://github.com/gin-gonic/gin/issues/932
type Env struct {
	Db      *elasticsearch.Client
	DbTyped *elasticsearch.TypedClient
}

type IGraphQueryDAO interface{}

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
	res, err := e.DbTyped.Search().
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
		if nodeQuery.HopLeft > 1 {
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
