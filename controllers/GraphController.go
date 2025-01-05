package controllers

import (
	"graph/models"
)

type IGraphService interface {
	GraphSearch(query models.GraphParam) models.GraphData
	NodeAttributeSearch(query models.NodeAttributeQueryParam) []interface{}
}

type IGraphController interface {
	GraphQuery(graphParam models.GraphParam) models.GraphData
	NodeAttributeQuery(queryParam models.NodeAttributeQueryParam) []interface{}
}

type GraphController struct {
	graphService IGraphService
}

func (graphController *GraphController) GraphQuery(graphParam models.GraphParam) models.GraphData {
	return graphController.graphService.GraphSearch(graphParam)
}

func (graphController *GraphController) NodeAttributeQuery(graphParam models.NodeAttributeQueryParam) []interface{} {
	return graphController.graphService.NodeAttributeSearch(graphParam)
}
