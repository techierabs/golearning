package client

import (
	"bytes"
	"net/http"
)

type BaseClient struct {
	HTTPClient *http.Client
	AuthClient *ApigeeAuthClient
}

// Execute restful calls injecting Apigee OAuth token systematically
func (b *BaseClient) DoRequest(method, url string, body []byte) (*http.Response, error) {
	token, err := b.AuthClient.GetValidToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json") // [cite: 64, 82]
	req.Header.Set("Authorization", "Bearer "+token)

	return b.HTTPClient.Do(req)
}
