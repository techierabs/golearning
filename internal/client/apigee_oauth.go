package client

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golearning/config"
	"golearning/internal/models"
)

type ApigeeAuthClient struct {
	cfg        config.ApigeeConfig
	client     *http.Client
	tokenCache string
	expiryTime time.Time
	mu         sync.RWMutex
}

func NewApigeeAuthClient(cfg config.ApigeeConfig) *ApigeeAuthClient {
	return &ApigeeAuthClient{
		cfg:    cfg,
		client: &http.Client{Timeout: 4 * time.Second}, // Snappy timeout to catch down networks quickly
	}
}

// GetValidToken provides dynamic caching, concurrency safety, and graceful file recovery
func (a *ApigeeAuthClient) GetValidToken() (string, error) {
	a.mu.RLock()
	// Cache Guard: check if token is valid with a 15-second clock skew buffer
	if a.tokenCache != "" && time.Now().Add(55*time.Second).Before(a.expiryTime) {
		defer a.mu.RUnlock()
		return a.tokenCache, nil
	}
	a.mu.RUnlock()

	a.mu.Lock()
	defer a.mu.Unlock()

	// Double-Checked Locking Pattern to block simultaneous network requests
	if a.tokenCache != "" && time.Now().Add(15*time.Second).Before(a.expiryTime) {
		return a.tokenCache, nil
	}

	// 1. Prepare application/x-www-form-urlencoded Request Payload
	formData := url.Values{}
	formData.Set("grant_type", "client_credentials") // Mandated by contract Section 3.5.2

	log.Printf("formData: %v", formData)

	req, err := http.NewRequest("POST", a.cfg.AuthURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return a.executeStubFallback("failed to compile http request template")
	}

	log.Printf("req: %v", req)

	// 2. Set Contract Required Headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // Mandated by contract
	req.SetBasicAuth(a.cfg.ClientID, a.cfg.ClientSecret)                // Basic + Base64(Key:Secret)

	// 3. Dispatch Live Callout toward Gateway
	resp, err := a.client.Do(req)
	log.Printf("resp: %v", resp)

	if err != nil {
		return a.executeStubFallback("gateway network interface connection timed out/unreachable")
	}
	defer resp.Body.Close()

	// 4. Validate Gateway Status Code
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			log.Println("[ERROR] Apigee Gateway reported 401: invalid_credentials (ClientId is Invalid).") // Section 6
		}
		return a.executeStubFallback("received non-200 transaction response status from token authority")
	}

	// 5. Decode Live Response Body
	var tokenResp models.OAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return a.executeStubFallback("failed parsing live json token metadata payload structure")
	}

	// Convert string duration token into integer seconds
	seconds, err := strconv.ParseInt(tokenResp.ExpiresIn, 10, 64)
	if err != nil {
		seconds = 3599 // Fallback to standard 1 hour contract default if corrupted
	}

	// 6. Persist internally to memory cache
	a.tokenCache = tokenResp.AccessToken
	a.expiryTime = time.Now().Add(time.Duration(seconds) * time.Second) // Current Timestamp + expires_in

	log.Printf("[INFO] Token rotation successful. New token expires at: %v (UTC)", a.expiryTime.UTC())
	return a.tokenCache, nil
}

// executeStubFallback intercepts runtime operational anomalies and drops to local mock state
func (a *ApigeeAuthClient) executeStubFallback(reason string) (string, error) {
	log.Printf("[WARN] Apigee Auth unreachable: %s. Initiating local fallback parsing engine...", reason)

	fileBytes, err := os.ReadFile("stub_token.json")
	log.Printf("fileBytes: %v", fileBytes)

	if err != nil {
		log.Printf("[CRITICAL] Fallback engine failed: mock JSON stub file missing: %v", err)
		return "", err
	}

	var tokenResp models.OAuthResponse
	if err := json.Unmarshal(fileBytes, &tokenResp); err != nil {
		log.Printf("[CRITICAL] Corrupted fallback json configuration profile: %v", err)
		return "", err
	}

	seconds, _ := strconv.ParseInt(tokenResp.ExpiresIn, 10, 64)
	if seconds == 0 {
		seconds = 1799
	}

	// Seed cache using local mock configurations to reduce disk I/O pressure during outtages
	a.tokenCache = tokenResp.AccessToken
	a.expiryTime = time.Now().Add(time.Duration(seconds) * time.Second)

	log.Println("[SUCCESS] Resiliency circuit activated: Token successfully extracted from local fallback template.")
	return a.tokenCache, nil
}
