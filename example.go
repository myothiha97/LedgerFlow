package main

import (
	"fmt"
	"slices"
)

func slicesExample() {
	// var s []string
	var s []string // this syntax []stirng mean initilazed a slice with the data type of string, i know its kinda confusing with array syntax but it is what it is in go lang
	fmt.Println("uninit:", s, s == nil, len(s) == 0)

	s = make([]string, 3) //  and this is where u assign or create an actual slice with
	fmt.Println("emp:", s, "len:", len(s), "cap:", cap(s))

	s[0] = "a"
	s[1] = "b"
	s[2] = "c"
	fmt.Println("set:", s)
	fmt.Println("get:", s[2])

	fmt.Println("len:", len(s))

	s = append(s, "d")
	s = append(s, "e", "f")
	fmt.Println("apd:", s)

	c := make([]string, len(s))
	copy(c, s)
	fmt.Println("cpy:", c)

	l := s[2:5]
	fmt.Println("sl1:", l)

	l = s[:5]
	fmt.Println("sl2:", l)

	l = s[2:]
	fmt.Println("sl3:", l)

	t := []string{"g", "h", "i"}
	fmt.Println("dcl:", t)

	t2 := []string{"g", "h", "i"}
	if slices.Equal(t, t2) {
		fmt.Println("t == t2")
	}

	twoD := make([][]int, 3)
	for i := range 3 {
		innerLen := i + 1
		twoD[i] = make([]int, innerLen)
		for j := range innerLen {
			twoD[i][j] = i + j
		}
	}
	fmt.Println("2d: ", twoD)
}

func arrayExample() {
	var a [5]int // initialized the array of len 5
	fmt.Println("emp 1:", a)

	a[4] = 100
	// fmt.Println("set:", a)
	// fmt.Println("get:", a[4])
	// fmt.Println("len:", len(a))

	// b := [5]int{1, 2, 3, 4, 5}
	// fmt.Println("e1 dcl:", b)
	//
	// b = [...]int{1, 2, 3, 4, 5}
	// fmt.Println("e2 dcl:", b)

	b := [...]int{100, 3: 400, 6: 500}
	b = [...]int{100, 1: 1, 2: 2, 3: 13, 4: 400, 5: 500, 6: 600}
	fmt.Println("idx:", b, " ", "len :", len(b), "value of zero index : ", b[0])

	// var twoD [2][3]int
	// for i := range 2 {
	// 	for j := range 3 {
	// 		twoD[i][j] = i
	// 	}
	// }

	// if we assigned only one var from range b , then i indicate the index of each looping
	// but if assigned to two var, the both first one indicate the index and 2nd one indicate the value of index in looping
	for i := range b {
		fmt.Println(i)
	}

	// fmt.Println("2d: ", twoD)

	twoD := [2][4]int{
		{1, 2, 3, 4},
		{1, 2, 3, 4},
	}
	fmt.Println("twoD", twoD)
}

func main() {
	// slicesExample()
	arrayExample()
}
