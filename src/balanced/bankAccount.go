package balanced

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"
)

const (
	bankAccountsUri          = "/v1/bank_accounts"
	BankAccountTypeChecking  = "checking"
	BankAccountTypeSavings   = "savings"
	VerificationTypeVerified = "verified"
	VerificationTypePending  = "pending"
)

type BankAccount struct {
	AccountNumber    string    `json:"account_number, omitempty"`
	AccountUri       string    `json:"account_uri, omitempty"`
	BankCode         string    `json:"bank_account, omitempty"`
	BankName         string    `json:"bank_name, omitempty"`
	CanDebit         bool      `json:"can_debit, omitempty"`
	CreatedAt        time.Time `json:"created_at, omitempty"`
	CreditsUri       string    `json:"credits_uri, omitempty"`
	Fingerprint      string    `json:"fingerprint, omitempty"`
	Id               string    `json:"id, omitempty"`
	IsValid          bool      `json:"is_valid, omitempty"`
	LastFour         string    `json:"last_four, omitempty"`
	Meta             MetaType  `json:"meta, omitempty"`
	Name             string    `json:"name, omitempty"`
	RoutingNumber    string    `json:"routing_number, omitempty"`
	Type             string    `json:"type, omitempty"`
	Uri              string    `json:"uri", omitempty`
	VerificationUri  string    `json:"verification_uri, omitempty"`
	VerificationsUri string    `json:"verifications_uri, omitempty"`
}

type ListOfBankAccounts struct {
	Items  []BankAccount `json:"items, omitempty"`
	Limit  int           `json:"limit, omitempty"`
	Offset int           `json:"offset, omitempty"`
	Total  int           `json:"total, omitempty"`
}

type Verification struct {
	Attempts          int    `json:"attempts, omitempty"`
	Id                string `json:"id, omitempty"`
	RemainingAttempts int    `json:"remaining_attempts, omitempty"`
	State             string `json:"pending, omitempty"`
	Uri               string `json:"uri, omitempty"`
}

type ListOfVerifications struct {
	FirstUri    string         `json:"first_uri, omitempty"`
	Items       []Verification `json:"items, omitempty"`
	LastUri     string         `json:"last_uri, omitempty"`
	Limit       int            `json:"limit, omitempty"`
	NextUri     string         `json:"next_uri, omitempty"`
	Offset      int            `json:"offset, omitempty"`
	PreviousUri string         `json:"previous_uri, omitempty"`
	Total       int            `json:"total, omitempty"`
	Uri         string         `json:"uri, omitempty"`
}

// You'll eventually want to be able to credit bank accounts without having to
// ask your users for their information over and over again. To do this, you'll
// need to create a bank account object.
// NOTE To debit a bank account you must first verify it.
// WARNING PCI Compliance required to use this functionality.
func CreateNewBankAccount(name, accountNumber, routingNumber,
	accountType string) (*BankAccount, error) {

	payload := url.Values{
		"name":           {name},
		"account_number": {accountNumber},
		"routing_number": {routingNumber},
		"type":           {accountType},
	}

	resp, err := post(bankAccountsUri, payload)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into BankAccount
	bankAccount := BankAccount{}
	if err := json.Unmarshal(resp, &bankAccount); err != nil {
		return nil, err
	}

	return &bankAccount, nil
}

// Retrieves the details of a bank account that has previously been created.
// Supply the uri that was returned from your previous request, and the
// corresponding bank account information will be returned. The same information
// is returned when creating the bank account.
// uri: In the form of /v1/bank_accounts/:bank_account_id
func RetrieveBankAccount(uri string) (*BankAccount, error) {
	resp, err := get(uri, nil)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into BankAccount
	bankAccount := BankAccount{}
	if err := json.Unmarshal(resp, &bankAccount); err != nil {
		return nil, err
	}

	return &bankAccount, nil
}

// Returns a list of bank accounts that you've created but haven't deleted.
func ListAllBankAccounts(limit, offset int) (*ListOfBankAccounts, error) {
	payload := url.Values{
		"limit":  {strconv.Itoa(limit)},
		"offset": {strconv.Itoa(offset)},
	}

	resp, err := get(bankAccountsUri, payload)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into ListOfBankAccounts
	listOfBankAccounts := ListOfBankAccounts{}
	if err := json.Unmarshal(resp, &listOfBankAccounts); err != nil {
		return nil, err
	}

	return &listOfBankAccounts, nil
}

// Permanently delete a bank account. It cannot be undone. All associated
// credits with a deleted bank account will not be affected.
// uri: In the form of /v1/bank_accounts/:bank_account_id
func DeleteBankAccount(uri string) error {
	_, err := delete(uri, nil)
	if err != nil {
		return err
	}

	return nil
}

// Creates a new bank account verification.
// uri: In the form of /v1/bank_accounts/:bank_account_id
func VerifyBankAccount(uri string) (*Verification, error) {
	resp, err := post(uri+"/verifications", nil)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into Verification
	verification := Verification{}
	if err := json.Unmarshal(resp, &verification); err != nil {
		return nil, err
	}

	return &verification, nil
}

// Retrieve a Verification for a Bank Account
// uri: /v1/bank_accounts/:bank_account_id/verifications/:verification_id
func RetrieveBankAccountVerification(uri string) (*Verification, error) {
	resp, err := get(uri, nil)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into Verification
	verification := Verification{}
	if err := json.Unmarshal(resp, &verification); err != nil {
		return nil, err
	}

	return &verification, nil
}

// List All Verifications for a Bank Account
// uri: In the form of /v1/bank_accounts/:bank_account_id/verifications
func ListAllBankAccountVerifications(uri string) (*ListOfVerifications, error) {
	resp, err := get(uri, nil)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into List of Verifications
	listOfVerifications := ListOfVerifications{}
	if err := json.Unmarshal(resp, &listOfVerifications); err != nil {
		return nil, err
	}

	return &listOfVerifications, nil
}

// Confirms the trial deposit amounts. For the test environment the trial
// deposit amounts are always 1 and 1.
// uri: /v1/bank_accounts/:bank_account_id/verifications/:verification_id
func ConfirmBankAccountVerification(uri string, amountOne,
	amountTwo int64) (*Verification, error) {

	payload := url.Values{
		"amount_1": {strconv.FormatInt(amountOne, 10)},
		"amount_2": {strconv.FormatInt(amountTwo, 10)},
	}

	resp, err := put(uri, payload)
	if err != nil {
		return nil, err
	}

	// Attempt to parse respo`nmnse into Verification
	verification := Verification{}
	if err := json.Unmarshal(resp, &verification); err != nil {
		return nil, err
	}

	return &verification, nil
}
