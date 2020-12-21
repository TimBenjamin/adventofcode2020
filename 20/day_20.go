package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var real_input_file = "./input.txt"
var test_input_file = "./test_input.txt"
var input_file = real_input_file

const N = "n"
const E = "e"
const S = "s"
const W = "w"

// describe tiles by their flipped status, how many times rotated, and what their n, e, s, w edges are like (which will vary depending on rotations and flip)
type Tile struct {
	id        string
	edges     map[string]string
	flipped   bool
	rotations int // clockwise
}

func (t *Tile) rotate() {
	if t.rotations < 4 {
		t.rotations++
	} else {
		t.rotations = 0
	}
	w := t.edges["w"]
	t.edges["w"] = t.edges["s"]
	t.edges["s"] = t.edges["e"]
	t.edges["e"] = t.edges["n"]
	t.edges["n"] = w
}

func (t *Tile) flip() {
	t.flipped = !t.flipped
	w := t.edges["w"]
	t.edges["w"] = t.edges["e"]
	t.edges["e"] = w
	// reverse n and s
	t.edges["n"] = reverse(t.edges["n"])
	t.edges["s"] = reverse(t.edges["s"])
}

func (t *Tile) reset() {
	if t.flipped {
		t.flip()
	}
	for {
		if t.rotations == 0 {
			break
		}
		t.rotate()
	}
}

func reverse(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}

