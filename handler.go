package httpx

import (
	"net/http"

	"go.uber.org/zap"
)

var logger, _ = zap.NewProduction()

type ResMap map[string]interface{}
type HandlerFunc func(w ResponseWriter, r *http.Request) error

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := NewResponseWriter(w)
	if err := h(res, r); err != nil {
		var logErr error
		switch ev := err.(type) {
		case Error:
			logErr = ev.Err
			_ = res.JSONErr(ev)
		case *Error:
			logErr = ev.Err
			_ = res.JSONErr(*ev)
		default:
			logErr = err
			_ = res.JSONErr(ErrInternalServer.WithErr(err))
		}

		if logErr != nil {
			// FIXME: inject a logger instead of using a global one
			logger.Error("request error", zap.Error(err))
		}
	}
}

func NativeHandlerFunc(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}
