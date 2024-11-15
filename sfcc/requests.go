package sfcc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sfcc/g/kv"
	"sfcc/g/log"
	"sort"
	"strings"
	"time"
)

var apiID string
var apiSecret string
var sfccAuthToken string
var sandboxList []SandboxInfo

type authResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func getSfccAuthToken() (string, error) {
	client := &http.Client{}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", apiID)
	data.Set("client_secret", apiSecret)

	req, err := http.NewRequest(
		http.MethodPost,
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

	log.Trace("Auth token response: %+v\n", authResp)

	return authResp.AccessToken, nil
}

func GetSfccAuthToken(sfccApiID string, sfccApiSecret string) string {
	apiID = sfccApiID
	apiSecret = sfccApiSecret

	value, err := kv.UseCachedResult(
		func() (string, error) { return getSfccAuthToken() },
		"sfccAuthToken",
		29*time.Minute+59*time.Second,
		false,
	)
	if err != nil {
		log.Fatalf("Error getting SFCC auth token: %v\n", err)
	}

	sfccAuthToken = value

	return value
}

func getSandboxListResponse() (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest(
		http.MethodGet,
		"https://admin.dx.commercecloud.salesforce.com/api/v1/sandboxes?include_deleted=false",
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+sfccAuthToken)
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error response: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	return string(body), nil
}

type sbVersions struct {
	App string `json:"app"`
	Web string `json:"web"`
}

type sbLinks struct {
	BM    string `json:"bm"`
	OCAPI string `json:"ocapi"`
	Impex string `json:"impex"`
	Code  string `json:"code"`
	Logs  string `json:"logs"`
}

type SandboxInfo struct {
	ID              string     `json:"id"`
	Realm           string     `json:"realm"`
	Instance        string     `json:"instance"`
	Versions        sbVersions `json:"versions"`
	ResourceProfile string     `json:"resourceProfile"`
	State           string     `json:"state"`
	CreatedAt       string     `json:"createdAt"`
	CreatedBy       string     `json:"createdBy"`
	HostName        string     `json:"hostName"`
	Links           sbLinks    `json:"links"`
}

func GetSandboxList(invalidate bool) []SandboxInfo {
	strBody, err := kv.UseCachedResult(
		func() (string, error) { return getSandboxListResponse() },
		"sandboxList",
		24*time.Hour*14,
		invalidate,
	)

	if err != nil {
		log.Fatalf("Error getting sandbox list: %v", err)
	}

	var response struct {
		Kind   string        `json:"kind"`
		Code   int           `json:"code"`
		Status string        `json:"status"`
		Data   []SandboxInfo `json:"data"`
	}

	err = json.Unmarshal([]byte(strBody), &response)
	if err != nil {
		log.Fatalf("Error unmarshalling sandbox list: %v", err)
	}

	sandboxList = response.Data

	sort.Slice(sandboxList, func(i, j int) bool {
		return sandboxList[i].HostName < sandboxList[j].HostName
	})

	return sandboxList
}
