package constant

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
