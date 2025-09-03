package collector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	BaseURL      string
	Username     string
	Password     string
	ClientID     string
	ClientSecret string
	AccessToken  string
	HTTPClient   *http.Client
}

// NewClient crea un cliente con los datos necesarios para autenticarse
func NewClient(baseURL, username, password, clientID, clientSecret string) *Client {
	return &Client{
		BaseURL:      baseURL,
		Username:     username,
		Password:     password,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		HTTPClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

// Login hace POST al endpoint /legacy/Api/access_token y guarda el token en memoria
func (c *Client) Login() error {
	url := c.BaseURL + "/legacy/Api/access_token"

	payload := map[string]string{
		"grant_type":    "password",
		"client_id":     c.ClientID,
		"client_secret": c.ClientSecret,
		"username":      c.Username,
		"password":      c.Password,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error al armar JSON: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error al crear request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error al ejecutar request: %w", err)
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(data))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return fmt.Errorf("error al parsear JSON: %w", err)
	}

	token, ok := result["access_token"].(string)
	if !ok {
		return fmt.Errorf("no access_token en respuesta: %v", result)
	}

	c.AccessToken = token
	return nil
}

// FetchAccounts hace GET a /legacy/Api/V8/module/Accounts usando el token en memoria
func (c *Client) FetchAccounts() ([]byte, error) {
	url := c.BaseURL + "/legacy/Api/V8/module/Accounts"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch accounts failed with status %d: %s", resp.StatusCode, string(data))
	}

	return data, nil
}
