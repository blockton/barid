package barid

import (
	"encoding/json"
	"time"
)

type APIResponse struct {
	Success bool            `json:"success"`
	Result  json.RawMessage `json:"result"`
}

type rawEmail struct {
	ID       string `json:"id"`
	To       string `json:"to_address"`
	From     string `json:"from_address"`
	Subject  string `json:"subject"`
	Received int64  `json:"received_at"`
}

type rawMessage struct {
	ID       string `json:"id"`
	To       string `json:"to_address"`
	From     string `json:"from_address"`
	Subject  string `json:"subject"`
	Received int64  `json:"received_at"`

	HTMLContent string `json:"html_content"`
	TextContent string `json:"text_content"`
}

type Email struct {
	ID       string    `json:"id"`
	To       string    `json:"to_address"`
	From     string    `json:"from_address"`
	Subject  string    `json:"subject"`
	Received time.Time `json:"received_at"`
}

type Message struct {
	ID       string    `json:"id"`
	To       string    `json:"to_address"`
	From     string    `json:"from_address"`
	Subject  string    `json:"subject"`
	Received time.Time `json:"received_at"`

	HTMLContent string `json:"html_content"`
	TextContent string `json:"text_content"`
}

// todo: add attachments field