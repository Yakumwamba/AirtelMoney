package airtelmoney

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Environment string

const (
	StagingEnvironment    Environment = "https://openapiuat.airtel.africa"
	ProductionEnvironment Environment = "https://openapi.airtel.africa"
)

type Client struct {
	ClientID     string
	ClientSecret string
	BaseURL      Environment
	HTTPClient   *http.Client
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type Payee struct {
	MSISDN     string `json:"msisdn"`
	WalletType string `json:"wallet_type"`
}

type Transaction struct {
	Amount string `json:"amount"`
	ID     string `json:"id"`
	Type   string `json:"type,omitempty"`
}

type Disbursement struct {
	Payee       Payee       `json:"payee"`
	Reference   string      `json:"reference"`
	Pin         string      `json:"pin"`
	Transaction Transaction `json:"transaction"`
}

type DisbursementResponse struct {
	Data struct {
		Transaction struct {
			ReferenceID   string `json:"reference_id"`
			AirtelMoneyID string `json:"airtel_money_id"`
			ID            string `json:"id"`
			Status        string `json:"status"`
			Message       string `json:"message"`
		} `json:"transaction"`
	} `json:"data"`
	Status struct {
		ResponseCode string `json:"response_code"`
		Code         string `json:"code"`
		Success      bool   `json:"success"`
		Message      string `json:"message"`
	} `json:"status"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error: %s (Code: %s)", e.Message, e.Code)
}

func NewClient(clientID, clientSecret string, env Environment) *Client {
	return &Client{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		BaseURL:      env,
		HTTPClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) Authenticate() (*AuthResponse, error) {
	url := fmt.Sprintf("%s/auth/oauth2/token", c.BaseURL)
	payload := map[string]string{
		"client_id":     c.ClientID,
		"client_secret": c.ClientSecret,
		"grant_type":    "client_credentials",
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiError APIError
		json.NewDecoder(resp.Body).Decode(&apiError)
		return nil, &apiError
	}

	var authResp AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	if err != nil {
		return nil, err
	}

	return &authResp, nil
}

func (c *Client) MakeDisbursement(d *Disbursement) (*DisbursementResponse, error) {
	url := fmt.Sprintf("%s/standard/v3/disbursements", c.BaseURL)
	jsonPayload, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	// TODO: Add required headers (Authorization, X-Country, X-Currency, etc.)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiError APIError
		json.NewDecoder(resp.Body).Decode(&apiError)
		return nil, &apiError
	}

	var disbResp DisbursementResponse
	err = json.NewDecoder(resp.Body).Decode(&disbResp)
	if err != nil {
		return nil, err
	}

	return &disbResp, nil
}
