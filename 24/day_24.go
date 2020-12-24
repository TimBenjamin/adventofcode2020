package main

import (
	"bufio"
	"fmt"
	"os"
)

var test_input_file = "./test_input.txt"
var real_input_file = "./input.txt"
var input_file = real_input_file
var input []string

// https://www.redblobgames.com/grids/hexagons/#coordinates
// I'll try the x,y,z version
type Move struct {
	x int
	y int
	z int
}

type Tile struct {
	x       int
	y       int
	z       int
	flipped bool
}

// ne -> x: +1, y: 0, z: -1
// nw -> x: 0, y: +1, z: -1
// e -> x: +1, y: -1, z: 0
// w -> x: -1, y: +1, z: 0
// se -> x: 0, y: -1, z: +1
// sw -> x: -1, y: 0, z: +1
var ne = Move{1, 0, -1}
var nw = Move{0, 1, -1}
var e = Move{1, -1, 0}
var w = Move{-1, 1, 0}
var se = Move{0, -1, 1}
var sw = Move{-1, 0, 1}

func (t *Tile) applyMove(move Move) {
	t.x += move.x
	t.y += move.y
	t.z += move.z
}

func (t *Tile) flip() {
	t.flipped = !t.flipped
}

func (t *Tile) matches(tile Tile) bool {
	if t.x == tile.x && t.y == tile.y && t.z == tile.z {
		return true
	}
	return false
}

// get a slice of tiles that would have adjacent coords to this tile
func (t *Tile) getAdjacent() []Tile {
	adjacent := []Tile{}
	moves := []Move{ne, nw, e, w, se, sw}
	var c Tile
	for _, m := range moves {
		c = Tile{t.x, t.y, t.z, false}
		c.applyMove(m)
		adjacent = append(adjacent, c)
	}
	return adjacent
}

func get_input() (sum int) {
	f, err := os.Open(input_file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		input = append(input, line)
	}

	return
}

var moves [][]Move

// parse each line of the input into a slice of Moves
// and append that slice onto "moves"
func parse_input() {
	for _, line := range input {
		set := []Move{}
		for i := 0; i < len(line); i++ {
			if line[i] == 'e' {
				set = append(set, e)
			} else if line[i] == 'w' {
				set = append(set, w)
			} else {
				// we assume the line syntax is OK, so i+1 should not run over the length!
				switch line[i : i+2] {
				case "ne":
					set = append(set, ne)
				case "nw":
					set = append(set, nw)
				case "se":
					set = append(set, se)
				case "sw":
					set = append(set, sw)
				}
				i++ // it was 2 characters
			}
		}
		moves = append(moves, set)
	}
}

func run_first() []Tile {
	// for each set of moves...
	tiles := []Tile{}
	for _, set := range moves {
		tile := Tile{0, 0, 0, false}
		for _, move := range set {
			tile.applyMove(move)
		}
		// is it the same as one we've already done?
		same := false
		for i := 0; i < len(tiles); i++ {
			if tiles[i].matches(tile) {
				tiles[i].flip()
				same = true
				break
			}
		}
		if !same {
			tile.flip()
			tiles = append(tiles, tile)
		}
	}
	fmt.Println("There are how many tiles to begin with:", len(tiles))
	fmt.Println("Day 0:", get_num_black(tiles))
	return tiles
}

func tile_exists(t Tile, tiles []Tile) bool {
	for _, n := range tiles {
		if t.matches(n) {
			return true
		}
	}
	return false
}

func part_2() {
	loops := 100 // 100
	day := 1     // day 1 is run_first()
	tiles := run_first()
	fmt.Println()
	for {
		// Every day, the tiles are all flipped according to the following rules:
		// - Any black tile with zero or more than 2 black tiles immediately adjacent to it is flipped to white.
		// - Any white tile with exactly 2 black tiles immediately adjacent to it is flipped to black.
		// Here, tiles immediately adjacent means the six tiles directly touching the tile in question.

		// part 1 just gave us some black tiles to start with, in an infinite sea of white tiles
		// so the second rule applies to the adjacent tiles of each black tile;
		// - if that adjacent tile is white (i.e. it does not exist in our list of black tiles)
		// - if it flips, it becomes black so we can add it to our master list of affected tiles

		// This all takes place at once so store some tiles to flip all at once:
		tiles_to_flip := []int{}

		// Maybe this would be easier if first I go through all the tiles I know about
		// and add adjacent white tiles to the existing black tiles, if they don't already exist
		// then I am guaranteed that each black tile is surrounded by either black or white tiles
		// surely doesn't matter if there are gaps, I only really care about black tiles and their white neighbours

		new_tiles := []Tile{} // just in case it adds them to the thing I'm already ranging over
		for _, t := range tiles {
			if !t.flipped {
				continue
			}
			adjacent := t.getAdjacent()
			for _, a := range adjacent {
				if !tile_exists(a, tiles) && !tile_exists(a, new_tiles) {
					new_tiles = append(new_tiles, a)
				}
			}
		}
		for _, n := range new_tiles {
			tiles = append(tiles, n)
		}

		// now go through all the tiles and apply the rules
		for i := 0; i < len(tiles); i++ {
			// NB getAdjacent just returns tiles that would theoretically be adjacent to the one we're looking at
			// they are fresh tiles though! Have to actually find them in "tiles" if we need to use them
			adjacent := tiles[i].getAdjacent()
			black_adjacent := 0
			for _, a := range adjacent {
				// so is the adjacent tile one that we know about?
				for j := 0; j < len(tiles); j++ {
					if tiles[j].matches(a) {
						if tiles[j].flipped {
							black_adjacent++
						}
					}
				}
			}
			if tiles[i].flipped {
				// we are on a black tile...
				if black_adjacent == 0 || black_adjacent > 2 {
					tiles_to_flip = append(tiles_to_flip, i)
				}
			} else {
				// we are on a white tile...
				if black_adjacent == 2 {
					tiles_to_flip = append(tiles_to_flip, i)
				}
			}
		}

		// finally apply the flipping operations all at once
		// tiles_to_flip just has the indexes in the "tiles" array of tiles that should be flipped
		for _, k := range tiles_to_flip {
			tiles[k].flip()
		}

		fmt.Printf("Day %v: %v\n", day, get_num_black(tiles))
		if day == loops {
			break
		}
		day++
	}
}

func get_num_black(tiles []Tile) int {
	// how many tiles are left with the black side up?
	black := 0
	for _, tile := range tiles {
		if tile.flipped {
			black++
		}
	}
	return black
}

// Run the program with the argument "2" to run part 2, or anything else for part 1.
func main() {
	get_input()
	parse_input()
	if len(os.Args) > 1 && os.Args[1] == "2" {
		// PART 2
		part_2()
	} else {
		// PART 1
		tiles := run_first()
		fmt.Println("The answer is:", get_num_black(tiles))
	}
}
