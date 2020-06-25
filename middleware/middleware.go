package middleware

import (
	"github.com/go-wyvern/leego"

	"reflect"
	"runtime"
)

type (
	// Skipper func return true or false
	Skipper func(c leego.Context) bool
)

// MiddlewareConfig config for route
var MiddlewareConfig = make(map[leego.Route]interface{})

// defaultSkipper
func defaultSkipper(c leego.Context) bool {
	return false
}

// defaultFormatLeeError
func defaultFormatLeeError(err error, middlewareName string) leego.LeeError {
	return err
}

// Middleware interface for mid
type Middleware interface {
	Skipper(c leego.Context) bool
	FormatLeeError(err error, middlewareName string) leego.LeeError
}

// DefaultMiddleware for default
type DefaultMiddleware struct{}

// Skipper func default mid method
func (d DefaultMiddleware) Skipper(c leego.Context) bool {
	return defaultSkipper(c)
}

// FormatLeeError  func default mid method
func (d DefaultMiddleware) FormatLeeError(err error, middlewareName string) leego.LeeError {
	return defaultFormatLeeError(err, middlewareName)
}

func handlerName(h leego.HandlerFunc) string {
	t := reflect.ValueOf(h).Type()
	if t.Kind() == reflect.Func {
		return runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	}
	return t.String()
}
