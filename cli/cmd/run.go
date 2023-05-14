package run

import (
	"fmt"

	"github.com/ayupov-ayaz/anti-brute-force/internal/server/http/handlers"

	"github.com/ayupov-ayaz/anti-brute-force/internal/modules/logger"

	redisstorage "github.com/ayupov-ayaz/anti-brute-force/internal/modules/storage/redis"

	redissdb "github.com/ayupov-ayaz/anti-brute-force/internal/modules/db/redis"

	grpcserver "github.com/ayupov-ayaz/anti-brute-force/internal/server/grpc"

	"github.com/ayupov-ayaz/anti-brute-force/config"
	"github.com/ayupov-ayaz/anti-brute-force/internal/app/checker"
	"github.com/ayupov-ayaz/anti-brute-force/internal/app/manager"
	"github.com/ayupov-ayaz/anti-brute-force/internal/modules/buckets"
	"github.com/ayupov-ayaz/anti-brute-force/internal/modules/iplist"
	httpserver "github.com/ayupov-ayaz/anti-brute-force/internal/server/http"
)

func Run() error {
	cfg, err := config.ParseConfig()
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	logger, err := logger.New(cfg.Logger)
	if err != nil {
		return fmt.Errorf("logger: %w", err)
	}

	redisClient, err := redissdb.NewRedisClient(cfg.Redis)
	if err != nil {
		return fmt.Errorf("redis client: %w", err)
	}

	storage := redisstorage.New(redisClient)

	blackList := iplist.New(cfg.IPList.BlackListAddr, storage)
	whiteList := iplist.New(cfg.IPList.WhiteListAddr, storage)
	ipBuckets := buckets.New()

	ipManager := manager.New(
		manager.WithResetter(ipBuckets),
		manager.WithBlackList(blackList),
		manager.WithWhiteList(whiteList),
		manager.WithLogger(logger))

	ipChecker := checker.New(
		checker.WithBuckets(ipBuckets),
		checker.WithWhiteList(whiteList),
		checker.WithBlackList(blackList),
		checker.WithLogger(logger))

	if cfg.UseGRPC() {
		server := grpcserver.New(
			grpcserver.WithManager(ipManager),
			grpcserver.WithChecker(ipChecker))

		if err := server.Start(cfg.GRPC.Port); err != nil {
			return fmt.Errorf("grpc server: %w", err)
		}
	}

	if cfg.UseHTTP() {
		http := httpserver.New(
			httpserver.WithChecker(handlers.NewChecker(ipChecker)),
			httpserver.WithManager(handlers.NewManager(ipManager, logger)))

		if err := http.Start(httpserver.NewFiber(), cfg.HTTP.Port); err != nil {
			return fmt.Errorf("http server: %w", err)
		}
	}

	return nil
}
