package models

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
	Queries            []string     `json:"queries"`
	Constraints        []string     `json:"constraints"`
	Hop                int          `json:"hop"`
	DocCount           int          `json:"doc_count"`
	NumberOfNeighbours int          `json:"number_of_neighbours"`
	Datasource         string       `json:"datasource"`
	Vertices           []FieldModel `json:"vertices"`
}

type NodeQueryModel struct {
	FromField          string   `json:"from_field"`
	ToField            string   `json:"to_field"`
	Values             string   `json:"values"`
	Constraints        []string `json:"constraints"`
	Datasource         string   `json:"datasource"`
	NumberOfNeighbours int      `json:"number_of_neighbours"`
	QuerySize          int      `json:"query_size"`
	HopLeft            int      `json:"hop_left"`
	CommonFieldName    string   `json:"common_field_name"`
	Reverse            bool     `json:"reverse"`
}
