package rest

import (
	"auth/internal/constants"
	"auth/internal/domain/models"
	"auth/internal/ports"
	"context"
	"fmt"
	"github.com/go-stack/stack"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

func Validate(h *helpers, authService ports.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			at, err := r.Cookie(constants.ACCESS_TOKEN)
			if err != nil {
				h.Error(rw, r, ErrorForbidden(err))
				return
			}

			rt, err := r.Cookie(constants.REFRESH_TOKEN)
			if err != nil {
				h.Error(rw, r, ErrorForbidden(err))
				return
			}

			user, err := authService.ValidateTokens(r.Context(), &models.TokenPair{
				AccessToken:  at.Value,
				RefreshToken: rt.Value,
			})
			if err != nil {
				h.Error(rw, r, ErrorForbidden(err))
				return
			}

			ctx := context.WithValue(r.Context(), constants.CTX_USER, user)

			next.ServeHTTP(rw, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func Logger(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			defer func() {
				logger.WithFields(
					logrus.Fields{
						"request_id":   r.Context().Value(constants.REQUEST_ID),
						"request_path": r.URL.Path,
						"status":       r.Response.Status,
						"method":       r.Method,
						"query":        r.URL.RawQuery,
						"ip":           r.RemoteAddr,
						"trace.id":     trace.SpanFromContext(r.Context()).SpanContext().TraceID().String(),
						"user-agent":   r.UserAgent(),
					}).Info("request completed")
			}()
			next.ServeHTTP(rw, r)
		}
		return http.HandlerFunc(fn)
	}
}

func RequestID(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get("X-Request-ID")
		if rid == "" {
			rid = uuid.NewV4().String()
		}
		ctx := context.WithValue(r.Context(), constants.REQUEST_ID, rid)
		rw.Header().Add("X-Request-ID", rid)

		next.ServeHTTP(rw, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func Recover(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			defer func() {
				if p := recover(); p != nil {
					err, ok := p.(error)
					if !ok {
						err = fmt.Errorf("%v", p)
					}

					var stackTrace stack.CallStack
					traces := stack.Trace().TrimRuntime()

					for i := 0; i < len(traces); i++ {
						t := traces[i]
						tFunc := t.Frame().Function

						if tFunc == "runtime.gopanic" || tFunc == "go.opentelemetry.io/otel/sdk/trace.(*span).End" {
							continue
						}

						if tFunc == "net/http.HandlerFunc.ServeHTTP" {
							break
						}
						stackTrace = append(stackTrace, t)
					}
					logger.WithFields(
						logrus.Fields{
							"trace.id":   trace.SpanFromContext(r.Context()).SpanContext().TraceID().String(),
							"request-id": r.Context().Value(constants.REQUEST_ID),
							"stack":      fmt.Sprintf("%+v", stackTrace),
						}).Panic(err)

					http.Error(rw, http.StatusText(http.StatusInternalServerError),
						http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(rw, r)
		}
		return http.HandlerFunc(fn)
	}
}

func Tracer(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		operation := r.Method + " " + r.URL.Path

		otelhttp.NewHandler(next, operation).ServeHTTP(rw, r)
	}
	return http.HandlerFunc(fn)
}
