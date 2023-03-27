package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_levenshtein_distance(t *testing.T) {
	tests := []struct {
		Distance int
		Desired  string
		Input    string
	}{
		{2, "flaw", "lawn"},
		{3, "Manhattan", "Manahaton"},
		{5, "Clear", ""},
		{5, "", "Clear"},
	}

	for _, test := range tests {
		assert.Equal(t, test.Distance, levenshtein_distance(test.Desired, test.Input))
	}
}

func Test_min(t *testing.T) {
	assert.Equal(t, 1, min(1, 2, 3))
	assert.Equal(t, 1, min(2, 1, 3))
	assert.Equal(t, 1, min(2, 3, 1))
}
