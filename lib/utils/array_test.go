package utils

import (
	"reflect"
	"testing"
)


func Test_PartitionChannel(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		size   int
		output [][]int
	}{
		{
			"incomplete",
			[]int{1},
			1,
			[][]int{{1}},
		},
		{
			"complete",
			[]int{1, 2, 3, 4},
			2,
			[][]int{{1, 2}, {3, 4}},
		},
		{
			"incomplete with remainder",
			[]int{1, 2, 3},
			2,
			[][]int{{1, 2}, {3}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testChannel := make(chan int, len(tt.input))
			go func() {
				for _, v := range tt.input {
					testChannel <- v
				}
				close(testChannel)
			}()

			var output [][]int
			PartitionChannel(testChannel, tt.size)(func(batch []int) bool {
				output = append(output, batch)
				return true
			})

			if !reflect.DeepEqual(output, tt.output) {
				t.Errorf("got %v, want %v", output, tt.output)
			}
		})
	}
}


func Test_Map(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		mapper func(int) int
		output []int
	}{
		{
			"double",
			[]int{1, 2, 3},
			func(x int) int { return x * 2 },
			[]int{2, 4, 6},
		},
		{
			"increment",
			[]int{1, 2, 3},
			func(x int) int { return x + 1 },
			[]int{2, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := Map(tt.input, tt.mapper)
			if !reflect.DeepEqual(output, tt.output) {
				t.Errorf("got %v, want %v", output, tt.output)
			}
		})
	}
}
