package userseeder

import (
	"context"
	"manga-go/internal/pkg/model"
	rolerepo "manga-go/internal/pkg/repo/role"
	userrepo "manga-go/internal/pkg/repo/user"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

const defaultAdminEmail = "admin@manga.com"
const defaultAdminPassword = "Admin@123456"

func adminEmail() string {
	if v := os.Getenv("SEED_ADMIN_EMAIL"); v != "" {
		return v
	}
	return defaultAdminEmail
}

func adminPassword() string {
	if v := os.Getenv("SEED_ADMIN_PASSWORD"); v != "" {
		return v
	}
	return defaultAdminPassword
}

type UserSeeder struct {
	userRepo *userrepo.UserRepository
	roleRepo *rolerepo.RoleRepo
}

func NewUserSeeder(userRepo *userrepo.UserRepository, roleRepo *rolerepo.RoleRepo) *UserSeeder {
	return &UserSeeder{userRepo: userRepo, roleRepo: roleRepo}
}

func (s *UserSeeder) Name() string {
	return "UserSeeder"
}

func (s *UserSeeder) Seed(ctx context.Context) error {
	email := adminEmail()
	hashed, err := bcrypt.GenerateFromPassword([]byte(adminPassword()), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := model.User{
		Name:     "Admin",
		Email:    email,
		Password: string(hashed),
	}
	if err := s.userRepo.DB.WithContext(ctx).
		Where(clause.Eq{Column: "email", Value: email}).
		FirstOrCreate(&user).Error; err != nil {
		return err
	}

	var adminRole model.Role
	if err := s.roleRepo.DB.WithContext(ctx).
		Where(clause.Eq{Column: "name", Value: "admin"}).
		First(&adminRole).Error; err != nil {
		return err
	}

	return s.userRepo.AssignRoles(ctx, user.ID, []*model.Role{&adminRole})
}
