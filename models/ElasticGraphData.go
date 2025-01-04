package models

type EsBucket struct {
	DocCount string `json:"doc_count"`
	Key      string `json:"key"`
}
