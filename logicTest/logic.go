package main

import "fmt"

func main() {
	fmt.Println("=== BAGIAN 1: LOGIC TEST ===\n")

	fmt.Println("Soal 1: Deret Fibonacci")
	fibonacci()
	fmt.Println("\n")

	fmt.Println("Soal 2: Pola Piramida")
	piramidaBintang()
	fmt.Println()

	fmt.Println("Soal 3: Pola Hollow Diamond")
	polaBintangLOGIC()
}

func fibonacci() {
	n := 10
	a, b := 1, 1

	for i := 0; i < n; i++ {
		if i < n-1 {
			fmt.Print(a, " ")
		} else {
			fmt.Print(a)
		}

		a, b = b, a+b
	}
}

func piramidaBintang() {
	n := 5
	mid := n/2 + 1

	for i := 1; i <= n; i++ {
		var spaces, stars int

		if i <= mid {

			spaces = mid - i
			stars = 2*i - 1
		} else {

			spaces = i - mid
			stars = 2*(n-i) + 1
		}

		for j := 0; j < spaces; j++ {
			fmt.Print(" ")
		}

		for j := 0; j < stars; j++ {
			fmt.Print("*")
		}

		fmt.Println()
	}
}

func polaBintangLOGIC() {
	height := 13
	width := 21
	mid := height / 2

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {

			centerX := width / 2

			if i <= mid {

				leftBorder := centerX - i
				rightBorder := centerX + i

				if i == mid {
					if j < leftBorder || j > rightBorder {
						fmt.Print("*")
					} else if j == leftBorder || j == rightBorder {
						fmt.Print("*")
					} else if (j-leftBorder)%2 == 1 {
						fmt.Print("-")
					} else {
						fmt.Print(" ")
					}
				} else {

					if j < leftBorder || j > rightBorder {
						fmt.Print("*")
					} else {
						fmt.Print(" ")
					}
				}
			} else {

				mirrorI := height - i - 1
				leftBorder := centerX - mirrorI
				rightBorder := centerX + mirrorI

				if j < leftBorder || j > rightBorder {
					fmt.Print("*")
				} else {
					fmt.Print(" ")
				}
			}
		}

		fmt.Println()
	}
}
