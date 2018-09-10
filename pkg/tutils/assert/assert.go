package assert

import (
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// Truef fails the test if condition is false.
func Truef(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		tb.Fatalf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
	}
}

// True fails the test if condition is false.
func True(tb testing.TB, condition bool) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		tb.Fatalf("\033[31m%s:%d: expected condition to be true \033[39m\n\n", filepath.Base(file), line)
	}
}

// Falsef fails the test if condition is true.
func Falsef(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if condition {
		_, file, line, _ := runtime.Caller(1)
		tb.Fatalf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
	}
}

// False fails the test if condition is true.
func False(tb testing.TB, condition bool) {
	if condition {
		_, file, line, _ := runtime.Caller(1)
		tb.Fatalf("\033[31m%s:%d: expected condition to be false \033[39m\n\n", filepath.Base(file), line)
	}
}

// Okf fails the test if err is not nil
func Okf(tb testing.TB, err error, msg string, v ...interface{}) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		tb.Fatalf("\033[31m%s:%d: unexpected error: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
	}
}

// Ok fails the test if err is not nil.
func Ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		tb.Fatalf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
	}
}

// Errorf fails the test if err is nil
func Errorf(tb testing.TB, err error, msg string, v ...interface{}) {
	if err == nil {
		_, file, line, _ := runtime.Caller(1)
		tb.Fatalf("\033[31m%s:%d: expected error: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
	}
}

// Error fails the test if err is nil.
func Error(tb testing.TB, err error) {
	if err == nil {
		_, file, line, _ := runtime.Caller(1)
		tb.Fatalf("\033[31m%s:%d: expected error\033[39m\n\n", filepath.Base(file), line)
	}
}

// Equals fails the test if expected is not equal to actual.
func Equals(tb testing.TB, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		_, file, line, _ := runtime.Caller(1)
		tb.Fatalf("\033[31m%s:%d:\n\n\texpected: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, expected, actual)
	}
}

// NE fails the test if expected is equal to actual.
func NE(tb testing.TB, expected, actual interface{}) {
	if reflect.DeepEqual(expected, actual) {
		_, file, line, _ := runtime.Caller(1)
		tb.Fatalf("\033[31m%s:%d:\n\n\texpected not equal: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, expected, actual)
	}
}

// Nil fails the test if value is not nil.
func Nil(tb testing.TB, value interface{}) {
	if !isNil(value) {
		_, file, line, _ := runtime.Caller(1)
		tb.Fatalf("\033[31m%s:%d:\n\n\texpected nil, got: %#v\033[39m\n\n", filepath.Base(file), line, value)
	}
}

// NotNil fails the test if value is nil.
func NotNil(tb testing.TB, value interface{}) {
	if isNil(value) {
		_, file, line, _ := runtime.Caller(1)
		tb.Fatalf("\033[31m%s:%d:\n\n\texpected value, got nil\033[39m\n\n", filepath.Base(file), line)
	}
}

// isNil checks if a specified object is nil or not.
func isNil(object interface{}) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)
	kind := value.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		return true
	}

	return false
}
