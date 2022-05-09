package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testCase struct {
	Input  Cards
	Output any
}

func TestNewCards(t *testing.T) {
	cards := NewCards()

	require.Len(t, cards, 312)
}

func TestCountValue(t *testing.T) {
	testCases := []testCase{
		{
			Input:  Cards{"♣5", "♦Q", "♥K", "♠J"},
			Output: 35,
		},
		{
			Input:  Cards{"♣A", "♦A", "♥2", "♠5"},
			Output: 19,
		},
		{
			Input:  Cards{"♣5", "♦5", "♥5", "♠A"},
			Output: 16,
		},
		{
			Input:  Cards{"♣A", "♦A", "♥A", "♠A"},
			Output: 14,
		},
		{
			Input:  Cards{"♣A", "♦A", "♥A", "♠A", "♥A", "♠A", "♥A", "♠A", "♥K", "♥3"},
			Output: 21,
		},
	}

	for _, testCase := range testCases {
		require.Equal(t, testCase.Output, testCase.Input.CountValue())
	}
}

func TestIsBlackjack(t *testing.T) {
	testCases := []testCase{
		{
			Input:  Cards{"♣5", "♦Q", "♥K", "♠J"},
			Output: false,
		},
		{
			Input:  Cards{"♣A", "♦A", "♥A", "♠A", "♥A", "♠A", "♥A", "♠A", "♥K", "♥3"},
			Output: false,
		},
		{
			Input:  Cards{"♣A", "♥K"},
			Output: true,
		},
		{
			Input:  Cards{"♣A", "♥10"},
			Output: true,
		},
	}

	for _, testCase := range testCases {
		require.Equal(t, testCase.Output, testCase.Input.IsBlackjack())
	}
}
