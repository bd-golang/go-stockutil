package httputil

import (
	"bytes"
	"encoding/json"
	"io"
)

type EncoderFunc func(interface{}) (io.Reader, error)
type DecoderFunc func(io.Reader, interface{}) error

func JSONEncoder(in interface{}) (io.Reader, error) {
	if data, err := json.Marshal(in); err == nil {
		return bytes.NewBuffer(data), nil
	} else {
		return nil, err
	}
}

func JSONDecoder(in io.Reader, out interface{}) error {
	return json.NewDecoder(in).Decode(out)
}
