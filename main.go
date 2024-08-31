package main

import (
	"fmt"
)

const boardSize int = 4
const winCount int = 3
const searchDepth int = 5

type boardState [boardSize][boardSize]int8

type vec2 struct {
	x int
	y int
}

type boardNode struct {
	board          boardState
	leaf           bool
	traversed      bool
	expanded       bool
	player         int8
	value          int
	parent         *boardNode
	valueListeners []*boardNode
	children       []*boardNode
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

// only works with 4x4 boards
func hashBoard(board boardState) uint32 {
	var buffer uint32
	buffer = 0
	for y := range 4 {
		for x := range 4 {
			buffer |= (uint32(board[x][y]+1) & 0b11) << ((y*4 + x) * 2)
		}
	}
	return buffer
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
				addMove(node, x, y, 1, 0)
				addMove(node, x, y, -1, 0)
				addMove(node, x, y, 0, 1)
				addMove(node, x, y, 0, -1)
			}
		}
	}
	node.expanded = true
}
func addMove(node *boardNode, x int, y int, xoffset int, yoffset int) {
	if pointInBoard(x+xoffset, y+yoffset) && node.board[x+xoffset][y+yoffset] == 0 {
		board := (*node).board
		board[x][y] = 0
		board[x+xoffset][y+yoffset] = node.player
		winner := checkWin(board, vec2{x: x + xoffset, y: y + yoffset})
		child := boardNode{
			parent:   node,
			board:    board,
			leaf:     winner == node.player,
			expanded: false,
			player:   -node.player,
			value:    int(winner),
		}
		if val, ok := visitedNodes[hashBoard(board)]; ok {
			parent := node.parent
			for parent != nil {
				if parent.board == board {
					fmt.Println("Repeated position in sequence. Breaking")
					return
				}
				parent = parent.parent
			}
			child.leaf = true
			val.valueListeners = append(val.valueListeners, &child)
			printBoard(board)
			fmt.Println("appended above as a value listener of")
			printBoard(val.board)
		} else {
			visitedNodes[hashBoard(board)] = &child
		}
		node.children = append(node.children, &child)

		updateListeners(int(winner), node)
	}
}
func updateListeners(value int, node *boardNode) {
	if value != 0 {
		listeners := node.parent.valueListeners
		for _, lis := range listeners {
			lis.value += value
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

var visitedNodes map[uint32]*boardNode

func generateBoardTree() {
	var node *boardNode

	node = &root

	viewedBoards := 0
	depth := 0
	for !root.traversed {
		viewedBoards++
		fmt.Println(viewedBoards)
		if depth > searchDepth {
			node.traversed = true
			if node.parent != nil {
				node = node.parent
				depth--
			}
			continue
		}
		if !node.expanded {
			expandBoard(node)
			node.expanded = true
		}
		allTraversed := true
		if node.children == nil {
			node.leaf = true
			node.traversed = true
			if node.parent != nil {
				node = node.parent
				depth--
			}
		} else {
			for i, c := range node.children {
				if !(c.leaf || c.traversed) {
					allTraversed = false
					node = node.children[i]
					depth++
					break
				}
			}
			if allTraversed {
				node.traversed = true
				if node.parent != nil {
					node = node.parent
					depth--
				}
			}
		}
	}
}

func backpropagateTree() {
	for _, node := range visitedNodes {
		if node.value != 0 {
			parent := node.parent
			value := -node.value
			for parent != nil {
				parent.value += value
				parent = parent.parent
				value = -value
			}
		}
	}
}

func main() {
	fmt.Println("Hello world. Toe solver v0")
	visitedNodes = make(map[uint32]*boardNode)
	generateBoardTree()
	fmt.Println("Board Generated. Backpropagating...")
	backpropagateTree()
	printBoard(root.board)
}
