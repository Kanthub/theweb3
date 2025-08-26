package node

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	tdtypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"

	"github.com/the-web3-contracts/vrf-node/bindings/bls"
	"github.com/the-web3-contracts/vrf-node/client"
	"github.com/the-web3-contracts/vrf-node/config"
	"github.com/the-web3-contracts/vrf-node/database"
	"github.com/the-web3-contracts/vrf-node/sign"
	wsClient "github.com/the-web3-contracts/vrf-node/ws/client"
)

type Node struct {
	wg   sync.WaitGroup
	done chan struct{}

	db *database.DB

	privateKey *ecdsa.PrivateKey
	from       common.Address
	ctx        context.Context
	cancle     context.CancelFunc
	stopChan   chan struct{}
	stopped    atomic.Bool

	wsClient *wsClient.WSClients
	KeyPairs *sign.KeyPair

	signTimeOut      time.Duration
	waitScanInterval time.Duration
	signRequestChan  chan tdtypes.RPCRequest
}

func NewNode(ctx context.Context, db *database.DB, privKey *ecdsa.PrivateKey, keyPairs *sign.KeyPair, shouldRegister bool, cfg *config.Config, shutdown context.CancelCauseFunc) (*Node, error) {
	from := crypto.PubkeyToAddress(privKey.PublicKey)
	pubKey := crypto.CompressPubkey(&privKey.PublicKey)
	pubkeyHex := hex.EncodeToString(pubKey)

	log.Info("public key", "publicKey", pubkeyHex)

	if shouldRegister {
		log.Info("register to operator")
		txKey, txOpt, err := registerOperator(ctx, cfg, privKey, keyPairs)
		if err != nil {
			log.Error("register operator fail", "err", err)
			return nil, err
		}
		log.Info("register success", "txkey", txKey.Hash().String(), "txOpt", txOpt.Hash().String())
	}

	wsCli, err := wsClient.NewWSClient(cfg.Node.WsAddr, "/ws", privKey, pubkeyHex)
	if err != nil {
		log.Error("new ws client fail", "err", err)
		return nil, err
	}

	return &Node{
		wg:   sync.WaitGroup{},
		done: make(chan struct{}),
		db:   db,

		privateKey: privKey,
		from:       from,
		ctx:        ctx,

		wsClient: wsCli,
		KeyPairs: keyPairs,

		signTimeOut:      cfg.Node.SignTimeOut,
		waitScanInterval: cfg.Node.WaitScanInterval,
		signRequestChan:  make(chan tdtypes.RPCRequest, 100),
	}, nil
}

func (n *Node) Start(ctx context.Context) error {
	n.wg.Add(1)
	go n.sign()
	return nil
}

func (n *Node) Stop(ctx context.Context) error {
	n.cancle()
	close(n.done)
	n.wg.Wait()
	n.stopped.Store(true)
	return nil
}

func (n *Node) Stopped() bool {
	return n.stopped.Load()
}

func (n *Node) sign() {
	defer n.wg.Done()
	log.Info("start to sign message")
	go func() {
		defer func() {
			log.Info("exit sign process")
		}()
		for {
			select {
			case <-n.stopChan:
				return
			case req := <-n.signRequestChan:
				log.Error("req", "req", req)
				return
			}
		}
	}()
}

func (n *Node) SignMessage(messageHash string) (*sign.Signature, error) {
	var bSign *sign.Signature
	log.Info("before sign", "messageHash", messageHash)
	bSign = n.KeyPairs.SignMessage(crypto.Keccak256Hash(common.Hex2Bytes(messageHash)))
	log.Info("after sign", "bSign", bSign)
	return bSign, nil
}

