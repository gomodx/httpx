package httpx

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
)

type ContentType string

const (
	TextHTML        ContentType = "text/html; charset=utf-8"
	ApplicationJSON ContentType = "application/json; charset=utf-8"
	ApplicationXML  ContentType = "application/xml; charset=utf-8"
)

var (
	xmlDeclaration = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"
)

type ResponseWriter interface {
	Status(code int) ResponseWriter
	ContentType(media ContentType) ResponseWriter
	Send(a ...any) error

	JSON(v any) error
	XML(v any) error

	JSONErr(err Error) error
	XMLErr(err Error) error

	Header() http.Header
	Write(bytes []byte) (int, error)
	WriteHeader(statusCode int)

	Redirect(r *http.Request, url string, code int)
	SetCookie(cookie http.Cookie) ResponseWriter
	DelCookie(cookie http.Cookie) ResponseWriter
	DelCookieByName(name string) ResponseWriter
}

type responseWriter struct {
	w           http.ResponseWriter
	code        int
	contentType ContentType
	sent        bool
}

func (r *responseWriter) DelCookie(cookie http.Cookie) ResponseWriter {
	cookie.Value = ""
	cookie.Path = "/"
	cookie.Expires = time.Unix(0, 0)
	cookie.HttpOnly = true
	return r.SetCookie(cookie)
}

func (r *responseWriter) DelCookieByName(name string) ResponseWriter {
	return r.DelCookie(http.Cookie{
		Name: name,
	})
}

func (r *responseWriter) SetCookie(cookie http.Cookie) ResponseWriter {
	http.SetCookie(r.w, &cookie)
	return r
}

func (r *responseWriter) Redirect(req *http.Request, url string, code int) {
	http.Redirect(r.w, req, url, code)
}

func NewResponseWriter(w http.ResponseWriter) ResponseWriter {
	return &responseWriter{w: w}
}

func (r *responseWriter) Status(code int) ResponseWriter {
	r.code = code
	return r
}

func (r *responseWriter) ContentType(media ContentType) ResponseWriter {
	if r.contentType == "" {
		r.contentType = media
		r.Header().Set("Content-Type", string(media))
		r.Header().Set("X-Content-Type-Options", "nosniff")
	}
	return r
}

func (r *responseWriter) Send(a ...any) error {
	r.sent = true
	_, err := fmt.Fprintln(r, a...)
	return err
}

func (r *responseWriter) JSON(v any) error {
	r.ContentType(ApplicationJSON)

	if err := json.NewEncoder(r).Encode(v); err != nil {
		return err
	}

	return nil
}

func (r *responseWriter) XML(v any) error {
	r.ContentType(ApplicationXML)

	if _, err := r.Write([]byte(xmlDeclaration)); err != nil {
		return err
	}

	if err := xml.NewEncoder(r).Encode(v); err != nil {
		return err
	}

	return nil
}

func (r *responseWriter) JSONErr(err Error) error {
	return r.Status(err.Status).JSON(err)
}

func (r *responseWriter) XMLErr(err Error) error {
	return r.Status(err.Status).XML(err)
}

func (r *responseWriter) Header() http.Header        { return r.w.Header() }
func (r *responseWriter) WriteHeader(statusCode int) { r.w.WriteHeader(statusCode) }
func (r *responseWriter) Write(bytes []byte) (int, error) {
	r.sent = true
	if r.code != 0 {
		r.WriteHeader(r.code)
	}
	return r.w.Write(bytes)
}
