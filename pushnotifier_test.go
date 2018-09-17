package pushnotifier

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

var (
	mux    *http.ServeMux
	client *Client
	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	serverUrl, _ := url.Parse(server.URL)

	client = NewClientWithBaseURL(serverUrl.String() + "/" + apiVersion)
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if want != r.Method {
		t.Errorf("Request method = %v, want %v", r.Method, want)
	}
}

func TestListDevices(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/devices",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w, `
				[
					{
						"id":    "ABC",
						"title": "example@example.org",
						"model": "E-Mail",
						"image": "https://devices.pushnotifier.de/virtual/E-Mail.png"
					},
					{
						"id":    "xYz",
						"title": "iPhone X",
						"model": "iPhone10,6",
						"image": "https://devices.pushnotifier.de/apple/iPhoneX.png"
					}
				]
			`)
		})

	devices, err := client.ListDevices()
	if err != nil {
		t.Errorf("ListDevices() returned error: %v", err)
	}

	want := []Device{
		Device{"ABC", "example@example.org", "E-Mail", "https://devices.pushnotifier.de/virtual/E-Mail.png"},
		Device{"xYz", "iPhone X", "iPhone10,6", "https://devices.pushnotifier.de/apple/iPhoneX.png"},
	}
	if !reflect.DeepEqual(devices, want) {
		t.Errorf("ListDevices() returned %+v, want %+v",
			devices, want)
	}
}
