package ethereum

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/pkg/errors"
	"github.com/the-web3/wallet-sign-s6/hsm"
	"github.com/the-web3/wallet-sign-s6/leveldb"
	"github.com/the-web3/wallet-sign-s6/ssm"
	"math/big"

	"github.com/the-web3/wallet-sign-s6/chain"
	"github.com/the-web3/wallet-sign-s6/config"
	"github.com/the-web3/wallet-sign-s6/protobuf/wallet"
)

const ChainName = "Ethereum"

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
		SignMethod: "ecdsa",
	}, nil
}

func (c ChainAdaptor) GetChainSchema(ctx context.Context, req *wallet.ChainSchemaRequest) (*wallet.ChainSchemaResponse, error) {
	es := EthereumSchema{
		RequestId: "0",
		DynamicFeeTx: Eip1559DynamicFeeTx{
			ChainId:              "",
			Nonce:                0,
			FromAddress:          common.Address{}.String(),
			ToAddress:            common.Address{}.String(),
			GasLimit:             0,
			Gas:                  0,
			MaxFeePerGas:         "0",
			MaxPriorityFeePerGas: "0",
			Amount:               "0",
			ContractAddress:      "",
		},
		ClassicFeeTx: LegacyFeeTx{
			ChainId:         "0",
			Nonce:           0,
			FromAddress:     common.Address{}.String(),
			ToAddress:       common.Address{}.String(),
			GasLimit:        0,
			GasPrice:        0,
			Amount:          "0",
			ContractAddress: "",
		},
	}
	b, err := json.Marshal(es)
	if err != nil {
		log.Error("marshal fail", "err", err)
	}
	return &wallet.ChainSchemaResponse{
		Code:    wallet.ReturnCode_SUCCESS,
		Message: "get ethereum sign schema success",
		Schema:  string(b),
	}, nil
}

func (c ChainAdaptor) CreateKeyPairsExportPublicKeyList(ctx context.Context, req *wallet.CreateKeyPairAndExportPublicKeyRequest) (*wallet.CreateKeyPairAndExportPublicKeyResponse, error) {
	var signer ssm.Signer
	signer = &ssm.ECDSASigner{}
	resp := &wallet.CreateKeyPairAndExportPublicKeyResponse{
		Code: wallet.ReturnCode_ERROR,
	}
	if req.KeyNum > 10000 {
		resp.Message = "Number must be less than 100000"
		return resp, nil
	}

	var keyList []leveldb.Key
	var retKeyList []*wallet.ExportPublicKey

	for counter := 0; counter < int(req.KeyNum); counter++ {
		priKeyStr, pubKeyStr, compressPubkeyStr, err := signer.CreateKeyPair()
		if err != nil {
			resp.Message = "create key pairs fail"
			return resp, nil
		}
		keyItem := leveldb.Key{
			PrivateKey: priKeyStr,
			Pubkey:     pubKeyStr,
		}
		pukItem := &wallet.ExportPublicKey{
			CompressPublicKey: compressPubkeyStr,
			PublicKey:         pubKeyStr,
		}
		retKeyList = append(retKeyList, pukItem)
		keyList = append(keyList, keyItem)
	}
	isOk := c.db.StoreKeys(keyList)
	if !isOk {
		log.Error("store keys fail", "isOk", isOk)
		return nil, errors.New("store keys fail")
	}
	resp.Code = wallet.ReturnCode_SUCCESS
	resp.Message = "create keys success"
	resp.PublicKeyList = retKeyList
	return resp, nil
}

