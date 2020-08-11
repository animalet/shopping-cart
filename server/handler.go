package server

import "net/http"

type CartHandler interface {
	Create(writer http.ResponseWriter, request *http.Request)
	Add(writer http.ResponseWriter, request *http.Request)
	Remove(writer http.ResponseWriter, request *http.Request)
	Total(writer http.ResponseWriter, request *http.Request)
}
