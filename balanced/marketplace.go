package balanced

const (
	marketplaceUri = "/v1/marketplaces"
)

type Marketplace struct {
	CallbacksUri        string   `json:"callbacks_uri,omitempty"`
	SupportEmailAddress string   `json:"support_email_address,omitempty"`
	EventsUri           string   `json:"events_uri,omitempty"`
	AccountsUri         string   `json:"accounts_uri,omitempty"`
	OwnerAccount        Account  `json:"owner_account,omitempty"`
	HoldsUri            string   `json:"holds_uri,omitempty"`
	Meta                MetaType `json:"meta,omitempty"`
	TransactionsUri     string   `json:"transactions_uri,omitempty"`
	BankAccountsUri     string   `json:"bank_accounts_uri,omitempty"`
	Id                  string   `json:"id,omitempty"`
	CreditsUri          string   `json:"credits_uri,omitempty"`
	CardsUri            string   `json:"cards_uri,omitempty"`
	InEscrow            int      `json:"in_escrow,omitempty"`
	DomainUrl           string   `json:"domain_url,omitempty"`
	Name                string   `json:"name,omitempty"`
	Uri                 string   `json:"uri,omitempty"`
	SupportPhoneNumber  string   `json:"support_phone_number,omitempty"`
	RefundsUri          string   `json:"refunds_uri,omitempty"`
	DebitsUri           string   `json:"debits_uri,omitempty"`
}
