package form

import "context"

type validatorFlag uint8

const (
	// Vnone only does a form binding
	Vnone validatorFlag = 0
	// Vfilter performs form binding and filter
	Vfilter validatorFlag = 1 << iota
)

// Validate performs filter and validation
/**
 * @param {interface{}} v - non nil pointer of any struct that will be binded
 * @param {*http.Request} r - non nil pointer of http.Request
 */
func Validate(v interface{}) error {
	return ValidateFlag(Vfilter, v)
}

// ValidateFlag performs filter (optional) and validation
/**
 * @param {flag} f - flag of action that will be perform (Vnone or Vfilter)
 * @param {interface{}} v - non nil pointer of any struct that will be validated
 */
func ValidateFlag(f validatorFlag, v interface{}) error {
	if f&Vfilter != 0 {
		err := _filter.Struct(context.Background(), v)
		if err != nil {
			return err
		}
	}
	return _validator.Struct(v)
}
