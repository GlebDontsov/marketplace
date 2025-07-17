package service

import (
	"errors"
	"marketplace/internal/dto"
	"marketplace/internal/repository"
	"strings"
	"time"
	"unicode"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidPassword    = errors.New("password must be at least 8 characters long")
	ErrInvalidLogin       = errors.New("login must be 3-20 characters long and contain only letters and numbers")
	ErrInvalidTitle       = errors.New("title must be 5-100 characters long")
	ErrInvalidDesc        = errors.New("description must be 10-1000 characters long")
	ErrInvalidImageURL    = errors.New("invalid image URL")
	ErrInvalidPrice       = errors.New("price must be positive")
)

const (
	jwtSecret = "your-secret-key"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

type AdvertisementService struct {
	repo repository.AdvertisementRepository
}

func NewAdvertisementService(repo repository.AdvertisementRepository) *AdvertisementService {
	return &AdvertisementService{repo: repo}
}

func (s *UserService) Register(login, password string) (*repository.User, error) {
	if len(login) < 3 || len(login) > 20 {
		return nil, ErrInvalidLogin
	}
	for _, c := range login {
		if !unicode.IsLetter(c) && !unicode.IsNumber(c) {
			return nil, ErrInvalidLogin
		}
	}

	if len(password) < 8 {
		return nil, ErrInvalidPassword
	}

	if _, err := s.repo.GetByLogin(login); err == nil {
		return nil, ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &repository.User{
		Login:    login,
		Password: string(hashedPassword),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(login, password string) (string, error) {
	user, err := s.repo.GetByLogin(login)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AdvertisementService) CreateAd(userID uint, title, description, imageURL string, price float64) (*repository.Advertisement, error) {
	if len(title) < 5 || len(title) > 100 {
		return nil, ErrInvalidTitle
	}
	if len(description) < 10 || len(description) > 1000 {
		return nil, ErrInvalidDesc
	}
	if !strings.HasPrefix(imageURL, "http") {
		return nil, ErrInvalidImageURL
	}
	if price <= 0 {
		return nil, ErrInvalidPrice
	}

	ad := &repository.Advertisement{
		Title:       title,
		Description: description,
		ImageURL:    imageURL,
		Price:       price,
		UserID:      userID,
	}

	if err := s.repo.Create(ad); err != nil {
		return nil, err
	}

	return ad, nil
}

func (s *AdvertisementService) GetAds(filter repository.AdFilter, currentUserID *uint) ([]dto.AdvertisementResponse, error) {
	if filter.SortBy != "created_at" && filter.SortBy != "price" {
		filter.SortBy = "created_at"
	}
	if filter.SortOrder != "asc" && filter.SortOrder != "desc" {
		filter.SortOrder = "desc"
	}

	ads, err := s.repo.GetAll(filter)
	if err != nil {
		return nil, err
	}

	var response []dto.AdvertisementResponse
	for _, ad := range ads {
		item := dto.AdvertisementResponse{
			ID:          ad.ID,
			Title:       ad.Title,
			Description: ad.Description,
			ImageURL:    ad.ImageURL,
			Price:       ad.Price,
			CreatedAt:   ad.CreatedAt,
			AuthorLogin: ad.User.Login,
		}

		if currentUserID != nil && ad.UserID == *currentUserID {
			item.IsOwner = true
		}

		response = append(response, item)
	}

	return response, nil
}
