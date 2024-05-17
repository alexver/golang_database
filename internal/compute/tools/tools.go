package tools

import "regexp"

const REGEXP_ARG_TEMPLATE = `^[\w\*\/]+$`

// function to check command aguments
func ValidateArgument(value string) bool {
	match, err := regexp.Match(REGEXP_ARG_TEMPLATE, []byte(value))

	return err == nil && match
}
