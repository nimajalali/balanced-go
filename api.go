package balanced

import (
	"fmt"
)

// Custom Error to handle balanced api responses. Implements
// error interface
type ApiError struct {
	Additional   string `json:"additional,omitempty"`
	CategoryType string `json:"category_type,omitempty"`
	CategoryCode string `json:"category_code,omitempty"`
	Description  string `json:"description,omitempty"`
	Extras       Extras `json:"extras,omitempty"`
	RequestId    string `json:"request_id,omitempty"`
	StatusCode   int    `json:"status_code,omitempty"`
	Status       string `json:"status,omitempty"`
}

type Extras map[string]string

func (e ApiError) Error() string {
	return fmt.Sprintf("StatusCode=\"%v\" Status=\"%v\" "+
		"RequestId=\"%v\" Description=\"%v\" Additional=\"%v\" "+
		"CategoryType=\"%v\" CategoryCode=\"%v\" Extras=\"%v\"",
		e.StatusCode, e.Status, e.RequestId, e.Description, e.Additional,
		e.CategoryType, e.CategoryCode, e.Extras)
}

type ApiDefaultResponse struct {
	_type string                       `json:"_type,omitempty"`
	_uris map[string]map[string]string `json:"_uris,omitempty"`
}
