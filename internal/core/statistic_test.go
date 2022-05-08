package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	firstUser  = "first_user"
	secondUser = "second_user"
	thirdUser  = "third_user"
	fourthUser = "fourth_user"
)

func TestSortByValue(t *testing.T) {
	stats := UsersStatistics{
		thirdUser:  -3,
		fourthUser: -15,
		secondUser: 0,
		firstUser:  5,
	}
	expected := SortedUserStatistics{
		{
			Username: firstUser,
			Points:   5,
		},
		{
			Username: secondUser,
			Points:   0,
		},
		{
			Username: thirdUser,
			Points:   -3,
		},
		{
			Username: fourthUser,
			Points:   -15,
		},
	}

	result := stats.SortByValue()

	require.Equal(t, expected, result)
}