func registerOperator(ctx context.Context, cfg *config.Config, privKey *ecdsa.PrivateKey, keyPairs *sign.KeyPair) (*types.Transaction, *types.Transaction, error) {
	ethCli, err := client.DialEthClientWithTimeout(ctx, cfg.Chain.ChainRpcUrl, false)
	if err != nil {
		log.Error("new eth client fail", "err", err)
		return nil, nil, err
	}

	blsRegistry, err := bls.NewBLSApkRegistry(common.HexToAddress(cfg.Chain.BlsRegistryAddress), ethCli)
	if err != nil {
		log.Error("new bls registry fail", "err", err)
		return nil, nil, err
	}

	blsParsed, err := bls.BLSApkRegistryMetaData.GetAbi()
	if err != nil {
		log.Error("fetch bls registry abi fail", "err", err)
		return nil, nil, err
	}

	rawBlsContract := bind.NewBoundContract(common.HexToAddress(cfg.Chain.BlsRegistryAddress), *blsParsed, ethCli, ethCli, ethCli)

	topts, err := client.NewTransactOpts(ctx, uint64(cfg.Chain.ChainId), privKey)
	if err != nil {
		log.Error("new transactopt fail", "err", err)
	}

	nodeAddr := crypto.PubkeyToAddress(privKey.PublicKey)
	latestBlock, _ := ethCli.BlockNumber(ctx)

	cOpts := &bind.CallOpts{
		BlockNumber: big.NewInt(int64(latestBlock)),
		From:        nodeAddr,
	}

	msg, err := blsRegistry.GetPubkeyRegMessageHash(cOpts, nodeAddr)
	if err != nil {
		log.Error("get public key register message hash fail", "err", err)
		return nil, nil, err
	}

	sigMsg := new(bn254.G1Affine).ScalarMultiplication(sign.NewG1Point(msg.Y, msg.Y).G1Affine, keyPairs.PrivKey.BigInt(new(big.Int)))

	blsParam := bls.IBLSApkRegistryPubkeyRegistrationParams{
		PubkeyRegistrationSignature: bls.BN254G1Point{
			X: sigMsg.X.BigInt(new(big.Int)),
			Y: sigMsg.Y.BigInt(new(big.Int)),
		},
		PubkeyG1: bls.BN254G1Point{
			X: keyPairs.GetPubKeyG1().X.BigInt(new(big.Int)),
			Y: keyPairs.GetPubKeyG1().Y.BigInt(new(big.Int)),
		},
		PubkeyG2: bls.BN254G2Point{
			X: [2]*big.Int{keyPairs.GetPubKeyG2().X.A1.BigInt(new(big.Int)), keyPairs.GetPubKeyG2().X.A0.BigInt(new(big.Int))},
			Y: [2]*big.Int{keyPairs.GetPubKeyG2().Y.A1.BigInt(new(big.Int)), keyPairs.GetPubKeyG2().Y.A0.BigInt(new(big.Int))},
		},
	}
	regBlsKey, err := blsRegistry.RegisterBLSPublicKey(topts, nodeAddr, blsParam, msg)
	if err != nil {
		log.Error("register bls public key fail", "err", err)
		return nil, nil, err
	}

	blsRegTx, err := rawBlsContract.RawTransact(topts, regBlsKey.Data())
	if err != nil {
		log.Error("raw bls tx fail", "err", err)
		return nil, nil, err
	}

	err = ethCli.SendTransaction(ctx, blsRegTx)
	if err != nil {
		log.Error("send raw transaction fail", "err", err)
		return nil, nil, err
	}

	_, err = client.GetTransactionReceipt(ctx, ethCli, blsRegTx.Hash())
	if err != nil {
		log.Error("get transaction receipt fail", "err", err)
		return nil, nil, err
	}

	regOperator, err := blsRegistry.RegisterOperator(topts, nodeAddr)
	if err != nil {
		log.Error("register operator fail", "err", err)
		return nil, nil, err
	}

	blsRegOptTx, err := rawBlsContract.RawTransact(topts, regOperator.Data())
	if err != nil {
		log.Error("raw bls tx fail", "err", err)
		return nil, nil, err
	}

	err = ethCli.SendTransaction(ctx, blsRegOptTx)
	if err != nil {
		log.Error("send raw transaction fail", "err", err)
		return nil, nil, err
	}

	_, err = client.GetTransactionReceipt(ctx, ethCli, blsRegOptTx.Hash())
	if err != nil {
		log.Error("get transaction receipt fail", "err", err)
	}
	return blsRegTx, blsRegOptTx, nil
}
