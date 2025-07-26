package barid

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

const BaseURL = "https://api.barid.site"

type Client struct {
	Email          string
	RequestTimeout time.Duration
}

func New(email string) *Client {
	return &Client{
		Email:          email,
		RequestTimeout: 10 * time.Second,
	}
}

func (client *Client) makeRequest(endpoint string) ([]byte, error) {
	requestURL := fmt.Sprintf("%s/%s", BaseURL, endpoint)

	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(response)

	request.SetRequestURI(requestURL)
	request.Header.SetMethod(fasthttp.MethodGet)
	request.Header.SetContentType("application/json")

	httpClient := &fasthttp.Client{}

	if err := httpClient.DoTimeout(request, response, client.RequestTimeout); err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}

	if response.StatusCode() < 200 || response.StatusCode() >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", response.StatusCode(), response.Body())
	}

	return response.Body(), nil
}

func (client *Client) GetEmailInbox(emailID string) (Message, error) {
	endpoint := fmt.Sprintf("inbox/%s", emailID)
	data, err := client.makeRequest(endpoint)
	if err != nil {
		return Message{}, err
	}

	var response struct {
		Success bool       `json:"success"`
		Message RawMessage `json:"result"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return Message{}, err
	}

	if !response.Success {
		return Message{}, fmt.Errorf("failed to get email inbox: %s", string(data))
	}

	return Message{
		ID:          response.Message.ID,
		To:          response.Message.To,
		From:        response.Message.From,
		Subject:     response.Message.Subject,
		Received:    time.Unix(response.Message.Received, 0),
		HTMLContent: response.Message.HTMLContent,
		TextContent: response.Message.TextContent,
	}, nil
}

func (client *Client) GetEmails(limit, offset int) ([]Email, error) {
	endpoint := fmt.Sprintf("emails/%s?limit=%d&offset=%d", client.Email, limit, offset)
	data, err := client.makeRequest(endpoint)
	if err != nil {
		return nil, err
	}

	var response struct {
		Success bool       `json:"success"`
		Emails  []RawEmail `json:"result"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, fmt.Errorf("failed to get emails: %s", string(data))
	}

	emails := make([]Email, len(response.Emails))
	for index, emailData := range response.Emails {
		emails[index] = Email{
			ID:       emailData.ID,
			To:       emailData.To,
			From:     emailData.From,
			Subject:  emailData.Subject,
			Received: time.Unix(emailData.Received, 0),
			_client:  client,
		}
	}

	return emails, nil
}

func (client *Client) GetAvailableDomains() ([]string, error) {
	responseBody, err := client.makeRequest("domains")
	if err != nil {
		return nil, err
	}

	var apiResponse struct {
		Success     bool     `json:"success"`
		DomainsList []string `json:"result"`
	}

	if err := json.Unmarshal(responseBody, &apiResponse); err != nil {
		return nil, err
	}

	if !apiResponse.Success {
		return nil, fmt.Errorf("request failed: %s", string(responseBody))
	}

	return apiResponse.DomainsList, nil
}

func (email *Email) GetInbox() (Message, error) {
	return email._client.GetEmailInbox(email.ID)
}
