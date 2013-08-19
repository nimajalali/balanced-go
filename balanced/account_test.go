package balanced

import (
	"testing"
)

func TestCreateAccount(t *testing.T) {
	account, err := CreateAccount()
	if err != nil {
		t.Fatal(err)
	}

	if len(account.Id) == 0 {
		t.Fatalf("Invalid account created. %v", account)
	}
}

func TestAddCardToAccount(t *testing.T) {
	// Create Account
	account, err := CreateAccount()
	if err != nil {
		t.Fatalf("Unable to create account: %v", err)
	}
	if len(account.Id) == 0 {
		t.Fatalf("Invalid account created. %v", account)
	}

	// Create Card, with minimum requrirements
	card, err := TokenizeCard(2020, 01, testVisaCard, "", "", "", "", "", "", "", "", nil)
	if err != nil {
		t.Fatalf("Unable to create card: %v", err)
	}

	// Add Card to Account
	account, err = AddCardToAccount(account.Uri, card.Uri)
	if err != nil {
		t.Fatalf("Unable to add card to account: %v", err)
	}

	// Get All Cards for Account
	list, err := ListAllCardsForUri(10, 0, account.CardsUri)
	if err != nil {
		t.Fatalf("Unable to get list of cards: %v", err)
	}

	// Verify first card is the card just added
	if len(list.Items) == 0 || list.Items[0].CreatedAt != card.CreatedAt {
		t.Fatal("Unable to add card to account, reading back failed.")
	}
}

func TestUnderwriteIndividual(t *testing.T) {
	merchant := &Merchant{
		PhoneNumber:   "+14089999999",
		Name:          "Timmy Q. CopyPasta",
		Dob:           "1989-12",
		PostalCode:    "94110",
		Type:          accountTypePerson,
		StreetAddress: "121 Skriptkid Row",
	}

	account, err := UnderwriteIndividual(merchant)
	if err != nil {
		t.Fatalf("Unable to underwrite individual: %v", err)
	}

	found := false
	for _, value := range account.Roles {
		if value == merchantRole {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Unable to underwrite individual, missing merchant role")
	}
}

func TestUnderwriteBusiness(t *testing.T) {
	merchant := &Merchant{
		PhoneNumber:   "+140899188155",
		Name:          "Skripts4Kids",
		TaxId:         "211111111",
		PostalCode:    "91111",
		Type:          accountTypeBusiness,
		StreetAddress: "555 VoidMain Road",
	}
	person := &Person{
		Name:          "Timmy Q. CopyPasta",
		Dob:           "1989-12",
		PostalCode:    "94110",
		StreetAddress: "121 Skriptkid Row",
	}

	account, err := UnderwriteBusiness(merchant, person)
	if err != nil {
		t.Fatalf("Unable to underwrite business: %v", err)
	}

	found := false
	for _, value := range account.Roles {
		if value == merchantRole {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Unable to underwrite business, missing merchant role")
	}
}
