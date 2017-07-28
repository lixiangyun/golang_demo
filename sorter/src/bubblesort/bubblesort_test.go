package bubblesort

import "testing"

func TestBubbleSort1(t *testing.T) {
	values := []int{5, 4, 3, 2, 1}

	values = BubbleSort(values)

	for i := 1; i < len(values); i++ {
		if values[i-1] > values[i] {
			t.Error("BubbleSort() failed. Got", values, "Expected 1 2 3 4 5.")
		}
	}
}

func TestBubbleSort2(t *testing.T) {
	values := []int{5, 5, 3, 2, 1}

	values = BubbleSort(values)

	for i := 1; i < len(values); i++ {
		if values[i-1] > values[i] {
			t.Error("BubbleSort() failed. Got", values, "Expected 1 2 3 5 5.")
		}
	}
}

func TestBubbleSort3(t *testing.T) {
	values := []int{5}

	values = BubbleSort(values)

	for i := 1; i < len(values); i++ {
		if values[i-1] > values[i] {
			t.Error("BubbleSort() failed. Got", values, "Expected 5.")
		}
	}
}

func BenchmarkSort(t *testing.B) {
	values := []int{5, 5, 3, 2, 1}

	BubbleSort(values)

}
