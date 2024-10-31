package mw

import "net/http"

// CheckAuth use mw.Logger
func CheckAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log := ExtractLog(r.Context(), "transport.rest.mw.InjectUserIdToCtx")
		if log == nil {
			panic("log is nil")
		}

		_, ok := ExtractUserID(r.Context())

		if !ok {
			//rest.Error(w, log, errs.E(errs.Unauthorized))
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
