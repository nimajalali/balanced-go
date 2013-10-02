package balanced

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

const (
	cardsUri = marketplaceUri + "/%v/cards"
)

type Card struct {
	Account         Account   `json:"account,omitempty"`
	Brand           string    `json:"brand,omitempty"`
	CanDebit        bool      `json:"can_debit,omitempty"`
	CardNumber      string    `json:"card_number,omitempty"`
	CardType        string    `json:"card_type,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	City            string    `json:"city,omitempty"`
	CountryCode     string    `json:"country_code,omitempty"`
	ExpirationMonth int       `json:"expiration_month,omitempty"`
	ExpirationYear  int       `json:"expiration_year,omitempty"`
	Hash            string    `json:"hash,omitempty"`
	Id              string    `json:"id,omitempty"`
	IsValid         bool      `json:"is_valid,omitempty"`
	LastFour        string    `json:"last_four,omitempty"`
	Meta            MetaType  `json:"meta,omitempty"`
	Name            string    `json:"name,omitempty"`
	PhoneNumber     string    `json:"phone_number,omitempty"`
	PostalCode      string    `json:"postal_code,omitempty"`
	SecurityCode    string    `json:"security_code,omitempty"`
	StreetAddress   string    `json:"street_address,omitempty"`
	Uri             string    `json:"uri,omitempty"`
}

type ListOfCards struct {
	FirstUri    string `json:"first_uri,omitempty"`
	Items       []Card `json:"items,omitempty"`
	LastUri     string `json:"last_uri,omitempty"`
	Limit       int    `json:"limit,omitempty"`
	NextUri     string `json:"next_uri,omitempty"`
	Offset      int    `json:"offset,omitempty"`
	PreviousUri string `json:"previous_uri,omitempty"`
	Total       int    `json:"total,omitempty"`
	Uri         string `json:"uri,omitempty"`
}

// Creates a new card
// WARNING PCI Compliance required to use this functionality.
func TokenizeCard(expirationYear, expirationMonth int, cardNumber, securityCode,
	name, phoneNumber, streetAddress, city, state, postalCode,
	countryCode string, meta MetaType) (card *Card, err error) {

	// Required fields
	payload := url.Values{
		"card_number":      {cardNumber},
		"expiration_year":  {strconv.Itoa(expirationYear)},
		"expiration_month": {strconv.Itoa(expirationMonth)},
	}

	// Add other fields, if not empty
	addToPayload(payload, "security_code", securityCode)
	addToPayload(payload, "name", name)
	addToPayload(payload, "phone_number", phoneNumber)
	addToPayload(payload, "street_address", streetAddress)
	addToPayload(payload, "city", city)
	addToPayload(payload, "state", state)
	addToPayload(payload, "postal_code", postalCode)
	addToPayload(payload, "countryCode", countryCode)
	for key, value := range meta {
		addToPayload(payload, "meta["+key+"]", value)
	}

	uri := fmt.Sprintf(cardsUri, marketplaceId)

	card = &Card{}
	err = post(uri, payload, card)

	return
}

// Retrieves the details of a card that has previously been created. Supply the
// uri that was returned from your previous request, and the corresponding card
// information will be returned. The same information is returned when creating
// the card.
func RetrieveCard(uri string) (card *Card, err error) {
	card = &Card{}
	err = get(uri, nil, card)

	return
}

// Returns a list of cards that you've created.
func ListAllCards(limit, offset int) (*ListOfCards, error) {
	uri := fmt.Sprintf(cardsUri, marketplaceId)

	return ListAllCardsForUri(limit, offset, uri)
}

// Returns a list of cards for a given uri
func ListAllCardsForUri(limit, offset int, uri string) (listOfCards *ListOfCards, err error) {
	payload := defaultPayload(limit, offset)

	listOfCards = &ListOfCards{}
	err = get(uri, payload, listOfCards)

	return
}

// Update information in a card
func UpdateCard(uri string, meta MetaType) (card *Card, err error) {
	payload := url.Values{}

	for key, value := range meta {
		addToPayload(payload, "meta["+key+"]", value)
	}

	card = &Card{}
	err = put(uri, payload, card)

	return
}

// Invalidating a card will mark the card as invalid, so it may not be charged.
func InvalidateCard(uri string) (card *Card, err error) {
	payload := url.Values{
		"is_valid": {"false"},
	}

	card = &Card{}
	err = put(uri, payload, card)

	return
}
