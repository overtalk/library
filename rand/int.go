package random

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// len(args)
// 0:随机
// 1->[0,args[0])
// 2-> [args[0],args[1])
func Int(args ...int) int {
	//rand.Seed(time.Now().UnixNano())

	switch len(args) {
	case 1:
		return rand.Intn(args[0])
	case 2:
		return rand.Intn(args[1]) + args[0]
	}
	return rand.Int()
}

func Int64(args ...int64) int64 {
	//rand.Seed(time.Now().UnixNano())

	switch len(args) {
	case 1:
		return rand.Int63n(args[0])
	case 2:
		return rand.Int63n(args[1]) + args[0]
	}
	return rand.Int63()
}
