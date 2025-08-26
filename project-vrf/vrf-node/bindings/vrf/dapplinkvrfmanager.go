// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// BN254G1Point is an auto generated low-level Go binding around an user-defined struct.
type BN254G1Point struct {
	X *big.Int
	Y *big.Int
}

// BN254G2Point is an auto generated low-level Go binding around an user-defined struct.
type BN254G2Point struct {
	X [2]*big.Int
	Y [2]*big.Int
}

// IBLSApkRegistryVrfNoSignerAndSignature is an auto generated low-level Go binding around an user-defined struct.
type IBLSApkRegistryVrfNoSignerAndSignature struct {
	NonSignerPubKeys   []BN254G1Point
	ApkG2              BN254G2Point
	Sigma              BN254G1Point
	TotalDappLinkStake *big.Int
	TotalBtcStake      *big.Int
}

// DappLinkVRFManagerMetaData contains all meta data concerning the DappLinkVRFManager contract.
var DappLinkVRFManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"blsRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIBLSApkRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"dappLinkAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"fulfillRandomWords\",\"inputs\":[{\"name\":\"_requestId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_randomWords\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"msgHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIBLSApkRegistry.VrfNoSignerAndSignature\",\"components\":[{\"name\":\"nonSignerPubKeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"totalDappLinkStake\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalBtcStake\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getRequestStatus\",\"inputs\":[{\"name\":\"_requestId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"fulfilled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"randomWords\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_dappLinkAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_blsRegistry\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"lastRequestId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestIds\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"requestMapping\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"fulfilled\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"requestRandomWords\",\"inputs\":[{\"name\":\"_requestId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_numWords\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setDappLink\",\"inputs\":[{\"name\":\"_dappLinkAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"FillRandomWords\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"randomWords\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RequestSent\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"_numWords\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"current\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x608060405234801561001057600080fd5b5061001961001e565b6100d0565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff161561006e5760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b03908116146100cd5780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b610e37806100df6000396000f3fe608060405234801561001057600080fd5b50600436106100cf5760003560e01c806393e4b8801161008c578063d8a4676f11610066578063d8a4676f146101da578063f0c28a41146101fb578063f2fde38b1461020e578063fc2a88c31461022157600080fd5b806393e4b880146101a1578063996869d0146101b4578063c0c53b8b146101c757600080fd5b80631b565710146100d45780631b739ef1146100e9578063715018a6146100fc57806382e215ab146101045780638796ba8c1461013c5780638da5cb5b1461015d575b600080fd5b6100e76100e2366004610aa1565b61022a565b005b6100e76100f7366004610b79565b61038d565b6100e7610462565b610127610112366004610b9b565b60046020526000908152604090205460ff1681565b60405190151581526020015b60405180910390f35b61014f61014a366004610b9b565b610476565b604051908152602001610133565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b03165b6040516001600160a01b039091168152602001610133565b600354610189906001600160a01b031681565b6100e76101c2366004610bd0565b610497565b6100e76101d5366004610bf2565b6104c1565b6101ed6101e8366004610b9b565b610603565b604051610133929190610c71565b600254610189906001600160a01b031681565b6100e761021c366004610bd0565b610679565b61014f60015481565b6002546001600160a01b031633146102895760405162461bcd60e51b815260206004820152601860248201527f446170704c696e6b5652462e6f6e6c79446170704c696e6b000000000000000060448201526064015b60405180910390fd5b60035460405163041c37a160e21b81526001600160a01b0390911690631070de84906102bd90869086908690600401610cdf565b606060405180830381865afa1580156102da573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102fe9190610da0565b50506040805180820182526001808252602080830188815260008a815260048352949094208351815460ff1916901515178155935180519394936103499385019291909101906107e5565b509050507ff3cb4deb0441dd096356debf166f879d78cadc19e4b94053c8bea6d3940de93a858560405161037e929190610de8565b60405180910390a15050505050565b6103956106b7565b604080518082018252600080825282518181526020808201855280840191825286835260048152939091208251815460ff191690151517815590518051929391926103e692600185019201906107e5565b505060008054600180820183559180527f290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e56301849055839055506040805183815260208101839052308183015290517fe697eb68c0228bd7d4e553246a2a86e8402d0895e45092ef8ae87b4cfd29f0169181900360600190a15050565b61046a6106b7565b6104746000610712565b565b6000818154811061048657600080fd5b600091825260209091200154905081565b61049f6106b7565b600280546001600160a01b0319166001600160a01b0392909216919091179055565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a008054600160401b810460ff16159067ffffffffffffffff166000811580156105075750825b905060008267ffffffffffffffff1660011480156105245750303b155b905081158015610532575080155b156105505760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561057a57845460ff60401b1916600160401b1785555b61058388610783565b600380546001600160a01b038089166001600160a01b03199283161790925560028054928a169290911691909117905583156105f957845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b5050505050505050565b600081815260046020908152604080832080546001909101805483518186028101860190945280845260609460ff90931693919283919083018282801561066957602002820191906000526020600020905b815481526020019060010190808311610655575b5050505050905091509150915091565b6106816106b7565b6001600160a01b0381166106ab57604051631e4fbdf760e01b815260006004820152602401610280565b6106b481610712565b50565b336106e97f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b0316146104745760405163118cdaa760e01b8152336004820152602401610280565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080546001600160a01b031981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a3505050565b61078b610794565b6106b4816107dd565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0054600160401b900460ff1661047457604051631afcd79f60e31b815260040160405180910390fd5b610681610794565b828054828255906000526020600020908101928215610820579160200282015b82811115610820578251825591602001919060010190610805565b5061082c929150610830565b5090565b5b8082111561082c5760008155600101610831565b634e487b7160e01b600052604160045260246000fd5b6040805190810167ffffffffffffffff8111828210171561087e5761087e610845565b60405290565b60405160a0810167ffffffffffffffff8111828210171561087e5761087e610845565b604051601f8201601f1916810167ffffffffffffffff811182821017156108d0576108d0610845565b604052919050565b600067ffffffffffffffff8211156108f2576108f2610845565b5060051b60200190565b60006040828403121561090e57600080fd5b61091661085b565b9050813581526020820135602082015292915050565b600082601f83011261093d57600080fd5b61094561085b565b80604084018581111561095757600080fd5b845b81811015610971578035845260209384019301610959565b509095945050505050565b60006080828403121561098e57600080fd5b61099661085b565b90506109a2838361092c565b81526109b1836040840161092c565b602082015292915050565b600061012082840312156109cf57600080fd5b6109d7610884565b9050813567ffffffffffffffff8111156109f057600080fd5b8201601f81018413610a0157600080fd5b80356020610a16610a11836108d8565b6108a7565b82815260069290921b83018101918181019087841115610a3557600080fd5b938201935b83851015610a5e57610a4c88866108fc565b82528282019150604085019450610a3a565b855250610a6d8686830161097c565b81850152505050610a818360a084016108fc565b604082015260e08201356060820152610100820135608082015292915050565b600080600080600060a08688031215610ab957600080fd5b8535945060208087013567ffffffffffffffff80821115610ad957600080fd5b818901915089601f830112610aed57600080fd5b8135610afb610a11826108d8565b81815260059190911b8301840190848101908c831115610b1a57600080fd5b938501935b82851015610b3857843582529385019390850190610b1f565b9850505060408901359550606089013594506080890135925080831115610b5e57600080fd5b5050610b6c888289016109bc565b9150509295509295909350565b60008060408385031215610b8c57600080fd5b50508035926020909101359150565b600060208284031215610bad57600080fd5b5035919050565b80356001600160a01b0381168114610bcb57600080fd5b919050565b600060208284031215610be257600080fd5b610beb82610bb4565b9392505050565b600080600060608486031215610c0757600080fd5b610c1084610bb4565b9250610c1e60208501610bb4565b9150610c2c60408501610bb4565b90509250925092565b60008151808452602080850194506020840160005b83811015610c6657815187529582019590820190600101610c4a565b509495945050505050565b8215158152604060208201526000610c8c6040830184610c35565b949350505050565b8060005b6002811015610cb7578151845260209384019390910190600101610c98565b50505050565b610cc8828251610c94565b6020810151610cda6040840182610c94565b505050565b838152600060208460208401526040606060408501526101808401855161012060608701528181518084526101a088019150602083019350600092505b80831015610d4d57610d3982855180518252602090810151910152565b928501926001929092019190840190610d1c565b5060208801519450610d626080880186610cbd565b604088015180516101008901526020015161012088015260608801516101408801526080909701516101609096019590955250939695505050505050565b6000808284036060811215610db457600080fd5b6040811215610dc257600080fd5b50610dcb61085b565b835181526020808501519082015260409093015192949293505050565b828152604060208201526000610c8c6040830184610c3556fea2646970667358221220be5dc8f2a3ea927ffe7b728906664cdafc70722a498d7516772de8406611878464736f6c63430008160033",
}

// DappLinkVRFManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use DappLinkVRFManagerMetaData.ABI instead.
var DappLinkVRFManagerABI = DappLinkVRFManagerMetaData.ABI

// DappLinkVRFManagerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DappLinkVRFManagerMetaData.Bin instead.
var DappLinkVRFManagerBin = DappLinkVRFManagerMetaData.Bin

// DeployDappLinkVRFManager deploys a new Ethereum contract, binding an instance of DappLinkVRFManager to it.
func DeployDappLinkVRFManager(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *DappLinkVRFManager, error) {
	parsed, err := DappLinkVRFManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DappLinkVRFManagerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DappLinkVRFManager{DappLinkVRFManagerCaller: DappLinkVRFManagerCaller{contract: contract}, DappLinkVRFManagerTransactor: DappLinkVRFManagerTransactor{contract: contract}, DappLinkVRFManagerFilterer: DappLinkVRFManagerFilterer{contract: contract}}, nil
}

// DappLinkVRFManager is an auto generated Go binding around an Ethereum contract.
type DappLinkVRFManager struct {
	DappLinkVRFManagerCaller     // Read-only binding to the contract
	DappLinkVRFManagerTransactor // Write-only binding to the contract
	DappLinkVRFManagerFilterer   // Log filterer for contract events
}

// DappLinkVRFManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type DappLinkVRFManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DappLinkVRFManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DappLinkVRFManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DappLinkVRFManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DappLinkVRFManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DappLinkVRFManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DappLinkVRFManagerSession struct {
	Contract     *DappLinkVRFManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// DappLinkVRFManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DappLinkVRFManagerCallerSession struct {
	Contract *DappLinkVRFManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// DappLinkVRFManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DappLinkVRFManagerTransactorSession struct {
	Contract     *DappLinkVRFManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// DappLinkVRFManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type DappLinkVRFManagerRaw struct {
	Contract *DappLinkVRFManager // Generic contract binding to access the raw methods on
}

// DappLinkVRFManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DappLinkVRFManagerCallerRaw struct {
	Contract *DappLinkVRFManagerCaller // Generic read-only contract binding to access the raw methods on
}

// DappLinkVRFManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DappLinkVRFManagerTransactorRaw struct {
	Contract *DappLinkVRFManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDappLinkVRFManager creates a new instance of DappLinkVRFManager, bound to a specific deployed contract.
func NewDappLinkVRFManager(address common.Address, backend bind.ContractBackend) (*DappLinkVRFManager, error) {
	contract, err := bindDappLinkVRFManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DappLinkVRFManager{DappLinkVRFManagerCaller: DappLinkVRFManagerCaller{contract: contract}, DappLinkVRFManagerTransactor: DappLinkVRFManagerTransactor{contract: contract}, DappLinkVRFManagerFilterer: DappLinkVRFManagerFilterer{contract: contract}}, nil
}

// NewDappLinkVRFManagerCaller creates a new read-only instance of DappLinkVRFManager, bound to a specific deployed contract.
func NewDappLinkVRFManagerCaller(address common.Address, caller bind.ContractCaller) (*DappLinkVRFManagerCaller, error) {
	contract, err := bindDappLinkVRFManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DappLinkVRFManagerCaller{contract: contract}, nil
}

// NewDappLinkVRFManagerTransactor creates a new write-only instance of DappLinkVRFManager, bound to a specific deployed contract.
func NewDappLinkVRFManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*DappLinkVRFManagerTransactor, error) {
	contract, err := bindDappLinkVRFManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DappLinkVRFManagerTransactor{contract: contract}, nil
}

// NewDappLinkVRFManagerFilterer creates a new log filterer instance of DappLinkVRFManager, bound to a specific deployed contract.
func NewDappLinkVRFManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*DappLinkVRFManagerFilterer, error) {
	contract, err := bindDappLinkVRFManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DappLinkVRFManagerFilterer{contract: contract}, nil
}

// bindDappLinkVRFManager binds a generic wrapper to an already deployed contract.
func bindDappLinkVRFManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DappLinkVRFManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DappLinkVRFManager *DappLinkVRFManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DappLinkVRFManager.Contract.DappLinkVRFManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DappLinkVRFManager *DappLinkVRFManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DappLinkVRFManager.Contract.DappLinkVRFManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DappLinkVRFManager *DappLinkVRFManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DappLinkVRFManager.Contract.DappLinkVRFManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DappLinkVRFManager *DappLinkVRFManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DappLinkVRFManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DappLinkVRFManager *DappLinkVRFManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DappLinkVRFManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DappLinkVRFManager *DappLinkVRFManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DappLinkVRFManager.Contract.contract.Transact(opts, method, params...)
}

