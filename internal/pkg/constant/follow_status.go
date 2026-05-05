package constant

import (
	"fmt"
	"math/rand"
	"strings"
)

type FollowStatus string

const (
	FollowStatusReading   FollowStatus = "reading"
	FollowStatusPlanned   FollowStatus = "planned"
	FollowStatusCompleted FollowStatus = "completed"
	FollowStatusDropped   FollowStatus = "dropped"
	FollowStatusFavorite  FollowStatus = "favorite"
)

func GetAllFollowStatuses() []FollowStatus {
	return []FollowStatus{
		FollowStatusReading,
		FollowStatusPlanned,
		FollowStatusCompleted,
		FollowStatusDropped,
		FollowStatusFavorite,
	}
}

func FollowStatusValidationMessage(field string) string {
	flowStatuses := GetAllFollowStatuses()
	var res strings.Builder
	res.WriteString(string(flowStatuses[0]))
	for _, status := range flowStatuses[1:] {
		res.WriteString(", " + string(status))
	}

	return fmt.Sprintf(
		"%s must be a valid follow status (%s)",
		field,
		res.String(),
	)
}

func GetAllowedFollowStatuses() map[FollowStatus]bool {
	return map[FollowStatus]bool{
		FollowStatusReading:   true,
		FollowStatusPlanned:   true,
		FollowStatusCompleted: true,
		FollowStatusDropped:   true,
		FollowStatusFavorite:  true,
	}
}

func FollowStatusRandom() FollowStatus {
	statuses := []FollowStatus{
		FollowStatusReading,
		FollowStatusPlanned,
		FollowStatusCompleted,
		FollowStatusDropped,
		FollowStatusFavorite,
	}

	return statuses[rand.Intn(len(statuses))]
}
