package run

import (
	"fmt"

	grpcserver "github.com/ayupov-ayaz/anti-bute-force/internal/server/grpc"

	"github.com/ayupov-ayaz/anti-bute-force/config"
	"github.com/ayupov-ayaz/anti-bute-force/internal/app/checker"
	"github.com/ayupov-ayaz/anti-bute-force/internal/app/manager"
	"github.com/ayupov-ayaz/anti-bute-force/internal/modules/buckets"
	"github.com/ayupov-ayaz/anti-bute-force/internal/modules/iplist"
	httpserver "github.com/ayupov-ayaz/anti-bute-force/internal/server/http"
)

func Run() error {
	cfg := config.ParseConfig()

	blackList := iplist.New(cfg.IPList.BlackListAddr)
	whiteList := iplist.New(cfg.IPList.WhiteListAddr)
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
