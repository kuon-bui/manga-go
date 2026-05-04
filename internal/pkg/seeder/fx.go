package seeder

import (
	"manga-go/internal/pkg/logger"
	activityseeder "manga-go/internal/pkg/seeder/activity"
	authorseeder "manga-go/internal/pkg/seeder/author"
	comicseeder "manga-go/internal/pkg/seeder/comic"
	genreseeder "manga-go/internal/pkg/seeder/genre"
	notificationseeder "manga-go/internal/pkg/seeder/notification"
	permissionseeder "manga-go/internal/pkg/seeder/permission"
	roleseeder "manga-go/internal/pkg/seeder/role"
	tagseeder "manga-go/internal/pkg/seeder/tag"
	translationgroupseeder "manga-go/internal/pkg/seeder/translation_group"
	userseeder "manga-go/internal/pkg/seeder/user"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

// SeederRunnerParams holds all individual seeders in the correct dependency order.
type SeederRunnerParams struct {
	fx.In

	Db                     *gorm.DB
	Logger                 *logger.Logger
	PermissionSeeder       *permissionseeder.PermissionSeeder
	RoleSeeder             *roleseeder.RoleSeeder
	UserSeeder             *userseeder.UserSeeder
	AuthorSeeder           *authorseeder.AuthorSeeder
	GenreSeeder            *genreseeder.GenreSeeder
	TagSeeder              *tagseeder.TagSeeder
	TranslationGroupSeeder *translationgroupseeder.TranslationGroupSeeder
	ComicSeeder            *comicseeder.ComicSeeder
	ActivitySeeder         *activityseeder.ActivitySeeder
	NotificationSeeder     *notificationseeder.NotificationSeeder
}

func newSeederRunner(p SeederRunnerParams) *SeederRunner {
	// Order is significant: permissions → roles → users → authors → genres → tags → translation groups → comics → activity → notifications
	seeders := []Seeder{
		p.PermissionSeeder,
		p.RoleSeeder,
		p.UserSeeder,
		p.AuthorSeeder,
		p.GenreSeeder,
		p.TagSeeder,
		p.TranslationGroupSeeder,
		p.ComicSeeder,
		p.ActivitySeeder,
		p.NotificationSeeder,
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
		translationgroupseeder.NewTranslationGroupSeeder,
		comicseeder.NewComicSeeder,
		activityseeder.NewActivitySeeder,
		notificationseeder.NewNotificationSeeder,
		newSeederRunner,
	),
)
