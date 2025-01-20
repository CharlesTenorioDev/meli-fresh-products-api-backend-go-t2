package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testStructure struct {
	ID   int
	Data string
}

func Test__GetBiggestId__WhenMapIsFilled(t *testing.T) {
	biggest := GetBiggestID(map[int]testStructure{
		1: {1, "Struct 1"},
		2: {1, "Struct 2"},
		3: {1, "Struct 3"},
	})
	require.Equal(t, 3, biggest)
}

func Test__GetBiggestId__WhenMapIsEmpty(t *testing.T) {
	biggest := GetBiggestID(map[int]testStructure{})
	require.Equal(t, 1, biggest)
}
