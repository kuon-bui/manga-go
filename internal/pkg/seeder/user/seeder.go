package userseeder

import (
	"context"
	"errors"
	"manga-go/internal/pkg/config"
	"manga-go/internal/pkg/model"
	rolerepo "manga-go/internal/pkg/repo/role"
	userrepo "manga-go/internal/pkg/repo/user"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const defaultAdminEmail = "admin@manga.com"
const defaultAdminPassword = "Admin@123456"

func adminEmail(cfg *config.Config) string {
	if cfg.Seeder.AdminEmail != "" {
		return cfg.Seeder.AdminEmail
	}

	return defaultAdminEmail
}

func adminPassword(cfg *config.Config) string {
	if cfg.Seeder.AdminPassword != "" {
		return cfg.Seeder.AdminPassword
	}

	return defaultAdminPassword
}

type UserSeeder struct {
	userRepo *userrepo.UserRepository
	roleRepo *rolerepo.RoleRepo
	config   *config.Config
}

func NewUserSeeder(userRepo *userrepo.UserRepository, roleRepo *rolerepo.RoleRepo, config *config.Config) *UserSeeder {
	return &UserSeeder{userRepo: userRepo, roleRepo: roleRepo, config: config}
}

func (s *UserSeeder) Name() string {
	return "UserSeeder"
}

func (s *UserSeeder) Seed(ctx context.Context) error {
	email := adminEmail(s.config)

	user, err := s.userRepo.FindOne(ctx, []any{clause.Eq{Column: "email", Value: email}}, nil)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		hashed, err := bcrypt.GenerateFromPassword([]byte(adminPassword(s.config)), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user = &model.User{
			Name:     "Admin",
			Email:    email,
			Password: string(hashed),
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return err
		}
	}

	adminRole, err := s.roleRepo.FindOne(ctx, []any{clause.Eq{Column: "name", Value: "admin"}}, nil)
	if err != nil {
		return err
	}

	return s.userRepo.AssignRoles(ctx, user.ID, []*model.Role{adminRole})
}
