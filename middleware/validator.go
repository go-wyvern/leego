package middleware

import (
	"github.com/go-wyvern/Leego"
	"github.com/go-wyvern/i18n"
	"github.com/go-wyvern/validator"
)

const ValidateName = "Validate"

type (
	ValidatorConfig struct {
		Skipper Skipper

		FormatLeegoError func(error, string) leego.LeegoError

		Name string

		Validate *validator.Validator
	}
)

var (
	DefaultValidatorConfig = ValidatorConfig{
		Skipper:          defaultSkipper,
		FormatLeegoError: defaultFormatLeegoError,
		Name:             ValidateName,
	}
)

func Validator(v *validator.Validator, m Middleware) leego.MiddlewareFunc {
	c := ValidatorConfig{
		Skipper:          m.Skipper,
		FormatLeegoError: m.FormatLeegoError,
		Name:             ValidateName,
		Validate:         v,
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
			if config.Validate == nil {
				return next(c)
			}
			err := validator.Validate(c.Request().FormParams(), config.Validate)
			if err != nil {
				if pErr, ok := err.(*validator.ParamsError); ok {
					pErr.Text = i18n.Translate(pErr.Text, c.Language())
					return config.FormatLeegoError(pErr.Tr(), config.Name)
				} else {
					return config.FormatLeegoError(err, config.Name)
				}
			}
			err = validator.UrlValidator(c.GetParamsMap(), config.Validate)
			if err != nil {
				if pErr, ok := err.(*validator.ParamsError); ok {
					pErr.Text = i18n.Translate(pErr.Text, c.Language())
					return config.FormatLeegoError(pErr.Tr(), config.Name)
				} else {
					return config.FormatLeegoError(err, config.Name)
				}
			}
			return next(c)
		}
	}
}
