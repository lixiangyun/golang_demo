package bubblesort

func BubbleSort(input []int) []int {
	//flag := true

	for i := 0; i < len(input)-1; i++ {
		for j := 0; j < len(input)-1; j++ {
			if input[j] > input[j+1] {
				input[j], input[j+1] = input[j+1], input[j]
			}
		}
	}

	return input
}
