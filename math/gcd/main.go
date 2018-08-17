package main

import (
	"fmt"
	"os"
	"strconv"
)

/* 迭代法（递推法）：欧几里得算法，计算最大公约数 */
func gcd1(m, n int) int {
	for {
		if m == 0 {
			return n
		}
		c := n % m
		n = m
		m = c
	}
}

/* 递归法：欧几里得算法，计算最大公约数 */
func gcd2(m, n int) int {
	if m == 0 {
		return n
	}
	return gcd2(n%m, m)
}

/* 连续整数试探算法，计算最大公约数 */
func gcd3(m, n int) int {
	if m > n {
		m, n = n, m
	}
	t := m
	for {
		if m%t == 0 && n%t == 0 {
			return t
		}
		t--
	}
}

func main() {

	parms := os.Args[1:]
	if len(parms) < 2 {
		fmt.Println("input parms num less then 2.")
		return
	}

	var nums []int

	for _, v := range parms {
		tmp, err := strconv.Atoi(v)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		nums = append(nums, tmp)
	}

	for _, v := range nums {
		fmt.Printf("%d ", v)
	}
	fmt.Println("")

	m := nums[0]
	for i := 1; i < len(nums); i++ {
		m = gcd1(m, nums[i])
	}
	fmt.Printf("gcd1 %d\r\n", m)

	m = nums[0]
	for i := 1; i < len(nums); i++ {
		m = gcd2(m, nums[i])
	}
	fmt.Printf("gcd2 %d\r\n", m)

	m = nums[0]
	for i := 1; i < len(nums); i++ {
		m = gcd3(m, nums[i])
	}
	fmt.Printf("gcd3 %d\r\n", m)
}
