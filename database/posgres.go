package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type TransactionLogger struct {
	DB *sql.DB
}

func InitDB(host string, port int, user, password, dbname, sslmode string) (*TransactionLogger, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Create log table if not exists
	query := `
	CREATE TABLE IF NOT EXISTS api_transaction_logs (
		id UUID PRIMARY KEY,
		interface_name VARCHAR(50),
		operation VARCHAR(50),
		request_payload TEXT,
		response_payload TEXT,
		http_status INT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(query); err != nil {
		return nil, err
	}

	return &TransactionLogger{DB: db}, nil
}

func (tl *TransactionLogger) LogTransaction(interf, op, req, resp string, status int) {
	query := `INSERT INTO api_transaction_logs (id, interface_name, operation, request_payload, response_payload, http_status, created_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := tl.DB.Exec(query, uuid.New(), interf, op, req, resp, status, time.Now())
	if err != nil {
		log.Printf("[ERROR] Failed to write log to Postgres: %v", err)
	}
}
