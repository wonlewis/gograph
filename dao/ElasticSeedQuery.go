package dao

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/validatequery"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"graph/models"
	"log"
	"strings"
)

type ElasticSeedQueryDAO struct {
	Db *elasticsearch.TypedClient
}

type ISeedQueryDAO interface {
	ValidateQueries(queries []string, datasource string) models.ValidationResponse
	ValidateFields(fields []models.FieldModel, datasource string) models.ValidationResponse
	GetSeedQueries(query models.GraphParam) []models.NodeQueryModel
}

func (e *ElasticSeedQueryDAO) ValidateQueries(queries []string, datasource string) models.ValidationResponse {
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

func (e *ElasticSeedQueryDAO) ValidateFields(fields []models.FieldModel, datasource string) models.ValidationResponse {
	var setOfFieldsToCheck map[string]struct{}
	for _, field := range fields {
		_, ok := setOfFieldsToCheck[field.FromField]
		if !ok {
			setOfFieldsToCheck[field.FromField] = struct{}{}
		}
		_, ok = setOfFieldsToCheck[field.ToField]
		if !ok {
			setOfFieldsToCheck[field.ToField] = struct{}{}
		}
	}
	listOfFieldsToCheck := make([]string, 0, len(setOfFieldsToCheck))
	stringOfFieldsToCheck := strings.Join(listOfFieldsToCheck, ",")
	for key, _ := range setOfFieldsToCheck {
		listOfFieldsToCheck = append(listOfFieldsToCheck, key)
	}
	res, err := e.Db.Indices.GetFieldMapping(stringOfFieldsToCheck).
		Index(datasource).
		Do(context.Background())
	if err != nil {
		return models.ValidationResponse{
			Validity:     false,
			InvalidField: "",
			ErrorMessage: models.ERR4,
		}
	}
	var r map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		log.Printf("Error decoding response: %s\n", err)
		return models.ValidationResponse{
			Validity:     false,
			InvalidField: "",
			ErrorMessage: models.ERR4,
		}
	}
	return models.ValidationResponse{
		Validity:     true,
		InvalidField: "",
		ErrorMessage: models.OK,
	}
}

func (e *ElasticSeedQueryDAO) GetSeedQueries(query models.GraphParam) ([]models.NodeQueryModel, error) {
	var ListOfQueries []types.Query
	if len(query.Queries) > 0 {
		for _, query := range query.Queries {
			encodedQuery := base64.StdEncoding.EncodeToString([]byte(query))
			wrapperQuery := &types.WrapperQuery{
				Query: encodedQuery,
			}
			queryWrapped := &types.Query{
				Wrapper: wrapperQuery,
			}
			ListOfQueries = append(ListOfQueries, *queryWrapped)
		}
	}
	var listOfConstraints []types.Query
	if len(query.Constraints) > 0 {
		for _, query := range query.Constraints {
			encodedQuery := base64.StdEncoding.EncodeToString([]byte(query))
			wrapperQuery := &types.WrapperQuery{
				Query: encodedQuery,
			}
			queryWrapped := &types.Query{
				Wrapper: wrapperQuery,
			}
			listOfConstraints = append(listOfConstraints, *queryWrapped)
		}
	}
	boolQuery := &types.BoolQuery{
		Should: ListOfQueries,
		Filter: listOfConstraints,
	}
	size := new(int)
	*size = query.DocCount
	res, err := e.Db.Search().
		Index(query.Datasource).
		Request(&search.Request{
			Query: &types.Query{
				Bool: boolQuery,
			},
			Size: size,
		}).Do(context.Background())
	if err != nil {
		log.Println("Error getting response: %s", err)
		return nil, err
	}
	var r map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&r)
	return nil, nil
}
