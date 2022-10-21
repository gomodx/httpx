package httpx

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/gorilla/securecookie"
)

const (
	XForwardedFor = "X-Forwarded-Proto"
)

func BodyReadAll(r *http.Request) ([]byte, error) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	// replace body so it can be read again later
	r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	return buf, nil
}

func ProxyBaseURL(r *http.Request) url.URL {
	port := os.Getenv("PROXY_PORT")
	scheme := r.Header.Get(XForwardedFor)
	host := r.Host

	if port != "" {
		host = net.JoinHostPort(host, port)
	}

	if scheme == "" {
		scheme = "http"
	}

	return url.URL{
		Host:   host,
		Scheme: scheme,
		Path:   os.Getenv("PROXY_PREFIX"),
	}
}

func ProxyURLFull(r *http.Request) url.URL {
	u := ProxyBaseURL(r)
	u.Path = path.Join(u.Path, r.URL.Path)
	u.RawQuery = r.URL.RawQuery
	return u
}

func IsSecure(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get(XForwardedFor), "https")
}

type SecureCookieParams struct {
	Cookie        *http.Cookie
	SigningSecret []byte
	EncryptSecret []byte
	Data          any
}

func EncryptCookie(opts SecureCookieParams) (*http.Cookie, error) {
	s := securecookie.New(opts.SigningSecret, opts.EncryptSecret)
	data, err := s.Encode(opts.Cookie.Name, opts.Data)
	if err != nil {
		return nil, err
	}

	opts.Cookie.Value = data
	return opts.Cookie, nil
}

func DecryptCookie[T any](opts SecureCookieParams) (*T, error) {
	v := new(T)
	s := securecookie.New(opts.SigningSecret, opts.EncryptSecret)
	if err := s.Decode(opts.Cookie.Name, opts.Cookie.Value, v); err != nil {
		return nil, err
	}
	return v, nil
}

func DefaultCookieSigningSecret() []byte {
	data, _ := hex.DecodeString(os.Getenv("COOKIE_SIG_SECRET"))
	return data
}

func DefaultCookieEncSecret() []byte {
	data, _ := hex.DecodeString(os.Getenv("COOKIE_SIG_SECRET"))
	return data
}
