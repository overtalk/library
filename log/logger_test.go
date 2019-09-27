package log_test

import (
	"fmt"
	"testing"

	. "web-layout/utils/log"
)

func TestLogger(t *testing.T) {
	err := fmt.Errorf("no password")
	Logger.Error("test",
		ErrorField(err),
		Field("key1", "value1"),
	)
}
