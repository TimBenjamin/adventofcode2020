package main

import (
	"bufio"
	"fmt"
	"os"
)

func part_1() (sum int) {
	f, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var line string
	for scanner.Scan() {
		line = scanner.Text()

		if len(line) == 0 || line == "" || line == "\n" {
			// TODO

		} else {

		}
	}
	// don't forget last one

	return
}

func part_2() (sum int) {
	f, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	return
}

func main() {
	fmt.Println("PART 1:")
	answer := part_1()
	fmt.Printf("The answer is %v\n", answer)

	fmt.Println("PART 2:")
	answer = part_2()
	fmt.Printf("The answer is %v\n", answer)
}
