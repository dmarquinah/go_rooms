package types

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Status  string `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
	Data    any    `json:"data,omitempty" bson:"data"`
}

type ErrorResponse struct {
	Status  string `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
	Code    int    `json:"cpde" bson:"code"`
}

func WriteSuccessResponse(w http.ResponseWriter, message string, data any) {
	response := SuccessResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func WriteErrorResponse(w http.ResponseWriter, message string, code int) {
	response := ErrorResponse{
		Status:  "error",
		Message: message,
		Code:    code,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func GetSuccessMessage(r *http.Request) string {
	switch method := r.Method; method {
	case "GET":
		return "Success retrieving data"

	case "POST":
		return "Success creating data"

	case "PUT":
		return "Success modifying data"

	case "DELETE":
		return "Success deleting data"

	default:
		return "Success"
	}
}
