
package barid

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"io"
	"github.com/go-faster/errors"
)

type apiActions int

const APIBase string = "https://api.barid.site"

const (
	getDomains apiActions = iota

	getEmails
	delEmails
	countMails

	getEmailInbox
	delEmailInbox
)

type API struct {
	Email string
}

func New(email string) *API {
	return &API{
		Email: email,
	}
}

func GenrateRandomEmail() *API {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")

	email := make([]rune, 7)
	for i := range email {
		email[i] = letters[rand.Intn(len(letters))]
	}
	email = append(email, '@')
	email = append(email, []rune("barid.site")...)

	return &API{
		Email: string(email),
	}
}

func (a *API) GetAvailableDomains() ([]string, error) {
	response, err := a.DoRequest(getDomains, nil)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var result APIResponse

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, errors.New("failed to decode domains")
	}
	if !result.Success {
		return nil, errors.New("something went wrong")
	}
	var domains []string

	if err := json.Unmarshal(result.Result, &domains); err != nil {
		return nil, errors.New("failed to decode domains")
	}
	return domains, nil
}

func (a *API) GetEmails() ([]Email, error) {
	response, err := a.DoRequest(getEmails, nil)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var result APIResponse

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, errors.New("failed to decode emails")
	}

	if !result.Success {
		return nil, errors.New("something went wrong")
	}

	var rawEmails []rawEmail
	if err := json.Unmarshal(result.Result, &rawEmails); err != nil {
		return nil, errors.New("failed to decode emails")
	}
	var emails []Email
	for _, rawEmail := range rawEmails {
		emails = append(emails, Email{
			ID:       rawEmail.ID,
			To:       rawEmail.To,
			From:     rawEmail.From,
			Subject:  rawEmail.Subject,
			Received: time.Unix(rawEmail.Received, 0),
		})
	}

	return emails, nil
}

func (a *API) DelEmails() (int, error) {
	response, err := a.DoRequest(delEmails, nil)
	if err != nil {
		return 0, err
	}

	defer response.Body.Close()

	var result struct {
		Success bool `json:"success"`
		Result  struct {
			DeletedCount int `json:"deleted_count"`
		} `json:"result"`
	}

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return 0, errors.New("failed to decode deleted count")
	}

	if !result.Success {
		return 0, errors.New("something went wrong")
	}

	return result.Result.DeletedCount, nil
}

func (a *API) GetEmailsCount() (int, error) {
	response, err := a.DoRequest(countMails, nil)
	if err != nil {
		return 0, err
	}

	defer response.Body.Close()

	var result struct {
		Success bool `json:"success"`
		Result  struct {
			Count int `json:"count"`
		} `json:"result"`
	}

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return 0, errors.New("failed to decode count")
	}

	if !result.Success {
		return 0, errors.New("something went wrong")
	}

	return result.Result.Count, nil
}

func (a *API) GetEmailInbox(emailID string) (*Message, error) {
	response, err := a.DoRequest(getEmailInbox, map[string]string{
		"ID": emailID,
	})
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var result struct {
		Success bool       `json:"success"`
		Result  rawMessage `json:"result"`
	}

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, errors.New("failed to decode message")
	}
	if !result.Success {
		return nil, errors.New("something went wrong")
	}

	return &Message{
		ID:       result.Result.ID,
		To:       result.Result.To,
		From:     result.Result.From,
		Subject:  result.Result.Subject,
		Received: time.Unix(result.Result.Received, 0),

		HTMLContent: result.Result.HTMLContent,
		TextContent: result.Result.TextContent,
	}, nil
}

func (a *API) DelEmailInbox(emailID string) (string, error) {
	response, err := a.DoRequest(delEmailInbox, map[string]string{
		"ID": emailID,
	})
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	var result struct {
		Success bool `json:"success"`
		Result  struct {
			Msg string `json:"message"`
		} `json:"result"`
	}

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return "", errors.New("failed to decode response")
	}

	if !result.Success {
		return "", errors.New("something went wrong")

	}
	return result.Result.Msg, nil
}

func (a *API) DoRequest(action apiActions, args map[string]string) (*http.Response, error) {
	var url string
	var method string

	switch action {
	case getEmails:
		method = "GET"
		url = fmt.Sprintf("%s/emails/%s", APIBase, a.Email)
	case delEmails:
		method = "DELETE"
		url = fmt.Sprintf("%s/emails/%s", APIBase, a.Email)
	case countMails:
		method = "GET"
		url = fmt.Sprintf("%s/emails/count/%s", APIBase, a.Email)
	case getDomains:
		method = "GET"
		url = fmt.Sprintf("%s/domains", APIBase)
	case getEmailInbox:
		method = "GET"
		url = fmt.Sprintf("%s/inbox/%s", APIBase, args["ID"])
	case delEmailInbox:
		method = "DELETE"
		url = fmt.Sprintf("%s/inbox/%s", APIBase, args["ID"])
	default:
		return nil, errors.New("unknown action")
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	r, err := http.NewRequest(method, url, nil)
	r.Header.Set("Content-Type", "application/json")

	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	response, err := client.Do(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request")
	}

	//defer response.Body.Close()

	if response != nil && response.StatusCode != 200 {
		body, _ := io.ReadAll(response.Body)
		return nil, errors.Errorf("[ %d ] - failed to send request:\n%s", response.StatusCode, string(body))
	}

	return response, nil
}
