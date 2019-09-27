package parse

// Bytes turn ( interface{} ) to ( []byte )
func Bytes(in interface{}) []byte {
	var ret []byte

	switch in.(type) {
	case []byte:
		ret = in.([]byte)
	case string:
		ret = []byte(in.(string))
	default:
		ret = nil
	}

	return ret
}
