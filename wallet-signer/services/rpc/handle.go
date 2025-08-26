package rpc

import (
	"context"

	"github.com/pkg/errors"

	"github.com/kanthub/wallet-sign/protobuf"
	"github.com/kanthub/wallet-sign/protobuf/wallet"
	"github.com/kanthub/wallet-sign/ssm"
)

func (s *RpcServer) GetChainMethod(ctx context.Context, in *wallet.ChainSignMethodRequest) (*wallet.ChainSignMethodResponse, error) {
	signMethodMap := make(map[string]string)

	signMethodMap["ethereum"] = "ecdsa"
	signMethodMap["bitcoin"] = "ecdsa"
	signMethodMap["solona"] = "eddsa"

	return &wallet.ChainSignMethodResponse{
		Code:       wallet.ReturnCode_SUCCESS,
		Message:    "get sign way success",
		SignMethod: signMethodMap[in.ChainName],
	}, nil
}

func (s *RpcServer) GetChainSchema(ctx context.Context, in *wallet.ChainSchemaRequest) (*wallet.ChainSchemaResponse, error) {
	resp := &wallet.ChainSchemaResponse{
		Code: wallet.ReturnCode_ERROR,
	}

	return resp, nil
}

func (s *RpcServer) SignTxMessage(ctx context.Context, in *wallet.BuildAndSignTransactionRequest) (*wallet.BuildAndSignBatchTransactionResponse, error) {
	resp := &wallet.BuildAndSignTransactionResponse{
		Code: wallet.ReturnCode_ERROR,
	}
	cryptoType, err := protobuf.ParseTransactionType(in.Type)
	if err != nil {
		resp.Msg = "input type error"
		return resp, nil
	}

	privKey, isOk := s.db.GetPrivKey(in.PublicKey)
	if !isOk {
		return nil, errors.New("get private key by public key fail")
	}

	var signature string
	var err2 error

	switch cryptoType {
	case protobuf.ECDSA:
		signature, err2 = ssm.SignECDSAMessage(privKey, in.MessageHash)
	case protobuf.EDDSA:
		signature, err2 = ssm.SignEdDSAMessage(privKey, in.MessageHash)
	default:
		return nil, errors.New("unsupported key type")
	}
	if err2 != nil {
		return nil, err2
	}
	resp.Msg = "sign tx message success"
	resp.Signature = signature
	resp.Code = wallet.ReturnCode_SUCCESS
	return resp, nil
}
