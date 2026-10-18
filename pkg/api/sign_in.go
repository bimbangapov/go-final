package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"go-final/pkg/config"
)

type SignInRequest struct {
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var (
		signInRequest SignInRequest
		buff          bytes.Buffer
	)

	_, err := buff.ReadFrom(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		res, err := json.Marshal(ErrorResponse{Error: "Неверный пароль"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write(res); err != nil {
			h.writeJSON(w, 400, ErrorResponse{Error: err.Error()})
			return
		}
		return
	}
	if err = json.Unmarshal(buff.Bytes(), &signInRequest); err != nil {
		h.writeJSON(w, 400, ErrorResponse{Error: err.Error()})
		return
	}

	if h.password != signInRequest.Password {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		res, err := json.Marshal(ErrorResponse{Error: "Неверный пароль"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write(res); err != nil {
			h.writeJSON(w, 400, ErrorResponse{Error: err.Error()})
			return
		}
		return
	}
	token, err := config.GetToken(h.password)
	if err != nil {
		h.writeJSON(w, 400, ErrorResponse{Error: err.Error()})
		return
	}

	h.writeJSON(w, 200, SignInResponse{Token: token})
}
