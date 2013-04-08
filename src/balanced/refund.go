package balanced

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

const (
	refundsUri = marketplaceUri + "/%v/refunds"
)

type Refund struct {
	Account              Account   `json:"account, omitempty"`
	Amount               int       `json:"amount, omitempty"`
	AppearsOnStatementAs string    `json:"appears_on_statement_as, omitempty"`
	CreatedAt            time.Time `json:"created_at, omitempty"`
	Debit                Debit     `json:"debit, omitempty"`
	Description          string    `json:"description, omitempty"`
	Fee                  string    `json:"fee, omitempty"`
	Id                   string    `json:"id, omitempty"`
	Meta                 MetaType  `json:"meta, omitempty"`
	TransactionNumber    string    `json:"transaction_number, omitempty"`
	Uri                  string    `json:"uri, omitempty"`
}

type ListOfRefunds struct {
	FirstUri    string   `json:"first_uri, omitempty"`
	Items       []Refund `json:"items, omitempty"`
	LastUri     string   `json:"last_uri, omitempty"`
	Limit       int      `json:"limit, omitempty"`
	NextUri     string   `json:"next_uri, omitempty"`
	Offset      int      `json:"offset, omitempty"`
	PreviousUri string   `json:"previous_uri, omitempty"`
	Total       int      `json:"total, omitempty"`
	Uri         string   `json:"uri, omitempty"`
}

// Issues a refund from a debit. You can either refund the full amount of the
// debit or you can issue a partial refund, where the amount is less than the
// charged amount.
func IssueRefund(description, debitUri string, amount int,
	meta MetaType) (*Refund, error) {

	payload := url.Values{}

	addToPayload(payload, "amount", strconv.Itoa(amount))
	addToPayload(payload, "description", description)
	addToPayload(payload, "debit_uri", debitsUri)

	for key, value := range meta {
		addToPayload(payload, "meta["+key+"]", value)
	}

	resp, err := post(fmt.Sprintf(refundsUri, marketplaceId), payload)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into Refund
	refund := Refund{}
	if err := json.Unmarshal(resp, &refund); err != nil {
		return nil, err
	}

	return &refund, nil
}

// Retrieves the details of a refund that you've previously created. Use the uri
// that was previously returned, and the corresponding refund information will
// be returned.
func RetrieveRefund(uri string) (*Refund, error) {
	resp, err := post(uri, nil)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into Refund
	refund := Refund{}
	if err := json.Unmarshal(resp, &refund); err != nil {
		return nil, err
	}

	return &refund, nil
}

// Returns a list of refunds you've previously created. The refunds are returned
// in sorted order, with the most recent refunds appearing first.
func ListAllRefunds(limit, offset int) (*ListOfRefunds, error) {
	payload := url.Values{
		"limit":  {strconv.Itoa(limit)},
		"offset": {strconv.Itoa(offset)},
	}

	resp, err := get(fmt.Sprintf(refundsUri, marketplaceId), payload)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into ListOfRefunds
	listOfRefunds := ListOfRefunds{}
	if err := json.Unmarshal(resp, &listOfRefunds); err != nil {
		return nil, err
	}

	return &listOfRefunds, nil
}

// Returns a list of refunds you've previously created against a specific
// account. The refunds are returned in sorted order, with the most recent
// refunds appearing first.
func ListAllRefundsForAccount(uri string, limit,
	offset int) (*ListOfRefunds, error) {

	payload := url.Values{
		"limit":  {strconv.Itoa(limit)},
		"offset": {strconv.Itoa(offset)},
	}

	resp, err := get(uri, payload)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into ListOfRefunds
	listOfRefunds := ListOfRefunds{}
	if err := json.Unmarshal(resp, &listOfRefunds); err != nil {
		return nil, err
	}

	return &listOfRefunds, nil
}

// Updates information about a refund
func UpdateRefund(uri, description string, meta MetaType) (*Refund, error) {
	payload := url.Values{}

	addToPayload(payload, "description", description)

	for key, value := range meta {
		addToPayload(payload, "meta["+key+"]", value)
	}

	resp, err := put(uri, payload)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into Refund
	refund := Refund{}
	if err := json.Unmarshal(resp, &refund); err != nil {
		return nil, err
	}

	return &refund, nil
}
