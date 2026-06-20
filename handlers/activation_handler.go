package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"database"
	"downstream"

	"github.com/gin-gonic/gin"
)

type ActivationHandler struct {
	Provider downstream.DownstreamProvider
	DBLogger *database.TransactionLogger
}

func NewActivationHandler(p downstream.DownstreamProvider, db *database.TransactionLogger) *ActivationHandler {
	return &ActivationHandler{Provider: p, DBLogger: db}
}

func (h *ActivationHandler) CreateProfileHandler(c *gin.Context) {
	var req downstream.CDILRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[WARN] Validation error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"responseCode": "1001", "responseDesc": "Invalid request parameter"}) // [cite: 102]
		return
	}

	reqBytes, _ := json.Marshal(req)
	resp, statusCode, err := h.Provider.CreateProfile(req)

	var respBytes []byte
	if resp != nil {
		respBytes, _ = json.Marshal(resp)
	}

	// Logging transaction asynchronously to Postgres to preserve low response latencies
	go h.DBLogger.LogTransaction(h.Provider.GetName(), "CreateProfile", string(reqBytes), string(respBytes), statusCode)

	if err != nil {
		c.JSON(statusCode, gin.H{"responseCode": "1000", "responseDesc": err.Error()}) // [cite: 102]
		return
	}

	c.JSON(statusCode, resp)
}

func (h *ActivationHandler) UpdateProfileHandler(c *gin.Context) {
	var req downstream.CDILRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"responseCode": "1001", "responseDesc": "Invalid request parameter"}) // [cite: 102]
		return
	}

	reqBytes, _ := json.Marshal(req)
	resp, statusCode, err := h.Provider.UpdateProfile(req)

	var respBytes []byte
	if resp != nil {
		respBytes, _ = json.Marshal(resp)
	}

	go h.DBLogger.LogTransaction(h.Provider.GetName(), "UpdateProfile", string(reqBytes), string(respBytes), statusCode)

	if err != nil {
		c.JSON(statusCode, gin.H{"responseCode": "1000", "responseDesc": err.Error()}) // [cite: 102]
		return
	}

	c.JSON(statusCode, resp)
}
