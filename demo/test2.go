package main

import "fmt"

func modify(array [5]int) (array2 [5]int) {
	for i, _ := range array {
		array[i] = i + i
	}

	fmt.Println("In modify(),array values is :", array)

	return array
}

func main2() {
	array := [5]int{1, 2, 3, 4, 5}

	array = modify(array)

	fmt.Println("In main(),array values is :", array)
}

var array = "hello world!"

func main1() {
	for i, v := range array {
		fmt.Printf("Array element[%d,%c]\r\n", i, v)
	}
}

// 数组切片用法
func main3() {
	var myArray [10]int = [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	var mySlice []int = myArray[2:7]

	fmt.Println("Elements of myArray: ")
	for _, v := range myArray {
		fmt.Print(v, " ")
	}

	fmt.Println("\nElements of mySlice: ")
	for _, v := range mySlice {
		fmt.Print(v, " ")
	}

	fmt.Println()
}

func main() {
	array1 := make([]int, 5)
	array2 := make([]int, 5, 10)
	array3 := []int{5, 4, 3, 2, 1}

	array4 := array2[5:]

	fmt.Println("array1 : ", array1, len(array1))
	fmt.Println("array2 : ", array2, len(array2), cap(array2))
	fmt.Println("array3 : ", array3, len(array3))

	array2 = append(array2, 1, 2, 3, 4, 5)

	array3 = append(array1, array3...)

	fmt.Println("array2 : ", array2, len(array2), cap(array2))

	fmt.Println("array3 : ", array3, len(array3), cap(array3))

	fmt.Println("array4 : ", array4, len(array4), cap(array4))

}
