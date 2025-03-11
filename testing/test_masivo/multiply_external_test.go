package testmasivo_test

import (
	"testing"
	testmasivo "testpro"
)

func TestMultiply(t *testing.T)  {
	want := 6
	got := testmasivo.Multiply(2,3)
	if got != want {
		t.Errorf("se esperaba %d, se obtuvo %d", want, got)
	}
}