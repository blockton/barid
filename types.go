package barid

import "time"

type RawEmail struct {
	ID       string `json:"id"`
	To       string `json:"to_address"`
	From     string `json:"from_address"`
	Subject  string `json:"subject"`
	Received int64  `json:"received_at"`
}

type RawMessage struct {
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
	_client  *Client   //
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
