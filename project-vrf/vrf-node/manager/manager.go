package manager

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"

	"math/big"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"

	"github.com/the-web3-contracts/vrf-node/bindings/bls"
	"github.com/the-web3-contracts/vrf-node/bindings/vrf"
	"github.com/the-web3-contracts/vrf-node/client"
	"github.com/the-web3-contracts/vrf-node/config"
	"github.com/the-web3-contracts/vrf-node/database"
	"github.com/the-web3-contracts/vrf-node/manager/router"
	"github.com/the-web3-contracts/vrf-node/manager/types"
	"github.com/the-web3-contracts/vrf-node/sign"
	"github.com/the-web3-contracts/vrf-node/ws/server"
)

var (
	errNotEnoughSignNode = errors.New("not enough available nodes to sign")
	errNotEnoughVoteNode = errors.New("not enough available nodes to vote")
)

type Manager struct {
	wg              sync.WaitGroup
	done            chan struct{}
	log             log.Logger
	db              *database.DB
	wsServer        server.IWebsocketManager
	NodeMembers     []string
	httpAddr        string
	httpServer      *http.Server
	sStakeUrl       string
	mu              sync.Mutex
	ctx             context.Context
	stopped         atomic.Bool
	ethChainID      uint64
	privateKey      *ecdsa.PrivateKey
	from            common.Address
	ethClient       *ethclient.Client
	vrfContract     *vrf.DappLinkVRFManager
	vrfContractAddr common.Address
	rawVrfContract  *bind.BoundContract
	barContract     *bls.BLSApkRegistry
	barContractAddr common.Address
	rawBarContract  *bind.BoundContract
	isFirstBatch    bool
	signTimeout     time.Duration
	fPTimeout       time.Duration
}

func NewManager(ctx context.Context, db *database.DB, wsServer server.IWebsocketManager, cfg *config.Config, priv *ecdsa.PrivateKey, shutdown context.CancelCauseFunc) (*Manager, error) {
	ethCli, err := client.DialEthClientWithTimeout(ctx, cfg.Chain.ChainRpcUrl, false)
	if err != nil {
		return nil, err
	}
	vrfContract, err := vrf.NewDappLinkVRFManager(common.HexToAddress(cfg.Chain.DappLinkVrfContractAddress), ethCli)
	if err != nil {
		return nil, err
	}
	fParsed, err := abi.JSON(strings.NewReader(
		vrf.DappLinkVRFFactoryMetaData.ABI,
	))
	if err != nil {
		return nil, err
	}
	rawVrfContract := bind.NewBoundContract(
		common.HexToAddress(cfg.Chain.DappLinkVrfContractAddress), fParsed, ethCli, ethCli,
		ethCli,
	)

	barContract, err := bls.NewBLSApkRegistry(common.HexToAddress(cfg.Chain.BlsRegistryAddress), ethCli)
	if err != nil {
		return nil, err
	}
	bParsed, err := bls.BLSApkRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	rawBarContract := bind.NewBoundContract(
		common.HexToAddress(cfg.Chain.BlsRegistryAddress), *bParsed, ethCli, ethCli,
		ethCli,
	)

	nodeMemberS := strings.Split(cfg.Manager.NodeMembers, ",")
	for _, nodeMember := range nodeMemberS {
		fmt.Println("nodeMember===", nodeMember)
		//if err := db.SetActiveMember(nodeMember); err != nil {
		//	return nil, fmt.Errorf("failed to set node member, err: %v", err)
		//}
	}

	return &Manager{
		done:            make(chan struct{}),
		db:              db,
		wsServer:        wsServer,
		httpAddr:        cfg.Manager.HttpAddr,
		NodeMembers:     nodeMemberS,
		ctx:             ctx,
		privateKey:      priv,
		from:            crypto.PubkeyToAddress(priv.PublicKey),
		signTimeout:     time.Second * 10,
		fPTimeout:       time.Second * 10,
		ethChainID:      uint64(cfg.Chain.ChainId),
		ethClient:       ethCli,
		vrfContract:     vrfContract,
		vrfContractAddr: common.HexToAddress(cfg.Chain.DappLinkVrfContractAddress),
		rawVrfContract:  rawVrfContract,
		barContract:     barContract,
		barContractAddr: common.HexToAddress(cfg.Chain.BlsRegistryAddress),
		rawBarContract:  rawBarContract,
	}, nil
}

func (m *Manager) Start(ctx context.Context) error {
	waitNodeTicker := time.NewTicker(5 * time.Second)
	var done bool
	for !done {
		select {
		case <-waitNodeTicker.C:
			availableNodes := m.availableNodes(m.NodeMembers)
			if len(availableNodes) < len(m.NodeMembers) {
				m.log.Warn("wait node to connect", "availableNodesNum", len(availableNodes), "connectedNodeNum", len(m.NodeMembers))
				continue
			} else {
				done = true
				break
			}
		}
	}

	registry := router.NewRegistry(m, m.db)
	r := gin.Default()
	registry.Register(r)

	var s *http.Server
	s = &http.Server{
		Addr:    m.httpAddr,
		Handler: r,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			m.log.Error("api server starts failed", "err", err)
		}
	}()
	m.httpServer = s

	m.wg.Add(1)
	go m.work()
	log.Info("manager is starting......")
	return nil
}

