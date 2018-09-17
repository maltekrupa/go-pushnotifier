package pushnotifier

import (
	"net/http"
)

type client struct {
	BaseURL   string
	UserAgent string

	httpClient *http.Client
}

func NewClient(baseUrl string) Client {
	client := client{}

	if baseUrl == nil {
		client.BaseURL = "https://api.pushnotifier.de/v2"
	} else {
		client.BaseURL = baseUrl
	}
	client.UserAgent = "go-pushnotifier 0.1.0"
	return client
}

func (c *Client) ListDevices() ([]Device, error) {
	req, err := c.newRequest("GET", "/devices", nil)
	if err != nil {
		return nil, err
	}
	var devices []Device
	_, err = c.do(req, &devices)
	return devices, err
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}