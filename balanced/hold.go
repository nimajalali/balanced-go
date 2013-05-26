package balanced

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

const (
	holdsUri = marketplaceUri + "/%v/holds"
)

type Hold struct {
	Account           Account   `json:"account,omitempty"`
	Amount            int       `json:"amount,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
	Description       string    `json:"description,omitempty"`
	ExpiresAt         time.Time `json:"expires_at,omitempty"`
	Fee               string    `json:"fee,omitempty"`
	Id                string    `json:"id,omitempty"`
	IsVoid            bool      `json:"is_void,omitempty"`
	Meta              MetaType  `json:"meta,omitempty"`
	Source            Card      `json:"source,omitempty"`
	TransactionNumber string    `json:"transaction_number,omitempty"`
	Uri               string    `json:"uri,omitempty"`
}

type ListOfHolds struct {
	FirstUri    string `json:"first_uri,omitempty"`
	Items       []Hold `json:"items,omitempty"`
	LastUri     string `json:"last_uri,omitempty"`
	Limit       int    `json:"limit,omitempty"`
	NextUri     string `json:"next_uri,omitempty"`
	Offset      int    `json:"offset,omitempty"`
	PreviousUri string `json:"previous_uri,omitempty"`
	Total       int    `json:"total,omitempty"`
	Uri         string `json:"uri,omitempty"`
}

// Creates a hold against a card. Returns a uri that can later be used to create
// a debit, up to the full amount of the hold.
func CreateNewHold(uri, accountUri, appearsOnStatementAs, description, sourceUri,
	cardUri string, amount int, meta MetaType) (hold *Hold, err error) {

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

	hold = &Hold{}
	err = post(uri, payload, hold)

	return
}

// Retrieves the details of a hold that you've previously created. Use the uri
// that was previously returned, and the corresponding hold information will be
// returned.
func RetrieveHold(uri string) (hold *Hold, err error) {
	hold = &Hold{}
	err = get(uri, nil, hold)

	return
}

// Returns a list of holds you've previously created. The holds are returned in
// sorted order, with the most recent holds appearing first.
func ListAllHolds(limit, offset int) (listOfHolds *ListOfHolds, err error) {
	payload := defaultPayload(limit, offset)

	uri := fmt.Sprintf(holdsUri, marketplaceId)

	listOfHolds = &ListOfHolds{}
	err = get(uri, payload, listOfHolds)

	return
}

// Returns a list of holds you've previously created. The holds are returned in
// sorted order, with the most recent holds appearing first.
func ListAllHoldsForAccount(uri string, limit, offset int) (listOfHolds *ListOfHolds, err error) {
	payload := defaultPayload(limit, offset)

	listOfHolds = &ListOfHolds{}
	err = get(uri, payload, listOfHolds)

	return
}

// Updates information about a hold
func UpdateHold(uri, description, appearsOnStatementAs string, isVoid bool, meta MetaType) (hold *Hold, err error) {
	payload := url.Values{}

	addToPayload(payload, "description", description)
	addToPayload(payload, "is_void", strconv.FormatBool(isVoid))
	addToPayload(payload, "appears_on_statement_as", appearsOnStatementAs)

	for key, value := range meta {
		addToPayload(payload, "meta["+key+"]", value)
	}

	hold = &Hold{}
	err = put(uri, payload, hold)

	return
}

// Captures a hold. This creates a debit.
func CaptureHold(uri, holdUri, description, appearsOnStatementAs string) (debit *Debit, err error) {
	payload := url.Values{}

	addToPayload(payload, "hold_uri", holdUri)
	addToPayload(payload, "description", description)
	addToPayload(payload, "appears_on_statement_as", appearsOnStatementAs)

	debit = &Debit{}
	err = post(uri, payload, debit)

	return
}

// Voids a hold. This cancels the hold. After voiding, the hold can no longer be
// captured. This operation is irreversible.
func VoidHold(uri, appearsOnStatementAs string, isVoid bool) (hold *Hold, err error) {
	payload := url.Values{}

	addToPayload(payload, "is_void", strconv.FormatBool(isVoid))
	addToPayload(payload, "appears_on_statement_as", appearsOnStatementAs)

	hold = &Hold{}
	err = put(uri, payload, hold)

	return
}
