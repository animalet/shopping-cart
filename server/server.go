package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"shopping-cart/server/storage"
)

func StartServer() {
	cartService := CartService{Storage: getStorage()}
	restHandler := Handler{CartService: &cartService}
	router := mux.NewRouter()
	router.Handle("/cart", dontPanic(http.HandlerFunc(restHandler.Create))).Methods("POST")
	router.Handle("/cart/{cartId}", dontPanic(http.HandlerFunc(restHandler.Remove))).Methods("DELETE")
	router.Handle("/cart/{cartId}/{code}", dontPanic(http.HandlerFunc(restHandler.Add))).Methods("PUT")
	router.Handle("/cart/{cartId}", dontPanic(http.HandlerFunc(restHandler.Total))).Methods("GET")

	log.Printf("Listening on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", router))
}

// If it is running inside our docker-compose.yml it will connect to Consul and use Redis as storage (for scalability!)
func getStorage() storage.Storage {
	var st storage.Storage
	if consulAddr, found := os.LookupEnv("CONSUL_ADDR"); found {
		log.Printf("Initializing Redis storage in %s", consulAddr)
		st = storage.NewRedisStorage(consulAddr)
	} else {
		log.Print("Environment variable CONSUL_ADDR not found. Defaulting to local map storage.")
		st = storage.NewMapStorage()
	}
	return st
}

func dontPanic(h http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				message := fmt.Sprintf("Panic handling request: %s\n%s", r, string(debug.Stack()))
				log.Print(message)
				http.Error(writer, message, http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(writer, r)
	})
}
