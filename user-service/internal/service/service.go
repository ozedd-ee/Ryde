package service

import (
	"context"
	"errors"
	"ryde/internal/data"
	"ryde/internal/models"
	"ryde/utils"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	userStore *data.UserStore
}

func NewUserService(userStore *data.UserStore) *UserService {
	return &UserService{
		userStore: userStore,
	}
}

func (s *UserService) SignUp(ctx context.Context, user *models.User) (*models.User, error) {
	exists, err := s.userStore.GetUserByEmail(ctx, user.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	} else if exists != nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword
	return s.userStore.CreateUser(ctx, user)
}

func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userStore.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	err = utils.CheckPasswordMatch(password, user.Password)
	if err != nil {
		return "", errors.New("incorrect password")
	}

	token, err := utils.GenerateJWT(user.ID.Hex())
	if err != nil {
		return "", errors.New("unable to generate JWT")
	}
	return token, nil
}

func (s *UserService) GetUser(ctx context.Context, userID string) (*models.User, error) {
	return s.userStore.GetUser(ctx, userID)
}
