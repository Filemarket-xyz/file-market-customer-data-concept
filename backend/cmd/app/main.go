package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/config"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/handler"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/repository"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/repository/postgres"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/repository/redis"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/server"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/service"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/ethclient/rpc"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/hash"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/jwtoken"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/logger"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/providers/dune"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/rand_manager"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/time_manager"
	rpc2 "github.com/ethereum/go-ethereum/rpc"
)

func main() {
	logging, err := logger.NewLogger()
	if err != nil {
		log.Panic(err)
	}
	defer logging.Sync()

	var cfgPath string

	flag.StringVar(&cfgPath, "cfg", "", "")
	flag.Parse()

	cfg, err := config.Init(cfgPath)
	if err != nil {
		logging.Panic(err)
	}
	ctx := context.Background()

	pool, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		logging.Panic(err)
	}

	redis := redis.New(ctx, cfg.Redis)

	rpcClient, err := rpc2.DialOptions(ctx, cfg.Service.RpcUrl)
	if err != nil {
		logging.Panic(err)
	}
	ethClient := rpc.NewEthClient(rpcClient)

	jwtokenManager := jwtoken.NewTokenManager(cfg.TokenManager.SigningKey)
	timeManager := time_manager.New(cfg.Locale)
	randManager := rand_manager.NewRandManager(time.Now().UnixNano())

	duneProvider := dune.NewDune(cfg.Dune, logging)

	repo, err := repository.NewRepository(cfg, redis, pool)
	if err != nil {
		logging.Panic(err)
	}

	service, err := service.NewService(ctx, repo, ethClient, jwtokenManager, hash.NewHashManager(), timeManager, randManager, duneProvider, cfg.Service, logging)
	if err != nil {
		logging.Panic(err)
	}

	router := handler.NewHandler(cfg.Handler, service, logging)

	srv := server.NewServer(cfg.Server, router.Init())

	go func() {
		if err = srv.ListenAndServe(); err != http.ErrServerClosed {
			logging.Panic(err)
		}
	}()
	logging.Infof("server listening on port %d\n", cfg.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	if err = srv.Shutdown(ctx); err != nil {
		logging.Panic(err)
	}

	service.Shutdown()
	if err := redis.Close(); err != nil {
		logging.Panic(err)
	}

}
