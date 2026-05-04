package userseeder

import (
	"errors"
	"fmt"
	"manga-go/internal/pkg/config"
	"manga-go/internal/pkg/model"
	rolerepo "manga-go/internal/pkg/repo/role"
	userrepo "manga-go/internal/pkg/repo/user"
	seederutil "manga-go/internal/pkg/seeder/util"

	"github.com/jaswdr/faker/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const defaultAdminEmail = "admin@manga.com"
const defaultAdminPassword = "Admin@123456"
const fakeUserCount = 18

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
	faker    faker.Faker
}

func NewUserSeeder(
	userRepo *userrepo.UserRepository,
	roleRepo *rolerepo.RoleRepo,
	config *config.Config,
	faker faker.Faker,
) *UserSeeder {
	return &UserSeeder{
		userRepo: userRepo,
		roleRepo: roleRepo,
		config:   config,
		faker:    faker,
	}
}

func (s *UserSeeder) Name() string {
	return "UserSeeder"
}

func (s *UserSeeder) Truncate(tx *gorm.DB) error {
	return seederutil.TruncateTables(tx, "users_roles", "users")
}

func (s *UserSeeder) Seed(tx *gorm.DB) error {
	email := adminEmail(s.config)
	createdAdmin := false

	user, err := s.userRepo.FindOneWithTransaction(tx, []any{clause.Eq{Column: "email", Value: email}}, nil)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		hashed, err := bcrypt.GenerateFromPassword([]byte(adminPassword(s.config)), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user = &model.User{}
		user.Fake(s.faker)
		user.Name = "Admin"
		user.Email = email
		user.Password = string(hashed)
		user.UserConfig = model.DefaultUserConfig()
		if err := s.userRepo.CreateWithTransaction(tx, user); err != nil {
			return err
		}
		createdAdmin = true
	}

	adminRole, err := s.roleRepo.FindOneWithTransaction(tx, []any{clause.Eq{Column: "name", Value: "admin"}}, nil)
	if err != nil {
		return err
	}
	readerRole, err := s.roleRepo.FindOneWithTransaction(tx, []any{clause.Eq{Column: "name", Value: "reader"}}, nil)
	if err != nil {
		return err
	}
	translatorRole, err := s.roleRepo.FindOneWithTransaction(tx, []any{clause.Eq{Column: "name", Value: "translator"}}, nil)
	if err != nil {
		return err
	}

	if err := s.fakeUsers(tx, fakeUserCount, readerRole, translatorRole); err != nil {
		return err
	}

	if !createdAdmin {
		return nil
	}

	return s.userRepo.AssignRolesWithTransaction(tx, user.ID, []*model.Role{adminRole})
}

func (s *UserSeeder) fakeUsers(tx *gorm.DB, count int, readerRole, translatorRole *model.Role) error {
	for index := 1; index <= count; index++ {
		email := fmt.Sprintf("seed-user-%02d@manga.local", index)
		user, err := s.userRepo.FindOneWithTransaction(tx, []any{clause.Eq{Column: "email", Value: email}}, nil)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = &model.User{}
			user.Fake(s.faker)
			user.Email = email
			user.Name = fmt.Sprintf("%s %02d", user.Name, index)
			if err := s.userRepo.CreateWithTransaction(tx, user); err != nil {
				return err
			}
		}

		roles := []*model.Role{readerRole}
		if index <= 6 {
			roles = []*model.Role{readerRole, translatorRole}
		}

		if err := s.userRepo.AssignRolesWithTransaction(tx, user.ID, roles); err != nil {
			return err
		}
	}

	return nil
}
