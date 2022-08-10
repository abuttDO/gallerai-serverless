package main

import "testing"

func Test_doMigrations(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test_doMigrations",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doMigrations()
		})
	}
}
