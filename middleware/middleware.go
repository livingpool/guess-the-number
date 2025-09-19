package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/livingpool/utils"
)

type LoggingConfig struct {
	DefaultLevel     slog.Level
	ServerErrorLevel slog.Level
	ClientErrorLevel slog.Level
}

type RequestId string

const RequestIdKey RequestId = "reqId"

type Middleware func(http.Handler) http.Handler

func CreateStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			next = xs[i](next)
		}

		return next
	}
}

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func Logging(logger *slog.Logger, config LoggingConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			reqId := uuid.New().String()

			wrapped := &wrappedWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			r = r.WithContext(context.WithValue(r.Context(), RequestIdKey, reqId))

			next.ServeHTTP(wrapped, r)

			level := config.DefaultLevel
			if wrapped.statusCode >= http.StatusInternalServerError {
				level = config.ServerErrorLevel
			} else if wrapped.statusCode >= http.StatusBadRequest {
				level = config.ClientErrorLevel
			}

			defer func() {
				logger.LogAttrs(
					r.Context(),
					level,
					strconv.Itoa(wrapped.statusCode),
					slog.String(string(RequestIdKey), reqId),
					slog.String("ip", utils.ReadUserIP(r).String()),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Duration("timeSpent", time.Since(start)),
				)
			}()
		})
	}
}
