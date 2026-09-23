package api

import (
	"log"
	"net/http"
)

type Handler struct {
	Logger     *log.Logger
	FileServer http.Handler
}

func (h *Handler) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Sec-Fetch-Dest") == "document" {
			h.Logger.Printf("Входящий запрос: метод=%s URL=%s", r.Method, r.RequestURI)
			log.Printf("Входящий запрос: метод=%s URL=%s\n", r.Method, r.RequestURI)
		}
		next.ServeHTTP(w, r)
	})
}
