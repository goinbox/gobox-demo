package validate

import (
	"github.com/go-playground/validator/v10"

	"gdemo/model/demo"
)

func init() {
	addValidateFuncMap(map[string]validator.Func{
		"demo_status": validateDemoStatus,
	})
}

func validateDemoStatus(fl validator.FieldLevel) bool {
	v := fl.Field().Int()
	if v == demo.StatusOnline || v == demo.StatusOffline {
		return true
	}
	return false
}
