package utils

import (
	"strings"

	"github.com/sirupsen/logrus"
)

func IsEmptyString(value string) bool {
	return len(strings.TrimSpace(value)) == 0
}

func IsErrorDoPanic(e error) {
	if e != nil {
		logrus.Panicln(e)
	}
}

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
