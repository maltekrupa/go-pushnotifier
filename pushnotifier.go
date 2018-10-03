package pushnotifier

import (
	"errors"
	"strings"

	"gopkg.in/resty.v1"
)

var apiVersion = "v2"
var baseURL = "https://api.pushnotifier.de/" + apiVersion
var userAgent = "go-pushnotifier 0.1.0"

type Client struct {
	Http      resty.Client
	BaseURL   string
	UserAgent string
}

type AuthError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type SendText struct {
	Devices []string `json:"devices"`
	Content string   `json:"content"`
}

type SendResponse struct {
	Success []string `json:"success"`
	Error   []string `json:"error"`
}

type Device struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Model string `json:"model"`
	Image string `json:"image"`
}

func NewClient() *Client {
	r := resty.New()
	r.SetError(AuthError{})
	r.SetHeaders(map[string]string{
		"User-Agent": userAgent,
	})

	return &Client{*r, baseURL, userAgent}
}

func (c Client) ListDevices() ([]Device, error) {
	var url strings.Builder
	url.WriteString(c.BaseURL)
	url.WriteString("/devices")

	var d []Device
	r, _ := c.Http.R().
		SetResult(&d).
		Get(url.String())

	if r.Error() != nil {
		msg := r.Error().(*AuthError)
		return nil, errors.New(msg.Message)
	}

	return d, nil
}
