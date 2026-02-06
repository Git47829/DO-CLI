//go:build ignore

package main

import (
	"fmt"
	"math"
)

func powerOf(x float64, y float64) float64 {
	if y == 0 {
		return 1
	} else if y == 1 {
		return x
	} else {
		_i := x
		for i := 1; float64(i) < y; i++ {
			x = x * _i
		}
		return x
	}
}

func add(x, y int) int {
	return x + y
}

func subtract(x, y int) int {
	return x - y
}

func multiply(x, y int) int {
	return x * y
}

func divide(x, y float32) float32 {
	return x / y
}

func pqFormula(x, y float64) (float64, float64) {
	var bruch float64 = -(x / 2)
	var wurzel float64 = math.Sqrt(powerOf(x/2, 2) - y)
	return bruch + wurzel, bruch - wurzel
}

func fibbonaci(number int) int {
	var f0 int = 0
	var f1 int = 1
	fn := 0

	for i := 1; i < number; i++ {
		fn = f1 + f0
		fmt.Println(fn)
		f0 = f1
		f1 = fn
	}

	return fn
}
