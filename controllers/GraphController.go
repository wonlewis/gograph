package controllers

import (
	"github.com/gin-gonic/gin"
	"graph/models"
	"graph/services"
)

type IGraphController interface {
	GraphQuery()
	NodeAttributeSearch()
}

type GraphController struct {
	GraphService services.IGraphService
}

func (graphController *GraphController) GraphQuery(c *gin.Context) {
	return graphController.GraphService.GraphSearch(graphParam)
}

func (graphController *GraphController) NodeAttributeSearch(c *gin.Context) {
	return graphController.GraphService.NodeAttributeSearch(graphParam)
}
