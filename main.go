package main

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

const (
	path = "/api/v1/key/{key}"
	port = ":" + "8080"
)

func main() {
	err := initializeTransactionLog()
	if err != nil {
		log.Fatal("init failed: ", err)
	}

	r := mux.NewRouter()
	r.HandleFunc(path, keyValuePutHandler).Methods("PUT")
	r.HandleFunc(path, keyValueGetHandler).Methods("GET")
	r.HandleFunc(path, keyValueDeleteHandler).Methods("DELETE")
	r.Use(loggingMiddleware)

	logrus.Info("started ...")
	log.Fatal(http.ListenAndServe(port, r))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Infof("request received: %v %v", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
