package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Player struct {
	ID    int
	Dice  []int
	Score int
}

type Move struct {
	FromPlayer int
	FromIndex  int
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var playersCount, diceCount int
	fmt.Print("Masukkan jumlah pemain: ")
	fmt.Scan(&playersCount)

	fmt.Print("Masukkan jumlah dadu per pemain: ")
	fmt.Scan(&diceCount)

	players := initializePlayers(playersCount, diceCount)

	round := 1
	for {
		fmt.Printf("\n==================\n")
		fmt.Printf("Giliran %d lempar dadu:\n", round)

		for i := range players {
			fmt.Printf("Pemain #%d (%d): %v\n", players[i].ID, players[i].Score, players[i].Dice)
		}

		processRound(players)

		if !checkIfGameContinues(players) {
			break
		}

		round++
	}

	fmt.Printf("\n==================\n")
	fmt.Printf("Game berakhir karena hanya pemain yang memiliki dadu.\n")
	winner := determineWinner(players)
	fmt.Printf("Game dimenangkan oleh pemain #%d karena memiliki poin lebih banyak dari pemain lainnya.\n", winner.ID)
}

func initializePlayers(playerCount, diceCount int) []Player {
	players := make([]Player, playerCount)
	for i := range players {
		players[i] = Player{
			ID:    i + 1,
			Dice:  make([]int, diceCount),
			Score: 0,
		}
		// Inisialisasi dadu
		rollDice(&players[i])
	}
	return players
}

func rollDice(player *Player) {
	player.Dice = make([]int, len(player.Dice))
	for i := range player.Dice {
		player.Dice[i] = rand.Intn(6) + 1
	}
}

func processRound(players []Player) {
	for i := range players {
		for j := len(players[i].Dice) - 1; j >= 0; j-- {
			switch players[i].Dice[j] {
			case 6:
				players[i].Dice = removeDice(players[i].Dice, j)
				players[i].Score++
			case 1:
				moveDiceToAdjacentPlayer(i, j, players)
			}
		}
	}

	fmt.Printf("Setelah evaluasi:\n")
	for i := range players {
		fmt.Printf("Pemain #%d (%d): %v\n", players[i].ID, players[i].Score, players[i].Dice)
	}

	for i := range players {
		if len(players[i].Dice) > 0 {
			rollDice(&players[i])
		}
	}
}

func moveDiceToAdjacentPlayer(currentPlayerIndex, diceIndex int, players []Player) {
	currentPlayer := &players[currentPlayerIndex]
	nextPlayerIndex := (currentPlayerIndex + 1) % len(players)

	if diceIndex == len(currentPlayer.Dice)-1 && currentPlayer.Dice[diceIndex] == 1 {
		// Jika pemain terakhir dan dadu bermata 1, pindahkan ke pemain pertama
		if nextPlayerIndex == 0 {
			nextPlayerIndex = len(players) - 1
		} else {
			nextPlayerIndex--
		}
	}

	nextPlayer := &players[nextPlayerIndex]

	if len(currentPlayer.Dice) > 0 && diceIndex < len(currentPlayer.Dice) {
		diceToMove := currentPlayer.Dice[diceIndex]

		// Periksa apakah pemain terkait memiliki dadu selain mata 1
		hasOtherDice := false
		for _, value := range currentPlayer.Dice {
			if value != 1 {
				hasOtherDice = true
				break
			}
		}

		// Hanya pindahkan dadu jika pemain memiliki dadu selain mata 1
		if hasOtherDice {
			currentPlayer.Dice = removeDice(currentPlayer.Dice, diceIndex)

			// Perbarui nextPlayerIndex jika dadu bermata 1 dan pemain terakhir
			if nextPlayerIndex == 0 && diceToMove == 1 {
				nextPlayerIndex = len(players) - 1
			}

			nextPlayer.Dice = append(nextPlayer.Dice, diceToMove)
		}
	}
}

func removeDice(dice []int, index int) []int {
	if len(dice) > 0 && index >= 0 && index < len(dice) {
		return append(dice[:index], dice[index+1:]...)
	}
	return dice
}

func checkIfGameContinues(players []Player) bool {
	activePlayers := 0
	for i := range players {
		if len(players[i].Dice) > 0 {
			activePlayers++
		}
	}
	return activePlayers > 1
}

func determineWinner(players []Player) Player {
	sort.Slice(players, func(i, j int) bool {
		return players[i].Score > players[j].Score
	})

	return players[0]
}
