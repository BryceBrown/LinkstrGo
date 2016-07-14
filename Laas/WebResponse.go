package Laas

import (
	"net/http"
	"encoding/json"
)


type WebResponse struct {}

func (resp *WebResponse) SendErrorMessage(code int, message string, w http.ResponseWriter){
	msg := ErrorMessage{message}
	json, _ := json.Marshal(msg)
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(json)
}

func (resp *WebResponse) RespondWithJson(json []byte,  w http.ResponseWriter){
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Write(json)
}

func (resp *WebResponse) SendAuthResponse(message string, token string, code int,  w http.ResponseWriter){
	response := AuthResponse{Message: message, Token: token}
	respJson, _ := json.Marshal(response)
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(respJson)
}

func (resp *WebResponse) SendUrlResponse(message string, url string, code int, w http.ResponseWriter){
	response := LinkResponse{Message: message, Url: url}
	respJson, _ := json.Marshal(response)
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(respJson)
}