package pushnotifier

import (
	"errors"
	"os"
	"strconv"
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

	AppToken string
}

type AuthError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Username  string `json:"username"`
	Avatar    string `json:"avatar"`
	AppToken  string `json:"app_token"`
	ExpiresAt string `json:"expires_at"`
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
	debug, _ := strconv.ParseBool(os.Getenv("PUSHNOTIFIER_DEBUG"))
	pkg := os.Getenv("PUSHNOTIFIER_PACKAGE")
	token := os.Getenv("PUSHNOTIFIER_TOKEN")
	username := os.Getenv("PUSHNOTIFIER_USERNAME")
	password := os.Getenv("PUSHNOTIFIER_PASSWORD")

	r := resty.New()
	r.SetError(AuthError{})
	r.SetDebug(debug)
	r.SetBasicAuth(pkg, token)
	r.SetHeaders(map[string]string{
		"User-Agent": userAgent,
	})

	c := &Client{*r, baseURL, userAgent, ""}
	c.login(username, password)

	return c
}

func (c *Client) login(username, password string) {
	var url strings.Builder
	url.WriteString(c.BaseURL)
	url.WriteString("/user/login")

	var l LoginResponse
	r, _ := c.Http.R().
		SetBody(LoginRequest{Username: username, Password: password}).
		SetResult(&l).
		Post(url.String())

	if r.Error() != nil {
		msg := r.Error().(*AuthError)
		errors.New(msg.Message)
	}

	c.AppToken = l.AppToken
}

func (c Client) ListDevices() ([]Device, error) {
	var url strings.Builder
	url.WriteString(c.BaseURL)
	url.WriteString("/devices")

	var d []Device
	r, _ := c.Http.R().
		SetHeader("X-AppToken", c.AppToken).
		SetResult(&d).
		Get(url.String())

	if r.Error() != nil {
		msg := r.Error().(*AuthError)
		return nil, errors.New(msg.Message)
	}

	return d, nil
}

func (c Client) SendText(devices []string, message string) (SendResponse, error) {
	var url strings.Builder
	url.WriteString(c.BaseURL)
	url.WriteString("/notifications/text")

	var s SendResponse
	r, _ := c.Http.R().
		SetHeader("X-AppToken", c.AppToken).
		SetBody(SendText{Devices: devices, Content: message}).
		SetResult(&s).
		Put(url.String())

	if r.StatusCode() == 400 {
		return s, errors.New("Malformed request")
	}
	if r.StatusCode() == 404 {
		return s, errors.New("One of the devices couldn't be found")
	}

	return s, nil
}
