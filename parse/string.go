package parse

import (
	"strconv"
)

// String turn ( interface{} ) to ( string )
func String(in interface{}) string {
	var ret string
	switch in.(type) {
	case string:
		ret = in.(string)
	case []uint8:
		ret = string(in.([]uint8))
	case int64:
		ret = strconv.FormatInt(in.(int64), 10)
	case int:
		ret = strconv.Itoa(in.(int))
	default:
		ret = ""
	}

	return ret
}
