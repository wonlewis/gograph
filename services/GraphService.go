package services

import (
	"github.com/gin-gonic/gin"
	"graph/dao"
	"graph/models"
)

type IGraphService interface {
	GraphSearch(query models.GraphParam) models.GraphData
	NodeAttributeSearch(query models.NodeAttributeQueryParam) []interface{}
}

type GraphService struct {
	seedQueryService ISeedQueryService
	graphQueryDAO    dao.IGraphQueryDAO
}

func (graphService *GraphService) GraphSearch(query models.GraphParam, c *gin.Context) (models.GraphData, models.ValidationResponse) {
	queryValidity := graphService.seedQueryService.ValidateQueries(query.Queries, query.Datasource)
	if !queryValidity.Validity {
		return models.GraphData{}, models.ValidationResponse{
			Validity:     false,
			ErrorMessage: models.ERR1,
		}
	}
	constraintsValidity := graphService.seedQueryService.ValidateQueries(query.Constraints, query.Datasource)
	if !constraintsValidity.Validity {
		return models.GraphData{}, models.ValidationResponse{
			Validity:     false,
			ErrorMessage: models.ERR2,
		}
	}
	fieldsValidity := graphService.seedQueryService.ValidateFields(query.Vertices, query.Datasource)
	if !fieldsValidity.Validity {
		return models.GraphData{}, models.ValidationResponse{
			Validity:     false,
			ErrorMessage: models.ValidationStatus(fieldsValidity.InvalidField),
		}
	}
	nodeQueries := graphService.seedQueryService.GetSeedQueries(query)
	graphStore := GraphStore{}
	return graphStore.BFS(nodeQueries, graphService.graphQueryDAO.BidirectionalQuery, graphService.graphQueryDAO.UnidirectionalQuery), models.ValidationResponse{}
}
