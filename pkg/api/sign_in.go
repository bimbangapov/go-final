package api

import (
	"bytes"
	"encoding/json"
	"go-final/pkg/config"
	"net/http"
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
		w.Write(res)
		return
	}
	err = json.Unmarshal(buff.Bytes(), &signInRequest)
	if h.password != signInRequest.Password {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		res, err := json.Marshal(ErrorResponse{Error: "Неверный пароль"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(res)
		return
	}
	token, err := config.GetToken(h.password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, false, SignInResponse{Token: token})
}
