package main

import "fmt"

func main() {
	for a := 1; a <= 9; a++ {
		for b := 0; b <= 9; b++ {
			fmt.Printf("%d%d (a*b+1)=%d ", a, b, a*b+1)
			if a*b+1 == 10*a+b {
				fmt.Print("!!!!!!!!!zzzzzzzzzzzz!!!!!!!!!")
			}
			fmt.Println()
		}
	}
}
