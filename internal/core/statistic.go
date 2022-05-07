package core

import "sort"

type UsersStatistics map[string]int

type SortedUserStatistics []userStatistics

type userStatistics struct {
	Username string
	Points   int
}

func (u UsersStatistics) SortByValue() SortedUserStatistics {
	var sortedUserStats SortedUserStatistics //nolint
	for username, points := range u {
		sortedUserStats = append(sortedUserStats, userStatistics{
			Username: username,
			Points:   points,
		})
	}

	sort.Slice(sortedUserStats, func(i, j int) bool {
		return sortedUserStats[i].Points > sortedUserStats[j].Points
	})

	return sortedUserStats
}
