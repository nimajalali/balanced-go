package balanced

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

const (
	holdsUri = marketplaceUri + "/%v/holds"
)

type Hold struct {
	Account           Account   `json:"account, omitempty"`
	Amount            int       `json:"amount, omitempty"`
	CreatedAt         time.Time `json:"created_at, omitempty"`
	Description       string    `json:"description, omitempty"`
	ExpiresAt         time.Time `json:"expires_at, omitempty"`
	Fee               string    `json:"fee, omitempty"`
	Id                string    `json:"id, omitempty"`
	IsVoid            bool      `json:"is_void, omitempty"`
	Meta              MetaType  `json:"meta, omitempty"`
	Source            Card      `json:"source, omitempty"`
	TransactionNumber string    `json:"transaction_number, omitempty"`
	Uri               string    `json:"uri, omitempty"`
}

type ListOfHolds struct {
	FirstUri    string `json:"first_uri, omitempty"`
	Items       []Hold `json:"items, omitempty"`
	LastUri     string `json:"last_uri, omitempty"`
	Limit       int    `json:"limit, omitempty"`
	NextUri     string `json:"next_uri, omitempty"`
	Offset      int    `json:"offset, omitempty"`
	PreviousUri string `json:"previous_uri, omitempty"`
	Total       int    `json:"total, omitempty"`
	Uri         string `json:"uri, omitempty"`
}

// Creates a hold against a card. Returns a uri that can later be used to create
// a debit, up to the full amount of the hold.
func CreateNewHold(uri, accountUri, appearsOnStatementAs, description,
	sourceUri, cardUri string, amount int, meta MetaType) (*Hold, error) {

	payload := url.Values{
		"amount": {strconv.Itoa(amount)},
	}

	addToPayload(payload, "account_uri", accountUri)
	addToPayload(payload, "appears_on_statement_as", appearsOnStatementAs)
	addToPayload(payload, "description", description)
	addToPayload(payload, "source_uri", sourceUri)
	addToPayload(payload, "card_uri", cardUri)

	for key, value := range meta {
		addToPayload(payload, "meta["+key+"]", value)
	}

	resp, err := post(uri, payload)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into Hold
	hold := Hold{}
	if err := json.Unmarshal(resp, &hold); err != nil {
		return nil, err
	}

	return &hold, nil
}

// Retrieves the details of a hold that you've previously created. Use the uri
// that was previously returned, and the corresponding hold information will be
// returned.
func RetrieveHold(uri string) (*Hold, error) {
	resp, err := get(uri, nil)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into Hold
	hold := Hold{}
	if err := json.Unmarshal(resp, &hold); err != nil {
		return nil, err
	}

	return &hold, nil
}

// Returns a list of holds you've previously created. The holds are returned in
// sorted order, with the most recent holds appearing first.
func ListAllHolds(limit, offset int) (*ListOfHolds, error) {
	payload := url.Values{
		"limit":  {strconv.Itoa(limit)},
		"offset": {strconv.Itoa(offset)},
	}

	resp, err := get(fmt.Sprintf(holdsUri, marketplaceId), payload)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into ListOfHolds
	listOfHolds := ListOfHolds{}
	if err := json.Unmarshal(resp, &listOfHolds); err != nil {
		return nil, err
	}

	return &listOfHolds, nil
}

// Returns a list of holds you've previously created. The holds are returned in
// sorted order, with the most recent holds appearing first.
func ListAllHoldsForAccount(uri string, limit,
	offset int) (*ListOfHolds, error) {

	payload := url.Values{
		"limit":  {strconv.Itoa(limit)},
		"offset": {strconv.Itoa(offset)},
	}

	resp, err := get(uri, payload)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into ListOfHolds
	listOfHolds := ListOfHolds{}
	if err := json.Unmarshal(resp, &listOfHolds); err != nil {
		return nil, err
	}

	return &listOfHolds, nil
}

// Updates information about a hold
func UpdateHold(uri, description, appearsOnStatementAs string, isVoid bool,
	meta MetaType) (*Hold, error) {

	payload := url.Values{}

	addToPayload(payload, "description", description)
	addToPayload(payload, "is_void", strconv.FormatBool(isVoid))
	addToPayload(payload, "appears_on_statement_as", appearsOnStatementAs)

	for key, value := range meta {
		addToPayload(payload, "meta["+key+"]", value)
	}

	resp, err := put(uri, payload)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into ListOfHolds
	hold := Hold{}
	if err := json.Unmarshal(resp, &hold); err != nil {
		return nil, err
	}

	return &hold, nil
}

// Captures a hold. This creates a debit.
func CaptureHold(uri, holdUri, description,
	appearsOnStatementAs string) (*Debit, error) {

	payload := url.Values{}

	addToPayload(payload, "hold_uri", holdUri)
	addToPayload(payload, "description", description)
	addToPayload(payload, "appears_on_statement_as", appearsOnStatementAs)

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

// Voids a hold. This cancels the hold. After voiding, the hold can no longer be
// captured. This operation is irreversible.
func VoidHold(uri, appearsOnStatementAs string, isVoid bool) (*Hold, error) {
	payload := url.Values{}

	addToPayload(payload, "is_void", strconv.FormatBool(isVoid))
	addToPayload(payload, "appears_on_statement_as", appearsOnStatementAs)

	resp, err := put(uri, payload)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into ListOfHolds
	hold := Hold{}
	if err := json.Unmarshal(resp, &hold); err != nil {
		return nil, err
	}

	return &hold, nil
}
