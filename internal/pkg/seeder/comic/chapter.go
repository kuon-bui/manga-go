package comicseeder

import "time"

func resolveChapterPublishedAt(comicSlug string, chapterIndex int) *time.Time {
	base := time.Date(2024, time.January, 1, 9, 0, 0, 0, time.UTC)
	offsetDays := checksumString(comicSlug) % 365
	publishedAt := base.AddDate(0, 0, offsetDays+(chapterIndex-1)*2)
	return &publishedAt
}

func checksumString(value string) int {
	checksum := 0
	for _, char := range value {
		checksum += int(char)
	}

	return checksum
}
