// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ticTacToeApp

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ChannelAllocation is an auto generated low-level Go binding around an user-defined struct.
type ChannelAllocation struct {
	Assets   []common.Address
	Balances [][]*big.Int
	Locked   []ChannelSubAlloc
}

// ChannelParams is an auto generated low-level Go binding around an user-defined struct.
type ChannelParams struct {
	ChallengeDuration *big.Int
	Nonce             *big.Int
	Participants      []common.Address
	App               common.Address
	LedgerChannel     bool
	VirtualChannel    bool
}

// ChannelState is an auto generated low-level Go binding around an user-defined struct.
type ChannelState struct {
	ChannelID [32]byte
	Version   uint64
	Outcome   ChannelAllocation
	AppData   []byte
	IsFinal   bool
}

// ChannelSubAlloc is an auto generated low-level Go binding around an user-defined struct.
type ChannelSubAlloc struct {
	ID       [32]byte
	Balances []*big.Int
	IndexMap []uint16
}

// AppABI is the input ABI used to generate the binding from.
const AppABI = "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"challengeDuration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"app\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"ledgerChannel\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"virtualChannel\",\"type\":\"bool\"}],\"internalType\":\"structChannel.Params\",\"name\":\"params\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"channelID\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"address[]\",\"name\":\"assets\",\"type\":\"address[]\"},{\"internalType\":\"uint256[][]\",\"name\":\"balances\",\"type\":\"uint256[][]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ID\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"balances\",\"type\":\"uint256[]\"},{\"internalType\":\"uint16[]\",\"name\":\"indexMap\",\"type\":\"uint16[]\"}],\"internalType\":\"structChannel.SubAlloc[]\",\"name\":\"locked\",\"type\":\"tuple[]\"}],\"internalType\":\"structChannel.Allocation\",\"name\":\"outcome\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structChannel.State\",\"name\":\"from\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"channelID\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"address[]\",\"name\":\"assets\",\"type\":\"address[]\"},{\"internalType\":\"uint256[][]\",\"name\":\"balances\",\"type\":\"uint256[][]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ID\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"balances\",\"type\":\"uint256[]\"},{\"internalType\":\"uint16[]\",\"name\":\"indexMap\",\"type\":\"uint16[]\"}],\"internalType\":\"structChannel.SubAlloc[]\",\"name\":\"locked\",\"type\":\"tuple[]\"}],\"internalType\":\"structChannel.Allocation\",\"name\":\"outcome\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structChannel.State\",\"name\":\"to\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"actorIdx\",\"type\":\"uint256\"}],\"name\":\"validTransition\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"}]"

// AppFuncSigs maps the 4-byte function signature to its string representation.
var AppFuncSigs = map[string]string{
	"0d1feb4f": "validTransition((uint256,uint256,address[],address,bool,bool),(bytes32,uint64,(address[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool),(bytes32,uint64,(address[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool),uint256)",
}

// App is an auto generated Go binding around an Ethereum contract.
type App struct {
	AppCaller     // Read-only binding to the contract
	AppTransactor // Write-only binding to the contract
	AppFilterer   // Log filterer for contract events
}

// AppCaller is an auto generated read-only Go binding around an Ethereum contract.
type AppCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AppTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AppTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AppFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AppFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AppSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AppSession struct {
	Contract     *App              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AppCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AppCallerSession struct {
	Contract *AppCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// AppTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AppTransactorSession struct {
	Contract     *AppTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AppRaw is an auto generated low-level Go binding around an Ethereum contract.
type AppRaw struct {
	Contract *App // Generic contract binding to access the raw methods on
}

// AppCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AppCallerRaw struct {
	Contract *AppCaller // Generic read-only contract binding to access the raw methods on
}

// AppTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AppTransactorRaw struct {
	Contract *AppTransactor // Generic write-only contract binding to access the raw methods on
}

// NewApp creates a new instance of App, bound to a specific deployed contract.
func NewApp(address common.Address, backend bind.ContractBackend) (*App, error) {
	contract, err := bindApp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &App{AppCaller: AppCaller{contract: contract}, AppTransactor: AppTransactor{contract: contract}, AppFilterer: AppFilterer{contract: contract}}, nil
}

// NewAppCaller creates a new read-only instance of App, bound to a specific deployed contract.
func NewAppCaller(address common.Address, caller bind.ContractCaller) (*AppCaller, error) {
	contract, err := bindApp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AppCaller{contract: contract}, nil
}

// NewAppTransactor creates a new write-only instance of App, bound to a specific deployed contract.
func NewAppTransactor(address common.Address, transactor bind.ContractTransactor) (*AppTransactor, error) {
	contract, err := bindApp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AppTransactor{contract: contract}, nil
}

// NewAppFilterer creates a new log filterer instance of App, bound to a specific deployed contract.
func NewAppFilterer(address common.Address, filterer bind.ContractFilterer) (*AppFilterer, error) {
	contract, err := bindApp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AppFilterer{contract: contract}, nil
}

// bindApp binds a generic wrapper to an already deployed contract.
func bindApp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AppABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_App *AppRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _App.Contract.AppCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_App *AppRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _App.Contract.AppTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_App *AppRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _App.Contract.AppTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_App *AppCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _App.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_App *AppTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _App.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_App *AppTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _App.Contract.contract.Transact(opts, method, params...)
}

