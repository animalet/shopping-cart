package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

// Rest handler for cart operations. Used in server
type Handler struct {
	CartService *CartService
}

func (r *Handler) Create(writer http.ResponseWriter, _ *http.Request) {
	cartId := r.CartService.Create()
	writer.WriteHeader(http.StatusCreated)
	_, _ = writer.Write([]byte(cartId))
}

func (r *Handler) Add(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	code, cartId := vars["code"], vars["cartId"]
	found := r.CartService.Add(cartId, code)
	if found {
		writer.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(writer, fmt.Sprintf("Cart (%s) or product (%s) not found", cartId, code), http.StatusNotFound)
	}
}

func (r *Handler) Remove(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	cartId := vars["cartId"]
	found := r.CartService.Remove(cartId)
	if found {
		writer.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(writer, fmt.Sprintf("Cart %s not found", cartId), http.StatusNotFound)
	}
}

func (r *Handler) Total(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	cartId := vars["cartId"]
	total, found := r.CartService.Total(cartId)
	if found {
		_, _ = writer.Write([]byte(total.String()))
	} else {
		http.Error(writer, fmt.Sprintf("Cart %s not found", cartId), http.StatusNotFound)
	}
}
