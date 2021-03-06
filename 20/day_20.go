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

var corner_tiles_real = []string{"2273", "2243", "2953", "1213"}
var corner_tiles_test = []string{"1951", "2971", "3079", "1171"}
var corner_tiles = corner_tiles_real

// I will map the tiles by their file ID number
// I need these to be manipulated by reference so use a pointer
var tiles = map[string]*Tile{}

// describe tiles by their flipped status, how many times rotated, and what their n, e, s, w edges are like (which will vary depending on rotations and flip)
type Tile struct {
	id        string
	edges     map[string]string
	flipped   bool
	rotations int      // clockwise
	data      []string // the full contents of the tile
}

func (t *Tile) rotate() {
	if t.rotations < 3 {
		t.rotations++
	} else {
		t.rotations = 0
	}
	w := t.edges[W]
	t.edges[W] = t.edges[S]
	t.edges[S] = reverse(t.edges[E]) // I'm counting top to bottom, left to right, so these get reversed!
	t.edges[E] = t.edges[N]
	t.edges[N] = reverse(w) // I'm counting top to bottom, left to right, so these get reversed!

	// rotate the data too...
	t.data = rotate_data(t.data)
}

func rotate_data(data []string) []string {
	rotated_data := []string{}
	for c := 0; c < len(data); c++ {
		col := get_col(c, data)
		rotated_data = append(rotated_data, reverse(col))
	}
	return rotated_data
}

func get_col(col int, data []string) string {
	col_slice := []string{}
	for _, row := range data {
		if col < len(row) {
			col_slice = append(col_slice, string(row[col]))
		}
	}
	col_string := strings.Join(col_slice, "")
	return col_string
}

// this flips on the vertical axis!
// a horizontal axis flip is the same as vertical flip + one rotate
func (t *Tile) flip() {
	t.flipped = !t.flipped
	w := t.edges[W]
	t.edges[W] = t.edges[E]
	t.edges[E] = w

	// reverse n and s
	t.edges[N] = reverse(t.edges[N])
	t.edges[S] = reverse(t.edges[S])

	// flip the data too...
	t.data = flip_data(t.data)
}

func flip_data(data []string) []string {
	flipped_data := []string{}
	for _, r := range data {
		flipped_data = append(flipped_data, reverse(r))
	}
	return flipped_data
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

func (t *Tile) print() {
	fmt.Println("Tile ID:", t.id)
	for s, p := range t.edges {
		fmt.Printf(" %v: %v\n", s, p)
	}
	fmt.Println(" flipped:", t.flipped)
	fmt.Println(" rotations:", t.rotations)
	fmt.Println(" data:")
	for _, r := range t.data {
		fmt.Println("   ", r)
	}
	fmt.Println("---------------")
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

// combine the edge elements into a string for comparison
func get_edge_from_tile_data(tile []string, edge string) string {
	var idx int
	if edge == N {
		return tile[0]
	} else if edge == S {
		return tile[len(tile)-1]
	} else if edge == W {
		idx = 0
	} else if edge == E {
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

// read the loaded input into a handy data structure
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
				tiles[tile_id] = &Tile{
					tile_id, edges, false, 0, tile_data,
				}
			}

		}
	}
}

// helper to check if a string is in a slice of strings
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

// used in part 1 only
/*
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
*/

// used for part 1 only
/*
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
*/

func find_matches_on_edge(tile *Tile, edge string, ignore []string) (*Tile, bool) {
	var target_edge string
	for check_id, check_tile := range tiles {
		if tile.id == check_id {
			// don't compare our tile with itself
			continue
		}
		if in_array(check_id, ignore) {
			// don't consider this one, it is already fixed elsewhere
			continue
		}
		// flip and rotate check_tile until we get a match (8 possible matches from each tile)
		// if no matches, reset it and next
		num_edges_to_match := 1 // might be more than one edge if we are filling a grid
		num_matches := 0
		// our target is the opposite edge to the one given
		switch edge {
		case N:
			target_edge = S
		case S:
			target_edge = N
		case E:
			target_edge = W
		case W:
			target_edge = E
		}
		for i := 0; i < 4; i++ {
			if check_tile.edges[target_edge] == tile.edges[edge] {
				//fmt.Printf("  tile %v:%v (f: %v, r: %v) matches our tile %v:%v - %v\n", check_tile.id, target_edge, f, r, tile.id, edge, check_tile.edges[target_edge])
				num_matches++
				if num_matches == num_edges_to_match {
					return check_tile, true
				}
			}
			check_tile.rotate()
		}
		check_tile.flip()
		for i := 0; i < 4; i++ {
			if check_tile.edges[target_edge] == tile.edges[edge] {
				//fmt.Printf("  tile %v:%v (f: %v, r: %v) matches our tile %v:%v - %v\n", check_tile.id, target_edge, f, r, tile.id, edge, check_tile.edges[target_edge])
				num_matches++
				if num_matches == num_edges_to_match {
					return check_tile, true
				}
			}
			check_tile.rotate()
		}
	}
	return nil, false
}

