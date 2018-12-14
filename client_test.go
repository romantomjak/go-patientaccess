package patientaccess

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
    "net/url"
    "reflect"
    "testing"
    "time"
)

var (
    mux *http.ServeMux
    client *Client
    server *httptest.Server
)

func setup() {
    mux = http.NewServeMux()
    server = httptest.NewServer(mux)

    client = NewClient()
    url, _ := url.Parse(server.URL)
    client.BaseURL = url
}

func teardown() {
    server.Close()
}

func assertHttpMethod(t *testing.T, got, want string) {
    t.Helper()
    if got != want {
        t.Errorf("got %+v, want %+v", got, want)
    }
}

func assertEqual(t *testing.T, got, want interface{}) {
    t.Helper()
    if !reflect.DeepEqual(got, want) {
        t.Errorf("got %+v, want %+v", got, want)
    }
}

func assertError(t *testing.T, got, want error) {
    t.Helper()
    if got != want {
        t.Errorf("got error '%s', want '%s'", got, want)
    }
}

func TestAccessTokenExpiresIn(t *testing.T) {
    tokenExpiresIn := time.Now().Add(time.Minute * 5).Format("2006-01-02T15:04:05.999999Z")
    jsonBlob := fmt.Sprintf(`{"access_token": "28d5cf150df203a0002f48395e380dff", "expires_in": "%s"}`, tokenExpiresIn)

    want := AccessToken {
        Token: "28d5cf150df203a0002f48395e380dff",
        ExpiresIn: 299,
    }

    var got AccessToken
    err := json.Unmarshal([]byte(jsonBlob), &got)
    if err != nil {
        t.Errorf("Token unmarshaling error: %v", err)
    }

    assertEqual(t, got, want)
}

func TestGetToken(t *testing.T) {
    setup()
    defer teardown()

    username := "roman"
    password := "sikr3t"

    tokenExpiresIn := time.Now().Add(time.Minute * 5).Format("2006-01-02T15:04:05.999999Z")
    jsonBlob := fmt.Sprintf(`{"accessToken": {"access_token": "28d5cf150df203a0002f48395e380dff", "expires_in": "%s"}}`, tokenExpiresIn)

    mux.HandleFunc("/authorization/signin", func(w http.ResponseWriter, r *http.Request) {
        assertHttpMethod(t, r.Method, "POST")

        want := map[string]string{
            "username": username,
            "password": password,
        }

        var got map[string]string
        json.NewDecoder(r.Body).Decode(&got)
        assertEqual(t, want, got)
        
        fmt.Fprint(w, jsonBlob)
    })

    want := &AccessToken {
        Token: "28d5cf150df203a0002f48395e380dff",
        ExpiresIn: 299,
    }

    got, _ := client.GetToken(username, password)

    assertEqual(t, got, want)
}

func TestGetTokenBadCredentials(t *testing.T) {
    setup()
    defer teardown()

    mux.HandleFunc("/authorization/signin", func(w http.ResponseWriter, r *http.Request) {
        assertHttpMethod(t, r.Method, "POST")
        fmt.Fprint(w, `{"accessToken": null}"`)
    })

    _, err := client.GetToken("roman", "sikr3t")

    assertError(t, err, ErrBadCredentials)
}

func TestGetTokenBadStatusCode(t *testing.T) {
    setup()
    defer teardown()

    _, err := client.GetToken("roman", "sikr3t")
    assertError(t, err, ErrBadStatusCode)
}
