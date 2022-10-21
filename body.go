package httpx

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

var (
	ErrNoHTTPBody = errors.New("no http body bytes to read")
)

func HasBody(r *http.Request) bool {
	return r.Body != http.NoBody
}

func BodyAs[T any](r *http.Request) (T, error) {
	v := new(T)

	if r.Body == http.NoBody {
		return *v, ErrNoHTTPBody
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return *v, err
	}

	// replace body so it can be read again later
	r.Body = ioutil.NopCloser(bytes.NewBuffer(data))

	err = json.Unmarshal(data, v)
	return *v, err
}
