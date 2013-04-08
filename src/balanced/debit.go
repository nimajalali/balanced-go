package balanced

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

const (
	debitsUri = marketplaceUri + "/%v/debits"
)

type Debit struct {
	Account              Account   `json:"account, omitempty"`
	Amount               int       `json:"amount, omitempty"`
	AppearsOnStatementAs string    `json:"appears_on_statement_as, omitempty"`
	AvailableAt          time.Time `json:"available_at"`
	CreatedAt            time.Time `json:"created_at, omitempty"`
	Description          string    `json:"description, omitempty"`
	Fee                  string    `json:"fee, omitempty"`
	Hold                 Hold      `json:"hold, omitempty"`
	Id                   string    `json:"id, omitempty"`
	Meta                 MetaType  `json:"meta, omitempty"`
	OnBehalfOf           string    `json:"on_behalf_of, omitempty"`
	RefundsUri           string    `json:"refunds_uri, omitempty"`
	Source               Card      `json:"source, omitempty"`
	Status               string    `json:"status, omitempty"`
	TransactionNumber    string    `json:"transaction_number, omitempty"`
	Uri                  string    `json:"uri, omitempty"`
}

type ListOfDebits struct {
	FirstUri    string  `json:"first_uri, omitempty"`
	Items       []Debit `json:"items, omitempty"`
	LastUri     string  `json:"last_uri, omitempty"`
	Limit       int     `json:"limit, omitempty"`
	NextUri     string  `json:"next_uri, omitempty"`
	Offset      int     `json:"offset, omitempty"`
	PreviousUri string  `json:"previous_uri, omitempty"`
	Total       int     `json:"total, omitempty"`
	Uri         string  `json:"uri, omitempty"`
}

// Debits an account. Returns a uri that can later be used to reference this
// debit. Successful creation of a debit using a card will return an associated
// hold mapping as part of the response. This hold was created and captured
// behind the scenes automatically. For ACH debits there is no corresponding
// hold.
func CreateNewDebit(uri, description, appearsOnStatementAs, accountUri,
	onBehalfOfUri, holdUri, sourceUri string, amount int,
	meta MetaType) (*Debit, error) {

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

	resp, err := post(uri, payload)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into Debit
	debit := Debit{}
	if err := json.Unmarshal(resp, &debit); err != nil {
		return nil, err
	}

	return &debit, nil
}

// Retrieves the details of a created debit.
func RetrieveDebit(uri string) (*Debit, error) {
	resp, err := get(uri, nil)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into Debit
	debit := Debit{}
	if err := json.Unmarshal(resp, &debit); err != nil {
		return nil, err
	}

	return &debit, nil
}

// Returns a list of debits you've previously created. The debits are returned
// in sorted order, with the most recent debits appearing first.
func ListAllDebits(limit, offset int) (*ListOfDebits, error) {
	payload := url.Values{
		"limit":  {strconv.Itoa(limit)},
		"offset": {strconv.Itoa(offset)},
	}

	resp, err := get(fmt.Sprintf(debitsUri, marketplaceId), payload)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into ListOfDebits
	listOfDebits := ListOfDebits{}
	if err := json.Unmarshal(resp, &listOfDebits); err != nil {
		return nil, err
	}

	return &listOfDebits, nil
}

// Returns a list of debits you've previously created against a specific account
// The debits_uri is a convenient uri provided so that you can simply issue a
// GET to the debits_uri. The debits are returned in sorted order, with the most
// recent debits appearing first.
func ListAllDebitsForAccount(uri string, limit,
	offset int) (*ListOfDebits, error) {

	payload := url.Values{
		"limit":  {strconv.Itoa(limit)},
		"offset": {strconv.Itoa(offset)},
	}

	resp, err := get(uri, payload)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into ListOfDebits
	listOfDebits := ListOfDebits{}
	if err := json.Unmarshal(resp, &listOfDebits); err != nil {
		return nil, err
	}

	return &listOfDebits, nil
}

func UpdateDebit(uri, description string, meta MetaType) (*Debit, error) {
	payload := url.Values{}

	addToPayload(payload, "description", description)

	for key, value := range meta {
		addToPayload(payload, "meta["+key+"]", value)
	}

	resp, err := put(uri, payload)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into Debit
	debit := Debit{}
	if err := json.Unmarshal(resp, &debit); err != nil {
		return nil, err
	}

	return &debit, nil
}

func RefundDebit(uri string) (*Refund, error) {
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
