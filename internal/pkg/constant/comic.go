package constant

import (
	"fmt"
	"strings"
)

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

func GetAllComicTypes() []ComicType {
	return []ComicType{
		ComicTypeManga,
		ComicTypeManhwa,
		ComicTypeManhua,
		ComicTypeComic,
		ComicTypeNovel,
	}
}

func ComicTypeValidationMessage(field string) string {
	types := GetAllComicTypes()
	var res strings.Builder
	res.WriteString(string(types[0]))
	for _, t := range types[1:] {
		res.WriteString(", " + string(t))
	}
	return fmt.Sprintf(
		"%s must be a valid comic type (%s)",
		field,
		res.String(),
	)
}

func GetAllowedComicTypes() map[ComicType]bool {
	return map[ComicType]bool{
		ComicTypeManga:  true,
		ComicTypeManhwa: true,
		ComicTypeManhua: true,
		ComicTypeComic:  true,
		ComicTypeNovel:  true,
	}
}

func GetAllComicStatuses() []ComicStatus {
	return []ComicStatus{
		ComicStatusOngoing,
		ComicStatusCompleted,
		ComicStatusHiatus,
		ComicStatusCancelled,
	}
}

func ComicStatusValidationMessage(field string) string {
	statuses := GetAllComicStatuses()
	var res strings.Builder
	res.WriteString(string(statuses[0]))
	for _, s := range statuses[1:] {
		res.WriteString(", " + string(s))
	}
	return fmt.Sprintf(
		"%s must be a valid comic status (%s)",
		field,
		res.String(),
	)
}

func GetAllowedComicStatuses() map[ComicStatus]bool {
	return map[ComicStatus]bool{
		ComicStatusOngoing:   true,
		ComicStatusCompleted: true,
		ComicStatusHiatus:    true,
		ComicStatusCancelled: true,
	}
}

func GetAllComicAgeRatings() []ComicAgeRating {
	return []ComicAgeRating{
		AgeRatingAll,
		AgeRating13Plus,
		AgeRating16Plus,
		AgeRating18Plus,
	}
}

func ComicAgeRatingValidationMessage(field string) string {
	ratings := GetAllComicAgeRatings()
	var res strings.Builder
	res.WriteString(string(ratings[0]))
	for _, r := range ratings[1:] {
		res.WriteString(", " + string(r))
	}

	return fmt.Sprintf(
		"%s must be a valid age rating (%s)",
		field,
		res.String(),
	)
}

func GetAllowedComicAgeRatings() map[ComicAgeRating]bool {
	return map[ComicAgeRating]bool{
		AgeRatingAll:    true,
		AgeRating13Plus: true,
		AgeRating16Plus: true,
		AgeRating18Plus: true,
	}
}
