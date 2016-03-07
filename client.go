package coprhd

import (
	"crypto/tls"
	"errors"
	"gopkg.in/jmcvetta/napping.v3"
	"net/http"
	"net/url"
	"strings"
)

const (
	LoginUri      = "login"
	ProxyTokenUri = "proxytoken"
)

var (
	ErrLoginFailed = errors.New("Login failed")
	ErrAuthFailed  = errors.New("Authorization failed")

	// httpClient is a runtime initialized client that ignores certificate errors
	httpClient *http.Client
)

type (
	Client struct {
		// session is the rest session
		s *napping.Session

		// host is the corphd api url with port (i.e, https://host:port)
		host string

		// proxyToken is the token used for authentication
		proxyToken string
	}

	ApiError struct {
		Code        int    `xml:"code"`
		Retryable   bool   `xml:"retryable"`
		Description string `xml:"description"`
		Details     string `xml:"details"`
	}
)

func init() {
	httpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

// NewClient returns a new coprhd rest api client using the specified proxy token
func NewClient(host string, token string) *Client {
	return &Client{
		newAuthSession(token),
		host,
		token,
	}
}

// GetProxyToken is a helper function for creating a new proxy token for the user
func GetProxyToken(host string, username string, password string) (string, error) {
	e := ApiError{}

	creds := url.UserPassword(username, password)

	s := &napping.Session{
		Userinfo: creds,
		Client:   httpClient,
	}

	path := buildUrl(host, LoginUri)

	resp, err := s.Get(path, nil, nil, &e)
	if err != nil {
		return "", err
	}
	if resp.Status() != http.StatusOK {
		return "", ErrLoginFailed
	}
	if e.Code > 0 {
		return "", errors.New(e.Description)
	}

	// this is the temporary login auth token
	token := resp.HttpResponse().Header.Get("X-SDS-AUTH-TOKEN")

	path = buildUrl(host, ProxyTokenUri)

	s = newAuthSession(token)

	resp, err = s.Get(path, nil, nil, &e)
	if err != nil {
		return "", err
	}
	if resp.Status() != http.StatusOK {
		return "", ErrAuthFailed
	}
	if e.Code > 0 {
		return "", errors.New(e.Description)
	}

	// this is the proxy token that should be use with the client
	token = resp.HttpResponse().Header.Get("X-SDS-AUTH-PROXY-TOKEN")

	return token, nil
}

func (this *Client) PathForResource(r string) string {
	return buildUrl(this.host, r)
}

func (this *Client) Get(r string, p *url.Values, result interface{}) error {
	e := ApiError{}
	path := this.PathForResource(r)
	resp, err := this.s.Get(path, p, result, &e)

	if err != nil {
		return err
	}
	if resp.Status() >= http.StatusBadRequest {
		return errors.New(http.StatusText(resp.Status()))
	}
	if e.Code > 0 {
		return errors.New(e.Description)
	}

	return nil
}

func buildUrl(host string, path string) string {
	if !strings.HasSuffix(host, "/") {
		host = host + "/"
	}
	if strings.HasPrefix(path, "/") {
		path = strings.TrimPrefix(path, "/")
	}
	return host + path
}

func newAuthSession(token string) *napping.Session {
	h := http.Header{}

	h.Set("X-SDS-AUTH-TOKEN", token)

	return &napping.Session{
		Client: httpClient,
		Header: &h,
	}
}
