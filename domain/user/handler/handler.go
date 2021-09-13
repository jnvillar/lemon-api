package userhandler

import (
	userconverter "lemonapp/domain/user/converter"
	usermodel "lemonapp/domain/user/model"
	usersservice "lemonapp/domain/user/service"
	uservalidator "lemonapp/domain/user/validator"
	walletservice "lemonapp/domain/wallet/service"
	"lemonapp/logger"
	"lemonapp/utils"

	"github.com/gin-gonic/gin"
)

type UsersHandler struct {
	usersService   usersservice.Service
	usersConverter *userconverter.Converter
	usersValidator *uservalidator.Validator
	walletService  walletservice.Service
}

func NewUsersHandler(
	usersService usersservice.Service,
	usersConverter *userconverter.Converter,
	usersValidator *uservalidator.Validator,
	walletService walletservice.Service, ) *UsersHandler {
	return &UsersHandler{
		usersService:   usersService,
		walletService:  walletService,
		usersConverter: usersConverter,
		usersValidator: usersValidator,
	}
}

func (b *UsersHandler) RegisterRoutes(router *gin.RouterGroup) {
	usersApi := router.Group("/users")
	usersApi.POST("", func(c *gin.Context) { b.createUser(c) })
	usersApi.GET("/:user_id", func(c *gin.Context) { b.getUser(c) })
	usersApi.GET("", func(c *gin.Context) { b.listUsers(c) })
}

func (b *UsersHandler) listUsers(c *gin.Context) {
	users, err := b.usersService.List(c)
	if err != nil {
		logger.Sugar.Errorf("error listing users %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	utils.ReturnJSONResponse(c.Writer, users)
}

func (b *UsersHandler) getUser(c *gin.Context) {
	user, err := b.usersService.GetByID(c, c.Param("user_id"))
	if err != nil {
		logger.Sugar.Errorf("error fetching user: %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	wallets, err := b.walletService.GetUserWallets(c, user.ID)
	if err != nil {
		logger.Sugar.Errorf("error fetching user wallets: %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	response := usermodel.NewUserResponse(user, wallets)

	utils.ReturnJSONResponse(c.Writer, response)
}

func (b *UsersHandler) createUser(c *gin.Context) {
	user, err := b.usersConverter.ConvertToUserFromRequest(c.Request)
	if err != nil {
		logger.Sugar.Errorf("error parsing user %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	if err := b.usersValidator.ValidateUser(user); err != nil {
		logger.Sugar.Errorf("error validating user %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	user, err = b.usersService.Create(c, usermodel.NewUser(user.FirstName, user.LastName, user.Alias, user.Email))
	if err != nil {
		logger.Sugar.Errorf("error creating user %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	_, err = b.walletService.CreateWallets(c, user.ID)
	if err != nil {
		logger.Sugar.Errorf("error creating user's wallets %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	utils.ReturnJSONResponse(c.Writer, user)
}