func (m *Manager) Stop(ctx context.Context) error {
	close(m.done)
	if err := m.httpServer.Shutdown(ctx); err != nil {
		m.log.Error("http server forced to shutdown", "err", err)
		return err
	}
	m.stopped.Store(true)
	m.log.Info("Server exiting")
	return nil
}

func (m *Manager) Stopped() bool {
	return m.stopped.Load()
}

func (m *Manager) work() {
	fpTicker := time.NewTicker(m.fPTimeout)
	defer m.wg.Done()

	for {
		select {
		case <-fpTicker.C:
			var signature *sign.G1Point
			var g2Point *sign.G2Point
			var NonSignerPubkeys []vrf.BN254G1Point

			request := types.SignMsgRequest{}

			res, err := m.SignMsgBatch(request)
			if errors.Is(err, errNotEnoughSignNode) || errors.Is(err, errNotEnoughVoteNode) {
				m.log.Error("not enough available nodes to sign or not enough available nodes to vote")
				continue
			} else if err != nil {
				m.log.Error("failed to sign msg", "err", err)
				continue
			}
			m.log.Info("success to sign msg", "signature", res.Signature)

			signature = res.Signature
			g2Point = res.G2Point

			for _, v := range res.NonSignerPubkeys {
				NonSignerPubkeys = append(NonSignerPubkeys, vrf.BN254G1Point{
					X: v.X.BigInt(new(big.Int)),
					Y: v.Y.BigInt(new(big.Int)),
				})
			}

			opts, err := client.NewTransactOpts(m.ctx, m.ethChainID, m.privateKey)
			if err != nil {
				m.log.Error("failed to new transact opts", "err", err)
				continue
			}
			vrfNonSignerAndSignature := vrf.IBLSApkRegistryVrfNoSignerAndSignature{
				NonSignerPubKeys: NonSignerPubkeys,
				ApkG2: vrf.BN254G2Point{
					X: [2]*big.Int{g2Point.X.A1.BigInt(new(big.Int)), g2Point.X.A0.BigInt(new(big.Int))},
					Y: [2]*big.Int{g2Point.Y.A1.BigInt(new(big.Int)), g2Point.Y.A0.BigInt(new(big.Int))},
				},
				Sigma: vrf.BN254G1Point{
					X: signature.X.BigInt(new(big.Int)),
					Y: signature.Y.BigInt(new(big.Int)),
				},
				TotalBtcStake:      big.NewInt(100),
				TotalDappLinkStake: big.NewInt(100),
			}

			var aaa []*big.Int
			var bb [32]byte

			tx, err := m.vrfContract.FulfillRandomWords(opts, big.NewInt(1), aaa, bb, big.NewInt(1), vrfNonSignerAndSignature)
			if err != nil {
				m.log.Error("failed to craft VerifyFinalitySignature transaction", "err", err)
				continue
			}
			rTx, err := m.rawVrfContract.RawTransact(opts, tx.Data())
			if err != nil {
				m.log.Error("failed to raw VerifyFinalitySignature transaction", "err", err)
				continue
			}
			err = m.ethClient.SendTransaction(m.ctx, tx)
			if err != nil {
				m.log.Error("failed to send VerifyFinalitySignature transaction", "err", err)
				continue
			}

			receipt, err := client.GetTransactionReceipt(m.ctx, m.ethClient, rTx.Hash())
			if err != nil {
				m.log.Error("failed to get verify finality transaction receipt", "err", err)
				continue
			}
			m.log.Info("success to send verify finality signature transaction", "tx_hash", receipt.TxHash.String())
		case <-m.done:
			return
		}
	}
}

func (m *Manager) SignMsgBatch(request types.SignMsgRequest) (*types.SignResult, error) {
	m.log.Info("received sign request")

	activeMember, err := m.db.Member.GetActiveMember()
	if err != nil {
		m.log.Error("failed to get active member from db", "err", err)
		return nil, err
	}

	availableNodes := m.availableNodes(activeMember)
	if len(availableNodes) == 0 {
		m.log.Warn("not enough sign node", "availableNodes", availableNodes)
		return nil, errNotEnoughSignNode
	}

	ctx := types.NewContext().WithAvailableNodes(availableNodes).WithRequestId(randomRequestId())

	var resp types.SignResult
	var signErr error
	resp, signErr = m.sign(ctx, request, types.SignMsgBatch)
	if signErr != nil {
		return nil, signErr
	}
	if resp.Signature == nil {
		return nil, errNotEnoughVoteNode
	}
	return &resp, nil
}

func (m *Manager) availableNodes(nodeMembers []string) []string {
	aliveNodes := m.wsServer.AliveNodes()
	log.Info("check available nodes", "expected", fmt.Sprintf("%v", nodeMembers), "alive nodes", fmt.Sprintf("%v", aliveNodes))
	availableNodes := make([]string, 0)
	for _, n := range aliveNodes {
		if ExistsIgnoreCase(nodeMembers, n) {
			availableNodes = append(availableNodes, n)
		}
	}
	return availableNodes
}

func randomRequestId() string {
	code := fmt.Sprintf("%04v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))
	return time.Now().Format("20060102150405") + code
}

func ExistsIgnoreCase(slice []string, target string) bool {
	for _, item := range slice {
		if strings.EqualFold(item, target) {
			return true
		}
	}
	return false
}
