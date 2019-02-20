package web

import (
	"encoding/json"
	"net/http"
)

type Codec interface {
	Encode(response http.ResponseWriter, values ...interface{}) error
	Decode(request *http.Request, values ...interface{}) error
}

type JSONCodec struct{}

func (j *JSONCodec) Encode(response http.ResponseWriter, values ...interface{}) error {
	encoder := json.NewEncoder(response)

	len := len(values)
	if len > 1 {
		return encoder.Encode(values)
	} else if len == 1 {
		return encoder.Encode(values[0])
	}

	return nil
}

func (j *JSONCodec) Decode(request *http.Request, values ...interface{}) error {
	decoder := json.NewDecoder(request.Body)

	len := len(values)

	if len > 1 {
		return decoder.Decode(&values)
	} else if len == 1 {
		return decoder.Decode(values[0])
	}

	return nil
}
