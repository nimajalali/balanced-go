package balanced

import (
	"time"
)

const (
	apiKeyUri = "/v1/api_keys"
)

type ApiKey struct {
	Merchant  Merchant  `json:"merchant,omitempty"`
	Secret    string    `json:"secret,omitempty"`
	Meta      MetaType  `json:"meta,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Uri       string    `json:"uri,omitempty"`
	Id        string    `json:"id,omitempty"`
}