// ValidTransition is a free data retrieval call binding the contract method 0x0d1feb4f.
//
// Solidity: function validTransition((uint256,uint256,address[],address,bool,bool) params, (bytes32,uint64,(address[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) from, (bytes32,uint64,(address[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) to, uint256 actorIdx) pure returns()
func (_App *AppCaller) ValidTransition(opts *bind.CallOpts, params ChannelParams, from ChannelState, to ChannelState, actorIdx *big.Int) error {
	var out []interface{}
	err := _App.contract.Call(opts, &out, "validTransition", params, from, to, actorIdx)

	if err != nil {
		return err
	}

	return err

}

// ValidTransition is a free data retrieval call binding the contract method 0x0d1feb4f.
//
// Solidity: function validTransition((uint256,uint256,address[],address,bool,bool) params, (bytes32,uint64,(address[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) from, (bytes32,uint64,(address[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) to, uint256 actorIdx) pure returns()
func (_App *AppSession) ValidTransition(params ChannelParams, from ChannelState, to ChannelState, actorIdx *big.Int) error {
	return _App.Contract.ValidTransition(&_App.CallOpts, params, from, to, actorIdx)
}

// ValidTransition is a free data retrieval call binding the contract method 0x0d1feb4f.
//
// Solidity: function validTransition((uint256,uint256,address[],address,bool,bool) params, (bytes32,uint64,(address[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) from, (bytes32,uint64,(address[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) to, uint256 actorIdx) pure returns()
func (_App *AppCallerSession) ValidTransition(params ChannelParams, from ChannelState, to ChannelState, actorIdx *big.Int) error {
	return _App.Contract.ValidTransition(&_App.CallOpts, params, from, to, actorIdx)
}

// ArrayABI is the input ABI used to generate the binding from.
const ArrayABI = "[]"

// ArrayBin is the compiled bytecode used for deploying new contracts.
var ArrayBin = "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220ebef2fb342270f5254def9fcaa031ba4a02645194c01a10e870b53c1cc12351064736f6c63430007060033"

