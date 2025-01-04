package service

import (
	"graph/models"
)

type BidirectionalQuery func(nodeQuery models.NodeQueryModel) models.QueryResultModel

type UnidirectionalQuery func(nodeQuery models.NodeQueryModel) models.QueryResultModel

type IGraphStore interface {
	BFS(seedQueries []models.NodeQueryModel, bidirectionalQuery BidirectionalQuery, unidirectionalQuery UnidirectionalQuery) models.GraphData
}

type GraphStore struct{}

func (graphStore GraphStore) BFS(seedQueries []models.NodeQueryModel, bidirectionalQuery BidirectionalQuery, unidirectionalQuery UnidirectionalQuery) models.GraphData {
	var queryQueue models.NodeQueries
	setOfVisitedNodes := make(map[models.VisitedNodeModel]struct{})
	graphData := models.GraphData{}
	for _, query := range seedQueries {
		queryQueue = append(queryQueue, query)
	}
	for len(queryQueue) > 0 {
		queryNode, queryQueue := queryQueue[0], queryQueue[1:]
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
			queryResult = bidirectionalQuery(queryNode)
		} else if queryNode.CommonFieldName == "" && !ok {
			unidirectionalQuery(queryNode)
		} else {
			queryResult = unidirectionalQuery(queryNode)
		}
		queryQueue.AddQueryOnlyIfUnique(queryResult.NodeQueries)
		graphData.Nodes.AddNodeOnlyIfUnique(queryResult.Nodes)
		graphData.Edges.AddEdgeOnlyIfUnique(queryResult.Edges)
	}
	return graphData
}
