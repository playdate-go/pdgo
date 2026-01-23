package main

import "github.com/playdate-go/pdgo"

// Conway's Game of Life rules:
// - Any live cell with fewer than two live neighbours dies (under-population)
// - Any live cell with two or three live neighbours lives on
// - Any live cell with more than three live neighbours dies (overcrowding)
// - Any dead cell with exactly three live neighbours becomes alive (reproduction)

// Step performs one generation step across the entire grid
func (g *Game) Step(frame, nextFrame []byte) {
	for y := 0; y < pdgo.LCDRows; y++ {
		// Calculate row indices with wrapping
		rowAboveY := (y - 1 + pdgo.LCDRows) % pdgo.LCDRows
		rowBelowY := (y + 1) % pdgo.LCDRows

		// Get row slices
		rowAbove := frame[rowAboveY*pdgo.LCDRowSize : (rowAboveY+1)*pdgo.LCDRowSize]
		row := frame[y*pdgo.LCDRowSize : (y+1)*pdgo.LCDRowSize]
		rowBelow := frame[rowBelowY*pdgo.LCDRowSize : (rowBelowY+1)*pdgo.LCDRowSize]
		outRow := nextFrame[y*pdgo.LCDRowSize : (y+1)*pdgo.LCDRowSize]

		processRow(rowAbove, row, rowBelow, outRow)
	}
}

// processRow processes one row of the Game of Life simulation
func processRow(rowAbove, row, rowBelow, outRow []byte) {
	var b byte = 0
	bitPos := 0x80

	for x := 0; x < pdgo.LCDColumns; x++ {
		// Count neighbors
		sum := rowSum(rowAbove, x) + middleRowSum(row, x) + rowSum(rowBelow, x)

		// Apply Game of Life rules:
		// - If sum is 3: cell becomes alive
		// - If sum is 2 and cell is alive: stays alive
		// - Otherwise: cell dies
		if sum == 3 || (isCellAlive(row, x) && sum == 2) {
			b |= byte(bitPos)
		}

		bitPos >>= 1

		if bitPos == 0 {
			outRow[x/8] = ^b // Invert because white=1, black=0 in Playdate
			b = 0
			bitPos = 0x80
		}
	}
}

// isCellAlive checks if cell at position x in row is alive (pixel is black = 0)
func isCellAlive(row []byte, x int) bool {
	return (row[x/8] & (0x80 >> (x % 8))) == 0
}

// cellValue returns the value of cell (1 = alive/black, 0 = dead/white)
func cellValue(row []byte, x int) int {
	return 1 - int((row[x/8]>>(7-(x%8)))&1)
}

// rowSum calculates sum of 3 cells in a row (with wrapping at edges)
func rowSum(row []byte, x int) int {
	switch {
	case x == 0:
		return cellValue(row, pdgo.LCDColumns-1) + cellValue(row, x) + cellValue(row, x+1)
	case x < pdgo.LCDColumns-1:
		return cellValue(row, x-1) + cellValue(row, x) + cellValue(row, x+1)
	default:
		return cellValue(row, x-1) + cellValue(row, x) + cellValue(row, 0)
	}
}

// middleRowSum calculates sum of 2 adjacent cells (excluding center cell)
func middleRowSum(row []byte, x int) int {
	switch {
	case x == 0:
		return cellValue(row, pdgo.LCDColumns-1) + cellValue(row, x+1)
	case x < pdgo.LCDColumns-1:
		return cellValue(row, x-1) + cellValue(row, x+1)
	default:
		return cellValue(row, x-1) + cellValue(row, 0)
	}
}
