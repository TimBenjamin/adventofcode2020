package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var seatmap [][]string

func get_seatmap() {
	f, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		row := strings.Split(line, "")
		seatmap = append(seatmap, row)
	}

	return
}

func print_seatmap() {
	fmt.Printf("\n")
	for i, row := range seatmap {
		fmt.Printf("%v:  ", i)
		for _, seat := range row {
			fmt.Printf(" %v ", seat)
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

// function to run when someone arrives at a seat.
// Rules:
/*
If a seat is empty (L) and there are no occupied seats adjacent to it, the seat becomes occupied.
If a seat is occupied (#) and four or more seats adjacent to it are also occupied, the seat becomes empty.
Otherwise, the seat's state does not change.
Floor never changes.
Adjacent = horizontally, vertically, diagonally
*/
func seat_decision(i int, j int) string {
	current_status := seatmap[i][j]
	new_status := current_status

	if current_status == "." {
		new_status = current_status
	} else {

		seats_to_test := []string{}

		// row above: left middle right
		if i > 0 {
			if j > 0 {
				seats_to_test = append(seats_to_test, seatmap[i-1][j-1])
			}
			seats_to_test = append(seats_to_test, seatmap[i-1][j])
			if j < len(seatmap[i])-1 {
				seats_to_test = append(seats_to_test, seatmap[i-1][j+1])
			}
		}
		// middle row: left middle right
		if j > 0 {
			seats_to_test = append(seats_to_test, seatmap[i][j-1])
		}
		// middle middle is the seat we are testing itself, so ignore!
		if j < len(seatmap[i])-1 {
			seats_to_test = append(seats_to_test, seatmap[i][j+1])
		}
		// row below: left middle right
		if i < len(seatmap)-1 {
			if j > 0 {
				seats_to_test = append(seats_to_test, seatmap[i+1][j-1])
			}
			seats_to_test = append(seats_to_test, seatmap[i+1][j])
			if j < len(seatmap[i])-1 {
				seats_to_test = append(seats_to_test, seatmap[i+1][j+1])
			}
		}

		var flip bool
		if current_status == "L" {
			// becomes # if there are NO OCCUPIED SEATS in seats_to_test
			flip = true
			for _, test_seat := range seats_to_test {
				if test_seat == "#" {
					flip = false
					break
				}
			}
		} else if current_status == "#" {
			// becomes L if FOUR OR MORE in seats_to_test are "#"
			count := 0
			flip = false
			for _, test_seat := range seats_to_test {
				if test_seat == "#" {
					count++
				}
			}
			if count >= 4 {
				flip = true
			}
		}

		// finally
		if flip {
			if current_status == "L" {
				new_status = "#"
			} else if current_status == "#" {
				new_status = "L"
			}
		}
	}
	return new_status
}

func part_1() (result int) {
	if len(seatmap) == 0 {
		get_seatmap()
	}

	round_counter := 1
	// infinite loop for the rounds, we'll break when there are no more changes.
	for {
		changes := 0
		var new_seatmap [][]string
		for i, row := range seatmap {
			var new_row []string
			for j, seat := range row {
				new_seat := seat
				if seat != "." {
					decision := seat_decision(i, j)
					if decision != seat {
						new_seat = decision
						changes++
					}
				}
				new_row = append(new_row, new_seat)
			}
			new_seatmap = append(new_seatmap, new_row)
		}
		copy(seatmap, new_seatmap)

		//fmt.Printf("Seatmap after round %v - there were %v changes:\n", round_counter, changes)
		//print_seatmap()

		if changes == 0 {
			break
		}
		round_counter++
	}

	// The return value is how many seats are occupied at the end.
	fmt.Printf("Ended after %v rounds\n", round_counter)
	return count_occupied_seats()
}

func count_occupied_seats() (num int) {
	for _, row := range seatmap {
		for _, seat := range row {
			if seat == "#" {
				num++
			}
		}
	}
	return
}

// function to run when someone arrives at a seat.
// Part 2 version!
// Rules:
/*
If a seat is empty (L) and there are no occupied seats adjacent to it, the seat becomes occupied.
If a seat is occupied (#) and FIVE or more seats adjacent to it are also occupied, the seat becomes empty.
Otherwise, the seat's state does not change.
Floor never changes.
Adjacent = horizontally, vertically, diagonally
*/
func part_2_seat_decision(i int, j int) string {
	current_status := seatmap[i][j]
	new_status := current_status

	if current_status == "." {
		new_status = current_status
	} else {

		seats_to_test := []string{}

		// look in the column above
		for r := i - 1; r >= 0; r-- {
			if seatmap[r][j] != "." {
				seats_to_test = append(seats_to_test, seatmap[r][j])
				break
			}
		}

		// look in the column below
		for r := i + 1; r < len(seatmap); r++ {
			if seatmap[r][j] != "." {
				seats_to_test = append(seats_to_test, seatmap[r][j])
				break
			}
		}

		// look left:
		for c := j - 1; c >= 0; c-- {
			if seatmap[i][c] != "." {
				seats_to_test = append(seats_to_test, seatmap[i][c])
				break
			}
		}

		// look right:
		for c := j + 1; c < len(seatmap[i]); c++ {
			if seatmap[i][c] != "." {
				seats_to_test = append(seats_to_test, seatmap[i][c])
				break
			}
		}

		// diagonal left and up:
		c := j - 1
	outerlu:
		for r := i - 1; r >= 0; r-- {
			if c < 0 {
				break outerlu
			}
			if seatmap[r][c] != "." {
				seats_to_test = append(seats_to_test, seatmap[r][c])
				break outerlu
			}
			c--
		}

		// diagonal left and down:
		c = j - 1
	outerld:
		for r := i + 1; r < len(seatmap); r++ {
			if c < 0 {
				break outerld
			}
			if seatmap[r][c] != "." {
				seats_to_test = append(seats_to_test, seatmap[r][c])
				break outerld
			}
			c--
		}

		// diagonal right and up:
		c = j + 1
	outerru:
		for r := i - 1; r >= 0; r-- {
			if c >= len(seatmap[i]) {
				break outerru
			}
			if seatmap[r][c] != "." {
				seats_to_test = append(seats_to_test, seatmap[r][c])
				break outerru
			}
			c++
		}

		// diagonal right and down:
		c = j + 1
	outerrd:
		for r := i + 1; r < len(seatmap); r++ {
			if c >= len(seatmap[i]) {
				break outerrd
			}
			if seatmap[r][c] != "." {
				seats_to_test = append(seats_to_test, seatmap[r][c])
				break outerrd
			}
			c++
		}

		var flip bool
		if current_status == "L" {
			// becomes # if there are NO OCCUPIED SEATS in seats_to_test
			flip = true
			for _, test_seat := range seats_to_test {
				if test_seat == "#" {
					flip = false
					break
				}
			}
		} else if current_status == "#" {
			// becomes L if FIVE OR MORE in seats_to_test are "#"
			count := 0
			flip = false
			for _, test_seat := range seats_to_test {
				if test_seat == "#" {
					count++
				}
			}
			if count >= 5 {
				flip = true
			}
		}

		// finally
		if flip {
			if current_status == "L" {
				new_status = "#"
			} else if current_status == "#" {
				new_status = "L"
			}
		}
	}
	return new_status
}

func part_2() (sum int) {
	if len(seatmap) == 0 {
		get_seatmap()
	}

	round_counter := 1
	// infinite loop for the rounds, we'll break when there are no more changes.
	for {

		changes := 0
		var new_seatmap [][]string
		for i, row := range seatmap {
			var new_row []string
			for j, seat := range row {
				new_seat := seat
				if seat != "." {
					decision := part_2_seat_decision(i, j)
					if decision != seat {
						new_seat = decision
						changes++
					}
				}
				new_row = append(new_row, new_seat)
			}
			new_seatmap = append(new_seatmap, new_row)
		}
		copy(seatmap, new_seatmap)

		//fmt.Printf("Seatmap after round %v:\n", round_counter)
		//print_seatmap()

		if changes == 0 {
			break
		}
		round_counter++
	}

	// The return value is how many seats are occupied at the end.
	fmt.Printf("Ended after %v rounds\n", round_counter)
	return count_occupied_seats()
}

// Run the program with the argument "2" to run part 2, or anything else for part 1.
func main() {
	if len(os.Args) > 1 && os.Args[1] == "2" {
		fmt.Println("PART 2:")
		answer := part_2()
		fmt.Printf("The answer is %v\n", answer)
	} else {
		fmt.Println("PART 1:")
		answer := part_1()
		fmt.Printf("The answer is %v\n", answer)
	}
}
