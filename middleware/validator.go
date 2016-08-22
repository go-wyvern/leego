package middleware

import (
	"github.com/go-wyvern/validator"
	"github.com/go-wyvern/Leego"
)

const ValidateName = "Validate"

type (
	ValidatorConfig struct {
		Skipper          Skipper

		FormatLeegoError func(error, string) leego.LeegoError

		Name             string

		Validate         *validator.Validator
	}
)

var (
	DefaultValidatorConfig = ValidatorConfig{
		Skipper: defaultSkipper,
		FormatLeegoError: defaultFormatLeegoError,
		Name:ValidateName,
	}
)

func Validator(m Middleware) leego.MiddlewareFunc {
	c := ValidatorConfig{
		Skipper:m.Skipper,
		FormatLeegoError:m.FormatLeegoError,
		Name:ValidateName,
	}
	return ValidatorWithConfig(c)
}

func ValidatorWithConfig(config ValidatorConfig) leego.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultValidatorConfig.Skipper
	}
	return func(next leego.HandlerFunc) leego.HandlerFunc {
		return func(c leego.Context) leego.LeegoError {
			if config.Skipper(c) {
				return next(c)
			}
			config.Validate = validator.Find(c.Request().Method(), c.Path())
			if config.Validate==nil{
				return next(c)
			}
			err := validator.Validate(c.Request().FormParams(), config.Validate)
			if err != nil {
				return config.FormatLeegoError(err, config.Name)
			}
			err = validator.UrlValidator(c.GetParamsMap(), config.Validate)
			if err != nil {
				return config.FormatLeegoError(err, config.Name)
			}
			return next(c)
		}
	}
}