func find_matches_east(tile_W *Tile, tile_N *Tile, ignore []string) (*Tile, bool) {
	for check_id, check_tile := range tiles {
		if tile_W.id == check_id || (tile_N != nil && tile_N.id == check_id) {
			// don't compare our tile with itself
			continue
		}
		if in_array(check_id, ignore) {
			// don't consider this one, it is already fixed elsewhere
			continue
		}
		// flip and rotate check_tile until we get a match (8 possible matches from each tile)
		// if no matches, reset it and next
		for flip := 0; flip < 2; flip++ {
			for i := 0; i < 4; i++ {
				if check_tile.edges[W] == tile_W.edges[E] {
					if tile_N == nil || check_tile.edges[N] == tile_N.edges[S] {
						return check_tile, true
					}
				}
				check_tile.rotate()
			}
			check_tile.flip()
		}
	}
	return nil, false
}

// a function to find matches from other tiles on this Tile's edge in the direction specified
// we might want to ignore certain tiles (as they are already fixed)
func get_column(start_tile *Tile, direction string, ignore []string) (column []*Tile) {
	column = append(column, start_tile)
	for col := 0; col < 12; col++ {
		tile := column[len(column)-1]
		// ignore everything already in this chain and in the ignore slice we were passed
		for _, t := range column {
			if !in_array(t.id, ignore) {
				ignore = append(ignore, t.id)
			}
		}
		matching_tile, found := find_matches_on_edge(tile, direction, ignore)
		if found {
			column = append(column, matching_tile)
		} else {
			// No more matches
			break
		}
	}
	return
}

// similar ...
// find tiles that match the E of the first element in the row, and the S of the corresponding element in the row above (if there is one)
func get_chain(grid [][]*Tile, row_num int, ignore []string) (chain []*Tile) {
	chain = append(chain, grid[row_num][0])
	for col := 1; col < 12; col++ {
		tile_W := chain[len(chain)-1]
		var tile_N *Tile = nil
		if row_num > 0 {
			tile_N = grid[row_num-1][col]
		}
		// ignore everything already in this chain and in the ignore slice we were passed
		for _, t := range chain {
			if !in_array(t.id, ignore) {
				ignore = append(ignore, t.id)
			}
		}
		matching_tile, found := find_matches_east(tile_W, tile_N, ignore)
		if found {
			chain = append(chain, matching_tile)
		} else {
			// No more matches
			break
		}
	}
	return
}

func show(tileset []*Tile) {
	for _, t := range tileset {
		fmt.Printf(" %v", t.id)
	}
	fmt.Println()
}

