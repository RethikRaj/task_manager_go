package common

import (
	"encoding/json"
	"net/http"

	"github.com/RethikRaj/task_manager_go/internal/errs"
)

func WriteJSONError(w http.ResponseWriter, message string, code string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errResponse := errs.APIError{
		Message: message,
		Status:  status,
		Code:    code,
	}

	json.NewEncoder(w).Encode(errResponse)
}
