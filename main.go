package pushnotifier

import (
	"net/http"
)

type client struct {
	BaseURL   string
	UserAgent string

	httpClient *http.Client
}

func NewClient() Client {
	client := client{}
	client.BaseURL = "https://api.pushnotifier.de/v2"
	client.UserAgent = "go-pushnotifier 0.1.0"
	return client
}

func NewClientWithBaseURL(baseUrl string) Client {
	client := client{}
	client.BaseURL = baseUrl
	client.UserAgent = "go-pushnotifier 0.1.0"
	return client
}
