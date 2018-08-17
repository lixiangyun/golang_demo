package main

import (
	"fmt"
	"os"
	"strconv"
)

/* 迭代法（递推法）：欧几里得算法，计算最大公约数 */
func gcd(m, n int) int {
	for {
		if m == 0 {
			return n
		}
		c := n % m
		n = m
		m = c
	}
}

type Item struct {
	SelectCnt   int
	WeightValue int
}

type Select struct {
	List []Item
	Idx  int

	Gcd       int
	MaxWeight int
	CurWeight int
}

func NewSelect(list []Item) *Select {
	s := &Select{List: make([]Item, len(list))}
	for idx, v := range list {
		s.List[idx] = v
		if s.MaxWeight < v.WeightValue {
			s.MaxWeight = v.WeightValue
		}
		if idx == 0 {
			s.Gcd = v.WeightValue
		} else {
			s.Gcd = gcd(s.Gcd, v.WeightValue)
		}
	}
	s.CurWeight = s.MaxWeight
	return s
}

func (s *Select) Reset() {
	s.CurWeight = s.MaxWeight
	s.Idx = 0
	for idx, _ := range s.List {
		s.List[idx].SelectCnt = 0
	}
}

func (s *Select) Choose() *Item {

	for {
		s.Idx = (s.Idx + 1) % len(s.List)
		if s.Idx == 0 {
			s.CurWeight = s.CurWeight - s.Gcd
			if s.CurWeight <= 0 {
				s.CurWeight = s.MaxWeight
			}
		}
		if s.List[s.Idx].WeightValue >= s.CurWeight {
			s.List[s.Idx].SelectCnt++
			return &s.List[s.Idx]
		}
	}
}

func main() {

	parms := os.Args[1:]
	if len(parms) < 2 {
		fmt.Println("input parms num less then 2.")
		return
	}

	list := make([]Item, 0)
	totalWeight := 0
	for _, v := range parms {
		tmp, err := strconv.Atoi(v)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		list = append(list, Item{WeightValue: tmp})
		totalWeight += tmp
	}
	slt := NewSelect(list)
	for totalcnt := 1000; totalcnt < 1000000000; totalcnt *= 10 {
		fmt.Printf("\r\nretry %d cnt\r\n", totalcnt)
		slt.Reset()
		for i := 0; i < totalcnt; i++ {
			slt.Choose()
		}
		for _, v := range slt.List {
			fmt.Printf("w:%d,%02.1f%% c:%d,%02.1f%% \r\n",
				v.WeightValue,
				100*float32(v.WeightValue)/float32(totalWeight),
				v.SelectCnt,
				100*float32(v.SelectCnt)/float32(totalcnt))
		}
	}
}
