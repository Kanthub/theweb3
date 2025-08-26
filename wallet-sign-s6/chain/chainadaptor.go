package chain

import (
	"context"

	"github.com/the-web3/wallet-sign-s6/protobuf/wallet"
)

type IChainAdaptor interface {
	GetChainSignMethod(ctx context.Context, req *wallet.ChainSignMethodRequest) (*wallet.ChainSignMethodResponse, error)
	GetChainSchema(ctx context.Context, req *wallet.ChainSchemaRequest) (*wallet.ChainSchemaResponse, error)
	CreateKeyPairsExportPublicKeyList(ctx context.Context, req *wallet.CreateKeyPairAndExportPublicKeyRequest) (*wallet.CreateKeyPairAndExportPublicKeyResponse, error)
	CreateKeyPairsWithAddresses(ctx context.Context, req *wallet.CreateKeyPairsWithAddressesRequest) (*wallet.CreateKeyPairsWithAddressesResponse, error)
	BuildAndSignTransaction(ctx context.Context, req *wallet.BuildAndSignTransactionRequest) (*wallet.BuildAndSignTransactionResponse, error)
	BuildAndSignBatchTransaction(ctx context.Context, req *wallet.BuildAndSignBatchTransactionRequest) (*wallet.BuildAndSignBatchTransactionResponse, error)
}
