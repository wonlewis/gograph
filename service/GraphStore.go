package service

import "graph/models"

type BidirectionalQuery func(nodeQuery models.NodeQueryModel) models.QueryResultModel

type UnidirectionalQuery func(nodeQuery models.NodeQueryModel) models.QueryResultModel

type IGraphStore interface {
	BFS(seedQueries []models.NodeQueryModel, bidirectionalQuery BidirectionalQuery, unidirectionalQuery UnidirectionalQuery) models.GraphData
}
