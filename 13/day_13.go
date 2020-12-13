package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// part 1 was solved using pencil and paper in a few minutes, doesn't need a program...

// test input:
//const input = "17,1"

//const input = "7,13,x,x,59,x,31,19" // solution: 1068781
//const input = "17,x,13,19" // solution: 3417, range 200-350
//const input = "67,7,59,61" // solution: 754018, range 10000-15000
//const input = "67,x,7,59,61" // solution: 779210, range 10000-15000
//const input = "67,7,x,59,61" // solution: 1261476
//const input = "1789,37,47,1889" // solution: 1202161486, range 500000-900000 (this sized batch is fine performance wise)

// my actual input:
const input = "41,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,37,x,x,x,x,x,911,x,x,x,x,x,x,x,x,x,x,x,x,13,17,x,x,x,x,x,x,x,x,23,x,x,x,x,x,29,x,827,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,19"

func main() {
	// we can jump in bigger increments than 1 to save time
	// first find where the first 2 buses line up
	// then we can keep jumping by the LCM of those 2 buses until we find a match for the 3rd bus
	// and so on
	// faster than the shitty brute force approach which I've kept below for posterity

	var buses []int
	var diffs []int // how far each number is from t, according to the problem text
	input_split := strings.Split(input, ",")
	for idx, bus := range input_split {
		if bus != "x" {
			b, _ := strconv.Atoi(bus)
			buses = append(buses, b)
			diffs = append(diffs, idx)
		}
	}
	fmt.Println("buses", buses)
	fmt.Println("diffs", diffs)

	// x is our answer for the time
	x := solve(buses[0], buses[1], diffs[1])
	jump := buses[0] * buses[1]

	// keep adding jump to x until we hit a (multiple of buses[2])+diffs[2]
	for b := 2; b < len(buses); b++ {
		for {
			x += jump
			test := x + diffs[b]
			if test%buses[b] == 0 {
				if b == len(buses)-1 {
					fmt.Println("The solution is:", x)
				} else {
					fmt.Println("got new x:", x)
					jump = jump * buses[b]
					fmt.Println("jump is now:", jump)
				}
				break
			}
		}

	}
}

// this isn't good using a cap of 1000, it should be an infinite loop with a break, but hey it works
func solve(bus_1 int, bus_2 int, diff int) (answer int) {
	var bus_1_values []int
	for i := 1; i <= 1000; i++ {
		bus_1_values = append(bus_1_values, i*bus_1)
	}
	j := 1
	for s := 0; s < 1000; s++ {
		test := j * bus_2
		for _, i := range bus_1_values {
			if i+diff == test {
				return i
			}
		}
		j++
	}
	fmt.Println("failed to find it after 1k rounds")
	return
}

func old_main() {
	// I think it reduces down to solving two equations only
	// we have an extra clue that t > 100000000000000

	var buses []int
	var diffs []int // how far each number is from t, according to the problem text
	input_split := strings.Split(input, ",")
	for idx, bus := range input_split {
		if bus != "x" {
			b, _ := strconv.Atoi(bus)
			buses = append(buses, b)
			diffs = append(diffs, idx)
		}
	}
	fmt.Println("buses", buses)
	fmt.Println("diffs", diffs)

	runs := 1000 // how many runs to do
	batch_size := 10000
	magic_offset := 1 // not sure how calculated, from reddit
	for i := 0; i < len(buses)-1; i++ {
		magic_offset *= buses[i]
	}
	fmt.Println("magic offset is:", magic_offset)

	// we can jump in bigger increments than 1 to save time
	lcm := buses[0]
	/*for i := 1; i < len(buses)-1; i++ {
		lcm *= buses[i]
	}*/
	fmt.Println("LCM:", lcm)

	for run := 0; run < runs; run++ {
		fmt.Println("Run:", run+1)

		var test_t_values []int // possible values for t

		for i := 0; i < batch_size; i++ {
			k := i * run
			test := magic_offset + (lcm * (k + 1))
			test_t_values = append(test_t_values, test)
		}
		fmt.Println("lowest t candidate:", test_t_values[0])
		fmt.Printf("BUS: %v - found %v values for t\n", buses[0], len(test_t_values))

		var candidate_t_values [][]int
		candidate_t_values = append(candidate_t_values, test_t_values)
		for bus_idx, bus := range buses {
			if bus_idx == 0 {
				continue
			}
			var bus_candidates []int
			for _, t := range test_t_values {
				diff := get_diff(t, bus)
				if diff == diffs[bus_idx] {
					// candidate for this bus found - but check it's also a candidate for the previous bus!
					for _, candidate := range candidate_t_values[bus_idx-1] {
						if candidate == t {
							//fmt.Printf("candidate for t found: %v\n", t)
							bus_candidates = append(bus_candidates, t)
							break
						}
					}
				}
			}
			fmt.Printf("BUS: %v - found %v values for t\n", bus, len(bus_candidates))
			candidate_t_values = append(candidate_t_values, bus_candidates)
			if len(bus_candidates) == 0 {
				// ending because len bus_candidates is zero, no point looking further
				break
			}
		}
		// the actual candidate should be in our final array...
		fmt.Println("Possible values for t:", candidate_t_values[len(candidate_t_values)-1])

		// it will be the first one that we care about
		if len(candidate_t_values[len(candidate_t_values)-1]) > 0 {
			fmt.Println("The solution is:", candidate_t_values[len(candidate_t_values)-1][0])
			break
		}
	}
}

func get_diff(t int, bus int) (diff int) {
	p := float64(t+bus) / float64(bus)
	q := int(math.Floor(p)) // q is the number of times t divides bus, +1
	r := q * bus
	return r - t
}
