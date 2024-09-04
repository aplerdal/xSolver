package main

import (
	"fmt"
	"slices"
	"time"
)

const boardSize int = 4
const winCount int = 3
const searchDepth int = 25

type boardState [boardSize][boardSize]int8

type vec2 struct {
	x int
	y int
}
type boardNode struct {
	board     boardState
	leaf      bool
	traversed bool
	expanded  bool
	player    int8
	parent    *boardNode
	children  []*boardNode
}
type boardInfo struct {
	board     boardState
	winner    int8
	player    int8
	leaf      bool
	generated bool
	children  []boardState
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
func expandBoard(b *boardInfo) {
	for y := range boardSize {
		for x := range boardSize {
			if b.board[x][y] == b.player {
				addMove(b, x, y, 1, 0)
				addMove(b, x, y, -1, 0)
				addMove(b, x, y, 0, 1)
				addMove(b, x, y, 0, -1)
			}
		}
	}
}
func addMove(b *boardInfo, x int, y int, xoffset int, yoffset int) {
	if pointInBoard(x+xoffset, y+yoffset) && b.board[x+xoffset][y+yoffset] == 0 {
		board := b.board
		board[x][y] = 0
		board[x+xoffset][y+yoffset] = b.player
		winner := checkWin(board, vec2{x: x + xoffset, y: y + yoffset})
		child := boardInfo{
			board:  board,
			winner: winner,
			player: -b.player,
		}
		if _, ok := boardMap[hashBoard(board)]; !ok {
			boardMap[hashBoard(board)] = child
		}
		b.children = append(b.children, board)
	}
}
func pointInBoard(x int, y int) bool {
	if x >= 0 && y >= 0 && x < boardSize && y < boardSize {
		return true
	}
	return false
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

var boardMap map[uint32]boardInfo
var root boardInfo = boardInfo{
	board: boardState{
		{1, 0, 0, -1},
		{-1, 0, 0, 1},
		{1, 0, 0, -1},
		{-1, 0, 0, 1},
	},
	player: 1,
	winner: 0,
}

func expandTreeNode(b boardInfo) {
	expandBoard(&b)
	b.generated = true
	boardMap[hashBoard(b.board)] = b
	for _, child := range b.children {
		if !boardMap[hashBoard(child)].generated {
			expandTreeNode(boardMap[hashBoard(child)])
		}
	}
}

func generateBoardTree() {
	boardMap = make(map[uint32]boardInfo)
	boardMap[hashBoard(root.board)] = root
	expandTreeNode(boardMap[hashBoard(root.board)])
}

func solveBoardTree() {
	DetermineState(hashBoard(root.board))
}

var parentStack []uint32

func DetermineState(hash uint32) {
	b := boardMap[hash]
	if slices.Contains(parentStack, hash) {
		b := boardMap[hash]
		b.winner = -2
		boardMap[hash] = b
	} else {
		parentStack = append(parentStack, hash)
		var matchValue int8
		if len(b.children) > 0 {
			matchValue = boardMap[hashBoard(b.children[0])].winner
		} else {
			matchValue = 0
		}
		matching := true
		winningMove := false
		if matchValue == 0 {
			matching = false
		}
		for _, c := range b.children {
			bi := boardMap[hashBoard(c)]
			if bi.winner == 0 {
				DetermineState(hashBoard(c))
				bi = boardMap[hashBoard(c)]
			}
			if bi.winner != matchValue {
				matching = false
			}
			if bi.winner == b.player {
				winningMove = true
				break
			}
		}
		if matching {
			b.winner = matchValue
		} else if winningMove {
			b.winner = b.player
		} else {
			b.winner = -2
		}
		boardMap[hash] = b
	}
	parentStack = parentStack[:len(parentStack)-1]
}

func main() {
	fmt.Println("Starting Tree Generation...")
	gstart := time.Now()
	generateBoardTree()
	fmt.Println("Generated Board Tree in", time.Since(gstart))
	fmt.Println("Boards generated:", len(boardMap))
	fmt.Println("Solving game...")
	sstart := time.Now()
	solveBoardTree()
	fmt.Println("Solved Board Tree in", time.Since(sstart))
	fmt.Println("The root node is", boardMap[hashBoard(root.board)].winner)
}
