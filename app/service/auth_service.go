package service

import (
	"awesomeProject/app/models"
	"awesomeProject/app/repository"
	"awesomeProject/app/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type AuthService interface {
	Login(c *gin.Context)
}

type AuthServiceImpl struct {
	authRepository repository.AuthRepository
}

func (u AuthServiceImpl) Login(c *gin.Context) {
	var request models.LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email := request.Email
	data, err := u.authRepository.FindUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email"})
		return
	}

	if !utils.CheckPasswordHash(request.Password, data.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		Name: data.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	var (
		JWTKey = []byte("example_secret_key_12345")
	)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWTKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	err = utils.AddToRedis(tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save token"})
		return
	}

	c.JSON(http.StatusOK, models.TokenResponse{AccessToken: tokenString, RefreshToken: "dasdasdad"})
}

func AuthServiceInit(authRepository repository.AuthRepository) *AuthServiceImpl {
	return &AuthServiceImpl{
		authRepository: authRepository,
	}
}
