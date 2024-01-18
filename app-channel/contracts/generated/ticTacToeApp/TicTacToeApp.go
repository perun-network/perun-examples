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
	_ = abi.ConvertType
)

// ChannelAllocation is an auto generated low-level Go binding around an user-defined struct.
type ChannelAllocation struct {
	Assets   []ChannelAsset
	Balances [][]*big.Int
	Locked   []ChannelSubAlloc
}

// ChannelAsset is an auto generated low-level Go binding around an user-defined struct.
type ChannelAsset struct {
	ChainID *big.Int
	Holder  common.Address
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
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"challengeDuration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"app\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"ledgerChannel\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"virtualChannel\",\"type\":\"bool\"}],\"internalType\":\"structChannel.Params\",\"name\":\"params\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"channelID\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"holder\",\"type\":\"address\"}],\"internalType\":\"structChannel.Asset[]\",\"name\":\"assets\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[][]\",\"name\":\"balances\",\"type\":\"uint256[][]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ID\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"balances\",\"type\":\"uint256[]\"},{\"internalType\":\"uint16[]\",\"name\":\"indexMap\",\"type\":\"uint16[]\"}],\"internalType\":\"structChannel.SubAlloc[]\",\"name\":\"locked\",\"type\":\"tuple[]\"}],\"internalType\":\"structChannel.Allocation\",\"name\":\"outcome\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structChannel.State\",\"name\":\"from\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"channelID\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"holder\",\"type\":\"address\"}],\"internalType\":\"structChannel.Asset[]\",\"name\":\"assets\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[][]\",\"name\":\"balances\",\"type\":\"uint256[][]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ID\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"balances\",\"type\":\"uint256[]\"},{\"internalType\":\"uint16[]\",\"name\":\"indexMap\",\"type\":\"uint16[]\"}],\"internalType\":\"structChannel.SubAlloc[]\",\"name\":\"locked\",\"type\":\"tuple[]\"}],\"internalType\":\"structChannel.Allocation\",\"name\":\"outcome\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structChannel.State\",\"name\":\"to\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"signerIdx\",\"type\":\"uint256\"}],\"name\":\"validTransition\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561000f575f80fd5b506114b98061001d5f395ff3fe608060405234801561000f575f80fd5b5060043610610029575f3560e01c80636d7eba0d1461002d575b5f80fd5b61004061003b366004610f46565b610042565b005b60026100516040860186610fd7565b90501461009e5760405162461bcd60e51b81526020600482015260166024820152756e756d626572206f66207061727469636970616e747360501b60448201526064015b60405180910390fd5b5f6100ac606085018561101d565b5f816100ba576100ba611060565b919091013560f81c9150600990506100d360015f611088565b6100dd9190611088565b60ff166100ed606085018561101d565b90501461012a5760405162461bcd60e51b815260206004820152600b60248201526a0c8c2e8c240d8cadccee8d60ab1b6044820152606401610095565b818160ff161461016f5760405162461bcd60e51b815260206004820152601060248201526f30b1ba37b9103737ba1039b4b3b732b960811b6044820152606401610095565b61017c606084018461101d565b5f8161018a5761018a611060565b919091013560f81c905060026101a1836001611088565b6101ab91906110a7565b60ff16146101e85760405162461bcd60e51b815260206004820152600a6024820152693732bc3a1030b1ba37b960b11b6044820152606401610095565b5f806101f5600182611088565b60ff1690505b600961020860015f611088565b6102129190611088565b60ff1681101561039957600261022b606087018761101d565b8381811061023b5761023b611060565b9050013560f81c60f81b60f81c60ff1611156102865760405162461bcd60e51b815260206004820152600a602482015269677269642076616c756560b01b6044820152606401610095565b610293606087018761101d565b828181106102a3576102a3611060565b909101356001600160f81b03191690506102c0606087018761101d565b838181106102d0576102d0611060565b9050013560f81c60f81b6001600160f81b03191614610391575f6102f7606088018861101d565b8381811061030757610307611060565b9050013560f81c60f81b60f81c60ff16146103505760405162461bcd60e51b81526020600482015260096024820152686f766572777269746560b81b6044820152606401610095565b811561038c5760405162461bcd60e51b815260206004820152600b60248201526a74776f20616374696f6e7360a81b6044820152606401610095565b600191505b6001016101fb565b505f80806103e66103ad606089018961101d565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284375f920191909152506107a492505050565b9194509250905082151561040060a0890160808a016110d4565b15151461043c5760405162461bcd60e51b815260206004820152600a60248201526966696e616c20666c616760b01b6044820152606401610095565b61051561044c60408901896110fa565b6104569080611118565b808060200260200160405190810160405280939291908181526020015f905b828210156104a157610492604083028601368190038101906111cc565b81526020019060010190610475565b50505050508980604001906104b691906110fa565b6104c09080611118565b808060200260200160405190810160405280939291908181526020015f905b8282101561050b576104fc604083028601368190038101906111cc565b815260200190600101906104df565b505050505061098f565b61056561052560408901896110fa565b610533906040810190610fd7565b61053c916112ba565b61054960408b018b6110fa565b610557906040810190610fd7565b610560916112ba565b610a34565b5f61057360408a018a6110fa565b610581906020810190610fd7565b61058a916113e7565b9050821561076b575f61059e836001611457565b9050815167ffffffffffffffff8111156105ba576105ba61115e565b6040519080825280602002602001820160405280156105ed57816020015b60608152602001906001900390816105d85790505b5091505f5b82518110156107685760408051600280825260608201835290916020830190803683370190505083828151811061062b5761062b611060565b602090810291909101015261064360408c018c6110fa565b610651906020810190610fd7565b8281811061066157610661611060565b90506020028101906106739190610fd7565b600181811061068457610684611060565b905060200201358b806040019061069b91906110fa565b6106a9906020810190610fd7565b838181106106b9576106b9611060565b90506020028101906106cb9190610fd7565b5f8181106106db576106db611060565b905060200201356106ec9190611470565b8382815181106106fe576106fe611060565b60200260200101518560ff168151811061071a5761071a611060565b6020026020010181815250505f83828151811061073957610739611060565b60200260200101518360ff168151811061075557610755611060565b60209081029190910101526001016105f2565b50505b61079861077b60408a018a6110fa565b610789906020810190610fd7565b610792916113e7565b82610ad4565b50505050505050505050565b60408051610160810182525f610100820181815260016101208401819052600261014085018190529184528451606081810187526003808352600460208481018290526005858b01819052818a019590955289518085018b52600680825260078284018190526008838e018190528c8e01939093528c518088018e528b815280850196909652858d01829052868c01959095528b518087018d52978852878301849052878c019490945260808a019690965289518085018b5287815280820195909552848a0186905260a089019490945288518084018a52878152808501829052808a019590955260c088019490945287519182018852938152908101919091529384015260e082019290925281908190815b6008811015610934575f806108e2888585600881106108d8576108d8611060565b6020020151610b74565b91509150811561092a575f1960ff82160161090a576001805f96509650965050505050610988565b60011960ff82160161092a57600180600196509650965050505050610988565b50506001016108b7565b505f5b855181101561097b575f60ff1686828151811061095657610956611060565b016020015160f81c14610973575f805f9450945094505050610988565b600101610937565b5060015f80935093509350505b9193909250565b80518251146109e05760405162461bcd60e51b815260206004820152601760248201527f41737365745b5d3a20756e657175616c206c656e6774680000000000000000006044820152606401610095565b5f5b8251811015610a2f57610a27838281518110610a0057610a00611060565b6020026020010151838381518110610a1a57610a1a611060565b6020026020010151610c4a565b6001016109e2565b505050565b8051825114610a855760405162461bcd60e51b815260206004820152601a60248201527f537562416c6c6f635b5d3a20756e657175616c206c656e6774680000000000006044820152606401610095565b5f5b8251811015610a2f57610acc838281518110610aa557610aa5611060565b6020026020010151838381518110610abf57610abf611060565b6020026020010151610cf9565b600101610a87565b8051825114610b255760405162461bcd60e51b815260206004820152601b60248201527f75696e743235365b5d5b5d3a20756e657175616c206c656e67746800000000006044820152606401610095565b5f5b8251811015610a2f57610b6c838281518110610b4557610b45611060565b6020026020010151838381518110610b5f57610b5f611060565b6020026020010151610d65565b600101610b27565b5f80808484826020020151610b8a60015f611088565b610b949190611088565b60ff1681518110610ba757610ba7611060565b01602001516001600160f81b031916905060015b6003811015610c38576001600160f81b0319821686868360038110610be257610be2611060565b6020020151610bf260015f611088565b610bfc9190611088565b60ff1681518110610c0f57610c0f611060565b01602001516001600160f81b03191614610c30575f80935093505050610c43565b600101610bbb565b506001925060f81c90505b9250929050565b8051825114610c945760405162461bcd60e51b8152602060048201526016602482015275105cdcd95d0e881d5b995c5d585b0818da185a5b925160521b6044820152606401610095565b80602001516001600160a01b031682602001516001600160a01b031614610cf55760405162461bcd60e51b815260206004820152601560248201527420b9b9b2ba1d103ab732b8bab0b6103437b63232b960591b6044820152606401610095565b5050565b8051825114610d415760405162461bcd60e51b815260206004820152601460248201527314dd58905b1b1bd8ce881d5b995c5d585b08125160621b6044820152606401610095565b610d5382602001518260200151610d65565b610cf582604001518260400151610e4a565b8051825114610db65760405162461bcd60e51b815260206004820152601960248201527f75696e743235365b5d3a20756e657175616c206c656e677468000000000000006044820152606401610095565b5f5b8251811015610a2f57818181518110610dd357610dd3611060565b6020026020010151838281518110610ded57610ded611060565b602002602001015114610e425760405162461bcd60e51b815260206004820152601760248201527f75696e743235365b5d3a20756e657175616c206974656d0000000000000000006044820152606401610095565b600101610db8565b8051825114610e9b5760405162461bcd60e51b815260206004820152601860248201527f75696e7431365b5d3a20756e657175616c206c656e67746800000000000000006044820152606401610095565b5f5b8251811015610a2f57818181518110610eb857610eb8611060565b602002602001015161ffff16838281518110610ed657610ed6611060565b602002602001015161ffff1614610f285760405162461bcd60e51b815260206004820152601660248201527575696e7431365b5d3a20756e657175616c206974656d60501b6044820152606401610095565b600101610e9d565b5f60a08284031215610f40575f80fd5b50919050565b5f805f8060808587031215610f59575f80fd5b843567ffffffffffffffff80821115610f70575f80fd5b9086019060c08289031215610f83575f80fd5b90945060208601359080821115610f98575f80fd5b610fa488838901610f30565b94506040870135915080821115610fb9575f80fd5b50610fc687828801610f30565b949793965093946060013593505050565b5f808335601e19843603018112610fec575f80fd5b83018035915067ffffffffffffffff821115611006575f80fd5b6020019150600581901b3603821315610c43575f80fd5b5f808335601e19843603018112611032575f80fd5b83018035915067ffffffffffffffff82111561104c575f80fd5b602001915036819003821315610c43575f80fd5b634e487b7160e01b5f52603260045260245ffd5b634e487b7160e01b5f52601160045260245ffd5b60ff81811683821601908111156110a1576110a1611074565b92915050565b5f60ff8316806110c557634e487b7160e01b5f52601260045260245ffd5b8060ff84160691505092915050565b5f602082840312156110e4575f80fd5b813580151581146110f3575f80fd5b9392505050565b5f8235605e1983360301811261110e575f80fd5b9190910192915050565b5f808335601e1984360301811261112d575f80fd5b83018035915067ffffffffffffffff821115611147575f80fd5b6020019150600681901b3603821315610c43575f80fd5b634e487b7160e01b5f52604160045260245ffd5b6040516060810167ffffffffffffffff811182821017156111955761119561115e565b60405290565b604051601f8201601f1916810167ffffffffffffffff811182821017156111c4576111c461115e565b604052919050565b5f604082840312156111dc575f80fd5b6040516040810181811067ffffffffffffffff821117156111ff576111ff61115e565b6040528235815260208301356001600160a01b038116811461121f575f80fd5b60208201529392505050565b5f67ffffffffffffffff8211156112445761124461115e565b5060051b60200190565b5f82601f83011261125d575f80fd5b8135602061127261126d8361122b565b61119b565b8083825260208201915060208460051b870101935086841115611293575f80fd5b602086015b848110156112af5780358352918301918301611298565b509695505050505050565b5f6112c761126d8461122b565b83815260208082019190600586811b8601368111156112e4575f80fd5b865b818110156113da57803567ffffffffffffffff80821115611305575f80fd5b818a01915060608236031215611319575f80fd5b611321611172565b823581528683013582811115611335575f80fd5b6113413682860161124e565b888301525060408084013583811115611358575f80fd5b939093019236601f85011261136b575f80fd5b8335925061137b61126d8461122b565b83815292871b84018801928881019036851115611396575f80fd5b948901945b848610156113c357853561ffff811681146113b4575f80fd5b8252948901949089019061139b565b9183019190915250885250509483019483016112e6565b5092979650505050505050565b5f6113f461126d8461122b565b80848252602080830192508560051b850136811115611411575f80fd5b855b8181101561144b57803567ffffffffffffffff811115611431575f80fd5b61143d36828a0161124e565b865250938201938201611413565b50919695505050505050565b60ff82811682821603908111156110a1576110a1611074565b808201808211156110a1576110a161107456fea2646970667358221220ad528d1f14dc0abe757f06b3c986b06bad96d24446e822cc9cb032d08e54b35b64736f6c63430008170033",
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
	parsed, err := TicTacToeAppMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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

// ValidTransition is a free data retrieval call binding the contract method 0x6d7eba0d.
//
// Solidity: function validTransition((uint256,uint256,address[],address,bool,bool) params, (bytes32,uint64,((uint256,address)[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) from, (bytes32,uint64,((uint256,address)[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) to, uint256 signerIdx) pure returns()
func (_TicTacToeApp *TicTacToeAppCaller) ValidTransition(opts *bind.CallOpts, params ChannelParams, from ChannelState, to ChannelState, signerIdx *big.Int) error {
	var out []interface{}
	err := _TicTacToeApp.contract.Call(opts, &out, "validTransition", params, from, to, signerIdx)

	if err != nil {
		return err
	}

	return err

}

// ValidTransition is a free data retrieval call binding the contract method 0x6d7eba0d.
//
// Solidity: function validTransition((uint256,uint256,address[],address,bool,bool) params, (bytes32,uint64,((uint256,address)[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) from, (bytes32,uint64,((uint256,address)[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) to, uint256 signerIdx) pure returns()
func (_TicTacToeApp *TicTacToeAppSession) ValidTransition(params ChannelParams, from ChannelState, to ChannelState, signerIdx *big.Int) error {
	return _TicTacToeApp.Contract.ValidTransition(&_TicTacToeApp.CallOpts, params, from, to, signerIdx)
}

// ValidTransition is a free data retrieval call binding the contract method 0x6d7eba0d.
//
// Solidity: function validTransition((uint256,uint256,address[],address,bool,bool) params, (bytes32,uint64,((uint256,address)[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) from, (bytes32,uint64,((uint256,address)[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) to, uint256 signerIdx) pure returns()
func (_TicTacToeApp *TicTacToeAppCallerSession) ValidTransition(params ChannelParams, from ChannelState, to ChannelState, signerIdx *big.Int) error {
	return _TicTacToeApp.Contract.ValidTransition(&_TicTacToeApp.CallOpts, params, from, to, signerIdx)
}
