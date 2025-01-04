package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"graph/dao"
	"graph/models"
	"net/http"
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
	queryValidity := graphService.seedQueryService.ValidateQuery(query.Queries, query.Datasource)
	if !queryValidity.Validity {
		return models.GraphData{}, models.ValidationResponse{
			Validity:     false,
			ErrorMessage: models.ERR1,
		}
	}
	constraintsValidity := graphService.seedQueryService.ValidateQuery(query.Constraints, query.Datasource)
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
	
}
