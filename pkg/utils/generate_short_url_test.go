package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUtils_GenerateEqualShortPaths(t *testing.T) {
	originalURL := "https://ozon.ru"

	shortPath1 := GenerateShortPath(originalURL, 0)
	shortPath2 := GenerateShortPath(originalURL, 0)

	assert.Equal(t, shortPath1, shortPath2, "Generated paths should be equal")
}

func TestUtils_GenerateNotEqualShortPaths(t *testing.T) {
	originalURL := "https://ozon.ru"

	shortPath1 := GenerateShortPath(originalURL, 0)
	shortPath2 := GenerateShortPath(originalURL, 1)

	assert.NotEqual(t, shortPath1, shortPath2, "Generated paths should not be equal")
}
