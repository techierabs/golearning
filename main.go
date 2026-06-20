package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golearning/config"
	"golearning/internal/client"

	"github.com/gin-gonic/gin"
)

func main() {

	// 1. Read the file from the root directory
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	// 2. Parse (Unmarshal) the YAML data into your struct
	log.Printf("yamlFile: %v", string(yamlFile))

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("[FATAL] Error reading config file: %v", err)
	}

	// Start Apigee OAuth Client Core Engine
	apigeeAuth := client.NewApigeeAuthClient(cfg.Apigee)
	log.Printf("[INFO] Using Apigee Auth Client %v", apigeeAuth)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Test endpoint to trigger and evaluate OAuth Lifecycle execution flow
	r.GET("/stub/authtoken", func(c *gin.Context) {
		token, err := apigeeAuth.GetValidToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":            "Success",
			"token":             token,
			"auth_header_value": fmt.Sprintf("Bearer %s", token),
		})
	})

	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("[SERVER] Token orchestration module active on %s", serverAddr)
	r.Run(serverAddr)
}
