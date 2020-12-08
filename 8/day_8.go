package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func get_program() (program [][]string) {
	f, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		/*
			The program looks like:
				acc +3
				jmp -3
				acc -99
		*/
		step := strings.Split(line, " ")
		program = append(program, step)
	}
	return
}

func check_log(log []int, step int) (int, error) {
	for _, p := range log {
		if p == step {
			return p, errors.New("Previously used!")
		}
	}
	return 0, nil
}

func run_program(program [][]string) (accumulator int, err error) {
	var log []int
	line := 0
	step := program[line]
	for {
		amount, _ := strconv.Atoi(step[1])
		switch step[0] {
		case "acc":
			accumulator += amount
			fmt.Printf("Line %v:\tACC\tChanging accumulator by %v (it is now %v)\n", line, amount, accumulator)
			line += 1
		case "jmp":
			fmt.Printf("Line %v:\tJMP\tMoving program location by %v steps\n", line, amount)
			line += amount
		case "nop":
			fmt.Printf("Line %v:\tNOP\tMoving to next step\n", line)
			line += 1
		}

		if line > len(program)-1 {
			fmt.Println("Success! Reached the end of the program")
			return
		}

		previous_step, log_error := check_log(log, line)
		if log_error != nil {
			fmt.Printf("Infinite loop detected, hit step %v for the second time!\n", previous_step)
			err = errors.New("Infinite loop detected")
			return
		}
		log = append(log, line)
		step = program[line]
	}
}

func part_1() int {
	program := get_program()
	accumulator, err := run_program(program)
	if err != nil {
		fmt.Println("Program failed with error:", err)
	}
	return accumulator
}

// Same as Part 1 but we allow one change of a JMP to a NOP (or vice versa) to see if that prevents a crash
func part_2() int {
	program := get_program()
	accumulator, err := run_program(program)
	var false_corrections []int
	var new_step []string
	if err != nil {
		fmt.Println("Program failed with error:", err)
		// restart the program but avoid an error correction on the first available NOP or JMP
		correction := false
		for line_number, step := range program {
			if len(false_corrections) > 0 && line_number < false_corrections[len(false_corrections)-1] {
				// keep going
				continue
			} else {
				if step[0] == "jmp" {
					new_step = []string{"nop", step[1]}
					correction = true
				} else if step[0] == "nop" {
					new_step = []string{"jmp", step[1]}
					correction = true
				} else {
					// keep going until we get a nop or acc to change
					continue
				}
			}
			if correction {
				new_program := get_program()
				new_program[line_number] = new_step
				false_corrections = append(false_corrections, line_number)
				fmt.Printf("Attempt re-run correcting line %v to: %v\n", line_number, new_step)
				accumulator, err := run_program(new_program)
				if err != nil {
					// we must try again
					fmt.Printf("No success modifying line %v - try again\n", line_number)
				} else {
					// success!
					return accumulator
				}
			}
		}
		fmt.Println("Reached the end of the program, no more corrections can be made...")

	}
	return accumulator
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
		fmt.Printf("The accumulator at stop is %v\n", answer)
	}
}
