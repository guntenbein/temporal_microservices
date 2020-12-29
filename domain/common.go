package domain

import (
	"reflect"
	"runtime"
	"strings"
)

func GetActivityName(i interface{}) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	parts := strings.Split(fullName, ".")
	lastPart := parts[len(parts)-1]
	parts2 := strings.Split(lastPart, "-")
	firstPart := parts2[0]
	return firstPart
}
