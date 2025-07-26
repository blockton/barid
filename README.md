# barid

Go client for [https://api.barid.site/](https://api.barid.site/).  

## Features

- Real-time updates
- Support multiple domains
- Retrieve full message payloads (plain-text & HTML bodies)
  
## Installation

```bash
go get github.com/blockton/barid
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/blockton/barid"
)

func main() {
	temp := barid.GenrateRandomEmail()
	fmt.Println("Created Random Email:", temp.Email)
	domains, err := temp.GetAvailableDomains()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Available Domains (%d):\n", len(domains))
	for _, domain := range domains {
		fmt.Println(domain)
	}

	count, err := temp.GetEmailsCount() // get total received emails count
	if err != nil {
		panic(err)
	}
	fmt.Println("total received emails:", count)
	
	emails, err := temp.GetEmails() // get received emails
	if err != nil {
		panic(err)
	}
	fmt.Printf("Received Emails (%d):\n", len(emails))
	for _, email := range emails {
		fmt.Println("ID:", email.ID)
		fmt.Println("From:", email.From)
		fmt.Println("Subject:", email.Subject)

		message, err := temp.GetEmailInbox(email.ID)
		if err != nil {
			panic(err)
		}

		fmt.Println("ID:", message.ID)
		fmt.Println("TO", message.To)
		fmt.Println("From:", message.From)
		fmt.Println("Subject:", message.Subject)
		fmt.Println("Received:", message.Received)
		fmt.Println("Text Content:")
		fmt.Println(message.TextContent)
		fmt.Println("HTML Content:")
		fmt.Println(message.HTMLContent)

		msg, err := temp.DelEmailInbox(email.ID) // delete email
		if err != nil {
			panic(err)
		}
		fmt.Println("Result:", msg)
	}

	deleted, err := temp.DelEmails() // delete all received emails
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deleted %d emails\n", deleted)
}

```
### Generate a Specified Email

```go
temp := barid.New("dev@wael.fun")
fmt.Println("Created Specified Email:", temp.Email)
```

### Generate a Random Email

```go
temp := barid.GenerateRandomEmail()
fmt.Println("Created Random Email:", temp.Email)
```

### List Available Domains

```go
domains, err := temp.GetAvailableDomains()
if err != nil {
  panic(err)
}
fmt.Printf("Available Domains (%d):\n", len(domains))
for _, d := range domains {
  fmt.Println("-", d)
}
```

### Get Emails Count

```go
count, err := temp.GetEmailsCount()
if err != nil {
  panic(err)
}
fmt.Println("Total received emails:", count)
```

### List Received Emails

```go
emails, err := temp.GetEmails()
if err != nil {
  panic(err)
}
fmt.Printf("Received Emails (%d):\n", len(emails))
for _, e := range emails {
  fmt.Println("ID      :", e.ID)
  fmt.Println("From    :", e.From)
  fmt.Println("Subject :", e.Subject)
  fmt.Println("Received:", e.Received.Format(time.RFC3339))
}
```

### Fetch a Single Email

```go
msg, err := temp.GetEmailInbox("your-email-id")
if err != nil {
  panic(err)
}
fmt.Println("ID         :", msg.ID)
fmt.Println("From       :", msg.From)
fmt.Println("To         :", msg.To)
fmt.Println("Subject    :", msg.Subject)
fmt.Println("Received   :", msg.Received.Format(time.RFC3339))
fmt.Println("Text Body  :\n", msg.TextContent)
fmt.Println("HTML Body  :\n", msg.HTMLContent)
```

### Delete a Single Email

```go
result, err := temp.DelEmailInbox("your-email-id")
if err != nil {
  panic(err)
}
fmt.Println("Delete result:", result)
```

### Delete All Emails

```go
deleted, err := temp.DelEmails()
if err != nil {
  panic(err)
}
fmt.Printf("Deleted %d emails\n", deleted)
```

## Resources

- **API Docs**: https://api.barid.site  
- **Web Client**: https://web.barid.site  
- **Repository**: https://github.com/vwh/temp-mail
