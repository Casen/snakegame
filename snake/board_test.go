package snake

import (
	"testing"

	"github.com/casen/snakegame/model"
)

var (
	eastVector  = model.Vector{X: 0, Y: 1}
	westVector  = model.Vector{X: 0, Y: -1}
	northVector = model.Vector{X: -1, Y: 0}
	southVector = model.Vector{X: 1, Y: 0}
)

func TestNewGameBoard(t *testing.T) {
	board := NewGameBoard(10, 10)
	if board == nil {
		t.Errorf("NewGameBoard() = %v; want %v", board, "not nil")
	}

	if board.snake == nil {
		t.Errorf("NewGameBoard() = %v; want %v", board.snake, "not nil")
	}

	if board.gameOver {
		t.Errorf("NewGameBoard() = %v; want %v", board.gameOver, "false")
	}
}

func TestCurrentState(t *testing.T) {
	board := NewBoard(
		10,
		10,
		NewSnake([]model.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 3}}, eastVector),
		model.Point{X: 0, Y: 4},
	)

	got := board.CurrentState()
	want := [11]float32{0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 0}

	if got != want {
		t.Errorf("CurrentState() = %v; want %v", got, want)
	}

}

func TestNextState(t *testing.T) {
	board := NewBoard(
		10,
		10,
		NewSnake([]model.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 3}}, eastVector),
		model.Point{X: 0, Y: 4},
	)

	gotCurrent := board.CurrentState()
	wantCurrent := [11]float32{0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 0}

	if gotCurrent != wantCurrent {
		t.Errorf("CurrentState() = %v; want %v", gotCurrent, wantCurrent)
	}

	// Generate a next state where the snake eats the food, and earns 1 point
	gotNext := board.NextState(eastVector)

	if gotNext[2] != 1 {
		t.Errorf("Expected danger on left. NextState() = %v; want [0,0,1 ...]", gotNext)
	}

	if gotNext[4] != 1 {
		t.Errorf("Expected snake moving right (east). NextState() = %v; want [0,0,1,0,1 ...]", gotNext)
	}

	if board.CurrentState() != gotCurrent {
		t.Errorf("CurrentState mutated by NextState: %v, %v", board.CurrentState(), gotNext)
	}

	if board.points != 0 {
		t.Errorf("Board points mutated by NextState: %v, expected 0", board.points)
	}

}

func TestEvaluateAction(t *testing.T) {
	type evaluateActionCase struct {
		testName     string
		board        *Board
		actionVector model.Vector
		wantReward   float32
		wantIsDone   bool
	}

	var testBoards []*Board = []*Board{
		NewBoard(
			10,
			10,
			NewSnake([]model.Point{{X: 1, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: 2}, {X: 1, Y: 3}}, eastVector),
			model.Point{X: 1, Y: 4},
		),
		NewBoard(
			10,
			10,
			NewSnake([]model.Point{{X: 3, Y: 0}, {X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}}, northVector),
			model.Point{X: 1, Y: 4},
		),
	}

	testCases := []evaluateActionCase{
		{
			"Snake moving east, action east, food one step to the east",
			testBoards[0], eastVector, 100, false,
		},
		{
			"Snake moving east, action north, food one step to the east",
			testBoards[0], northVector, -1, false,
		},
		{
			"Snake moving east, action south, food one step to the east",
			testBoards[0], southVector, -1, false,
		},
		{
			"Snake moving north, action east, wall one step to the north, and one step to the west",
			testBoards[1], eastVector, -1, false,
		},
		{
			"Snake moving north, action north, wall one step to the north, and one step to the west",
			testBoards[1], northVector, -100, true,
		},
		{
			"Snake moving north, action west, wall one step to the north, and one step to the west",
			testBoards[1], westVector, -100, true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			gotReward, gotIsDone := tc.board.EvaluateAction(tc.actionVector)
			if tc.wantReward != gotReward {
				t.Errorf("Expected '%f', but got '%f'", tc.wantReward, gotReward)
			}

			if tc.wantIsDone != gotIsDone {
				t.Errorf("Expected '%t', but got '%t'", tc.wantIsDone, gotIsDone)
			}
		})
	}
}

func TestValidMoves(t *testing.T) {
	type validMoveTestCase struct {
		name         string
		board        *Board
		actionVector model.Vector
		want         bool
	}
	board := NewBoard(
		10,
		10,
		NewSnake([]model.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 3}}, eastVector),
		model.Point{X: 0, Y: 4},
	)

	var testCases []validMoveTestCase = []validMoveTestCase{
		{"Move east", board, eastVector, true},
		{"Move west", board, westVector, false},
		{"Move north", board, northVector, true},
		{"Move south", board, southVector, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			nextLocation := tc.board.NextLocation(tc.actionVector)
			got := tc.board.MoveIsValid(nextLocation)
			if tc.want != got {
				t.Errorf("Expected '%t', but got '%t'", tc.want, got)
			}
		})
	}

}

func TestMove(t *testing.T) {
	board := NewBoard(
		20,
		20,
		NewSnake([]model.Point{{X: 8, Y: 13}, {X: 9, Y: 13}, {X: 10, Y: 13}, {X: 10, Y: 14}, {X: 11, Y: 14}}, southVector),
		model.Point{X: 0, Y: 4},
	)
	snakeHead := board.snake.Head()

	// If we move north, the snake should not move, since north is an invalid move
	board.Move(northVector)
	snakeHead = board.snake.Head()
	if snakeHead.X != 11 || snakeHead.Y != 14 {
		t.Errorf("Expected snake head to be at (11, 14), but got %v", snakeHead)
	}

	// If we move west, the snake head should move west
	board.Move(westVector)
	snakeHead = board.snake.Head()
	if snakeHead.X != 11 || snakeHead.Y != 13 {
		t.Errorf("Expected snake head to be at (11, 13), but got %v", snakeHead)
	}

}