// Part 2 working...
// I know the corner tiles from part 1, these are the only tiles that have 2 matching edges (using find_all_matches())
// {"2273", "2243", "2953", "1213"}
// starting with one of the corners, try and fill in an edge of the big picture
func run() {
	// Through trial and error I know that this corner tile works without flipping / rotating
	corner_tile := "2273"

	// the grid we will put the fixed tiles in:
	grid := make([][]*Tile, 12)

	// the grid will start with a corner tile
	header := make([]*Tile, 12)
	header[0] = tiles[corner_tile]
	grid[0] = header

	// find a row running E beginning on specified start tile:
	ignore := []string{}
	header = get_chain(grid, 0, ignore)
	for _, t := range header {
		ignore = append(ignore, t.id)
	}
	for _, t := range tiles {
		if !in_array(t.id, ignore) {
			t.reset()
		}
	}
	grid[0] = header

	// OK that worked to find an edge row.
	// These tiles should be able to form the heads of columns.
	// So we should be able to find one chain for each, using the S edge, and end up with a 12x12 grid!

	// now get the left column
	start_tile := header[0]
	column := get_column(start_tile, S, ignore)

	for j := 1; j < len(grid); j++ {
		grid[j] = append(grid[j], column[j])
		ignore = append(ignore, column[j].id)
	}
	for _, t := range tiles {
		if !in_array(t.id, ignore) {
			t.reset()
		}
	}

	// Now get all the other rows based on the first column
	// these need to match both E of the tile to the left, and S of the tile above!
	for row_num := 1; row_num < 12; row_num++ {
		row := get_chain(grid, row_num, ignore)
		for j := 1; j < 12; j++ {
			grid[row_num] = append(grid[row_num], row[j])
			ignore = append(ignore, row[j].id)
		}
		for _, t := range tiles {
			if !in_array(t.id, ignore) {
				t.reset()
			}
		}
	}

	fmt.Println()
	visualise_grid(grid)
	fmt.Println()

	// Well at least it gives the grid...
	/*
		GRID:
		2273 1543 2143 3041 1759 1667 1621 2441 1867 1531 3259 1213
		2617 2593 3623 3203 2699 3023 2711 1979 1787 2099 3187 3511
		2969 2011 2297 3229 1571 3929 2357 1931 2591 2237 3329 1289
		1823 3221 3391 2579 1229 2749 3541 1259 2551 3853 1933 2677
		1523 3163 3407 1049 1993 3251 3767 1597 3631 1019 2207 2549
		2141 3373 1663 1087 3863 1181 3209 3643 1097 1801 3769 1307
		1697 1567 1907 1637 2287 3313 1879 3923 1063 1601 1471 1583
		1747 1303 3533 2903 2957 1733 1447 1657 1831 1873 2003 1489
		1039 2621 3571 2281 3539 3037 2843 3557 2609 3833 2161 3559
		2521 2467 3469 3011 1187 1009 2063 2753 2477 1399 3821 1481
		3299 2341 1913 1321 2909 2411 3467 3491 3343 1409 1549 2707
		2953 2371 1117 3877 2111 1607 2333 3881 3067 2087 2663 2243
		------------------------------------------------------------
	*/

	// Let's see the picture
	// we need to strip the edges from each tile.data first
	// the tile data is 10x10 so stripping the outside should make it 8x8
	for _, tile := range tiles {
		new_data := []string{}
		for i, d := range tile.data {
			if i == 0 || i == len(tile.data)-1 {
				continue
			}
			new_data = append(new_data, d[1:len(d)-1])
		}
		tile.data = new_data
	}

	picture := []string{}
	for _, row := range grid {
		for i := 0; i < len(tiles[corner_tile].data); i++ { // all tiles are the same, just picking any one here
			r := []string{} // row of the picture
			for _, tile := range row {
				r = append(r, tile.data[i])
			}
			pic_row := strings.Join(r, "")
			picture = append(picture, pic_row)
		}
	}

	// the grid is 12 x 12
	// so we should end up with a thing that is (12x8) by (12x8) i.e. 96x96
	fmt.Printf("Picture is %v by %v\n", len(picture[0]), len(picture))
	fmt.Printf("Grid is %v by %v\n", len(grid[0]), len(grid))

	// Have a look at the picture:
	for _, p := range picture {
		fmt.Println(p)
	}

	// OK now we are looking for this pattern in this big image:
	// Actually it seems that the # just have to be there. The "." can be "#" or "."
	/*
		..................#.
		#....##....##....###
		.#..#..#..#..#..#...
	*/

	// We might have to rotate or flip the damn image first.
	// the instructions say we might have to.
	// Blatantly we have to rotate / flip it.
	// fortunately we have some functions to more or less do that in the Tile...

	// testing - find_dragons() works
	/*
		picture = make([]string, 3)
		picture[0] = "..................#."
		picture[1] = "#....##....##....###"
		picture[2] = ".#..#..#..#..#..#..."
	*/

	// test 2:
	/*
		picture = []string{}
		picture = append(picture, ".#.#..#.##...#.##..#####")
		picture = append(picture, "###....#.#....#..#......")
		picture = append(picture, "##.##.###.#.#..######...")
		picture = append(picture, "###.#####...#.#####.#..#")
		picture = append(picture, "##.#....#.##.####...#.##")
		picture = append(picture, "...########.#....#####.#")
		picture = append(picture, "....#..#...##..#.#.###..")
		picture = append(picture, ".####...#..#.....#......")
		picture = append(picture, "#..#.##..#..###.#.##....")
		picture = append(picture, "#.####..#.####.#.#.###..")
		picture = append(picture, "###.#.#...#.######.#..##")
		picture = append(picture, "#.####....##..########.#")
		picture = append(picture, "##..##.#...#...#.#.#.#..")
		picture = append(picture, "...#..#..#.#.##..###.###")
		picture = append(picture, ".#.#....#.##.#...###.##.")
		picture = append(picture, "###.#...#..#.##.######..")
		picture = append(picture, ".#.#.###.##.##.#..#.##..")
		picture = append(picture, ".####.###.#...###.#..#.#")
		picture = append(picture, "..#.#..#..#.#.#.####.###")
		picture = append(picture, "#..####...#.#.#.###.###.")
		picture = append(picture, "#####..#####...###....##")
		picture = append(picture, "#.##..#..#...#..####...#")
		picture = append(picture, ".#.###..##..##..####.##.")
		picture = append(picture, "...###...##...#...#..###")
		picture = flip_data(picture)
		picture = rotate_data(picture)
		picture = rotate_data(picture)
		picture = rotate_data(picture)
	*/

	fmt.Printf("Picture is %v by %v\n", len(picture[0]), len(picture))
	num_dragons := find_dragons(picture)
	if num_dragons == 0 {
	outer:
		for flip := 0; flip < 2; flip++ {
			for rot := 0; rot < 4; rot++ {
				picture = rotate_data(picture)
				num_dragons = find_dragons(picture)
				if num_dragons > 0 {
					break outer
				}
			}
			picture = flip_data(picture)
			num_dragons = find_dragons(picture)
			if num_dragons > 0 {
				break outer
			}
		}
	}

	// this gives us 31 dragons... (this is 10 dragons short, either that or I am counting 150 hashes too many!)
	fmt.Println("Found how many dragons:", num_dragons)
	// there are 15 x "#" in each dragon
	// we need the total number of "#" in the picture, minus the number of # that are in dragons.
	// so let's count the total of # in the picture...
	num_hashes := 0
	for _, row := range picture {
		num_hashes += strings.Count(row, "#")
	}
	fmt.Println("Number of hashes:", num_hashes)
	fmt.Println("Hashes minus dragons, the final answer:", num_hashes-(num_dragons*15))
}

