package server

import (
	"database/sql"
	"fmt"

	transactionconverter "lemonapp/domain/transaction/converter"
	"lemonapp/logger"

	_ "github.com/go-sql-driver/mysql"

	healthHandler "lemonapp/domain/health/handler"
	transactionhandler "lemonapp/domain/transaction/handler"
	transactionrepository "lemonapp/domain/transaction/repository"
	transactionservice "lemonapp/domain/transaction/service"
	transferconverter "lemonapp/domain/transfer/converter"
	transferhandler "lemonapp/domain/transfer/handler"
	transferrepository "lemonapp/domain/transfer/repository"
	transferservice "lemonapp/domain/transfer/service"
	trasnfervalidator "lemonapp/domain/transfer/validator"
	userconverter "lemonapp/domain/user/converter"
	userhandler "lemonapp/domain/user/handler"
	userrepository "lemonapp/domain/user/repository"
	userservice "lemonapp/domain/user/service"
	uservalidator "lemonapp/domain/user/validator"
	wallethandler "lemonapp/domain/wallet/handler"
	walletrepository "lemonapp/domain/wallet/repository"
	walletservice "lemonapp/domain/wallet/service"
	"lemonapp/handler"

	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

type Server struct {
	usersService         userservice.Service
	usersConverter       *userconverter.Converter
	usersValidator       *uservalidator.Validator
	walletService        walletservice.Service
	transactionsService  transactionservice.Service
	transactionConverter *transactionconverter.Converter
	transferService      transferservice.Service
	transferConverter    *transferconverter.Converter
	transferValidator    *trasnfervalidator.Validator
}

type Runnable struct{}

func NewRunnable() *Runnable {
	return &Runnable{}
}

func (r *Runnable) Cmd() *cobra.Command {
	cmd.Run = func(_ *cobra.Command, _ []string) {
		server := r.Run(&Options{
			LogLevel:                logLevelArg,
			SQLDatabase:             sqlDatabase,
			SQLDatabaseMaxOpenConns: sqlDatabaseMaxOpenConns,
			SQLUser:                 sqlUser,
			SQLPassword:             sqlPassword,
		})
		server.Start()
	}
	return cmd
}

func (r *Runnable) Run(options *Options) *Server {
	return newServer(options)
}

func newServer(options *Options) *Server {
	err := logger.Initialize(options.LogLevel, "1.0.0")
	if err != nil {
		panic(err)
	}

	db := mustGetDb(options)

	walletRepository := walletrepository.NewRepository(db)
	walletService := walletservice.NewServiceImpl(walletRepository)
	usersRepository := userrepository.NewRepository(db)
	usersService := userservice.NewServiceImpl(usersRepository)
	transactionsRepository := transactionrepository.NewRepository(db)
	transactionsService := transactionservice.NewServiceImpl(transactionsRepository)
	transferRepository := transferrepository.NewRepository(db, walletRepository, transactionsRepository)
	transferService := transferservice.NewServiceImpl(transferRepository)

	app := &Server{
		usersService:         usersService,
		usersConverter:       userconverter.NewUsersConverter(),
		usersValidator:       uservalidator.NewUserValidator(),
		walletService:        walletService,
		transactionsService:  transactionsService,
		transactionConverter: transactionconverter.NewTransactionConverter(),
		transferService:      transferService,
		transferValidator:    trasnfervalidator.NewTransferValidator(),
		transferConverter:    transferconverter.NewTransferConverter(),
	}
	return app
}

func (s *Server) Start() {
	router := gin.Default()
	s.registerRoutes(router)
}

func (s *Server) registerRoutes(router *gin.Engine) {

	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowedMethods:  []string{"*"},
		AllowedHeaders:  []string{"*"},
	}))

	api := router.Group("/api")

	appHandlers := []handler.Handler{
		healthHandler.NewHealthHandler(),
		userhandler.NewUsersHandler(s.usersService, s.usersConverter, s.usersValidator, s.walletService),
		wallethandler.NewWalletHandler(s.usersService, s.walletService),
		transactionhandler.NewTransactionsHandler(s.usersService, s.walletService, s.transactionsService, s.transactionConverter),
		transferhandler.NewTransferHandler(s.usersService, s.walletService, s.transferService, s.transferConverter, s.transferValidator),
	}

	for _, h := range appHandlers {
		h.RegisterRoutes(api)
	}

	if err := router.Run(); err != nil {
		panic(err)
	}
}

func mustGetDb(options *Options) *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s?parseTime=true", options.SQLUser, options.SQLPassword, options.SQLDatabase))
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(options.SQLDatabaseMaxOpenConns)
	return db
}
