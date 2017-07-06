package main

import "fmt"
import "errors"

// Personinfo

type PersonInfo struct {
	ID      string
	Name    string
	Address string
}

func example(x int) (int, int) {
	if x == 0 {
		return 1, 0
	} else {
		return 0, 1
	}
}

func myfunc() {
	i := 0

HERE:

	fmt.Println(i)
	i++

	if i < 10 {
		goto HERE
	}
}

func myfun2() {
JLOOP:

	for j := 0; j < 5; j++ {
		for i := 0; i < 10; i++ {
			if i > 5 {
				break JLOOP
			}
			fmt.Println(j, i)
		}
	}
}

func add(a int, b int) (sum int, err error) {
	if a < 0 || b < 0 {
		err = errors.New("should be non-negative numbers")
		return -1, err
	}

	return a + b, nil
}

func argsmore(args ...int) {
	for _, v := range args {
		fmt.Print(v)
	}
	fmt.Println()
}

func main() {
	var personDB map[string]PersonInfo

	personDB = map[string]PersonInfo{
		"12": PersonInfo{"12", "Jack", "Room 1212,..."},
	}

	// personDB := make(map[string]PersonInfo)

	// input same data

	personDB["12345"] = PersonInfo{"12345", "Tom", "Room 1212,..."}
	personDB["1"] = PersonInfo{"1", "Jack", "Room 2121,...."}

	person, ok := personDB["12345"]

	// found one person success!

	if ok {
		fmt.Println("Found the person", person.Name, person.ID, person.Address)
	} else {
		fmt.Println("Can not found the person ID", "123")
	}

	person, ok = personDB["12"]

	// found one person failed!

	if ok {
		fmt.Println("Found the person", person.Name, person.ID, person.Address)
	} else {
		fmt.Println("Can not found the person ID", "123")
	}

	delete(personDB, "12")

	person, ok = personDB["12"]

	// found one person failed!

	if ok {
		fmt.Println("Found the person", person.Name, person.ID, person.Address)

		//return
	} else {
		fmt.Println("Can not found the person ID", "123")

		//return
	}

	myfunc()

	myfun2()

	sum, err := add(-1, -1)
	if nil != err {
		fmt.Println("err : ", err.Error())
	} else {
		fmt.Println("sum = ", sum)
	}

	sum, err = add(1, 1)
	if nil != err {
		fmt.Println("err : ", err.Error())
	} else {
		fmt.Println("sum = ", sum)
	}

	argsmore(1, 2, 3)

	var array = []int{9, 8, 7}

	argsmore(array...)

}
