package server

import "strconv"

// Option provides options for the select template
type Option struct {
	Value    int
	Text     string
	Selected bool
}

func createOptions(stringIndex string, options ...string) []Option {
	selectedIndex, err := strconv.Atoi(stringIndex)
	if err != nil {
		selectedIndex = 0
	}

	output := make([]Option, len(options))
	for index, option := range options {
		selected := false
		if index == selectedIndex {
			selected = true
		}
		output[index] = Option{index, option, selected}
	}
	return output
}

func createLabels(length int) []int {
	labels := make([]int, length)
	for i := 0; i < len(labels); i++ {
		labels[i] = i
	}
	return labels
}
