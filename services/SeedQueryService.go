package services

import (
	"graph/dao"
	"graph/models"
)

type ISeedQueryService interface {
	ValidateQueries(query models.GraphParam, datasource string) models.ValidationResponse
	ValidateFields(fields []models.FieldModel, datasource string) models.ValidationResponse
	GetSeedQueries(query models.GraphParam) ([]models.NodeQueryModel, error)
}

type SeedQueryService struct {
	SeedQueryDAO dao.ISeedQueryDAO
}

func (s *SeedQueryService) ValidateQueries(query models.GraphParam, datasource string) models.ValidationResponse {
	return s.SeedQueryDAO.ValidateQueries(query.Queries, datasource)
}

func (s *SeedQueryService) ValidateFields(fields []models.FieldModel, datasource string) models.ValidationResponse {
	return s.SeedQueryDAO.ValidateFields(fields, datasource)
}

func (s *SeedQueryService) GetSeedQueries(query models.GraphParam) ([]models.NodeQueryModel, error) {
	return s.SeedQueryDAO.GetSeedQueries(query)
}
