package barid

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-faster/errors"
)

type apiActions int

const APIBase string = "https://api.barid.site"

const (
	GetDomains apiActions = iota

	GetEmails
	DelEmails
	CountMails

	GetEmailInbox
	DelEmailInbox
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
	response, err := a.DoRequest(GetDomains, nil)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	domains := make([]string, 0)
	if err := json.NewDecoder(response.Body).Decode(&domains); err != nil {
		return nil, errors.New("failed to decode domains")
	}
	return domains, nil
}

func (a *API) DoRequest(action apiActions, args map[string]string) (*http.Response, error) {
	var url string
	var method string

	switch action {
	case GetEmails:
		method = "GET"
		url = fmt.Sprintf("%s/emails/%s", APIBase, a.Email)
	case DelEmails:
		method = "DELETE"
		url = fmt.Sprintf("%s/emails/%s", APIBase, a.Email)
	case CountMails:
		method = "GET"
		url = fmt.Sprintf("%s/emails/count/%s", APIBase, a.Email)
	case GetDomains:
		method = "GET"
		url = fmt.Sprintf("%s/domains", APIBase)
	case GetEmailInbox:
		method = "GET"
		url = fmt.Sprintf("%s/inbox/%s", APIBase, args["ID"])
	case DelEmailInbox:
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

	if err != nil || (response != nil && response.StatusCode != 200) {
		return nil, errors.Errorf("[ %d ] - failed to send request", response.StatusCode)
	}

	return response, nil
}
