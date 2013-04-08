package balanced

import (
	"testing"
)

func TestCredit(t *testing.T) {
	credit := creditNewBankAccount(t)
	t.Log(credit)
}

func creditNewBankAccount(t *testing.T) *Credit {
	// Create Card

	// Capture Debit

	// Credit Bank
	bankAccount := &BankAccount{
		Name:          "Johann Bernoulli",
		AccountNumber: "9900000001",
		RoutingNumber: "121000358",
		Type:          BankAccountTypeChecking,
	}

	credit, err := CreditNewBankAccount(50, "Testing credit", bankAccount)
	if err != nil {
		t.Fatalf("Failed to create credit: %v", err)
	}

	if len(credit.CreatedAt.String()) == 0 {
		t.Fatal("Failed to create credit.")
	}

	return credit
}
