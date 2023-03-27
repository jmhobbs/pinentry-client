package main

import (
	"strings"
)

func quality(input string) int {
	qlty := 100 - (levenshtein_distance("unicorn", strings.ToLower(input)) * 14)
	if qlty < 0 {
		return 0
	}
	return qlty
}

func levenshtein_distance(desired, target string) int {
	rows := len(desired) + 1
	cols := len(target) + 1

	distributions := make([][]int, rows)
	for r := 0; r < rows; r++ {
		distributions[r] = make([]int, cols)
		distributions[r][0] = r
	}
	for c := 1; c < cols; c++ {
		distributions[0][c] = c
	}

	for col := 1; col < cols; col++ {
		for row := 1; row < rows; row++ {
			cost := 1
			if desired[row-1] == target[col-1] {
				cost = 0
			}
			distributions[row][col] = min(
				distributions[row-1][col]+1,
				distributions[row][col-1]+1,
				distributions[row-1][col-1]+cost,
			)
		}
	}

	return distributions[rows-1][cols-1]
}

func min(a, b, c int) int {
	if a < b && a < c {
		return a
	}
	if b < a && b < c {
		return b
	}
	return c
}
