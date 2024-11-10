package sfcc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sfcc/g/kv"
	"strings"
	"time"
)

type authResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func getSfccAuthToken(apiID string, apiSecret string) (string, error) {
	client := &http.Client{}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", apiID)
	data.Set("client_secret", apiSecret)

	req, err := http.NewRequest(
		"POST",
		"https://account.demandware.com:443/dwsso/oauth2/access_token",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	var authResp authResponse
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	if err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	fmt.Printf("Auth token response: %+v\n", authResp)

	return authResp.AccessToken, nil
}

func GetSfccAuthToken(apiID string, apiSecret string) string {
	value, err := kv.UseCachedResult(
		func() (string, error) { return getSfccAuthToken(apiID, apiSecret) },
		"sfccAuthToken",
		29*time.Minute+59*time.Second,
		false,
	)
	if err != nil {
		fmt.Errorf("Error getting SFCC auth token: %w", err)
		os.Exit(1)
	}

	return value
}
