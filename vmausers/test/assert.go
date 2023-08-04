package test

import (
	"reflect"
	"testing"
)

// AssertEqual checks if values are equal
func AssertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}
	// debug.PrintStack()
	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}

// AssertNotEqual checks if values are equal
func AssertNotEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		return
	}
	// debug.PrintStack()
	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}

// AssertFail always returns fail
func AssertFail(t *testing.T, a interface{}, b interface{}) {
	// debug.PrintStack()
	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}

// AssertSuccess always returns success
func AssertSuccess(t *testing.T, a interface{}, b interface{}) {
	return
	// debug.PrintStack()
}
