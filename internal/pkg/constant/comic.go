package constant

// ComicType represents the type of comic.
type ComicType string

const (
	ComicTypeManga  ComicType = "manga"
	ComicTypeManhwa ComicType = "manhwa"
	ComicTypeManhua ComicType = "manhua"
	ComicTypeComic  ComicType = "comic"
	ComicTypeNovel  ComicType = "novel"
)

// ComicStatus represents the publication status of a comic.
type ComicStatus string

const (
	ComicStatusOngoing   ComicStatus = "ongoing"
	ComicStatusCompleted ComicStatus = "completed"
	ComicStatusHiatus    ComicStatus = "hiatus"
	ComicStatusCancelled ComicStatus = "cancelled"
)

// ComicAgeRating represents the age rating of a comic.
type ComicAgeRating string

const (
	AgeRatingAll    ComicAgeRating = "ALL"
	AgeRating13Plus ComicAgeRating = "T"
	AgeRating16Plus ComicAgeRating = "16+"
	AgeRating18Plus ComicAgeRating = "18+"
)
