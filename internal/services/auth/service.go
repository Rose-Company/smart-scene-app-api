package auth

import (
	"context"
	"smart-scene-app-api/common"
	"smart-scene-app-api/internal/models"
	userModel "smart-scene-app-api/internal/models/user"
	user "smart-scene-app-api/internal/repositories/user"
	"smart-scene-app-api/pkg/jwt"
	"smart-scene-app-api/server"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service interface {
	Login(email, password string) (*userModel.User, string, error)
	Register(email, password, fullName string) (*userModel.User, string, error)
}

type authService struct {
	sc      server.ServerContext
	useRepo user.Repository
}

func NewAuthService(sc server.ServerContext) Service {
	return &authService{
		sc:      sc,
		useRepo: user.NewRepository(sc.DB()),
	}
}

func (s *authService) Login(email, password string) (*userModel.User, string, error) {
	user, err := s.useRepo.GetDetailByConditions(s.sc.Ctx(), func(tx *gorm.DB) {
		tx.Where("email = ?", email)
	})
	if err != nil {
		return nil, "", common.ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", common.ErrInvalidPassword
	}

	token, err := jwt.GenerateToken(user.ID, user.RoleID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *authService) Register(email, password, fullName string) (*userModel.User, string, error) {
	existingUser, err := s.useRepo.GetDetailByConditions(s.sc.Ctx(), func(tx *gorm.DB) {
		tx.Where("email = ?", email)
	})
	if err == nil && existingUser != nil {
		return nil, "", common.ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	UserRoleID, err := s.GetRoleIDByName(s.sc.Ctx(), "user")
	if err != nil {
		return nil, "", err
	}

	newUser := &userModel.User{
		Base: models.Base{
			ID: uuid.New(),
		},
		Email:    email,
		Password: string(hashedPassword),
		FullName: fullName,
		RoleID:   UserRoleID,
		Status:   "active",
	}

	_, err = s.useRepo.Create(s.sc.Ctx(), newUser)
	if err != nil {
		return nil, "", err
	}

	token, err := jwt.GenerateToken(newUser.ID, newUser.RoleID)
	if err != nil {
		return nil, "", err
	}

	return newUser, token, nil
}


func (s *authService) GetRoleIDByName(context context.Context,name string) (uuid.UUID, error) {
	user, err := s.useRepo.GetDetailByConditions(context, func(tx *gorm.DB) {
		tx.Where("name = ?", name)
	})
	if err != nil {
		return uuid.Nil, err
	}
	return user.RoleID, nil
}