package constant

// ComicType represents the type of comic.
type ComicType string

const (
	ComicTypeManga  ComicType = "manga"
	ComicTypeManhwa ComicType = "manhwa"
	ComicTypeManhua ComicType = "manhua"
	ComicTypeComic  ComicType = "comic"
)

// ComicStatus represents the publication status of a comic.
type ComicStatus string

const (
	ComicStatusOngoing   ComicStatus = "ongoing"
	ComicStatusCompleted ComicStatus = "completed"
	ComicStatusHiatus    ComicStatus = "hiatus"
	ComicStatusCancelled ComicStatus = "cancelled"
)
