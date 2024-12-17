package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type EnvelopeJson struct {
	Data any `json:"data"`
}

type ErrorEnvelopeJson struct {
	Error struct {
		Method    string `json:"method"`
		Path      string `json:"path"`
		Timestamp int64  `json:"timestamp"`
		Message   string `json:"message"`
	} `json:"error"`
}

func getTagValidationError(tag string) string {
	switch tag {
	case "required":
		return "required field is empty"
	case "max":
		return "text is exceeding length"
	default:
		return ""
	}
}

func ValidateStruct(payload any) error {
	var sb strings.Builder
	var err error

	if err = validate.Struct(payload); err == nil {
		return nil
	}

	err_list := err.(validator.ValidationErrors)
	if len(err_list) == 0 {
		return nil
	}

	err = err_list[0]
	sb.WriteString(err.(validator.FieldError).Field() + ": " + getTagValidationError(err.(validator.FieldError).Tag()) + ";")
	for i := 1; i < len(err_list); i++ {
		err = err_list[i]
		sb.WriteString(" " + err.(validator.FieldError).Field() + ": " + getTagValidationError(err.(validator.FieldError).Tag()) + ";")
	}

	return errors.New(sb.String())
}

func ReadJson(w http.ResponseWriter, r *http.Request, payload any) error {
	maxBytes := 1_048_578
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(payload)
}

func WriteJson(w http.ResponseWriter, status int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	return encoder.Encode(payload)
}

func WriteJsonData(w http.ResponseWriter, status int, payload any) error {
	response := EnvelopeJson{
		Data: payload,
	}

	return WriteJson(w, status, response)
}

func WriteJsonError(w http.ResponseWriter, r *http.Request, status int, message string) error {
	response := ErrorEnvelopeJson{
		Error: struct {
			Method    string `json:"method"`
			Path      string `json:"path"`
			Timestamp int64  `json:"timestamp"`
			Message   string `json:"message"`
		}{
			Method:    r.Method,
			Path:      r.URL.Path,
			Timestamp: time.Now().Unix(),
			Message:   message,
		},
	}

	return WriteJson(w, status, response)
}
