package model

import "testing"

func TestTableNameContracts(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{name: "Author", got: (Author{}).TableName(), want: "authors"},
		{name: "Chapter", got: (Chapter{}).TableName(), want: "chapters"},
		{name: "Comic", got: (Comic{}).TableName(), want: "comics"},
		{name: "ComicAuthor", got: (ComicAuthor{}).TableName(), want: "comic_authors"},
		{name: "ComicGenre", got: (ComicGenre{}).TableName(), want: "comic_genres"},
		{name: "ComicTag", got: (ComicTag{}).TableName(), want: "comic_tags"},
		{name: "Comment", got: (Comment{}).TableName(), want: "comments"},
		{name: "Genre", got: (Genre{}).TableName(), want: "genres"},
		{name: "Page", got: (Page{}).TableName(), want: "pages"},
		{name: "Permission", got: (Permission{}).TableName(), want: "permissions"},
		{name: "Rating", got: (Rating{}).TableName(), want: "ratings"},
		{name: "ReadingHistory", got: (ReadingHistory{}).TableName(), want: "reading_histories"},
		{name: "ReadingProgress", got: (ReadingProgress{}).TableName(), want: "reading_progresses"},
		{name: "Role", got: (Role{}).TableName(), want: "roles"},
		{name: "Tag", got: (Tag{}).TableName(), want: "tags"},
		{name: "TranslationGroup", got: (TranslationGroup{}).TableName(), want: "translation_groups"},
		{name: "User", got: (User{}).TableName(), want: "users"},
		{name: "UserComicRead", got: (UserComicRead{}).TableName(), want: "user_comic_reads"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("expected %s table name = %s, got %s", tt.name, tt.want, tt.got)
			}
		})
	}
}
