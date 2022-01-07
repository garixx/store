package main

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
)

const (
	path = "/api/v1/key/{key}"
	port = ":" + "8080"
)

func main() {
	err := initializeTransactionLog()
	if err != nil {
		os.Exit(1)
	}

	r := mux.NewRouter()
	r.HandleFunc(path, keyValuePutHandler).Methods("PUT")
	r.HandleFunc(path, keyValueGetHandler).Methods("GET")
	r.HandleFunc(path, keyValueDeleteHandler).Methods("DELETE")
	r.Use(loggingMiddleware)

	log.Fatal(http.ListenAndServe(port, r))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Infof("request received: %v %v", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
