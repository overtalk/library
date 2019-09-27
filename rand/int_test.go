package random_test

import (
	"testing"

	. "web-layout/utils/rand"
)

func TestInt(t *testing.T) {
	t.Log(Int())
	t.Log(Int(2))
	t.Log(Int(1, 2))
}

func BenchmarkInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int()
	}
}
