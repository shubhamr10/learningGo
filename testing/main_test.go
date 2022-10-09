package main

import "testing"

//testing - manual test & tablel test

var tests = []struct {
	name     string
	dividend float32
	divisor  float32
	expected float32
	isErr    bool
}{
	{"valid-data", 10, 5, 2, false},
	{"invalid-data", 100, 0, 0, true},
	{"invalid-data", 100, 0, 0, true},
}

func TestDivision(t *testing.T) {
	// Table testing
	for _, tt := range tests {
		results, err := divide(tt.dividend, tt.divisor)
		if tt.isErr {
			if err == nil {
				t.Error("expected an error but did not get one")
			}
		} else {
			if err != nil {
				t.Error("didnot expect an error but got one", err.Error())
			}
		}

		if results != tt.expected {
			t.Errorf("expected %f, but got %f", tt.expected, results)
		}
	}
}

// Manual test
//func TestDivide(t *testing.T) {
//	_, err := divide(10, 1)
//	if err != nil {
//		t.Error("got an error we should not have!")
//	}
//}
//
//func TestBadDivide(t *testing.T) {
//	_, err := divide(10, 0)
//	if err == nil {
//		t.Error("Should get an error we should not have!")
//	}
//}

// to check the coverage in test
// go test -coverprofile=coverage.out && go tool -html=coverage.out
