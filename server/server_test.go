package server

import (
	"gowatts/pvwatts"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want Server
	}{
		{"creates new http server", &httpServer{pvwatts.New()}},
	}
	for _, tt := range tests {
		test := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, test.want) {
				t.Errorf("New() = %v, want %v", got, test.want)
			}
		})
	}
}