func find_dragons(picture []string) int {
	num_dragons := 0
	test := 0
	for i, row := range picture {
		if i < 2 {
			continue // no need to look in the top 2 rows, as we'll look for the "base" of the dragon
		}
		/*
			..................#.
			#....##....##....###
			.#..#..#..#..#..#...
		*/
		// the # must match, the bits in between can be either . or #
		// the pattern just has to fit in amongst the rest
		// pattern for each row:
		re_bottom := regexp.MustCompile(`.\#.{2}\#.{2}\#.{2}\#.{2}\#.{2}\#.{3}`)
		// There are overlapping dragons!! (found out the hard way)
		// Go doesn't do overlapping regexp matches
		// find the first match, then truncate the row to start 1 character after this match
		// and stop when there are no more matches
		tmp := row
		snip := 0
		for {
			match := re_bottom.FindStringIndex(tmp)
			if match != nil {
				test++
				tmp = tmp[match[0]+1:]
				m_1 := snip + match[0]
				m_2 := snip + match[0] + 20
				middle_segment := picture[i-1][m_1:m_2]
				snip += match[0] + 1
				match_2, _ := regexp.MatchString(`\#.{4}\#{2}.{4}\#{2}.{4}\#{3}`, middle_segment)
				if match_2 {
					top_segment := picture[i-2][m_1:m_2]
					match_1, _ := regexp.MatchString(`.{18}\#.`, top_segment)
					if match_1 {
						num_dragons++
					}
				}
			} else {
				break
			}
		}
	}
	return num_dragons
}

func visualise_grid(grid [][]*Tile) {
	fmt.Println("GRID:")
	for _, row := range grid {
		for _, tile := range row {
			fmt.Printf(" %v", tile.id)
		}
		fmt.Println()
	}
	fmt.Println("------------------------------------------------------------")
}

/*
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
*/

// Run the program with the argument "2" to run part 2, or anything else for part 1.
func main() {
	parse_input(get_input())
	//fmt.Printf("There are %v tiles\n", len(tiles))
	//part_1_answer := part_1()
	//fmt.Println("Part 1 answer:", part_1_answer)
	run()
}
