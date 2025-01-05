package services

import (
	"graph/dao"
	"graph/models"
)

type ISeedQueryService interface {
	ValidateQueries(query []string, datasource string) models.ValidationResponse
	ValidateFields(fields []models.FieldModel, datasource string) models.ValidationResponse
	GetSeedQueries(query models.GraphParam) []models.NodeQueryModel
}

type SeedQueryService struct {
	dao.ISeedQueryDAO
}

func (s *SeedQueryService) ValidateQueries(query models.GraphParam, datasource string) models.ValidationResponse {
	return s.ISeedQueryDAO.ValidateQueries(query.Queries, datasource)
}

func (s *SeedQueryService) ValidateFields(fields []models.FieldModel, datasource string) models.ValidationResponse {
	return s.ISeedQueryDAO.ValidateFields(fields, datasource)
}

func (s *SeedQueryService) GetSeedQueries(query models.GraphParam) []models.NodeQueryModel {
	return s.ISeedQueryDAO.GetSeedQueries(query)
}
