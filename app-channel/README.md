# app-channel
## Requirements
- Node.js 20.x or higher
- npm or yarn
- Hardhat 2.22.x

## Setup
Install dependencies:
```sh
npm install
```

or 
```sh
yarn add --dev hardhat
```


## Running the Ethereum Testnet

### Using Hardhat (Recommended)
Run the Hardhat node with pre-configured accounts:
```sh
npx hardhat node
```

This will start a local Ethereum node on http://127.0.0.1:8545 with the following accounts:

- Deployer: `0x79ea8f62d97bc0591a4224c1725fca6b00de5b2cea286fe2e0bb35c5e76be46e` (10 ETH)

- Alice: `0x1af2e950272dd403de7a5760d41c6e44d92b6d02797e51810795ff03cc2cda4f` (10 ETH)

- Bob: `0xf63d7d8e930bccd74e93cf5662fde2c28fd8be95edb70c73f1bdd863d07f412e` (10 ETH)

The node mines a new block every 5 seconds and runs on chain ID `1337`.



### Using Docker
Build and run the Hardhat node in Docker:
```sh
docker build -t app-channel-node .
docker run -p 8545:8545 app-channel-node
```



## Running the Application
Then run
```
go run .
```