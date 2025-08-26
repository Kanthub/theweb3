package grpc

import (
	"context"
	"fmt"
	"github.com/the-web3/rpc-server/database"
	"net"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/the-web3/rpc-server/proto/market"
)

const MaxRecvMessageSize = 1024 * 1024 * 30000

type MarketRpcConfig struct {
	Host string
	Port int
}

type MarketRpcService struct {
	db *database.DB
	*MarketRpcConfig
	market.UnimplementedMarketPriceServiceServer
	stopped atomic.Bool
}

func NewMarketRpcService(db *database.DB, conf *MarketRpcConfig) (*MarketRpcService, error) {
	return &MarketRpcService{
		db:              db,
		MarketRpcConfig: conf,
	}, nil
}

func (mps *MarketRpcService) Start(ctx context.Context) error {
	go func(ms *MarketRpcService) {
		rpcAddr := fmt.Sprintf("%s:%d", ms.MarketRpcConfig.Host, ms.MarketRpcConfig.Port)
		listener, err := net.Listen("tcp", rpcAddr)
		if err != nil {
			log.Error("Could not start tcp listener. ")
		}

		opt := grpc.MaxRecvMsgSize(MaxRecvMessageSize)

		gs := grpc.NewServer(
			opt,
			grpc.ChainUnaryInterceptor(
				nil,
			),
		)

		reflection.Register(gs)
		market.RegisterMarketPriceServiceServer(gs, ms)

		log.Info("grpc info", "addr", listener.Addr())

		if err := gs.Serve(listener); err != nil {
			log.Error("start rpc server fail", "err", err)
		}
	}(mps)
	return nil
}

func (mps *MarketRpcService) Stop(ctx context.Context) error {
	mps.stopped.Store(true)
	return nil
}

func (mps *MarketRpcService) Stopped() bool {
	return mps.stopped.Load()
}