func get_input() (input [][]string) {
	f, err := os.Open(input_file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var line string
	var tile_data []string
	for scanner.Scan() {
		line = scanner.Text()
		tile_data = append(tile_data, line)
		if len(line) == 0 || line == "" || line == "\n" {
			input = append(input, tile_data)
			tile_data = []string{}
		}
	}
	// don't forget last one
	input = append(input, tile_data)
	return
}

// map the tiles by their file ID number
var tiles = map[string]Tile{}

// combine the edge elements into a string for comparison
func get_edge_from_tile_data(tile []string, edge string) string {
	var idx int
	if edge == "n" {
		return tile[0]
	} else if edge == "s" {
		return tile[len(tile)-1]
	} else if edge == "e" {
		idx = 0
	} else if edge == "w" {
		idx = len(tile[0]) - 1
	} else {
		panic(errors.New("Unexpected edge"))
	}
	edge_data := []string{}
	for _, tile_row := range tile {
		tile_row_split := strings.Split(tile_row, "")
		edge_data = append(edge_data, tile_row_split[idx])
	}
	return strings.Join(edge_data, "")
}

func parse_input(input [][]string) {
	for _, tile := range input {
		var tile_id string
		for row_num, row := range tile {
			matched, _ := regexp.MatchString(`\d+`, row)
			if matched {
				re := regexp.MustCompile(`(\d+)`)
				tile_id = re.FindStringSubmatch(row)[0]
				var tile_data []string
				for i := row_num + 1; i <= 10; i++ {
					tile_data = append(tile_data, tile[i])
				}
				edges := map[string]string{}
				edges["n"] = get_edge_from_tile_data(tile_data, "n")
				edges["e"] = get_edge_from_tile_data(tile_data, "e")
				edges["s"] = get_edge_from_tile_data(tile_data, "s")
				edges["w"] = get_edge_from_tile_data(tile_data, "w")
				tiles[tile_id] = Tile{
					tile_id, edges, false, 0,
				}
			}

		}
	}
}

func in_array(val string, array []string) (ok bool) {
	for _, i := range array {
		if ok = i == val; ok {
			return
		}
	}
	return
}

// finds tiles that match this tile's specified edge, and says which of their edges it is that matches
func find_matches(tile Tile, edge string) map[string]string {
	matches := map[string]string{}
	for check_id, check_tile := range tiles {
		if tile.id == check_id {
			continue
		}
		if tile.edges[edge] == check_tile.edges["n"] {
			matches[check_id] = "n"
		}
		if tile.edges[edge] == check_tile.edges["e"] {
			matches[check_id] = "e"
		}
		if tile.edges[edge] == check_tile.edges["s"] {
			matches[check_id] = "s"
		}
		if tile.edges[edge] == check_tile.edges["w"] {
			matches[check_id] = "w"
		}
	}
	return matches
}

func run() {

	immune_tiles := []string{} // ones that we match to a corner, don't rotate/flip them pls

	// try to find just the NE corner
	// want one that has matches on S:n and W:e (only!)
	// That "n" must have matches on n, w, and s only (not e)
	// That "e" must have matches on e, s, and w only (not n)
	fmt.Println("Look for candidates for the NE corner")
	for id, tile := range tiles {
		// don't want anything that matches on the N or E sides
		possible := true
		candidates_N := find_matches(tile, N)
		if len(candidates_N) > 0 {
			possible = false
		}
		candidates_E := find_matches(tile, E)
		if len(candidates_E) > 0 {
			possible = false
		}
		if !possible {
			continue
		}
		// now
		condition_1 := false
		//var candidate_1 int
		condition_2 := false
		//var candidate_2 int
		candidates_S := find_matches(tile, S)
		for _, edge := range candidates_S {
			if edge == N {
				condition_1 = true
			}
		}
		candidates_W := find_matches(tile, W)
		for _, edge := range candidates_W {
			if edge == E {
				condition_2 = true
			}
		}
		if condition_1 && condition_2 {
			fmt.Printf("  Tile %v works for the NE corner, matches S and W\n", id)
			immune_tiles = append(immune_tiles, id)
		}
	}

	// OK now find the SE corner
	fmt.Println("Look for candidates for the SE corner")
	for id, tile := range tiles {
		if in_array(id, immune_tiles) {
			continue
		}
		tile.flip()
		possible := true
		candidates_S := find_matches(tile, S)
		if len(candidates_S) > 0 {
			possible = false
		}
		candidates_E := find_matches(tile, E)
		if len(candidates_E) > 0 {
			possible = false
		}
		if !possible {
			continue
		}
		// now
		condition_1 := false
		condition_2 := false
		candidates_N := find_matches(tile, N)
		for _, edge := range candidates_N {
			if edge == S {
				condition_1 = true
			}
		}
		candidates_W := find_matches(tile, W)
		for _, edge := range candidates_W {
			if edge == E {
				condition_2 = true
			}
		}
		if condition_1 && condition_2 {
			fmt.Printf("  Tile %v works for the SE corner, matches N and W\n", id)
			immune_tiles = append(immune_tiles, id)
		}
		tile.reset()
	}

	// OK now find the SW corner
	fmt.Println("Look for candidates for the SW corner")
	for id, tile := range tiles {
		if in_array(id, immune_tiles) {
			continue
		}
		possible := true
		candidates_S := find_matches(tile, S)
		if len(candidates_S) > 0 {
			possible = false
		}
		candidates_W := find_matches(tile, W)
		if len(candidates_W) > 0 {
			possible = false
		}
		if !possible {
			continue
		}
		// now
		condition_1 := false
		condition_2 := false
		candidates_N := find_matches(tile, N)
		for _, edge := range candidates_N {
			if edge == S {
				condition_1 = true
			}
		}
		candidates_E := find_matches(tile, E)
		for _, edge := range candidates_E {
			if edge == W {
				condition_2 = true
			}
		}
		if condition_1 && condition_2 {
			fmt.Printf("  Tile %v works for the SE corner, matches N and W\n", id)
			immune_tiles = append(immune_tiles, id)
		}
	}

	// OK now find the SW corner
	fmt.Println("Look for candidates for the NW corner")
	for id, tile := range tiles {
		if in_array(id, immune_tiles) {
			continue
		}
		possible := true
		candidates_N := find_matches(tile, N)
		if len(candidates_N) > 0 {
			possible = false
		}
		candidates_W := find_matches(tile, W)
		if len(candidates_W) > 0 {
			possible = false
		}
		if !possible {
			continue
		}
		// now
		condition_1 := false
		condition_2 := false
		candidates_S := find_matches(tile, S)
		for _, edge := range candidates_S {
			if edge == N {
				condition_1 = true
			}
		}
		candidates_E := find_matches(tile, E)
		for _, edge := range candidates_E {
			if edge == W {
				condition_2 = true
			}
		}
		if condition_1 && condition_2 {
			fmt.Printf("  Tile %v works for the SE corner, matches N and W\n", id)
			immune_tiles = append(immune_tiles, id)
		}
	}

	// Part 1 - we only need corner tiles
	// "multiply the IDs of the four corner tiles together"
	// so that's tiles that only match 2 other tiles, SE, SW, NE, NW.
}

// Run the program with the argument "2" to run part 2, or anything else for part 1.
func main() {
	parse_input(get_input())
	fmt.Printf("There are %v tiles\n", len(tiles))
	run()
}
