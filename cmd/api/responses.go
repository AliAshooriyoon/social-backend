package main

import (
	"encoding/json"
	"net/http"
)

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	r.Body = http.MaxBytesReader(w, r.Body, int64(1_048_578))
	decode := json.NewDecoder(r.Body)
	decode.DisallowUnknownFields()
	return decode.Decode(data)
}
