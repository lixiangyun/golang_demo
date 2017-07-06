package qsort

func quicksort(values []int, left, right int) {
	tmp := values[left]

	p := left

	i, j := left, right

	for i <= j {
		for j >= p && values[j] >= tmp {
			j--
		}

		if j >= p {
			values[p] = values[j]
			p = j
		}

		if values[i] <= tmp && i <= p {
			i++
		}

		if i <= p {
			values[p] = values[i]
			p = i
		}
	}

	values[p] = tmp

	if p-left > 1 {
		quicksort(values, left, p-1)
	}

	if right-p > 1 {
		quicksort(values, p+1, right)
	}
}

func QuickSort(values []int) []int {
	quicksort(values, 0, len(values)-1)

	return values
}
