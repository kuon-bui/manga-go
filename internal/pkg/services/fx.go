package services

import (
	authorservice "manga-go/internal/pkg/services/author"
	chapterservice "manga-go/internal/pkg/services/chapter"
	comicservice "manga-go/internal/pkg/services/comic"
	fileservice "manga-go/internal/pkg/services/file"
	genreservice "manga-go/internal/pkg/services/genre"
	permissionservice "manga-go/internal/pkg/services/permission"
	roleservice "manga-go/internal/pkg/services/role"
	tagservice "manga-go/internal/pkg/services/tag"
	translationgroupservice "manga-go/internal/pkg/services/translation_group"
	userservice "manga-go/internal/pkg/services/user"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"services",
	userservice.Module,
	authorservice.Module,
	genreservice.Module,
	fileservice.Module,
	tagservice.Module,
	comicservice.Module,
	chapterservice.Module,
	translationgroupservice.Module,
	roleservice.Module,
	permissionservice.Module,
)
