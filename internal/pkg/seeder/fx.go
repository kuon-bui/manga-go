package seeder

import (
	"manga-go/internal/pkg/logger"
	authorseeder "manga-go/internal/pkg/seeder/author"
	comicseeder "manga-go/internal/pkg/seeder/comic"
	genreseeder "manga-go/internal/pkg/seeder/genre"
	permissionseeder "manga-go/internal/pkg/seeder/permission"
	roleseeder "manga-go/internal/pkg/seeder/role"
	tagseeder "manga-go/internal/pkg/seeder/tag"
	userseeder "manga-go/internal/pkg/seeder/user"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

// SeederRunnerParams holds all individual seeders in the correct dependency order.
type SeederRunnerParams struct {
	fx.In

	Db               *gorm.DB
	Logger           *logger.Logger
	PermissionSeeder *permissionseeder.PermissionSeeder
	RoleSeeder       *roleseeder.RoleSeeder
	UserSeeder       *userseeder.UserSeeder
	AuthorSeeder     *authorseeder.AuthorSeeder
	GenreSeeder      *genreseeder.GenreSeeder
	TagSeeder        *tagseeder.TagSeeder
	ComicSeeder      *comicseeder.ComicSeeder
}

func newSeederRunner(p SeederRunnerParams) *SeederRunner {
	// Order is significant: permissions → roles → users → authors → genres → tags → comics
	seeders := []Seeder{
		p.PermissionSeeder,
		p.RoleSeeder,
		p.UserSeeder,
		p.AuthorSeeder,
		p.GenreSeeder,
		p.TagSeeder,
		p.ComicSeeder,
	}
	return NewSeederRunner(seeders, p.Logger, p.Db)
}

var Module = fx.Module(
	"seeder",
	fx.Provide(
		permissionseeder.NewPermissionSeeder,
		roleseeder.NewRoleSeeder,
		userseeder.NewUserSeeder,
		authorseeder.NewAuthorSeeder,
		genreseeder.NewGenreSeeder,
		tagseeder.NewTagSeeder,
		comicseeder.NewComicSeeder,
		newSeederRunner,
	),
)
