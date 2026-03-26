package services

import (
	authorservice "manga-go/internal/pkg/services/author"
	fileservice "manga-go/internal/pkg/services/file"
	genreservice "manga-go/internal/pkg/services/genre"
	userserivce "manga-go/internal/pkg/services/user"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"services",
	userserivce.Module,
	authorservice.Module,
	genreservice.Module,
	fileservice.Module,
)
