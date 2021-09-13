package main

import "testing"

func TestFoo(t *testing.T) {
	var tests = []struct {
		name string
		c    int
	}{
		{
			"ok",
			0,
		},
		{
			"error",
			1,
		},
		{
			"fatal",
			2,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.c == 1 {
				t.Error("error")
			}
			if tt.c == 2 {
				t.Fatal("fail")
			}
		})
	}
}
