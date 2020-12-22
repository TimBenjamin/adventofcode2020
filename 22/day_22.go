package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func get_log(p1 []int, p2 []int) (log string) {
	tmp := []string{}
	for _, c := range p1 {
		tmp = append(tmp, strconv.Itoa(c))
	}
	log += strings.Join(tmp, ",")
	log += ":"
	for _, c := range p2 {
		tmp = append(tmp, strconv.Itoa(c))
	}
	log += strings.Join(tmp, ",")
	return
}

func log_round(previous_turns []string, player_1 []int, player_2 []int) []string {
	previous_turns = append(previous_turns, get_log(player_1, player_2))
	return previous_turns
}

// returns true if a previous round had the same cards, same order as the current decks
func check_previous_rounds(previous_turns []string, player_1 []int, player_2 []int) bool {
	if len(previous_turns) == 0 {
		return false
	}
	test_log := get_log(player_1, player_2)
	for _, log := range previous_turns {
		if log == test_log {
			fmt.Printf("previous round condition matched, test: %v against log: %v\n", test_log, log)
			return true
		}
	}
	return false
}

// can also use this in part 2 to get a unique number to save the round
func get_score(cards []int) (score int) {
	/*
		The bottom card in their deck is worth the value of the card multiplied by 1,
		the second-from-the-bottom card is worth the value of the card multiplied by 2, and so on.
		With 10 cards, the top card is worth the value on the card multiplied by 10.
		// test result: 306
	*/
	multiplier := 1
	for i := len(cards) - 1; i >= 0; i-- {
		score += cards[i] * multiplier
		multiplier++
	}
	return
}

func do_winner(player_1 []int, player_2 []int, p1_card int, p2_card int, winner int) ([]int, []int) {
	if winner == 1 {
		player_1 = append(player_1, p1_card)
		player_1 = append(player_1, p2_card)
		player_1 = player_1[1:]
		if len(player_2) > 0 {
			player_2 = player_2[1:]
		}
	} else {
		player_2 = append(player_2, p2_card)
		player_2 = append(player_2, p1_card)
		player_2 = player_2[1:]
		if len(player_1) > 0 {
			player_1 = player_1[1:]
		}
	}
	fmt.Printf("Player %v wins the turn\n", winner)
	return player_1, player_2
}

func take_turn(player_1 []int, player_2 []int, part_2 bool) ([]int, []int) {
	p1_card := player_1[0]
	p2_card := player_2[0]
	if part_2 {
		/*
			If both players have at least as many cards remaining in their deck as the value of the card
			they just drew, the winner of the round is determined by playing a new game of Recursive Combat
		*/
		if len(player_1) > p1_card && len(player_2) > p2_card {
			fmt.Printf("Subgame! remaining cards condition - start a new game (p1:%v, p2:%v)\n", p1_card, p2_card)
			// play a new game
			/*
				To play a sub-game of Recursive Combat, each player creates a new deck by making a copy of
				the next cards in their deck
				(the quantity of cards copied is equal to the number on the card they drew to trigger the sub-game).
			*/
			new_player_1 := []int{}
			for i := 1; i < len(player_1); i++ {
				new_player_1 = append(new_player_1, player_1[i])
			}
			new_player_2 := []int{}
			for i := 1; i < len(player_2); i++ {
				new_player_2 = append(new_player_2, player_2[i])
			}
			current_game++
			total_games++
			fmt.Println("Start game", current_game)
			subgame_score, winner := run(new_player_1, new_player_2, part_2)
			fmt.Printf("Player %v won the subgame, the score for game %v was %v\n", winner, current_game, subgame_score)
			current_game--
			fmt.Println("Anyway, back to game", current_game)
			// this round is determined by the winner of that subgame we just had
			player_1, player_2 = do_winner(player_1, player_2, p1_card, p2_card, winner)
			fmt.Println()
			return player_1, player_2
		}
	}
	fmt.Println(" p1 plays:", p1_card)
	fmt.Println(" p2 plays:", p2_card)
	if p1_card > p2_card {
		player_1, player_2 = do_winner(player_1, player_2, p1_card, p2_card, 1)
	} else if p2_card > p1_card {
		player_1, player_2 = do_winner(player_1, player_2, p1_card, p2_card, 2)
	} else {
		panic(errors.New("Draw! not possible?"))
	}
	fmt.Println()
	return player_1, player_2
}

var current_game = 1
var total_turns = 0
var total_games = 1

func run(player_1 []int, player_2 []int, part_2 bool) (score int, winner int) {
	var turns = 0
	var previous_turns = []string{}
	for {
		turns++
		total_turns++
		fmt.Printf("--- Game %v, turn %v---\n", current_game, turns)
		fmt.Println(" p1 deck:", player_1)
		fmt.Println(" p2 deck:", player_2)
		if part_2 {
			/*
				Before either player deals a card, if there was a previous round in this game that had
				exactly the same cards in the same order in the same players' decks, the game instantly
				ends in a win for player 1.
			*/
			if check_previous_rounds(previous_turns, player_1, player_2) {
				fmt.Println("Previous round condition - win for player 1")
				score = get_score(player_1)
				return score, 1
			} else {
				// log this round
				previous_turns = log_round(previous_turns, player_1, player_2)
			}
		}
		player_1, player_2 = take_turn(player_1, player_2, part_2)
		if len(player_1) == 0 {
			fmt.Println("Player 2 was the winner")
			score = get_score(player_2)
			return score, 2
		} else if len(player_2) == 0 {
			fmt.Println("Player 1 was the winner")
			score = get_score(player_1)
			return score, 1
		}
	}
}

// Run the program with the argument "2" to run part 2, or anything else for part 1.
func main() {
	var player_1, player_2 []int

	// test data
	player_1 = []int{9, 2, 6, 3, 1}
	player_2 = []int{5, 8, 4, 7, 10}

	// part 2 test data
	//player_1 = []int{43, 19}
	//player_2 = []int{2, 29, 14}

	// real data
	//player_1 = []int{44, 47, 29, 31, 10, 40, 50, 27, 35, 30, 38, 11, 14, 9, 42, 1, 26, 24, 6, 13, 8, 15, 21, 18, 4}
	//player_2 = []int{17, 22, 28, 34, 32, 23, 3, 19, 36, 12, 45, 37, 46, 39, 49, 43, 25, 33, 2, 41, 48, 7, 5, 16, 20}

	if len(os.Args) > 1 && os.Args[1] == "2" {
		score, winner := run(player_1, player_2, true)
		fmt.Printf("The winner was player %v and the score was %v\n", winner, score)
	} else {
		score, winner := run(player_1, player_2, false)
		fmt.Printf("The winner was player %v and the score was %v\n", winner, score)
	}
	fmt.Printf("Games: %v / turns: %v\n", total_games, total_turns)
}
