package solana

import (
	"context"
	"encoding/json"

	"github.com/ethereum/go-ethereum/log"

	"github.com/the-web3/wallet-sign-s6/chain"
	"github.com/the-web3/wallet-sign-s6/config"
	"github.com/the-web3/wallet-sign-s6/hsm"
	"github.com/the-web3/wallet-sign-s6/leveldb"
	"github.com/the-web3/wallet-sign-s6/protobuf/wallet"
)

const ChainName = "Solana"

type ChainAdaptor struct {
	db        *leveldb.Keys
	HsmClient *hsm.HsmClient
}

func NewChainAdaptor(conf *config.Config, db *leveldb.Keys, hsmCli *hsm.HsmClient) (chain.IChainAdaptor, error) {
	return &ChainAdaptor{
		db:        db,
		HsmClient: hsmCli,
	}, nil
}

func (c ChainAdaptor) GetChainSignMethod(ctx context.Context, req *wallet.ChainSignMethodRequest) (*wallet.ChainSignMethodResponse, error) {
	return &wallet.ChainSignMethodResponse{
		Code:       wallet.ReturnCode_SUCCESS,
		Message:    "get sign method success",
		SignMethod: "eddsa",
	}, nil
}

func (c ChainAdaptor) GetChainSchema(ctx context.Context, req *wallet.ChainSchemaRequest) (*wallet.ChainSchemaResponse, error) {
	ss := SolanaSchema{
		Nonce:           "",
		GasPrice:        "",
		GasTipCap:       "",
		GasFeeCap:       "",
		Gas:             0,
		ContractAddress: "",
		FromAddress:     "",
		ToAddress:       "",
		TokenId:         "",
		Value:           "",
	}
	b, err := json.Marshal(ss)
	if err != nil {
		log.Error("marshal fail", "err", err)
	}
	return &wallet.ChainSchemaResponse{
		Code:    wallet.ReturnCode_SUCCESS,
		Message: "get solana sign schema success",
		Schema:  string(b),
	}, nil
}

func (c ChainAdaptor) CreateKeyPairsExportPublicKeyList(ctx context.Context, req *wallet.CreateKeyPairAndExportPublicKeyRequest) (*wallet.CreateKeyPairAndExportPublicKeyResponse, error) {
	panic("implement me")
}

func (c ChainAdaptor) CreateKeyPairsWithAddresses(ctx context.Context, req *wallet.CreateKeyPairsWithAddressesRequest) (*wallet.CreateKeyPairsWithAddressesResponse, error) {
	panic("implement me")
}

func (c ChainAdaptor) BuildAndSignTransaction(ctx context.Context, req *wallet.BuildAndSignTransactionRequest) (*wallet.BuildAndSignTransactionResponse, error) {
	panic("implement me")
}

func (c ChainAdaptor) BuildAndSignBatchTransaction(ctx context.Context, req *wallet.BuildAndSignBatchTransactionRequest) (*wallet.BuildAndSignBatchTransactionResponse, error) {
	panic("implement me")
}
