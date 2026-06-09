package httpapi

import (
	"encoding/json"
	"errors"
	"io"
)

var errExtraJSONValue = errors.New("request body must contain a single JSON value")

func decodeStrictJSON(body io.Reader, target any) error {
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		return errExtraJSONValue
	}
	return nil
}
