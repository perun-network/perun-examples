# payment-channel

## Setup
Install dependencies:
```sh
cd multiledger-channel/chain1
npm install
```

Install offckb using:
```sh
npm install -g @offckb/cli
```
and make all scripts executable:
```sh
cd devnet
chmod +x ./setup-devnet.sh
chmod +x ./print_accounts.sh
chmod +x ./fund_omni_accounts.sh
chmod +x ./deploy_contracts.sh
chmod +x ./sudt_helper.sh
cd ..
```
Initialize the submodule.
```sh
git submodule update --init --recursive
```

## Run the Example
Start the local CKB devnet:
```sh
cd devnet
make dev
```

Start the Ethereum local node:
```sh
npx hardhat node --port 8545
```

Then run
```sh
go run .
```