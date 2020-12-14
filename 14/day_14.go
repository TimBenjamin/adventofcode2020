package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var programs [][]string

func get_programs() {
	f, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var line string
	var program []string
	for scanner.Scan() {
		line = scanner.Text()
		match, _ := regexp.MatchString(`^mask\s=`, line)
		if match {
			// start of a new program
			if len(program) > 0 {
				programs = append(programs, program)
			}
			program = []string{}
		}
		program = append(program, line)
	}
	// last one
	programs = append(programs, program)
}

func get_mask(program []string) (mask []string) {
	for _, line := range program {
		pair := strings.Split(line, " = ")
		if pair[0] == "mask" {
			_mask := strings.Split(pair[1], "")
			for _, m := range _mask {
				mask = append(mask, m)
			}
			break
		}
	}
	// check we have the mask
	if len(mask) == 0 {
		panic(errors.New("Could not find the mask in the program, aborting"))
	}
	return
}

func get_location(line string) (loc int, decimal_int int) {
	pair := strings.Split(line, " = ")
	// get the location:
	re := regexp.MustCompile(`^mem\[(\d+)\]`)
	match := re.FindStringSubmatch(pair[0])
	_loc := match[1]
	loc, _ = strconv.Atoi(_loc)
	decimal_int, _ = strconv.Atoi(pair[1])
	return
}

func read_binary(binary []string) (result int) {
	p := 1
	result = 0
	for i := len(binary) - 1; i >= 0; i-- {
		if binary[i] == "1" {
			result += p
		}
		p *= 2
	}
	return
}

func apply_mask(mask []string, padded_binary string) (split_binary []string) {
	split_binary = strings.Split(padded_binary, "")
	for i := 0; i < len(mask); i++ {
		if mask[i] == "1" {
			split_binary[i] = "1"
		} else if mask[i] == "0" {
			split_binary[i] = "0"
		}
	}
	return
}

// part 2 version
func apply_v2_mask(mask []string, binary_loc string) (masked_binary []string) {
	// we have to padd the location to length of mask, though
	binary_loc = lpad(binary_loc, "0", len(mask))
	masked_binary = strings.Split(binary_loc, "")
	//fmt.Println("mask:", mask)
	//fmt.Println("loc: ", masked_binary)
	bit_loc := len(masked_binary) - 1
	for i := len(mask) - 1; i >= 0; i-- {
		mask_bit := mask[i]
		//fmt.Printf("mask_bit at %v is %v - matching sb bit is %v\n", bit_loc, mask_bit, masked_binary[bit_loc])
		if mask_bit == "1" {
			masked_binary[bit_loc] = "1"
		} else if mask[i] == "0" {
			// unchanged
		} else {
			// here we have the weird behaviour of X
			masked_binary[bit_loc] = "X"
		}
		bit_loc--
		if bit_loc < 0 {
			break
		}
	}
	//fmt.Println("done:", masked_binary)
	return
}

// part 2
// add the variants as dictated by the positions of the X's
func add_locations(masked_location []string, value int) {
	//X1101X
	// first get the decimal value with all the X's replaced by 0's
	//011010 = 26
	// variations will be:
	// +1
	// +32
	// +1 +32
	var lower_bound []string
	var adds []int
	for i, bit := range masked_location {
		if bit == "X" {
			lower_bound = append(lower_bound, "0")
			add := int(math.Pow(2, float64(len(masked_location)-1-i)))
			adds = append(adds, add)
			// also add this number to each of the existing numbers in adds, and add those too
			for _, a := range adds {
				if a != add {
					adds = append(adds, a+add)
				}
			}
		} else {
			lower_bound = append(lower_bound, bit)
		}
	}

	base := read_binary(lower_bound)
	add_to_memory(base, value)

	// now the same with all the additional variants
	for _, add := range adds {
		location := base + add
		add_to_memory(location, value)
	}
}

func add_to_memory(location int, value int) {
	for i, d := range dec_locations {
		if d == location {
			// replace the corresponding value
			dec_values[i] = value
			return
		}
	}
	// didn't find it, so
	dec_locations = append(dec_locations, location)
	dec_values = append(dec_values, value)
}

// part 2
var dec_locations []int // the names of the locations (they are ints)
var dec_values []int    // the values we will add

func part_2() (answer int) {
	for _, program := range programs {
		mask := get_mask(program)
		for i := 1; i < len(program); i++ {
			loc, value := get_location(program[i])

			// in this part we want to write value to a bunch of different memory locations
			// these locations are based on applying the mask to "loc"
			binary_loc := strconv.FormatInt(int64(loc), 2)

			// now apply the mask to the binary location, using the funky v2 version
			masked_location := apply_v2_mask(mask, binary_loc)

			// get all the combinations that this triggers
			add_locations(masked_location, value)
		}
	}
	// now we should just be able to read off all those locations and add them
	for _, v := range dec_values {
		answer += v
	}
	return
}

// part 1
// we  need to put the result of the program into a specified memory location
// in the input the location can be 5 digits, but the space is not contiguous, so I'll treat it as key/value pairs
var memory [][]int

func part_1() (answer int) {
	for _, program := range programs {
		mask := get_mask(program)
		fmt.Println("mask:", mask)
		for i := 1; i < len(program); i++ {
			loc, decimal_int := get_location(program[i])
			fmt.Printf("loc: %v, value: %v\n", loc, decimal_int)

			binary := strconv.FormatInt(int64(decimal_int), 2)
			// pad it to 36 "bits" first to make masking easier
			padded_binary := lpad(binary, "0", len(mask))

			// now apply the mask to the value
			masked_binary := apply_mask(mask, padded_binary)

			// now read the split_binary into a decimal!
			result := read_binary(masked_binary)
			fmt.Println("result:", result)

			// OK we now need to put "result" into memory location "loc"
			// possibly replacing an existing value
			write_memory(loc, result)

		}
	}
	// finally sum up all the values we have in memory
	for _, memloc := range memory {
		answer += memloc[1]
	}
	return
}

func write_memory(loc int, result int) {
	existing_loc := false
	for i := 0; i < len(memory); i++ {
		if memory[i][0] == loc {
			memory[i][1] = result // overwrite
			existing_loc = true
			break
		}
	}
	if !existing_loc {
		memloc := []int{loc, result}
		memory = append(memory, memloc)
	}
}

func lpad(s string, pad string, plength int) string {
	for i := len(s); i < plength; i++ {
		s = pad + s
	}
	return s
}

// Run the program with the argument "2" to run part 2, or anything else for part 1.
func main() {
	if len(programs) == 0 {
		get_programs()
	}

	var answer int
	if len(os.Args) > 1 && os.Args[1] == "2" {
		answer = part_2()
	} else {
		answer = part_1()
	}

	fmt.Printf("The answer is %v\n", answer)
	// part 1: 8471403462063
}
