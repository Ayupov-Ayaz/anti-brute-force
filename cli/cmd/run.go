package run

import (
	"fmt"

	"github.com/spf13/cobra"

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

var (
	port    int
	useGRPC bool
	runCmd  = &cobra.Command{
		Use:   "run -p [port] -g [use grpc]",
		Short: "run server",
		RunE:  run,
		Long: `Run HTTP or GRPC server.
Example: ./anti-brute-force run -p 8080
Example: ./anti-brute-force run -p 8080 -g true`,
	}
)

func init() {
	runCmd.Flags().IntVarP(&port, "port", "p", 8080, "port")
	runCmd.Flags().BoolVarP(&useGRPC, "use grpc", "g", false, "use grpc")
}

func Execute() error {
	return runCmd.Execute()
}

func run(_ *cobra.Command, _ []string) error {
	cfg, err := config.ParseConfig(port, useGRPC)
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	zLogger, err := logger.New(cfg.Logger)
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
		manager.WithLogger(zLogger))

	ipChecker := checker.New(
		checker.WithBuckets(ipBuckets),
		checker.WithWhiteList(whiteList),
		checker.WithBlackList(blackList),
		checker.WithLogger(zLogger))

	port := cfg.Server.Port
	fmt.Println(cfg.Server)

	if useGRPC {
		server := grpcserver.New(
			grpcserver.WithManager(ipManager),
			grpcserver.WithChecker(ipChecker))

		if err := server.Start(port); err != nil {
			return fmt.Errorf("grpc server: %w", err)
		}
	} else {
		server := httpserver.New(
			httpserver.WithChecker(ipChecker),
			httpserver.WithManager(ipManager))
		if err := server.Start(port); err != nil {
			return fmt.Errorf("http server: %w", err)
		}
	}

	return nil
}
