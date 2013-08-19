package balanced

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

const (
	debitsUri = marketplaceUri + "/%v/debits"
)

type Debit struct {
	Account              Account   `json:"account,omitempty"`
	Amount               int       `json:"amount,omitempty"`
	AppearsOnStatementAs string    `json:"appears_on_statement_as,omitempty"`
	AvailableAt          time.Time `json:"available_at,omitempty"`
	CreatedAt            time.Time `json:"created_at,omitempty"`
	Description          string    `json:"description,omitempty"`
	Fee                  string    `json:"fee,omitempty"`
	Hold                 Hold      `json:"hold,omitempty"`
	Id                   string    `json:"id,omitempty"`
	Meta                 MetaType  `json:"meta,omitempty"`
	OnBehalfOf           string    `json:"on_behalf_of,omitempty"`
	RefundsUri           string    `json:"refunds_uri,omitempty"`
	Source               Card      `json:"source,omitempty"`
	Status               string    `json:"status,omitempty"`
	TransactionNumber    string    `json:"transaction_number,omitempty"`
	Uri                  string    `json:"uri,omitempty"`
}

type ListOfDebits struct {
	FirstUri    string  `json:"first_uri,omitempty"`
	Items       []Debit `json:"items,omitempty"`
	LastUri     string  `json:"last_uri,omitempty"`
	Limit       int     `json:"limit,omitempty"`
	NextUri     string  `json:"next_uri,omitempty"`
	Offset      int     `json:"offset,omitempty"`
	PreviousUri string  `json:"previous_uri,omitempty"`
	Total       int     `json:"total,omitempty"`
	Uri         string  `json:"uri,omitempty"`
}

// Debits an account. Returns a uri that can later be used to reference this
// debit. Successful creation of a debit using a card will return an associated
// hold mapping as part of the response. This hold was created and captured
// behind the scenes automatically. For ACH debits there is no corresponding
// hold.
func CreateNewDebit(uri, description, appearsOnStatementAs, accountUri,
	onBehalfOfUri, holdUri, sourceUri string, amount int,
	meta MetaType) (debit *Debit, err error) {

	payload := url.Values{}

	addToPayload(payload, "amount", strconv.Itoa(amount))
	addToPayload(payload, "description", description)
	addToPayload(payload, "appears_on_statement_as", appearsOnStatementAs)
	addToPayload(payload, "account_uri", accountUri)
	addToPayload(payload, "on_behalf_of_uri", onBehalfOfUri)
	addToPayload(payload, "hold_uri", holdUri)
	addToPayload(payload, "source_uri", sourceUri)

	for key, value := range meta {
		addToPayload(payload, "meta["+key+"]", value)
	}

	debit = &Debit{}
	err = post(uri, payload, debit)

	return
}

// Retrieves the details of a created debit.
func RetrieveDebit(uri string) (debit *Debit, err error) {
	debit = &Debit{}
	err = get(uri, nil, debit)

	return
}

// Returns a list of debits you've previously created. The debits are returned
// in sorted order, with the most recent debits appearing first.
func ListAllDebits(limit, offset int) (listOfDebits *ListOfDebits, err error) {
	payload := defaultPayload(limit, offset)

	uri := fmt.Sprintf(debitsUri, marketplaceId)

	listOfDebits = &ListOfDebits{}
	err = get(uri, payload, listOfDebits)

	return
}

// Returns a list of debits you've previously created against a specific account
// The debits_uri is a convenient uri provided so that you can simply issue a
// GET to the debits_uri. The debits are returned in sorted order, with the most
// recent debits appearing first.
func ListAllDebitsForAccount(uri string, limit, offset int) (listOfDebits *ListOfDebits, err error) {
	payload := defaultPayload(limit, offset)

	listOfDebits = &ListOfDebits{}
	err = get(uri, payload, listOfDebits)

	return
}

func UpdateDebit(uri, description string, meta MetaType) (debit *Debit, err error) {
	payload := url.Values{}

	addToPayload(payload, "description", description)

	for key, value := range meta {
		addToPayload(payload, "meta["+key+"]", value)
	}

	debit = &Debit{}
	err = put(uri, payload, debit)

	return
}

func RefundDebit(uri string) (refund *Refund, err error) {
	refund = &Refund{}
	err = post(uri, nil, refund)

	return
}
