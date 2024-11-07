package internal

import "net/http"

type Step func(http.HandlerFunc) http.HandlerFunc

// Prepend the middlewares to the handler in the order they are provided,
// and return the resulting (chained) handler.
func Chain(h http.HandlerFunc, m ...Step) http.HandlerFunc {
	if len(m) == 0 {
		return h
	}
	handler := h
	for i := len(m) - 1; i >= 0; i-- {
		handler = m[i](handler)
	}
	return handler
}

func NoCache(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, max-age=0")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
		}
		next(w, r)
	}
}
