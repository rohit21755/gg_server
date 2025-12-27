package main

import (
	"log"
	"net/http"
)

func internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s path:%s error: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Bad request: %s path:%s error: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found: %s path:%s error: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusNotFound, err.Error())
}

func conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Bad request: %s path:%s error: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusConflict, err.Error())
}

func unauthorizedResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Unauthorized: %s path:%s error: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusUnauthorized, "invalid credentials")
}
