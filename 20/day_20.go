package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var real_input_file = "./input.txt"
var test_input_file = "./test_input.txt"
var input_file = real_input_file

const N = "n"
const E = "e"
const S = "s"
const W = "w"

var corner_tiles_real = []string{"2273", "2243", "2953", "1213"}
var corner_tiles_test = []string{"1951", "2971", "3079", "1171"}
var corner_tiles = corner_tiles_real

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
	w := t.edges[W]
	t.edges[W] = t.edges[S]
	t.edges[S] = t.edges[E]
	t.edges[E] = t.edges[N]
	t.edges[N] = w
}

func (t *Tile) flip() {
	t.flipped = !t.flipped
	w := t.edges[W]
	t.edges[W] = t.edges[E]
	t.edges[E] = w
	// reverse n and s
	t.edges[N] = reverse(t.edges[N])
	t.edges[S] = reverse(t.edges[S])
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
	if edge == N {
		return tile[0]
	} else if edge == S {
		return tile[len(tile)-1]
	} else if edge == E {
		idx = 0
	} else if edge == W {
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
				edges[N] = get_edge_from_tile_data(tile_data, N)
				edges[E] = get_edge_from_tile_data(tile_data, E)
				edges[S] = get_edge_from_tile_data(tile_data, S)
				edges[W] = get_edge_from_tile_data(tile_data, W)
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
// if we offered a N edge, need matching S edges (and vice versa)
// if we offered a E edge, need matching W edges (and vice versa)
func find_matches(tile Tile, edge string) map[string]string {
	matches := map[string]string{}
	for check_id, check_tile := range tiles {
		if tile.id == check_id {
			continue
		}
		if edge == N && tile.edges[N] == check_tile.edges[S] {
			matches[check_id] = S
			//fmt.Printf(" tile %v:N [%v] matches tile %v:S [%v]\n", tile.id, tile.edges[N], check_id, tiles[check_id].edges[S])
		}
		if edge == S && tile.edges[S] == check_tile.edges[N] {
			matches[check_id] = N
			//fmt.Printf(" tile %v:S [%v] matches tile %v:N [%v]\n", tile.id, tile.edges[S], check_id, tiles[check_id].edges[N])
		}
		if edge == E && tile.edges[E] == check_tile.edges[W] {
			matches[check_id] = W
			//fmt.Printf(" tile %v:E [%v] matches tile %v:W [%v]\n", tile.id, tile.edges[E], check_id, tiles[check_id].edges[W])
		}
		if edge == W && tile.edges[W] == check_tile.edges[E] {
			matches[check_id] = E
			//fmt.Printf(" tile %v:W [%v] matches tile %v:E [%v]\n", tile.id, tile.edges[W], check_id, tiles[check_id].edges[E])
		}
	}
	return matches
}

func find_all_matches(tile Tile, edge string) map[string]string {
	matches := map[string]string{}
	sides := []string{N, E, S, W}
	for check_id, check_tile := range tiles {
		if tile.id == check_id {
			continue
		}
		for _, side := range sides {
			if check_tile.edges[edge] == tile.edges[side] {
				matches[check_id] = side
				fmt.Printf(" tile %v:%v [%v] matches tile %v:%v [%v]\n", tile.id, side, tile.edges[side], check_id, edge, check_tile.edges[edge])
			}
		}
		// check the reversed version of each side
		for _, side := range sides {
			if reverse(check_tile.edges[edge]) == tile.edges[side] {
				matches[check_id] = side
				fmt.Printf(" tile %v:%v [%v] matches reversed-edge tile %v:%v [%v]\n", tile.id, side, tile.edges[side], check_id, edge, reverse(check_tile.edges[edge]))
			}
		}
	}
	return matches
}

func examine_tile(tile Tile) int {
	fmt.Println("Examine:", tile.id)
	matches_n := find_all_matches(tile, N)
	matches_e := find_all_matches(tile, E)
	matches_s := find_all_matches(tile, S)
	matches_w := find_all_matches(tile, W)
	if len(matches_n) > 0 {
		fmt.Printf("   %v matches on N\n", len(matches_n))
		for k, m := range matches_n {
			fmt.Printf("     %v:%v\n", k, m)
		}
	}
	if len(matches_e) > 0 {
		fmt.Printf("   %v matches on E\n", len(matches_e))
		for k, m := range matches_e {
			fmt.Printf("     %v:%v\n", k, m)
		}
	}
	if len(matches_s) > 0 {
		fmt.Printf("   %v matches on S\n", len(matches_s))
		for k, m := range matches_s {
			fmt.Printf("     %v:%v\n", k, m)
		}
	}
	if len(matches_w) > 0 {
		fmt.Printf("   %v matches on W\n", len(matches_w))
		for k, m := range matches_w {
			fmt.Printf("     %v:%v\n", k, m)
		}
	}
	return len(matches_n) + len(matches_e) + len(matches_s) + len(matches_w)
}

// Part 2 working...
// I know the corner tiles from part 1
// {"2273", "2243", "2953", "1213"}
// first find out which corners they are and in which rotations
// my functions tell me which the matching edges are so we can pencil in something
func run() {
	tile := tiles["2953"]
	tile.rotate()
	tile.rotate()
	tile.rotate()
	examine_tile(tile)
}

// NW and SW could be the other way round.
// I think the next step would be to find all tiles that have exactly 3 matching edges, these would be my edge tiles.
// And then those that have 4 matching edges, these are the inner tiles.
// Total pain to work out how they fit together though.
// Going to have to give up here!
var solution = [][]string{
	{"2243", "2707-FR1", "", "", "", "", "", "", "", "", "1543", "2273"},
	{"2663", "", "", "", "", "", "", "", "", "", "", "2617"},
	{"", "", "", "", "", "", "", "", "", "", "", ""},
	{"", "", "", "", "", "", "", "", "", "", "", ""},
	{"", "", "", "", "", "", "", "", "", "", "", ""},
	{"", "", "", "", "", "", "", "", "", "", "", ""},
	{"", "", "", "", "", "", "", "", "", "", "", ""},
	{"", "", "", "", "", "", "", "", "", "", "", ""},
	{"", "", "", "", "", "", "", "", "", "", "", ""},
	{"", "", "", "", "", "", "", "", "", "", "", ""},
	{"2371-R2", "", "", "", "", "", "", "", "", "", "", "3259-R3"},
	{"2953-R3", "", "", "", "", "", "", "", "", "", "3511-F", "1213"},
}

func part_1() (product int) {
	// examine the tiles ... how many edges line up with other tiles, for each?
	// known corner tile has 2 matching edges regardless of rotation
	// this works to produce exactly 4 corner tiles from both test and real data.
	corner_tiles := []string{}
	product = 1
	for _, tile := range tiles {
		num_matches := examine_tile(tile)
		if num_matches == 2 {
			fmt.Printf("Tile %v has %v matching edges\n", tile.id, num_matches)
			corner_tiles = append(corner_tiles, tile.id)
			id, _ := strconv.Atoi(tile.id)
			product *= id
		}
	}
	fmt.Println("The corner tiles are:", corner_tiles)
	// Part 1 - "multiply the IDs of the four corner tiles together"
	return
}

// Run the program with the argument "2" to run part 2, or anything else for part 1.
func main() {
	parse_input(get_input())
	//fmt.Printf("There are %v tiles\n", len(tiles))
	//part_1_answer := part_1()
	//fmt.Println("Part 1 answer:", part_1_answer)
	run()
}
