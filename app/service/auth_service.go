package service

import (
	"awesomeProject/app/models"
	"awesomeProject/app/repository"
	"awesomeProject/app/utils"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"time"
)

type AuthService interface {
	Login(c *gin.Context)
	RefreshToken(c *gin.Context)
	Logout(c *gin.Context)
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
		fmt.Println("Password check failed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	} else {
		fmt.Println("Password check succeeded")
	}

	accessTokenExpirationTime := time.Now().Add(24 * time.Hour)
	accessClaims := &models.Claims{
		Id: data.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessTokenExpirationTime.Unix(),
		},
	}

	accessToken, err := utils.GenerateToken(accessTokenExpirationTime, *accessClaims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshTokenExpirationTime := time.Now().Add(7 * 24 * time.Hour)
	refreshClaims := &models.Claims{
		Id: data.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshTokenExpirationTime.Unix(),
		},
	}

	refreshToken, err := utils.GenerateToken(refreshTokenExpirationTime, *refreshClaims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	err = utils.AddToRedis(strconv.Itoa(data.ID)+"_access", accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save access token"})
		return
	}

	err = utils.AddToRedis(strconv.Itoa(data.ID)+"_refresh", refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save refresh token"})
		return
	}

	c.JSON(http.StatusOK, models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (u AuthServiceImpl) RefreshToken(c *gin.Context) {
	var request models.TokenRefreshRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshToken := request.RefreshToken

	claims := &models.Claims{}
	tkn, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil || !tkn.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	redisRefreshToken, err := utils.GetFromRedis(strconv.Itoa(claims.Id) + "_refresh")

	if err != nil || redisRefreshToken != refreshToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not found or invalid"})
		return
	}

	accessTokenExpirationTime := time.Now().Add(24 * time.Hour)
	accessClaims := &models.Claims{
		Id: claims.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessTokenExpirationTime.Unix(),
		},
	}

	newAccessToken, err := utils.GenerateToken(accessTokenExpirationTime, *accessClaims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new access token"})
		return
	}

	err = utils.AddToRedis(strconv.Itoa(claims.Id)+"_access", newAccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save new access token"})
		return
	}

	c.JSON(http.StatusOK, models.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: refreshToken,
	})
}

func (u AuthServiceImpl) Logout(c *gin.Context) {
	var request models.LogoutRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshToken := request.RefreshToken

	err := utils.RemoveFromRedis(refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove token from Redis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func AuthServiceInit(authRepository repository.AuthRepository) *AuthServiceImpl {
	return &AuthServiceImpl{
		authRepository: authRepository,
	}
}
