package balanced

import (
	"testing"
)

const (
	testVisaCard          = "4111111111111111"
	testVisaSecCode       = "123"
	testAmexCard          = "341111111111111"
	testAmexSecCode       = "1234"
	testMasterCard        = "5105105105105100"
	testMasterSecCode     = "123"
	testCancelledCard     = "4222222222222220"
	testInsufficientFunds = "4444444444444448"
)

func TestCard(t *testing.T) {
	card := tokenizeCard(t)
	card = retrieveCard(t, card)

	listAllCards(t, card)
	updateCard(t, card)
	invalidateCard(t, card)
}

func tokenizeCard(t *testing.T) *Card {
	card, err := TokenizeCard(2020, 12, testVisaCard, testVisaSecCode,
		"Peter Sherman", "1234567890", "42 Wallaby Way", "Sydney", "CA", "92617",
		"US", nil)

	if err != nil {
		t.Fatalf("Failed to tokenize card: %v", err)
	}

	if len(card.Uri) == 0 {
		t.Fatal("Failed to tokenize card")
	}

	return card
}

func retrieveCard(t *testing.T, c *Card) *Card {
	card, err := RetrieveCard(c.Uri)
	if err != nil {
		t.Fatalf("Failed to retrieve card: %v", err)
	}

	if len(card.Uri) == 0 {
		t.Fatal("Failed to retrieve card")
	}

	return card
}

func listAllCards(t *testing.T, c *Card) {
	list, err := ListAllCards(10, 0)
	if err != nil {
		t.Fatalf("Failed to retrieve list of cards: %v", err)
	}

	if len(list.Items) == 0 {
		t.Fatal("Invalid list of cards.")
	}
}

func updateCard(t *testing.T, c *Card) {
	meta := MetaType{
		"testKey": "testValue",
	}

	card, err := UpdateCard(c.Uri, meta)
	if err != nil {
		t.Fatalf("Failed to update card: %v", err)
	}

	if len(card.Meta) == 0 {
		t.Fatal("Failed to update card")
	}
}

func invalidateCard(t *testing.T, c *Card) {
	card, err := InvalidateCard(c.Uri)
	if err != nil {
		t.Fatalf("Failed to invalidate card: %v", err)
	}

	if card.IsValid != false {
		t.Fatal("failed to invalidate card")
	}
}
