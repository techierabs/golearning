package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server      ServerConfig     `yaml:"server"`
	Database    DBConfig         `yaml:"database"`
	Apigee      ApigeeConfig     `yaml:"apigee"`
	Downstreams DownstreamConfig `yaml:"downstreams"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type ApigeeConfig struct {
	AuthURL      string `yaml:"auth_url"`      // e.g., https://t2-prod.apigw.networks.au.singtelgroup.net/oauth/v1/accesstoken
	ClientID     string `yaml:"client_id"`     // Stored Consumer Key
	ClientSecret string `yaml:"client_secret"` // Stored Consumer Secret
}

type DownstreamConfig struct {
	CDIL SystemConfig `yaml:"cdil"`
}

type SystemConfig struct {
	Host      string `yaml:"host"`
	BasePath  string `yaml:"base_path"`
	TimeoutMS int    `yaml:"timeout_ms"`
}

func LoadConfig(path string) (*Config, error) {
	log.Printf("Loading path: %s", path)
	//log.Printf("Loading config from path: %s", &Config)

	var cfg Config
	file, err := os.Open(path)
	log.Printf("Loading file: %s", file)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	log.Printf("decoder %v", decoder)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	log.Printf("&cfg %v", &cfg)
	return &cfg, nil
}
