package balanced

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

const (
	refundsUri = marketplaceUri + "/%v/refunds"
)

type Refund struct {
	Account              Account   `json:"account,omitempty"`
	Amount               int       `json:"amount,omitempty"`
	AppearsOnStatementAs string    `json:"appears_on_statement_as,omitempty"`
	CreatedAt            time.Time `json:"created_at,omitempty"`
	Debit                Debit     `json:"debit,omitempty"`
	Description          string    `json:"description,omitempty"`
	Fee                  string    `json:"fee,omitempty"`
	Id                   string    `json:"id,omitempty"`
	Meta                 MetaType  `json:"meta,omitempty"`
	TransactionNumber    string    `json:"transaction_number,omitempty"`
	Uri                  string    `json:"uri,omitempty"`
}

type ListOfRefunds struct {
	FirstUri    string   `json:"first_uri,omitempty"`
	Items       []Refund `json:"items,omitempty"`
	LastUri     string   `json:"last_uri,omitempty"`
	Limit       int      `json:"limit,omitempty"`
	NextUri     string   `json:"next_uri,omitempty"`
	Offset      int      `json:"offset,omitempty"`
	PreviousUri string   `json:"previous_uri,omitempty"`
	Total       int      `json:"total,omitempty"`
	Uri         string   `json:"uri,omitempty"`
}

// Issues a refund from a debit. You can either refund the full amount of the
// debit or you can issue a partial refund, where the amount is less than the
// charged amount.
func IssueRefund(description, debitUri string, amount int, meta MetaType) (refund *Refund, err error) {
	payload := url.Values{}

	addToPayload(payload, "amount", strconv.Itoa(amount))
	addToPayload(payload, "description", description)
	addToPayload(payload, "debit_uri", debitsUri)

	for key, value := range meta {
		addToPayload(payload, "meta["+key+"]", value)
	}

	uri := fmt.Sprintf(refundsUri, marketplaceId)

	refund = &Refund{}
	err = post(uri, payload, refund)

	return
}

// Retrieves the details of a refund that you've previously created. Use the uri
// that was previously returned, and the corresponding refund information will
// be returned.
func RetrieveRefund(uri string) (refund *Refund, err error) {
	refund = &Refund{}
	err = post(uri, nil, refund)

	return
}

// Returns a list of refunds you've previously created. The refunds are returned
// in sorted order, with the most recent refunds appearing first.
func ListAllRefunds(limit, offset int) (listOfRefunds *ListOfRefunds, err error) {
	payload := defaultPayload(limit, offset)

	uri := fmt.Sprintf(refundsUri, marketplaceId)

	listOfRefunds = &ListOfRefunds{}
	err = get(uri, payload, listOfRefunds)

	return
}

// Returns a list of refunds you've previously created against a specific
// account. The refunds are returned in sorted order, with the most recent
// refunds appearing first.
func ListAllRefundsForAccount(uri string, limit, offset int) (listOfRefunds *ListOfRefunds, err error) {
	payload := defaultPayload(limit, offset)

	listOfRefunds = &ListOfRefunds{}
	err = get(uri, payload, listOfRefunds)

	return
}

// Updates information about a refund
func UpdateRefund(uri, description string, meta MetaType) (refund *Refund, err error) {
	payload := url.Values{}

	addToPayload(payload, "description", description)

	for key, value := range meta {
		addToPayload(payload, "meta["+key+"]", value)
	}

	refund = &Refund{}
	err = put(uri, payload, refund)

	return
}
