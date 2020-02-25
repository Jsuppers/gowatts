package server

import (
	"reflect"
	"testing"
)

func Test_createOptions(t *testing.T) {
	firstOption := "Fixed - Open Rack"
	secondOption := "Fixed - Roof Mounted"
	type args struct {
		stringIndex string
		options     []string
	}
	tests := []struct {
		name string
		args args
		want []Option
	}{
		{
			"sets first index as selected if no index is given",
			args{options: []string{firstOption, secondOption}},
			[]Option{{0, firstOption, true}, {1, secondOption, false}},
		},
		{
			"correctly sets the selected index",
			args{stringIndex: "1", options: []string{firstOption, secondOption}},
			[]Option{{0, firstOption, false}, {1, secondOption, true}},
		},
		{
			"sets first index as selected if index is not a number",
			args{stringIndex: "not a number", options: []string{firstOption, secondOption}},
			[]Option{{0, firstOption, true}, {1, secondOption, false}},
		},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			if got := createOptions(test.args.stringIndex, test.args.options...); !reflect.DeepEqual(got, test.want) {
				t.Errorf("createOptions() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_createLabels(t *testing.T) {
	tests := []struct {
		name   string
		length int
		want   []int
	}{
		{"correctly creates a list of integers", 3, []int{0, 1, 2}},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			if got := createLabels(test.length); !reflect.DeepEqual(got, test.want) {
				t.Errorf("createLabels() = %v, want %v", got, test.want)
			}
		})
	}
}
