package helper

import (
	"regexp"
	"strings"
)

func CamelToSnake(input string) string {
	regex := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := regex.ReplaceAllString(input, "${1}_${2}")
	return strings.ToLower(snake)
}
