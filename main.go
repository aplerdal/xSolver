package main

import (
	"fmt"
)

const boardSize int = 4
const winCount int = 3

type boardState [boardSize][boardSize]int8

type vec2 struct {
	x int
	y int
}

type boardNode struct {
	board     boardState
	leaf      bool
	traversed bool
	player    int8
	value     int
	parent    *boardNode
	children  []*boardNode
}

// Checks for a winner given the move and the board state. Returns the winning player
func checkWin(board boardState, position vec2) int8 {
	player := board[position.x][position.y]
	// Check horizontal
	count := 0
	for piece := range boardSize {
		if board[piece][position.y] == player {
			count++
			if count >= winCount {
				return player
			}
		} else {
			count = 0
		}
	}
	// Check vertical
	count = 0
	for piece := range boardSize {
		if board[position.x][piece] == player {
			count++
			if count >= winCount {
				return player
			}
		} else {
			count = 0
		}
	}
	// Check Down Right Diagonal
	dr_start := vec2{x: position.x - min(position.x, position.y), y: position.y - min(position.x, position.y)}
	dr_limit := boardSize - max(dr_start.x, dr_start.y)
	count = 0
	for piece := range dr_limit {
		if board[dr_start.x+piece][dr_start.y+piece] == player {
			count++
			if count >= winCount {
				return player
			}
		} else {
			if piece >= winCount-1 {
				break
			}
			count = 0
		}
	}
	// Check Up Right Diagonal
	ur_start := vec2{x: position.x - min(position.x, (boardSize-1)-position.y), y: position.y + min(position.x, (boardSize-1)-position.y)}
	ur_limit := boardSize - max(ur_start.x, (boardSize-1)-ur_start.y)
	count = 0
	for piece := range ur_limit {
		if board[ur_start.x+piece][ur_start.y-piece] == player {
			count++
			if count >= winCount {
				return player
			}
		} else {
			if piece >= winCount-1 {
				break
			}
			count = 0
		}
	}

	return 0
}

func printBoard(board boardState) {
	for y := range boardSize {
		for x := range boardSize {
			fmt.Print(numToPiece(board[x][y]))
		}
		fmt.Print("\n")
	}
}
func numToPiece(num int8) string {
	if num < 0 {
		return "\u001B[44m\u001B[37;1m  \u001B[0m"
	} else if num > 0 {
		return "\u001B[41m\u001B[37;1m  \u001B[0m"
	} else {
		return "  "
	}
}
func expandBoard(node *boardNode) {
	for y := range boardSize {
		for x := range boardSize {
			if node.board[x][y] == node.player {
				checkMove(node, x, y, 1, 0)
				checkMove(node, x, y, -1, 0)
				checkMove(node, x, y, 0, 1)
				checkMove(node, x, y, 0, -1)
			}
		}
	}
}
func checkMove(node *boardNode, x int, y int, xoffset int, yoffset int) {
	if pointInBoard(x+xoffset, y+yoffset) && node.board[x+xoffset][y+yoffset] == 0 {
		board := (*node).board
		board[x][y] = 0
		board[x+xoffset][y+yoffset] = node.player
		value := checkWin(board, vec2{x: x + xoffset, y: y + yoffset})
		child := boardNode{
			parent: node,
			board:  board,
			leaf:   value == node.player,
			player: -node.player,
			value:  int(value),
		}
		node.children = append(node.children, &child)

		backpropagateValue(int(value), node)
	}
}
func backpropagateValue(value int, node *boardNode) {
	if value != 0 {
		parent := node.parent
		for parent != nil {
			value = -value
			parent.value += value
			parent = parent.parent
		}
	}
}
func pointInBoard(x int, y int) bool {
	if x >= 0 && y >= 0 && x < boardSize && y < boardSize {
		return true
	}
	return false
}

var root boardNode = boardNode{
	parent: nil,
	board: boardState{
		{1, 0, 0, -1},
		{-1, 0, 0, 1},
		{1, 0, 0, -1},
		{-1, 0, 0, 1},
	},
	player: 1,
	value:  0,
	leaf:   false,
}

func main() {
	fmt.Println("Hello world. Toe solver v0")
	var node *boardNode
	node = &root
	expandBoard(node)

	for root.traversed == false {

	}

	printBoard(node.board)
}
