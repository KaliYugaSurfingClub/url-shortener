package mw

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"time"
)

func NewLogger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log.Info("logger middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			attr := slog.String("request_id", middleware.GetReqID(r.Context()))
			ctx := InjectLog(r.Context(), log.With(attr))

			t1 := time.Now()

			defer func() {
				log.Info(
					"request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(t1).String()),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.String("remote_addr", r.RemoteAddr),
					slog.String("user_agent", r.UserAgent()),
				)
			}()

			next.ServeHTTP(ww, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

type logKey struct{}

func ExtractLog(ctx context.Context, operation string) *slog.Logger {
	if ctx == nil {
		return nil
	}

	log, ok := ctx.Value(logKey{}).(*slog.Logger)
	if !ok {
		return nil
	}

	return log.With("op", operation)
}

func InjectLog(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, logKey{}, log)
}

func ErrAttr(err error) slog.Attr {
	return slog.String("error", err.Error())
}
