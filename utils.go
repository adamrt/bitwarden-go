package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type NullString struct {
	sql.NullString
}

// NullString MarshalJSON interface redefinition
func (r NullString) MarshalJSON() ([]byte, error) {
	if r.Valid {
		return json.Marshal(r.String)
	} else {
		return json.Marshal(nil)
	}
}

func JSONResponse(w http.ResponseWriter, d interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(d)

	if err != nil {
		log.Printf("JSONResponse: %s", err)
		InternalError(w)
	}
}

type JSONErrResponse struct {
	ValidationErrors []string
	Object           string
}

func JSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	messages := []string{message}
	// ignoring the possible error because this really cant fail
	json.NewEncoder(w).Encode(JSONErrResponse{messages, "error"})
}

func InternalError(w http.ResponseWriter) {
	JSONError(w, http.StatusInternalServerError, "Internal Server Error")
}
