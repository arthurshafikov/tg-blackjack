package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testCase struct {
	Input    any
	Expected any
}

type testCaseWithCards struct {
	Input    Cards
	Expected any
}

func TestNewCards(t *testing.T) {
	testCases := []testCase{
		{
			Input:    6,
			Expected: 312,
		},
		{
			Input:    1,
			Expected: 52,
		},
		{
			Input:    10,
			Expected: 520,
		},
	}

	for _, testCase := range testCases {
		require.Len(t, NewCards(testCase.Input.(int)), testCase.Expected.(int))
	}
}

func TestCountValue(t *testing.T) {
	testCases := []testCaseWithCards{
		{
			Input:    Cards{"♣5", "♦Q", "♥K", "♠J"},
			Expected: 35,
		},
		{
			Input:    Cards{"♣A", "♦A", "♥2", "♠5"},
			Expected: 19,
		},
		{
			Input:    Cards{"♣5", "♦5", "♥5", "♠A"},
			Expected: 16,
		},
		{
			Input:    Cards{"♣A", "♦A", "♥A", "♠A"},
			Expected: 14,
		},
		{
			Input:    Cards{"♣A", "♦A", "♥A", "♠A", "♥A", "♠A", "♥A", "♠A", "♥K", "♥3"},
			Expected: 21,
		},
	}

	for _, testCase := range testCases {
		require.Equal(t, testCase.Expected, testCase.Input.CountValue())
	}
}

func TestIsBlackjack(t *testing.T) {
	testCases := []testCaseWithCards{
		{
			Input:    Cards{"♣5", "♦Q", "♥K", "♠J"},
			Expected: false,
		},
		{
			Input:    Cards{"♣A", "♦A", "♥A", "♠A", "♥A", "♠A", "♥A", "♠A", "♥K", "♥3"},
			Expected: false,
		},
		{
			Input:    Cards{"♣A", "♥K"},
			Expected: true,
		},
		{
			Input:    Cards{"♣A", "♥10"},
			Expected: true,
		},
	}

	for _, testCase := range testCases {
		require.Equal(t, testCase.Expected, testCase.Input.IsBlackjack())
	}
}
