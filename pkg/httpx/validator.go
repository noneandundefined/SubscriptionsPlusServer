package httpx

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator"
)

var Validate *validator.Validate
var passwordRegex = regexp.MustCompile(`^[A-Za-z0-9!"#$%&'()*+\-./:;<=>?@[\\\]^_{|}~]{0,8}$`)

func init() {
	Validate = validator.New()

	_ = Validate.RegisterValidation("optional_uuid", func(fl validator.FieldLevel) bool {
		uuid := fl.Field().String()
		if uuid == "" {
			return true
		}

		return validator.New().Var(uuid, "uuid") == nil
	})

	_ = Validate.RegisterValidation("dpass", func(fl validator.FieldLevel) bool {
		pass := fl.Field().String()
		if pass == "" {
			return true
		}

		return passwordRegex.MatchString(pass)
	})
}

func ValidateMsg(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			switch fieldError.Tag() {
			case "required":
				return fmt.Sprintf("field %s is required", fieldError.Field())
			case "len":
				return fmt.Sprintf("field %s must be %s characters long", fieldError.Field(), fieldError.Param())
			case "numeric":
				return fmt.Sprintf("field %s must contain only numbers", fieldError.Field())
			case "uuid":
				return fmt.Sprintf("field %s must be a valid UUID", fieldError.Field())
			case "gt":
				return fmt.Sprintf("field %s must be greater than %s", fieldError.Field(), fieldError.Param())
			case "max":
				return fmt.Sprintf("field %s cannot exceed %s characters", fieldError.Field(), fieldError.Param())
			case "min":
				return fmt.Sprintf("field %s must be at least %s characters", fieldError.Field(), fieldError.Param())
			case "optional_uuid":
				return fmt.Sprintf("field %s must be a valid UUID or empty", fieldError.Field())
			default:
				return fmt.Sprintf("validation error in field %s", fieldError.Field())
			}
		}
	}

	return "request validation error"
}
