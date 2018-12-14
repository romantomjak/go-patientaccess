package patientaccess

import (
    "bytes"
    "encoding/json"
    "errors"
    "net/http"
    "net/url"
    "strings"
    "time"
)

const (
    defaultBaseURL = "https://api.patientaccess.com/api"
    userAgent = "go-patientaccess/0.1 (+https://github.com/romantomjak/go-patientaccess)"
)

var (
    ErrBadCredentials = errors.New("bad credentials")
    ErrBadStatusCode = errors.New("bad HTTP status code")
)

// Client manages communication with Patient Access API
type Client struct {
    // HTTP client used to communicate with PA API
    client *http.Client

    // Base URL for API requests.
    BaseURL *url.URL

    // User agent for client
    UserAgent string
}

// AccessToken contains the Token used to authenticate with PA API
// and number of seconds until it expires.
type AccessToken struct {
    // The actual token
    Token string `json:"access_token"`

    // Custom field that converts a timestamp into seconds
    ExpiresIn ExpiresIn `json:"expires_in"`
}

type ExpiresIn int64

func (e *ExpiresIn) UnmarshalJSON(b []byte) error {
   var ts time.Time
   if err := json.Unmarshal(b, &ts); err != nil {
       return err
   }
   *e = ExpiresIn(time.Until(ts).Seconds())
   return nil
}

// AuthResponse is a minimal structure for holding authentication response data
type AuthResponse struct {
    AccessToken AccessToken `json:"accessToken"`
}

// Returns a new Patient Access API client
func NewClient() *Client {
    baseURL, _ := url.Parse(defaultBaseURL)
    return &Client{client: http.DefaultClient, BaseURL: baseURL, UserAgent: userAgent}
}

// Creates a new HTTP request with all necessary HTTP headers
func (c *Client) NewRequest(method, apiPath string, body interface{}) (*http.Request, error) {
    loc, err := joinPaths(c.BaseURL, apiPath)
    if err != nil {
        return nil, err
    }

    buf := new(bytes.Buffer)
    if body != nil {
        err = json.NewEncoder(buf).Encode(body)
        if err != nil {
            return nil, err
        }
    }

    req, err := http.NewRequest(method, loc.String(), buf)
    if err != nil {
        return nil, err
    }

    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("Accept", "application/json, text/plain, */*")
    req.Header.Add("User-Agent", c.UserAgent)
    return req, nil
}

// Obtains a new PA API access token
func (c *Client) GetToken(username, password string) (token *AccessToken, err error) {
    form := map[string]string{
        "username": username,
        "password": password,
    }

    req, err := c.NewRequest("POST", "/authorization/signin", form)
    if err != nil {
        return nil, err
    }

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode < 200 || resp.StatusCode > 299 {
        return nil, ErrBadStatusCode
    }

    var authResp AuthResponse
    err = json.NewDecoder(resp.Body).Decode(&authResp)
    if err != nil {
        return nil, err
    }

    if (authResp.AccessToken.Token == "") {
        return nil, ErrBadCredentials
    }
    
    return &authResp.AccessToken, nil
}

func joinPaths(baseURL *url.URL, path string) (*url.URL, error) {
    u, err := url.Parse(path)
    if err != nil {
        return nil, err
    }

    if !strings.HasSuffix(baseURL.Path, "/") {
        baseURL.Path = baseURL.Path + "/"
    }

    if strings.HasPrefix(u.Path, "/") {
        u.Path = u.Path[1:]
    }

    return baseURL.ResolveReference(u), nil
}
