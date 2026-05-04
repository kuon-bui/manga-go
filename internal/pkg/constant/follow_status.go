package constant

import "math/rand"

type FollowStatus string

const (
	FollowStatusReading   FollowStatus = "reading"
	FollowStatusPlanned   FollowStatus = "planned"
	FollowStatusCompleted FollowStatus = "completed"
	FollowStatusDropped   FollowStatus = "dropped"
	FollowStatusFavorite  FollowStatus = "favorite"
)

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
