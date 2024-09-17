package service

import (
	"awesomeProject/app/constant"
	"awesomeProject/app/domain/dao"
	"awesomeProject/app/pkg"
	"awesomeProject/app/repository"
	"awesomeProject/app/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type UserService interface {
	GetMe(c *gin.Context)
	GetAllUser(c *gin.Context)
	GetUserById(c *gin.Context)
	AddUserData(c *gin.Context)
	UpdateUserData(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type UserServiceImpl struct {
	userRepository repository.UserRepository
}

func (u UserServiceImpl) GetMe(c *gin.Context) {
	userIdAny := c.MustGet("Id").(string)

	userId, err := uuid.Parse(userIdAny)
	if err != nil {
		log.Println("Failed to parse user ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}
	user, err := u.userRepository.FindUserById(userId)
	if err != nil {
		log.Println("Happened error when get data from database. Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (u UserServiceImpl) UpdateUserData(c *gin.Context) {
	defer pkg.PanicHandler(c)

	log.Info("start to execute program update user data by id")
	userID := c.Param("userID")

	var request dao.User
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicException(constant.InvalidRequest)
	}

	parsedId, err := uuid.Parse(userID)
	if err != nil {
		log.Error("Happened error when parse user id. Error", err)
		pkg.PanicException(constant.InvalidRequest)
	}

	data, err := u.userRepository.FindUserById(parsedId)
	if err != nil {
		log.Error("Happened error when get data from database. Error", err)
		pkg.PanicException(constant.DataNotFound)
	}

	data.Email = request.Email
	data.Name = request.Name
	data.Status = request.Status
	u.userRepository.Save(&data)

	if err != nil {
		log.Error("Happened error when updating data to database. Error", err)
		pkg.PanicException(constant.UnknownError)
	}

	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, data))
}

func (u UserServiceImpl) GetUserById(c *gin.Context) {
	defer pkg.PanicHandler(c)

	log.Info("start to execute program get user by id")
	userID := c.Param("userID")
	parsedId, err := uuid.Parse(userID)
	if err != nil {
		log.Error("Happened error when parse user id. Error", err)
		pkg.PanicException(constant.InvalidRequest)
	}
	data, err := u.userRepository.FindUserById(parsedId)
	if err != nil {
		log.Error("Happened error when get data from database. Error", err)
		pkg.PanicException(constant.DataNotFound)
	}
	result := dao.UserResponse{
		ID:    data.ID,
		Name:  data.Name,
		Email: data.Email,
	}
	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, result))
}

func (u UserServiceImpl) AddUserData(c *gin.Context) {
	defer pkg.PanicHandler(c)

	if u.userRepository == nil {
		log.Error("userRepository is nil")
		pkg.PanicException(constant.UnknownError)
	}

	log.Info("start to execute program add data user")
	var request dao.User
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicException(constant.InvalidRequest)
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := dao.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: hashedPassword,
		Status:   request.Status,
	}
	data, err := u.userRepository.Save(&user)
	if err != nil {
		log.Error("Happened error when saving data to database. Error", err)
		pkg.PanicException(constant.UnknownError)
	}

	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, data))
}

func (u UserServiceImpl) GetAllUser(c *gin.Context) {
	defer pkg.PanicHandler(c)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	log.Info("start to execute get all data user")

	data, err := u.userRepository.FindAllUser(limit, offset)
	if err != nil {
		log.Error("Happened Error when find all user data. Error: ", err)
		pkg.PanicException(constant.UnknownError)
	}

	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, data))
}

func (u UserServiceImpl) DeleteUser(c *gin.Context) {
	defer pkg.PanicHandler(c)

	log.Info("start to execute delete data user by id")
	userID := c.Param("userID")
	parsedId, err := uuid.Parse(userID)
	if err != nil {
		log.Error("Happened error when parse user id. Error", err)
		pkg.PanicException(constant.InvalidRequest)
	}
	errRepo := u.userRepository.DeleteUserById(parsedId)
	if errRepo != nil {
		log.Error("Happened Error when try delete data user from DB. Error:", err)
		pkg.PanicException(constant.UnknownError)
	}

	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, pkg.Null()))
}

func UserServiceInit(userRepository repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		userRepository: userRepository,
	}
}
