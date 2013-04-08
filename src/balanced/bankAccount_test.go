package balanced

import (
	"testing"
)

func TestBankAccount(t *testing.T) {
	// Test Bank Account Creation
	bankAccount := createNewBankAccount(t)

	// Test Retrieving Bank Accounts
	bankAccount = retrieveBankAccount(t, bankAccount)
	listAllBankAccounts(t, bankAccount)

	// Test Bank Account Verifications
	verification := verifyBankAccount(t, bankAccount)
	retrieveBankAccountVerification(t, verification)
	listAllBankACcountVerifications(t, bankAccount, verification)

	// Test Deleting a Bank Account
	deleteBankAccount(t, bankAccount)
}

func createNewBankAccount(t *testing.T) *BankAccount {
	bankAccount, err := CreateNewBankAccount(
		"Johann Bernoulli",
		"9900000001",
		"121000358",
		BankAccountTypeChecking)

	if err != nil {
		t.Fatalf("Failed to create bank account: %v", err)
	}

	if len(bankAccount.Uri) == 0 {
		t.Fatalf("Invalid Bank Account created: %v", bankAccount)
	}

	return bankAccount
}

func retrieveBankAccount(t *testing.T, bankAccount *BankAccount) *BankAccount {
	retrievedAccount, err := RetrieveBankAccount(bankAccount.Uri)
	if err != nil {
		t.Fatalf("Failed to retrieve bank account: %v", err)
	}

	if retrievedAccount.Uri != bankAccount.Uri {
		t.Fatalf("Incorrect bank account retrieved")
	}

	return retrievedAccount
}

func listAllBankAccounts(t *testing.T, bankAccount *BankAccount) {
	list, err := ListAllBankAccounts(10, 0)
	if err != nil {
		t.Fatalf("Failed to retrieve list of bank accounts: %v", err)
	}

	// Verify first bank account is the account just added
	if list.Items[0].CreatedAt != bankAccount.CreatedAt {
		t.Fatal("Unable to add bank account, reading back failed.")
	}
}

func deleteBankAccount(t *testing.T, bankAccount *BankAccount) {
	err := DeleteBankAccount(bankAccount.Uri)
	if err != nil {
		t.Fatalf("Failed to delete bank account: %v", err)
	}

	// Verify bank account was deleted
	list, err := ListAllBankAccounts(10, 0)
	if len(list.Items) == 0 || list.Items[0].CreatedAt == bankAccount.CreatedAt {
		t.Fatal("Unable to delete bank account, reading back failed.")
	}
}

func verifyBankAccount(t *testing.T, bankAccount *BankAccount) *Verification {
	verification, err := VerifyBankAccount(bankAccount.Uri)
	if err != nil {
		t.Fatalf("Failed to generate verification for bank account: ", err)
	}

	if len(verification.Id) == 0 {
		t.Fatalf("Invalid verification for Bank Account created: ",
			verification)
	}

	return verification
}

func retrieveBankAccountVerification(t *testing.T, v *Verification) {
	retrievedVerification, err := RetrieveBankAccountVerification(v.Uri)
	if err != nil {
		t.Fatalf("Failed to retrieve verification for bank account: ", err)
	}

	if retrievedVerification.Id != v.Id {
		t.Fatalf("Invalid verification for Bank Account retrieved: ",
			retrievedVerification)
	}
}

func listAllBankACcountVerifications(t *testing.T, b *BankAccount,
	v *Verification) {

	list, err := ListAllBankAccountVerifications(b.VerificationsUri)
	if err != nil {
		t.Fatalf("Unable to get list of bank account verifications: %v", err)
	}

	// Verify first verification is the verification just added
	if len(list.Items) == 0 || list.Items[0].Id != v.Id {
		t.Fatal("Invalid verification for Bank Account, reading back failed.")
	}
}

func confirmBankAccountVerification(t *testing.T, v *Verification) {
	confirmedVerification, err := ConfirmBankAccountVerification(v.Uri, 1, 1)
	if err != nil {
		t.Fatalf("Failed to confirm verification for bank account: ", err)
	}

	if confirmedVerification.State != VerificationTypeVerified {
		t.Fatal("Invalid confirmation of verification. Statys is not verified")
	}
}
