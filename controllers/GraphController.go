package controllers

import (
	"github.com/gin-gonic/gin"
	"graph/models"
	"graph/services"
	"log"
	"net/http"
)

type IGraphController interface {
	GraphQuery()
	NodeAttributeSearch()
}

type GraphController struct {
	GraphService services.IGraphService
}

func (graphController *GraphController) GraphQuery(c *gin.Context) {
	var request models.GraphParam
	if err := c.BindJSON(&request); err != nil {
		log.Println("Invalid json request: %s", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	graphData, validationResponse := graphController.GraphService.GraphSearch(request)
	if validationResponse.Validity == false {
		c.IndentedJSON(http.StatusBadRequest, validationResponse)
	} else {
		c.IndentedJSON(http.StatusOK, graphData)
	}
}

func (graphController *GraphController) NodeAttributeSearch(c *gin.Context) {
	var request models.NodeAttributeQueryParam
	if err := c.BindJSON(&request); err != nil {
		log.Println("Invalid json request: %s", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := graphController.GraphService.NodeAttributeSearch(request)
	c.IndentedJSON(http.StatusOK, result)
}