func (c ChainAdaptor) CreateKeyPairsWithAddresses(ctx context.Context, req *wallet.CreateKeyPairsWithAddressesRequest) (*wallet.CreateKeyPairsWithAddressesResponse, error) {
	var signer ssm.Signer
	signer = &ssm.ECDSASigner{}

	resp := &wallet.CreateKeyPairsWithAddressesResponse{
		Code: wallet.ReturnCode_ERROR,
	}
	if req.KeyNum > 10000 {
		resp.Message = "Number must be less than 100000"
		return resp, nil
	}
	var keyList []leveldb.Key
	var retKeyWithAddressList []*wallet.ExportPublicKeyWithAddress
	for counter := 0; counter < int(req.KeyNum); counter++ {
		priKeyStr, pubKeyStr, compressPubkeyStr, err := signer.CreateKeyPair()
		if err != nil {
			resp.Message = "create key pairs fail"
			return resp, nil
		}
		keyItem := leveldb.Key{
			PrivateKey: priKeyStr,
			Pubkey:     pubKeyStr,
		}
		publicKeyBytes, err := hex.DecodeString(pubKeyStr)
		pukAddressItem := &wallet.ExportPublicKeyWithAddress{
			CompressPublicKey: compressPubkeyStr,
			PublicKey:         pubKeyStr,
			Address:           common.BytesToAddress(crypto.Keccak256(publicKeyBytes[1:])[12:]).String(),
		}
		retKeyWithAddressList = append(retKeyWithAddressList, pukAddressItem)
		keyList = append(keyList, keyItem)
	}
	isOk := c.db.StoreKeys(keyList)
	if !isOk {
		log.Error("store keys fail", "isOk", isOk)
		return nil, errors.New("store keys fail")
	}
	resp.Code = wallet.ReturnCode_SUCCESS
	resp.Message = "create keys with address success"
	resp.PublicKeyAddresses = retKeyWithAddressList
	return resp, nil
}

func (c ChainAdaptor) BuildAndSignTransaction(ctx context.Context, req *wallet.BuildAndSignTransactionRequest) (*wallet.BuildAndSignTransactionResponse, error) {
	var signer ssm.Signer
	signer = &ssm.ECDSASigner{}

	resp := &wallet.BuildAndSignTransactionResponse{
		Code: wallet.ReturnCode_ERROR,
	}

	dFeeTx, _, err := c.buildDynamicFeeTx(req.TxBase64Body)
	if err != nil {
		return nil, err
	}

	rawTx, err := CreateEip1559UnSignTx(dFeeTx, dFeeTx.ChainID)
	if err != nil {
		log.Error("create un sign tx fail", "err", err)
		resp.Message = "get un sign tx fail"
		return resp, nil
	}

	privKey, isOk := c.db.GetPrivKey(req.PublicKey)
	if !isOk {
		log.Error("get private key by public key fail", "err", err)
		resp.Message = "get private key by public key fail"
		return resp, nil
	}

	signature, err := signer.SignMessage(privKey, rawTx)
	if err != nil {
		log.Error("sign transaction fail", "err", err)
		resp.Message = "sign transaction fail"
		return resp, nil
	}

	inputSignatureByteList, err := hex.DecodeString(signature)
	if err != nil {
		log.Error("decode signature failed", "err", err)
		resp.Message = "decode signature failed"
		return resp, nil
	}

	eip1559Signer, signedTx, signAndHandledTx, txHash, err := CreateEip1559SignedTx(dFeeTx, inputSignatureByteList, dFeeTx.ChainID)
	if err != nil {
		log.Error("create signed tx fail", "err", err)
		resp.Message = "create signed tx fail"
		return resp, nil
	}
	log.Info("sign transaction success",
		"eip1559Signer", eip1559Signer,
		"signedTx", signedTx,
		"signAndHandledTx", signAndHandledTx,
		"txHash", txHash,
	)
	resp.Code = wallet.ReturnCode_SUCCESS
	resp.Message = "sign whole transaction success"
	resp.SignedTx = signAndHandledTx
	resp.TxHash = txHash
	resp.TxMessageHash = rawTx
	return resp, nil
}

func (c ChainAdaptor) BuildAndSignBatchTransaction(ctx context.Context, req *wallet.BuildAndSignBatchTransactionRequest) (*wallet.BuildAndSignBatchTransactionResponse, error) {
	panic("implement me")
}

