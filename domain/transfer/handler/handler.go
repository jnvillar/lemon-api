package transferhandler

import (
	"context"

	transferconverter "lemonapp/domain/transfer/converter"
	transferservice "lemonapp/domain/transfer/service"
	trasnfervalidator "lemonapp/domain/transfer/validator"
	usersservice "lemonapp/domain/user/service"
	walletmodel "lemonapp/domain/wallet/model"
	walletservice "lemonapp/domain/wallet/service"
	"lemonapp/logger"
	"lemonapp/utils"

	"github.com/gin-gonic/gin"
)

type TransferHandler struct {
	transferService   transferservice.Service
	transferConverter *transferconverter.Converter
	walletService     walletservice.Service
	usersService      usersservice.Service
	transferValidator *trasnfervalidator.Validator
}

func NewTransferHandler(
	userService usersservice.Service,
	walletService walletservice.Service,
	transferService transferservice.Service,
	transferConverter *transferconverter.Converter,
	transferValidator *trasnfervalidator.Validator) *TransferHandler {
	return &TransferHandler{
		walletService:     walletService,
		usersService:      userService,
		transferService:   transferService,
		transferConverter: transferConverter,
		transferValidator: transferValidator,
	}
}

func (b *TransferHandler) RegisterRoutes(router *gin.RouterGroup) {
	walletApi := router.Group("/users/:user_id/wallets/:wallet_id")
	walletApi.POST("/transfer", func(c *gin.Context) { b.transfer(c) })
	walletApi.POST("/deposit", func(c *gin.Context) { b.deposit(c) })
	walletApi.POST("/extraction", func(c *gin.Context) { b.extraction(c) })
}

func (b *TransferHandler) deposit(c *gin.Context) {
	transfer, err := b.transferConverter.ConvertTransferFromRequest(c.Request)
	if err != nil {
		logger.Sugar.Errorf("error parsing trasnfer: %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	if err := b.transferValidator.ValidateDeposit(transfer); err != nil {
		logger.Sugar.Errorf("invalid trasnfer: %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	user, err := b.usersService.GetByID(c, c.Param("user_id"))
	if err != nil {
		logger.Sugar.Errorf("error fetching user: %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	wallet, err := b.walletService.GetUserWallet(c, user.ID, c.Param("wallet_id"))
	if err != nil {
		logger.Sugar.Errorf("error fetching wallets: %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	transfer.Currency = wallet.GetBaseWallet().Currency

	err = b.transferService.Deposit(c, transfer, wallet)
	if err != nil {
		logger.Sugar.Errorf("error doing deposit: %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	utils.ReturnJSONResponse(c.Writer, transfer)
}

func (b *TransferHandler) extraction(c *gin.Context) {
	transfer, err := b.transferConverter.ConvertTransferFromRequest(c.Request)
	if err != nil {
		logger.Sugar.Errorf("error parsing trasnfer: %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	if err := b.transferValidator.ValidateDeposit(transfer); err != nil {
		logger.Sugar.Errorf("invalid trasnfer: %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	user, err := b.usersService.GetByID(c, c.Param("user_id"))
	if err != nil {
		logger.Sugar.Errorf("error fetching user: %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	wallet, err := b.walletService.GetUserWallet(c, user.ID, c.Param("wallet_id"))
	if err != nil {
		logger.Sugar.Errorf("error fetching wallets: %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	transfer.Currency = wallet.GetBaseWallet().Currency

	err = b.transferService.Extraction(c, transfer, wallet)
	if err != nil {
		logger.Sugar.Errorf("error doing deposit: %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	utils.ReturnJSONResponse(c.Writer, transfer)
}

func (b *TransferHandler) getTransferWallets(ctx context.Context, userID, walletFrom, walletTo string) (walletmodel.Wallet, walletmodel.Wallet, error) {
	wf, err := b.walletService.GetUserWallet(ctx, userID, walletFrom)
	if err != nil {
		return nil, nil, err
	}

	wt, err := b.walletService.GetByID(ctx, walletTo)
	if err != nil {
		return nil, nil, err
	}

	return wf, wt, nil
}

func (b *TransferHandler) transfer(c *gin.Context) {
	transfer, err := b.transferConverter.ConvertTransferFromRequest(c.Request)
	if err != nil {
		logger.Sugar.Errorf("error parsing trasnfer: %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	if err := b.transferValidator.ValidateTransfer(transfer); err != nil {
		logger.Sugar.Errorf("invalid trasnfer: %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	user, err := b.usersService.GetByID(c, c.Param("user_id"))
	if err != nil {
		logger.Sugar.Errorf("error fetching user: %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	walletFrom, walletTo, err := b.getTransferWallets(c, user.ID, c.Param("wallet_id"), transfer.WalletTo)
	if err != nil {
		logger.Sugar.Errorf("error fetching wallets: %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	transfer.WalletTo = walletFrom.GetBaseWallet().ID
	transfer.WalletFrom = walletTo.GetBaseWallet().ID
	transfer.Currency = walletFrom.GetBaseWallet().Currency

	err = b.transferService.Transfer(c, transfer, walletFrom, walletTo)
	if err != nil {
		logger.Sugar.Errorf("error doing trasnfer: %v", err)
		utils.ReturnError(c.Writer, err)
		return
	}

	utils.ReturnJSONResponse(c.Writer, transfer)
}
