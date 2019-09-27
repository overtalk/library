package gate

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

var (
	// errInvalidHandlerType defines the error that the type of handler is not func
	errInvalidHandlerType = errors.New("handler should be function type")
	// errInvalidInputParamsNum defines the error that number of input args is invalid
	errInvalidInputParamsNum = errors.New("handler function require 1 or 2 input parameters")
	// errInvalidOutputParamsNum defines the error that number of output args is invalid
	errInvalidOutputParamsNum = errors.New("handler function require 1 output parameters")
)

// wrap wrap general handler to gin handler
func wrap(f interface{}) gin.HandlerFunc {
	t := reflect.TypeOf(f)

	// check the Type of handler
	if t.Kind() != reflect.Func {
		panic(errInvalidHandlerType)
	}

	// check the number of input args
	numIn := t.NumIn()
	if numIn < 1 || numIn > 2 {
		panic(errInvalidInputParamsNum)
	}

	// check the number of output args
	numOut := t.NumOut()
	if numOut != 1 {
		panic(errInvalidOutputParamsNum)
	}

	return func(c *gin.Context) {
		// handler input
		// first arg should be context.Context
		inValues := []reflect.Value{
			reflect.ValueOf(c),
		}

		// if necessary
		if numIn == 2 {
			req := newReqInstance(t.In(1))
			if err := c.Bind(req); err != nil {
				fmt.Println(err)
				// TODO: err handle
			}

			inValues = append(inValues, reflect.ValueOf(req))
		}

		ret := reflect.ValueOf(f).Call(inValues)
		c.JSON(http.StatusOK, ret[0].Interface())
	}
}

func newReqInstance(t reflect.Type) interface{} {
	switch t.Kind() {
	case reflect.Ptr, reflect.Interface:
		return newReqInstance(t.Elem())
	default:
		return reflect.New(t).Interface()
	}
}
