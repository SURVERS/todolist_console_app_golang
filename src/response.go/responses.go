package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type JSONResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func WriteJSONResponse(w http.ResponseWriter, success bool, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	response := JSONResponse{
		Success: success,
		Message: message,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"success":false,"message":"Ошибка формирования JSON"}`)
		return
	}

	w.WriteHeader(statusCode)
	fmt.Fprint(w, string(jsonData))
}
