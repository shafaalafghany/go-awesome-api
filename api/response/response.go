package response

import (
	"encoding/json"
	"net/http"

	apierror "awesome-api/api/error"
)

func GenerateResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(data)
}

func Error(w http.ResponseWriter, e apierror.Error) {
	GenerateResponse(w, e.HttpStatus, e)
}

func ValidationError(w http.ResponseWriter, err apierror.UnprocessableEntity) {
	GenerateResponse(w, err.HttpStatus, err)
}
