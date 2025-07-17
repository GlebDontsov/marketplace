package controller

import (
	"marketplace/internal/dto"
	"marketplace/internal/repository"
	"marketplace/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService: userService}
}

type AdvertisementController struct {
	adService service.AdvertisementService
}

func NewAdvertisementController(adService service.AdvertisementService) *AdvertisementController {
	return &AdvertisementController{adService: adService}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with login and password
// @Tags auth
// @Accept  json
// @Produce json
// @Param input body dto.RegisterRequest true "User credentials"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /register [post]
func (c *UserController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	user, err := c.userService.Register(req.Login, req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.UserResponse{
		ID:    user.ID,
		Login: user.Login,
	})
}

// Login godoc
// @Summary Login user
// @Description Login user with credentials
// @Tags auth
// @Accept  json
// @Produce json
// @Param input body dto.LoginRequest true "User credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /login [post]
func (c *UserController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	token, err := c.userService.Login(req.Login, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "invalid credentials"})
		return
	}

	ctx.JSON(http.StatusOK, dto.LoginResponse{Token: token})
}

// CreateAd godoc
// @Summary Create new advertisement
// @Description Create new advertisement (only for authorized users)
// @Tags advertisement
// @Accept  json
// @Produce json
// @securityDefinitions.apikey BearerAuth
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param input body dto.CreateAdRequest true "Advertisement data"
// @Success 201 {object} dto.AdvertisementResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /ads [post]
func (c *AdvertisementController) CreateAd(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	var req dto.CreateAdRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	ad, err := c.adService.CreateAd(
		userID.(uint),
		req.Title,
		req.Description,
		req.ImageURL,
		req.Price,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.AdvertisementResponse{
		ID:          ad.ID,
		Title:       ad.Title,
		Description: ad.Description,
		ImageURL:    ad.ImageURL,
		Price:       ad.Price,
		CreatedAt:   ad.CreatedAt,
	})
}

// GetAds godoc
// @Summary Get advertisements
// @Description Get list of advertisements with filters and pagination
// @Tags advertisement
// @Accept  json
// @Produce json
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @Param Authorization header string false "Bearer token (optional)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param sort_by query string false "Sort field (created_at or price)" default(created_at)
// @Param sort_order query string false "Sort order (asc or desc)" default(desc)
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Success 200 {array} dto.AdvertisementResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /ads [get]
func (c *AdvertisementController) GetAds(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	sortBy := ctx.DefaultQuery("sort_by", "created_at")
	sortOrder := ctx.DefaultQuery("sort_order", "desc")
	minPrice, _ := strconv.ParseFloat(ctx.Query("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(ctx.Query("max_price"), 64)

	var userIDPtr *uint
	if userID, exists := ctx.Get("userID"); exists {
		if id, ok := userID.(uint); ok {
			userIDPtr = &id
		}
	}

	filter := repository.AdFilter{
		Page:      page,
		Limit:     limit,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		MinPrice:  minPrice,
		MaxPrice:  maxPrice,
	}

	ads, err := c.adService.GetAds(filter, userIDPtr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, ads)
}
