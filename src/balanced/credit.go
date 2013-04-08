package balanced

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"
)

const (
	creditsUri = "/v1/credits"
)

type Credit struct {
	Account           Account     `json:"account, omitempty"`
	Amount            int         `json:"amount, omitempty"`
	AvailableAt       time.Time   `json:"available_at"`
	BankAccount       BankAccount `json:"bank_account, omitempty"`
	CreatedAt         time.Time   `json:"created_at, omitempty"`
	Description       string      `json:"description, omitempty"`
	Destination       BankAccount `json:"destination, omitempty"`
	Fee               string      `json:"fee, omitempty"`
	Id                string      `json:"id, omitempty"`
	IsVoid            bool        `json:"is_void, omitempty"`
	Meta              MetaType    `json:"meta, omitempty"`
	Status            string      `json:"status, omitempty"`
	Source            Card        `json:"source, omitempty"`
	TransactionNumber string      `json:"transaction_number, omitempty"`
	Uri               string      `json:"uri, omitempty"`
}

type ListOfCredits struct {
	FirstUri    string   `json:"first_uri, omitempty"`
	Items       []Credit `json:"items, omitempty"`
	LastUri     string   `json:"last_uri, omitempty"`
	Limit       int      `json:"limit, omitempty"`
	NextUri     string   `json:"next_uri, omitempty"`
	Offset      int      `json:"offset, omitempty"`
	PreviousUri string   `json:"previous_uri, omitempty"`
	Total       int      `json:"total, omitempty"`
	Uri         string   `json:"uri, omitempty"`
}

// To credit a new bank account, you simply pass the amount along with the bank
// account details. We do not store this bank account when you create a credit
// this way, so you can safely assume that the information has been deleted.
// WARNING PCI Compliance required to use this functionality.
func CreditNewBankAccount(amount int, description string,
	bankAccount *BankAccount) (*Credit, error) {

	// Required values
	payload := url.Values{
		"amount":                       {strconv.Itoa(amount)},
		"bank_account[name]":           {bankAccount.Name},
		"bank_account[account_number]": {bankAccount.AccountNumber},
		"bank_account[routing_number]": {bankAccount.RoutingNumber},
		"bank_account[type]":           {bankAccount.Type},
	}

	addToPayload(payload, "description", description)

	for key, value := range bankAccount.Meta {
		addToPayload(payload, "bank_account[meta["+key+"]]", value)
	}

	resp, err := post(creditsUri, payload)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into Credit
	credit := Credit{}
	if err := json.Unmarshal(resp, &credit); err != nil {
		return nil, err
	}

	return &credit, nil
}

// To credit an existing bank account, you simply pass the amount to the nested
// credit endpoint of a bank account. The credits_uri is a convenient uri
// provided so that you can simply issue a POST with the amount and a credit
// shall be created.
func CreditExistingBankAccount(uri, description string,
	amount int) (*Credit, error) {

	// Required values
	payload := url.Values{
		"amount": {strconv.Itoa(amount)},
	}

	addToPayload(payload, "description", description)

	resp, err := post(uri, payload)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into Credit
	credit := Credit{}
	if err := json.Unmarshal(resp, &credit); err != nil {
		return nil, err
	}

	return &credit, nil
}

// Retrieves the details of a credit that you've previously created. Use the uri
// that was previously returned, and the corresponding credit information will
// be returned.
func RetrieveCredit(uri string) (*Credit, error) {
	resp, err := get(uri, nil)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into Credit
	credit := Credit{}
	if err := json.Unmarshal(resp, &credit); err != nil {
		return nil, err
	}

	return &credit, nil
}

// Returns a list of credits you've previously created. The credits are returned
// in sorted order, with the most recent credits appearing first.
func ListAllCredits(limit, offset int) (*ListOfCredits, error) {
	payload := url.Values{
		"limit":  {strconv.Itoa(limit)},
		"offset": {strconv.Itoa(offset)},
	}

	resp, err := get(creditsUri, payload)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into Credit
	listOfCredits := ListOfCredits{}
	if err := json.Unmarshal(resp, &listOfCredits); err != nil {
		return nil, err
	}

	return &listOfCredits, nil
}

// Returns a list of credits you've previously created to a specific bank
// account. The credits_uri is a convenient uri provided so that you can simply
// issue a GET to the credits_uri. The credits are returned in sorted order,
// with the most recent credits appearing first.
func ListAllCreditsForBankAccount(uri string, limit,
	offset int) (*ListOfCredits, error) {

	payload := url.Values{
		"limit":  {strconv.Itoa(limit)},
		"offset": {strconv.Itoa(offset)},
	}

	resp, err := get(uri, payload)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into Credit
	listOfCredits := ListOfCredits{}
	if err := json.Unmarshal(resp, &listOfCredits); err != nil {
		return nil, err
	}

	return &listOfCredits, nil
}

func CreateNewCreditForAccount(uri, description, appearsOnStatementAs,
	destinationUri, bankAccountUri string, amount int,
	meta MetaType) (*Credit, error) {

	// Required values
	payload := url.Values{
		"amount": {strconv.Itoa(amount)},
	}

	addToPayload(payload, "description", description)
	addToPayload(payload, "appears_on_statement_as", appearsOnStatementAs)
	addToPayload(payload, "destination_uri", destinationUri)
	addToPayload(payload, "bank_account_uri", bankAccountsUri)

	resp, err := post(uri, payload)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into Credit
	credit := Credit{}
	if err := json.Unmarshal(resp, &credit); err != nil {
		return nil, err
	}

	return &credit, nil
}

func ListAllCreditsForAccount(uri string, limit,
	offset int) (*ListOfCredits, error) {

	payload := url.Values{
		"limit":  {strconv.Itoa(limit)},
		"offset": {strconv.Itoa(offset)},
	}

	resp, err := get(uri, payload)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into list of credits
	listOfCredits := ListOfCredits{}
	if err := json.Unmarshal(resp, &listOfCredits); err != nil {
		return nil, err
	}

	return &listOfCredits, nil
}
