package responsewriter

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		zap.L().Error("failed to encode successful response", zap.Error(err))
	}
}
