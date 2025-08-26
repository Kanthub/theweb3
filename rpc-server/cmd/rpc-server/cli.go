package main

import (
	"context"
	"github.com/the-web3/rpc-server/services/rest"

	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"

	"github.com/the-web3/rpc-server/common/cliapp"
	"github.com/the-web3/rpc-server/common/opio"
	"github.com/the-web3/rpc-server/config"
	"github.com/the-web3/rpc-server/database"
	flags2 "github.com/the-web3/rpc-server/flags"
	"github.com/the-web3/rpc-server/services/grpc"
	"github.com/the-web3/rpc-server/tasker"
)

func runTask(ctx *cli.Context, shutdown context.CancelCauseFunc) (cliapp.Lifecycle, error) {
	log.Info("run market price task")
	cfg := config.NewConfig(ctx)
	db, err := database.NewDB(ctx.Context, cfg.MasterDB)
	if err != nil {
		log.Error("failed to connect to database", "err", err)
		return nil, err
	}
	return tasker.NewMarketPriceTasker(context.Background(), db, &cfg, shutdown)
}

func runRpc(ctx *cli.Context, shutdown context.CancelCauseFunc) (cliapp.Lifecycle, error) {
	log.Info("start grpc service")
	cfg := config.NewConfig(ctx)

	grpcConf := &grpc.MarketRpcConfig{
		Host: cfg.RpcServer.Host,
		Port: cfg.RpcServer.Port,
	}

	db, err := database.NewDB(ctx.Context, cfg.SlaveDB)
	if err != nil {
		log.Error("failed to connect to database", "err", err)
		return nil, err
	}

	return grpc.NewMarketRpcService(db, grpcConf)
}

func runMigrations(ctx *cli.Context) error {
	ctx.Context = opio.CancelOnInterrupt(ctx.Context)
	log.Info("running migrations...")
	cfg := config.NewConfig(ctx)
	db, err := database.NewDB(ctx.Context, cfg.MasterDB)
	if err != nil {
		log.Error("failed to connect to database", "err", err)
		return err
	}
	defer func(db *database.DB) {
		err := db.Close()
		if err != nil {
			log.Error("fail to close database", "err", err)
		}
	}(db)
	return db.ExecuteSQLMigration(cfg.Migrations)
}

func RunConfigSymbols(ctx *cli.Context) error {
	log.Info("init system symbols")
	cfg := config.NewConfig(ctx)
	db, err := database.NewDB(ctx.Context, cfg.MasterDB)
	if err != nil {
		log.Error("failed to connect to database", "err", err)
		return err
	}
	assetExchange, err := tasker.NewMarketAssetAndExchange(db)
	if err != nil {
		log.Error("new market asset and exchange fail", "err", err)
		return err
	}
	return assetExchange.ConfigAssetAndExchange()
}

func runRestApi(ctx *cli.Context, shutdown context.CancelCauseFunc) (cliapp.Lifecycle, error) {
	log.Info("running api...")
	cfg := config.NewConfig(ctx)
	return rest.NewApi(ctx.Context, &cfg)
}

func NewCli() *cli.App {
	flags := flags2.Flags
	return &cli.App{
		Version:              "0.0.1",
		Description:          "An market services with rpc",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:        "task",
				Flags:       flags,
				Description: "Run market services task",
				Action:      cliapp.LifecycleCmd(runTask),
			},
			{
				Name:        "rpc",
				Flags:       flags,
				Description: "Run rpc services",
				Action:      cliapp.LifecycleCmd(runRpc),
			},
			{
				Name:        "migrate",
				Flags:       flags,
				Description: "Run database migrations",
				Action:      runMigrations,
			},
			{
				Name:        "init-symbol",
				Flags:       flags,
				Description: "init symbol in database",
				Action:      RunConfigSymbols,
			},
			{
				Name:        "api",
				Flags:       flags,
				Description: "run rest api",
				Action:      cliapp.LifecycleCmd(runRestApi),
			},
			{
				Name:        "version",
				Description: "Show project version",
				Action: func(ctx *cli.Context) error {
					cli.ShowVersion(ctx)
					return nil
				},
			},
		},
	}
}
