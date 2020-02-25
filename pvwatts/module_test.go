package pvwatts

import (
	"os"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name  string
		setup func()
		want  API
	}{
		{"test no api key", func() {}, &pvWatts{"DEMO_KEY"}},
		{"test sets api key", func() { os.Setenv("PVWATTS_API_KEY", "testKey") }, &pvWatts{"testKey"}},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			test.setup()
			defer os.Unsetenv("PVWATTS_API_KEY")
			if got := New(); !reflect.DeepEqual(got, test.want) {
				t.Errorf("New() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_getDataSet(t *testing.T) {
	type args struct {
		longitude string
		latitude  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test unknown longitude", args{"fake longitude", "30"}, "intl"},
		{"test unknown latitude", args{"30", "fake latitude"}, "intl"},
		{"uses nsrdb at New York", args{"-73.935242", "40.730610"}, "nsrdb"},
		{"uses nsrdb at New Delhi", args{"77.216721", "28.644800"}, "nsrdb"},
		{"uses intl at Auckland", args{"174.763336", "-36.848461"}, "intl"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDataSet(tt.args.longitude, tt.args.latitude); got != tt.want {
				t.Errorf("getDataSet() = %v, want %v", got, tt.want)
			}
		})
	}
}
