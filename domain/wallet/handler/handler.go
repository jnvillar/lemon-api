package wallterhandler

import (
	usersservice "lemonapp/domain/user/service"
	walletservice "lemonapp/domain/wallet/service"
	"lemonapp/logger"
	"lemonapp/utils"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	walletService walletservice.Service
	usersService  usersservice.Service
}

func NewWalletHandler(userService usersservice.Service, walletService walletservice.Service) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
		usersService:  userService,
	}
}

func (b *WalletHandler) RegisterRoutes(router *gin.RouterGroup) {
	walletApi := router.Group("/users/:user_id/wallets")
	walletApi.GET("", func(c *gin.Context) { b.getUserWallets(c) })
}

func (b *WalletHandler) getUserWallets(c *gin.Context) {
	user, err := b.usersService.GetByID(c, c.Param("user_id"))
	if err != nil {
		logger.Sugar.Errorf("error fetching user %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	wallets, err := b.walletService.GetUserWallets(c, user.ID)
	if err != nil {
		logger.Sugar.Errorf("error creating user %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	utils.ReturnJSONResponse(c.Writer, wallets)
}