func (c ChainAdaptor) buildDynamicFeeTx(base64Tx string) (*types.DynamicFeeTx, *Eip1559DynamicFeeTx, error) {
	// 1. Decode base64 string
	txReqJsonByte, err := base64.StdEncoding.DecodeString(base64Tx)
	if err != nil {
		log.Error("decode string fail", "err", err)
		return nil, nil, err
	}

	// 2. Unmarshal JSON to struct
	var dynamicFeeTx Eip1559DynamicFeeTx
	if err := json.Unmarshal(txReqJsonByte, &dynamicFeeTx); err != nil {
		log.Error("parse json fail", "err", err)
		return nil, nil, err
	}

	// 3. Convert string values to big.Int
	chainID := new(big.Int)
	maxPriorityFeePerGas := new(big.Int)
	maxFeePerGas := new(big.Int)
	amount := new(big.Int)

	if _, ok := chainID.SetString(dynamicFeeTx.ChainId, 10); !ok {
		return nil, nil, fmt.Errorf("invalid chain ID: %s", dynamicFeeTx.ChainId)
	}
	if _, ok := maxPriorityFeePerGas.SetString(dynamicFeeTx.MaxPriorityFeePerGas, 10); !ok {
		return nil, nil, fmt.Errorf("invalid max priority fee: %s", dynamicFeeTx.MaxPriorityFeePerGas)
	}
	if _, ok := maxFeePerGas.SetString(dynamicFeeTx.MaxFeePerGas, 10); !ok {
		return nil, nil, fmt.Errorf("invalid max fee: %s", dynamicFeeTx.MaxFeePerGas)
	}
	if _, ok := amount.SetString(dynamicFeeTx.Amount, 10); !ok {
		return nil, nil, fmt.Errorf("invalid amount: %s", dynamicFeeTx.Amount)
	}

	// 4. Handle addresses and data
	toAddress := common.HexToAddress(dynamicFeeTx.ToAddress)
	var finalToAddress common.Address
	var finalAmount *big.Int
	var buildData []byte
	log.Info("contract address check",
		"contractAddress", dynamicFeeTx.ContractAddress,
		"isEthTransfer", isEthTransfer(&dynamicFeeTx),
	)

	// 5. Handle contract interaction vs direct transfer
	if isEthTransfer(&dynamicFeeTx) {
		finalToAddress = toAddress
		finalAmount = amount
	} else {
		contractAddress := common.HexToAddress(dynamicFeeTx.ContractAddress)
		buildData = BuildErc20Data(toAddress, amount)
		finalToAddress = contractAddress
		finalAmount = big.NewInt(0)
	}

	// 6. Create dynamic fee transaction
	dFeeTx := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     dynamicFeeTx.Nonce,
		GasTipCap: maxPriorityFeePerGas,
		GasFeeCap: maxFeePerGas,
		Gas:       dynamicFeeTx.GasLimit,
		To:        &finalToAddress,
		Value:     finalAmount,
		Data:      buildData,
	}

	return dFeeTx, &dynamicFeeTx, nil

}

func isEthTransfer(tx *Eip1559DynamicFeeTx) bool {
	if tx.ContractAddress == "" ||
		tx.ContractAddress == "0x0000000000000000000000000000000000000000" ||
		tx.ContractAddress == "0x00" {
		return true
	}
	return false
}

/*
{
  "code": "SUCCESS",
  "message": "create keys with address success",
  "public_key_addresses": [
    {
      "public_key": "04bc07a90250df4be52d23310d34cde52d57f4d95f9826cedb594cf3ae002f285ab9708f3303ba526fb3bcb68896568922ef6673ea64e4c49ac57a9c784cfdb811",
      "compress_public_key": "03bc07a90250df4be52d23310d34cde52d57f4d95f9826cedb594cf3ae002f285a",
      "address": "0x360Ccc76AE5DA207d0960cF81f2589d9d4d5F26D"
    },
    {
      "public_key": "042db546d77bf9427b9dd6b2ce1ce5343ab7c83aa70d73f3472a61b8abe972665d5b3aaf73c75555728547c7adbe3d069139b8d423a4ecf8d9d3f83e30088da05f",
      "compress_public_key": "032db546d77bf9427b9dd6b2ce1ce5343ab7c83aa70d73f3472a61b8abe972665d",
      "address": "0x86C79442fAEce848b194050Bb1dEDD0b8Cad7487"
    }
  ]
}

{
    "chain_id": "11155111",
    "nonce": 0,
    "from_address": "0x360Ccc76AE5DA207d0960cF81f2589d9d4d5F26D",
    "to_address": "0x45Bd8ea16cFEB0D937a2D98cBEb0300e3E689Fe7",
    "gas_limit": 21000,
    "gas": 2000000,
    "max_fee_per_gas": "327993150328",
    "max_priority_fee_per_gas": "32799315032",
    "amount": "100000000000000000",
    "contract_address": "0x00"
}
*/