// BlsRegistry is a free data retrieval call binding the contract method 0x93e4b880.
//
// Solidity: function blsRegistry() view returns(address)
func (_DappLinkVRFManager *DappLinkVRFManagerCaller) BlsRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DappLinkVRFManager.contract.Call(opts, &out, "blsRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BlsRegistry is a free data retrieval call binding the contract method 0x93e4b880.
//
// Solidity: function blsRegistry() view returns(address)
func (_DappLinkVRFManager *DappLinkVRFManagerSession) BlsRegistry() (common.Address, error) {
	return _DappLinkVRFManager.Contract.BlsRegistry(&_DappLinkVRFManager.CallOpts)
}

// BlsRegistry is a free data retrieval call binding the contract method 0x93e4b880.
//
// Solidity: function blsRegistry() view returns(address)
func (_DappLinkVRFManager *DappLinkVRFManagerCallerSession) BlsRegistry() (common.Address, error) {
	return _DappLinkVRFManager.Contract.BlsRegistry(&_DappLinkVRFManager.CallOpts)
}

// DappLinkAddress is a free data retrieval call binding the contract method 0xf0c28a41.
//
// Solidity: function dappLinkAddress() view returns(address)
func (_DappLinkVRFManager *DappLinkVRFManagerCaller) DappLinkAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DappLinkVRFManager.contract.Call(opts, &out, "dappLinkAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DappLinkAddress is a free data retrieval call binding the contract method 0xf0c28a41.
//
// Solidity: function dappLinkAddress() view returns(address)
func (_DappLinkVRFManager *DappLinkVRFManagerSession) DappLinkAddress() (common.Address, error) {
	return _DappLinkVRFManager.Contract.DappLinkAddress(&_DappLinkVRFManager.CallOpts)
}

// DappLinkAddress is a free data retrieval call binding the contract method 0xf0c28a41.
//
// Solidity: function dappLinkAddress() view returns(address)
func (_DappLinkVRFManager *DappLinkVRFManagerCallerSession) DappLinkAddress() (common.Address, error) {
	return _DappLinkVRFManager.Contract.DappLinkAddress(&_DappLinkVRFManager.CallOpts)
}

// GetRequestStatus is a free data retrieval call binding the contract method 0xd8a4676f.
//
// Solidity: function getRequestStatus(uint256 _requestId) view returns(bool fulfilled, uint256[] randomWords)
func (_DappLinkVRFManager *DappLinkVRFManagerCaller) GetRequestStatus(opts *bind.CallOpts, _requestId *big.Int) (struct {
	Fulfilled   bool
	RandomWords []*big.Int
}, error) {
	var out []interface{}
	err := _DappLinkVRFManager.contract.Call(opts, &out, "getRequestStatus", _requestId)

	outstruct := new(struct {
		Fulfilled   bool
		RandomWords []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fulfilled = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.RandomWords = *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// GetRequestStatus is a free data retrieval call binding the contract method 0xd8a4676f.
//
// Solidity: function getRequestStatus(uint256 _requestId) view returns(bool fulfilled, uint256[] randomWords)
func (_DappLinkVRFManager *DappLinkVRFManagerSession) GetRequestStatus(_requestId *big.Int) (struct {
	Fulfilled   bool
	RandomWords []*big.Int
}, error) {
	return _DappLinkVRFManager.Contract.GetRequestStatus(&_DappLinkVRFManager.CallOpts, _requestId)
}

// GetRequestStatus is a free data retrieval call binding the contract method 0xd8a4676f.
//
// Solidity: function getRequestStatus(uint256 _requestId) view returns(bool fulfilled, uint256[] randomWords)
func (_DappLinkVRFManager *DappLinkVRFManagerCallerSession) GetRequestStatus(_requestId *big.Int) (struct {
	Fulfilled   bool
	RandomWords []*big.Int
}, error) {
	return _DappLinkVRFManager.Contract.GetRequestStatus(&_DappLinkVRFManager.CallOpts, _requestId)
}

// LastRequestId is a free data retrieval call binding the contract method 0xfc2a88c3.
//
// Solidity: function lastRequestId() view returns(uint256)
func (_DappLinkVRFManager *DappLinkVRFManagerCaller) LastRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DappLinkVRFManager.contract.Call(opts, &out, "lastRequestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastRequestId is a free data retrieval call binding the contract method 0xfc2a88c3.
//
// Solidity: function lastRequestId() view returns(uint256)
func (_DappLinkVRFManager *DappLinkVRFManagerSession) LastRequestId() (*big.Int, error) {
	return _DappLinkVRFManager.Contract.LastRequestId(&_DappLinkVRFManager.CallOpts)
}

// LastRequestId is a free data retrieval call binding the contract method 0xfc2a88c3.
//
// Solidity: function lastRequestId() view returns(uint256)
func (_DappLinkVRFManager *DappLinkVRFManagerCallerSession) LastRequestId() (*big.Int, error) {
	return _DappLinkVRFManager.Contract.LastRequestId(&_DappLinkVRFManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DappLinkVRFManager *DappLinkVRFManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DappLinkVRFManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DappLinkVRFManager *DappLinkVRFManagerSession) Owner() (common.Address, error) {
	return _DappLinkVRFManager.Contract.Owner(&_DappLinkVRFManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DappLinkVRFManager *DappLinkVRFManagerCallerSession) Owner() (common.Address, error) {
	return _DappLinkVRFManager.Contract.Owner(&_DappLinkVRFManager.CallOpts)
}

// RequestIds is a free data retrieval call binding the contract method 0x8796ba8c.
//
// Solidity: function requestIds(uint256 ) view returns(uint256)
func (_DappLinkVRFManager *DappLinkVRFManagerCaller) RequestIds(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _DappLinkVRFManager.contract.Call(opts, &out, "requestIds", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RequestIds is a free data retrieval call binding the contract method 0x8796ba8c.
//
// Solidity: function requestIds(uint256 ) view returns(uint256)
func (_DappLinkVRFManager *DappLinkVRFManagerSession) RequestIds(arg0 *big.Int) (*big.Int, error) {
	return _DappLinkVRFManager.Contract.RequestIds(&_DappLinkVRFManager.CallOpts, arg0)
}

// RequestIds is a free data retrieval call binding the contract method 0x8796ba8c.
//
// Solidity: function requestIds(uint256 ) view returns(uint256)
func (_DappLinkVRFManager *DappLinkVRFManagerCallerSession) RequestIds(arg0 *big.Int) (*big.Int, error) {
	return _DappLinkVRFManager.Contract.RequestIds(&_DappLinkVRFManager.CallOpts, arg0)
}

// RequestMapping is a free data retrieval call binding the contract method 0x82e215ab.
//
// Solidity: function requestMapping(uint256 ) view returns(bool fulfilled)
func (_DappLinkVRFManager *DappLinkVRFManagerCaller) RequestMapping(opts *bind.CallOpts, arg0 *big.Int) (bool, error) {
	var out []interface{}
	err := _DappLinkVRFManager.contract.Call(opts, &out, "requestMapping", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// RequestMapping is a free data retrieval call binding the contract method 0x82e215ab.
//
// Solidity: function requestMapping(uint256 ) view returns(bool fulfilled)
func (_DappLinkVRFManager *DappLinkVRFManagerSession) RequestMapping(arg0 *big.Int) (bool, error) {
	return _DappLinkVRFManager.Contract.RequestMapping(&_DappLinkVRFManager.CallOpts, arg0)
}

// RequestMapping is a free data retrieval call binding the contract method 0x82e215ab.
//
// Solidity: function requestMapping(uint256 ) view returns(bool fulfilled)
func (_DappLinkVRFManager *DappLinkVRFManagerCallerSession) RequestMapping(arg0 *big.Int) (bool, error) {
	return _DappLinkVRFManager.Contract.RequestMapping(&_DappLinkVRFManager.CallOpts, arg0)
}

// FulfillRandomWords is a paid mutator transaction binding the contract method 0x1b565710.
//
// Solidity: function fulfillRandomWords(uint256 _requestId, uint256[] _randomWords, bytes32 msgHash, uint256 referenceBlockNumber, ((uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint256,uint256) params) returns()
func (_DappLinkVRFManager *DappLinkVRFManagerTransactor) FulfillRandomWords(opts *bind.TransactOpts, _requestId *big.Int, _randomWords []*big.Int, msgHash [32]byte, referenceBlockNumber *big.Int, params IBLSApkRegistryVrfNoSignerAndSignature) (*types.Transaction, error) {
	return _DappLinkVRFManager.contract.Transact(opts, "fulfillRandomWords", _requestId, _randomWords, msgHash, referenceBlockNumber, params)
}

// FulfillRandomWords is a paid mutator transaction binding the contract method 0x1b565710.
//
// Solidity: function fulfillRandomWords(uint256 _requestId, uint256[] _randomWords, bytes32 msgHash, uint256 referenceBlockNumber, ((uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint256,uint256) params) returns()
func (_DappLinkVRFManager *DappLinkVRFManagerSession) FulfillRandomWords(_requestId *big.Int, _randomWords []*big.Int, msgHash [32]byte, referenceBlockNumber *big.Int, params IBLSApkRegistryVrfNoSignerAndSignature) (*types.Transaction, error) {
	return _DappLinkVRFManager.Contract.FulfillRandomWords(&_DappLinkVRFManager.TransactOpts, _requestId, _randomWords, msgHash, referenceBlockNumber, params)
}

// FulfillRandomWords is a paid mutator transaction binding the contract method 0x1b565710.
//
// Solidity: function fulfillRandomWords(uint256 _requestId, uint256[] _randomWords, bytes32 msgHash, uint256 referenceBlockNumber, ((uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint256,uint256) params) returns()
func (_DappLinkVRFManager *DappLinkVRFManagerTransactorSession) FulfillRandomWords(_requestId *big.Int, _randomWords []*big.Int, msgHash [32]byte, referenceBlockNumber *big.Int, params IBLSApkRegistryVrfNoSignerAndSignature) (*types.Transaction, error) {
	return _DappLinkVRFManager.Contract.FulfillRandomWords(&_DappLinkVRFManager.TransactOpts, _requestId, _randomWords, msgHash, referenceBlockNumber, params)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address initialOwner, address _dappLinkAddress, address _blsRegistry) returns()
func (_DappLinkVRFManager *DappLinkVRFManagerTransactor) Initialize(opts *bind.TransactOpts, initialOwner common.Address, _dappLinkAddress common.Address, _blsRegistry common.Address) (*types.Transaction, error) {
	return _DappLinkVRFManager.contract.Transact(opts, "initialize", initialOwner, _dappLinkAddress, _blsRegistry)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address initialOwner, address _dappLinkAddress, address _blsRegistry) returns()
func (_DappLinkVRFManager *DappLinkVRFManagerSession) Initialize(initialOwner common.Address, _dappLinkAddress common.Address, _blsRegistry common.Address) (*types.Transaction, error) {
	return _DappLinkVRFManager.Contract.Initialize(&_DappLinkVRFManager.TransactOpts, initialOwner, _dappLinkAddress, _blsRegistry)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address initialOwner, address _dappLinkAddress, address _blsRegistry) returns()
func (_DappLinkVRFManager *DappLinkVRFManagerTransactorSession) Initialize(initialOwner common.Address, _dappLinkAddress common.Address, _blsRegistry common.Address) (*types.Transaction, error) {
	return _DappLinkVRFManager.Contract.Initialize(&_DappLinkVRFManager.TransactOpts, initialOwner, _dappLinkAddress, _blsRegistry)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DappLinkVRFManager *DappLinkVRFManagerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DappLinkVRFManager.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DappLinkVRFManager *DappLinkVRFManagerSession) RenounceOwnership() (*types.Transaction, error) {
	return _DappLinkVRFManager.Contract.RenounceOwnership(&_DappLinkVRFManager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DappLinkVRFManager *DappLinkVRFManagerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _DappLinkVRFManager.Contract.RenounceOwnership(&_DappLinkVRFManager.TransactOpts)
}

// RequestRandomWords is a paid mutator transaction binding the contract method 0x1b739ef1.
//
// Solidity: function requestRandomWords(uint256 _requestId, uint256 _numWords) returns()
func (_DappLinkVRFManager *DappLinkVRFManagerTransactor) RequestRandomWords(opts *bind.TransactOpts, _requestId *big.Int, _numWords *big.Int) (*types.Transaction, error) {
	return _DappLinkVRFManager.contract.Transact(opts, "requestRandomWords", _requestId, _numWords)
}

// RequestRandomWords is a paid mutator transaction binding the contract method 0x1b739ef1.
//
// Solidity: function requestRandomWords(uint256 _requestId, uint256 _numWords) returns()
func (_DappLinkVRFManager *DappLinkVRFManagerSession) RequestRandomWords(_requestId *big.Int, _numWords *big.Int) (*types.Transaction, error) {
	return _DappLinkVRFManager.Contract.RequestRandomWords(&_DappLinkVRFManager.TransactOpts, _requestId, _numWords)
}

// RequestRandomWords is a paid mutator transaction binding the contract method 0x1b739ef1.
//
// Solidity: function requestRandomWords(uint256 _requestId, uint256 _numWords) returns()
func (_DappLinkVRFManager *DappLinkVRFManagerTransactorSession) RequestRandomWords(_requestId *big.Int, _numWords *big.Int) (*types.Transaction, error) {
	return _DappLinkVRFManager.Contract.RequestRandomWords(&_DappLinkVRFManager.TransactOpts, _requestId, _numWords)
}

// SetDappLink is a paid mutator transaction binding the contract method 0x996869d0.
//
// Solidity: function setDappLink(address _dappLinkAddress) returns()
func (_DappLinkVRFManager *DappLinkVRFManagerTransactor) SetDappLink(opts *bind.TransactOpts, _dappLinkAddress common.Address) (*types.Transaction, error) {
	return _DappLinkVRFManager.contract.Transact(opts, "setDappLink", _dappLinkAddress)
}

// SetDappLink is a paid mutator transaction binding the contract method 0x996869d0.
//
// Solidity: function setDappLink(address _dappLinkAddress) returns()
func (_DappLinkVRFManager *DappLinkVRFManagerSession) SetDappLink(_dappLinkAddress common.Address) (*types.Transaction, error) {
	return _DappLinkVRFManager.Contract.SetDappLink(&_DappLinkVRFManager.TransactOpts, _dappLinkAddress)
}

// SetDappLink is a paid mutator transaction binding the contract method 0x996869d0.
//
// Solidity: function setDappLink(address _dappLinkAddress) returns()
func (_DappLinkVRFManager *DappLinkVRFManagerTransactorSession) SetDappLink(_dappLinkAddress common.Address) (*types.Transaction, error) {
	return _DappLinkVRFManager.Contract.SetDappLink(&_DappLinkVRFManager.TransactOpts, _dappLinkAddress)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DappLinkVRFManager *DappLinkVRFManagerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _DappLinkVRFManager.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DappLinkVRFManager *DappLinkVRFManagerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DappLinkVRFManager.Contract.TransferOwnership(&_DappLinkVRFManager.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DappLinkVRFManager *DappLinkVRFManagerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DappLinkVRFManager.Contract.TransferOwnership(&_DappLinkVRFManager.TransactOpts, newOwner)
}

// DappLinkVRFManagerFillRandomWordsIterator is returned from FilterFillRandomWords and is used to iterate over the raw logs and unpacked data for FillRandomWords events raised by the DappLinkVRFManager contract.
type DappLinkVRFManagerFillRandomWordsIterator struct {
	Event *DappLinkVRFManagerFillRandomWords // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DappLinkVRFManagerFillRandomWordsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DappLinkVRFManagerFillRandomWords)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DappLinkVRFManagerFillRandomWords)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DappLinkVRFManagerFillRandomWordsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DappLinkVRFManagerFillRandomWordsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DappLinkVRFManagerFillRandomWords represents a FillRandomWords event raised by the DappLinkVRFManager contract.
type DappLinkVRFManagerFillRandomWords struct {
	RequestId   *big.Int
	RandomWords []*big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterFillRandomWords is a free log retrieval operation binding the contract event 0xf3cb4deb0441dd096356debf166f879d78cadc19e4b94053c8bea6d3940de93a.
//
// Solidity: event FillRandomWords(uint256 requestId, uint256[] randomWords)
func (_DappLinkVRFManager *DappLinkVRFManagerFilterer) FilterFillRandomWords(opts *bind.FilterOpts) (*DappLinkVRFManagerFillRandomWordsIterator, error) {

	logs, sub, err := _DappLinkVRFManager.contract.FilterLogs(opts, "FillRandomWords")
	if err != nil {
		return nil, err
	}
	return &DappLinkVRFManagerFillRandomWordsIterator{contract: _DappLinkVRFManager.contract, event: "FillRandomWords", logs: logs, sub: sub}, nil
}

// WatchFillRandomWords is a free log subscription operation binding the contract event 0xf3cb4deb0441dd096356debf166f879d78cadc19e4b94053c8bea6d3940de93a.
//
// Solidity: event FillRandomWords(uint256 requestId, uint256[] randomWords)
func (_DappLinkVRFManager *DappLinkVRFManagerFilterer) WatchFillRandomWords(opts *bind.WatchOpts, sink chan<- *DappLinkVRFManagerFillRandomWords) (event.Subscription, error) {

	logs, sub, err := _DappLinkVRFManager.contract.WatchLogs(opts, "FillRandomWords")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DappLinkVRFManagerFillRandomWords)
				if err := _DappLinkVRFManager.contract.UnpackLog(event, "FillRandomWords", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFillRandomWords is a log parse operation binding the contract event 0xf3cb4deb0441dd096356debf166f879d78cadc19e4b94053c8bea6d3940de93a.
//
// Solidity: event FillRandomWords(uint256 requestId, uint256[] randomWords)
func (_DappLinkVRFManager *DappLinkVRFManagerFilterer) ParseFillRandomWords(log types.Log) (*DappLinkVRFManagerFillRandomWords, error) {
	event := new(DappLinkVRFManagerFillRandomWords)
	if err := _DappLinkVRFManager.contract.UnpackLog(event, "FillRandomWords", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DappLinkVRFManagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the DappLinkVRFManager contract.
type DappLinkVRFManagerInitializedIterator struct {
	Event *DappLinkVRFManagerInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DappLinkVRFManagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DappLinkVRFManagerInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DappLinkVRFManagerInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DappLinkVRFManagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DappLinkVRFManagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DappLinkVRFManagerInitialized represents a Initialized event raised by the DappLinkVRFManager contract.
type DappLinkVRFManagerInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_DappLinkVRFManager *DappLinkVRFManagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*DappLinkVRFManagerInitializedIterator, error) {

	logs, sub, err := _DappLinkVRFManager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &DappLinkVRFManagerInitializedIterator{contract: _DappLinkVRFManager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_DappLinkVRFManager *DappLinkVRFManagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *DappLinkVRFManagerInitialized) (event.Subscription, error) {

	logs, sub, err := _DappLinkVRFManager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DappLinkVRFManagerInitialized)
				if err := _DappLinkVRFManager.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_DappLinkVRFManager *DappLinkVRFManagerFilterer) ParseInitialized(log types.Log) (*DappLinkVRFManagerInitialized, error) {
	event := new(DappLinkVRFManagerInitialized)
	if err := _DappLinkVRFManager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DappLinkVRFManagerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the DappLinkVRFManager contract.
type DappLinkVRFManagerOwnershipTransferredIterator struct {
	Event *DappLinkVRFManagerOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DappLinkVRFManagerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DappLinkVRFManagerOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DappLinkVRFManagerOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DappLinkVRFManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DappLinkVRFManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DappLinkVRFManagerOwnershipTransferred represents a OwnershipTransferred event raised by the DappLinkVRFManager contract.
type DappLinkVRFManagerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DappLinkVRFManager *DappLinkVRFManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DappLinkVRFManagerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DappLinkVRFManager.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DappLinkVRFManagerOwnershipTransferredIterator{contract: _DappLinkVRFManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DappLinkVRFManager *DappLinkVRFManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DappLinkVRFManagerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DappLinkVRFManager.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DappLinkVRFManagerOwnershipTransferred)
				if err := _DappLinkVRFManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DappLinkVRFManager *DappLinkVRFManagerFilterer) ParseOwnershipTransferred(log types.Log) (*DappLinkVRFManagerOwnershipTransferred, error) {
	event := new(DappLinkVRFManagerOwnershipTransferred)
	if err := _DappLinkVRFManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DappLinkVRFManagerRequestSentIterator is returned from FilterRequestSent and is used to iterate over the raw logs and unpacked data for RequestSent events raised by the DappLinkVRFManager contract.
type DappLinkVRFManagerRequestSentIterator struct {
	Event *DappLinkVRFManagerRequestSent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DappLinkVRFManagerRequestSentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DappLinkVRFManagerRequestSent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DappLinkVRFManagerRequestSent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DappLinkVRFManagerRequestSentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DappLinkVRFManagerRequestSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DappLinkVRFManagerRequestSent represents a RequestSent event raised by the DappLinkVRFManager contract.
type DappLinkVRFManagerRequestSent struct {
	RequestId *big.Int
	NumWords  *big.Int
	Current   common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRequestSent is a free log retrieval operation binding the contract event 0xe697eb68c0228bd7d4e553246a2a86e8402d0895e45092ef8ae87b4cfd29f016.
//
// Solidity: event RequestSent(uint256 requestId, uint256 _numWords, address current)
func (_DappLinkVRFManager *DappLinkVRFManagerFilterer) FilterRequestSent(opts *bind.FilterOpts) (*DappLinkVRFManagerRequestSentIterator, error) {

	logs, sub, err := _DappLinkVRFManager.contract.FilterLogs(opts, "RequestSent")
	if err != nil {
		return nil, err
	}
	return &DappLinkVRFManagerRequestSentIterator{contract: _DappLinkVRFManager.contract, event: "RequestSent", logs: logs, sub: sub}, nil
}

// WatchRequestSent is a free log subscription operation binding the contract event 0xe697eb68c0228bd7d4e553246a2a86e8402d0895e45092ef8ae87b4cfd29f016.
//
// Solidity: event RequestSent(uint256 requestId, uint256 _numWords, address current)
func (_DappLinkVRFManager *DappLinkVRFManagerFilterer) WatchRequestSent(opts *bind.WatchOpts, sink chan<- *DappLinkVRFManagerRequestSent) (event.Subscription, error) {

	logs, sub, err := _DappLinkVRFManager.contract.WatchLogs(opts, "RequestSent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DappLinkVRFManagerRequestSent)
				if err := _DappLinkVRFManager.contract.UnpackLog(event, "RequestSent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRequestSent is a log parse operation binding the contract event 0xe697eb68c0228bd7d4e553246a2a86e8402d0895e45092ef8ae87b4cfd29f016.
//
// Solidity: event RequestSent(uint256 requestId, uint256 _numWords, address current)
func (_DappLinkVRFManager *DappLinkVRFManagerFilterer) ParseRequestSent(log types.Log) (*DappLinkVRFManagerRequestSent, error) {
	event := new(DappLinkVRFManagerRequestSent)
	if err := _DappLinkVRFManager.contract.UnpackLog(event, "RequestSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
