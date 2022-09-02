package rest

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

//func (h *handlers.Handler) setTokensInCookie(w http.ResponseWriter, r *http.Request, accessToken, refreshToken string) {
//	accessCookie := http.Cookie{
//		HttpOnly: true,
//		Name:     "access",
//		Value:    accessToken,
//	}
//	refreshCookie := http.Cookie{
//		HttpOnly: true,
//		Name:     "refresh",
//		Value:    refreshToken,
//	}
//
//	http.SetCookie(w, &accessCookie)
//	http.SetCookie(w, &refreshCookie)
//	return
//}

type helpers struct {
	logger *logrus.Logger
}

func NewHelpers(logger *logrus.Logger) *helpers {
	return &helpers{
		logger: logger,
	}
}

func (h *helpers) ResponseOk(w http.ResponseWriter, r *http.Request, v interface{}) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		h.Error(w, r, err)
	}
	//h.logger.Println(v)
}

func (h *helpers) Error(w http.ResponseWriter, r *http.Request, err error) {
	span := trace.SpanFromContext(r.Context())
	span.RecordError(err)

	h.logger.Error(err)
	e, ok := err.(HttpError)
	if !ok {
		http.Error(w, fmt.Sprintf("{\"error\": \"%s\"}", err), http.StatusInternalServerError)
	}
	http.Error(w, fmt.Sprintf("{\"error\": \"%s\"}", err), e.StatusCode)
	return

}
