package middleware

import (
	"github.com/go-wyvern/leego"

	"reflect"
	"runtime"
)

type (
	Skipper func(c leego.Context) bool
)

var MiddlewareConfig = make(map[leego.Route]interface{})

func defaultSkipper(c leego.Context) bool {
	return false
}

func defaultFormatLeegoError(err error, middlewareName string) leego.LeegoError {
	return err
}

type Middleware interface {
	Skipper(c leego.Context) bool
	FormatLeegoError(err error, middlewareName string) leego.LeegoError
}

type DefaultMiddleware struct{}

func (d DefaultMiddleware) Skipper(c leego.Context) bool {
	return defaultSkipper(c)
}

func (d DefaultMiddleware) FormatLeegoError(err error, middlewareName string) leego.LeegoError {
	return defaultFormatLeegoError(err, middlewareName)
}

func handlerName(h leego.HandlerFunc) string {
	t := reflect.ValueOf(h).Type()
	if t.Kind() == reflect.Func {
		return runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	}
	return t.String()
}
