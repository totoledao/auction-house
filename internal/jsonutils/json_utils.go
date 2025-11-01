package jsonutils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/totoledao/auction-house/internal/validator"
)

func EncodeJson[T any](w http.ResponseWriter, r *http.Request, statusCode int, data T) error {
	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("\nERROR EncodeJson failed: %v\n", err)
		return err
	}

	return err
}

func DecodeValidJson[T validator.Validator](r *http.Request) (T, map[string]string, error) {
	var data T

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Printf("\nERROR DecodeValidJson failed: %v\n", err)
		return data, nil, err
	}

	problems := data.Valid(r.Context())
	if len(problems) > 0 {
		return data, problems, fmt.Errorf("invalid %T: %d problems", data, len(problems))
	}

	return data, nil, nil
}

func DecodeJson[T any](r *http.Request) (T, error) {
	var data T

	err := json.NewDecoder(r.Body).Decode(&data)
	log.Printf("\nERROR DecodeJson failed: %v\n", err)
	if err != nil {
		return data, err
	}

	return data, nil
}
