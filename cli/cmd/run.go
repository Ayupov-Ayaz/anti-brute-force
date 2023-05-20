package run

import (
	"fmt"

	"github.com/ayupov-ayaz/anti-brute-force/internal/modules/validator"

	"github.com/ayupov-ayaz/anti-brute-force/internal/app/checker"

	"github.com/ayupov-ayaz/anti-brute-force/cli/cmd/internal"

	"github.com/spf13/cobra"

	"github.com/ayupov-ayaz/anti-brute-force/internal/server/http/handlers"

	"github.com/ayupov-ayaz/anti-brute-force/internal/modules/logger"

	redisstorage "github.com/ayupov-ayaz/anti-brute-force/internal/modules/storage/redis"

	redissdb "github.com/ayupov-ayaz/anti-brute-force/internal/modules/db/redis"

	"github.com/ayupov-ayaz/anti-brute-force/config"
	"github.com/ayupov-ayaz/anti-brute-force/internal/app/manager"
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
Example: ./anti-brute-force run -p 8080 -g`,
	}
)

func init() {
	runCmd.Flags().IntVarP(&port, "port", "p", 0, "port")
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

	valid := validator.New()
	if err := valid.Validate(cfg); err != nil {
		return fmt.Errorf("validate config: %w", err)
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
	authLimiter := internal.NewAuthRateLimiter(cfg.Limiter, redisClient, zLogger)

	managerHandler := handlers.NewManager(manager.New(whiteList, blackList, authLimiter), valid, zLogger)
	checkerHandler := handlers.NewChecker(checker.New(whiteList, blackList, authLimiter), valid, zLogger)

	server := httpserver.New(managerHandler, checkerHandler, zLogger)

	if err := server.Start(httpserver.NewFiber(), cfg.Server.Port); err != nil {
		return fmt.Errorf("http server: %w", err)
	}

	return nil
}
