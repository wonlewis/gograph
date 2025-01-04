package models

import "github.com/elastic/go-elasticsearch/v8/typedapi/types"

type NodeModel struct {
	FieldName  string `json:"field_name"`
	FieldValue string `json:"field_value"`
	Datasource string `json:"datasource"`
}

type EdgeModel struct {
	ToFieldName    string `json:"to_field_name"`
	ToFieldValue   string `json:"to_field_value"`
	FromFieldName  string `json:"from_field_name"`
	FromFieldValue string `json:"from_field_value"`
	Datasource     string `json:"datasource"`
	Frequency      int    `json:"frequency"`
}

type FieldModel struct {
	FromField       string `json:"from_field"`
	ToField         string `json:"to_field"`
	CommonFieldName string `json:"common_field_name"`
}

type GraphParam struct {
	Queries            []string     `json:"queries" binding:"required"`
	Constraints        []string     `json:"constraints" binding:"required"`
	Hop                int          `json:"hop" binding:"required"`
	DocCount           int          `json:"doc_count" binding:"required"`
	NumberOfNeighbours int          `json:"number_of_neighbours" binding:"required"`
	Datasource         string       `json:"datasource" binding:"required"`
	Vertices           []FieldModel `json:"vertices" binding:"required"`
}

type NodeQueryModel struct {
	FromField          string        `json:"from_field"`
	ToField            string        `json:"to_field"`
	Value              string        `json:"value"`
	Constraints        []types.Query `json:"constraints"`
	Datasource         string        `json:"datasource" binding:"required"`
	NumberOfNeighbours int           `json:"number_of_neighbours"`
	QuerySize          int           `json:"query_size"`
	HopLeft            int           `json:"hop_left"`
	CommonFieldName    string        `json:"common_field_name"`
	Reverse            bool          `json:"reverse"`
}

type ValidationResponse struct {
	Validity     bool   `json:"validity"`
	InvalidField string `json:"invalid_field"`
	ErrorMessage ValidationStatus
}

type GraphData struct {
	Nodes []NodeModel `json:"nodes"`
	Edges []EdgeModel `json:"edges"`
}

type QueryResultModel struct {
	Nodes       []NodeModel      `json:"nodes"`
	Edges       []EdgeModel      `json:"edges"`
	NodeQueries []NodeQueryModel `json:"node_queries"`
}

type NodeAttributeQueryParam struct {
	Value      string `json:"value" binding:"required"`
	FieldName  string `json:"field_name" binding:"required"`
	Datasource string `json:"datasource" binding:"required"`
}
