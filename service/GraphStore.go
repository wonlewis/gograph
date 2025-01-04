package service

import (
	"github.com/gammazero/deque"
	"graph/models"
)

type BidirectionalQuery func(nodeQuery models.NodeQueryModel) models.QueryResultModel

type UnidirectionalQuery func(nodeQuery models.NodeQueryModel) models.QueryResultModel

type IGraphStore interface {
	BFS(seedQueries []models.NodeQueryModel, bidirectionalQuery BidirectionalQuery, unidirectionalQuery UnidirectionalQuery) models.GraphData
}

type GraphStore struct{}

func (graphStore GraphStore) BFS(seedQueries []models.NodeQueryModel, bidirectionalQuery BidirectionalQuery) models.GraphData {
	var queryQueue deque.Deque[models.NodeQueryModel]
	setOfVisitedNodes := make(map[models.NodeModel]struct{})
	graphData := models.GraphData{}
	for _, query := range seedQueries {
		queryQueue.PushBack(query)
	}
	for queryQueue.Len() > 0 {
		queryNode := queryQueue.PopFront()
		queryResult := models.QueryResultModel{}
		queryNodeToNodeModel := models.NodeModel{
			FieldName:  queryNode.CommonFieldName,
			FieldValue: "",
			Datasource: "",
		}
		if queryNode.CommonFieldName != "" && 
	}
}
