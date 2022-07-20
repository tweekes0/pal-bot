package testing

import (
	"errors"
	"reflect"
	"testing"
)

func AssertError(t testing.TB, got, expected error) {
	if got != expected {
		for errors.Unwrap(got) != nil {
			got = errors.Unwrap(got)
		}

		if got != expected {
			t.Fatalf("got: %v, expected: %v", got, expected)
		}
	}
}

func AssertType(t testing.TB, got, expected interface{}) {
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got: %v, expected: %v", got, expected)
	}
}
