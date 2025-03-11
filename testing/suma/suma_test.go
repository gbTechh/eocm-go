package suma

import "testing"

func TestAdd(t *testing.T) {
	want := 5
	//t.Logf("want value: %d\n", want)
	got := Add(2,3)
	//t.Logf("got value: %d\n", got)
	if got != want {
		t.Errorf("Se esperaba %d, se obtuvo %d", want, got)
	}
	//t.Log("Terminó la prueba Add")
}

func TestAddMultiple(t *testing.T) {
	want := 6
	got := AddMultiple(1,2,3)
	if got != want {
		t.Errorf("Se esperaba %d, se obtuvo %d", want, got)
	}
}