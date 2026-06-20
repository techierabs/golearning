package client

import (
	"fmt"
	"golearning/config"
	"net/http"
	"time"
)

type CDILClient struct {
	Base *BaseClient
	Cfg  config.SystemConfig
}

func NewCDILClient(sysCfg config.SystemConfig, auth *ApigeeAuthClient) *CDILClient {
	return &CDILClient{
		Cfg: sysCfg,
		Base: &BaseClient{
			AuthClient: auth,
			HTTPClient: &http.Client{Timeout: time.Duration(sysCfg.TimeoutMS) * time.Millisecond},
		},
	}
}

func (c *CDILClient) PostProfile(payload []byte) (*http.Response, error) {
	fullURL := fmt.Sprintf("%s%s/profile/channel/HELIX", c.Cfg.Host, c.Cfg.BasePath) // [cite: 63]
	return c.Base.DoRequest("POST", fullURL, payload)                                // [cite: 62]
}
