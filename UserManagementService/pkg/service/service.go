package service

import (
	"advancedProg2Final/UserManagementService/pkg/entity"
	"advancedProg2Final/UserManagementService/pkg/repository"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

func (s *UserService) ValidateToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte("sfjsdifjsjfsf2131edsad!@@!#@!$sadas4324#@$@#$@#$#@^&&*()*&)"), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserByID(id int64) (*entity.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) CreateUser(user *entity.User) (int64, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	return s.repo.Save(user)
}

func (s *UserService) Authenticate(username, password string) (*entity.User, error) {
	user, err := s.repo.Authenticate(username, password)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUser(id int) (*entity.User, error) {
	user, err := s.repo.GetByID(int64(id))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateUser(user *entity.User) error {
	existingUser, err := s.GetUser(int(user.ID))
	if err != nil {
		return err
	}
	existingUser.Email = user.Email
	err = s.repo.SaveUser(existingUser)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) DeleteUser(id int) error {
	_, err := s.GetUser(id)
	if err != nil {
		return err
	}
	err = s.repo.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}
