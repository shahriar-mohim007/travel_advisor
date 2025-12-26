package helpers

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  int         `json:"-"`
	Data    interface{} `json:"data,omitempty"`
	Message interface{} `json:"message,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func (r *Response) Render(w http.ResponseWriter) error {
	bb, err := json.Marshal(r)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	if r.Status != 0 {
		w.WriteHeader(r.Status)
	}
	_, err = w.Write(bb)
	return err
}
