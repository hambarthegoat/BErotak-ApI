package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

const errBodyTidakValid = "body harus berupa JSON object yang valid dengan tipe field yang benar"

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, pesan string) {
	writeJSON(w, status, map[string]string{"error": pesan})
}

func parseID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "id harus berupa angka")
		return 0, false
	}
	return id, true
}

func parseJSON(r *http.Request, dst any) error {
	return json.NewDecoder(r.Body).Decode(dst)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := parseJSON(r, dst); err != nil {
		writeError(w, http.StatusBadRequest, errBodyTidakValid)
		return false
	}
	return true
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotFound, "endpoint tidak ditemukan")
}
