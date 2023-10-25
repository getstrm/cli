package test

import (
	"reflect"
	"testing"
)

func TestDummy(t *testing.T) {
	in := 1
	out := 1
	if !reflect.DeepEqual(in, out) {
		t.Fail()
	}
}
