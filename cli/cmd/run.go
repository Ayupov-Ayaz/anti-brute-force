package run

import (
	"fmt"

	redisstorage "github.com/ayupov-ayaz/anti-brute-force/internal/modules/storage/redis"

	redisstore "github.com/ayupov-ayaz/anti-brute-force/internal/modules/db/redis"

	grpcserver "github.com/ayupov-ayaz/anti-brute-force/internal/server/grpc"

	"github.com/ayupov-ayaz/anti-brute-force/config"
	"github.com/ayupov-ayaz/anti-brute-force/internal/app/checker"
	"github.com/ayupov-ayaz/anti-brute-force/internal/app/manager"
	"github.com/ayupov-ayaz/anti-brute-force/internal/modules/buckets"
	"github.com/ayupov-ayaz/anti-brute-force/internal/modules/iplist"
	httpserver "github.com/ayupov-ayaz/anti-brute-force/internal/server/http"
)

func Run() error {
	cfg := config.ParseConfig()

	redisClient, err := redisstore.NewRedisClient(cfg.Redis)
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
		manager.WithWhiteList(whiteList))

	ipChecker := checker.New(
		checker.WithBuckets(ipBuckets),
		checker.WithBlackList(whiteList),
		checker.WithBlackList(blackList))

	if cfg.UseGRPC() {
		server := grpcserver.New(
			grpcserver.WithManager(ipManager),
			grpcserver.WithChecker(ipChecker))

		if err := server.Start(cfg.GRPC.Port); err != nil {
			return fmt.Errorf("grpc server: %w", err)
		}
	}

	if cfg.UseHTTP() {
		server := httpserver.New(
			httpserver.WithChecker(ipChecker),
			httpserver.WithManager(ipManager))
		if err := server.Start(cfg.HTTP.Port); err != nil {
			return fmt.Errorf("http server: %w", err)
		}
	}

	return nil
}
