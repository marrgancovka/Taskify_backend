package responser

import (
	"encoding/json"
	"net/http"
)

type MessageResponse struct {
	Msg string `json:"msg"`
}

func Send200(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	resp, err := json.Marshal(v)
	if err != nil {
		return
	}
	_, _ = w.Write(resp)
}

func Send201(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	resp, err := json.Marshal(v)
	if err != nil {
		return
	}
	_, _ = w.Write(resp)
}

func Send500(w http.ResponseWriter, msg string) {
	resp, err := json.Marshal(MessageResponse{msg})
	if err != nil {
		return
	}
	w.WriteHeader(500)
	_, _ = w.Write(resp)
}

func Send400(w http.ResponseWriter, msg string) {
	resp, err := json.Marshal(MessageResponse{msg})
	if err != nil {
		return
	}
	w.WriteHeader(400)
	_, _ = w.Write(resp)
}

func Send404(w http.ResponseWriter, msg string) {
	resp, err := json.Marshal(MessageResponse{msg})
	if err != nil {
		return
	}
	w.WriteHeader(404)
	_, _ = w.Write(resp)
}
