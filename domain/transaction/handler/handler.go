package handler

import (
	transactionconverter "lemonapp/domain/transaction/converter"
	transactionservice "lemonapp/domain/transaction/service"
	usersservice "lemonapp/domain/user/service"
	walletService "lemonapp/domain/wallet/service"
	"lemonapp/logger"
	"lemonapp/utils"

	"github.com/gin-gonic/gin"
)

type TransactionsHandler struct {
	transactionService   transactionservice.Service
	transactionConverter *transactionconverter.Converter
	walletService        walletService.Service
	usersService         usersservice.Service
}

func NewTransactionsHandler(
	usersService usersservice.Service,
	walletService walletService.Service,
	transactionService transactionservice.Service,
	transactionConverter *transactionconverter.Converter) *TransactionsHandler {
	return &TransactionsHandler{
		usersService:         usersService,
		walletService:        walletService,
		transactionService:   transactionService,
		transactionConverter: transactionConverter,
	}
}

func (b *TransactionsHandler) RegisterRoutes(router *gin.RouterGroup) {
	walletTransactionsApi := router.Group("/users/:user_id/wallets/:wallet_id/transactions")
	walletTransactionsApi.GET("", func(c *gin.Context) { b.getWalletTransactions(c) })

	userTransactionsApi := router.Group("/users/:user_id/transactions")
	userTransactionsApi.GET("", func(c *gin.Context) { b.getUserTransactions(c) })
}

func (b *TransactionsHandler) getUserTransactions(c *gin.Context) {
	user, err := b.usersService.GetByID(c, c.Param("user_id"))
	if err != nil {
		logger.Sugar.Errorf("error fetching user: %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	searchParams, err := b.transactionConverter.ConvertSearchParamsFromRequest(c.Request)
	if err != nil {
		logger.Sugar.Errorf("error parsing search params: %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}
	searchParams.UserID = user.ID

	transactions, err := b.transactionService.Search(c, searchParams)
	if err != nil {
		logger.Sugar.Errorf("error fetching user transactions: %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	utils.ReturnJSONResponse(c.Writer, transactions)
}

func (b *TransactionsHandler) getWalletTransactions(c *gin.Context) {
	user, err := b.usersService.GetByID(c, c.Param("user_id"))
	if err != nil {
		logger.Sugar.Errorf("error fetching user: %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	wallet, err := b.walletService.GetUserWallet(c, user.ID, c.Param("wallet_id"))
	if err != nil {
		logger.Sugar.Errorf("error fetching wallet: %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	searchParams, err := b.transactionConverter.ConvertSearchParamsFromRequest(c.Request)
	if err != nil {
		logger.Sugar.Errorf("error parsing search params: %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	searchParams.UserID = user.ID
	searchParams.WalletID = wallet.GetBaseWallet().ID

	transactions, err := b.transactionService.Search(c, searchParams)
	if err != nil {
		logger.Sugar.Errorf("error fetching wallet transaction: %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	utils.ReturnJSONResponse(c.Writer, transactions)
}
