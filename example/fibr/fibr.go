package main

import "fmt"

func fibonacciRecursion(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacciRecursion(n-1) + fibonacciRecursion(n-2)
}

func main() {
	fmt.Println(fibonacciRecursion(10)) // 55
}
