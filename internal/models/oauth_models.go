package models

// OAuthResponse reflects the exact NaaS Apigee contract attributes (Section 3.5.2)
type OAuthResponse struct {
	AccessToken    string   `json:"access_token"`          // The Bearer token value
	TokenType      string   `json:"token_type"`            // E.g., BearerToken
	ExpiresIn      string   `json:"expires_in"`            // Remaining seconds (Returned as String)
	IssuedAt       string   `json:"issued_at"`             // Epoch timestamp generated
	ClientID       string   `json:"client_id"`             // Consumer Key
	Organization   string   `json:"organization_name"`     // Apigee Org
	ApiProductList []string `json:"api_product_list_json"` // Products authorized
	Status         string   `json:"status"`                // Approval status
}
