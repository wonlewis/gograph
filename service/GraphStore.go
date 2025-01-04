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

func (graphStore GraphStore) BFS(seedQueries []models.NodeQueryModel, bidirectionalQuery BidirectionalQuery, unidirectionalQuery UnidirectionalQuery) models.GraphData {
	var queryQueue deque.Deque[models.NodeQueryModel]
	setOfVisitedNodes := make(map[models.VisitedNodeModel]struct{})
	graphData := models.GraphData{}
	for _, query := range seedQueries {
		queryQueue.PushBack(query)
	}
	for queryQueue.Len() > 0 {
		queryNode := queryQueue.PopFront()
		queryResult := models.QueryResultModel{}
		queryNodeToNodeModel := models.VisitedNodeModel{
			FromField:       queryNode.FromField,
			ToField:         queryNode.ToField,
			Value:           queryNode.Value,
			Datasource:      queryNode.Datasource,
			CommonFieldName: queryNode.CommonFieldName,
		}
		_, ok := setOfVisitedNodes[queryNodeToNodeModel]
		if queryNode.CommonFieldName != "" && !ok {
			bidirectionalQuery(queryNode)
		} else if queryNode.CommonFieldName == "" && !ok {
			unidirectionalQuery(queryNode)
		} else {
			queryResult = unidirectionalQuery(queryNode)
		}
		graphData.Nodes
	}
}

func (models.GraphNodes) AddNodeOnlyIfUnique(graphNodes []models.NodeModel) {

}
