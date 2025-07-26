package main

import (
	"time"

	"github.com/blockton/barid"
)

func main() {
	barid := barid.New("tests@wael.fun")

	emailSummaries, err := barid.GetEmails(10, 0)
	if err != nil {
		println(err)
		return
	}

	if len(emailSummaries) == 0 {
		println("No emails found")
		return
	}

	firstEmail := emailSummaries[0]
	println(firstEmail.ID)
	println(firstEmail.To)
	println(firstEmail.From)
	println(firstEmail.Subject)
	println(firstEmail.Received.Format(time.RFC3339))

	availableDomains, err := barid.GetAvailableDomains()
	if err != nil {
		println("error", err)
	}

	println(availableDomains)

	inbox, err := firstEmail.GetInbox()
	if err != nil {
		println(err)
		return
	}

	println(inbox.ID)
	println(inbox.To)
	println(inbox.From)
	println(inbox.Subject)
	println(inbox.Received.Format(time.RFC3339))

	println(inbox.HTMLContent)
	println(inbox.TextContent)
}
