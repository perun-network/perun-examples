# multi-ledger-channel

## Requirements
- Node.js 20.x or higher
- npm
- Hardhat 2.22.x

## Setup
Install dependencies:
```sh
cd multiledger-channel/chain1
npm install
```

## Running the Ethereum Testnet

### Using Hardhat (Recommended)

On one terminal:
```sh
cd multiledger-channel/chain1
npx hardhat node --port 8545
```

On another terminal:
```sh
cd multiledger-channel/chain1
npx hardhat node --port 8546
```

### Using Docker
Build and run both Hardhat nodes:
```sh
docker build -t hardhat-chain1 -f chain1/Dockerfile chain1
docker build -t hardhat-chain2 -f chain2/Dockerfile chain2

docker run -d --name chain1 -p 8545:8545 hardhat-chain1
docker run -d --name chain2 -p 8546:8545 hardhat-chain2
```

To stop the chains:
```sh
docker stop chain1 chain2
docker rm chain1 chain2
```

## Running the Application
With both chains running, execute:
```
go run .
```