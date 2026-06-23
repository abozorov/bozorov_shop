package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abozorov/bozorov_shop/internal/config"
	"github.com/abozorov/bozorov_shop/internal/handlers"
	orderhandler "github.com/abozorov/bozorov_shop/internal/handlers/order"
	userhandler "github.com/abozorov/bozorov_shop/internal/handlers/user"
	orderrepo "github.com/abozorov/bozorov_shop/internal/repo/order"
	userrepo "github.com/abozorov/bozorov_shop/internal/repo/user"
	orderservice "github.com/abozorov/bozorov_shop/internal/service/order"
	userservice "github.com/abozorov/bozorov_shop/internal/service/user"
	"github.com/abozorov/bozorov_shop/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	// load logger
	logger, err := logger.NewLogger(true)
	if err != nil {
		log.Fatal("Eror creating logger %w", err)
	}

	// load config
	cfg, err := config.NewConfig("internal/config/config.env")
	if err != nil {
		logger.Fatal("Error config load %w", zap.String("error:", err.Error()))
	}

	// connectiong db
	db, err := pgxpool.New(context.Background(), fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	))

	if err != nil {
		logger.Fatal("Error config load", zap.String("error:", err.Error()))
	}
	// create layers
	userRepo := userrepo.NewUserRepo(db)
	orderRepo := orderrepo.NewOrderRepo(db)

	userService := userservice.NewUserService(userRepo, orderRepo)
	orderService := orderservice.NewOrderService(userRepo, orderRepo)

	userHandlers := userhandler.NewUserHandler(userService, logger)
	orderHandlers := orderhandler.NewOrderHandler(orderService, logger)

	// create router
	router := handlers.NewRouter(userHandlers, orderHandlers)
	server := &http.Server{
		Addr:    cfg.HttpHost,
		Handler: router,
	}

	// statrt server
	go func() {
		logger.Info(fmt.Sprintf("Server started localhost:%s started", server.Addr))
		err := server.ListenAndServe()
		if err != nil {
			logger.Error("main", zap.Error(err))
			return
		}
	}()

	// gracefull shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	logger.Info("Shutdown server started")
	stopCtx, stopCancle := context.WithTimeout(context.Background(), time.Second*5)
	defer stopCancle()

	server.Shutdown(stopCtx)

	logger.Info("Server shutdown completed")

}
