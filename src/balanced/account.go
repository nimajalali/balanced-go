package balanced

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"
)

const (
	accountsUri         = marketplaceUri + "/%v/accounts"
	accountTypePerson   = "person"
	accountTypeBusiness = "business"
	merchantRole        = "merchant"
)

type Account struct {
	BankAccountsUri string    `json:"bank_accounts_uri, omitempty"`
	CardsUri        string    `json:"cards_uri, omitempty"`
	CreatedAt       time.Time `json:"created_at, omitempty"`
	CreditsUri      string    `json:"credits_uri, omitempty"`
	DebitsUri       string    `json:"debits_uri, omitempty"`
	EmailAddress    string    `json:"email_address, omitempty"`
	HoldsUri        string    `json:"holds_uri, omitempty"`
	Id              string    `json:"id, omitempty"`
	Meta            MetaType  `json:"meta, omitempty"`
	Name            string    `json:"name, omitempty"`
	RefundsUri      string    `json:"refunds_uri, omitempty"`
	Roles           []string  `json:"roles, omitempty"`
	TransactionsUri string    `json:"transactions_uri, omitempty"`
	Uri             string    `json:"uri, omitempty"`
}

type Merchant struct {
	PhoneNumber   string    `json:"phone_number, omitempty"`
	Type          string    `json:"type, omitempty"`
	EmailAddress  string    `json:"email_address, omitempty"`
	Meta          MetaType  `json:"meta, omitempty"`
	TaxId         string    `json:"tax_id, omitempty"`
	Dob           string    `json:"dob, omitempty"`
	Name          string    `json:"name, omitempty"`
	City          string    `json:"city, omitempty"`
	PostalCode    string    `json:"postal_code, omitempty"`
	StreetAddress string    `json:"street_address, omitempty"`
	CountryCode   string    `json:"country_code, omitempty"`
	CreatedAt     time.Time `json:"created_at, omitempty"`
	Uri           string    `json:"uri, omitempty"`
	AccountsUri   string    `json:"accounts_uri, omitempty"`
	Balance       int       `json:"balance, omitempty"`
	ApiKeysUri    string    `json:"api_keys_uri, omitempty"`
	Id            string    `json:"id, omitempty"`
}

type Person struct {
	Name          string `json:"name, omitempty"`
	Dob           string `json:"dob, omitempty"`
	City          string `json:"city, omitempty"`
	PostalCode    string `json:"postal_code, omitempty"`
	StreetAddress string `json:"street_address, omitempty"`
	CountryCode   string `json:"country_code, omitempty"`
	TaxId         string `json:"tax_id, omitempty"`
}

// Accounts help facilitate managing multiple credit cards, debit cards, and
// bank accounts along with different financial transaction operations, i.e.
// refunds, debits, credits.
func CreateAccount() (*Account, error) {
	resp, err := post(fmt.Sprintf(accountsUri, marketplaceId), nil)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into Account
	account := Account{}
	if err := json.Unmarshal(resp, &account); err != nil {
		return nil, err
	}

	return &account, nil
}

// Adding a card to an account activates the ability to debit an account,
// more specifically, charging a card.You can add multiple cards to an account.
// Balanced associates a buyer role to signify whether or not an account has a
// valid credit card, to acquire funds from.
func AddCardToAccount(uri, cardUri string) (*Account, error) {
	payload := url.Values{
		"card_uri": {cardUri},
	}

	resp, err := put(uri, payload)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into Account
	account := Account{}
	if err := json.Unmarshal(resp, &account); err != nil {
		return nil, err
	}

	return &account, nil
}

// Adding a bank account to an account activates the ability to credit an
// account, or in this case, initiate a next-day ACH payment.Balanced does not
// associate a role to signify whether or not an account has a valid bank
// account to send money to.
func AddBankAccountToAccount(uri, bankAccountUri string) (*Account, error) {
	payload := url.Values{
		"bank_account_uri": {bankAccountUri},
	}

	resp, err := put(uri, payload)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into Account
	account := Account{}
	if err := json.Unmarshal(resp, &account); err != nil {
		return nil, err
	}

	return &account, nil
}

// A person, or an individual, is a US based individual or a sole proprietor.
// Balanced associates a merchant role to signify whether or not an account has
// been underwritten.
// WARNING PCI Compliance required to use this functionality.
func UnderwriteIndividual(merchant *Merchant) (*Account, error) {
	// Required Parameters
	payload := url.Values{
		"merchant[phone_number]":   {merchant.PhoneNumber},
		"merchant[name]":           {merchant.Name},
		"merchant[dob]":            {merchant.Dob},
		"merchant[postal_code]":    {merchant.PostalCode},
		"merchant[type]":           {accountTypePerson},
		"merchant[street_address]": {merchant.StreetAddress},
	}

	// Optional Parameters
	addToPayload(payload, "merchant[email_address]", merchant.EmailAddress)
	addToPayload(payload, "merchant[tax_id]", merchant.TaxId)
	addToPayload(payload, "merchant[name]", merchant.Name)
	addToPayload(payload, "merchant[city]", merchant.City)
	addToPayload(payload, "merchant[country_code]", merchant.CountryCode)
	for key, value := range merchant.Meta {
		addToPayload(payload, "merchant[meta["+key+"]]", value)
	}

	resp, err := post(fmt.Sprintf(accountsUri, marketplaceId), payload)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Attempt to parse response into Account
	account := Account{}
	if err := json.Unmarshal(resp, &account); err != nil {
		return nil, err
	}

	return &account, nil
}

// Balanced associates a merchant role to signify whether or not an account has
// been underwritten.
// WARNING PCI Compliance required to use this functionality.
func UnderwriteBusiness(merchant *Merchant, person *Person) (*Account, error) {
	// Required Parameters
	payload := url.Values{
		"merchant[phone_number]":           {merchant.PhoneNumber},
		"merchant[name]":                   {merchant.Name},
		"merchant[tax_id]":                 {merchant.TaxId},
		"merchant[postal_code]":            {merchant.PostalCode},
		"merchant[type]":                   {accountTypeBusiness},
		"merchant[street_address]":         {merchant.StreetAddress},
		"merchant[person[name]]":           {person.Name},
		"merchant[person[dob]]":            {person.Dob},
		"merchant[person[postal_code]]":    {person.PostalCode},
		"merchant[person[street_address]]": {person.StreetAddress},
	}

	// Add optional merchant parameters to payload
	addToPayload(payload, "merchant[email_address]", merchant.EmailAddress)
	addToPayload(payload, "merchant[city]", merchant.City)
	addToPayload(payload, "merchant[country_code]", merchant.CountryCode)
	for key, value := range merchant.Meta {
		addToPayload(payload, "merchant[meta["+key+"]]", value)
	}

	// Add optional person parameters to payload
	addToPayload(payload, "merchant[person[city]]", person.City)
	addToPayload(payload, "merchant[person[country_code]]", person.CountryCode)
	addToPayload(payload, "merchant[person[tax_id]]", person.TaxId)

	resp, err := post(fmt.Sprintf(accountsUri, marketplaceId), payload)
	if err != nil {
		return nil, err
	}

	// Attempt to parse response into Account
	account := Account{}
	if err := json.Unmarshal(resp, &account); err != nil {
		return nil, err
	}

	return &account, nil
}
