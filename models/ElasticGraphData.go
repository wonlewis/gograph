package models

type EsBucket struct {
	DocCount string `json:"doc_count"`
	Key      string `json:"key"`
}

type ValidationStatus string

const (
	OK   ValidationStatus = "no error"
	ERR1                  = "query is invalid"
	ERR2                  = "constraint is invalid"
	ERR3                  = "field does not exist"
)
