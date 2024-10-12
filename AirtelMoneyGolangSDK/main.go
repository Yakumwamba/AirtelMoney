package main

import (
	"fmt"
	"log"

	"github.com/yakumwamba/airtel-money-project/airtelmoney"
)

func main() {
	client := airtelmoney.NewClient("YOUR_CLIENT_ID", "YOUR_CLIENT_SECRET", airtelmoney.StagingEnvironment)

	// Authenticate
	token, err := client.Authenticate()
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	fmt.Printf("Authentication successful. Access Token: %s\n", token.AccessToken)

	// Example disbursement
	disbursement := &airtelmoney.Disbursement{
		Payee: airtelmoney.Payee{
			MSISDN:     "75****26",
			WalletType: "NORMAL",
		},
		Reference: "TestDisbursement",
		Pin:       "EncryptedPIN",
		Transaction: airtelmoney.Transaction{
			Amount: "1000",
			ID:     "TEST123",
			Type:   "B2C",
		},
	}

	response, err := client.MakeDisbursement(disbursement)
	if err != nil {
		log.Fatalf("Disbursement failed: %v", err)
	}

	fmt.Printf("Disbursement successful. Transaction ID: %s\n", response.Data.Transaction.ID)
}
