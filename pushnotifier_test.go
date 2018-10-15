package pushnotifier

import (
	"os"
	"strconv"
	"testing"

	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v1"
)

var c *Client
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

	r := SetupHttpClient(token, pkg, debug)

	c = &Client{*r, "https://" + domain + apiVersion, userAgent, username, password, ""}

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

	c.Login()

	st.Expect(t, c.AppToken, "1234")
	st.Expect(t, gock.IsDone(), true)
}

// func TestClient(t *testing.T) {
// 	defer gock.Off() // Flush pending mocks after test execution
//
// 	gock.New("http://server.com").
// 		Get("/bar").
// 		Reply(200).
// 		JSON(map[string]string{"foo": "bar"})
//
// 	st.Expect(t, nil, '1')
// }
