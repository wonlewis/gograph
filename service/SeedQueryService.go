package service

import (
	"graph/models"
)

type ISeedQueryService interface {
	ValidateQueries(query []string, datasource string) models.ValidationResponse
	ValidateFields(fields []models.FieldModel, datasource string) models.ValidationResponse
	GetSeedQueries(query models.GraphParam) []models.NodeQueryModel
}
