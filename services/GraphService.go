package services

import (
	"fmt"
	"graph/dao"
	"graph/models"
)

type IGraphService interface {
	GraphSearch(query models.GraphParam) (models.GraphData, models.ValidationResponse)
	NodeAttributeSearch(query models.NodeAttributeQueryParam) map[string]interface{}
}

type GraphService struct {
	SeedQueryService ISeedQueryService
	GraphQueryDAO    dao.IGraphQueryDAO
}

func (graphService *GraphService) GraphSearch(query models.GraphParam) (models.GraphData, models.ValidationResponse) {
	queryValidity := graphService.SeedQueryService.ValidateQueries(query, query.Datasource)
	if !queryValidity.Validity {
		return models.GraphData{}, models.ValidationResponse{
			Validity:     false,
			ErrorMessage: models.ERR1,
		}
	}
	constraintsValidity := graphService.SeedQueryService.ValidateQueries(query, query.Datasource)
	if !constraintsValidity.Validity {
		return models.GraphData{}, models.ValidationResponse{
			Validity:     false,
			ErrorMessage: models.ERR2,
		}
	}
	fieldsValidity := graphService.SeedQueryService.ValidateFields(query.Vertices, query.Datasource)
	if !fieldsValidity.Validity {
		return models.GraphData{}, models.ValidationResponse{
			Validity:     false,
			ErrorMessage: models.ValidationStatus(fieldsValidity.InvalidField),
		}
	}
	nodeQueries, err := graphService.SeedQueryService.GetSeedQueries(query)
	if err != nil {
		return models.GraphData{}, models.ValidationResponse{
			Validity:     false,
			ErrorMessage: models.ValidationStatus(fmt.Sprintf("%s", err)),
		}
	}
	var graphStore = GraphStore{}
	return graphStore.BFS(nodeQueries, graphService.GraphQueryDAO.BidirectionalQuery, graphService.GraphQueryDAO.UnidirectionalQuery), models.ValidationResponse{}
}

func (graphService *GraphService) NodeAttributeSearch(query models.NodeAttributeQueryParam) map[string]interface{} {
	return graphService.GraphQueryDAO.NodeAttributeQuery(query)
}
