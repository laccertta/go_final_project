package api

import (
	"encoding/json"
	"net/http"
)

func writeJSONStatus(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode JSON: "+err.Error(), http.StatusInternalServerError)
	}
}
