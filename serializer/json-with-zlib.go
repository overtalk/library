package serializer

func EncodeWithZlib(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return Compress(data)
}

func DecodeWithZlib(data []byte, v interface{}) error {
	decodedData, err := Decompress(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(decodedData, v)
}
