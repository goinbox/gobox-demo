package validate

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

var validateFuncMap = map[string]validator.Func{}

func addValidateFuncMap(fm map[string]validator.Func) {
	for k, f := range fm {
		validateFuncMap[k] = f
	}
}

func Init() error {
	validate = validator.New()
	for tag, fn := range validateFuncMap {
		err := validate.RegisterValidation(tag, fn)
		if err != nil {
			return err
		}
	}

	return nil
}

func Validator() *validator.Validate {
	return validate
}
