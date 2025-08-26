package chaindispatcher

import (
	"context"
	"runtime/debug"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ethereum/go-ethereum/log"

	"github.com/the-web3/wallet-sign-s6/chain"
	"github.com/the-web3/wallet-sign-s6/chain/bitcoin"
	"github.com/the-web3/wallet-sign-s6/chain/ethereum"
	"github.com/the-web3/wallet-sign-s6/chain/solana"
	"github.com/the-web3/wallet-sign-s6/config"
	"github.com/the-web3/wallet-sign-s6/hsm"
	"github.com/the-web3/wallet-sign-s6/leveldb"
	"github.com/the-web3/wallet-sign-s6/protobuf/wallet"
)

const AccessToken = "DappLinkTheWeb3202402290001"

type CommonRequest interface {
	GetConsumerToken() string
	GetChainName() string
}

type ChainType = string

type CommonReply = wallet.ChainSignMethodResponse

type ChainDispatcher struct {
	registry map[string]chain.IChainAdaptor
}

func NewChainDispatcher(conf *config.Config) (*ChainDispatcher, error) {
	dispatcher := ChainDispatcher{
		registry: make(map[ChainType]chain.IChainAdaptor),
	}
	chainAdaptorFactoryMap := map[ChainType]func(conf *config.Config, db *leveldb.Keys, hsmCli *hsm.HsmClient) (chain.IChainAdaptor, error){
		bitcoin.ChainName:  bitcoin.NewChainAdaptor,
		ethereum.ChainName: ethereum.NewChainAdaptor,
		solana.ChainName:   solana.NewChainAdaptor,
	}
	supportedChains := []string{
		bitcoin.ChainName,
		ethereum.ChainName,
		solana.ChainName,
	}

	db, err := leveldb.NewKeyStore(conf.LevelDbPath)
	if err != nil {
		log.Error("new key store level db", "err", err)
		return nil, err
	}
	var hsmClient *hsm.HsmClient
	var errHsmCli error

	if conf.HsmEnable {
		hsmClient, errHsmCli = hsm.NewHSMClient(context.Background(), conf.KeyPath, conf.KeyName)
		if errHsmCli != nil {
			log.Error("new hsm client fail", "err", errHsmCli)
			return nil, errHsmCli
		}
	}
	for _, c := range conf.Chains {
		if factory, ok := chainAdaptorFactoryMap[c]; ok {
			adaptor, err := factory(conf, db, hsmClient)
			if err != nil {
				log.Crit("failed to setup chain", "chain", c, "error", err)
			}
			dispatcher.registry[c] = adaptor
		} else {
			log.Error("unsupported chain", "chain", c, "supportedChains", supportedChains)
		}
	}
	return &dispatcher, nil
}

func (d *ChainDispatcher) Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Error("panic error", "msg", e)
			log.Debug(string(debug.Stack()))
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
		}
	}()
	pos := strings.LastIndex(info.FullMethod, "/")
	method := info.FullMethod[pos+1:]
	consumerToken := req.(CommonRequest).GetConsumerToken()
	chainName := req.(CommonRequest).GetChainName()
	log.Info(method, "chain", chainName, "consumerToken", consumerToken, "req", req)
	resp, err = handler(ctx, req)
	log.Debug("Finish handling", "resp", resp, "err", err)
	return
}

func (d *ChainDispatcher) preHandler(req interface{}) (resp *CommonReply) {
	consumerToken := req.(CommonRequest).GetConsumerToken()
	log.Debug("consumer token", "consumerToken", consumerToken, "req", req)
	if consumerToken != AccessToken {
		return &CommonReply{
			Code:    wallet.ReturnCode_ERROR,
			Message: "consumer token is error",
		}
	}
	chainName := req.(CommonRequest).GetChainName()
	log.Debug("chain name", "chain", chainName, "req", req)
	if _, ok := d.registry[chainName]; !ok {
		return &CommonReply{
			Code:    wallet.ReturnCode_ERROR,
			Message: "unsupported chain",
		}
	}
	return nil
}

func (d *ChainDispatcher) GetChainSignMethod(ctx context.Context, request *wallet.ChainSignMethodRequest) (*wallet.ChainSignMethodResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet.ChainSignMethodResponse{
			Code:    resp.Code,
			Message: resp.Message,
		}, nil
	}
	return d.registry[request.ChainName].GetChainSignMethod(ctx, request)
}

func (d *ChainDispatcher) GetChainSchema(ctx context.Context, request *wallet.ChainSchemaRequest) (*wallet.ChainSchemaResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet.ChainSchemaResponse{
			Code:    resp.Code,
			Message: resp.Message,
		}, nil
	}
	return d.registry[request.ChainName].GetChainSchema(ctx, request)
}

func (d *ChainDispatcher) CreateKeyPairsExportPublicKeyList(ctx context.Context, request *wallet.CreateKeyPairAndExportPublicKeyRequest) (*wallet.CreateKeyPairAndExportPublicKeyResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet.CreateKeyPairAndExportPublicKeyResponse{
			Code:    resp.Code,
			Message: resp.Message,
		}, nil
	}
	return d.registry[request.ChainName].CreateKeyPairsExportPublicKeyList(ctx, request)
}

func (d *ChainDispatcher) CreateKeyPairsWithAddresses(ctx context.Context, request *wallet.CreateKeyPairsWithAddressesRequest) (*wallet.CreateKeyPairsWithAddressesResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet.CreateKeyPairsWithAddressesResponse{
			Code:    resp.Code,
			Message: resp.Message,
		}, nil
	}
	return d.registry[request.ChainName].CreateKeyPairsWithAddresses(ctx, request)
}

func (d *ChainDispatcher) BuildAndSignTransaction(ctx context.Context, request *wallet.BuildAndSignTransactionRequest) (*wallet.BuildAndSignTransactionResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet.BuildAndSignTransactionResponse{
			Code:    resp.Code,
			Message: resp.Message,
		}, nil
	}
	return d.registry[request.ChainName].BuildAndSignTransaction(ctx, request)
}

func (d *ChainDispatcher) BuildAndSignBatchTransaction(ctx context.Context, request *wallet.BuildAndSignBatchTransactionRequest) (*wallet.BuildAndSignBatchTransactionResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet.BuildAndSignBatchTransactionResponse{
			Code:    resp.Code,
			Message: resp.Message,
		}, nil
	}
	return d.registry[request.ChainName].BuildAndSignBatchTransaction(ctx, request)
}
