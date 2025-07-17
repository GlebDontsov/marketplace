package repository

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Login    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Ads      []Advertisement
}

type Advertisement struct {
	gorm.Model
	Title       string  `gorm:"size:100;not null"`
	Description string  `gorm:"size:1000;not null"`
	ImageURL    string  `gorm:"not null"`
	Price       float64 `gorm:"not null"`
	UserID      uint    `gorm:"not null"`
	User        User    `gorm:"foreignKey:UserID"`
}

type UserRepository interface {
	Create(user *User) error
	GetByLogin(login string) (*User, error)
	GetByID(id uint) (*User, error)
}

type AdvertisementRepository interface {
	Create(ad *Advertisement) error
	GetAll(filter AdFilter) ([]Advertisement, error)
	GetByUserID(userID uint) ([]Advertisement, error)
}

type AdFilter struct {
	Page      int
	Limit     int
	SortBy    string
	SortOrder string
	MinPrice  float64
	MaxPrice  float64
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

type advertisementRepository struct {
	db *gorm.DB
}

func NewAdvertisementRepository(db *gorm.DB) AdvertisementRepository {
	return &advertisementRepository{db: db}
}

func (r *userRepository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByLogin(login string) (*User, error) {
	var user User
	err := r.db.Where("login = ?", strings.ToLower(login)).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByID(id uint) (*User, error) {
	var user User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *advertisementRepository) Create(ad *Advertisement) error {
	return r.db.Create(ad).Error
}

func (r *advertisementRepository) GetAll(filter AdFilter) ([]Advertisement, error) {
	var ads []Advertisement

	query := r.db.Model(&Advertisement{}).Preload("User")

	if filter.MinPrice > 0 {
		query = query.Where("price >= ?", filter.MinPrice)
	}
	if filter.MaxPrice > 0 {
		query = query.Where("price <= ?", filter.MaxPrice)
	}

	sortOrder := "DESC"
	if strings.ToLower(filter.SortOrder) == "asc" {
		sortOrder = "ASC"
	}

	sortField := "created_at"
	if filter.SortBy == "price" {
		sortField = "price"
	}

	query = query.Order(sortField + " " + sortOrder)

	if filter.Page > 0 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset).Limit(filter.Limit)
	}

	err := query.Find(&ads).Error
	if err != nil {
		return nil, err
	}

	return ads, nil
}

func (r *advertisementRepository) GetByUserID(userID uint) ([]Advertisement, error) {
	var ads []Advertisement
	err := r.db.Where("user_id = ?", userID).Find(&ads).Error
	if err != nil {
		return nil, err
	}
	return ads, nil
}
