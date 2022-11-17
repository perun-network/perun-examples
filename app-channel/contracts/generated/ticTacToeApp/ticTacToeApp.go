// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ticTacToeApp

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

// TicTacToeAppMetaData contains all meta data concerning the TicTacToeApp contract.
var TicTacToeAppMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"challengeDuration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"app\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"ledgerChannel\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"virtualChannel\",\"type\":\"bool\"}],\"internalType\":\"structChannel.Params\",\"name\":\"params\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"channelID\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"address[]\",\"name\":\"assets\",\"type\":\"address[]\"},{\"internalType\":\"uint256[][]\",\"name\":\"balances\",\"type\":\"uint256[][]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ID\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"balances\",\"type\":\"uint256[]\"},{\"internalType\":\"uint16[]\",\"name\":\"indexMap\",\"type\":\"uint16[]\"}],\"internalType\":\"structChannel.SubAlloc[]\",\"name\":\"locked\",\"type\":\"tuple[]\"}],\"internalType\":\"structChannel.Allocation\",\"name\":\"outcome\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structChannel.State\",\"name\":\"from\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"channelID\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"address[]\",\"name\":\"assets\",\"type\":\"address[]\"},{\"internalType\":\"uint256[][]\",\"name\":\"balances\",\"type\":\"uint256[][]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ID\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"balances\",\"type\":\"uint256[]\"},{\"internalType\":\"uint16[]\",\"name\":\"indexMap\",\"type\":\"uint16[]\"}],\"internalType\":\"structChannel.SubAlloc[]\",\"name\":\"locked\",\"type\":\"tuple[]\"}],\"internalType\":\"structChannel.Allocation\",\"name\":\"outcome\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structChannel.State\",\"name\":\"to\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"signerIdx\",\"type\":\"uint256\"}],\"name\":\"validTransition\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506111c1806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c80630d1feb4f14610030575b600080fd5b61004361003e366004610d36565b610045565b005b60026100546040860186610f98565b90501461007c5760405162461bcd60e51b815260040161007390610f0c565b60405180910390fd5b600061008b6060850185610fdf565b60008161009457fe5b919091013560f81c9150600a90506100af6060850185610fdf565b9050146100ce5760405162461bcd60e51b815260040161007390610f3c565b818160ff16146100f05760405162461bcd60e51b815260040161007390610e40565b6100fd6060840184610fdf565b60008161010657fe5b9091013560f81c6001838101161490506101325760405162461bcd60e51b815260040161007390610e8d565b600060015b600a81101561024f57600261014f6060870187610fdf565b8381811061015957fe5b9050013560f81c60f81b60f81c60ff1611156101875760405162461bcd60e51b815260040161007390610dee565b6101946060870187610fdf565b8281811061019e57fe5b909101356001600160f81b03191690506101bb6060870187610fdf565b838181106101c557fe5b9050013560f81c60f81b6001600160f81b031916146102475760006101ed6060880188610fdf565b838181106101f757fe5b9050013560f81c60f81b60f81c60ff16146102245760405162461bcd60e51b815260040161007390610e6a565b81156102425760405162461bcd60e51b815260040161007390610dc9565b600191505b600101610137565b506000808061029e6102646060890189610fdf565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506105bc92505050565b919450925090508215156102b860a0890160808a01610d0f565b1515146102d75760405162461bcd60e51b815260040161007390610eb1565b6103716102e76040890189611024565b6102f19080610f98565b808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152506103309250505060408b018b611024565b61033a9080610f98565b808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152506107a492505050565b6103c16103816040890189611024565b61038f906040810190610f98565b610398916110d3565b6103a560408b018b611024565b6103b3906040810190610f98565b6103bc916110d3565b61089f565b60006103d060408a018a611024565b6103de906020810190610f98565b6103e791611085565b9050821561058357805160018390039067ffffffffffffffff8111801561040d57600080fd5b5060405190808252806020026020018201604052801561044157816020015b606081526020019060019003908161042c5790505b50915060005b82518110156105805760408051600280825260608201835290916020830190803683370190505083828151811061047a57fe5b602090810291909101015261049260408c018c611024565b6104a0906020810190610f98565b828181106104aa57fe5b90506020028101906104bc9190610f98565b60018181106104c757fe5b905060200201358b80604001906104de9190611024565b6104ec906020810190610f98565b838181106104f657fe5b90506020028101906105089190610f98565b600081811061051357fe5b905060200201350183828151811061052757fe5b60200260200101518560ff168151811061053d57fe5b602002602001018181525050600083828151811061055757fe5b60200260200101518360ff168151811061056d57fe5b6020908102919091010152600101610447565b50505b6105b061059360408a018a611024565b6105a1906020810190610f98565b6105aa91611085565b82610904565b50505050505050505050565b60408051610160810182526000610100820181815260016101208401819052600261014085018190529184528451606081810187526003808352600460208481018290526005858b01819052818a019590955289518085018b52600680825260078284018190526008838e018190528c8e01939093528c518088018e528b815280850196909652858d01829052868c01959095528b518087018d52978852878301849052878c019490945260808a019690965289518085018b5287815280820195909552848a0186905260a089019490945288518084018a52878152808501829052808a019590955260c088019490945287519182018852938152908101919091529384015260e082019290925281908190815b600881101561074a576000806106f6888585600881106106ec57fe5b6020020151610969565b9150915081156107405760ff8116600114156107205760018060009650965096505050505061079d565b60ff8116600214156107405760018060019650965096505050505061079d565b50506001016106d0565b5060005b855181101561078f57600060ff1686828151811061076857fe5b016020015160f81c14610787576000806000945094509450505061079d565b60010161074e565b506001600080935093509350505b9193909250565b80518251146107fa576040805162461bcd60e51b815260206004820152601960248201527f616464726573735b5d3a20756e657175616c206c656e67746800000000000000604482015290519081900360640190fd5b60005b825181101561089a5781818151811061081257fe5b60200260200101516001600160a01b031683828151811061082f57fe5b60200260200101516001600160a01b031614610892576040805162461bcd60e51b815260206004820152601760248201527f616464726573735b5d3a20756e657175616c206974656d000000000000000000604482015290519081900360640190fd5b6001016107fd565b505050565b80518251146108c05760405162461bcd60e51b815260040161007390610f61565b60005b825181101561089a576108fc8382815181106108db57fe5b60200260200101518383815181106108ef57fe5b6020026020010151610a11565b6001016108c3565b80518251146109255760405162461bcd60e51b815260040161007390610ed5565b60005b825181101561089a5761096183828151811061094057fe5b602002602001015183838151811061095457fe5b6020026020010151610a5a565b600101610928565b60008080848482602002015160016000010160ff168151811061098857fe5b01602001516001600160f81b031916905060015b60038110156109ff576001600160f81b03198216868683600381106109bd57fe5b602002015160016000010160ff16815181106109d557fe5b01602001516001600160f81b031916146109f757600080935093505050610a0a565b60010161099c565b506001925060f81c90505b9250929050565b8051825114610a325760405162461bcd60e51b815260040161007390610e12565b610a4482602001518260200151610a5a565b610a5682604001518260400151610b3e565b5050565b8051825114610ab0576040805162461bcd60e51b815260206004820152601960248201527f75696e743235365b5d3a20756e657175616c206c656e67746800000000000000604482015290519081900360640190fd5b60005b825181101561089a57818181518110610ac857fe5b6020026020010151838281518110610adc57fe5b602002602001015114610b36576040805162461bcd60e51b815260206004820152601760248201527f75696e743235365b5d3a20756e657175616c206974656d000000000000000000604482015290519081900360640190fd5b600101610ab3565b8051825114610b94576040805162461bcd60e51b815260206004820152601860248201527f75696e7431365b5d3a20756e657175616c206c656e6774680000000000000000604482015290519081900360640190fd5b60005b825181101561089a57818181518110610bac57fe5b602002602001015161ffff16838281518110610bc457fe5b602002602001015161ffff1614610c1b576040805162461bcd60e51b815260206004820152601660248201527575696e7431365b5d3a20756e657175616c206974656d60501b604482015290519081900360640190fd5b600101610b97565b600082601f830112610c33578081fd5b81356020610c48610c4383611067565b611043565b8281528181019085830183850287018401881015610c64578586fd5b855b85811015610c9157813561ffff81168114610c7f578788fd5b84529284019290840190600101610c66565b5090979650505050505050565b600082601f830112610cae578081fd5b81356020610cbe610c4383611067565b8281528181019085830183850287018401881015610cda578586fd5b855b85811015610c9157813584529284019290840190600101610cdc565b600060a08284031215610d09578081fd5b50919050565b600060208284031215610d20578081fd5b81358015158114610d2f578182fd5b9392505050565b60008060008060808587031215610d4b578283fd5b843567ffffffffffffffff80821115610d62578485fd5b9086019060c08289031215610d75578485fd5b90945060208601359080821115610d8a578485fd5b610d9688838901610cf8565b94506040870135915080821115610dab578384fd5b50610db887828801610cf8565b949793965093946060013593505050565b6020808252600b908201526a74776f20616374696f6e7360a81b604082015260600190565b6020808252600a9082015269677269642076616c756560b01b604082015260600190565b60208082526014908201527314dd58905b1b1bd8ce881d5b995c5d585b08125160621b604082015260600190565b60208082526010908201526f30b1ba37b9103737ba1039b4b3b732b960811b604082015260600190565b6020808252600990820152686f766572777269746560b81b604082015260600190565b6020808252600a90820152693732bc3a1030b1ba37b960b11b604082015260600190565b6020808252600a908201526966696e616c20666c616760b01b604082015260600190565b6020808252601b908201527f75696e743235365b5d5b5d3a20756e657175616c206c656e6774680000000000604082015260600190565b6020808252601690820152756e756d626572206f66207061727469636970616e747360501b604082015260600190565b6020808252600b908201526a0c8c2e8c240d8cadccee8d60ab1b604082015260600190565b6020808252601a908201527f537562416c6c6f635b5d3a20756e657175616c206c656e677468000000000000604082015260600190565b6000808335601e19843603018112610fae578283fd5b83018035915067ffffffffffffffff821115610fc8578283fd5b6020908101925081023603821315610a0a57600080fd5b6000808335601e19843603018112610ff5578283fd5b83018035915067ffffffffffffffff82111561100f578283fd5b602001915036819003821315610a0a57600080fd5b60008235605e19833603018112611039578182fd5b9190910192915050565b60405181810167ffffffffffffffff8111828210171561105f57fe5b604052919050565b600067ffffffffffffffff82111561107b57fe5b5060209081020190565b6000611093610c4384611067565b8381526020808201919084845b878110156110c7576110b53683358901610c9e565b855293820193908201906001016110a0565b50919695505050505050565b60006110e1610c4384611067565b8381526020808201919084845b878110156110c75781358701606080823603121561110a578788fd5b604080519182019167ffffffffffffffff808411828510171561112957fe5b92825283358152868401359280841115611141578a8bfd5b61114d36858701610c9e565b8883015282850135935080841115611163578a8bfd5b5061117036848601610c23565b918101919091528752505093820193908201906001016110ee56fea264697066735822122044800f789b90fd7beb3d3050eafee3ae00bffb9f3b9087b671f24666c0f952bb64736f6c63430007060033",
}

// TicTacToeAppABI is the input ABI used to generate the binding from.
// Deprecated: Use TicTacToeAppMetaData.ABI instead.
var TicTacToeAppABI = TicTacToeAppMetaData.ABI

// TicTacToeAppBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TicTacToeAppMetaData.Bin instead.
var TicTacToeAppBin = TicTacToeAppMetaData.Bin

// DeployTicTacToeApp deploys a new Ethereum contract, binding an instance of TicTacToeApp to it.
func DeployTicTacToeApp(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TicTacToeApp, error) {
	parsed, err := TicTacToeAppMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TicTacToeAppBin), backend)
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
