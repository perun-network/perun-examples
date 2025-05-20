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
	Backends []*big.Int
	Balances [][]*big.Int
	Locked   []ChannelSubAlloc
}

// ChannelAsset is an auto generated low-level Go binding around an user-defined struct.
type ChannelAsset struct {
	ChainID   *big.Int
	EthHolder common.Address
	CcHolder  []byte
}

// ChannelParams is an auto generated low-level Go binding around an user-defined struct.
type ChannelParams struct {
	ChallengeDuration *big.Int
	Nonce             *big.Int
	Participants      []ChannelParticipant
	App               common.Address
	LedgerChannel     bool
	VirtualChannel    bool
}

// ChannelParticipant is an auto generated low-level Go binding around an user-defined struct.
type ChannelParticipant struct {
	EthAddress common.Address
	CcAddress  []byte
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
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"challengeDuration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"ethAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"ccAddress\",\"type\":\"bytes\"}],\"internalType\":\"structChannel.Participant[]\",\"name\":\"participants\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"app\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"ledgerChannel\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"virtualChannel\",\"type\":\"bool\"}],\"internalType\":\"structChannel.Params\",\"name\":\"params\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"channelID\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"ethHolder\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"ccHolder\",\"type\":\"bytes\"}],\"internalType\":\"structChannel.Asset[]\",\"name\":\"assets\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"backends\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[][]\",\"name\":\"balances\",\"type\":\"uint256[][]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ID\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"balances\",\"type\":\"uint256[]\"},{\"internalType\":\"uint16[]\",\"name\":\"indexMap\",\"type\":\"uint16[]\"}],\"internalType\":\"structChannel.SubAlloc[]\",\"name\":\"locked\",\"type\":\"tuple[]\"}],\"internalType\":\"structChannel.Allocation\",\"name\":\"outcome\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structChannel.State\",\"name\":\"from\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"channelID\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"ethHolder\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"ccHolder\",\"type\":\"bytes\"}],\"internalType\":\"structChannel.Asset[]\",\"name\":\"assets\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"backends\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[][]\",\"name\":\"balances\",\"type\":\"uint256[][]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ID\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"balances\",\"type\":\"uint256[]\"},{\"internalType\":\"uint16[]\",\"name\":\"indexMap\",\"type\":\"uint16[]\"}],\"internalType\":\"structChannel.SubAlloc[]\",\"name\":\"locked\",\"type\":\"tuple[]\"}],\"internalType\":\"structChannel.Allocation\",\"name\":\"outcome\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structChannel.State\",\"name\":\"to\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"signerIdx\",\"type\":\"uint256\"}],\"name\":\"validTransition\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b506114ec8061001c5f395ff3fe608060405234801561000f575f5ffd5b5060043610610029575f3560e01c8063f7530b411461002d575b5f5ffd5b61004061003b366004610f01565b610042565b005b60026100516040860186610f98565b90501461009e5760405162461bcd60e51b81526020600482015260166024820152756e756d626572206f66207061727469636970616e747360501b60448201526064015b60405180910390fd5b5f6100ac6060850185610fdd565b5f816100ba576100ba61101f565b919091013560f81c9150600990506100d360015f611047565b6100dd9190611047565b60ff166100ed6060850185610fdd565b90501461012a5760405162461bcd60e51b815260206004820152600b60248201526a0c8c2e8c240d8cadccee8d60ab1b6044820152606401610095565b818160ff161461016f5760405162461bcd60e51b815260206004820152601060248201526f30b1ba37b9103737ba1039b4b3b732b960811b6044820152606401610095565b61017c6060840184610fdd565b5f8161018a5761018a61101f565b919091013560f81c905060026101a1836001611047565b6101ab9190611066565b60ff16146101e85760405162461bcd60e51b815260206004820152600a6024820152693732bc3a1030b1ba37b960b11b6044820152606401610095565b5f806101f5600182611047565b60ff1690505b600961020860015f611047565b6102129190611047565b60ff1681101561039957600261022b6060870187610fdd565b8381811061023b5761023b61101f565b9050013560f81c60f81b60f81c60ff1611156102865760405162461bcd60e51b815260206004820152600a602482015269677269642076616c756560b01b6044820152606401610095565b6102936060870187610fdd565b828181106102a3576102a361101f565b909101356001600160f81b03191690506102c06060870187610fdd565b838181106102d0576102d061101f565b9050013560f81c60f81b6001600160f81b03191614610391575f6102f76060880188610fdd565b838181106103075761030761101f565b9050013560f81c60f81b60f81c60ff16146103505760405162461bcd60e51b81526020600482015260096024820152686f766572777269746560b81b6044820152606401610095565b811561038c5760405162461bcd60e51b815260206004820152600b60248201526a74776f20616374696f6e7360a81b6044820152606401610095565b600191505b6001016101fb565b505f80806103e66103ad6060890189610fdd565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284375f9201919091525061071292505050565b9194509250905082151561040060a0890160808a01611093565b15151461043c5760405162461bcd60e51b815260206004820152600a60248201526966696e616c20666c616760b01b6044820152606401610095565b61048461044c60408901896110b9565b6104569080610f98565b61045f91611165565b61046c60408b018b6110b9565b6104769080610f98565b61047f91611165565b6108fd565b6104d461049460408901896110b9565b6104a2906060810190610f98565b6104ab916112f8565b6104b860408b018b6110b9565b6104c6906060810190610f98565b6104cf916112f8565b6109a2565b5f6104e260408a018a6110b9565b6104f0906040810190610f98565b6104f991611428565b905082156106d9575f61050d83600161148a565b905081516001600160401b03811115610528576105286110d7565b60405190808252806020026020018201604052801561055b57816020015b60608152602001906001900390816105465790505b5091505f5b82518110156106d6576040805160028082526060820183529091602083019080368337019050508382815181106105995761059961101f565b60209081029190910101526105b160408c018c6110b9565b6105bf906040810190610f98565b828181106105cf576105cf61101f565b90506020028101906105e19190610f98565b60018181106105f2576105f261101f565b905060200201358b806040019061060991906110b9565b610617906040810190610f98565b838181106106275761062761101f565b90506020028101906106399190610f98565b5f8181106106495761064961101f565b9050602002013561065a91906114a3565b83828151811061066c5761066c61101f565b60200260200101518560ff16815181106106885761068861101f565b6020026020010181815250505f8382815181106106a7576106a761101f565b60200260200101518360ff16815181106106c3576106c361101f565b6020908102919091010152600101610560565b50505b6107066106e960408a018a6110b9565b6106f7906040810190610f98565b61070091611428565b82610a42565b50505050505050505050565b60408051610160810182525f610100820181815260016101208401819052600261014085018190529184528451606081810187526003808352600460208481018290526005858b01819052818a019590955289518085018b52600680825260078284018190526008838e018190528c8e01939093528c518088018e528b815280850196909652858d01829052868c01959095528b518087018d52978852878301849052878c019490945260808a019690965289518085018b5287815280820195909552848a0186905260a089019490945288518084018a52878152808501829052808a019590955260c088019490945287519182018852938152908101919091529384015260e082019290925281908190815b60088110156108a2575f5f610850888585600881106108465761084661101f565b6020020151610ae2565b915091508115610898575f1960ff821601610878576001805f965096509650505050506108f6565b60011960ff821601610898576001806001965096509650505050506108f6565b5050600101610825565b505f5b85518110156108e9575f60ff168682815181106108c4576108c461101f565b016020015160f81c146108e1575f5f5f94509450945050506108f6565b6001016108a5565b5060015f5f935093509350505b9193909250565b805182511461094e5760405162461bcd60e51b815260206004820152601760248201527f41737365745b5d3a20756e657175616c206c656e6774680000000000000000006044820152606401610095565b5f5b825181101561099d5761099583828151811061096e5761096e61101f565b60200260200101518383815181106109885761098861101f565b6020026020010151610bb8565b600101610950565b505050565b80518251146109f35760405162461bcd60e51b815260206004820152601a60248201527f537562416c6c6f635b5d3a20756e657175616c206c656e6774680000000000006044820152606401610095565b5f5b825181101561099d57610a3a838281518110610a1357610a1361101f565b6020026020010151838381518110610a2d57610a2d61101f565b6020026020010151610cb4565b6001016109f5565b8051825114610a935760405162461bcd60e51b815260206004820152601b60248201527f75696e743235365b5d5b5d3a20756e657175616c206c656e67746800000000006044820152606401610095565b5f5b825181101561099d57610ada838281518110610ab357610ab361101f565b6020026020010151838381518110610acd57610acd61101f565b6020026020010151610d20565b600101610a95565b5f80808484826020020151610af860015f611047565b610b029190611047565b60ff1681518110610b1557610b1561101f565b01602001516001600160f81b031916905060015b6003811015610ba6576001600160f81b0319821686868360038110610b5057610b5061101f565b6020020151610b6060015f611047565b610b6a9190611047565b60ff1681518110610b7d57610b7d61101f565b01602001516001600160f81b03191614610b9e575f5f935093505050610bb1565b600101610b29565b506001925060f81c90505b9250929050565b8051825114610bfb5760405162461bcd60e51b815260206004820152600f60248201526e1d5b995c5d585b0818da185a5b9251608a1b6044820152606401610095565b80602001516001600160a01b031682602001516001600160a01b031614610c585760405162461bcd60e51b81526020600482015260116024820152703ab732b8bab0b61032ba342437b63232b960791b6044820152606401610095565b80604001518051906020012082604001518051906020012014610cb05760405162461bcd60e51b815260206004820152601060248201526f3ab732b8bab0b61031b1a437b63232b960811b6044820152606401610095565b5050565b8051825114610cfc5760405162461bcd60e51b815260206004820152601460248201527314dd58905b1b1bd8ce881d5b995c5d585b08125160621b6044820152606401610095565b610d0e82602001518260200151610d20565b610cb082604001518260400151610e05565b8051825114610d715760405162461bcd60e51b815260206004820152601960248201527f75696e743235365b5d3a20756e657175616c206c656e677468000000000000006044820152606401610095565b5f5b825181101561099d57818181518110610d8e57610d8e61101f565b6020026020010151838281518110610da857610da861101f565b602002602001015114610dfd5760405162461bcd60e51b815260206004820152601760248201527f75696e743235365b5d3a20756e657175616c206974656d0000000000000000006044820152606401610095565b600101610d73565b8051825114610e565760405162461bcd60e51b815260206004820152601860248201527f75696e7431365b5d3a20756e657175616c206c656e67746800000000000000006044820152606401610095565b5f5b825181101561099d57818181518110610e7357610e7361101f565b602002602001015161ffff16838281518110610e9157610e9161101f565b602002602001015161ffff1614610ee35760405162461bcd60e51b815260206004820152601660248201527575696e7431365b5d3a20756e657175616c206974656d60501b6044820152606401610095565b600101610e58565b5f60a08284031215610efb575f5ffd5b50919050565b5f5f5f5f60808587031215610f14575f5ffd5b84356001600160401b03811115610f29575f5ffd5b850160c08188031215610f3a575f5ffd5b935060208501356001600160401b03811115610f54575f5ffd5b610f6087828801610eeb565b93505060408501356001600160401b03811115610f7b575f5ffd5b610f8787828801610eeb565b949793965093946060013593505050565b5f5f8335601e19843603018112610fad575f5ffd5b8301803591506001600160401b03821115610fc6575f5ffd5b6020019150600581901b3603821315610bb1575f5ffd5b5f5f8335601e19843603018112610ff2575f5ffd5b8301803591506001600160401b0382111561100b575f5ffd5b602001915036819003821315610bb1575f5ffd5b634e487b7160e01b5f52603260045260245ffd5b634e487b7160e01b5f52601160045260245ffd5b60ff818116838216019081111561106057611060611033565b92915050565b5f60ff83168061108457634e487b7160e01b5f52601260045260245ffd5b8060ff84160691505092915050565b5f602082840312156110a3575f5ffd5b813580151581146110b2575f5ffd5b9392505050565b5f8235607e198336030181126110cd575f5ffd5b9190910192915050565b634e487b7160e01b5f52604160045260245ffd5b604051606081016001600160401b038111828210171561110d5761110d6110d7565b60405290565b604051601f8201601f191681016001600160401b038111828210171561113b5761113b6110d7565b604052919050565b5f6001600160401b0382111561115b5761115b6110d7565b5060051b60200190565b5f61117761117284611143565b611113565b8381526020810190600585901b840136811115611192575f5ffd5b845b818110156112885780356001600160401b038111156111b1575f5ffd5b860160603682900312156111c3575f5ffd5b6111cb6110eb565b8135815260208201356001600160a01b03811681146111e8575f5ffd5b602082015260408201356001600160401b03811115611205575f5ffd5b919091019036601f830112611218575f5ffd5b81356001600160401b03811115611231576112316110d7565b611244601f8201601f1916602001611113565b818152366020838601011115611258575f5ffd5b816020850160208301375f6020838301015280604084015250508086525050602084019350602081019050611194565b509095945050505050565b5f82601f8301126112a2575f5ffd5b81356112b061117282611143565b8082825260208201915060208360051b8601019250858311156112d1575f5ffd5b602085015b838110156112ee5780358352602092830192016112d6565b5095945050505050565b5f61130561117284611143565b8381526020810190600585901b840136811115611320575f5ffd5b845b818110156112885780356001600160401b0381111561133f575f5ffd5b86016060368290031215611351575f5ffd5b6113596110eb565b8135815260208201356001600160401b03811115611375575f5ffd5b61138136828501611293565b60208301525060408201356001600160401b0381111561139f575f5ffd5b919091019036601f8301126113b2575f5ffd5b81356113c061117282611143565b8082825260208201915060208360051b8601019250368311156113e1575f5ffd5b6020850194505b8285101561141257843561ffff81168114611401575f5ffd5b8252602094850194909101906113e8565b6040840152505085525060209384019301611322565b5f61143561117284611143565b8381526020810190600585901b840136811115611450575f5ffd5b845b818110156112885780356001600160401b0381111561146f575f5ffd5b61147b36828901611293565b85525060209384019301611452565b60ff828116828216039081111561106057611060611033565b808201808211156110605761106061103356fea2646970667358221220c9c8e469f26d979598bb9811df8c27c60f981d09c31821da042b28130cb61da964736f6c634300081b0033",
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

// ValidTransition is a free data retrieval call binding the contract method 0xf7530b41.
//
// Solidity: function validTransition((uint256,uint256,(address,bytes)[],address,bool,bool) params, (bytes32,uint64,((uint256,address,bytes)[],uint256[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) from, (bytes32,uint64,((uint256,address,bytes)[],uint256[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) to, uint256 signerIdx) pure returns()
func (_TicTacToeApp *TicTacToeAppCaller) ValidTransition(opts *bind.CallOpts, params ChannelParams, from ChannelState, to ChannelState, signerIdx *big.Int) error {
	var out []interface{}
	err := _TicTacToeApp.contract.Call(opts, &out, "validTransition", params, from, to, signerIdx)

	if err != nil {
		return err
	}

	return err

}

// ValidTransition is a free data retrieval call binding the contract method 0xf7530b41.
//
// Solidity: function validTransition((uint256,uint256,(address,bytes)[],address,bool,bool) params, (bytes32,uint64,((uint256,address,bytes)[],uint256[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) from, (bytes32,uint64,((uint256,address,bytes)[],uint256[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) to, uint256 signerIdx) pure returns()
func (_TicTacToeApp *TicTacToeAppSession) ValidTransition(params ChannelParams, from ChannelState, to ChannelState, signerIdx *big.Int) error {
	return _TicTacToeApp.Contract.ValidTransition(&_TicTacToeApp.CallOpts, params, from, to, signerIdx)
}

// ValidTransition is a free data retrieval call binding the contract method 0xf7530b41.
//
// Solidity: function validTransition((uint256,uint256,(address,bytes)[],address,bool,bool) params, (bytes32,uint64,((uint256,address,bytes)[],uint256[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) from, (bytes32,uint64,((uint256,address,bytes)[],uint256[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) to, uint256 signerIdx) pure returns()
func (_TicTacToeApp *TicTacToeAppCallerSession) ValidTransition(params ChannelParams, from ChannelState, to ChannelState, signerIdx *big.Int) error {
	return _TicTacToeApp.Contract.ValidTransition(&_TicTacToeApp.CallOpts, params, from, to, signerIdx)
}
