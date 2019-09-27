package serializer

import (
	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigDefault

// Encode data with the json serializer
func Encode(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Decode data with the json serializer
func Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
