package balanced

import (
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
	ApiDefaultResponse
	AccountNumber    string    `json:"account_number,omitempty"`
	AccountType      string    `json:"account_type,omitempty"`
	BankCode         string    `json:"bank_account,omitempty"`
	BankName         string    `json:"bank_name,omitempty"`
	CanDebit         bool      `json:"can_debit,omitempty"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
	CreditsUri       string    `json:"credits_uri,omitempty"`
	DebitsUri        string    `json:"debits_uri,omitempty"`
	Fingerprint      string    `json:"fingerprint,omitempty"`
	Id               string    `json:"id,omitempty"`
	Meta             MetaType  `json:"meta,omitempty"`
	Name             string    `json:"name,omitempty"`
	RoutingNumber    string    `json:"routing_number,omitempty"`
	Type             string    `json:"type,omitempty"`
	Uri              string    `json:"uri,omitempty"`
	VerificationUri  string    `json:"verification_uri,omitempty"`
	VerificationsUri string    `json:"verifications_uri,omitempty"`
}

type ListOfBankAccounts struct {
	ApiDefaultResponse
	Items  []BankAccount `json:"items,omitempty"`
	Limit  int           `json:"limit,omitempty"`
	Offset int           `json:"offset,omitempty"`
	Total  int           `json:"total,omitempty"`
	Uri    string        `json:"uri,omitempty"`
}

type Verification struct {
	ApiDefaultResponse
	Attempts          int    `json:"attempts,omitempty"`
	Id                string `json:"id,omitempty"`
	RemainingAttempts int    `json:"remaining_attempts,omitempty"`
	State             string `json:"pending,omitempty"`
	Uri               string `json:"uri,omitempty"`
}

type ListOfVerifications struct {
	ApiDefaultResponse
	FirstUri    string         `json:"first_uri,omitempty"`
	Items       []Verification `json:"items,omitempty"`
	LastUri     string         `json:"last_uri,omitempty"`
	Limit       int            `json:"limit,omitempty"`
	NextUri     string         `json:"next_uri,omitempty"`
	Offset      int            `json:"offset,omitempty"`
	PreviousUri string         `json:"previous_uri,omitempty"`
	Total       int            `json:"total,omitempty"`
	Uri         string         `json:"uri,omitempty"`
}

// You'll eventually want to be able to credit bank accounts without having to
// ask your users for their information over and over again. To do this, you'll
// need to create a bank account object.
// NOTE To debit a bank account you must first verify it.
// WARNING PCI Compliance required to use this functionality.
func CreateNewBankAccount(name, accountNumber, routingNumber, accountType string) (bankAccount *BankAccount, err error) {
	payload := url.Values{
		"name":           {name},
		"account_number": {accountNumber},
		"routing_number": {routingNumber},
		"type":           {accountType},
	}

	bankAccount = &BankAccount{}
	err = post(bankAccountsUri, payload, bankAccount)

	return
}

// Retrieves the details of a bank account that has previously been created.
// Supply the uri that was returned from your previous request, and the
// corresponding bank account information will be returned. The same information
// is returned when creating the bank account.
// uri: In the form of /v1/bank_accounts/:bank_account_id
func RetrieveBankAccount(uri string) (bankAccount *BankAccount, err error) {
	bankAccount = &BankAccount{}
	err = get(uri, nil, bankAccount)

	return
}

// Returns a list of bank accounts that you've created but haven't deleted.
func ListAllBankAccounts(limit, offset int) (listOfBankAccounts *ListOfBankAccounts, err error) {
	payload := defaultPayload(limit, offset)

	listOfBankAccounts = &ListOfBankAccounts{}
	err = get(bankAccountsUri, payload, listOfBankAccounts)

	return
}

// Permanently delete a bank account. It cannot be undone. All associated
// credits with a deleted bank account will not be affected.
// uri: In the form of /v1/bank_accounts/:bank_account_id
func DeleteBankAccount(uri string) (err error) {
	err = delete(uri, nil, nil)

	return
}

// Creates a new bank account verification.
// uri: In the form of /v1/bank_accounts/:bank_account_id
func VerifyBankAccount(uri string) (verification *Verification, err error) {
	uri += "/verifications"

	verification = &Verification{}
	err = post(uri, nil, verification)

	return
}

// Retrieve a Verification for a Bank Account
// uri: /v1/bank_accounts/:bank_account_id/verifications/:verification_id
func RetrieveBankAccountVerification(uri string) (verification *Verification, err error) {
	verification = &Verification{}
	err = get(uri, nil, verification)

	return
}

// List All Verifications for a Bank Account
// uri: In the form of /v1/bank_accounts/:bank_account_id/verifications
func ListAllBankAccountVerifications(uri string) (listOfVerifications *ListOfVerifications, err error) {
	listOfVerifications = &ListOfVerifications{}
	err = get(uri, nil, listOfVerifications)

	return
}

// Confirms the trial deposit amounts. For the test environment the trial
// deposit amounts are always 1 and 1.
// uri: /v1/bank_accounts/:bank_account_id/verifications/:verification_id
func ConfirmBankAccountVerification(uri string, amountOne, amountTwo int64) (verification *Verification, err error) {
	payload := url.Values{
		"amount_1": {strconv.FormatInt(amountOne, 10)},
		"amount_2": {strconv.FormatInt(amountTwo, 10)},
	}

	verification = &Verification{}
	err = put(uri, payload, verification)

	return
}
