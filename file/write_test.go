package file_test

import (
	"testing"

	. "github.com/qinhan-shu/go-utils/file"
)

func TestWrite(t *testing.T) {
	path := "xxx"
	writeBytes := []byte("qinhan")
	if err := Write(path, writeBytes); err != nil {
		t.Error(err)
		return
	}
}
