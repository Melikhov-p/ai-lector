package responser

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func (r *Responser) WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		r.log.Error("Error encoding JSON", zap.Error(err))
	}
}

func (r *Responser) WriteError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	writeBytes, err := w.Write([]byte(err.Error()))
	if err != nil {
		r.log.Error("Error encoding JSON", zap.Error(err))
	} else {
		if writeBytes == 0 {
			r.log.Warn("error is 0")
		}
	}
}
