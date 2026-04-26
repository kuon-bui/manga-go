package constant

type FollowStatus string

const (
	FollowStatusReading   FollowStatus = "reading"
	FollowStatusPlanned   FollowStatus = "planned"
	FollowStatusCompleted FollowStatus = "completed"
	FollowStatusDropped   FollowStatus = "dropped"
	FollowStatusFavorite  FollowStatus = "favorite"
)
