package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func make_seat_id(row, seat int) (seat_id int) {
	return row*8 + seat
}

func get_seat(line string) (row, seat int) {
	cur_max := 128
	cur_min := 0
	if len(line) == 10 {
		letters := strings.Split(line, "")
		for i := 0; i < 8; i++ {
			p := cur_min + ((cur_max - cur_min) / 2)
			if letters[i] == "F" {
				cur_max = p
			} else if letters[i] == "B" {
				cur_min = p
			}
		}
		row := cur_min
		seat_max := 8
		seat_min := 0
		for i := 7; i < 10; i++ {
			p := seat_min + ((seat_max - seat_min) / 2)
			if letters[i] == "L" {
				seat_max = p
			} else if letters[i] == "R" {
				seat_min = p
			}
		}
		seat := seat_min
		return row, seat
	}
	return
}

func part_1() (seat_id int) {
	f, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		row, seat := get_seat(line)
		test_seat_id := make_seat_id(row, seat)
		if test_seat_id > seat_id {
			seat_id = test_seat_id
		}
	}

	return
}

func part_2() (my_seat_id int) {
	f, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var line string
	min_row := 127
	max_row := 0
	m := [128][8]int{}
	for scanner.Scan() {
		line = scanner.Text()
		row, seat := get_seat(line)
		if row < min_row {
			min_row = row
		}
		if row > max_row {
			max_row = row
		}
		m[row][seat] = 1
	}

	// somewhere between min row and max row there is a row with a missing seat...
	var my_row int
	var my_seat int
	for i := min_row + 1; i < max_row; i++ {
		cur_row := m[i]
		for p, _ := range cur_row {
			if cur_row[p] != 1 {
				my_row = i
				my_seat = p
				break
			}
		}
	}

	if my_row > 0 && my_seat > 0 {
		my_seat_id = make_seat_id(my_row, my_seat)
	}
	return
}

func main() {
	fmt.Println("PART 1:")
	seat_id := part_1()
	fmt.Printf("The highest seat ID is %v\n", seat_id)

	fmt.Println("PART 2:")
	my_seat_id := part_2()
	fmt.Printf("My seat id is %v\n", my_seat_id)
}
