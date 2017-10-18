package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/horae-app/api/Calendar"
	"github.com/horae-app/api/Cassandra"
	"github.com/horae-app/api/Company"
	"github.com/horae-app/api/Contact"
	"github.com/horae-app/api/Device"
	"log"
	"net/http"
	"time"
)

type healthcheckResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
}

func main() {
	CassandraSession := Cassandra.Session
	defer CassandraSession.Close()

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/healthcheck", Middlewares(healthcheck)).Methods("GET")
	router.HandleFunc("/company/new", Middlewares(Company.Post)).Methods("POST")
	router.HandleFunc("/company/auth", Middlewares(Company.Auth)).Methods("POST")

	router.HandleFunc("/{companyId}/contact/new", Middlewares(Contact.Post)).Methods("POST")
	router.HandleFunc("/{companyId}/contact/{contactId}", Middlewares(Contact.Delete)).Methods("DELETE")
	router.HandleFunc("/{companyId}/contact", Middlewares(Contact.List)).Methods("GET")

	router.HandleFunc("/{companyId}/calendar/{contactId}/new", Middlewares(Calendar.Post)).Methods("POST")
	router.HandleFunc("/{companyId}/calendar/{calendarId}", Middlewares(Calendar.Delete)).Methods("DELETE")
	router.HandleFunc("/{companyId}/calendar/{calendarId}/notify", Middlewares(Calendar.Notify)).Methods("GET")
	router.HandleFunc("/{companyId}/calendar", Middlewares(Calendar.List)).Methods("GET")
	router.HandleFunc("/{companyId}/calendar/{calendarId}/confirm", Middlewares(Calendar.Confirm)).Methods("GET")
	router.HandleFunc("/{companyId}/calendar/{calendarId}/cancel", Middlewares(Calendar.Cancel)).Methods("GET")

	router.HandleFunc("/contact/auth", Middlewares(Contact.UserAuth)).Methods("POST")
	router.HandleFunc("/contact/calendar", Middlewares(Contact.Calendar)).Methods("POST")

	router.HandleFunc("/device/register", Middlewares(Device.Post)).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(healthcheckResponse{Status: "OK", Code: 200})
}

func Middlewares(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ResponseHeaderJsonMiddleware(w)

		if !RequestPostBodyEmptyValidationMiddleware(w, r) {
			return
		}

		LoggerMiddleware(h, w, r)
	}
}

func LoggerMiddleware(h http.HandlerFunc, w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	h.ServeHTTP(w, r)

	log.Printf(
		"%s\t%s\t%s",
		r.Method,
		r.RequestURI,
		time.Since(start),
	)
}

func ResponseHeaderJsonMiddleware(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func RequestPostBodyEmptyValidationMiddleware(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		if r.ContentLength == 0 {
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Please send a request body"})
			return false
		}
	}

	return true
}

type ErrorResponse struct {
	Error string
}
