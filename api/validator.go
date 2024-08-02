package api

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"regexp"
)

func RegexValidator(pattern string) func(fl validator.FieldLevel) bool {

	return func(fl validator.FieldLevel) bool {
		switch v := fl.Field(); v.Kind() {
		case reflect.String:
			re := regexp.MustCompile(pattern)
			return re.MatchString(v.String())
		default:
			return false
		}
	}

}

func InitValidator() *validator.Validate {
	validate := validator.New()

	validate.RegisterValidation("retailer", RegexValidator(`^[\w\s\-&]+$`))
	validate.RegisterValidation("description", RegexValidator(`^[\w\s\-]+$`))
	validate.RegisterValidation("currency", RegexValidator(`^\d+\.\d{2}$`))
	validate.RegisterValidation("time", RegexValidator(`^([01][0-9]|2[0-3]):[0-5][0-9]$`))

	return validate
}
