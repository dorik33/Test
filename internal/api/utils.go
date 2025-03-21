package api

import (
	"encoding/json"
	"net/http"
)

func (api *API) sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Message{
		StatusCode: statusCode,
		Message:    message,
		IsError:    true,
	})
}