// DeployArray deploys a new Ethereum contract, binding an instance of Array to it.
func DeployArray(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Array, error) {
	parsed, err := abi.JSON(strings.NewReader(ArrayABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ArrayBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Array{ArrayCaller: ArrayCaller{contract: contract}, ArrayTransactor: ArrayTransactor{contract: contract}, ArrayFilterer: ArrayFilterer{contract: contract}}, nil
}

// Array is an auto generated Go binding around an Ethereum contract.
type Array struct {
	ArrayCaller     // Read-only binding to the contract
	ArrayTransactor // Write-only binding to the contract
	ArrayFilterer   // Log filterer for contract events
}

// ArrayCaller is an auto generated read-only Go binding around an Ethereum contract.
type ArrayCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ArrayTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ArrayTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ArrayFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ArrayFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ArraySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ArraySession struct {
	Contract     *Array            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ArrayCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ArrayCallerSession struct {
	Contract *ArrayCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ArrayTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ArrayTransactorSession struct {
	Contract     *ArrayTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ArrayRaw is an auto generated low-level Go binding around an Ethereum contract.
type ArrayRaw struct {
	Contract *Array // Generic contract binding to access the raw methods on
}

// ArrayCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ArrayCallerRaw struct {
	Contract *ArrayCaller // Generic read-only contract binding to access the raw methods on
}

// ArrayTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ArrayTransactorRaw struct {
	Contract *ArrayTransactor // Generic write-only contract binding to access the raw methods on
}

// NewArray creates a new instance of Array, bound to a specific deployed contract.
func NewArray(address common.Address, backend bind.ContractBackend) (*Array, error) {
	contract, err := bindArray(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Array{ArrayCaller: ArrayCaller{contract: contract}, ArrayTransactor: ArrayTransactor{contract: contract}, ArrayFilterer: ArrayFilterer{contract: contract}}, nil
}

// NewArrayCaller creates a new read-only instance of Array, bound to a specific deployed contract.
func NewArrayCaller(address common.Address, caller bind.ContractCaller) (*ArrayCaller, error) {
	contract, err := bindArray(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ArrayCaller{contract: contract}, nil
}

// NewArrayTransactor creates a new write-only instance of Array, bound to a specific deployed contract.
func NewArrayTransactor(address common.Address, transactor bind.ContractTransactor) (*ArrayTransactor, error) {
	contract, err := bindArray(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ArrayTransactor{contract: contract}, nil
}

// NewArrayFilterer creates a new log filterer instance of Array, bound to a specific deployed contract.
func NewArrayFilterer(address common.Address, filterer bind.ContractFilterer) (*ArrayFilterer, error) {
	contract, err := bindArray(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ArrayFilterer{contract: contract}, nil
}

// bindArray binds a generic wrapper to an already deployed contract.
func bindArray(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ArrayABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Array *ArrayRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Array.Contract.ArrayCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Array *ArrayRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Array.Contract.ArrayTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Array *ArrayRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Array.Contract.ArrayTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Array *ArrayCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Array.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Array *ArrayTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Array.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Array *ArrayTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Array.Contract.contract.Transact(opts, method, params...)
}

// ChannelABI is the input ABI used to generate the binding from.
const ChannelABI = "[]"

// ChannelBin is the compiled bytecode used for deploying new contracts.
var ChannelBin = "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea264697066735822122058368571959499bf5c002554352cc9b02e2980a95a85e4230bdcf4f08c4d604d64736f6c63430007060033"

// DeployChannel deploys a new Ethereum contract, binding an instance of Channel to it.
func DeployChannel(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Channel, error) {
	parsed, err := abi.JSON(strings.NewReader(ChannelABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ChannelBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Channel{ChannelCaller: ChannelCaller{contract: contract}, ChannelTransactor: ChannelTransactor{contract: contract}, ChannelFilterer: ChannelFilterer{contract: contract}}, nil
}

// Channel is an auto generated Go binding around an Ethereum contract.
type Channel struct {
	ChannelCaller     // Read-only binding to the contract
	ChannelTransactor // Write-only binding to the contract
	ChannelFilterer   // Log filterer for contract events
}

// ChannelCaller is an auto generated read-only Go binding around an Ethereum contract.
type ChannelCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChannelTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ChannelTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChannelFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ChannelFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChannelSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ChannelSession struct {
	Contract     *Channel          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ChannelCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ChannelCallerSession struct {
	Contract *ChannelCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// ChannelTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ChannelTransactorSession struct {
	Contract     *ChannelTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ChannelRaw is an auto generated low-level Go binding around an Ethereum contract.
type ChannelRaw struct {
	Contract *Channel // Generic contract binding to access the raw methods on
}

// ChannelCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ChannelCallerRaw struct {
	Contract *ChannelCaller // Generic read-only contract binding to access the raw methods on
}

// ChannelTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ChannelTransactorRaw struct {
	Contract *ChannelTransactor // Generic write-only contract binding to access the raw methods on
}

// NewChannel creates a new instance of Channel, bound to a specific deployed contract.
func NewChannel(address common.Address, backend bind.ContractBackend) (*Channel, error) {
	contract, err := bindChannel(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Channel{ChannelCaller: ChannelCaller{contract: contract}, ChannelTransactor: ChannelTransactor{contract: contract}, ChannelFilterer: ChannelFilterer{contract: contract}}, nil
}

// NewChannelCaller creates a new read-only instance of Channel, bound to a specific deployed contract.
func NewChannelCaller(address common.Address, caller bind.ContractCaller) (*ChannelCaller, error) {
	contract, err := bindChannel(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChannelCaller{contract: contract}, nil
}

// NewChannelTransactor creates a new write-only instance of Channel, bound to a specific deployed contract.
func NewChannelTransactor(address common.Address, transactor bind.ContractTransactor) (*ChannelTransactor, error) {
	contract, err := bindChannel(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChannelTransactor{contract: contract}, nil
}

// NewChannelFilterer creates a new log filterer instance of Channel, bound to a specific deployed contract.
func NewChannelFilterer(address common.Address, filterer bind.ContractFilterer) (*ChannelFilterer, error) {
	contract, err := bindChannel(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChannelFilterer{contract: contract}, nil
}

// bindChannel binds a generic wrapper to an already deployed contract.
func bindChannel(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ChannelABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Channel *ChannelRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Channel.Contract.ChannelCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Channel *ChannelRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Channel.Contract.ChannelTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Channel *ChannelRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Channel.Contract.ChannelTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Channel *ChannelCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Channel.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Channel *ChannelTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Channel.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Channel *ChannelTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Channel.Contract.contract.Transact(opts, method, params...)
}

// ECDSAABI is the input ABI used to generate the binding from.
const ECDSAABI = "[]"

// ECDSABin is the compiled bytecode used for deploying new contracts.
var ECDSABin = "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220bec3b453040fee1745af167d578c2ea9fa390f85a7ebf8bc3a249e57813f48bd64736f6c63430007060033"

// DeployECDSA deploys a new Ethereum contract, binding an instance of ECDSA to it.
func DeployECDSA(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ECDSA, error) {
	parsed, err := abi.JSON(strings.NewReader(ECDSAABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ECDSABin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ECDSA{ECDSACaller: ECDSACaller{contract: contract}, ECDSATransactor: ECDSATransactor{contract: contract}, ECDSAFilterer: ECDSAFilterer{contract: contract}}, nil
}

// ECDSA is an auto generated Go binding around an Ethereum contract.
type ECDSA struct {
	ECDSACaller     // Read-only binding to the contract
	ECDSATransactor // Write-only binding to the contract
	ECDSAFilterer   // Log filterer for contract events
}

// ECDSACaller is an auto generated read-only Go binding around an Ethereum contract.
type ECDSACaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ECDSATransactor is an auto generated write-only Go binding around an Ethereum contract.
type ECDSATransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ECDSAFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ECDSAFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ECDSASession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ECDSASession struct {
	Contract     *ECDSA            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ECDSACallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ECDSACallerSession struct {
	Contract *ECDSACaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ECDSATransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ECDSATransactorSession struct {
	Contract     *ECDSATransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ECDSARaw is an auto generated low-level Go binding around an Ethereum contract.
type ECDSARaw struct {
	Contract *ECDSA // Generic contract binding to access the raw methods on
}

// ECDSACallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ECDSACallerRaw struct {
	Contract *ECDSACaller // Generic read-only contract binding to access the raw methods on
}

// ECDSATransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ECDSATransactorRaw struct {
	Contract *ECDSATransactor // Generic write-only contract binding to access the raw methods on
}

// NewECDSA creates a new instance of ECDSA, bound to a specific deployed contract.
func NewECDSA(address common.Address, backend bind.ContractBackend) (*ECDSA, error) {
	contract, err := bindECDSA(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ECDSA{ECDSACaller: ECDSACaller{contract: contract}, ECDSATransactor: ECDSATransactor{contract: contract}, ECDSAFilterer: ECDSAFilterer{contract: contract}}, nil
}

// NewECDSACaller creates a new read-only instance of ECDSA, bound to a specific deployed contract.
func NewECDSACaller(address common.Address, caller bind.ContractCaller) (*ECDSACaller, error) {
	contract, err := bindECDSA(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ECDSACaller{contract: contract}, nil
}

// NewECDSATransactor creates a new write-only instance of ECDSA, bound to a specific deployed contract.
func NewECDSATransactor(address common.Address, transactor bind.ContractTransactor) (*ECDSATransactor, error) {
	contract, err := bindECDSA(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ECDSATransactor{contract: contract}, nil
}

// NewECDSAFilterer creates a new log filterer instance of ECDSA, bound to a specific deployed contract.
func NewECDSAFilterer(address common.Address, filterer bind.ContractFilterer) (*ECDSAFilterer, error) {
	contract, err := bindECDSA(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ECDSAFilterer{contract: contract}, nil
}

// bindECDSA binds a generic wrapper to an already deployed contract.
func bindECDSA(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ECDSAABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ECDSA *ECDSARaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ECDSA.Contract.ECDSACaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ECDSA *ECDSARaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ECDSA.Contract.ECDSATransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ECDSA *ECDSARaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ECDSA.Contract.ECDSATransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ECDSA *ECDSACallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ECDSA.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ECDSA *ECDSATransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ECDSA.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ECDSA *ECDSATransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ECDSA.Contract.contract.Transact(opts, method, params...)
}

// SafeMathABI is the input ABI used to generate the binding from.
const SafeMathABI = "[]"

// SafeMathBin is the compiled bytecode used for deploying new contracts.
var SafeMathBin = "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220b5a22cf2ed655f097bcfcbecbac6ad233327b0f9bde3535f9a52829c5c1114fe64736f6c63430007060033"

// DeploySafeMath deploys a new Ethereum contract, binding an instance of SafeMath to it.
func DeploySafeMath(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SafeMath, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SafeMathBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// SafeMath is an auto generated Go binding around an Ethereum contract.
type SafeMath struct {
	SafeMathCaller     // Read-only binding to the contract
	SafeMathTransactor // Write-only binding to the contract
	SafeMathFilterer   // Log filterer for contract events
}

// SafeMathCaller is an auto generated read-only Go binding around an Ethereum contract.
type SafeMathCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SafeMathTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SafeMathFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SafeMathSession struct {
	Contract     *SafeMath         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafeMathCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SafeMathCallerSession struct {
	Contract *SafeMathCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// SafeMathTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SafeMathTransactorSession struct {
	Contract     *SafeMathTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// SafeMathRaw is an auto generated low-level Go binding around an Ethereum contract.
type SafeMathRaw struct {
	Contract *SafeMath // Generic contract binding to access the raw methods on
}

// SafeMathCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SafeMathCallerRaw struct {
	Contract *SafeMathCaller // Generic read-only contract binding to access the raw methods on
}

// SafeMathTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SafeMathTransactorRaw struct {
	Contract *SafeMathTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSafeMath creates a new instance of SafeMath, bound to a specific deployed contract.
func NewSafeMath(address common.Address, backend bind.ContractBackend) (*SafeMath, error) {
	contract, err := bindSafeMath(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// NewSafeMathCaller creates a new read-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathCaller(address common.Address, caller bind.ContractCaller) (*SafeMathCaller, error) {
	contract, err := bindSafeMath(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathCaller{contract: contract}, nil
}

// NewSafeMathTransactor creates a new write-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathTransactor(address common.Address, transactor bind.ContractTransactor) (*SafeMathTransactor, error) {
	contract, err := bindSafeMath(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathTransactor{contract: contract}, nil
}

// NewSafeMathFilterer creates a new log filterer instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathFilterer(address common.Address, filterer bind.ContractFilterer) (*SafeMathFilterer, error) {
	contract, err := bindSafeMath(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafeMathFilterer{contract: contract}, nil
}

// bindSafeMath binds a generic wrapper to an already deployed contract.
func bindSafeMath(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.SafeMathCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transact(opts, method, params...)
}

// SigABI is the input ABI used to generate the binding from.
const SigABI = "[]"

// SigBin is the compiled bytecode used for deploying new contracts.
var SigBin = "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220e7eab097e26aa1458f5efa6ba68afa258230754eba24b22ee9b488deea8c50cf64736f6c63430007060033"

// DeploySig deploys a new Ethereum contract, binding an instance of Sig to it.
func DeploySig(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Sig, error) {
	parsed, err := abi.JSON(strings.NewReader(SigABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SigBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Sig{SigCaller: SigCaller{contract: contract}, SigTransactor: SigTransactor{contract: contract}, SigFilterer: SigFilterer{contract: contract}}, nil
}

// Sig is an auto generated Go binding around an Ethereum contract.
type Sig struct {
	SigCaller     // Read-only binding to the contract
	SigTransactor // Write-only binding to the contract
	SigFilterer   // Log filterer for contract events
}

// SigCaller is an auto generated read-only Go binding around an Ethereum contract.
type SigCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SigTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SigTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SigFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SigFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SigSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SigSession struct {
	Contract     *Sig              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SigCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SigCallerSession struct {
	Contract *SigCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// SigTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SigTransactorSession struct {
	Contract     *SigTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SigRaw is an auto generated low-level Go binding around an Ethereum contract.
type SigRaw struct {
	Contract *Sig // Generic contract binding to access the raw methods on
}

// SigCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SigCallerRaw struct {
	Contract *SigCaller // Generic read-only contract binding to access the raw methods on
}

// SigTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SigTransactorRaw struct {
	Contract *SigTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSig creates a new instance of Sig, bound to a specific deployed contract.
func NewSig(address common.Address, backend bind.ContractBackend) (*Sig, error) {
	contract, err := bindSig(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Sig{SigCaller: SigCaller{contract: contract}, SigTransactor: SigTransactor{contract: contract}, SigFilterer: SigFilterer{contract: contract}}, nil
}

// NewSigCaller creates a new read-only instance of Sig, bound to a specific deployed contract.
func NewSigCaller(address common.Address, caller bind.ContractCaller) (*SigCaller, error) {
	contract, err := bindSig(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SigCaller{contract: contract}, nil
}

// NewSigTransactor creates a new write-only instance of Sig, bound to a specific deployed contract.
func NewSigTransactor(address common.Address, transactor bind.ContractTransactor) (*SigTransactor, error) {
	contract, err := bindSig(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SigTransactor{contract: contract}, nil
}

// NewSigFilterer creates a new log filterer instance of Sig, bound to a specific deployed contract.
func NewSigFilterer(address common.Address, filterer bind.ContractFilterer) (*SigFilterer, error) {
	contract, err := bindSig(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SigFilterer{contract: contract}, nil
}

// bindSig binds a generic wrapper to an already deployed contract.
func bindSig(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SigABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Sig *SigRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Sig.Contract.SigCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Sig *SigRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Sig.Contract.SigTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Sig *SigRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Sig.Contract.SigTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Sig *SigCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Sig.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Sig *SigTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Sig.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Sig *SigTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Sig.Contract.contract.Transact(opts, method, params...)
}

// TicTacToeAppABI is the input ABI used to generate the binding from.
const TicTacToeAppABI = "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"challengeDuration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"app\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"ledgerChannel\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"virtualChannel\",\"type\":\"bool\"}],\"internalType\":\"structChannel.Params\",\"name\":\"params\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"channelID\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"address[]\",\"name\":\"assets\",\"type\":\"address[]\"},{\"internalType\":\"uint256[][]\",\"name\":\"balances\",\"type\":\"uint256[][]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ID\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"balances\",\"type\":\"uint256[]\"},{\"internalType\":\"uint16[]\",\"name\":\"indexMap\",\"type\":\"uint16[]\"}],\"internalType\":\"structChannel.SubAlloc[]\",\"name\":\"locked\",\"type\":\"tuple[]\"}],\"internalType\":\"structChannel.Allocation\",\"name\":\"outcome\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structChannel.State\",\"name\":\"from\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"channelID\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"address[]\",\"name\":\"assets\",\"type\":\"address[]\"},{\"internalType\":\"uint256[][]\",\"name\":\"balances\",\"type\":\"uint256[][]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ID\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"balances\",\"type\":\"uint256[]\"},{\"internalType\":\"uint16[]\",\"name\":\"indexMap\",\"type\":\"uint16[]\"}],\"internalType\":\"structChannel.SubAlloc[]\",\"name\":\"locked\",\"type\":\"tuple[]\"}],\"internalType\":\"structChannel.Allocation\",\"name\":\"outcome\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structChannel.State\",\"name\":\"to\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"signerIdx\",\"type\":\"uint256\"}],\"name\":\"validTransition\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"}]"

// TicTacToeAppFuncSigs maps the 4-byte function signature to its string representation.
var TicTacToeAppFuncSigs = map[string]string{
	"0d1feb4f": "validTransition((uint256,uint256,address[],address,bool,bool),(bytes32,uint64,(address[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool),(bytes32,uint64,(address[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool),uint256)",
}

// TicTacToeAppBin is the compiled bytecode used for deploying new contracts.
var TicTacToeAppBin = "0x608060405234801561001057600080fd5b5061120a806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c80630d1feb4f14610030575b600080fd5b61004361003e366004610d57565b610045565b005b60026100546040860186610fe1565b90501461007c5760405162461bcd60e51b815260040161007390610f55565b60405180910390fd5b600061008b6060850185611028565b60008161009457fe5b919091013560f81c9150600a90506100af6060850185611028565b9050146100ce5760405162461bcd60e51b815260040161007390610f85565b818160ff16146100f05760405162461bcd60e51b815260040161007390610e89565b6100fd6060840184611028565b60008161010657fe5b9091013560f81c6001838101161490506101325760405162461bcd60e51b815260040161007390610ed6565b600060015b600a81101561024f57600261014f6060870187611028565b8381811061015957fe5b9050013560f81c60f81b60f81c60ff1611156101875760405162461bcd60e51b815260040161007390610e37565b6101946060870187611028565b8281811061019e57fe5b909101356001600160f81b03191690506101bb6060870187611028565b838181106101c557fe5b9050013560f81c60f81b6001600160f81b031916146102475760006101ed6060880188611028565b838181106101f757fe5b9050013560f81c60f81b60f81c60ff16146102245760405162461bcd60e51b815260040161007390610eb3565b81156102425760405162461bcd60e51b815260040161007390610dea565b600191505b600101610137565b506000808061029e6102646060890189611028565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506105dd92505050565b919450925090508215156102b860a0890160808a01610d30565b1515146102d75760405162461bcd60e51b815260040161007390610efa565b6103716102e7604089018961106d565b6102f19080610fe1565b808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152506103309250505060408b018b61106d565b61033a9080610fe1565b808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152506107c592505050565b6103c1610381604089018961106d565b61038f906040810190610fe1565b6103989161111c565b6103a560408b018b61106d565b6103b3906040810190610fe1565b6103bc9161111c565b6108c0565b60006103d060408a018a61106d565b6103de906020810190610fe1565b6103e7916110ce565b905082156105a45760ff8216156104105760405162461bcd60e51b815260040161007390610e0f565b805160018390039067ffffffffffffffff8111801561042e57600080fd5b5060405190808252806020026020018201604052801561046257816020015b606081526020019060019003908161044d5790505b50915060005b82518110156105a15760408051600280825260608201835290916020830190803683370190505083828151811061049b57fe5b60209081029190910101526104b360408c018c61106d565b6104c1906020810190610fe1565b828181106104cb57fe5b90506020028101906104dd9190610fe1565b60018181106104e857fe5b905060200201358b80604001906104ff919061106d565b61050d906020810190610fe1565b8381811061051757fe5b90506020028101906105299190610fe1565b600081811061053457fe5b905060200201350183828151811061054857fe5b60200260200101518560ff168151811061055e57fe5b602002602001018181525050600083828151811061057857fe5b60200260200101518360ff168151811061058e57fe5b6020908102919091010152600101610468565b50505b6105d16105b460408a018a61106d565b6105c2906020810190610fe1565b6105cb916110ce565b82610925565b50505050505050505050565b60408051610160810182526000610100820181815260016101208401819052600261014085018190529184528451606081810187526003808352600460208481018290526005858b01819052818a019590955289518085018b52600680825260078284018190526008838e018190528c8e01939093528c518088018e528b815280850196909652858d01829052868c01959095528b518087018d52978852878301849052878c019490945260808a019690965289518085018b5287815280820195909552848a0186905260a089019490945288518084018a52878152808501829052808a019590955260c088019490945287519182018852938152908101919091529384015260e082019290925281908190815b600881101561076b576000806107178885856008811061070d57fe5b602002015161098a565b9150915081156107615760ff811660011415610741576001806000965096509650505050506107be565b60ff811660021415610761576001806001965096509650505050506107be565b50506001016106f1565b5060005b85518110156107b057600060ff1686828151811061078957fe5b016020015160f81c146107a857600080600094509450945050506107be565b60010161076f565b506001600080935093509350505b9193909250565b805182511461081b576040805162461bcd60e51b815260206004820152601960248201527f616464726573735b5d3a20756e657175616c206c656e67746800000000000000604482015290519081900360640190fd5b60005b82518110156108bb5781818151811061083357fe5b60200260200101516001600160a01b031683828151811061085057fe5b60200260200101516001600160a01b0316146108b3576040805162461bcd60e51b815260206004820152601760248201527f616464726573735b5d3a20756e657175616c206974656d000000000000000000604482015290519081900360640190fd5b60010161081e565b505050565b80518251146108e15760405162461bcd60e51b815260040161007390610faa565b60005b82518110156108bb5761091d8382815181106108fc57fe5b602002602001015183838151811061091057fe5b6020026020010151610a32565b6001016108e4565b80518251146109465760405162461bcd60e51b815260040161007390610f1e565b60005b82518110156108bb5761098283828151811061096157fe5b602002602001015183838151811061097557fe5b6020026020010151610a7b565b600101610949565b60008080848482602002015160016000010160ff16815181106109a957fe5b01602001516001600160f81b031916905060015b6003811015610a20576001600160f81b03198216868683600381106109de57fe5b602002015160016000010160ff16815181106109f657fe5b01602001516001600160f81b03191614610a1857600080935093505050610a2b565b6001016109bd565b506001925060f81c90505b9250929050565b8051825114610a535760405162461bcd60e51b815260040161007390610e5b565b610a6582602001518260200151610a7b565b610a7782604001518260400151610b5f565b5050565b8051825114610ad1576040805162461bcd60e51b815260206004820152601960248201527f75696e743235365b5d3a20756e657175616c206c656e67746800000000000000604482015290519081900360640190fd5b60005b82518110156108bb57818181518110610ae957fe5b6020026020010151838281518110610afd57fe5b602002602001015114610b57576040805162461bcd60e51b815260206004820152601760248201527f75696e743235365b5d3a20756e657175616c206974656d000000000000000000604482015290519081900360640190fd5b600101610ad4565b8051825114610bb5576040805162461bcd60e51b815260206004820152601860248201527f75696e7431365b5d3a20756e657175616c206c656e6774680000000000000000604482015290519081900360640190fd5b60005b82518110156108bb57818181518110610bcd57fe5b602002602001015161ffff16838281518110610be557fe5b602002602001015161ffff1614610c3c576040805162461bcd60e51b815260206004820152601660248201527575696e7431365b5d3a20756e657175616c206974656d60501b604482015290519081900360640190fd5b600101610bb8565b600082601f830112610c54578081fd5b81356020610c69610c64836110b0565b61108c565b8281528181019085830183850287018401881015610c85578586fd5b855b85811015610cb257813561ffff81168114610ca0578788fd5b84529284019290840190600101610c87565b5090979650505050505050565b600082601f830112610ccf578081fd5b81356020610cdf610c64836110b0565b8281528181019085830183850287018401881015610cfb578586fd5b855b85811015610cb257813584529284019290840190600101610cfd565b600060a08284031215610d2a578081fd5b50919050565b600060208284031215610d41578081fd5b81358015158114610d50578182fd5b9392505050565b60008060008060808587031215610d6c578283fd5b843567ffffffffffffffff80821115610d83578485fd5b9086019060c08289031215610d96578485fd5b90945060208601359080821115610dab578485fd5b610db788838901610d19565b94506040870135915080821115610dcc578384fd5b50610dd987828801610d19565b949793965093946060013593505050565b6020808252600b908201526a74776f20616374696f6e7360a81b604082015260600190565b6020808252600e908201526d06861732077696e6e657220213d360941b604082015260600190565b6020808252600a9082015269677269642076616c756560b01b604082015260600190565b60208082526014908201527314dd58905b1b1bd8ce881d5b995c5d585b08125160621b604082015260600190565b60208082526010908201526f30b1ba37b9103737ba1039b4b3b732b960811b604082015260600190565b6020808252600990820152686f766572777269746560b81b604082015260600190565b6020808252600a90820152693732bc3a1030b1ba37b960b11b604082015260600190565b6020808252600a908201526966696e616c20666c616760b01b604082015260600190565b6020808252601b908201527f75696e743235365b5d5b5d3a20756e657175616c206c656e6774680000000000604082015260600190565b6020808252601690820152756e756d626572206f66207061727469636970616e747360501b604082015260600190565b6020808252600b908201526a0c8c2e8c240d8cadccee8d60ab1b604082015260600190565b6020808252601a908201527f537562416c6c6f635b5d3a20756e657175616c206c656e677468000000000000604082015260600190565b6000808335601e19843603018112610ff7578283fd5b83018035915067ffffffffffffffff821115611011578283fd5b6020908101925081023603821315610a2b57600080fd5b6000808335601e1984360301811261103e578283fd5b83018035915067ffffffffffffffff821115611058578283fd5b602001915036819003821315610a2b57600080fd5b60008235605e19833603018112611082578182fd5b9190910192915050565b60405181810167ffffffffffffffff811182821017156110a857fe5b604052919050565b600067ffffffffffffffff8211156110c457fe5b5060209081020190565b60006110dc610c64846110b0565b8381526020808201919084845b87811015611110576110fe3683358901610cbf565b855293820193908201906001016110e9565b50919695505050505050565b600061112a610c64846110b0565b8381526020808201919084845b8781101561111057813587016060808236031215611153578788fd5b604080519182019167ffffffffffffffff808411828510171561117257fe5b9282528335815286840135928084111561118a578a8bfd5b61119636858701610cbf565b88830152828501359350808411156111ac578a8bfd5b506111b936848601610c44565b9181019190915287525050938201939082019060010161113756fea264697066735822122083747e473efd745af0ea3c4af2021e679f44172fc1297a01049044aea209a9ff64736f6c63430007060033"

// DeployTicTacToeApp deploys a new Ethereum contract, binding an instance of TicTacToeApp to it.
func DeployTicTacToeApp(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TicTacToeApp, error) {
	parsed, err := abi.JSON(strings.NewReader(TicTacToeAppABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TicTacToeAppBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TicTacToeApp{TicTacToeAppCaller: TicTacToeAppCaller{contract: contract}, TicTacToeAppTransactor: TicTacToeAppTransactor{contract: contract}, TicTacToeAppFilterer: TicTacToeAppFilterer{contract: contract}}, nil
}

// TicTacToeApp is an auto generated Go binding around an Ethereum contract.
type TicTacToeApp struct {
	TicTacToeAppCaller     // Read-only binding to the contract
	TicTacToeAppTransactor // Write-only binding to the contract
	TicTacToeAppFilterer   // Log filterer for contract events
}

// TicTacToeAppCaller is an auto generated read-only Go binding around an Ethereum contract.
type TicTacToeAppCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TicTacToeAppTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TicTacToeAppTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TicTacToeAppFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TicTacToeAppFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TicTacToeAppSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TicTacToeAppSession struct {
	Contract     *TicTacToeApp     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TicTacToeAppCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TicTacToeAppCallerSession struct {
	Contract *TicTacToeAppCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// TicTacToeAppTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TicTacToeAppTransactorSession struct {
	Contract     *TicTacToeAppTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// TicTacToeAppRaw is an auto generated low-level Go binding around an Ethereum contract.
type TicTacToeAppRaw struct {
	Contract *TicTacToeApp // Generic contract binding to access the raw methods on
}

// TicTacToeAppCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TicTacToeAppCallerRaw struct {
	Contract *TicTacToeAppCaller // Generic read-only contract binding to access the raw methods on
}

// TicTacToeAppTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TicTacToeAppTransactorRaw struct {
	Contract *TicTacToeAppTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTicTacToeApp creates a new instance of TicTacToeApp, bound to a specific deployed contract.
func NewTicTacToeApp(address common.Address, backend bind.ContractBackend) (*TicTacToeApp, error) {
	contract, err := bindTicTacToeApp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TicTacToeApp{TicTacToeAppCaller: TicTacToeAppCaller{contract: contract}, TicTacToeAppTransactor: TicTacToeAppTransactor{contract: contract}, TicTacToeAppFilterer: TicTacToeAppFilterer{contract: contract}}, nil
}

// NewTicTacToeAppCaller creates a new read-only instance of TicTacToeApp, bound to a specific deployed contract.
func NewTicTacToeAppCaller(address common.Address, caller bind.ContractCaller) (*TicTacToeAppCaller, error) {
	contract, err := bindTicTacToeApp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TicTacToeAppCaller{contract: contract}, nil
}

// NewTicTacToeAppTransactor creates a new write-only instance of TicTacToeApp, bound to a specific deployed contract.
func NewTicTacToeAppTransactor(address common.Address, transactor bind.ContractTransactor) (*TicTacToeAppTransactor, error) {
	contract, err := bindTicTacToeApp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TicTacToeAppTransactor{contract: contract}, nil
}

// NewTicTacToeAppFilterer creates a new log filterer instance of TicTacToeApp, bound to a specific deployed contract.
func NewTicTacToeAppFilterer(address common.Address, filterer bind.ContractFilterer) (*TicTacToeAppFilterer, error) {
	contract, err := bindTicTacToeApp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TicTacToeAppFilterer{contract: contract}, nil
}

// bindTicTacToeApp binds a generic wrapper to an already deployed contract.
func bindTicTacToeApp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TicTacToeAppABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TicTacToeApp *TicTacToeAppRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TicTacToeApp.Contract.TicTacToeAppCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TicTacToeApp *TicTacToeAppRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TicTacToeApp.Contract.TicTacToeAppTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TicTacToeApp *TicTacToeAppRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TicTacToeApp.Contract.TicTacToeAppTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TicTacToeApp *TicTacToeAppCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TicTacToeApp.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TicTacToeApp *TicTacToeAppTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TicTacToeApp.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TicTacToeApp *TicTacToeAppTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TicTacToeApp.Contract.contract.Transact(opts, method, params...)
}

// ValidTransition is a free data retrieval call binding the contract method 0x0d1feb4f.
//
// Solidity: function validTransition((uint256,uint256,address[],address,bool,bool) params, (bytes32,uint64,(address[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) from, (bytes32,uint64,(address[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) to, uint256 signerIdx) pure returns()
func (_TicTacToeApp *TicTacToeAppCaller) ValidTransition(opts *bind.CallOpts, params ChannelParams, from ChannelState, to ChannelState, signerIdx *big.Int) error {
	var out []interface{}
	err := _TicTacToeApp.contract.Call(opts, &out, "validTransition", params, from, to, signerIdx)

	if err != nil {
		return err
	}

	return err

}

// ValidTransition is a free data retrieval call binding the contract method 0x0d1feb4f.
//
// Solidity: function validTransition((uint256,uint256,address[],address,bool,bool) params, (bytes32,uint64,(address[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) from, (bytes32,uint64,(address[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) to, uint256 signerIdx) pure returns()
func (_TicTacToeApp *TicTacToeAppSession) ValidTransition(params ChannelParams, from ChannelState, to ChannelState, signerIdx *big.Int) error {
	return _TicTacToeApp.Contract.ValidTransition(&_TicTacToeApp.CallOpts, params, from, to, signerIdx)
}

// ValidTransition is a free data retrieval call binding the contract method 0x0d1feb4f.
//
// Solidity: function validTransition((uint256,uint256,address[],address,bool,bool) params, (bytes32,uint64,(address[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) from, (bytes32,uint64,(address[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) to, uint256 signerIdx) pure returns()
func (_TicTacToeApp *TicTacToeAppCallerSession) ValidTransition(params ChannelParams, from ChannelState, to ChannelState, signerIdx *big.Int) error {
	return _TicTacToeApp.Contract.ValidTransition(&_TicTacToeApp.CallOpts, params, from, to, signerIdx)
}
