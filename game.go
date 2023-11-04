package main

import "math/rand"

func initBoard(size int) [][]int {

	state := emptyBoard(size)

	for i := range state {
		for j := range state[i] {
			if rand.Float64() < density {
				state[i][j] = 1
			}
		}
	}

	return state
}

func emptyBoard(size int) [][]int {

	state := make([][]int, size)
	for i := range state {
		state[i] = make([]int, size)
	}
	return state
}

func countNeighbors(state [][]int, row int, col int) int {
	count := 0
	for i := row - 1; i <= row+1; i++ {
		for j := col - 1; j <= col+1; j++ {
			if i >= 0 && i < len(state) && j >= 0 && j < len(state) {
				if i == row && j == col {
					continue
				}
				count += state[i][j]
			}
		}
	}
	return count
}

func runCycle(game Game) Game {
	state := game.Board
	// Create a new 2D slice with the same dimensions as state.
	newState := make([][]int, len(state))
	for i := range state {
		newState[i] = make([]int, len(state[i]))
		copy(newState[i], state[i])
	}

	for i := 0; i < game.Dimension; i++ {
		for j := 0; j < game.Dimension; j++ {
			alive := newState[i][j] == 1
			neighbors := countNeighbors(newState, i, j)
			if alive {
				if neighbors < 2 || neighbors > 3 {
					newState[i][j] = 0
				}
			} else {
				if neighbors == 3 {
					newState[i][j] = 1
				}
			}
		}
	}

	dead := areEqual(state, newState)

	return Game{Id: game.Id, Board: newState, Dimension: game.Dimension, Dead: dead}
}

func areEqual(a, b [][]int) bool {
	for i := range a {
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}

	return true
}
