package client

import (
	"fmt"
	"golearning/config"
	"net/http"
	"time"
)

type EDA2Client struct {
	Base *BaseClient
	Cfg  config.SystemConfig
}

func NewEDA2Client(sysCfg config.SystemConfig, auth *ApigeeAuthClient) *EDA2Client {
	return &EDA2Client{
		Cfg: sysCfg,
		Base: &BaseClient{
			AuthClient: auth,
			HTTPClient: &http.Client{Timeout: time.Duration(sysCfg.TimeoutMS) * time.Millisecond},
		},
	}
}

func (e *EDA2Client) ProcessEvent(payload []byte) (*http.Response, error) {
	fullURL := fmt.Sprintf("%s%s/events", e.Cfg.Host, e.Cfg.BasePath)
	return e.Base.DoRequest("POST", fullURL, payload)
}
