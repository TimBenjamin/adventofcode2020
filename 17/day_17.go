package main

import (
	"fmt"
	"strconv"
	"strings"
)

// my input
/*
##......
.##...#.
.#######
..###.##
.#.###..
..#.####
##.####.
##..#.##
*/
var real_input = [][]string{
	[]string{"#", "#", ".", ".", ".", ".", ".", "."},
	[]string{".", "#", "#", ".", ".", ".", "#", "."},
	[]string{".", "#", "#", "#", "#", "#", "#", "#"},
	[]string{".", ".", "#", "#", "#", ".", "#", "#"},
	[]string{".", "#", ".", "#", "#", "#", ".", "."},
	[]string{".", ".", "#", ".", "#", "#", "#", "#"},
	[]string{"#", "#", ".", "#", "#", "#", "#", "."},
	[]string{"#", "#", ".", ".", "#", ".", "#", "#"},
}

// test input
/*
.#.
..#
###
*/
var test_input = [][]string{
	[]string{".", "#", "."},
	[]string{".", ".", "#"},
	[]string{"#", "#", "#"},
}

// key: z,y,x (outer to inner)
var test_input_map = map[string]bool{
	"0,0,1": true,
	"0,1,2": true,
	"0,2,0": true,
	"0,2,1": true,
	"0,2,2": true,
}

var input_map = map[string]bool{}
var input = real_input
var turn = 0

func prepare_input() {
	z := "0"
	var map_key string
	for i, y := range input {
		for j, x := range y {
			if x == "#" {
				map_key = z + "," + strconv.Itoa(i) + "," + strconv.Itoa(j)
				input_map[map_key] = true
			}
		}
	}
}

func print_map(map_to_print map[string]bool) {
	for key, value := range map_to_print {
		fmt.Printf("%v:\t%v\n", key, value)
	}
}

func get_coords(coordinate string) []int {
	coordinate_split := strings.Split(coordinate, ",")
	z, _ := strconv.Atoi(coordinate_split[0])
	y, _ := strconv.Atoi(coordinate_split[1])
	x, _ := strconv.Atoi(coordinate_split[2])
	return []int{z, y, x}
}

func get_coordinate(coords []int) string {
	coord_strings := []string{}
	coord_strings = append(coord_strings, strconv.Itoa(coords[0]))
	coord_strings = append(coord_strings, strconv.Itoa(coords[1]))
	coord_strings = append(coord_strings, strconv.Itoa(coords[2]))
	return strings.Join(coord_strings, ",")
}

// we'll then find out how many of these are known to be active in input_map
func get_surrounding_coordinates(coordinate string) []string {
	coords := get_coords(coordinate)
	x_coord := coords[2]
	y_coord := coords[1]
	z_coord := coords[0]
	coordinates := []string{}
	for z := -1; z <= 1; z++ {
		for y := -1; y <= 1; y++ {
			for x := -1; x <= 1; x++ {
				if x == 0 && y == 0 && z == 0 {
					// skip the one we are looking at
					continue
				}
				coordinate := get_coordinate([]int{z + z_coord, y + y_coord, x + x_coord})
				coordinates = append(coordinates, coordinate)
			}
		}
	}
	return coordinates
}

func get_surrounding_active_cells(coordinate string) []string {
	coords := get_coords(coordinate)
	x_coord := coords[2]
	y_coord := coords[1]
	z_coord := coords[0]
	active_cells := []string{}
	for z := -1; z <= 1; z++ {
		for y := -1; y <= 1; y++ {
			for x := -1; x <= 1; x++ {
				if x == 0 && y == 0 && z == 0 {
					// skip the one we are looking at
					continue
				}
				check_key := get_coordinate([]int{z + z_coord, y + y_coord, x + x_coord})
				_, ok := input_map[check_key]
				if ok {
					active_cells = append(active_cells, check_key)
				}
			}
		}
	}
	return active_cells
}

func solve() int {
	for {
		// turn:
		turn++

		// the rules
		/*
		   Simultaneously on each round:
		   - If a cube is active and exactly 2 or 3 of its neighbors are also active, the cube remains active.
		   - Otherwise, the cube becomes inactive.

		   - If a cube is inactive but exactly 3 of its neighbors are active, the cube becomes active.
		   - Otherwise, the cube remains inactive.
		*/
		input_map_copy := map[string]bool{} // start this turn with a blank map, only add what we have to
		for coordinate, _ := range input_map {
			fmt.Printf("============== %v ==============\n", coordinate)
			active_around_cells := get_surrounding_active_cells(coordinate)
			active_count := len(active_around_cells)
			fmt.Println("Found how many active around:", active_count)
			// delete or keep this active location per the rules
			// we are starting with an empty map, so just add it if the count is good
			// - If a cube is active and exactly 2 or 3 of its neighbors are also active, the cube remains active.
			// - Otherwise, the cube becomes inactive.
			if active_count == 2 || active_count == 3 {
				input_map_copy[coordinate] = true
			}
		}
		fmt.Println("Active cells because of the first rule in this round:", len(input_map_copy))
		print_map(input_map_copy)

		// now we have identified the cells which are still active in this round
		// we can find new cells that should be active according to the second rule
		// for each active cell, see if any of the inactive cells around it have exactly 3 active neighbours
		// - If a cube is inactive but exactly 3 of its neighbors are active, the cube becomes active.
		// - Otherwise, the cube remains inactive.
		keys_to_add := []string{}
		for coordinate, _ := range input_map {
			surrounding_cells := get_surrounding_coordinates(coordinate)
			// which of these are INACTIVE?
			for _, cell := range surrounding_cells {
				_, ok := input_map[cell]
				if !ok {
					// each inactive neighbour, we need to check
					// if there are exactly 3 active, then "cell" needs to be made active
					count := 0
					surrounding_cells := get_surrounding_coordinates(cell)
					for _, check_cell := range surrounding_cells {
						_, check_ok := input_map[check_cell]
						if check_ok {
							count++
						}
					}
					if count == 3 {
						// I'm not sure if input_map_copy gets modified before or after this current loop
						// so to be safe, keep the new keys seperate until the end....
						keys_to_add = append(keys_to_add, cell)
					}
				}
			}
		}

		// now add them to the map
		for _, k := range keys_to_add {
			input_map_copy[k] = true
		}

		input_map = map[string]bool{}
		for k, v := range input_map_copy {
			input_map[k] = v
		}

		fmt.Println("After turn, number of active cells", len(input_map))
		print_map(input_map)

		// when we reach the specified number of turns
		if turn == 6 {
			break
		}
	}

	return len(input_map)
}

// Run the program with the argument "2" to run part 2, or anything else for part 1.
func main() {
	prepare_input()
	print_map(input_map)
	answer := solve()
	fmt.Printf("The answer is %v\n", answer)
}
