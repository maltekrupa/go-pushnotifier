package pushnotifier

import (
	"os"
	"strconv"
	"testing"

	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v1"
)

var pkg, token, username, password string
var debug bool

func TestMain(m *testing.M) {
	os.Setenv("PUSHNOTIFIER_DEBUG", "false")
	debug, _ = strconv.ParseBool(os.Getenv("PUSHNOTIFIER_DEBUG"))

	os.Setenv("PUSHNOTIFIER_PACKAGE", "com.test.foo")
	pkg = os.Getenv("PUSHNOTIFIER_PACKAGE")

	os.Setenv("PUSHNOTIFIER_TOKEN", "TOKEN1234")
	token = os.Getenv("PUSHNOTIFIER_TOKEN")

	os.Setenv("PUSHNOTIFIER_USERNAME", "username_foo")
	username = os.Getenv("PUSHNOTIFIER_USERNAME")

	os.Setenv("PUSHNOTIFIER_PASSWORD", "password_bar")
	password = os.Getenv("PUSHNOTIFIER_PASSWORD")

	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

func TestRestyClient(t *testing.T) {
	r := SetupHttpClient(token, pkg, debug)

	st.Expect(t, r.Debug, debug)
	// We send the package and token as basic auth values
	st.Expect(t, r.UserInfo.Username, pkg)
	st.Expect(t, r.UserInfo.Password, token)
	st.Expect(t, r.Header["User-Agent"][0], userAgent)
}

func TestClient(t *testing.T) {
	c := NewClient(username, password, token, pkg, debug)

	st.Expect(t, c.BaseURL, "https://"+domain+apiVersion)
	st.Expect(t, c.UserAgent, userAgent)
	st.Expect(t, c.Username, username)
	st.Expect(t, c.Password, password)
}

func TestClientFromEnv(t *testing.T) {
	c := NewClientFromEnv()

	st.Expect(t, c.BaseURL, "https://"+domain+apiVersion)
	st.Expect(t, c.UserAgent, userAgent)
	st.Expect(t, c.Username, username)
	st.Expect(t, c.Password, password)
}

func TestLogin(t *testing.T) {
	defer gock.Off()

	gock.New("https://" + domain).
		Post(apiVersion + "/user/login").
		Reply(200).
		JSON(map[string]string{
			"username":   "foo",
			"avatar":     "bar",
			"app_token":  "1234",
			"expires_at": "5678",
		})

	c := NewClientFromEnv()
	c.Login()

	st.Expect(t, c.AppToken, "1234")
	st.Expect(t, gock.IsDone(), true)
}

func TestListDevices(t *testing.T) {
	defer gock.Off()

	gock.New("https://" + domain).
		Post(apiVersion + "/user/login").
		Reply(200).
		JSON(LoginResponse{
			Username:  "foo",
			Avatar:    "bar",
			AppToken:  "1234",
			ExpiresAt: "5678",
		})

	device_response := make([]Device, 0)
	device_response = append(device_response, Device{
		Id:    "ABC",
		Title: "example@example.org",
		Model: "E-Mail",
		Image: "https://devices.pushnotifier.de/virtual/E-Mail.png",
	})

	gock.New("https://" + domain).
		Get(apiVersion + "/devices").
		Reply(200).
		JSON(device_response)

	c := NewClientFromEnv()
	c.Login()
	devices, err := c.ListDevices()

	st.Expect(t, err, nil)
	st.Expect(t, devices, device_response)
	st.Expect(t, len(devices), 1)
	st.Expect(t, gock.IsDone(), true)
}

func TestSendText(t *testing.T) {
	defer gock.Off()

	gock.New("https://" + domain).
		Post(apiVersion + "/user/login").
		Reply(200).
		JSON(LoginResponse{
			Username:  "foo",
			Avatar:    "bar",
			AppToken:  "1234",
			ExpiresAt: "5678",
		})

	gock.New("https://" + domain).
		Put(apiVersion + "/notifications/text").
		Reply(200).
		JSON(SendResponse{
			Success: []string{"foo"},
			Error:   []string{"bar"},
		})

	c := NewClientFromEnv()
	c.Login()
	resp, err := c.SendText([]string{"foo"}, "test message")

	st.Expect(t, err, nil)
	st.Expect(t, resp, SendResponse{[]string{"foo"}, []string{"bar"}})
	st.Expect(t, gock.IsDone(), true)
}
