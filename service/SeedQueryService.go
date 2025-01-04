package service

import (
	"graph/models"
)

type ISeedQueryService interface {
	ValidateQuery(query []string, datasource string) models.ValidationResponse
	ValidateField(field string, datasource string) models.ValidationResponse
	GetSeedQueries(query models.GraphParam) []models.NodeQueryModel
}
