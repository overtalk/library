package file_test

import (
	"testing"

	. "web-layout/utils/file"
)

func TestDir(t *testing.T) {
	path := "xxx"
	if Exists(path) {
		if IsDir(path) {
			t.Logf("path[%s] is a dir", path)
			return
		}
		t.Logf("path[%s] is a file", path)
		return
	}
	t.Logf("path[%s] not exists", path)
}
