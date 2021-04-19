package validation

import (
	fmt "fmt"

	color "github.com/fatih/color"
	validator "github.com/go-playground/validator/v10"
	conditional "github.com/mileusna/conditional"
)

// ValidatorUtilsValidateStructOptions options for ValidatorUtils.validateStruct
type ValidatorUtilsValidateStructOptions struct {
	PrefixError  string
	Validate     *validator.Validate
	TheStruct    interface{}
	PanicOnError bool
}

// ValidatorUtils utils for the main return
type ValidatorUtils interface {
	ValidateStruct(x ValidatorUtilsValidateStructOptions) *validator.ValidationErrors
}

// ValidatorUtilsInstance implements ValidatorUtils
type ValidatorUtilsInstance struct {
	ValidatorUtils
}

// ValidateStruct implementation of ValidatorUtils.ValidateStruct
func (x ValidatorUtilsInstance) ValidateStruct(options ValidatorUtilsValidateStructOptions) *validator.ValidationErrors {
	errs := options.Validate.Struct(options.TheStruct)

	if errs != nil {
		errsValidation := errs.(validator.ValidationErrors)

		for _, err := range errsValidation {
			// fmt.Println(err.Namespace())
			// fmt.Println(err.Field())
			// fmt.Println(err.StructNamespace())
			// fmt.Println(err.StructField())
			// fmt.Println(err.Tag())
			// fmt.Println(err.ActualTag())
			// fmt.Println(err.Kind())
			// fmt.Println(err.Type())
			// fmt.Println(err.Value())
			// fmt.Println(err.Param())

			actualTag := err.ActualTag()
			colorError := color.New(color.FgRed)

			colorError.Println(
				options.PrefixError +
					err.StructNamespace() +
					" is " +
					actualTag +
					" " +
					conditional.String(actualTag == "oneof", err.Param(), err.Kind().String()) +
					" - got: \"" +
					fmt.Sprintf("%v", err.Value()) +
					"\"",
			)
		}

		if options.PanicOnError {
			panic("Struct validation failed.")
		}

		return &errsValidation
	}

	return nil
}

// GetValidatorReturn returned from GetValidator
type GetValidatorReturn struct {
	Validate *validator.Validate
	Utils    *ValidatorUtilsInstance
}

// GetValidator validation library for the app
func GetValidator() GetValidatorReturn {
	return GetValidatorReturn{
		Validate: validator.New(),
		Utils:    &ValidatorUtilsInstance{},
	}
}
